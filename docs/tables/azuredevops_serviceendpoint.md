# Table: azuredevops_serviceendpoint

Service endpoints are a way for Azure DevOps to connect to external systems or services. They're a bundle of properties securely stored by Azure DevOps, which includes but isn't limited to the following properties:

- Service name
- Description
- Server URL
- Certificates or tokens
- User names and passwords

## Examples

### Basic info

```sql
select
  id,
  name,
  is_ready,
  is_shared,
  owner,
  type
from
  azuredevops_serviceendpoint;
```

### List service endpoints which are not ready

```sql
select
  id,
  name,
  is_ready,
  is_shared,
  owner,
  type
from
  azuredevops_serviceendpoint
where
  not is_ready;
```

### List shared service endpoints

```sql
select
  id,
  name,
  is_ready,
  is_shared,
  owner,
  type
from
  azuredevops_serviceendpoint
where
  is_shared;
```

### List service endpoints owned by library

```sql
select
  id,
  name,
  is_ready,
  is_shared,
  owner,
  type
from
  azuredevops_serviceendpoint
where
  owner = 'Library';
```

### Get creator details of the service endpoints

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
  azuredevops_serviceendpoint;
```
