---
title: "Steampipe Table: azuredevops_advanced_security_enablement - Query Azure DevOps Advanced Security Enablement using SQL"
description: "Allows users to query whether GitHub Advanced Security for Azure DevOps (Secret Protection and Code Security) is enabled on each repository of a project."
---

# Table: azuredevops_advanced_security_enablement - Query Azure DevOps Advanced Security Enablement using SQL

GitHub Advanced Security for Azure DevOps adds security scanning to Azure Repos. Secret Protection scans repositories for committed secrets and can block pushes that contain them (push protection). Code Security adds dependency scanning and code scanning with CodeQL.

Each feature is turned on per repository, and each one records when it was last turned on or off.

## Table Usage Guide

The `azuredevops_advanced_security_enablement` table returns one row per repository. As a security engineer, you can use this table to find repositories without secret scanning or push protection, and to check where dependency and code scanning are enabled.

**Important Notes**
- You must specify `project` (name or ID) in the `where` clause. `repository_id` is optional; leave it out to read every repository in the project.
- An organization without Advanced Security reports every flag false rather than an error.
- Disabled repositories are not read and report false, with `is_disabled` true.
- The personal access token needs the `Code (Read)` scope (`vso.code`) to list repositories and the `Advanced Security (Read)` scope (`vso.advsec`). Without `vso.advsec` the query fails with an error that starts `azuredevops: permission denied`.

## Examples

### Basic info
List each repository of a project with the Advanced Security features enabled on it.

```sql+postgres
select
  repository_name,
  enabled,
  secret_protection_enabled,
  block_pushes,
  code_security_enabled,
  codeql_enabled
from
  azuredevops_advanced_security_enablement
where
  project = 'my-project';
```

```sql+sqlite
select
  repository_name,
  enabled,
  secret_protection_enabled,
  block_pushes,
  code_security_enabled,
  codeql_enabled
from
  azuredevops_advanced_security_enablement
where
  project = 'my-project';
```

### Repositories without push protection
Find active repositories where pushes containing secrets are not blocked.

```sql+postgres
select
  repository_name,
  secret_protection_enabled,
  block_pushes
from
  azuredevops_advanced_security_enablement
where
  project = 'my-project'
  and not is_disabled
  and not (secret_protection_enabled and block_pushes);
```

```sql+sqlite
select
  repository_name,
  secret_protection_enabled,
  block_pushes
from
  azuredevops_advanced_security_enablement
where
  project = 'my-project'
  and is_disabled = 0
  and not (secret_protection_enabled = 1 and block_pushes = 1);
```

### Repositories without dependency or code scanning
Find active repositories where Code Security is off or CodeQL default setup is not enabled.

```sql+postgres
select
  repository_name,
  code_security_enabled,
  dependency_scanning_injection_enabled,
  codeql_enabled
from
  azuredevops_advanced_security_enablement
where
  project = 'my-project'
  and not is_disabled
  and (
    not code_security_enabled
    or not codeql_enabled
  );
```

```sql+sqlite
select
  repository_name,
  code_security_enabled,
  dependency_scanning_injection_enabled,
  codeql_enabled
from
  azuredevops_advanced_security_enablement
where
  project = 'my-project'
  and is_disabled = 0
  and (
    code_security_enabled = 0
    or codeql_enabled = 0
  );
```

### Repositories where Secret Protection was changed in the last 30 days
See which repositories had secret scanning turned on or off recently.

```sql+postgres
select
  repository_name,
  secret_protection_enabled,
  secret_protection_last_changed_date
from
  azuredevops_advanced_security_enablement
where
  project = 'my-project'
  and secret_protection_last_changed_date > now() - interval '30 days';
```

```sql+sqlite
select
  repository_name,
  secret_protection_enabled,
  secret_protection_last_changed_date
from
  azuredevops_advanced_security_enablement
where
  project = 'my-project'
  and secret_protection_last_changed_date > datetime('now', '-30 days');
```

### Count open alerts on each scanned repository
Join with `azuredevops_advanced_security_alert` to count the active alerts on repositories that have Advanced Security enabled; a count of zero means the repository is scanned and has no open alerts.

```sql+postgres
select
  e.repository_name,
  count(a.alert_id) as active_alert_count
from
  azuredevops_advanced_security_enablement as e
  left join azuredevops_advanced_security_alert as a
    on a.project = e.project
    and a.repository_id = e.repository_id
    and a.state = 'active'
where
  e.project = 'my-project'
  and e.enabled
group by
  e.repository_name
order by
  active_alert_count desc;
```

```sql+sqlite
select
  e.repository_name,
  count(a.alert_id) as active_alert_count
from
  azuredevops_advanced_security_enablement as e
  left join azuredevops_advanced_security_alert as a
    on a.project = e.project
    and a.repository_id = e.repository_id
    and a.state = 'active'
where
  e.project = 'my-project'
  and e.enabled = 1
group by
  e.repository_name
order by
  active_alert_count desc;
```
