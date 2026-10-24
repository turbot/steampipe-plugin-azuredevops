package azuredevops

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const buildGeneralSettingsAPIVersion = "7.1"

// https://learn.microsoft.com/en-us/rest/api/azure/devops/build/general-settings/get
func tableAzureDevOpsProjectPipelineSettings(_ context.Context) *plugin.Table {
	setting := func(name, description string) *plugin.Column {
		return &plugin.Column{Name: name, Description: description, Type: proto.ColumnType_BOOL}
	}
	return &plugin.Table{
		Name:        "azuredevops_project_pipeline_settings",
		Description: "Pipeline general settings of an Azure DevOps project: job authorization scope, settable variables and pull requests from forks.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project", Require: plugin.Required},
			},
			Hydrate: listProjectPipelineSettings,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "project", Description: "Name or ID of the project, as given in the qualifier.", Type: proto.ColumnType_STRING, Transform: transform.FromQual("project")},
			setting("audit_enforce_settable_var", "Audits runs that set variables not marked settable at queue time."),
			setting("builds_enabled_for_forks", "Pull requests from forked repositories can trigger builds."),
			setting("disable_classic_build_pipeline_creation", "Classic build pipelines cannot be created."),
			setting("disable_classic_pipeline_creation", "Classic build and release pipelines cannot be created."),
			setting("disable_classic_release_pipeline_creation", "Classic release pipelines cannot be created."),
			setting("enable_shell_tasks_args_sanitizing", "Arguments of shell tasks are sanitized."),
			setting("enable_shell_tasks_args_sanitizing_audit", "Unsanitized shell task arguments are audited."),
			setting("enforce_job_auth_scope", "Limits the job authorization scope of non-release pipelines to the current project."),
			setting("enforce_job_auth_scope_for_forks", "Builds of forks run with a job token scoped to the project."),
			setting("enforce_job_auth_scope_for_releases", "Limits the job authorization scope of release pipelines to the current project."),
			setting("enforce_no_access_to_secrets_from_forks", "Builds of forks cannot read secrets."),
			setting("enforce_referenced_repo_scoped_token", "Limits the job token to the repositories a YAML pipeline explicitly references."),
			setting("enforce_settable_var", "Only variables marked settable at queue time can be set at queue time."),
			setting("fork_protection_enabled", "Fork protection settings apply to pull requests from forks."),
			setting("is_comment_required_for_pull_request", "A comment is required before a pull request from a fork is built."),
			setting("publish_pipeline_metadata", "Pipelines record metadata."),
			setting("require_comments_for_non_team_member_and_non_contributors", "A comment is required before a pull request from a non-team member or non-contributor is built."),
			setting("require_comments_for_non_team_members_only", "A comment is required before a pull request from a non-team member is built."),
			{Name: "settings", Description: "All general settings as Azure DevOps returns them.", Type: proto.ColumnType_JSON},
			setting("status_badges_are_private", "Status badges are not readable anonymously."),
		}),
	}
}

// ProjectPipelineSettings holds the general settings. A field the API omits stays null rather
// than reading as false.
type ProjectPipelineSettings struct {
	EnforceJobAuthScope                               *bool                  `json:"enforceJobAuthScope"`
	EnforceJobAuthScopeForReleases                    *bool                  `json:"enforceJobAuthScopeForReleases"`
	EnforceReferencedRepoScopedToken                  *bool                  `json:"enforceReferencedRepoScopedToken"`
	EnforceSettableVar                                *bool                  `json:"enforceSettableVar"`
	AuditEnforceSettableVar                           *bool                  `json:"auditEnforceSettableVar"`
	BuildsEnabledForForks                             *bool                  `json:"buildsEnabledForForks"`
	ForkProtectionEnabled                             *bool                  `json:"forkProtectionEnabled"`
	EnforceJobAuthScopeForForks                       *bool                  `json:"enforceJobAuthScopeForForks"`
	EnforceNoAccessToSecretsFromForks                 *bool                  `json:"enforceNoAccessToSecretsFromForks"`
	IsCommentRequiredForPullRequest                   *bool                  `json:"isCommentRequiredForPullRequest"`
	RequireCommentsForNonTeamMembersOnly              *bool                  `json:"requireCommentsForNonTeamMembersOnly"`
	RequireCommentsForNonTeamMemberAndNonContributors *bool                  `json:"requireCommentsForNonTeamMemberAndNonContributors"`
	PublishPipelineMetadata                           *bool                  `json:"publishPipelineMetadata"`
	StatusBadgesArePrivate                            *bool                  `json:"statusBadgesArePrivate"`
	DisableClassicPipelineCreation                    *bool                  `json:"disableClassicPipelineCreation"`
	DisableClassicBuildPipelineCreation               *bool                  `json:"disableClassicBuildPipelineCreation"`
	DisableClassicReleasePipelineCreation             *bool                  `json:"disableClassicReleasePipelineCreation"`
	EnableShellTasksArgsSanitizing                    *bool                  `json:"enableShellTasksArgsSanitizing"`
	EnableShellTasksArgsSanitizingAudit               *bool                  `json:"enableShellTasksArgsSanitizingAudit"`
	Settings                                          map[string]interface{} `json:"-"`
}

func listProjectPipelineSettings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	project := d.EqualsQualString("project")
	if project == "" {
		return nil, nil
	}
	client, err := restClientFor(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_project_pipeline_settings.listProjectPipelineSettings", "connection_error", err)
		return nil, err
	}
	settings, err := readProjectPipelineSettings(ctx, client, project)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_project_pipeline_settings.listProjectPipelineSettings", "api_error", err)
		return nil, err
	}
	if settings != nil {
		d.StreamListItem(ctx, *settings)
	}
	return nil, nil
}

func readProjectPipelineSettings(ctx context.Context, client *restClient, project string) (*ProjectPipelineSettings, error) {
	endpoint := projectURL(client.baseURL, project, "build/generalsettings") + "?" + apiVersion(buildGeneralSettingsAPIVersion).Encode()
	var raw map[string]interface{}
	err := client.getJSON(ctx, endpoint, scopeProject, &raw)
	if isRestNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	settings := &ProjectPipelineSettings{Settings: raw}
	if err := remarshal(raw, settings); err != nil {
		return nil, err
	}
	return settings, nil
}
