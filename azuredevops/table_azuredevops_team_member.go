package azuredevops

import (
	"context"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/webapi"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureDevOpsTeamMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_team_member",
		Description: "Retrieve information about your team members.",
		List: &plugin.ListConfig{
			ParentHydrate: listTeams,
			Hydrate:       listTeamMembers,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The member id.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_team_admin",
				Description: "Check if the member is the team admin.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "project_id",
				Description: "The project id.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "team_id",
				Description: "The project id.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "descriptor",
				Description: "The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url",
				Description: "This url is the full route to the source resource of this graph subject.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_deleted_in_origin",
				Description: "Check if the member is already deleted.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "links",
				Description: "This field contains zero or more interesting links about the graph subject. These links may be invoked to obtain additional relationships or more detailed information about this graph subject.",
				Type:        proto.ColumnType_JSON,
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
		},
	}
}

type Member struct {
	webapi.IdentityRef
	IsTeamAdmin *bool
	ProjectId   string
	TeamId      string
}

func listTeamMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team := h.Item.(core.WebApiTeam)

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_team_member.listTeamMembers", "connection_error", err)
		return nil, err
	}

	client, err := core.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_team_member.listTeamMembers", "client_error", err)
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

	input := core.GetTeamMembersWithExtendedPropertiesArgs{
		ProjectId: types.String(team.ProjectId.String()),
		TeamId:    types.String(team.Id.String()),
		Top:       types.Int(maxLimit),
	}

	members, err := client.GetTeamMembersWithExtendedProperties(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_team_member.listTeamMembers", "api_error", err)
		return nil, err
	}

	for _, member := range *members {
		d.StreamListItem(ctx, Member{*member.Identity, member.IsTeamAdmin, team.ProjectId.String(), team.Id.String()})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
