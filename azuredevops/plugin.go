package azuredevops

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin creates this (azuredevops) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-azuredevops",
		DefaultTransform: transform.FromCamel(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"azuredevops_build":                 tableAzureDevOpsBuild(ctx),
			"azuredevops_build_definition":      tableAzureDevOpsBuildDefinition(ctx),
			"azuredevops_dashboard":             tableAzureDevOpsDashboard(ctx),
			"azuredevops_git_repository":        tableAzureDevOpsGitRepository(ctx),
			"azuredevops_git_repository_branch": tableAzureDevOpsGitRepositoryBranch(ctx),
			"azuredevops_group":                 tableAzureDevOpsGroup(ctx),
			"azuredevops_pipeline":              tableAzureDevOpsPipeline(ctx),
			"azuredevops_project":               tableAzureDevOpsProject(ctx),
			"azuredevops_release":               tableAzureDevOpsRelease(ctx),
			"azuredevops_serviceendpoint":       tableAzureDevOpsServiceEndpoint(ctx),
			"azuredevops_team":                  tableAzureDevOpsTeam(ctx),
			"azuredevops_team_member":           tableAzureDevOpsTeamMember(ctx),
			"azuredevops_user":                  tableAzureDevOpsUser(ctx),
		},
	}
	return p
}
