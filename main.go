package main

import (
	"github.com/turbot/steampipe-plugin-azuredevops/azuredevops"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: azuredevops.Plugin,
	})
}
