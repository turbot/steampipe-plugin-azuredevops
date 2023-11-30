---
title: "Steampipe Table: azuredevops_build - Query Azure DevOps Builds using SQL"
description: "Allows users to query Azure DevOps Builds, providing insights into the build history, status, and details associated with each build."
---

# Table: azuredevops_build - Query Azure DevOps Builds using SQL

Azure DevOps is a Microsoft product that provides version control, reporting, requirements management, project management, automated builds, testing and release management capabilities. It covers the entire application lifecycle and enables DevOps capabilities. Azure DevOps can be used for a variety of application types.

## Table Usage Guide

The `azuredevops_build` table provides insights into the builds within Azure DevOps. As a DevOps engineer, explore build-specific details through this table, including build status, the build process, and associated metadata. Utilize it to uncover information about builds, such as those with failed status, the details of the build process, and the verification of build metadata.

## Examples

### Basic info
Explore the various builds in your Azure DevOps project, identifying their quality, status, and priority. This can help you understand which builds are kept indefinitely and potentially streamline your project management efforts.

```sql
select
  id,
  build_number,
  quality,
  project_id,
  status,
  keep_forever,
  priority
from
  azuredevops_build;
```

### List postponed builds
Explore which build projects in Azure DevOps have been postponed. This is useful for prioritizing and managing development workflow.

```sql
select
  id,
  build_number,
  quality,
  project_id,
  status,
  keep_forever,
  priority
from
  azuredevops_build
where
  status = 'postponed';
```

### List high priority builds
Determine the areas in which high priority builds are being used within your Azure DevOps projects. This can help prioritize resources and track project progress more effectively.

```sql
select
  id,
  build_number,
  quality,
  project_id,
  status,
  keep_forever,
  priority
from
  azuredevops_build
where
  priority = 'high';
```

### List builds which should be skipped by retention policies
Assess the elements within your Azure DevOps projects to pinpoint specific builds that have been marked to bypass retention policies. This can be beneficial for understanding which builds are being retained indefinitely, aiding in project management and resource allocation.

```sql
select
  id,
  build_number,
  quality,
  project_id,
  status,
  keep_forever,
  priority
from
  azuredevops_build
where
  keep_forever;
```

### List builds without repository
Analyze the settings to understand which Azure DevOps builds lack an associated repository. This can be useful for identifying potential configuration issues or orphaned builds.

```sql
select
  id,
  build_number,
  quality,
  project_id,
  status,
  keep_forever,
  priority
from
  azuredevops_build
where
  repository_id is null;
```

### List deleted builds
Explore which builds have been deleted in your Azure DevOps project. This can be useful for auditing purposes or for understanding the lifecycle of your builds.

```sql
select
  id,
  build_number,
  quality,
  project_id,
  status,
  keep_forever,
  priority
from
  azuredevops_build
where
  deleted;
```

### List builds associated with a particular project
Assess the elements within a specific project to identify the associated builds and their details, such as their quality, status, and priority. This can help in project management and in making decisions about resource allocation and task prioritization.

```sql
select
  id,
  build_number,
  quality,
  project_id,
  status,
  keep_forever,
  priority
from
  azuredevops_build
where
  project ->> 'name' = 'private_project';
```