---
title: "Steampipe Table: azuredevops_project_pipeline_settings - Query Azure DevOps Project Pipeline Settings using SQL"
description: "Allows users to query the pipeline general settings of an Azure DevOps project, such as job authorization scope, settable variables and fork protection."
---

# Table: azuredevops_project_pipeline_settings - Query Azure DevOps Project Pipeline Settings using SQL

Each Azure DevOps project has pipeline general settings that limit what a pipeline's job token can reach, which variables can be set when a run is queued, and how pull requests from forks are built. They also control whether classic pipelines can be created and whether shell task arguments are sanitized.

## Table Usage Guide

The `azuredevops_project_pipeline_settings` table returns one row per project. As a DevOps or security engineer, you can use this table to check that job tokens are scoped to the project, that builds of forks cannot read secrets, and that classic pipelines are turned off.

**Important Notes**
- You must specify `project` (name or ID) in the `where` clause. The `project` column echoes the value you give.
- A setting the API does not report is null rather than false. `settings` holds every setting exactly as Azure DevOps returns it.
- Organization-level pipeline settings have no public REST API. The values here are the project's, as the project API reports them.
- To check fork settings per pipeline, read `triggers` on `azuredevops_build_definition`: a `pullRequest` trigger carries `forks.enabled` and `forks.allowSecrets`.
- The personal access token needs the `Project and Team (Read)` scope (`vso.project`).

## Examples

### Basic info
Get the main job authorization and variable settings of a project.

```sql+postgres
select
  project,
  enforce_job_auth_scope,
  enforce_job_auth_scope_for_releases,
  enforce_referenced_repo_scoped_token,
  enforce_settable_var,
  status_badges_are_private
from
  azuredevops_project_pipeline_settings
where
  project = 'my-project';
```

```sql+sqlite
select
  project,
  enforce_job_auth_scope,
  enforce_job_auth_scope_for_releases,
  enforce_referenced_repo_scoped_token,
  enforce_settable_var,
  status_badges_are_private
from
  azuredevops_project_pipeline_settings
where
  project = 'my-project';
```

### Projects whose job tokens are not limited to the project
Check whether build or release pipelines get a job token that can reach other projects in the organization.

```sql+postgres
select
  project,
  enforce_job_auth_scope,
  enforce_job_auth_scope_for_releases,
  enforce_referenced_repo_scoped_token
from
  azuredevops_project_pipeline_settings
where
  project = 'my-project'
  and not (enforce_job_auth_scope and enforce_job_auth_scope_for_releases);
```

```sql+sqlite
select
  project,
  enforce_job_auth_scope,
  enforce_job_auth_scope_for_releases,
  enforce_referenced_repo_scoped_token
from
  azuredevops_project_pipeline_settings
where
  project = 'my-project'
  and not (enforce_job_auth_scope and enforce_job_auth_scope_for_releases);
```

### Fork protection
Review how pull requests from forked repositories are built and whether those builds can read secrets.

```sql+postgres
select
  project,
  builds_enabled_for_forks,
  fork_protection_enabled,
  enforce_no_access_to_secrets_from_forks,
  enforce_job_auth_scope_for_forks,
  is_comment_required_for_pull_request
from
  azuredevops_project_pipeline_settings
where
  project = 'my-project';
```

```sql+sqlite
select
  project,
  builds_enabled_for_forks,
  fork_protection_enabled,
  enforce_no_access_to_secrets_from_forks,
  enforce_job_auth_scope_for_forks,
  is_comment_required_for_pull_request
from
  azuredevops_project_pipeline_settings
where
  project = 'my-project';
```

### Check whether classic pipelines can still be created
Confirm that new classic build and release pipelines are blocked in a project that has moved to YAML.

```sql+postgres
select
  project,
  disable_classic_pipeline_creation,
  disable_classic_build_pipeline_creation,
  disable_classic_release_pipeline_creation
from
  azuredevops_project_pipeline_settings
where
  project = 'my-project';
```

```sql+sqlite
select
  project,
  disable_classic_pipeline_creation,
  disable_classic_build_pipeline_creation,
  disable_classic_release_pipeline_creation
from
  azuredevops_project_pipeline_settings
where
  project = 'my-project';
```

### Settings of every project in the organization
Join with `azuredevops_project` to read the settings of all projects at once.

```sql+postgres
select
  p.name as project_name,
  s.enforce_job_auth_scope,
  s.enforce_settable_var,
  s.enforce_no_access_to_secrets_from_forks,
  s.enable_shell_tasks_args_sanitizing
from
  azuredevops_project as p
  join azuredevops_project_pipeline_settings as s on s.project = p.name;
```

```sql+sqlite
select
  p.name as project_name,
  s.enforce_job_auth_scope,
  s.enforce_settable_var,
  s.enforce_no_access_to_secrets_from_forks,
  s.enable_shell_tasks_args_sanitizing
from
  azuredevops_project as p
  join azuredevops_project_pipeline_settings as s on s.project = p.name;
```
