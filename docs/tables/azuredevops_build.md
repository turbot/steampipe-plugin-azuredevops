# Table: azuredevops_build

Represents Azure DevOps builds.

## Examples

### Basic info

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
