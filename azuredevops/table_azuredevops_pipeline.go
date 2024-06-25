package azuredevops

import (
	"context"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/pipelines"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureDevOpsPipeline(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_pipeline",
		Description: "Retrieve information about your pipelines.",
		List: &plugin.ListConfig{
			ParentHydrate: listProjects,
			Hydrate:       listPipelines,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project_id", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "project_id"}),
			Hydrate:    getPipeline,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "Pipeline ID.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "name",
				Description: "Pipeline name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_id",
				Description: "ID of the project this pipeline belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "configuration_type",
				Description: "Type of the pipeline configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Configuration.Type"),
			},
			{
				Name:        "folder",
				Description: "Pipeline folder.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "revision",
				Description: "Revision number.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url",
				Description: "URL of the pipeline.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "links",
				Description: "The class to represent a collection of REST reference links.",
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

type Pipeline struct {
	pipelines.Pipeline
	ProjectId string
}

func listPipelines(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project := h.Item.(core.TeamProjectReference)
	project_id := d.EqualsQuals["project_id"].GetStringValue()

	// check if the provided project_id is not matching with the parentHydrate
	if project_id != "" && project_id != project.Id.String() {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_pipeline.listPipelines", "connection_error", err)
		return nil, err
	}
	client := pipelines.NewClient(ctx, connection)

	// Limiting the results
	maxLimit := 1000
	if d.QueryContext.Limit != nil {
		limit := int(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := pipelines.ListPipelinesArgs{
		Project: types.String(project.Id.String()),
		Top:     types.Int(maxLimit),
	}

	pipelines, err := client.ListPipelines(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_pipeline.listPipelines", "api_error", err)
		return nil, err
	}

	for _, pipeline := range *pipelines {
		d.StreamListItem(ctx, Pipeline{pipeline, project.Id.String()})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getPipeline(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pipelineId := d.EqualsQuals["id"].GetInt64Value()
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	// Check if projectId is empty
	if projectId == "" {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_pipeline.getPipeline", "connection_error", err)
		return nil, err
	}
	client := pipelines.NewClient(ctx, connection)

	input := pipelines.GetPipelineArgs{
		Project:    types.String(projectId),
		PipelineId: types.Int(int(pipelineId)),
	}

	pipeline, err := client.GetPipeline(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_pipeline.getPipeline", "api_error", err)
		return nil, err
	}

	return Pipeline{*pipeline, projectId}, nil
}
