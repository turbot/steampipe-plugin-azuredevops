---
title: "Steampipe Table: azuredevops_user - Query Azure DevOps Users using SQL"
description: "Allows users to query Azure DevOps Users, specifically user details, providing insights into user roles, access, and permissions."
---

# Table: azuredevops_user - Query Azure DevOps Users using SQL

Azure DevOps is a Microsoft product that provides version control, reporting, requirements management, project management, automated builds, lab management, testing and release management capabilities. It covers the entire application lifecycle and enables DevOps capabilities. Azure DevOps can be used for any kind of application regardless of the framework, platform, or cloud.

The following types of users can join your Azure DevOps Services organization for free:

- Five users who get Basic features, such as version control, tools for Agile, Java, build, release, and more.
- Unlimited users who get Stakeholder features, such as working with your backlog, work items, and queries.
- Unlimited Visual Studio subscribers who also get Basic or Basic + Test Plan features, depending on their subscription level.

## Table Usage Guide

The `azuredevops_user` table provides insights into users within Azure DevOps. As a DevOps engineer or system administrator, explore user-specific details through this table, including roles, access levels, and associated metadata. Utilize it to uncover information about users, such as their access permissions, the projects they are associated with, and their activity patterns.

## Examples

### Basic info
Explore which users are active in your Azure DevOps environment. This can help in managing user access and understanding the distribution of users across different domains and origins.

```sql+postgres
select
  principal_name,
  display_name,
  membership_state,
  domain,
  origin
from
  azuredevops_user;
```

```sql+sqlite
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
Explore which users in your Azure DevOps environment are inactive. This can help maintain system security by identifying potential unused or unnecessary accounts.

```sql+postgres
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

```sql+sqlite
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
Explore which Azure DevOps users have their origin set to 'aad', enabling you to understand their source and manage user access effectively. This is particularly useful in large organizations where user management can be complex.

```sql+postgres
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

```sql+sqlite
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

### List the users that have been deleted in the identity provider
Explore the users who have been removed from the identity provider. This query is useful for auditing and maintaining security compliance by tracking changes in user access.

```sql+postgres
select
  principal_name,
  display_name,
  membership_state,
  domain,
  origin
from
  azuredevops_user
where
  is_deleted_in_origin;
```

```sql+sqlite
select
  principal_name,
  display_name,
  membership_state,
  domain,
  origin
from
  azuredevops_user
where
  is_deleted_in_origin = 1;
```

### List users who are part of `Project Collection Administrators` group
Explore which users are part of a specific administrative group within a project collection. This is useful for auditing purposes, allowing you to ensure only authorized individuals have administrative access.

```sql+postgres
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

```sql+sqlite
select
  principal_name,
  display_name,
  membership_state,
  domain,
  origin
from
  azuredevops_user,
  json_each(memberships) as m
where
  json_extract(m.value, '$.containerDescriptor') in
  (
    select
      descriptor
    from
      azuredevops_group
    where
      display_name = 'Project Collection Administrators'
  );
```