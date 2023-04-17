# Table: azuredevops_project

Azure DevOps Project helps you launch an app on an Azure App Service of your choice in a few quick steps and set you up with everything you need for developing, deploying, and monitoring your app. Creating a DevOps Project provisions Azure resources and comes with a Git code repository, Application Insights integration and a continuous delivery pipeline setup to deploy to Azure. The DevOps Project dashboard lets you monitor code commits, builds and, deployments, from a single view in the Azure portal.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

```sql
select
  id,
  name,
  jsonb_pretty(capabilities -> 'versioncontrol') as version_control,
  jsonb_pretty(capabilities -> 'processTemplate') as process_template
from
  azuredevops_project;
```

### Get project default team details

```sql
select
  id,
  name,
  default_team ->> 'id' as default_team_id,
  default_team ->> 'name' as default_team_name,
  default_team ->> 'url' as default_team_url
from
  azuredevops_project;
```

### List project properties

```sql
select
  id,
  name,
  state,
  visibility,
  jsonb_pretty(properties) as properties
from
  azuredevops_project;
```
