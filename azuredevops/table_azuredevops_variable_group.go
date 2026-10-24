package azuredevops

import (
	"context"
	"regexp"
	"sort"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const variableGroupsAPIVersion = "7.1"

// https://learn.microsoft.com/en-us/rest/api/azure/devops/distributedtask/variablegroups/get-variable-groups
func tableAzureDevOpsVariableGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_variable_group",
		Description: "Variable groups in an Azure DevOps project, with counts of secret variables and the names of plain variables that look like secrets. Variable values are never read.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project", Require: plugin.Required},
			},
			Hydrate: listVariableGroups,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "project", Description: "Name or ID of the project, as given in the qualifier.", Type: proto.ColumnType_STRING, Transform: transform.FromQual("project")},
			{Name: "description", Description: "Description of the variable group.", Type: proto.ColumnType_STRING},
			{Name: "id", Description: "ID of the variable group.", Type: proto.ColumnType_INT},
			{Name: "is_shared", Description: "True if the variable group is shared with other projects.", Type: proto.ColumnType_BOOL},
			{Name: "key_vault_linked", Description: "True if the variable group reads its variables from an Azure key vault.", Type: proto.ColumnType_BOOL},
			{Name: "key_vault_name", Description: "Name of the linked key vault.", Type: proto.ColumnType_STRING},
			{Name: "name", Description: "Name of the variable group.", Type: proto.ColumnType_STRING},
			{Name: "plaintext_variable_names_suspected", Description: "Names of variables not marked secret whose names look like credentials (password, token, key...). Empty for key vault groups.", Type: proto.ColumnType_JSON},
			{Name: "secret_variable_count", Description: "Number of variables marked secret.", Type: proto.ColumnType_INT},
			{Name: "type", Description: "Vsts for variables stored in Azure DevOps, AzureKeyVault for a group linked to a key vault.", Type: proto.ColumnType_STRING},
			{Name: "variable_count", Description: "Number of variables in the group.", Type: proto.ColumnType_INT},
			{Name: "variable_names", Description: "Names of all variables in the group, sorted.", Type: proto.ColumnType_JSON},
			{Name: "title", Description: "Title of the resource.", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name")},
		}),
	}
}

// variableGroupResponse leaves out every variable's value: only names and the secret flag are decoded.
type variableGroupResponse struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	IsShared     bool   `json:"isShared"`
	ProviderData *struct {
		Vault string `json:"vault"`
	} `json:"providerData"`
	Variables map[string]struct {
		IsSecret bool `json:"isSecret"`
	} `json:"variables"`
}

type VariableGroup struct {
	Id                              int
	Name                            string
	Description                     string
	Type                            string
	IsShared                        bool
	KeyVaultLinked                  bool
	KeyVaultName                    *string
	VariableCount                   int
	SecretVariableCount             int
	VariableNames                   []string
	PlaintextVariableNamesSuspected []string
}

func listVariableGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	project := d.EqualsQualString("project")
	if project == "" {
		return nil, nil
	}
	client, err := restClientFor(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_variable_group.listVariableGroups", "connection_error", err)
		return nil, err
	}
	err = readVariableGroups(ctx, client, project, func(g VariableGroup) bool {
		d.StreamListItem(ctx, g)
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_variable_group.listVariableGroups", "api_error", err)
		return nil, err
	}
	return nil, nil
}

func readVariableGroups(ctx context.Context, client *restClient, project string, emit func(VariableGroup) bool) error {
	params := apiVersion(variableGroupsAPIVersion)
	return getList(ctx, client, projectURL(client.baseURL, project, "distributedtask/variablegroups"), params, scopeVariableGroups, func(g variableGroupResponse) bool {
		return emit(summarizeVariableGroup(g))
	})
}

func summarizeVariableGroup(g variableGroupResponse) VariableGroup {
	out := VariableGroup{
		Id: g.Id, Name: g.Name, Description: g.Description, Type: g.Type, IsShared: g.IsShared,
		KeyVaultLinked:                  strings.EqualFold(g.Type, "AzureKeyVault"),
		VariableNames:                   []string{},
		PlaintextVariableNamesSuspected: []string{},
	}
	if out.KeyVaultLinked && g.ProviderData != nil && g.ProviderData.Vault != "" {
		vault := g.ProviderData.Vault
		out.KeyVaultName = &vault
	}
	for name, v := range g.Variables {
		out.VariableCount++
		out.VariableNames = append(out.VariableNames, name)
		if v.IsSecret {
			out.SecretVariableCount++
		} else if !out.KeyVaultLinked && looksLikeSecretName(name) {
			out.PlaintextVariableNamesSuspected = append(out.PlaintextVariableNamesSuspected, name)
		}
	}
	sort.Strings(out.VariableNames)
	sort.Strings(out.PlaintextVariableNamesSuspected)
	return out
}

var (
	// Matched against the name lowercased with separators removed, so DB_PASSWORD and dbPassword agree.
	secretNameFragments = []string{
		"password", "passwd", "passphrase", "secret", "token", "apikey", "accesskey", "privatekey",
		"connectionstring", "connstring", "credential", "sastoken", "saskey", "accountkey", "signingkey",
	}
	// Short words that only count as a whole word: "pat" must not match "path".
	secretNameWords  = map[string]bool{"pwd": true, "pw": true, "pat": true, "pass": true, "sas": true}
	nameWordBoundary = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	nameSeparators   = regexp.MustCompile(`[^a-zA-Z0-9]+`)
)

func looksLikeSecretName(name string) bool {
	squashed := strings.ToLower(nameSeparators.ReplaceAllString(name, ""))
	for _, f := range secretNameFragments {
		if strings.Contains(squashed, f) {
			return true
		}
	}
	words := strings.Fields(strings.ToLower(nameSeparators.ReplaceAllString(nameWordBoundary.ReplaceAllString(name, "$1 $2"), " ")))
	for _, w := range words {
		if secretNameWords[w] {
			return true
		}
	}
	if n := len(words); n > 1 && (words[n-1] == "key" || words[n-1] == "keys") {
		return true
	}
	return false
}
