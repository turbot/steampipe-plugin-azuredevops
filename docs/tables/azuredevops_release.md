# Table: azuredevops_release

A release represents continuous delivery in Azure DevOps. A pipeline usually takes code, builds it, tests, and creates an artifact. Release pipelines takes the artifact and deploys it. So release pipelines are primarily split up into two components, the artifact and stages.

## Examples

### Basic info

```sql
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

```sql
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

### List abandoned releases

```sql
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

```sql
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

```sql
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
