package azuredevops

import (
	"context"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureDevOpsTeam(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_team",
		Description: "Retrieve information about your teams.",
		List: &plugin.ListConfig{
			Hydrate: listTeams,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "project_id"}),
			Hydrate:    getTeam,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "Team (Identity) Guid. A Team Foundation ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "Team name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_id",
				Description: "The project id.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_name",
				Description: "The project name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Team description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "identity_url",
				Description: "Identity REST API Url to this team.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url",
				Description: "Team REST API Url.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "identity",
				Description: "Team identity.",
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

func listTeams(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_team.listTeams", "connection_error", err)
		return nil, err
	}

	client, err := core.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_team.listTeams", "client_error", err)
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

	input := core.GetAllTeamsArgs{
		Top:            types.Int(maxLimit),
		ExpandIdentity: types.Bool(true),
	}

	teams, err := client.GetAllTeams(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_team.listTeams", "api_error", err)
		return nil, err
	}

	for _, team := range *teams {
		d.StreamListItem(ctx, team)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getTeam(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	teamId := d.EqualsQuals["id"].GetStringValue()
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	// Check if projectId or teamId is empty
	if projectId == "" || teamId == "" {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_team.getTeam", "connection_error", err)
		return nil, err
	}

	client, err := core.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_team.getTeam", "client_error", err)
		return nil, err
	}

	input := core.GetTeamArgs{
		ProjectId:      types.String(projectId),
		TeamId:         types.String(teamId),
		ExpandIdentity: types.Bool(true),
	}

	team, err := client.GetTeam(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_team.getTeam", "api_error", err)
		return nil, err
	}

	return team, nil
}
