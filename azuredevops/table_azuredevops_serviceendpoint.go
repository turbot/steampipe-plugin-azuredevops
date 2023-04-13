package azuredevops

import (
	"context"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureDevOpsServiceEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_serviceendpoint",
		Description: "Retrieve information about your service endpoints.",
		List: &plugin.ListConfig{
			ParentHydrate: listProjects,
			Hydrate:       listServiceEndpoints,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project_id", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
				{Name: "owner", Require: plugin.Optional},
				{Name: "id", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "project_id"}),
			Hydrate:    getServiceEndpoint,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "Gets the identifier of this endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "Gets the friendly name of the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_ready",
				Description: "EndPoint state indicator.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_shared",
				Description: "Indicates whether service endpoint is shared with other projects or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "project_id",
				Description: "ID of the project this service endpoint belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Gets the description of endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner",
				Description: "Owner of the endpoint. Supported values are library and agentcloud.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Gets the type of the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url",
				Description: "Gets the url of the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "administrators_group",
				Description: "Gets the identity reference for the administrators group of the service endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "authorization",
				Description: "Gets the authorization data for talking to the endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "created_by",
				Description: "Gets the identity reference for the user who created the Service endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "data",
				Description: "The service endpoint data.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "operation_status",
				Description: "Error message during creation/deletion of endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "readers_group",
				Description: "Gets the identity reference for the readers group of the service endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "service_endpoint_project_references",
				Description: "All other project references where the service endpoint is shared.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

type ServiceEndpoint struct {
	serviceendpoint.ServiceEndpoint
	ProjectId string
}

func listServiceEndpoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project := h.Item.(core.TeamProjectReference)
	project_id := d.EqualsQuals["project_id"].GetStringValue()

	// check if the provided project_id is not matching with the parentHydrate
	if project_id != "" && project_id != project.Id.String() {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_serviceendpoint.listServiceEndpoints", "connection_error", err)
		return nil, err
	}
	client, err := serviceendpoint.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_serviceendpoint.listServiceEndpoints", "client_error", err)
		return nil, err
	}

	input := serviceendpoint.GetServiceEndpointsArgs{
		Project:        types.String(project.Id.String()),
		IncludeFailed:  types.Bool(true),
		IncludeDetails: types.Bool(true),
	}
	if d.EqualsQuals["type"] != nil {
		input.Type = types.String(d.EqualsQuals["type"].GetStringValue())
	}
	if d.EqualsQuals["owner"] != nil {
		input.Owner = types.String(d.EqualsQuals["owner"].GetStringValue())
	}
	if d.EqualsQuals["id"] != nil {
		endpointId, _ := uuid.Parse(d.EqualsQuals["id"].GetStringValue())
		input.EndpointIds = &[]uuid.UUID{endpointId}
	}

	serviceEndpoints, err := client.GetServiceEndpoints(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_serviceendpoint.listServiceEndpoints", "api_error", err)
		return nil, err
	}

	for _, serviceEndpoint := range *serviceEndpoints {
		d.StreamListItem(ctx, ServiceEndpoint{serviceEndpoint, project.Id.String()})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getServiceEndpoint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	// Check if projectId or id is empty
	if projectId == "" || id == "" {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_repository_branch.getRepositoryBranch", "connection_error", err)
		return nil, err
	}
	client, err := serviceendpoint.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_repository_branch.getRepositoryBranch", "client_error", err)
		return nil, err
	}

	endpointId, _ := uuid.Parse(id)
	input := serviceendpoint.GetServiceEndpointDetailsArgs{
		Project:    types.String(projectId),
		EndpointId: &endpointId,
	}

	serviceEndpoint, err := client.GetServiceEndpointDetails(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_repository_branch.getRepositoryBranch", "api_error", err)
		return nil, err
	}

	return ServiceEndpoint{*serviceEndpoint, projectId}, nil
}
