# Table: azuredevops_dashboard

Dashboards are customizable interactive signboards that provide real-time information. Dashboards are associated with a team or a project and display configurable charts and widgets.

## Examples

### Basic info

```sql
select
  id,
  name,
  dashboard_scope,
  project_id,
  owner_id,
  group_id
from
  azuredevops_dashboard;
```

### List dashboards with project scope

```sql
select
  id,
  name,
  dashboard_scope,
  project_id,
  owner_id,
  group_id
from
  azuredevops_dashboard
where
  dashboard_scope = 'project';
```

### List dashboards of a particular project

```sql
select
  d.id as dashboard_id,
  d.name,
  d.dashboard_scope,
  d.project_id,
  d.owner_id,
  d.group_id
from
  azuredevops_dashboard as d,
  azuredevops_project as p
where
  d.project_id = p.id
  and p.name = 'private_project';
```

### Get owner details of a particular dashboard

```sql
select
  t.id as team_id,
  t.name,
  t.project_id,
  t.project_name
from
  azuredevops_dashboard as d,
  azuredevops_team as t
where
  d.owner_id = t.id
  and d.name = 'test';
```
