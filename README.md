![image](https://hub.steampipe.io/images/plugins/turbot/azuredevops-social-graphic.png)

# Azure DevOps Plugin for Steampipe

Use SQL to query projects, groups, builds and more from Azure DevOps.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/azuredevops)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/azuredevops/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-azuredevops/issues)

## Quick start

Download and install the latest Azure DevOps plugin:

```bash
steampipe plugin install azuredevops
```

Configure your [credentials](https://hub.steampipe.io/plugins/turbot/azuredevops#credentials) and [config file](https://hub.steampipe.io/plugins/turbot/azuredevops#configuration).

Configure the Organization URL and Personal Access Token in `~/.steampipe/config/azuredevops.spc`:

```hcl
connection "azuredevops" {
  plugin = "azuredevops"
  # Authentication information
  organization_url = "https://dev.azure.com/test"
  personal_access_token = "wf3hahidy7i7fkzmeqr3e6fbjwuspabpo766grp7hl4o65v2"
}
```

Or through environment variables:

```sh
export AZDO_ORG_SERVICE_URL=https://dev.azure.com/test
export AZDO_PERSONAL_ACCESS_TOKEN=wf3hahidy7i7fkzmeqr3e6fbjwuspabpo766grp7hl4o65v2
```

Run steampipe:

```shell
steampipe query
```

List your Azure DevOps projects:

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

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). Contributions to the plugin are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-azuredevops/blob/main/LICENSE). Contributions to the plugin documentation are subject to the [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-azuredevops/blob/main/docs/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Azure DevOps Plugin](https://github.com/turbot/steampipe-plugin-azuredevops/labels/help%20wanted)
