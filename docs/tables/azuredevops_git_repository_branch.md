# Table: azuredevops_git_repository_branch

Represents Azure DevOps git repository branches.

## Examples

### Basic info

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
