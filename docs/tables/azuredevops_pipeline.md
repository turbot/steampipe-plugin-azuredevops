---
title: "Steampipe Table: azuredevops_pipeline - Query Azure DevOps Pipelines using SQL"
description: "Allows users to query Azure DevOps Pipelines, specifically the details of each pipeline, providing insights into pipeline configurations, status, and associated metadata."
---

# Table: azuredevops_pipeline - Query Azure DevOps Pipelines using SQL

Azure DevOps Pipelines is a cloud service that you can use to automatically build, test, and deploy your code project to any platform. It helps you catch bugs early in the development process and deploy updates more frequently and with higher quality. It supports continuous integration (CI) and continuous delivery (CD) to constantly and consistently test and build your code and ship it to any target.

## Table Usage Guide

The `azuredevops_pipeline` table provides insights into pipelines within Azure DevOps. As a DevOps engineer, explore pipeline-specific details through this table, including configurations, status, and associated metadata. Utilize it to uncover information about pipelines, such as those with specific configurations, the status of each pipeline, and the verification of pipeline settings.

## Examples

### Basic info
Explore which Azure DevOps pipelines are currently active, identifying their unique identifiers and names. This can be particularly useful for project managers and developers to gain insights into the status and revision history of their ongoing projects.

```sql
select
  id,
  name,
  configuration_type,
  project_id,
  folder,
  revision
from
  azuredevops_pipeline;
```

### List yaml based pipelines
Explore which Azure DevOps pipelines are based on YAML configurations. This can be useful in identifying pipelines that follow this specific format for potential updates or troubleshooting.

```sql
select
  id,
  name,
  configuration_type,
  project_id,
  folder,
  revision
from
  azuredevops_pipeline
where
  configuration_type = 'yaml';
```

### List pipelines associated with a particular project
Explore which pipelines are associated with a specific project in Azure DevOps. This is useful for managing and organizing project resources efficiently.

```sql
select
  l.id as pipeline_id,
  l.name as pipeline_name,
  configuration_type,
  project_id,
  folder,
  l.revision
from
  azuredevops_pipeline as l,
  azuredevops_project as p
where
  l.project_id = p.id
  and p.name = 'private_project';
```