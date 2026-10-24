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
		},
		ConnectionKeyColumns: []plugin.ConnectionKeyColumn{
			{
				Name:    "organization",
				Hydrate: getOrganization,
			},
		},
		TableMap: map[string]*plugin.Table{
			"azuredevops_advanced_security_alert":      tableAzureDevOpsAdvancedSecurityAlert(ctx),
			"azuredevops_advanced_security_enablement": tableAzureDevOpsAdvancedSecurityEnablement(ctx),
			"azuredevops_audit_stream":                 tableAzureDevOpsAuditStream(ctx),
			"azuredevops_build":                        tableAzureDevOpsBuild(ctx),
			"azuredevops_build_definition":             tableAzureDevOpsBuildDefinition(ctx),
			"azuredevops_dashboard":                    tableAzureDevOpsDashboard(ctx),
			"azuredevops_git_branch_policy":            tableAzureDevOpsGitBranchPolicy(ctx),
			"azuredevops_git_item":                     tableAzureDevOpsGitItem(ctx),
			"azuredevops_git_repository":               tableAzureDevOpsGitRepository(ctx),
			"azuredevops_git_repository_branch":        tableAzureDevOpsGitRepositoryBranch(ctx),
			"azuredevops_group":                        tableAzureDevOpsGroup(ctx),
			"azuredevops_pipeline":                     tableAzureDevOpsPipeline(ctx),
			"azuredevops_project":                      tableAzureDevOpsProject(ctx),
			"azuredevops_project_pipeline_settings":    tableAzureDevOpsProjectPipelineSettings(ctx),
			"azuredevops_release":                      tableAzureDevOpsRelease(ctx),
			"azuredevops_serviceendpoint":              tableAzureDevOpsServiceEndpoint(ctx),
			"azuredevops_team":                         tableAzureDevOpsTeam(ctx),
			"azuredevops_team_member":                  tableAzureDevOpsTeamMember(ctx),
			"azuredevops_user":                         tableAzureDevOpsUser(ctx),
			"azuredevops_user_entitlement":             tableAzureDevOpsUserEntitlement(ctx),
			"azuredevops_variable_group":               tableAzureDevOpsVariableGroup(ctx),
		},
	}
	return p
}
