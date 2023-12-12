---
title: "Steampipe Table: azuredevops_team_member - Query Azure DevOps Team Members using SQL"
description: "Allows users to query Team Members in Azure DevOps, specifically the members' details, providing insights into team structure and roles."
---

# Table: azuredevops_team_member - Query Azure DevOps Team Members using SQL

Azure DevOps is a service within Microsoft Azure that supports development teams with version control, reporting, requirements management, project management, automated builds, lab management, testing, and release management capabilities. It provides a rich ecosystem for managing multi-stage, multi-environment, and multi-provider pipelines. Team Members in Azure DevOps are individuals who are part of a particular team, and their details, roles, and permissions can be managed and queried.

## Table Usage Guide

The `azuredevops_team_member` table provides insights into team members within Azure DevOps. As a project manager or team lead, explore member-specific details through this table, including roles, permissions, and associated metadata. Utilize it to manage and understand team structure, roles, and permissions, and to ensure the right individuals have access to the right resources.

## Examples

### Basic info
Explore which team members hold administrative roles within your Azure DevOps project. This query could be beneficial for management or auditing purposes, providing a quick overview of team structures and roles in the project.

```sql+postgres
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

```sql+sqlite
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
Explore which team members hold admin status in your Azure DevOps setup. This can help manage permissions and roles within your projects.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that include team members who have been removed from your Azure DevOps team. This can be useful for auditing changes to team composition or identifying potential security issues.

```sql+postgres
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

```sql+sqlite
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
Explore which team members are part of a specific project in Azure DevOps. This is useful for project managers who need to quickly understand the composition of their project team, including who the team admin is.

```sql+postgres
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

```sql+sqlite
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
Explore which team members belong to a specific team in Azure DevOps, and determine their roles within the team. This can be particularly useful in understanding team composition and identifying team administrators.

```sql+postgres
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

```sql+sqlite
select
  m.id as member_id,
  display_name,
  is_team_admin,
  team_id
from
  azuredevops_team_member as m
join
  azuredevops_team as t on m.team_id = t.id
where
  t.name = 'private_project Team';
```