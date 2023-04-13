package azuredevops

import (
	"context"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/identity"
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
		Columns: []*plugin.Column{
			{
				Name:        "principal_name",
				Description: "This is the PrincipalName of this graph member from the source provider. The source provider may change this field over time and it is not guaranteed to be immutable for the life of the graph member by VSTS.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provider_display_name",
				Description: "This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.",
				Type:        proto.ColumnType_STRING,
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

func listGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.listGroups", "connection_error", err)
		return nil, err
	}

	client, err := identity.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.listGroups", "client_error", err)
		return nil, err
	}

	input := identity.ListGroupsArgs{}

	// for {
	groups, err := client.ListGroups(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_group.listGroups", "api_error", err)
		return nil, err
	}

	for _, group := range *groups {
		d.StreamListItem(ctx, group)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	// 	continuationToken := *groups.ContinuationToken
	// 	if continuationToken[0] == "" {
	// 		break
	// 	}
	// 	input.ContinuationToken = types.String(continuationToken[0])
	// }

	return nil, nil
}
