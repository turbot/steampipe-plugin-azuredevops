package azuredevops

import (
	"context"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/accounts"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureDevOpsAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_account",
		Description: "Retrieve information about your accounts.",
		List: &plugin.ListConfig{
			Hydrate:    listAccounts,
			KeyColumns: plugin.AnyColumn([]string{"owner_id", "member_id"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "account_id",
				Description: "Identifier for an account",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_name",
				Description: "Name for an account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "Name for an account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("owner_id"),
			},
			{
				Name:        "member_id",
				Description: "Name for an account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("member_id"),
			},
			{
				Name:        "account_owner",
				Description: "Name for an account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_status",
				Description: "Current account status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_type",
				Description: "Type of account: Personal, Organization.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_uri",
				Description: "Uri for an account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_by",
				Description: "Who created the account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_date",
				Description: "Date account was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreatedDate.Time"),
			},
			{
				Name:        "has_moved",
				Description: "Check if the account has moved.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_updated_by",
				Description: "Identity of last person to update the account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_date",
				Description: "Date account was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastUpdatedDate.Time"),
			},
			{
				Name:        "namespace_id",
				Description: "Namespace for an account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "new_collection_id",
				Description: "New collection for an account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization_name",
				Description: "Organization that created the account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_reason",
				Description: "Reason for current status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "properties",
				Description: "Extended properties.",
				Type:        proto.ColumnType_JSON,
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountName"),
			},
		},
	}
}

func listAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_account.listAccounts", "connection_error", err)
		return nil, err
	}

	client, err := accounts.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_account.listAccounts", "client_error", err)
		return nil, err
	}

	input := accounts.GetAccountsArgs{}

	if d.EqualsQuals["owner_id"] != nil {
		ownerId, _ := uuid.Parse(d.EqualsQuals["owner_id"].GetStringValue())
		input.OwnerId = &ownerId
	} else if d.EqualsQuals["member_id"] != nil {
		memberId, _ := uuid.Parse(d.EqualsQuals["member_id"].GetStringValue())
		input.MemberId = &memberId
	}

	accounts, err := client.GetAccounts(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_account.listAccounts", "api_error", err)
		return nil, err
	}

	for _, account := range *accounts {
		d.StreamListItem(ctx, account)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
