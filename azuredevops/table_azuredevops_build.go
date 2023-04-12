package azuredevops

import (
	"context"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureDevOpsBuild(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_build",
		Description: "Retrieve information about your builds.",
		List: &plugin.ListConfig{
			ParentHydrate: listProjects,
			Hydrate:       listBuilds,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project_id", Require: plugin.Optional},
				{Name: "build_number", Require: plugin.Optional},
				{Name: "reason", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "result", Require: plugin.Optional},
				{Name: "deleted", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "project_id"}),
			Hydrate:    getBuild,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the build.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "build_number",
				Description: "The build number/name of the build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_id",
				Description: "ID of the project this build belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Project.Id"),
			},
			{
				Name:        "quality",
				Description: "The quality of the xaml build (good, bad, etc.).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "build_number_revision",
				Description: "The build number revision.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "deleted",
				Description: "Indicates whether the build has been deleted.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "deleted_date",
				Description: "The date the build was deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DeletedDate.Time"),
			},
			{
				Name:        "deleted_reason",
				Description: "The description of how the build was deleted.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "finish_time",
				Description: "The time that the build was completed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("FinishTime.Time"),
			},
			{
				Name:        "keep_forever",
				Description: "Indicates whether the build should be skipped by retention policies.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_changed_date",
				Description: "The date the build was last changed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastChangedDate.Time"),
			},
			{
				Name:        "parameters",
				Description: "The parameters for the build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "priority",
				Description: "The build's priority.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "queue_options",
				Description: "Additional options for queueing the build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "queue_position",
				Description: "The current position of the build in the queue.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "queue_time",
				Description: "The time that the build was queued.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("QueueTime.Time"),
			},
			{
				Name:        "reason",
				Description: "The reason that the build was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "result",
				Description: "The build result.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "retained_by_release",
				Description: "Indicates whether the build is retained by a release.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "source_branch",
				Description: "The source branch.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_version",
				Description: "The source version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_time",
				Description: "The time that the build was started.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("StartTime.Time"),
			},
			{
				Name:        "uri",
				Description: "The URI of the build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url",
				Description: "The REST URL of the build.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "agent_specification",
				Description: "The agent specification for the build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "controller",
				Description: "The build controller. This is only set if the definition type is Xaml.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "definition",
				Description: "The definition associated with the build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "deleted_by",
				Description: "The identity of the process or person that deleted the build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "demands",
				Description: "A list of demands that represents the agent capabilities required by this build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "last_changed_by",
				Description: "The identity representing the process or person that last changed the build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "links",
				Description: "The class to represent a collection of REST reference links.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "logs",
				Description: "Information about the build logs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "orchestration_plan",
				Description: "The orchestration plan for the build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "plans",
				Description: "Orchestration plans associated with the build (build, cleanup).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "project",
				Description: "The team project.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "properties",
				Description: "The build properties.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "queue",
				Description: "The queue. This is only set if the definition type is Build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "repository",
				Description: "The repository.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "requested_by",
				Description: "The identity that queued the build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "requested_for",
				Description: "The identity on whose behalf the build was queued.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags",
				Description: "The build tags.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "trigger_info",
				Description: "Sourceprovider-specific information about what triggered the build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "triggered_by_build",
				Description: "The build that triggered this build via a Build completion trigger.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "validation_results",
				Description: "Represents the result of validating a build request.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

func listBuilds(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project := h.Item.(core.TeamProjectReference)
	project_id := d.EqualsQuals["project_id"].GetStringValue()

	// check if the provided project_id is not matching with the parentHydrate
	if project_id != "" && project_id != project.Id.String() {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_build.listBuilds", "connection_error", err)
		return nil, err
	}
	client, err := build.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_build.listBuilds", "client_error", err)
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

	input := build.GetBuildsArgs{
		Project: types.String(project.Id.String()),
		Top:     types.Int(maxLimit),
	}
	if d.EqualsQuals["build_number"] != nil {
		input.BuildNumber = types.String(d.EqualsQuals["build_number"].GetStringValue())
	}
	if d.EqualsQuals["reason"] != nil {
		reason := build.BuildReason(d.EqualsQuals["reason"].GetStringValue())
		input.ReasonFilter = &reason
	}
	if d.EqualsQuals["status"] != nil {
		status := build.BuildStatus(d.EqualsQuals["status"].GetStringValue())
		input.StatusFilter = &status
	}
	if d.EqualsQuals["result"] != nil {
		result := build.BuildResult(d.EqualsQuals["result"].GetStringValue())
		input.ResultFilter = &result
	}
	if d.EqualsQuals["id"] != nil {
		id := []int{int(d.EqualsQuals["id"].GetInt64Value())}
		input.BuildIds = &id
	}

	for {
		builds, err := client.GetBuilds(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("azuredevops_build.listBuilds", "api_error", err)
			return nil, err
		}

		for _, build := range builds.Value {
			d.StreamListItem(ctx, build)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		input.ContinuationToken = types.String(builds.ContinuationToken)
		if builds.ContinuationToken == "" {
			break
		}
	}

	return nil, nil
}

func getBuild(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	buildId := d.EqualsQuals["id"].GetInt64Value()
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	// Check if projectId is empty
	if projectId == "" {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_build.getBuild", "connection_error", err)
		return nil, err
	}
	client, err := build.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_build.getBuild", "client_error", err)
		return nil, err
	}

	input := build.GetBuildArgs{
		Project: types.String(projectId),
		BuildId: types.Int(int(buildId)),
	}

	build, err := client.GetBuild(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_build.getBuild", "api_error", err)
		return nil, err
	}

	return build, nil
}
