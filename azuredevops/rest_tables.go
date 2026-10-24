package azuredevops

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// PAT scopes named in permission-denied errors, as "<scope id> (<name in the PAT dialog>)".
const (
	scopeCode                        = "vso.code (Code: Read)"
	scopeProject                     = "vso.project (Project and Team: Read)"
	scopeBuild                       = "vso.build (Build: Read)"
	scopeVariableGroups              = "vso.variablegroups_read (Variable Groups: Read)"
	scopeMemberEntitlement           = "vso.memberentitlementmanagement (Member Entitlement Management: Read)"
	scopeAuditLog                    = "vso.auditlog (Audit Log: Read)"
	scopeAdvancedSecurity            = "vso.advsec (Advanced Security: Read)"
	gitRepositoriesAPIVersion        = "7.1"
	advancedSecurityAlertsAPIVersion = "7.2-preview.1"
)

func restClientFor(ctx context.Context, d *plugin.QueryData) (*restClient, error) {
	connection, err := getConnection(ctx, d)
	if err != nil {
		return nil, err
	}
	return newRestClient(connection.BaseUrl, connection.AuthorizationString), nil
}

func apiVersion(v string) url.Values {
	params := url.Values{}
	params.Set("api-version", v)
	return params
}

// adoTime reads an Azure DevOps timestamp, which reports "never" as 0001-01-01.
type adoTime struct{ t *time.Time }

func (a *adoTime) UnmarshalJSON(b []byte) error {
	var s *string
	if err := json.Unmarshal(b, &s); err != nil || s == nil || *s == "" {
		return err
	}
	t, err := time.Parse(time.RFC3339Nano, *s)
	if err != nil {
		return nil
	}
	if t.Year() > 1 {
		a.t = &t
	}
	return nil
}

func (a adoTime) ptr() *time.Time { return a.t }

type gitRepositoryRef struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	IsDisabled bool   `json:"isDisabled"`
}

// repositoriesInProject returns the one repository named by repositoryID, or every repository in the
// project when it is empty. A repository that does not exist yields none.
func repositoriesInProject(ctx context.Context, client *restClient, project, repositoryID string) ([]gitRepositoryRef, error) {
	params := apiVersion(gitRepositoriesAPIVersion)
	if repositoryID != "" {
		repo := gitRepositoryRef{}
		err := client.getJSON(ctx, projectURL(client.baseURL, project, "git/repositories/"+url.PathEscape(repositoryID))+"?"+params.Encode(), scopeCode, &repo)
		if isRestNotFound(err) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		// Rows echo the qualifier, because Postgres re-applies it to every row.
		repo.Id = repositoryID
		return []gitRepositoryRef{repo}, nil
	}
	var repos []gitRepositoryRef
	err := getList(ctx, client, projectURL(client.baseURL, project, "git/repositories"), params, scopeCode, func(r gitRepositoryRef) bool {
		repos = append(repos, r)
		return true
	})
	return repos, err
}

// remarshal copies a decoded JSON value into a typed struct.
func remarshal(from, to interface{}) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, to)
}
