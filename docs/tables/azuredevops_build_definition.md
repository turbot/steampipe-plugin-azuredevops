# Table: azuredevops_build_definition

The build definitions are Azure DevOps tasks which build your project within your git repository.

## Examples

### Basic info

```sql
select
  id,
  name,
  quality,
  project_id,
  repository_id,
  created_date
from
  azuredevops_build_definition;
```

### List build definitions where badge is enabled

```sql
select
  id,
  name,
  quality,
  project_id,
  repository_id,
  created_date
from
  azuredevops_build_definition
where
  badge_enabled;
```

### List build definitions where builds can be queued

```sql
select
  id,
  name,
  quality,
  project_id,
  repository_id,
  created_date
from
  azuredevops_build_definition
where
  queue_status = 'enabled';
```

### List build definitions of a particular project

```sql
select
  id,
  name,
  quality,
  project_id,
  repository_id,
  created_date
from
  azuredevops_build_definition
where
  project ->> 'name' = 'private_project';
```

### List build definitions without repository

```sql
select
  id,
  name,
  quality,
  project_id,
  repository_id,
  created_date
from
  azuredevops_build_definition
where
  repository_id is null;
```

### Get the author details of a particular build definition

```sql
select
  authored_by ->> 'id' as author_id,
  authored_by ->> 'displayName' as display_name,
  authored_by ->> 'uniqueName' as unique_name,
  authored_by ->> 'descriptor' as descriptor,
  authored_by ->> 'url' as url,
  authored_by ->> 'imageUrl' as image_url
from
  azuredevops_build_definition
where
  name = 'private_project';
```
