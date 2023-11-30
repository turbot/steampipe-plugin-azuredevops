---
title: "Steampipe Table: azuredevops_dashboard - Query Azure DevOps Dashboards using SQL"
description: "Allows users to query Azure DevOps Dashboards, providing insights into the layout and widgets of each dashboard."
---

# Table: azuredevops_dashboard - Query Azure DevOps Dashboards using SQL

Azure DevOps Dashboards are a customizable interactive signboard that provides real-time information, analytics, and insights. They are a place to find quick answers and information about the work items and builds. Dashboards are customizable and can contain charts, graphs, and lists, among other things.

## Table Usage Guide

The `azuredevops_dashboard` table provides insights into the dashboards within Azure DevOps. As a project manager or team lead, explore dashboard-specific details through this table, including layout, widgets, and associated metadata. Utilize it to uncover information about the dashboards, such as their configuration, the widgets they contain, and their layout.

## Examples

### Basic info
Discover the segments that are linked to your Azure DevOps dashboard, including the project and group it belongs to, as well as its owner. This can help in managing and organizing your projects more efficiently.

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
Discover the dashboards that are specifically scoped to projects in Azure DevOps. This is useful for assessing the distribution of resources and managing project-specific data.

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
Discover the segments that belong to a specific project by identifying all associated dashboards. This is beneficial in understanding the structure and distribution of resources within a particular project.

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
Explore the teams associated with a specific Azure DevOps dashboard to understand its ownership and associated project details. This can be useful for auditing purposes or to manage permissions and access controls.

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