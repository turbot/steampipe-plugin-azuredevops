package azuredevops

import (
	"context"
	"encoding/json"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const policyConfigurationsAPIVersion = "7.1"

// https://learn.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/list
func tableAzureDevOpsGitBranchPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_git_branch_policy",
		Description: "Branch and repository policies configured in an Azure DevOps project, one row per policy scope.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project", Require: plugin.Required},
			},
			Hydrate: listGitBranchPolicies,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "project", Description: "Name or ID of the project, as given in the qualifier.", Type: proto.ColumnType_STRING, Transform: transform.FromQual("project")},
			{Name: "created_date", Description: "When this revision of the policy was created.", Type: proto.ColumnType_TIMESTAMP},
			{Name: "id", Description: "The policy configuration ID. Repeats when a policy has more than one scope.", Type: proto.ColumnType_INT},
			{Name: "is_blocking", Description: "True if the policy is required, false if it is optional.", Type: proto.ColumnType_BOOL},
			{Name: "is_deleted", Description: "True if the policy has been soft-deleted.", Type: proto.ColumnType_BOOL},
			{Name: "is_enabled", Description: "True if the policy is enabled.", Type: proto.ColumnType_BOOL},
			{Name: "is_enterprise_managed", Description: "True if the policy needs Manage Enterprise Policies permission to change.", Type: proto.ColumnType_BOOL},
			{Name: "match_kind", Description: "How ref_name matches: Exact, Prefix or DefaultBranch.", Type: proto.ColumnType_STRING},
			{Name: "ref_name", Description: "Branch this scope applies to, e.g. refs/heads/main, or a prefix when match_kind is Prefix. Null for a DefaultBranch scope and for repository-wide policies.", Type: proto.ColumnType_STRING},
			{Name: "repository_id", Description: "Repository this scope applies to. Null for a scope that covers every repository in the project.", Type: proto.ColumnType_STRING},
			{Name: "revision", Description: "Revision of the policy configuration.", Type: proto.ColumnType_INT},
			{Name: "settings", Description: "The policy settings as Azure DevOps returns them, including the full scope array.", Type: proto.ColumnType_JSON},
			{Name: "type_display_name", Description: "Display name of the policy type.", Type: proto.ColumnType_STRING},
			{Name: "type_id", Description: "ID of the policy type, e.g. fa4e907d-c16b-4a4c-9dfa-4906e5d171dd for minimum number of reviewers.", Type: proto.ColumnType_STRING},
			{Name: "url", Description: "REST URL of the policy configuration.", Type: proto.ColumnType_STRING},
			{Name: "title", Description: "Title of the resource.", Type: proto.ColumnType_STRING, Transform: transform.FromField("TypeDisplayName")},
		}),
	}
}

type policyConfiguration struct {
	Id   int `json:"id"`
	Type struct {
		Id          string `json:"id"`
		DisplayName string `json:"displayName"`
	} `json:"type"`
	IsEnabled           bool            `json:"isEnabled"`
	IsBlocking          bool            `json:"isBlocking"`
	IsDeleted           bool            `json:"isDeleted"`
	IsEnterpriseManaged bool            `json:"isEnterpriseManaged"`
	Revision            int             `json:"revision"`
	CreatedDate         *time.Time      `json:"createdDate"`
	Url                 string          `json:"url"`
	Settings            json.RawMessage `json:"settings"`
}

type policyScope struct {
	RepositoryId *string `json:"repositoryId"`
	RefName      *string `json:"refName"`
	MatchKind    *string `json:"matchKind"`
}

type GitBranchPolicy struct {
	Id                  int
	TypeId              string
	TypeDisplayName     string
	IsEnabled           bool
	IsBlocking          bool
	IsDeleted           bool
	IsEnterpriseManaged bool
	RepositoryId        *string
	RefName             *string
	MatchKind           *string
	Settings            interface{}
	Revision            int
	CreatedDate         *time.Time
	Url                 string
}

func listGitBranchPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	project := d.EqualsQualString("project")
	if project == "" {
		return nil, nil
	}
	client, err := restClientFor(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_branch_policy.listGitBranchPolicies", "connection_error", err)
		return nil, err
	}
	err = readGitBranchPolicies(ctx, client, project, func(p GitBranchPolicy) bool {
		d.StreamListItem(ctx, p)
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_branch_policy.listGitBranchPolicies", "api_error", err)
		return nil, err
	}
	return nil, nil
}

// readGitBranchPolicies reads every policy in the project. Repository and branch are left to
// Postgres to filter: a server-side repository filter would drop the project-wide policies, whose
// null repository_id then fails the qualifier.
func readGitBranchPolicies(ctx context.Context, client *restClient, project string, emit func(GitBranchPolicy) bool) error {
	params := apiVersion(policyConfigurationsAPIVersion)
	keepGoing := true
	return getList(ctx, client, projectURL(client.baseURL, project, "policy/configurations"), params, scopeCode, func(pc policyConfiguration) bool {
		for _, row := range branchPolicyRows(pc) {
			if keepGoing = emit(row); !keepGoing {
				break
			}
		}
		return keepGoing
	})
}

func branchPolicyRows(pc policyConfiguration) []GitBranchPolicy {
	base := GitBranchPolicy{
		Id:                  pc.Id,
		TypeId:              pc.Type.Id,
		TypeDisplayName:     pc.Type.DisplayName,
		IsEnabled:           pc.IsEnabled,
		IsBlocking:          pc.IsBlocking,
		IsDeleted:           pc.IsDeleted,
		IsEnterpriseManaged: pc.IsEnterpriseManaged,
		Revision:            pc.Revision,
		CreatedDate:         pc.CreatedDate,
		Url:                 pc.Url,
	}
	var settings struct {
		Scope []policyScope `json:"scope"`
	}
	if len(pc.Settings) > 0 {
		_ = json.Unmarshal(pc.Settings, &base.Settings)
		_ = json.Unmarshal(pc.Settings, &settings)
	}
	if len(settings.Scope) == 0 {
		return []GitBranchPolicy{base}
	}
	rows := make([]GitBranchPolicy, 0, len(settings.Scope))
	for _, s := range settings.Scope {
		row := base
		row.RepositoryId, row.RefName, row.MatchKind = s.RepositoryId, s.RefName, s.MatchKind
		rows = append(rows, row)
	}
	return rows
}
