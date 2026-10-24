---
title: "Steampipe Table: azuredevops_audit_stream - Query Azure DevOps Audit Streams using SQL"
description: "Allows users to query the audit streams that forward an Azure DevOps organization's audit log to a SIEM or log store."
---

# Table: azuredevops_audit_stream - Query Azure DevOps Audit Streams using SQL

Azure DevOps auditing records changes made across an organization, such as permission changes, deleted resources and policy updates. Audit streams send these audit events to an external destination, such as Splunk, Azure Monitor Logs or Azure Event Grid, so they can be retained beyond the built-in audit log and used for alerting.

A stream can only exist while auditing is on for the organization.

## Table Usage Guide

The `azuredevops_audit_stream` table lists the organization's audit streams. As a security engineer, you can use this table to check that audit events leave Azure DevOps for retention and alerting, and to find streams that have stopped sending events.

**Important Notes**
- Lists the whole organization; there are no required qualifiers.
- The stream's consumer inputs, which hold the destination's credentials, are never exposed.
- The personal access token needs the `Audit Log (Read)` scope (`vso.auditlog`), and the token's owner needs the View audit log permission. Without them the query fails with an error that starts `azuredevops: permission denied`.

## Examples

### Basic info
List the audit streams of the organization with their destination and status.

```sql+postgres
select
  id,
  display_name,
  consumer_type,
  status,
  created_time,
  updated_time
from
  azuredevops_audit_stream;
```

```sql+sqlite
select
  id,
  display_name,
  consumer_type,
  status,
  created_time,
  updated_time
from
  azuredevops_audit_stream;
```

### Enabled audit streams
List the streams that are currently sending audit events.

```sql+postgres
select
  display_name,
  consumer_type,
  updated_time
from
  azuredevops_audit_stream
where
  status = 'enabled';
```

```sql+sqlite
select
  display_name,
  consumer_type,
  updated_time
from
  azuredevops_audit_stream
where
  status = 'enabled';
```

### Audit streams that are not sending events
Find streams that a user or Azure DevOps has disabled, with the reason Azure DevOps gives.

```sql+postgres
select
  display_name,
  consumer_type,
  status,
  status_reason,
  updated_time
from
  azuredevops_audit_stream
where
  status <> 'enabled';
```

```sql+sqlite
select
  display_name,
  consumer_type,
  status,
  status_reason,
  updated_time
from
  azuredevops_audit_stream
where
  status <> 'enabled';
```

### Count audit streams by destination type
See which kinds of destination receive the organization's audit events.

```sql+postgres
select
  consumer_type,
  count(*) as stream_count
from
  azuredevops_audit_stream
group by
  consumer_type;
```

```sql+sqlite
select
  consumer_type,
  count(*) as stream_count
from
  azuredevops_audit_stream
group by
  consumer_type;
```
