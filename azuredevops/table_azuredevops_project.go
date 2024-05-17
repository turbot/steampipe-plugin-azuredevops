package azuredevops

import (
	"context"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureDevOpsProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_project",
		Description: "Retrieve information about your projects.",
		List: &plugin.ListConfig{
			Hydrate: listProjects,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "state",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getProject,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "Project identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "Project name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "Project state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "visibility",
				Description: "Project visibility.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "abbreviation",
				Description: "Project abbreviation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The project's description (if any).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_update_time",
				Description: "Project last update time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastUpdateTime.Time"),
			},
			{
				Name:        "revision",
				Description: "Project revision.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "url",
				Description: "Url to the full version of the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capabilities",
				Description: "Set of capabilities this project has (such as process template & version control).",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getProject,
			},
			{
				Name:        "default_team",
				Description: "The shallow ref to the default team.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getProject,
			},
			{
				Name:        "links",
				Description: "The links to other objects related to this object.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getProject,
			},
			{
				Name:        "properties",
				Description: "Get a collection of team project properties.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     GetProjectProperties,
				Transform:   transform.FromValue(),
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

func listProjects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_project.listProjects", "connection_error", err)
		return nil, err
	}

	client, err := core.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_project.listProjects", "client_error", err)
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

	input := core.GetProjectsArgs{
		Top: types.Int(maxLimit),
	}

	if d.EqualsQuals["state"] != nil {
		state := core.ProjectState(d.EqualsQuals["state"].GetStringValue())
		input.StateFilter = &state
	}

	for {
		response, err := client.GetProjects(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("azuredevops_project.listProjects", "api_error", err)
			return nil, err
		}

		for _, project := range response.Value {
			d.StreamListItem(ctx, project)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		input.ContinuationToken = types.String(response.ContinuationToken)
		if response.ContinuationToken == "" {
			break
		}
	}

	return nil, nil
}

func getProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		id = h.Item.(core.TeamProjectReference).Id.String()
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
	}

	// Check if id is empty
	if id == "" {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_project.getProject", "connection_error", err)
		return nil, err
	}

	client, err := core.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_project.getProject", "client_error", err)
		return nil, err
	}

	input := core.GetProjectArgs{
		ProjectId:           types.String(id),
		IncludeCapabilities: types.Bool(true),
	}

	project, err := client.GetProject(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_project.getProject", "api_error", err)
		return nil, err
	}

	return project, nil
}

func GetProjectProperties(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectId := getProjectId(h.Item)

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_project.GetProjectProperties", "connection_error", err)
		return nil, err
	}

	client, err := core.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_project.GetProjectProperties", "client_error", err)
		return nil, err
	}

	input := core.GetProjectPropertiesArgs{
		ProjectId: projectId,
	}

	properties, err := client.GetProjectProperties(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_project.GetProjectProperties", "api_error", err)
		return nil, err
	}

	return properties, nil
}

func getProjectId(item interface{}) *uuid.UUID {
	switch item := item.(type) {
	case core.TeamProjectReference:
		return item.Id
	case *core.TeamProject:
		return item.Id
	}

	return nil
}
