---
title: "Steampipe Table: azuredevops_git_repository_branch - Query Azure DevOps Git Repositories using SQL"
description: "Allows users to query Git Repositories in Azure DevOps, specifically the branches within each repository, providing insights into version control and development workflows."
---

# Table: azuredevops_git_repository_branch - Query Azure DevOps Git Repositories using SQL

Azure DevOps is a service within Microsoft Azure that supports development teams in planning work, collaborating on code development, and deploying applications. Git Repositories in Azure DevOps provide version control and allow for collaborative code development. The branches within these repositories represent independent lines of development that can be created, merged, or deleted.

## Table Usage Guide

The `azuredevops_git_repository_branch` table provides insights into branches within Git Repositories in Azure DevOps. As a DevOps engineer, explore branch-specific details through this table, including the branch name, repository it belongs to, and its commit history. Utilize it to manage and track the development process across different branches, ensuring code integrity and efficient workflows.

## Examples

### Basic info
Analyze the settings to understand the status of different branches in your Azure DevOps Git repository. This would be useful to assess the progress of different projects or features by comparing the number of commits ahead or behind the base version.

```sql
select
  name,
  repository_id,
  ahead_count,
  behind_count,
  is_base_version
from
  azuredevops_git_repository_branch;
```

### List base version branches
Explore which branches serve as the base versions in your Azure DevOps Git repositories. This can help you understand the structure of your repositories and monitor the development progress in terms of commits ahead or behind the base version.

```sql
select
  name,
  repository_id,
  ahead_count,
  behind_count,
  is_base_version
from
  azuredevops_git_repository_branch
where
  is_base_version;
```

### Get current commit details of main branch
Explore the most recent changes made to the main branch of your project. This can help you understand the nature of the changes, who made them, and their potential impact on the project.

```sql
select
  commit ->> 'commitId' as commit_id,
  commit ->> 'comment' as comment,
  commit -> 'committer' as committer,
  commit -> 'changeCounts' as change_counts,
  commit -> 'parents' as parents
from
  azuredevops_git_repository_branch
where
  name = 'main';
```

### List branches of a particular repository
Explore the different branches within a specific repository to gain insights into their respective ahead and behind counts, as well as their base version status. This can be particularly useful for managing and tracking changes in a complex development environment.

```sql
select
  b.name as branch_name,
  repository_id,
  ahead_count,
  behind_count,
  is_base_version
from
  azuredevops_git_repository_branch as b,
  azuredevops_git_repository as r
where
  b.repository_id = r.id
  and r.name = 'test_repo';
```