# Table: azuredevops_team

You create a team in Azure DevOps that corresponds to a group of project members focused on specific products, services, or feature areas. You add teams to provide them the tools they need to manage their backlog, plan sprints, configure dashboards, define alerts, and set team favorites.

## Examples

### Basic info

```sql
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

```sql
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

### List inactive teams

```sql
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
