package azuredevops

import (
	"context"
	"errors"
	"os"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type azureDevOpsConfig struct {
	OrganizationURL     *string `hcl:"organization_url"`
	PersonalAccessToken *string `hcl:"personal_access_token"`
}

func ConfigInstance() interface{} {
	return &azureDevOpsConfig{}
}

func GetConfig(connection *plugin.Connection) azureDevOpsConfig {
	if connection == nil || connection.Config == nil {
		return azureDevOpsConfig{}
	}

	config, _ := connection.Config.(azureDevOpsConfig)

	return config
}

func getConnection(ctx context.Context, d *plugin.QueryData) (*azuredevops.Connection, error) {
	azureDevOpsConfig := GetConfig(d.Connection)

	organizationURL := os.Getenv("AZDO_ORG_SERVICE_URL")
	personalAccessToken := os.Getenv("AZDO_PERSONAL_ACCESS_TOKEN")

	if azureDevOpsConfig.OrganizationURL != nil {
		organizationURL = *azureDevOpsConfig.OrganizationURL
	}
	if azureDevOpsConfig.PersonalAccessToken != nil {
		personalAccessToken = *azureDevOpsConfig.PersonalAccessToken
	}

	if organizationURL != "" && personalAccessToken != "" {
		connection := azuredevops.NewPatConnection(organizationURL, personalAccessToken)
		return connection, nil
	}

	return nil, errors.New("'organization_url' and 'personal_access_token' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe.")
}
