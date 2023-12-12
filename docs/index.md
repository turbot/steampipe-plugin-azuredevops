---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/azuredevops.svg"
brand_color: "#0090F1"
display_name: "Azure DevOps"
short_name: "azuredevops"
description: "Steampipe plugin to query projects, groups, builds and more from Azure DevOps."
og_description: "Query Azure DevOps with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/azuredevops-social-graphic.png"
engines: ["steampipe", "sqlite", "postgres", "export"]
---

# Azure DevOps + Steampipe

[Azure DevOps](https://dev.azure.com) is a software as a service (SaaS) platform that provides DevOps practices and tools for the end-to-end software life cycle.

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

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

## Documentation

- **[Table definitions & examples →](/plugins/turbot/azuredevops/tables)**

## Quick start

### Install

Download and install the latest Azure DevOps plugin:

```sh
steampipe plugin install azuredevops
```

### Credentials

| Item        | Description                                                                                                                                                                                                                                                                                                                                            |
| ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Credentials | Azure DevOps requires an [Organization URL](https://learn.microsoft.com/en-us/azure/devops/extend/develop/work-with-urls?view=azure-devops&tabs=http) and a [Personal Access Token](https://learn.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate?view=azure-devops&tabs=Windows) for all requests. |
| Permissions | Personal Access Tokens have the same permissions as the user who creates them, and if the user permissions change, the Personal Access Token permissions also change.                                                                                                                                                                                  |
| Radius      | Each connection represents a single Azure DevOps Installation.                                                                                                                                                                                                                                                                                         |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/azuredevops.spc`)<br />2. Credentials specified in environment variables, e.g., `AZDO_ORG_SERVICE_URL` and `AZDO_PERSONAL_ACCESS_TOKEN`.                                                                                                                                |

### Configuration

Installing the latest azuredevops plugin will create a config file (`~/.steampipe/config/azuredevops.spc`) with a single connection named `azuredevops`:

Configure your account details in `~/.steampipe/config/azuredevops.spc`:

```hcl
connection "azuredevops" {
  plugin = "azuredevops"

  # `organization_url`: Azure DevOps Organization URL. (Required)
  # For more information on the Organization URL, please see https://learn.microsoft.com/en-us/azure/devops/extend/develop/work-with-urls?view=azure-devops&tabs=http.
  # Can also be set with the AZDO_ORG_SERVICE_URL environment variable.
  # organization_url = "https://dev.azure.com/test"

  # `personal_access_token`: Azure DevOps Personal Access Token. (Required)
  # For more information on the Personal Access Token, please see https://learn.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate?view=azure-devops&tabs=Windows.
  # Can also be set with the AZDO_PERSONAL_ACCESS_TOKEN environment variable.
  # personal_access_token = "wf3hahidy7i7fkzmeqr3e6fbjwuspabpo766grp7hl4o65v2"
}
```

Alternatively, you can also use the standard Azure DevOps environment variables to obtain credentials **only if other arguments (`organization_url` and `personal_access_token`) are not specified** in the connection:

```sh
export AZDO_ORG_SERVICE_URL=https://dev.azure.com/test
export AZDO_PERSONAL_ACCESS_TOKEN=wf3hahidy7i7fkzmeqr3e6fbjwuspabpo766grp7hl4o65v2
```


