---
title: "Steampipe Table: azuredevops_group - Query Azure DevOps Groups using SQL"
description: "Allows users to query Azure DevOps Groups, providing detailed insights into group membership, permissions, and associated metadata."
---

# Table: azuredevops_group - Query Azure DevOps Groups using SQL

Azure DevOps Groups is a feature in Azure DevOps that allows for the management of user and group permissions across the platform. It enables administrators to control access to resources such as repositories, pipelines, and boards. Azure DevOps Groups simplifies the process of managing permissions by allowing users to be added to or removed from multiple resources at once.

## Table Usage Guide

The `azuredevops_group` table provides insights into groups within Azure DevOps. As a DevOps engineer, explore group-specific details through this table, including group membership, permissions, and associated metadata. Utilize it to manage user permissions across multiple resources effectively and ensure the right level of access for each user.

## Examples

### Basic info
Explore which Azure DevOps groups are active within your domain. This is useful for understanding user access and permissions within your organization.

```sql+postgres
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

```sql+sqlite
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
Identify groups within your Azure DevOps that are currently inactive. This could be useful for optimizing resource use, ensuring security, or maintaining an organized workspace.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that are classified as empty groups in the Azure DevOps platform. This can be useful to identify and clean up unused or redundant groups, enhancing the efficiency of your resource management.

```sql+postgres
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

```sql+sqlite
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
Determine the parent group details associated with a specific group in Azure DevOps. This can be useful for understanding group hierarchies and membership states, especially when managing user access and permissions.

```sql+postgres
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

```sql+sqlite
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
      json_extract(m.value, '$.containerDescriptor')
    from
      azuredevops_group,
      json_each(memberships) as m
    where
      display_name = 'Build Administrators'
  );
```