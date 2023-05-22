# Table: azuredevops_team_member

Represents Azure DevOps team members. By default, team members inherit the permissions afforded to members of the project Contributors group. Members of this group can add and modify source code, create and delete test runs, and create and modify work items. Team members can collaborate on a Git project or check in work to the team's code base.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  is_team_admin,
  project_id,
  team_id,
  url
from
  azuredevops_team_member;
```

### List team members who are admins

```sql
select
  id,
  display_name,
  is_team_admin,
  project_id,
  team_id,
  url
from
  azuredevops_team_member
where
  is_team_admin;
```

### List deleted team members

```sql
select
  id,
  display_name,
  is_team_admin,
  project_id,
  team_id,
  url
from
  azuredevops_team_member
where
  is_deleted_in_origin;
```

### List team members of a particular project

```sql
select
  m.id as member_id,
  display_name,
  is_team_admin,
  project_id,
  team_id
from
  azuredevops_team_member as m,
  azuredevops_project as p
where
  m.project_id = p.id
  and p.name = 'private_project';
```

### List team members of a particular team

```sql
select
  m.id as member_id,
  display_name,
  is_team_admin,
  team_id
from
  azuredevops_team_member as m,
  azuredevops_team as t
where
  m.team_id = t.id
  and t.name = 'private_project Team';
```
