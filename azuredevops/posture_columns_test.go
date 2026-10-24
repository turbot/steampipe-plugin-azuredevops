package azuredevops

import (
	"context"
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Every column must resolve to a field of the row its hydrate returns; a misspelt one reads as NULL
// on every row instead of failing.
func TestEveryColumnResolvesToARowField(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		table    *plugin.Table
		row      interface{}
		hydrated map[string]interface{}
	}{
		{tableAzureDevOpsGitBranchPolicy(ctx), GitBranchPolicy{}, nil},
		{tableAzureDevOpsProjectPipelineSettings(ctx), ProjectPipelineSettings{}, nil},
		{tableAzureDevOpsServiceEndpoint(ctx), ServiceEndpoint{}, map[string]interface{}{
			"all_pipelines_authorized": pipelinePermissions{}, "authorized_pipeline_ids": pipelinePermissions{},
		}},
		{tableAzureDevOpsVariableGroup(ctx), VariableGroup{}, nil},
		{tableAzureDevOpsUserEntitlement(ctx), UserEntitlement{}, nil},
		{tableAzureDevOpsAuditStream(ctx), AuditStream{}, nil},
		{tableAzureDevOpsAdvancedSecurityEnablement(ctx), AdvancedSecurityEnablement{}, nil},
		{tableAzureDevOpsAdvancedSecurityAlert(ctx), AdvancedSecurityAlert{}, nil},
	} {
		for _, col := range tc.table.Columns {
			if col.Name == "organization" {
				continue
			}
			row := tc.row
			if h, ok := tc.hydrated[col.Name]; ok {
				row = h
			} else if col.Hydrate != nil {
				t.Errorf("%s.%s: hydrated column missing from the test", tc.table.Name, col.Name)
				continue
			}
			field := strcase.ToCamel(col.Name)
			if col.Transform != nil {
				names, isField := col.Transform.Transforms[0].Param.([]string)
				if !isField {
					continue
				}
				field = names[0]
			}
			if _, ok := helpers.GetNestedFieldValueFromInterface(row, field); !ok {
				t.Errorf("%s.%s: %T has no field %s", tc.table.Name, col.Name, row, field)
			}
		}
	}
}
