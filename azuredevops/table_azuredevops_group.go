package azuredevops

import (
	"context"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/graph"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureDevOpsGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_group",
		Description: "Retrieve information about your groups.",
		List: &plugin.ListConfig{
			Hydrate: listGroups,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("descriptor"),
			Hydrate:    getGroup,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "principal_name",
				Description: "This is the PrincipalName of this graph member from the source provider. The source provider may change this field over time and it is not guaranteed to be immutable for the life of the graph member by VSTS.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "membership_state",
				Description: "When true, the membership is active.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getGroupMembershipState,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "domain",
				Description: "This represents the name of the container of origin for a graph member. (For MSA this is Windows Live ID, for AD the name of the domain, for AAD the tenantID of the directory, for VSTS groups the ScopeId, etc).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "origin",
				Description: "The type of source provider for the origin identifier (ex:AD, AAD, MSA).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A short phrase to help human readers disambiguate groups with similar names.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "descriptor",
				Description: "The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "legacy_descriptor",
				Description: "The legacy descriptor is here in case you need to access old version IMS using identity descriptor.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mail_address",
				Description: "The email address of record for a given graph member. This may be different than the principal name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "origin_id",
				Description: "The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subject_kind",
				Description: "This field identifies the type of the graph subject (ex: Group, Scope, User).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url",
				Description: "This url is the full route to the source resource of this graph subject.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "links",
				Description: "This field contains zero or more interesting links about the graph subject. These links may be invoked to obtain additional relationships or more detailed information about this graph subject.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "memberships",
				Description: "Get all the memberships where this descriptor is a member in the relationship.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGroupMemberships,
				Transform:   transform.FromValue(),
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PrincipalName"),
			},
		}),
	}
}

func listGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.listGroups", "connection_error", err)
		return nil, err
	}

	client, err := graph.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.listGroups", "client_error", err)
		return nil, err
	}

	input := graph.ListGroupsArgs{}

	for {
		groups, err := client.ListGroups(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("azuredevops_group.listGroups", "api_error", err)
			return nil, err
		}

		for _, group := range *groups.GraphGroups {
			d.StreamListItem(ctx, group)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		continuationToken := *groups.ContinuationToken
		if continuationToken[0] == "" {
			break
		}
		input.ContinuationToken = types.String(continuationToken[0])
	}

	return nil, nil
}

func getGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	descriptor := d.EqualsQuals["descriptor"].GetStringValue()

	// Check if descriptor is empty
	if descriptor == "" {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.getGroup", "connection_error", err)
		return nil, err
	}

	client, err := graph.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.getGroup", "client_error", err)
		return nil, err
	}

	input := graph.GetGroupArgs{
		GroupDescriptor: types.String(descriptor),
	}

	output, err := client.GetGroup(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.getGroup", "api_error", err)
		return nil, err
	}

	return *output, nil
}

func getGroupMembershipState(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	group := h.Item.(graph.GraphGroup)
	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.getGroupMembershipState", "connection_error", err)
		return nil, err
	}

	client, err := graph.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.getGroupMembershipState", "client_error", err)
		return nil, err
	}

	input := graph.GetMembershipStateArgs{
		SubjectDescriptor: group.Descriptor,
	}

	output, err := client.GetMembershipState(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.getGroupMembershipState", "api_error", err)
		return nil, err
	}

	return output.Active, nil
}

func getGroupMemberships(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	group := h.Item.(graph.GraphGroup)
	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.getGroupMemberships", "connection_error", err)
		return nil, err
	}

	client, err := graph.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.getGroupMemberships", "client_error", err)
		return nil, err
	}

	input := graph.ListMembershipsArgs{
		SubjectDescriptor: group.Descriptor,
	}

	output, err := client.ListMemberships(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.getGroupMemberships", "api_error", err)
		return nil, err
	}

	return output, nil
}
