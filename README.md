![image](https://hub.steampipe.io/images/plugins/turbot/azuredevops-social-graphic.png)

# Azure DevOps Plugin for Steampipe

Use SQL to query projects, groups, builds and more from Azure DevOps.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/azuredevops)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/azuredevops/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-azuredevops/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install azuredevops
```

Configure the Organization URL and Personal Access Token in `~/.steampipe/config/azuredevops.spc`:

```hcl
connection "azuredevops" {
  plugin = "azuredevops"

  # Azure DevOps Organization URL
  # For more information on the Organization URL, please see https://learn.microsoft.com/en-us/azure/devops/extend/develop/work-with-urls?view=azure-devops&tabs=http.
  # Can also be set with the AZDO_ORG_SERVICE_URL environment variable.
  # organization_url = "https://dev.azure.com/test"

  # Azure DevOps Personal Access Token
  # For more information on the Personal Access Token, please see https://learn.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate?view=azure-devops&tabs=Windows.
  # Can also be set with the AZDO_PERSONAL_ACCESS_TOKEN environment variable.
  # personal_access_token = "wf3hahidy7i7fkzmeqr3e6fbjwuspabpo766grp7hl4o65v2"
}
```

Run steampipe:

```shell
steampipe query
```

Query your project:

```sql
select
  id,
  name,
  state,
  visibility
from
  azuredevops_project;
```

```
+--------------------------------------+-----------------+------------+------------+
| id                                   | name            | state      | visibility |
+--------------------------------------+-----------------+------------+------------+
| bdcdf70b-7757-4253-b36c-33c08ca07dbb | test-1          | wellFormed | private    |
| 9a1f26ce-c715-4ef4-b557-503fdb6be55a | private_project | wellFormed | private    |
+--------------------------------------+-----------------+------------+------------+
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-azuredevops.git
cd steampipe-plugin-azuredevops
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/azuredevops.spc
```

Try it!

```
steampipe query
> .inspect azuredevops
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-azuredevops/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Azure DevOps Plugin](https://github.com/turbot/steampipe-plugin-azuredevops/labels/help%20wanted)
