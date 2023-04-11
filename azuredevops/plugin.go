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
		// DefaultIgnoreConfig: &plugin.IgnoreConfig{
		// 	ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"404"}),
		// },
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			// "azuredevops_build":            tableAzureDevOpsBuild(ctx),
			// "azuredevops_build_definition": tableAzureDevOpsBuildDefinition(ctx),
			// "azuredevops_git_repository":   tableAzureDevOpsGetRepository(ctx),
			// "azuredevops_pipeline":         tableAzureDevOpsPipeline(ctx),
			"azuredevops_project":     tableAzureDevOpsProject(ctx),
			"azuredevops_account":     tableAzureDevOpsAccount(ctx),
			"azuredevops_team":        tableAzureDevOpsTeam(ctx),
			"azuredevops_team_member": tableAzureDevOpsTeamMember(ctx),
		},
	}
	return p
}
