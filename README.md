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

## Engines

This plugin is available for the following engines:

| Engine        | Description
|---------------|------------------------------------------
| [Steampipe](https://steampipe.io/docs) | The Steampipe CLI exposes APIs and services as a high-performance relational database, giving you the ability to write SQL-based queries to explore dynamic data. Mods extend Steampipe's capabilities with dashboards, reports, and controls built with simple HCL. The Steampipe CLI is a turnkey solution that includes its own Postgres database, plugin management, and mod support.
| [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview) | Steampipe Postgres FDWs are native Postgres Foreign Data Wrappers that translate APIs to foreign tables. Unlike Steampipe CLI, which ships with its own Postgres server instance, the Steampipe Postgres FDWs can be installed in any supported Postgres database version.
| [SQLite Extension](https://steampipe.io/docs/steampipe_sqlite/overview) | Steampipe SQLite Extensions provide SQLite virtual tables that translate your queries into API calls, transparently fetching information from your API or service as you request it.
| [Export](https://steampipe.io/docs/steampipe_export/overview) | Steampipe Plugin Exporters provide a flexible mechanism for exporting information from cloud services and APIs. Each exporter is a stand-alone binary that allows you to extract data using Steampipe plugins without a database.
| [Turbot Pipes](https://turbot.com/pipes/docs) | Turbot Pipes is the only intelligence, automation & security platform built specifically for DevOps. Pipes provide hosted Steampipe database instances, shared dashboards, snapshots, and more.

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

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Steampipe](https://steampipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #steampipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Azure DevOps Plugin](https://github.com/turbot/steampipe-plugin-azuredevops/labels/help%20wanted)
