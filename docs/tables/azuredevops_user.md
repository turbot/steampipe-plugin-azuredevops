# Table: azuredevops_user

Represents Azure DevOps users. Users can include human users, service accounts, and service principals. The following types of users can join your Azure DevOps Services organization for free:

- Five users who get Basic features, such as version control, tools for Agile, Java, build, release, and more.
- Unlimited users who get Stakeholder features, such as working with your backlog, work items, and queries.
- Unlimited Visual Studio subscribers who also get Basic or Basic + Test Plan features, depending on their subscription level.

## Examples

### Basic info

```sql
select
  principal_name,
  display_name,
  membership_state,
  domain,
  origin
from
  azuredevops_user;
```

### List inactive users

```sql
select
  principal_name,
  display_name,
  membership_state,
  domain,
  origin
from
  azuredevops_user
where
  not membership_state;
```

### List users from aad

```sql
select
  principal_name,
  display_name,
  membership_state,
  domain,
  origin
from
  azuredevops_user
where
  origin = 'aad';
```

### List users who are part of `Project Collection Administrators` group

```sql
select
  principal_name,
  display_name,
  membership_state,
  domain,
  origin
from
  azuredevops_user,
  jsonb_array_elements(memberships) as m
where
  m ->> 'containerDescriptor' in
  (
    select
      descriptor
    from
      azuredevops_group
    where
      display_name = 'Project Collection Administrators'
  );
```
