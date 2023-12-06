---
title: "Steampipe Table: azuredevops_git_repository - Query Azure DevOps Git Repositories using SQL"
description: "Allows users to query Git Repositories in Azure DevOps, providing details on each repository's project, name, size, and other associated metadata."
---

# Table: azuredevops_git_repository - Query Azure DevOps Git Repositories using SQL

Azure DevOps is a service from Microsoft that provides developer services for support teams to plan work, collaborate on code development, and build and deploy applications. Git Repositories in Azure DevOps are part of the version control system that allows developers to collaborate on code development. It provides a place for teams to store, version, and share source code.

## Table Usage Guide

The `azuredevops_git_repository` table provides insights into Git Repositories within Azure DevOps. As a DevOps engineer, explore repository-specific details through this table, including project association, repository name, size, and other associated metadata. Utilize it to uncover information about repositories, such as their default branch, web URL, and the verification of their visibility and fork status.

## Examples

### Basic info
Gain insights into the general information of your AzureDevOps Git repositories, such as whether they are forks or not and their size. This can help with project management and resource allocation.

```sql+postgres
select
  id,
  name,
  default_branch,
  is_fork,
  project_id,
  size
from
  azuredevops_git_repository;
```

```sql+sqlite
select
  id,
  name,
  default_branch,
  is_fork,
  project_id,
  size
from
  azuredevops_git_repository;
```

### List forked repositories
Explore which repositories in your Azure DevOps have been forked. This can help identify duplicated or shared project resources, aiding in project management and resource allocation.

```sql+postgres
select
  id,
  name,
  default_branch,
  is_fork,
  project_id,
  size
from
  azuredevops_git_repository
where
  is_fork;
```

```sql+sqlite
select
  id,
  name,
  default_branch,
  is_fork,
  project_id,
  size
from
  azuredevops_git_repository
where
  is_fork;
```

### List repositories of a particular project
Explore the different repositories associated with a specific project to gain insights into their default branches, size, and whether they are a fork of another repository. This is useful for managing and understanding the structure of a particular project in Azure DevOps.

```sql+postgres
select
  id,
  name,
  default_branch,
  is_fork,
  project_id,
  size
from
  azuredevops_git_repository
where
  project ->> 'name' = 'private_project';
```

```sql+sqlite
select
  id,
  name,
  default_branch,
  is_fork,
  project_id,
  size
from
  azuredevops_git_repository
where
  json_extract(project, '$.name') = 'private_project';
```