package azuredevops

import (
	"context"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/git"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureDevOpsGitRepositoryBranch(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_git_repository_branch",
		Description: "Retrieve information about your repository branches.",
		List: &plugin.ListConfig{
			ParentHydrate: listGitRepositories,
			Hydrate:       listGitRepositoryBranches,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_id", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "repository_id"}),
			Hydrate:    getRepositoryBranch,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the ref.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository_id",
				Description: "The repository id.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ahead_count",
				Description: "Number of commits ahead.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "behind_count",
				Description: "Number of commits behind.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "is_base_version",
				Description: "True if this is the result for the base version.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "commit",
				Description: "Current commit.",
				Type:        proto.ColumnType_JSON,
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		},
	}
}

type Branch struct {
	git.GitBranchStats
	RepositoryId string
}

func listGitRepositoryBranches(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	repo := h.Item.(git.GitRepository)
	repository_id := d.EqualsQuals["repository_id"].GetStringValue()

	// check if the provided repository_id is not matching with the parentHydrate
	if repository_id != "" && repository_id != repo.Id.String() {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_repository_branch.listGitRepositoryBranches", "connection_error", err)
		return nil, err
	}
	client, err := git.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_repository_branch.listGitRepositoryBranches", "client_error", err)
		return nil, err
	}

	input := git.GetBranchesArgs{
		RepositoryId: types.String(repo.Id.String()),
	}

	branches, err := client.GetBranches(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_repository_branch.listGitRepositoryBranches", "api_error", err)
		return nil, err
	}

	for _, branch := range *branches {
		d.StreamListItem(ctx, Branch{branch, repo.Id.String()})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getRepositoryBranch(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	repositoryId := d.EqualsQuals["repository_id"].GetStringValue()
	branchName := d.EqualsQuals["name"].GetStringValue()

	// Check if repositoryId or branchName is empty
	if repositoryId == "" || branchName == "" {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_repository_branch.getRepositoryBranch", "connection_error", err)
		return nil, err
	}
	client, err := git.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_repository_branch.getRepositoryBranch", "client_error", err)
		return nil, err
	}

	input := git.GetBranchArgs{
		RepositoryId: types.String(repositoryId),
		Name:         types.String(branchName),
	}

	branch, err := client.GetBranch(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_repository_branch.getRepositoryBranch", "api_error", err)
		return nil, err
	}

	return Branch{*branch, repositoryId}, nil
}
