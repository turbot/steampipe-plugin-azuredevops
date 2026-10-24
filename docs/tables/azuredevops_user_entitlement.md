---
title: "Steampipe Table: azuredevops_user_entitlement - Query Azure DevOps User Entitlements using SQL"
description: "Allows users to query the users of an Azure DevOps organization with their access level, status and last access date."
---

# Table: azuredevops_user_entitlement - Query Azure DevOps User Entitlements using SQL

A user entitlement in Azure DevOps records what a user is allowed to use in an organization: the license (access level) they hold, where that license comes from, the status of their membership and when they last accessed the organization.

Access levels such as Basic, Stakeholder or a Visual Studio subscription decide which features a user can reach, and each paid access level consumes a license.

## Table Usage Guide

The `azuredevops_user_entitlement` table lists every user in the organization with their entitlement. As a security engineer, you can use this table to find inactive users who still hold access, review guest accounts, and see how licenses are assigned across the organization.

**Important Notes**
- Lists the whole organization; there are no required qualifiers.
- `last_accessed_date` is null for a user who has never accessed the organization.
- The personal access token needs the `Member Entitlement Management (Read)` scope (`vso.memberentitlementmanagement`). Without it the query fails with an error that starts `azuredevops: permission denied`.

## Examples

### Basic info
List the users of the organization with their access level, status and when they last accessed it.

```sql+postgres
select
  display_name,
  principal_name,
  user_type,
  access_level,
  status,
  last_accessed_date
from
  azuredevops_user_entitlement;
```

```sql+sqlite
select
  display_name,
  principal_name,
  user_type,
  access_level,
  status,
  last_accessed_date
from
  azuredevops_user_entitlement;
```

### Active users who have not accessed the organization for 90 days
Find users who still hold access but have not used it recently, or have never used it.

```sql+postgres
select
  principal_name,
  access_level,
  last_accessed_date
from
  azuredevops_user_entitlement
where
  status = 'active'
  and (
    last_accessed_date is null
    or last_accessed_date < now() - interval '90 days'
  );
```

```sql+sqlite
select
  principal_name,
  access_level,
  last_accessed_date
from
  azuredevops_user_entitlement
where
  status = 'active'
  and (
    last_accessed_date is null
    or last_accessed_date < datetime('now', '-90 days')
  );
```

### List guest users
Review the guest accounts that have been added to the organization.

```sql+postgres
select
  display_name,
  principal_name,
  origin,
  access_level,
  date_created
from
  azuredevops_user_entitlement
where
  user_type = 'guest';
```

```sql+sqlite
select
  display_name,
  principal_name,
  origin,
  access_level,
  date_created
from
  azuredevops_user_entitlement
where
  user_type = 'guest';
```

### Count users by access level
See how licenses are distributed across the organization.

```sql+postgres
select
  access_level,
  licensing_source,
  count(*) as user_count
from
  azuredevops_user_entitlement
group by
  access_level,
  licensing_source
order by
  user_count desc;
```

```sql+sqlite
select
  access_level,
  licensing_source,
  count(*) as user_count
from
  azuredevops_user_entitlement
group by
  access_level,
  licensing_source
order by
  user_count desc;
```

### Licensed users whose membership is no longer active
Join with `azuredevops_user` on the graph descriptor to find users who still hold an access level while their membership is inactive.

```sql+postgres
select
  e.principal_name,
  e.access_level,
  e.status,
  e.last_accessed_date
from
  azuredevops_user_entitlement as e
  join azuredevops_user as u on u.descriptor = e.descriptor
where
  not u.membership_state
  and e.account_license_type <> 'none';
```

```sql+sqlite
select
  e.principal_name,
  e.access_level,
  e.status,
  e.last_accessed_date
from
  azuredevops_user_entitlement as e
  join azuredevops_user as u on u.descriptor = e.descriptor
where
  u.membership_state = 0
  and e.account_license_type <> 'none';
```
