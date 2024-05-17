package azuredevops

import (
	"context"
	"os"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func commonColumns(c []*plugin.Column) []*plugin.Column {
	return append([]*plugin.Column{
		{
			Name:        "organization",
			Description: "The name of the organization.",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getOrganization,
			Transform:   transform.FromValue(),
		},
	}, c...)
}

// if the caching is required other than per connection, build a cache key for the call and use it in Memoize.
var getOrganizationMemoized = plugin.HydrateFunc(getOrganizationUncached).Memoize(memoize.WithCacheKeyFunction(getOrganizationCacheKey))

// declare a wrapper hydrate function to call the memoized function
// - this is required when a memoized function is used for a column definition
func getOrganization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getOrganizationMemoized(ctx, d, h)
}

// Build a cache key for the call to getOrganizationCacheKey.
func getOrganizationCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := "getOrganization"
	return key, nil
}

func getOrganizationUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	organizationURL := os.Getenv("AZDO_ORG_SERVICE_URL")
	org := GetConfig(d.Connection)
	if org.OrganizationURL != nil {
		organizationURL = *org.OrganizationURL
	}

	orgName := strings.Split(organizationURL, "https://dev.azure.com/")[1]

	return orgName, nil
}
