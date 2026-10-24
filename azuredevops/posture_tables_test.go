package azuredevops

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

// Response shapes below follow the Azure DevOps REST API reference examples for each endpoint.

func routes(t *testing.T, handlers map[string]func(w http.ResponseWriter, r *http.Request)) (*fakeADO, *restClient) {
	t.Helper()
	f := newFakeADO(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		h, ok := handlers[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h(w, r)
	})
	var slept []time.Duration
	return f, f.client(&slept)
}

func body(s string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) { _, _ = fmt.Fprint(w, s) }
}

func status(code int, msg string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(code)
		_, _ = fmt.Fprintf(w, `{"message":%q}`, msg)
	}
}

func assertPermissionDenied(t *testing.T, err error, scope string) {
	t.Helper()
	if err == nil {
		t.Fatal("no error, so a missing scope would read as an empty result")
	}
	var statusErr *restStatusError
	if !errors.As(err, &statusErr) || !statusErr.PermissionDenied() {
		t.Fatalf("not a permission error: %v", err)
	}
	if !strings.HasPrefix(err.Error(), permissionDeniedPrefix) || !strings.Contains(err.Error(), scope) {
		t.Errorf("error %q should start %q and name %q", err, permissionDeniedPrefix, scope)
	}
}

func TestRefusedRequestsAreClearPermissionErrors(t *testing.T) {
	signIn := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", "https://spsprodcus4.vssps.visualstudio.com/_signin?realm=dev.azure.com")
		w.WriteHeader(http.StatusFound)
	}
	for name, h := range map[string]func(http.ResponseWriter, *http.Request){
		"401":              status(http.StatusUnauthorized, "TF400813: The user is not authorized to access this resource."),
		"403":              status(http.StatusForbidden, "Access denied."),
		"203 sign-in page": func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNonAuthoritativeInfo) },
		"302 to sign-in":   signIn,
	} {
		t.Run(name, func(t *testing.T) {
			_, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){"/myorg/_apis/audit/streams": h})
			err := readAuditStreams(context.Background(), c, func(AuditStream) bool { return true })
			assertPermissionDenied(t, err, "vso.auditlog")
		})
	}
}

func TestOtherFailuresAreNotPermissionErrors(t *testing.T) {
	_, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/_apis/audit/streams": status(http.StatusInternalServerError, "boom"),
	})
	err := readAuditStreams(context.Background(), c, func(AuditStream) bool { return true })
	if err == nil || strings.HasPrefix(err.Error(), permissionDeniedPrefix) {
		t.Errorf("err = %v", err)
	}
}

func TestListsFollowTheContinuationTokenAcrossPages(t *testing.T) {
	f, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/p/_apis/policy/configurations": func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("continuationToken") == "" {
				w.Header().Set("x-ms-continuationtoken", "page2")
				_, _ = fmt.Fprint(w, `{"count":1,"value":[{"id":1,"type":{"id":"t"}}]}`)
				return
			}
			_, _ = fmt.Fprint(w, `[{"id":2,"type":{"id":"t"}}]`)
		},
	})
	var ids []int
	if err := readGitBranchPolicies(context.Background(), c, "p", func(p GitBranchPolicy) bool {
		ids = append(ids, p.Id)
		return true
	}); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(ids) != "[1 2]" || len(f.requests) != 2 {
		t.Errorf("ids %v over %d requests", ids, len(f.requests))
	}
	if got := f.requests[1].URL.Query().Get("continuationToken"); got != "page2" {
		t.Errorf("second page token = %q", got)
	}
}

func TestServiceURLMovesOrganizationServicesToTheirHost(t *testing.T) {
	for in, want := range map[string]string{
		"https://dev.azure.com/contoso":    "https://vsaex.dev.azure.com/contoso",
		"https://dev.azure.com/contoso/":   "https://vsaex.dev.azure.com/contoso",
		"https://contoso.visualstudio.com": "https://vsaex.dev.azure.com/contoso",
		"http://127.0.0.1:5000/myorg":      "http://127.0.0.1:5000/myorg",
	} {
		if got := serviceURL(in, "vsaex"); got != want {
			t.Errorf("serviceURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBranchPolicyHasOneRowPerScope(t *testing.T) {
	f, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/My Project/_apis/policy/configurations": body(`{"count":3,"value":[
		  {"id":18,"isEnabled":true,"isBlocking":true,"revision":1,"createdDate":"2015-02-23T12:51:06.8887894Z",
		   "type":{"id":"fa4e907d-c16b-4a4c-9dfa-4906e5d171dd","displayName":"Minimum number of reviewers"},
		   "settings":{"minimumApproverCount":2,"creatorVoteCounts":false,"resetOnSourcePush":true,"allowDownvotes":false,
		     "scope":[{"refName":null,"matchKind":"DefaultBranch","repositoryId":null},
		              {"refName":"refs/heads/main","matchKind":"Exact","repositoryId":"r1"}]}},
		  {"id":19,"isEnabled":true,"isBlocking":false,
		   "type":{"id":"0609b952-1397-4640-95ec-e00a01b2c241","displayName":"Build"},
		   "settings":{"buildDefinitionId":5,"validDuration":720,"scope":[{"refName":"refs/heads/releases/","matchKind":"Prefix","repositoryId":null}]}},
		  {"id":20,"isEnabled":true,"isBlocking":true,
		   "type":{"id":"2e26e725-8201-4edd-8bf5-978563c34a80","displayName":"File size restriction"},
		   "settings":{"maximumGitBlobSizeInBytes":1048576}}]}`),
	})
	var rows []GitBranchPolicy
	if err := readGitBranchPolicies(context.Background(), c, "My Project", func(p GitBranchPolicy) bool {
		rows = append(rows, p)
		return true
	}); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 4 {
		t.Fatalf("got %d rows, want 4 (two scopes, one scope, no scope)", len(rows))
	}
	if q := f.requests[0].URL.Query(); q.Get("api-version") != "7.1" || q.Get("repositoryId") != "" {
		t.Errorf("query = %v; the repository filter is Postgres's job", q)
	}
	def := rows[0]
	if def.RepositoryId != nil || def.RefName != nil || def.MatchKind == nil || *def.MatchKind != "DefaultBranch" {
		t.Errorf("default-branch scope = %+v", def)
	}
	exact := rows[1]
	if exact.Id != 18 || *exact.RepositoryId != "r1" || *exact.RefName != "refs/heads/main" || *exact.MatchKind != "Exact" {
		t.Errorf("exact scope = %+v", exact)
	}
	settings := def.Settings.(map[string]interface{})
	if settings["minimumApproverCount"] != float64(2) || settings["resetOnSourcePush"] != true || len(settings["scope"].([]interface{})) != 2 {
		t.Errorf("settings = %v", settings)
	}
	if def.CreatedDate == nil || def.TypeDisplayName != "Minimum number of reviewers" || !def.IsBlocking {
		t.Errorf("row = %+v", def)
	}
	build := rows[2].Settings.(map[string]interface{})
	if build["buildDefinitionId"] != float64(5) || build["validDuration"] != float64(720) {
		t.Errorf("build settings = %v", build)
	}
	if noScope := rows[3]; noScope.Id != 20 || noScope.RepositoryId != nil || noScope.MatchKind != nil {
		t.Errorf("unscoped policy = %+v", noScope)
	}
}

func TestPipelineSettingsKeepUnreportedSettingsNull(t *testing.T) {
	_, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/p/_apis/build/generalsettings": body(`{"enforceJobAuthScope":true,"enforceJobAuthScopeForReleases":false,
		  "enforceReferencedRepoScopedToken":true,"enforceSettableVar":true,"enforceNoAccessToSecretsFromForks":true,"newSetting":true}`),
	})
	s, err := readProjectPipelineSettings(context.Background(), c, "p")
	if err != nil {
		t.Fatal(err)
	}
	if !*s.EnforceJobAuthScope || *s.EnforceJobAuthScopeForReleases || !*s.EnforceNoAccessToSecretsFromForks {
		t.Errorf("settings = %+v", s)
	}
	if s.BuildsEnabledForForks != nil {
		t.Error("a setting the API did not report reads as false instead of null")
	}
	if s.Settings["newSetting"] != true {
		t.Error("settings column lost a field the struct does not know")
	}
}

func TestServiceEndpointPipelinePermissions(t *testing.T) {
	f, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/p/_apis/pipelines/pipelinepermissions/endpoint/5e47a0d8-c745-44f8-8f93-784f18ff31c4": body(`{"resource":{"type":"endpoint","id":"5e47a0d8-c745-44f8-8f93-784f18ff31c4"},
		  "allPipelines":{"authorized":true,"authorizedOn":"2020-01-02T08:59:48.77Z"},
		  "pipelines":[{"id":6,"authorized":true},{"id":7,"authorized":false},{"id":15,"authorized":true}]}`),
	})
	perms, err := readPipelinePermissions(context.Background(), c, "p", "endpoint", "5e47a0d8-c745-44f8-8f93-784f18ff31c4")
	if err != nil {
		t.Fatal(err)
	}
	if !perms.AllPipelinesAuthorized || fmt.Sprint(perms.AuthorizedPipelineIds) != "[6 15]" {
		t.Errorf("perms = %+v", perms)
	}
	if v := f.requests[0].URL.Query().Get("api-version"); v != "7.1-preview.1" {
		t.Errorf("api-version = %q", v)
	}
}

func TestServiceEndpointWithoutAllPipelinesIsNotOpenToAll(t *testing.T) {
	_, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/p/_apis/pipelines/pipelinepermissions/endpoint/e1": body(`{"resource":{"type":"endpoint","id":"e1"},"pipelines":[]}`),
	})
	perms, err := readPipelinePermissions(context.Background(), c, "p", "endpoint", "e1")
	if err != nil {
		t.Fatal(err)
	}
	if perms.AllPipelinesAuthorized || perms.AuthorizedPipelineIds == nil || len(perms.AuthorizedPipelineIds) != 0 {
		t.Errorf("perms = %+v", perms)
	}
}

func TestVariableGroupsCountSecretsAndNeverCarryValues(t *testing.T) {
	_, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/p/_apis/distributedtask/variablegroups": body(`{"count":2,"value":[
		  {"id":1,"type":"Vsts","name":"app-prod","variables":{
		    "DB_PASSWORD":{"value":"hunter2-plaintext"},
		    "ApiToken":{"value":null,"isSecret":true},
		    "BuildPath":{"value":"/src"},
		    "GitHubPAT":{"value":"ghp_plaintextvalue"}}},
		  {"id":2,"type":"AzureKeyVault","name":"kv","providerData":{"serviceEndpointId":"se","vault":"prod-kv"},
		   "variables":{"StoragePassword":{"enabled":true,"contentType":""}}}]}`),
	})
	var rows []VariableGroup
	if err := readVariableGroups(context.Background(), c, "p", func(g VariableGroup) bool {
		rows = append(rows, g)
		return true
	}); err != nil {
		t.Fatal(err)
	}
	vsts, kv := rows[0], rows[1]
	if vsts.VariableCount != 4 || vsts.SecretVariableCount != 1 || vsts.KeyVaultLinked {
		t.Errorf("vsts = %+v", vsts)
	}
	if got := fmt.Sprint(vsts.PlaintextVariableNamesSuspected); got != "[DB_PASSWORD GitHubPAT]" {
		t.Errorf("suspected = %s", got)
	}
	if !kv.KeyVaultLinked || kv.KeyVaultName == nil || *kv.KeyVaultName != "prod-kv" || len(kv.PlaintextVariableNamesSuspected) != 0 {
		t.Errorf("key vault group = %+v", kv)
	}
	out, _ := json.Marshal(rows)
	if strings.Contains(string(out), "hunter2") || strings.Contains(string(out), "ghp_") {
		t.Errorf("a variable value reached the row: %s", out)
	}
}

func TestLooksLikeSecretName(t *testing.T) {
	for name, want := range map[string]bool{
		"DB_PASSWORD": true, "dbPassword": true, "client_secret": true, "NPM_TOKEN": true, "ApiKey": true,
		"AZURE_STORAGE_KEY": true, "SigningKey": true, "GitHubPAT": true, "SMTP_PASS": true, "sql-pwd": true,
		"ConnectionString": true, "AWS_ACCESS_KEY_ID": true, "SasToken": true,
		"BuildConfiguration": false, "path": false, "PATH_PREFIX": false, "KeyVaultName": false, "Environment": false,
		"Key": false, "passRate": true, "compass": false,
	} {
		if got := looksLikeSecretName(name); got != want {
			t.Errorf("looksLikeSecretName(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestUserEntitlementsPageAndReportNeverAsNull(t *testing.T) {
	f, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/_apis/userentitlements": func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("continuationToken") == "" {
				_, _ = fmt.Fprint(w, `{"items":[{"id":"u1","user":{"displayName":"Ana","principalName":"ana@contoso.com","metaType":"member","subjectKind":"user","origin":"aad"},
				  "accessLevel":{"licenseDisplayName":"Basic","accountLicenseType":"express","licensingSource":"account","status":"active"},
				  "lastAccessedDate":"2026-07-01T10:00:00.123Z","dateCreated":"2024-01-01T00:00:00Z"}],"continuationToken":"c2","totalCount":2}`)
				return
			}
			_, _ = fmt.Fprint(w, `{"members":[{"id":"u2","user":{"displayName":"Ben","metaType":"guest"},
			  "accessLevel":{"licenseDisplayName":"Stakeholder","status":"pending"},
			  "lastAccessedDate":"0001-01-01T00:00:00Z","dateCreated":"2026-10-01T00:00:00Z"}],"continuationToken":null}`)
		},
	})
	var rows []UserEntitlement
	if err := readUserEntitlements(context.Background(), c, func(u UserEntitlement) bool {
		rows = append(rows, u)
		return true
	}); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || len(f.requests) != 2 {
		t.Fatalf("%d rows over %d requests", len(rows), len(f.requests))
	}
	ana, ben := rows[0], rows[1]
	if ana.PrincipalName != "ana@contoso.com" || ana.AccessLevel != "Basic" || ana.Status != "active" || ana.UserType != "member" ||
		ana.LastAccessedDate == nil || ana.LastAccessedDate.Month() != time.July {
		t.Errorf("ana = %+v", ana)
	}
	if ben.LastAccessedDate != nil || ben.DateCreated == nil || ben.UserType != "guest" {
		t.Errorf("ben = %+v; 0001-01-01 means never", ben)
	}
	if q := f.requests[0].URL.Query(); q.Get("api-version") != "7.1" {
		t.Errorf("query = %v", q)
	}
}

func TestUserEntitlementsWithoutScopeIsAPermissionError(t *testing.T) {
	_, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/_apis/userentitlements": status(http.StatusUnauthorized, "TF400813"),
	})
	err := readUserEntitlements(context.Background(), c, func(UserEntitlement) bool { return true })
	assertPermissionDenied(t, err, "vso.memberentitlementmanagement")
}

func TestAuditStreamsReadStatusAndDropCredentials(t *testing.T) {
	_, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/_apis/audit/streams": body(`[
		  {"id":1,"consumerType":"Splunk","displayName":"https://splunk:8088",
		   "consumerInputs":{"SplunkUrl":"https://splunk:8088","SplunkEventCollectorToken":"tok-secret"},
		   "status":1,"statusReason":null,"updatedTime":"/Date(1565715592526)/"},
		  {"id":2,"consumerType":"AzureMonitorLogs","displayName":"law","status":"disabledBySystem",
		   "statusReason":"Invalid credentials","createdTime":"2026-01-02T03:04:05Z"}]`),
	})
	var rows []AuditStream
	if err := readAuditStreams(context.Background(), c, func(s AuditStream) bool {
		rows = append(rows, s)
		return true
	}); err != nil {
		t.Fatal(err)
	}
	if rows[0].Status != "enabled" || rows[0].UpdatedTime == nil || rows[0].UpdatedTime.Year() != 2019 {
		t.Errorf("splunk = %+v", rows[0])
	}
	if rows[1].Status != "disabledBySystem" || rows[1].StatusReason == nil || rows[1].CreatedTime == nil {
		t.Errorf("law = %+v", rows[1])
	}
	out, _ := json.Marshal(rows)
	if strings.Contains(string(out), "tok-secret") {
		t.Errorf("stream credentials reached the row: %s", out)
	}
}

func TestAdvancedSecurityEnablementPerRepository(t *testing.T) {
	f, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/p/_apis/git/repositories": body(`{"count":4,"value":[
		  {"id":"r1","name":"api"},{"id":"r2","name":"web"},{"id":"r3","name":"old","isDisabled":true},{"id":"r4","name":"legacy"}]}`),
		"/myorg/p/_apis/management/repositories/r1/enablement": body(`{"projectId":"p","repositoryId":"r1",
		  "secretProtectionFeatures":{"secretProtectionEnabled":true,"blockPushes":true,"secretProtectionEnablementLastChangedDate":"2026-05-01T00:00:00Z"},
		  "codeSecurityFeatures":{"codeSecurityEnabled":true,"dependencyScanningInjectionEnabled":true,"codeQLEnabled":false}}`),
		"/myorg/p/_apis/management/repositories/r4/enablement": body(`{"advSecEnabled":true}`),
	})
	var rows []AdvancedSecurityEnablement
	if err := readAdvancedSecurityEnablement(context.Background(), c, "p", "", func(e AdvancedSecurityEnablement) bool {
		rows = append(rows, e)
		return true
	}); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 4 {
		t.Fatalf("got %d rows", len(rows))
	}
	api, web, old, legacy := rows[0], rows[1], rows[2], rows[3]
	if !api.Enabled || !api.SecretProtectionEnabled || !api.BlockPushes || !api.CodeSecurityEnabled || !api.DependencyScanningInjectionEnabled ||
		api.CodeqlEnabled || api.SecretProtectionLastChangedDate == nil || api.RepositoryName != "api" {
		t.Errorf("api = %+v", api)
	}
	if web.Enabled || web.SecretProtectionEnabled || web.CodeSecurityEnabled {
		t.Errorf("web (404, no Advanced Security) = %+v; want false, not an error", web)
	}
	if !old.IsDisabled || old.Enabled {
		t.Errorf("disabled repo = %+v", old)
	}
	if !legacy.Enabled || !legacy.SecretProtectionEnabled || !legacy.CodeSecurityEnabled {
		t.Errorf("legacy single switch = %+v", legacy)
	}
	for _, r := range f.requests {
		if strings.Contains(r.URL.Path, "/r3/") {
			t.Error("read enablement of a disabled repository")
		}
		if strings.HasSuffix(r.URL.Path, "/enablement") && (r.URL.Query().Get("includeAllProperties") != "true" || r.URL.Query().Get("api-version") != "7.2-preview.3") {
			t.Errorf("enablement query = %v", r.URL.Query())
		}
	}
}

func TestAdvancedSecurityEnablementWithoutScopeIsAPermissionError(t *testing.T) {
	_, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/p/_apis/git/repositories/api":                   body(`{"id":"8d1f-guid","name":"api"}`),
		"/myorg/p/_apis/management/repositories/api/enablement": status(http.StatusUnauthorized, "TF400813"),
	})
	err := readAdvancedSecurityEnablement(context.Background(), c, "p", "api", func(AdvancedSecurityEnablement) bool { return true })
	assertPermissionDenied(t, err, "vso.advsec")
}

func TestAdvancedSecurityEnablementEchoesTheRepositoryQualifier(t *testing.T) {
	_, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/p/_apis/git/repositories/api": body(`{"id":"8d1f-guid","name":"api"}`),
	})
	var rows []AdvancedSecurityEnablement
	if err := readAdvancedSecurityEnablement(context.Background(), c, "p", "api", func(e AdvancedSecurityEnablement) bool {
		rows = append(rows, e)
		return true
	}); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].RepositoryId != "api" || rows[0].RepositoryName != "api" {
		t.Errorf("rows = %+v", rows)
	}
}

func TestAdvancedSecurityAlerts(t *testing.T) {
	f, c := routes(t, map[string]func(http.ResponseWriter, *http.Request){
		"/myorg/p/_apis/git/repositories": body(`{"value":[{"id":"r1","name":"api"},{"id":"r2","name":"web"}]}`),
		"/myorg/p/_apis/alert/repositories/r1/alerts": func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("continuationToken") == "" {
				w.Header().Set("x-ms-continuationtoken", "next")
				_, _ = fmt.Fprint(w, `{"count":1,"value":[{"alertId":11,"alertType":"secret","state":"active","severity":"critical",
				  "title":"Azure Storage key","truncatedSecret":"AbC***","firstSeenDate":"2026-09-01T00:00:00Z"}]}`)
				return
			}
			_, _ = fmt.Fprint(w, `{"count":1,"value":[{"alertId":12,"alertType":"dependency","state":"active","severity":"high","title":"lodash"}]}`)
		},
	})
	var rows []AdvancedSecurityAlert
	err := readAdvancedSecurityAlerts(context.Background(), c, advancedSecurityAlertRequest{Project: "p", State: "active"}, func(a AdvancedSecurityAlert) bool {
		rows = append(rows, a)
		return true
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].AlertId != 11 || rows[0].RepositoryId != "r1" || rows[0].FirstSeenDate == nil || rows[1].AlertType != "dependency" {
		t.Errorf("rows = %+v (r2 answers 404 and should add none)", rows)
	}
	out, _ := json.Marshal(rows)
	if strings.Contains(string(out), "AbC") {
		t.Errorf("secret material reached the row: %s", out)
	}
	for _, r := range f.requests {
		if strings.HasSuffix(r.URL.Path, "/alerts") && r.URL.Query().Get("criteria.states") != "active" {
			t.Errorf("state qualifier not sent: %v", r.URL.Query())
		}
	}
}
