---
title: "Steampipe Table: azuredevops_release - Query Azure DevOps Releases using SQL"
description: "Allows users to query Azure DevOps Releases, providing insights into the release pipelines, stages, and deployment status."
---

# Table: azuredevops_release - Query Azure DevOps Releases using SQL

Azure DevOps is a suite of development tools, services, and features that enables teams to plan, develop, test, and deliver software more efficiently. The Releases in Azure DevOps provide a consistent and reliable delivery pipeline, managing the stages of deployment and tracking the status of each. It offers the ability to automate deployments, monitor the health of the pipeline, and roll back if necessary.

## Table Usage Guide

The `azuredevops_release` table provides insights into the release pipelines within Azure DevOps. As a DevOps engineer, you can explore details about each release, including the stages, deployment status, and associated metadata. Use this table to manage and monitor your software delivery pipeline, ensuring a reliable and efficient deployment process.

## Examples

### Basic info
Explore which projects have been created in Azure DevOps and their current status, allowing you to understand the lifecycle and longevity of each project for better management.

```sql+postgres
select
  id,
  name,
  status,
  project_id,
  created_on,
  keep_forever
from
  azuredevops_release;
```

```sql+sqlite
select
  id,
  name,
  status,
  project_id,
  created_on,
  keep_forever
from
  azuredevops_release;
```

### List releases which should be skipped by retention policies
Assess the elements within your Azure DevOps releases that are earmarked to be kept indefinitely, allowing you to identify potential areas for data management and storage optimization. This is particularly useful in managing retention policies and ensuring efficient use of resources.

```sql+postgres
select
  id,
  name,
  status,
  project_id,
  created_on,
  keep_forever
from
  azuredevops_release
where
  keep_forever;
```

```sql+sqlite
select
  id,
  name,
  status,
  project_id,
  created_on,
  keep_forever
from
  azuredevops_release
where
  keep_forever = 1;
```

### List abandoned releases
Identify instances where releases have been abandoned in Azure DevOps. This aids in understanding project progress and identifying potential bottlenecks or areas for improvement.

```sql+postgres
select
  id,
  name,
  status,
  project_id,
  created_on,
  keep_forever
from
  azuredevops_release
where
  status = 'abandoned';
```

```sql+sqlite
select
  id,
  name,
  status,
  project_id,
  created_on,
  keep_forever
from
  azuredevops_release
where
  status = 'abandoned';
```

### List manual releases
Explore which project releases in Azure DevOps have been manually initiated. This can help in assessing the frequency of manual interventions and their impact on the overall project timeline.

```sql+postgres
select
  id,
  name,
  status,
  project_id,
  created_on,
  keep_forever
from
  azuredevops_release
where
  reason = 'manual';
```

```sql+sqlite
select
  id,
  name,
  status,
  project_id,
  created_on,
  keep_forever
from
  azuredevops_release
where
  reason = 'manual';
```

### Get creator details of a particular release
Explore which individual created a specific release in the Azure DevOps platform. This is useful for accountability and tracking changes within the project management lifecycle.

```sql+postgres
select
  name,
  created_by ->> 'id' as id,
  created_by ->> 'displayName' as display_name,
  created_by ->> 'uniqueName' as unique_name,
  created_by ->> 'descriptor' as descriptor,
  created_by ->> 'url' as url,
  created_by ->> 'imageUrl' as image_url
from
  azuredevops_release
where
  name = 'Release-1';
```

```sql+sqlite
select
  name,
  json_extract(created_by, '$.id') as id,
  json_extract(created_by, '$.displayName') as display_name,
  json_extract(created_by, '$.uniqueName') as unique_name,
  json_extract(created_by, '$.descriptor') as descriptor,
  json_extract(created_by, '$.url') as url,
  json_extract(created_by, '$.imageUrl') as image_url
from
  azuredevops_release
where
  name = 'Release-1';
```