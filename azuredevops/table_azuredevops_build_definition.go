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

func tableAzureDevOpsBuildDefinition(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_build_definition",
		Description: "Retrieve information about your build definitions.",
		List: &plugin.ListConfig{
			ParentHydrate: listProjects,
			Hydrate:       listBuildDefinitions,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
				{Name: "project_id", Require: plugin.Optional},
				{Name: "repository_id", Require: plugin.Optional},
				{Name: "repository_type", Require: plugin.Optional},
				{Name: "path", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "project_id"}),
			Hydrate:    getBuildDefinition,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the referenced definition.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "name",
				Description: "The name of the referenced definition.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "quality",
				Description: "The quality of the definition document (draft, etc.).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_id",
				Description: "ID of the project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Project.Id"),
			},
			{
				Name:        "repository_id",
				Description: "ID of the repository.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Repository.Id"),
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "badge_enabled",
				Description: "Indicates whether badges are enabled for this definition.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "build_number_format",
				Description: "The build number format.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "comment",
				Description: "A save-time comment for the definition.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "created_date",
				Description: "The date this version of the definition was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreatedDate.Time"),
			},
			{
				Name:        "description",
				Description: "The description.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "drop_location",
				Description: "The drop location for the definition.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "job_authorization_scope",
				Description: "The job authorization scope for builds queued against this definition.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "job_cancel_timeout_in_minutes",
				Description: "The job cancel timeout (in minutes) for builds cancelled by user for this definition.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "job_timeout_in_minutes",
				Description: "The job execution timeout (in minutes) for builds queued against this definition.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "path",
				Description: "The folder path of the definition.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "queue_status",
				Description: "A value that indicates whether builds can be queued against this definition.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository_type",
				Description: "Type of the repository.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Repository.Type"),
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "revision",
				Description: "The definition revision number.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "type",
				Description: "The type of the definition.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uri",
				Description: "The definition's URI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url",
				Description: "The REST URL of the definition.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "authored_by",
				Description: "The author of the definition.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "demands",
				Description: "A list of demands that represents the agent capabilities required by this build.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "draft_of",
				Description: "A reference to the definition that this definition is a draft of, if this is a draft definition.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "drafts",
				Description: "The list of drafts associated with this definition, if this is not a draft definition.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "latest_build",
				Description: "Data representation of a build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "latest_completed_build",
				Description: "Data representation of a latest completed build.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "links",
				Description: "The class to represent a collection of REST reference links.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metrics",
				Description: "Represents metadata about builds in the system.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "options",
				Description: "Represents the application of an optional behavior to a build definition.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "process",
				Description: "The build process.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "process_parameters",
				Description: "The process parameters for this definition.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "project",
				Description: "A reference to the project.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "properties",
				Description: "The class represents a property bag as a collection of key-value pairs. Values of all primitive types (any type with a TypeCode != TypeCode.Object) except for DBNull are accepted. Values of type Byte[], Int32, Double, DateType and String preserve their type, other primitives are returned as a String. Byte[] expected as base64 encoded string.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "queue",
				Description: "The default queue for builds run against this definition.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "repository",
				Description: "The repository.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "retention_rules",
				Description: "Represents a retention policy for a build definition.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "tags",
				Description: "The build definition tags.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "triggers",
				Description: "Represents a trigger for a build definition.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "variable_groups",
				Description: "Represents a variable group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBuildDefinition,
			},
			{
				Name:        "variables",
				Description: "Represents a variable used by a build definition.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBuildDefinition,
			},
		},
	}
}

func listBuildDefinitions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project := h.Item.(core.TeamProjectReference)
	project_id := d.EqualsQuals["project_id"].GetStringValue()

	// check if the provided project_id is not matching with the parentHydrate
	if project_id != "" && project_id != project.Id.String() {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_build_definition.listBuildDefinitions", "connection_error", err)
		return nil, err
	}
	client, err := build.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_build_definition.listBuildDefinitions", "client_error", err)
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

	input := build.GetDefinitionsArgs{
		Project:              types.String(project.Id.String()),
		Top:                  types.Int(maxLimit),
		IncludeAllProperties: types.Bool(true),
		IncludeLatestBuilds:  types.Bool(true),
	}
	if d.EqualsQuals["name"] != nil {
		input.Name = types.String(d.EqualsQuals["name"].GetStringValue())
	}
	if d.EqualsQuals["repository_id"] != nil {
		input.RepositoryId = types.String(d.EqualsQuals["repository_id"].GetStringValue())
	}
	if d.EqualsQuals["repository_type"] != nil {
		input.RepositoryType = types.String(d.EqualsQuals["repository_type"].GetStringValue())
	}
	if d.EqualsQuals["path"] != nil {
		input.Path = types.String(d.EqualsQuals["path"].GetStringValue())
	}
	if d.EqualsQuals["id"] != nil {
		id := []int{int(d.EqualsQuals["id"].GetInt64Value())}
		input.DefinitionIds = &id
	}

	for {
		definitions, err := client.GetDefinitions(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("azuredevops_build_definition.listBuildDefinitions", "api_error", err)
			return nil, err
		}

		for _, build := range definitions.Value {
			d.StreamListItem(ctx, build)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		input.ContinuationToken = types.String(definitions.ContinuationToken)
		if definitions.ContinuationToken == "" {
			break
		}
	}

	return nil, nil
}

func getBuildDefinition(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var projectId string
	var definitionId int
	if h.Item != nil {
		definitionId = *h.Item.(build.BuildDefinitionReference).Id
		projectId = h.Item.(build.BuildDefinitionReference).Project.Id.String()
	} else {
		definitionId = int(d.EqualsQuals["id"].GetInt64Value())
		projectId = d.EqualsQuals["project_id"].GetStringValue()
	}

	// Check if projectId is empty
	if projectId == "" {
		return nil, nil
	}

	connection, err := getConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_build_definition.getBuildDefinition", "connection_error", err)
		return nil, err
	}
	client, err := build.NewClient(ctx, connection)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_build_definition.getBuildDefinition", "client_error", err)
		return nil, err
	}

	input := build.GetDefinitionArgs{
		Project:      types.String(projectId),
		DefinitionId: types.Int(int(definitionId)),
	}

	definition, err := client.GetDefinition(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_build_definition.getBuildDefinition", "api_error", err)
		return nil, err
	}

	return definition, nil
}
