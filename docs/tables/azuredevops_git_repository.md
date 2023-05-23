# Table: azuredevops_git_repository

A Git repository, or repo, is a folder that Git tracks changes in. There can be any number of repos on a computer, each stored in their own folder. Each Git repo on a system is independent, so changes saved in one Git repo don't affect the contents of another.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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
