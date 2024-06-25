package azuredevops

import (
	"context"
	"strconv"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/release"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureDevOpsRelease(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_release",
		Description: "Retrieve information about your releases.",
		List: &plugin.ListConfig{
			Hydrate: listReleases,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project_id", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "id", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "project_id"}),
			Hydrate:    getRelease,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "Gets the unique identifier of this release.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "name",
				Description: "The release name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The release status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_id",
				Description: "ID of the project this release belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectReference.Id"),
			},
			{
				Name:        "comment",
				Description: "Gets comment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_on",
				Description: "Gets date on which it got created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreatedOn.Time"),
			},
			{
				Name:        "definition_snapshot_revision",
				Description: "Gets revision number of definition snapshot.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "Gets description of release.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "keep_forever",
				Description: "Whether to exclude the release from retention policies.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "logs_container_url",
				Description: "Gets logs container url.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "modified_on",
				Description: "Gets date on which it got modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ModifiedOn.Time"),
			},
			{
				Name:        "pool_name",
				Description: "Gets pool name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reason",
				Description: "Gets reason of release.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "release_definition_revision",
				Description: "Gets the release definition revision.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "release_name_format",
				Description: "Gets release name format.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "triggering_artifact_alias",
				Description: "Gets triggering artifact alias.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "artifacts",
				Description: "Gets the list of artifacts.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "created_by",
				Description: "Gets the identity who created.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "created_for",
				Description: "Gets the identity for whom release was created.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "environments",
				Description: "Gets list of environments.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "links",
				Description: "The class to represent a collection of REST reference links.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "modified_by",
				Description: "Gets the identity who modified.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "project_reference",
				Description: "Gets project reference.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "properties",
				Description: "The release properties.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "release_definition",
				Description: "Gets releaseDefinitionReference which specifies the reference of the release definition to which this release is associated.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags",
				Description: "Gets list of tags.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "variable_groups",
				Description: "Gets the list of variable groups.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "variables",
				Description: "Gets or sets the dictionary of variables.",
				Type:        proto.ColumnType_JSON,
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

func listReleases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_release.listReleases", "connection_error", err)
		return nil, err
	}
	client, err := release.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_release.listReleases", "client_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := 1000
	if d.QueryContext.Limit != nil {
		limit := int(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := release.GetReleasesArgs{
		Top: types.Int(maxLimit),
	}
	if d.EqualsQuals["project_id"] != nil {
		input.Project = types.String(d.EqualsQuals["project_id"].GetStringValue())
	}
	if d.EqualsQuals["name"] != nil {
		input.SearchText = types.String(d.EqualsQuals["name"].GetStringValue())
	}
	if d.EqualsQuals["status"] != nil {
		status := release.ReleaseStatus(d.EqualsQuals["status"].GetStringValue())
		input.StatusFilter = &status
	}
	if d.EqualsQuals["id"] != nil {
		id := []int{int(d.EqualsQuals["id"].GetInt64Value())}
		input.ReleaseIdFilter = &id
	}

	for {
		releases, err := client.GetReleases(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("azuredevops_release.listReleases", "api_error", err)
			return nil, err
		}

		for _, release := range releases.Value {
			d.StreamListItem(ctx, release)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		token, _ := strconv.Atoi(releases.ContinuationToken)
		input.ContinuationToken = &token
		if releases.ContinuationToken == "" {
			break
		}
	}

	return nil, nil
}

func getRelease(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	releaseId := d.EqualsQuals["id"].GetInt64Value()
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	// Check if projectId is empty
	if projectId == "" {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_release.getRelease", "connection_error", err)
		return nil, err
	}
	client, err := release.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_release.getRelease", "client_error", err)
		return nil, err
	}

	input := release.GetReleaseArgs{
		Project:   types.String(projectId),
		ReleaseId: types.Int(int(releaseId)),
	}

	release, err := client.GetRelease(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_release.getRelease", "api_error", err)
		return nil, err
	}

	return release, nil
}
