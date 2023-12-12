---
title: "Steampipe Table: azuredevops_project - Query Azure DevOps Projects using SQL"
description: "Allows users to query Azure DevOps Projects. It provides information about the project's name, description, visibility, version control, and process template."
---

# Table: azuredevops_project - Query Azure DevOps Projects using SQL

Azure DevOps is a Microsoft product that provides version control, reporting, requirements management, project management, automated builds, lab management, testing and release management capabilities. It covers the entire application lifecycle and enables DevOps capabilities. Azure DevOps can be used for a variety of application types.

## Table Usage Guide

The `azuredevops_project` table provides insights into projects within Azure DevOps. As a project manager or developer, explore project-specific details through this table, including project name, description, visibility, version control, and process template. Utilize it to uncover information about projects, such as those with public visibility, the version control used, and the process template associated with each project.

## Examples

### Basic info
Explore which Azure DevOps projects are currently active, their visibility status, and when they were last updated. This information can help assess the current state of your projects and identify any that may require attention.

```sql+postgres
select
  id,
  name,
  state,
  visibility,
  abbreviation,
  last_update_time
from
  azuredevops_project;
```

```sql+sqlite
select
  id,
  name,
  state,
  visibility,
  abbreviation,
  last_update_time
from
  azuredevops_project;
```

### List public projects
Discover the segments that are public in your Azure DevOps projects, allowing you to assess the elements within your setup that are visible to all users. This can help you maintain appropriate access controls and security measures.

```sql+postgres
select
  id,
  name,
  state,
  visibility,
  abbreviation,
  last_update_time
from
  azuredevops_project
where
  visibility = 'public';
```

```sql+sqlite
select
  id,
  name,
  state,
  visibility,
  abbreviation,
  last_update_time
from
  azuredevops_project
where
  visibility = 'public';
```

### List projects which are in the `createPending` state
Discover the segments that are pending creation within your projects. This can aid in understanding project progress and managing resources effectively.

```sql+postgres
select
  id,
  name,
  state,
  visibility,
  abbreviation,
  last_update_time
from
  azuredevops_project
where
  state = 'createPending';
```

```sql+sqlite
select
  id,
  name,
  state,
  visibility,
  abbreviation,
  last_update_time
from
  azuredevops_project
where
  state = 'createPending';
```

### Show project capabilities
Explore the capabilities of your projects to understand the version control and process template settings. This can help you manage and optimize your project settings in Azure DevOps.

```sql+postgres
select
  id,
  name,
  jsonb_pretty(capabilities -> 'versioncontrol') as version_control,
  jsonb_pretty(capabilities -> 'processTemplate') as process_template
from
  azuredevops_project;
```

```sql+sqlite
select
  id,
  name,
  capabilities as version_control,
  capabilities as process_template
from
  azuredevops_project;
```

### Get project default team details
Gain insights into the default team details associated with various projects to better understand team structure and project management within Azure DevOps.

```sql+postgres
select
  id,
  name,
  default_team ->> 'id' as default_team_id,
  default_team ->> 'name' as default_team_name,
  default_team ->> 'url' as default_team_url
from
  azuredevops_project;
```

```sql+sqlite
select
  id,
  name,
  json_extract(default_team, '$.id') as default_team_id,
  json_extract(default_team, '$.name') as default_team_name,
  json_extract(default_team, '$.url') as default_team_url
from
  azuredevops_project;
```

### List project properties
Explore the various properties of your projects in Azure DevOps. This is useful for gaining insights into project details like their state and visibility settings.

```sql+postgres
select
  id,
  name,
  state,
  visibility,
  jsonb_pretty(properties) as properties
from
  azuredevops_project;
```

```sql+sqlite
select
  id,
  name,
  state,
  visibility,
  properties
from
  azuredevops_project;
```