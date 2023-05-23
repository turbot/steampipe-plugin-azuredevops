package azuredevops

import (
	"context"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/dashboard"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureDevOpsDashboard(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_dashboard",
		Description: "Retrieve information about your dashboards.",
		List: &plugin.ListConfig{
			ParentHydrate: listProjects,
			Hydrate:       listDashboards,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project_id", Require: plugin.Optional},
				{Name: "group_id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "ID of the dashboard. Provided by service at creation time.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "Name of the Dashboard.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dashboard_scope",
				Description: "Entity to which the dashboard is scoped.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "group_id",
				Description: "ID of the group for a dashboard. For team-scoped dashboards, this is the unique identifier for the team associated with the dashboard. For project-scoped dashboards this property is empty.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "ID of the owner for a dashboard. For team-scoped dashboards, this is the unique identifier for the team associated with the dashboard. For project-scoped dashboards, this is the unique identifier for the user identity associated with the dashboard.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_id",
				Description: "ID of the project this dashboard belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Description of the dashboard.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "Server defined version tracking value, used for edit collision detection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ETag"),
			},
			{
				Name:        "position",
				Description: "Position of the dashboard, within a dashboard group. If unset at creation time, position is decided by the service.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "refresh_interval",
				Description: "Interval for client to automatically refresh the dashboard. Expressed in minutes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "url",
				Description: "URL of the dashboard.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "links",
				Description: "The class to represent a collection of REST reference links.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "widgets",
				Description: "The set of Widgets on the dashboard.",
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

type Dashboard struct {
	dashboard.Dashboard
	ProjectId string
}

func listDashboards(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project := h.Item.(core.TeamProjectReference)
	project_id := d.EqualsQuals["project_id"].GetStringValue()

	// check if the provided project_id is not matching with the parentHydrate
	if project_id != "" && project_id != project.Id.String() {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_dashboard.listDashboards", "connection_error", err)
		return nil, err
	}
	client, err := dashboard.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_dashboard.listDashboards", "client_error", err)
		return nil, err
	}

	input := dashboard.GetDashboardsByProjectArgs{
		Project: types.String(project.Id.String()),
	}
	if d.EqualsQuals["group_id"] != nil {
		input.Team = types.String(d.EqualsQuals["group_id"].GetStringValue())
	}

	dashboards, err := client.GetDashboardsByProject(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_dashboard.listDashboards", "api_error", err)
		return nil, err
	}

	for _, dashboard := range *dashboards {
		d.StreamListItem(ctx, Dashboard{dashboard, project.Id.String()})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
