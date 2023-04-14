---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/azuredevops.svg"
brand_color: "#0090f1"
display_name: "Azure DevOps"
short_name: "azuredevops"
description: "Steampipe plugin to query projects, groups, builds and more from Azure DevOps."
og_description: "Query Azure DevOps with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/azuredevops-social-graphic.png"
---

# Azure DevOps + Steampipe

[Azure DevOps](https://dev.azure.com) is a software as a service (SaaS) platform that provides DevOps practices and tools for the end-to-end software life cycle.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

Get Azure DevOps project details:

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

## Get started

### Install

Download and install the latest Azure DevOps plugin:

```bash
steampipe plugin install azuredevops
```

### Configuration

Installing the latest azuredevops plugin will create a config file (`~/.steampipe/config/azuredevops.spc`) with a single connection named `azuredevops`:

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

- `organization_url` - Azure DevOps Organization URL. Can also be set with the `AZDO_ORG_SERVICE_URL` environment variable.
- `personal_access_token` - Azure DevOps Personal Access Token. Can also be set with the `AZDO_PERSONAL_ACCESS_TOKEN` environment variable.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-azuredevops
- Community: [Slack Channel](https://steampipe.io/community/join)
