# Table: azuredevops_group

All users added to Azure DevOps are added to one or more default security groups. Security groups are assigned permissions, which either allow or deny access to a feature or task. Members of a security group inherit the permissions assigned to the group.

## Examples

### Basic info

```sql
select
  principal_name,
  display_name,
  membership_state,
  domain,
  origin,
  description
from
  azuredevops_group;
```

### List inactive groups

```sql
select
  principal_name,
  display_name,
  membership_state,
  domain,
  origin,
  description
from
  azuredevops_group
where
  not membership_state;
```

### List empty groups

```sql
select
  principal_name,
  display_name,
  membership_state,
  domain,
  origin,
  description
from
  azuredevops_group
where
  memberships = '[]';
```

### Get parent group detail of a particular group

```sql
select
  principal_name,
  display_name,
  membership_state,
  domain,
  origin,
  description
from
  azuredevops_group
where
  descriptor in
  (
    select
      m ->> 'containerDescriptor'
    from
      azuredevops_group,
      jsonb_array_elements(memberships) as m
    where
      display_name = 'Build Administrators'
  );
```
