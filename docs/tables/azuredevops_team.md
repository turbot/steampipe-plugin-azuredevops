---
title: "Steampipe Table: azuredevops_team - Query Azure DevOps Teams using SQL"
description: "Allows users to query Azure DevOps Teams, providing valuable insights into team details and configurations."
---

# Table: azuredevops_team - Query Azure DevOps Teams using SQL

Azure DevOps Teams is a feature within Microsoft's Azure DevOps that enables collaboration and organization for software development projects. It provides a shared workspace for development teams to plan, track, and discuss work across the entire development process. Azure DevOps Teams supports efficient project management and team collaboration.

## Table Usage Guide

The `azuredevops_team` table provides insights into team configurations within Azure DevOps. As a project manager or team lead, explore team-specific details through this table, including team settings, members, and associated projects. Utilize it to manage team configurations, track team progress, and facilitate efficient team collaboration.

## Examples

### Basic info
Explore the teams within your Azure DevOps organization to understand their associated projects and access points. This can help establish an overview of your organization's structure and streamline project management processes.

```sql+postgres
select
  id,
  name,
  project_name,
  project_id,
  identity_url,
  url
from
  azuredevops_team;
```

```sql+sqlite
select
  id,
  name,
  project_name,
  project_id,
  identity_url,
  url
from
  azuredevops_team;
```

### List teams of a particular project
Discover the teams associated with a specific project in Azure DevOps. This can be particularly useful to understand the structure and distribution of teams for project management purposes.

```sql+postgres
select
  t.id as team_id,
  t.name as team_name,
  project_name,
  project_id,
  identity_url,
  t.url as url
from
  azuredevops_team as t,
  azuredevops_project as p
where
  t.project_id = p.id
  and p.name = 'private_project';
```

```sql+sqlite
select
  t.id as team_id,
  t.name as team_name,
  p.name as project_name,
  p.id as project_id,
  p.identity_url,
  t.url as url
from
  azuredevops_team as t,
  azuredevops_project as p
where
  t.project_id = p.id
  and p.name = 'private_project';
```

### List inactive teams
Uncover the details of teams that are currently inactive in Azure DevOps. This can be particularly useful in managing resources and ensuring efficient project allocation.

```sql+postgres
select
  id,
  name,
  project_name,
  project_id,
  identity_url,
  url
from
  azuredevops_team
where
  identity ->> 'isActive' = 'false';
```

```sql+sqlite
select
  id,
  name,
  project_name,
  project_id,
  identity_url,
  url
from
  azuredevops_team
where
  json_extract(identity, '$.isActive') = 'false';
```