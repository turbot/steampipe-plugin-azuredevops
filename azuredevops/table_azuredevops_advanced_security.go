package azuredevops

import (
	"context"
	"net/url"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

const advancedSecurityEnablementAPIVersion = "7.2-preview.3"

// https://learn.microsoft.com/en-us/rest/api/azure/devops/advancedsecurity/repo-enablement/get
func tableAzureDevOpsAdvancedSecurityEnablement(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_advanced_security_enablement",
		Description: "Whether GitHub Advanced Security for Azure DevOps (Secret Protection and Code Security) is enabled on each repository of a project.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project", Require: plugin.Required},
				{Name: "repository_id", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
			Hydrate: listAdvancedSecurityEnablement,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "project", Description: "Name or ID of the project, as given in the qualifier.", Type: proto.ColumnType_STRING, Transform: transform.FromQual("project")},
			{Name: "repository_id", Description: "ID of the repository, or the name or ID given in the qualifier.", Type: proto.ColumnType_STRING},
			{Name: "block_pushes", Description: "True if pushes that contain secrets are blocked (push protection).", Type: proto.ColumnType_BOOL},
			{Name: "code_security_enabled", Description: "True if Code Security (code and dependency scanning) is enabled.", Type: proto.ColumnType_BOOL},
			{Name: "code_security_last_changed_date", Description: "When Code Security was last turned on or off.", Type: proto.ColumnType_TIMESTAMP},
			{Name: "codeql_enabled", Description: "True if CodeQL default setup is enabled.", Type: proto.ColumnType_BOOL},
			{Name: "dependency_scanning_injection_enabled", Description: "True if dependency scanning is injected into pipelines automatically.", Type: proto.ColumnType_BOOL},
			{Name: "enabled", Description: "True if Secret Protection or Code Security is enabled. False when the organization has no Advanced Security.", Type: proto.ColumnType_BOOL},
			{Name: "is_disabled", Description: "True if the repository itself is disabled; its enablement is not read and reports false.", Type: proto.ColumnType_BOOL},
			{Name: "repository_name", Description: "Name of the repository.", Type: proto.ColumnType_STRING},
			{Name: "secret_protection_enabled", Description: "True if Secret Protection (secret scanning) is enabled.", Type: proto.ColumnType_BOOL},
			{Name: "secret_protection_last_changed_date", Description: "When Secret Protection was last turned on or off.", Type: proto.ColumnType_TIMESTAMP},
			{Name: "title", Description: "Title of the resource.", Type: proto.ColumnType_STRING, Transform: transform.FromField("RepositoryName")},
		}),
	}
}

type repoEnablementResponse struct {
	SecretProtectionFeatures *struct {
		SecretProtectionEnabled                   bool     `json:"secretProtectionEnabled"`
		BlockPushes                               *bool    `json:"blockPushes"`
		SecretProtectionEnablementLastChangedDate *adoTime `json:"secretProtectionEnablementLastChangedDate"`
	} `json:"secretProtectionFeatures"`
	CodeSecurityFeatures *struct {
		CodeSecurityEnabled                   bool     `json:"codeSecurityEnabled"`
		DependencyScanningInjectionEnabled    bool     `json:"dependencyScanningInjectionEnabled"`
		CodeQLEnabled                         bool     `json:"codeQLEnabled"`
		CodeSecurityEnablementLastChangedDate *adoTime `json:"codeSecurityEnablementLastChangedDate"`
	} `json:"codeSecurityFeatures"`
	// Before Secret Protection and Code Security were sold separately there was one switch.
	AdvSecEnabled *bool `json:"advSecEnabled"`
}

type AdvancedSecurityEnablement struct {
	RepositoryId                       string
	RepositoryName                     string
	IsDisabled                         bool
	Enabled                            bool
	SecretProtectionEnabled            bool
	BlockPushes                        bool
	CodeSecurityEnabled                bool
	DependencyScanningInjectionEnabled bool
	CodeqlEnabled                      bool
	SecretProtectionLastChangedDate    *time.Time
	CodeSecurityLastChangedDate        *time.Time
}

func listAdvancedSecurityEnablement(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	project := d.EqualsQualString("project")
	if project == "" {
		return nil, nil
	}
	client, err := restClientFor(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_advanced_security_enablement.listAdvancedSecurityEnablement", "connection_error", err)
		return nil, err
	}
	err = readAdvancedSecurityEnablement(ctx, client, project, d.EqualsQualString("repository_id"), func(e AdvancedSecurityEnablement) bool {
		d.StreamListItem(ctx, e)
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_advanced_security_enablement.listAdvancedSecurityEnablement", "api_error", err)
		return nil, err
	}
	return nil, nil
}

func readAdvancedSecurityEnablement(ctx context.Context, client *restClient, project, repositoryID string, emit func(AdvancedSecurityEnablement) bool) error {
	repos, err := repositoriesInProject(ctx, client, project, repositoryID)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		row := AdvancedSecurityEnablement{RepositoryId: repo.Id, RepositoryName: repo.Name, IsDisabled: repo.IsDisabled}
		if !repo.IsDisabled {
			if err := readRepoEnablement(ctx, client, project, repo.Id, &row); err != nil {
				return err
			}
		}
		if !emit(row) {
			return nil
		}
	}
	return nil
}

// readRepoEnablement fills row from the enablement API. A 404, which an organization without
// Advanced Security returns, leaves everything false.
func readRepoEnablement(ctx context.Context, client *restClient, project, repositoryID string, row *AdvancedSecurityEnablement) error {
	params := apiVersion(advancedSecurityEnablementAPIVersion)
	params.Set("includeAllProperties", "true")
	endpoint := serviceURL(client.baseURL, "advsec") + "/" + url.PathEscape(project) +
		"/_apis/management/repositories/" + url.PathEscape(repositoryID) + "/enablement?" + params.Encode()
	var resp repoEnablementResponse
	err := client.getJSON(ctx, endpoint, scopeAdvancedSecurity, &resp)
	if isRestNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if sp := resp.SecretProtectionFeatures; sp != nil {
		row.SecretProtectionEnabled = sp.SecretProtectionEnabled
		row.BlockPushes = sp.BlockPushes != nil && *sp.BlockPushes
		if sp.SecretProtectionEnablementLastChangedDate != nil {
			row.SecretProtectionLastChangedDate = sp.SecretProtectionEnablementLastChangedDate.ptr()
		}
	}
	if cs := resp.CodeSecurityFeatures; cs != nil {
		row.CodeSecurityEnabled = cs.CodeSecurityEnabled
		row.DependencyScanningInjectionEnabled = cs.DependencyScanningInjectionEnabled
		row.CodeqlEnabled = cs.CodeQLEnabled
		if cs.CodeSecurityEnablementLastChangedDate != nil {
			row.CodeSecurityLastChangedDate = cs.CodeSecurityEnablementLastChangedDate.ptr()
		}
	}
	if resp.SecretProtectionFeatures == nil && resp.CodeSecurityFeatures == nil && resp.AdvSecEnabled != nil {
		row.SecretProtectionEnabled = *resp.AdvSecEnabled
		row.CodeSecurityEnabled = *resp.AdvSecEnabled
	}
	row.Enabled = row.SecretProtectionEnabled || row.CodeSecurityEnabled
	return nil
}

// https://learn.microsoft.com/en-us/rest/api/azure/devops/advancedsecurity/alerts/list
func tableAzureDevOpsAdvancedSecurityAlert(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_advanced_security_alert",
		Description: "Advanced Security alerts (dependency, secret and code scanning) on the repositories of a project. Secret material is never exposed.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project", Require: plugin.Required},
				{Name: "repository_id", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "state", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "alert_type", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
			Hydrate: listAdvancedSecurityAlerts,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "project", Description: "Name or ID of the project, as given in the qualifier.", Type: proto.ColumnType_STRING, Transform: transform.FromQual("project")},
			{Name: "repository_id", Description: "ID of the repository, or the name or ID given in the qualifier.", Type: proto.ColumnType_STRING},
			{Name: "alert_type", Description: "Type of the alert: dependency, secret, code, aiCode, malware or unknown.", Type: proto.ColumnType_STRING},
			{Name: "state", Description: "State of the alert: active, dismissed, fixed, autoDismissed or unknown.", Type: proto.ColumnType_STRING},
			{Name: "alert_id", Description: "ID of the alert, unique in the organization.", Type: proto.ColumnType_INT},
			{Name: "first_seen_date", Description: "When the alert was first detected.", Type: proto.ColumnType_TIMESTAMP},
			{Name: "fixed_date", Description: "When the alert was seen fixed.", Type: proto.ColumnType_TIMESTAMP},
			{Name: "git_ref", Description: "Branch the alert was found on.", Type: proto.ColumnType_STRING},
			{Name: "introduced_date", Description: "When the vulnerability was introduced.", Type: proto.ColumnType_TIMESTAMP},
			{Name: "last_seen_date", Description: "When the alert was last detected.", Type: proto.ColumnType_TIMESTAMP},
			{Name: "severity", Description: "Severity of the alert: critical, high, medium, low, note, warning, error or undefined.", Type: proto.ColumnType_STRING},
			{Name: "title", Description: "Title of the alert.", Type: proto.ColumnType_STRING},
		}),
	}
}

// advancedSecurityAlertResponse decodes only summary fields; truncatedSecret and the locations are left out.
type advancedSecurityAlertResponse struct {
	AlertId        int64    `json:"alertId"`
	AlertType      string   `json:"alertType"`
	State          string   `json:"state"`
	Severity       string   `json:"severity"`
	Title          string   `json:"title"`
	GitRef         string   `json:"gitRef"`
	FirstSeenDate  *adoTime `json:"firstSeenDate"`
	LastSeenDate   *adoTime `json:"lastSeenDate"`
	FixedDate      *adoTime `json:"fixedDate"`
	IntroducedDate *adoTime `json:"introducedDate"`
}

type AdvancedSecurityAlert struct {
	RepositoryId   string
	AlertId        int64
	AlertType      string
	State          string
	Severity       string
	Title          string
	GitRef         string
	FirstSeenDate  *time.Time
	LastSeenDate   *time.Time
	FixedDate      *time.Time
	IntroducedDate *time.Time
}

type advancedSecurityAlertRequest struct {
	Project      string
	RepositoryID string
	State        string
	AlertType    string
}

func listAdvancedSecurityAlerts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	req := advancedSecurityAlertRequest{
		Project:      d.EqualsQualString("project"),
		RepositoryID: d.EqualsQualString("repository_id"),
		State:        d.EqualsQualString("state"),
		AlertType:    d.EqualsQualString("alert_type"),
	}
	if req.Project == "" {
		return nil, nil
	}
	client, err := restClientFor(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_advanced_security_alert.listAdvancedSecurityAlerts", "connection_error", err)
		return nil, err
	}
	err = readAdvancedSecurityAlerts(ctx, client, req, func(a AdvancedSecurityAlert) bool {
		d.StreamListItem(ctx, a)
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_advanced_security_alert.listAdvancedSecurityAlerts", "api_error", err)
		return nil, err
	}
	return nil, nil
}

// readAdvancedSecurityAlerts lists alerts per repository. A repository without Advanced Security
// answers 404 and contributes no rows.
func readAdvancedSecurityAlerts(ctx context.Context, client *restClient, req advancedSecurityAlertRequest, emit func(AdvancedSecurityAlert) bool) error {
	repos, err := repositoriesInProject(ctx, client, req.Project, req.RepositoryID)
	if err != nil {
		return err
	}
	base := serviceURL(client.baseURL, "advsec") + "/" + url.PathEscape(req.Project) + "/_apis/alert/repositories/"
	for _, repo := range repos {
		if repo.IsDisabled {
			continue
		}
		params := apiVersion(advancedSecurityAlertsAPIVersion)
		if req.State != "" {
			params.Set("criteria.states", req.State)
		}
		if req.AlertType != "" {
			params.Set("criteria.alertType", req.AlertType)
		}
		stopped := false
		err := getList(ctx, client, base+url.PathEscape(repo.Id)+"/alerts", params, scopeAdvancedSecurity, func(a advancedSecurityAlertResponse) bool {
			if !emit(alertRow(repo.Id, a)) {
				stopped = true
			}
			return !stopped
		})
		if err != nil && !isRestNotFound(err) {
			return err
		}
		if stopped {
			return nil
		}
	}
	return nil
}

func alertRow(repositoryID string, a advancedSecurityAlertResponse) AdvancedSecurityAlert {
	row := AdvancedSecurityAlert{
		RepositoryId: repositoryID, AlertId: a.AlertId, AlertType: a.AlertType, State: a.State,
		Severity: a.Severity, Title: a.Title, GitRef: a.GitRef,
	}
	for dst, src := range map[**time.Time]*adoTime{
		&row.FirstSeenDate: a.FirstSeenDate, &row.LastSeenDate: a.LastSeenDate,
		&row.FixedDate: a.FixedDate, &row.IntroducedDate: a.IntroducedDate,
	} {
		if src != nil {
			*dst = src.ptr()
		}
	}
	return row
}
