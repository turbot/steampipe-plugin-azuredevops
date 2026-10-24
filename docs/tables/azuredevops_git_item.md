---
title: "Steampipe Table: azuredevops_git_item - Query Azure DevOps Git repository files using SQL"
description: "Allows users to list the files and folders of an Azure DevOps Git repository at a branch, tag or commit, and read the content of a single file."
---

# Table: azuredevops_git_item - Query Azure DevOps Git repository files using SQL

Azure DevOps Repos hosts Git repositories. The Git Items API returns the tree of a repository at any branch, tag or commit in a single request, and the content of individual files.

## Table Usage Guide

The `azuredevops_git_item` table lists the files in a repository and reads individual files without cloning it. As a DevOps or security engineer, you can use it to find dependency manifests and lockfiles across repositories, or to check that a required file (such as `CODEOWNERS` or a pipeline definition) exists and contains what it should.

**Important Notes**
- You must specify `project` and `repository_id` in the `where` or `join` clause to query this table. Both accept a name or an ID.
- Give `path` to read one file and its `content`, or leave it out to list the whole subtree under `scope_path` (the repository root by default) in one request. `path` must name a file, and `path` and `scope_path` cannot be combined.
- `version` takes a branch, tag or commit and defaults to the repository's default branch. `version_type` (`branch`, `tag` or `commit`) defaults to `commit` for a 40-character hash and `branch` otherwise.
- Paths start with `/`. Listings do not report file sizes; `size` is set when a text file is read with `path`.
- `content` and `size` are null for folders and binary files.
- Requests honor Azure DevOps rate limiting: a throttled request is retried after the `Retry-After` it names, and a response that reports the limit is near delays the next request to the organization.
- The personal access token needs the `Code (Read)` scope.

## Examples

### Basic info
List every file and folder in the default branch of a repository, with the Git object each one points to.

```sql+postgres
select
  path,
  git_object_type,
  is_folder,
  object_id,
  commit_id
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo';
```

```sql+sqlite
select
  path,
  git_object_type,
  is_folder,
  object_id,
  commit_id
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo';
```

### List the files under a folder on a branch
List only the files, not the folders, below one folder of a branch.

```sql+postgres
select
  path,
  object_id
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo'
  and version = 'develop'
  and scope_path = '/src'
  and not is_folder;
```

```sql+sqlite
select
  path,
  object_id
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo'
  and version = 'develop'
  and scope_path = '/src'
  and is_folder = 0;
```

### Find dependency manifests and lockfiles
Locate the package manifests and lockfiles in a repository, for example to build a dependency inventory.

```sql+postgres
select
  path
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo'
  and git_object_type = 'blob'
  and (
    path like '%/package-lock.json'
    or path like '%/packages.lock.json'
    or path like '%/go.mod'
    or path like '%/requirements.txt'
  );
```

```sql+sqlite
select
  path
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo'
  and git_object_type = 'blob'
  and (
    path like '%/package-lock.json'
    or path like '%/packages.lock.json'
    or path like '%/go.mod'
    or path like '%/requirements.txt'
  );
```

### Read one file at a commit
Read the content of a single file as it was at a given commit.

```sql+postgres
select
  path,
  size,
  content
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo'
  and version = '23d0bc5b128a10056dc68afece360d8a0fabb014'
  and path = '/package.json';
```

```sql+sqlite
select
  path,
  size,
  content
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo'
  and version = '23d0bc5b128a10056dc68afece360d8a0fabb014'
  and path = '/package.json';
```

### List repositories in a project without a CODEOWNERS file
Check every repository of a project for a `CODEOWNERS` file at the root of its default branch.

```sql+postgres
select
  r.name as repository_name
from
  azuredevops_project as p
  join azuredevops_git_repository as r on r.project_id = p.id
where
  p.name = 'my-project'
  and not exists (
    select
      1
    from
      azuredevops_git_item as i
    where
      i.project = p.name
      and i.repository_id = r.id
      and i.path = '/CODEOWNERS'
  );
```

```sql+sqlite
select
  r.name as repository_name
from
  azuredevops_project as p
  join azuredevops_git_repository as r on r.project_id = p.id
where
  p.name = 'my-project'
  and not exists (
    select
      1
    from
      azuredevops_git_item as i
    where
      i.project = p.name
      and i.repository_id = r.id
      and i.path = '/CODEOWNERS'
  );
```
