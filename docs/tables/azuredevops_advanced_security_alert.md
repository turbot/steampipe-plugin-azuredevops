---
title: "Steampipe Table: azuredevops_advanced_security_alert - Query Azure DevOps Advanced Security Alerts using SQL"
description: "Allows users to query the Advanced Security alerts (dependency, secret and code scanning) on the repositories of an Azure DevOps project."
---

# Table: azuredevops_advanced_security_alert - Query Azure DevOps Advanced Security Alerts using SQL

GitHub Advanced Security for Azure DevOps raises alerts for vulnerable dependencies, committed secrets and code weaknesses found in Azure Repos. Each alert has a type, a severity and a state, and records when it was first and last seen and when it was fixed.

## Table Usage Guide

The `azuredevops_advanced_security_alert` table lists alerts per repository. As a security engineer, you can use this table to count open critical alerts, track how long alerts stay open, and find committed secrets that still need to be rotated.

**Important Notes**
- You must specify `project` (name or ID) in the `where` clause. `repository_id`, `state` and `alert_type` are optional and are sent to Azure DevOps as filters.
- Azure DevOps returns alerts on the default branch only, and high-confidence secret alerts only.
- A repository without Advanced Security contributes no rows. Join `azuredevops_advanced_security_enablement` to tell "no alerts" from "not scanned".
- Secret material (`truncatedSecret`, validation fingerprints) is never exposed.
- The `title` column holds the title of the alert.
- The personal access token needs the `Code (Read)` scope (`vso.code`) and the `Advanced Security (Read)` scope (`vso.advsec`).

## Examples

### Basic info
List the alerts on the repositories of a project with their type, severity and state.

```sql+postgres
select
  repository_id,
  alert_id,
  alert_type,
  severity,
  state,
  title,
  first_seen_date
from
  azuredevops_advanced_security_alert
where
  project = 'my-project';
```

```sql+sqlite
select
  repository_id,
  alert_id,
  alert_type,
  severity,
  state,
  title,
  first_seen_date
from
  azuredevops_advanced_security_alert
where
  project = 'my-project';
```

### Open alerts by repository and severity
Count the active alerts of each type and severity on every repository of a project.

```sql+postgres
select
  repository_id,
  alert_type,
  severity,
  count(*) as alert_count
from
  azuredevops_advanced_security_alert
where
  project = 'my-project'
  and state = 'active'
group by
  repository_id,
  alert_type,
  severity;
```

```sql+sqlite
select
  repository_id,
  alert_type,
  severity,
  count(*) as alert_count
from
  azuredevops_advanced_security_alert
where
  project = 'my-project'
  and state = 'active'
group by
  repository_id,
  alert_type,
  severity;
```

### Open secret alerts
List committed secrets that are still active, so they can be rotated and removed.

```sql+postgres
select
  repository_id,
  alert_id,
  title,
  git_ref,
  first_seen_date
from
  azuredevops_advanced_security_alert
where
  project = 'my-project'
  and alert_type = 'secret'
  and state = 'active';
```

```sql+sqlite
select
  repository_id,
  alert_id,
  title,
  git_ref,
  first_seen_date
from
  azuredevops_advanced_security_alert
where
  project = 'my-project'
  and alert_type = 'secret'
  and state = 'active';
```

### Critical and high alerts open for more than 30 days
Find serious alerts that have stayed open past a remediation deadline.

```sql+postgres
select
  repository_id,
  alert_id,
  alert_type,
  severity,
  title,
  first_seen_date
from
  azuredevops_advanced_security_alert
where
  project = 'my-project'
  and state = 'active'
  and severity in ('critical', 'high')
  and first_seen_date < now() - interval '30 days'
order by
  first_seen_date;
```

```sql+sqlite
select
  repository_id,
  alert_id,
  alert_type,
  severity,
  title,
  first_seen_date
from
  azuredevops_advanced_security_alert
where
  project = 'my-project'
  and state = 'active'
  and severity in ('critical', 'high')
  and first_seen_date < datetime('now', '-30 days')
order by
  first_seen_date;
```

### Repositories with open critical or high alerts
Join with `azuredevops_advanced_security_enablement` to show each repository's name and whether push protection is on, next to its count of serious open alerts.

```sql+postgres
select
  e.repository_name,
  e.block_pushes,
  count(*) as critical_or_high_alert_count
from
  azuredevops_advanced_security_alert as a
  join azuredevops_advanced_security_enablement as e
    on e.project = a.project
    and e.repository_id = a.repository_id
where
  a.project = 'my-project'
  and a.state = 'active'
  and a.severity in ('critical', 'high')
group by
  e.repository_name,
  e.block_pushes
order by
  critical_or_high_alert_count desc;
```

```sql+sqlite
select
  e.repository_name,
  e.block_pushes,
  count(*) as critical_or_high_alert_count
from
  azuredevops_advanced_security_alert as a
  join azuredevops_advanced_security_enablement as e
    on e.project = a.project
    and e.repository_id = a.repository_id
where
  a.project = 'my-project'
  and a.state = 'active'
  and a.severity in ('critical', 'high')
group by
  e.repository_name,
  e.block_pushes
order by
  critical_or_high_alert_count desc;
```
