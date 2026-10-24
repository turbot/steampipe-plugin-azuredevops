---
title: "Steampipe Table: azuredevops_variable_group - Query Azure DevOps Variable Groups using SQL"
description: "Allows users to query the variable groups of an Azure DevOps project, counting secret variables and flagging plain variables whose names look like credentials."
---

# Table: azuredevops_variable_group - Query Azure DevOps Variable Groups using SQL

Variable groups in Azure DevOps Pipelines store values that are shared across pipelines. Variables can be marked secret, or a group can be linked to an Azure key vault so its variables are read from the vault. Values not marked secret are stored and shown in plain text.

## Table Usage Guide

The `azuredevops_variable_group` table lists the variable groups in a project. As a DevOps or security engineer, you can use this table to find credentials kept as plain variables instead of secrets or a key vault, and to see which groups are shared with other projects.

**Important Notes**
- You must specify `project` (name or ID) in the `where` clause. The `project` column echoes the value you give.
- Variable values are never read into the table: only names and the secret flag are decoded.
- `plaintext_variable_names_suspected` lists variables not marked secret whose names look like credentials (password, secret, token, API key, connection string, PAT, names ending in key). It is a name heuristic and can include false positives. It is always empty for key vault groups, whose variables come from the vault.
- The personal access token needs the `Variable Groups (Read)` scope (`vso.variablegroups_read`).

## Examples

### Basic info
List the variable groups of a project with how many of their variables are secret.

```sql+postgres
select
  id,
  name,
  type,
  is_shared,
  variable_count,
  secret_variable_count
from
  azuredevops_variable_group
where
  project = 'my-project';
```

```sql+sqlite
select
  id,
  name,
  type,
  is_shared,
  variable_count,
  secret_variable_count
from
  azuredevops_variable_group
where
  project = 'my-project';
```

### Variable groups with credential-like plain variables
Find groups holding variables whose names suggest a credential but that are not marked secret.

```sql+postgres
select
  name,
  secret_variable_count,
  plaintext_variable_names_suspected
from
  azuredevops_variable_group
where
  project = 'my-project'
  and jsonb_array_length(plaintext_variable_names_suspected) > 0;
```

```sql+sqlite
select
  name,
  secret_variable_count,
  plaintext_variable_names_suspected
from
  azuredevops_variable_group
where
  project = 'my-project'
  and json_array_length(plaintext_variable_names_suspected) > 0;
```

### List variable groups linked to a key vault
Identify the groups that read their variables from an Azure key vault, and which vault.

```sql+postgres
select
  name,
  key_vault_name,
  variable_count
from
  azuredevops_variable_group
where
  project = 'my-project'
  and key_vault_linked;
```

```sql+sqlite
select
  name,
  key_vault_name,
  variable_count
from
  azuredevops_variable_group
where
  project = 'my-project'
  and key_vault_linked = 1;
```

### List variable groups shared with other projects
Find groups that pipelines in other projects can also use.

```sql+postgres
select
  name,
  description,
  type,
  variable_count
from
  azuredevops_variable_group
where
  project = 'my-project'
  and is_shared;
```

```sql+sqlite
select
  name,
  description,
  type,
  variable_count
from
  azuredevops_variable_group
where
  project = 'my-project'
  and is_shared = 1;
```

### Find the variable groups that define a variable
Search the variable names of every group for one variable.

```sql+postgres
select
  g.name as variable_group_name,
  v as variable_name
from
  azuredevops_variable_group as g,
  jsonb_array_elements_text(g.variable_names) as v
where
  g.project = 'my-project'
  and v = 'DATABASE_HOST';
```

```sql+sqlite
select
  g.name as variable_group_name,
  v.value as variable_name
from
  azuredevops_variable_group as g,
  json_each(g.variable_names) as v
where
  g.project = 'my-project'
  and v.value = 'DATABASE_HOST';
```
