---
title: "Steampipe Table: azuredevops_git_branch_policy - Query Azure DevOps Git Branch Policies using SQL"
description: "Allows users to query the branch and repository policies of an Azure DevOps project, such as required reviewers, build validation, comment resolution and linked work items."
---

# Table: azuredevops_git_branch_policy - Query Azure DevOps Git Branch Policies using SQL

Azure DevOps branch policies protect branches by requiring pull requests to meet conditions before they merge: a minimum number of reviewers, a passing build, resolved comments or linked work items. Repository policies, such as a file size restriction, apply to every push to a repository.

## Table Usage Guide

The `azuredevops_git_branch_policy` table lists every policy configured in a project, one row per policy scope. As a DevOps or security engineer, you can use this table to check that default branches require reviewed, built pull requests, and to find policies that are disabled or only optional.

**Important Notes**
- You must specify `project` (name or ID) in the `where` clause. The `project` column echoes the value you give.
- A policy can cover several scopes, so there is one row per scope and `id` can repeat. `repository_id`, `ref_name` and `match_kind` come from that scope.
- A scope that covers every repository in the project has a null `repository_id`. A default-branch scope has `match_kind = 'DefaultBranch'` and a null `ref_name`.
- `repository_id` and `ref_name` are filtered by Postgres, not sent to Azure DevOps. Filtering on `repository_id = '...'` therefore drops project-wide policies that also apply to that repository, so use `repository_id is null or repository_id = '...'`.
- `settings` holds the policy settings exactly as Azure DevOps returns them, e.g. `minimumApproverCount`, `creatorVoteCounts`, `resetOnSourcePush`, `allowDownvotes` for reviewers, or `buildDefinitionId` and `validDuration` for build validation.
- The personal access token needs the `Code (Read)` scope (`vso.code`).

## Examples

### Basic info
List the policies of a project with the branch and repository each one applies to.

```sql+postgres
select
  id,
  type_display_name,
  is_enabled,
  is_blocking,
  repository_id,
  ref_name,
  match_kind
from
  azuredevops_git_branch_policy
where
  project = 'my-project';
```

```sql+sqlite
select
  id,
  type_display_name,
  is_enabled,
  is_blocking,
  repository_id,
  ref_name,
  match_kind
from
  azuredevops_git_branch_policy
where
  project = 'my-project';
```

### Enabled, blocking minimum-reviewer policies and their settings
Check how many approvals each required reviewer policy asks for, and whether the author's own vote counts.

```sql+postgres
select
  id,
  repository_id,
  ref_name,
  match_kind,
  settings ->> 'minimumApproverCount' as minimum_approvers,
  settings ->> 'creatorVoteCounts' as creator_vote_counts,
  settings ->> 'resetOnSourcePush' as reset_on_push
from
  azuredevops_git_branch_policy
where
  project = 'my-project'
  and type_display_name = 'Minimum number of reviewers'
  and is_enabled
  and is_blocking;
```

```sql+sqlite
select
  id,
  repository_id,
  ref_name,
  match_kind,
  json_extract(settings, '$.minimumApproverCount') as minimum_approvers,
  json_extract(settings, '$.creatorVoteCounts') as creator_vote_counts,
  json_extract(settings, '$.resetOnSourcePush') as reset_on_push
from
  azuredevops_git_branch_policy
where
  project = 'my-project'
  and type_display_name = 'Minimum number of reviewers'
  and is_enabled = 1
  and is_blocking = 1;
```

### List policies that are disabled or optional
Find policies that are configured but do not block a pull request from merging.

```sql+postgres
select
  id,
  type_display_name,
  repository_id,
  ref_name,
  is_enabled,
  is_blocking
from
  azuredevops_git_branch_policy
where
  project = 'my-project'
  and (not is_enabled or not is_blocking);
```

```sql+sqlite
select
  id,
  type_display_name,
  repository_id,
  ref_name,
  is_enabled,
  is_blocking
from
  azuredevops_git_branch_policy
where
  project = 'my-project'
  and (is_enabled = 0 or is_blocking = 0);
```

### Policies that apply to one repository, including project-wide ones
List every policy that covers a repository, counting the project-wide policies that have no `repository_id`.

```sql+postgres
select
  id,
  type_display_name,
  ref_name,
  match_kind
from
  azuredevops_git_branch_policy
where
  project = 'my-project'
  and (repository_id is null or repository_id = '49c1d4d0-be28-4f20-9e0a-4ecfafc39257');
```

```sql+sqlite
select
  id,
  type_display_name,
  ref_name,
  match_kind
from
  azuredevops_git_branch_policy
where
  project = 'my-project'
  and (repository_id is null or repository_id = '49c1d4d0-be28-4f20-9e0a-4ecfafc39257');
```

### List repositories whose default branch has no required reviewer policy
Join with `azuredevops_git_repository` to find repositories where pull requests into the default branch can merge without review.

```sql+postgres
select
  r.name as repository_name,
  r.default_branch
from
  azuredevops_project as p
  join azuredevops_git_repository as r on r.project_id = p.id
where
  p.name = 'my-project'
  and not exists (
    select
      1
    from
      azuredevops_git_branch_policy as bp
    where
      bp.project = p.name
      and bp.type_display_name = 'Minimum number of reviewers'
      and bp.is_enabled
      and bp.is_blocking
      and (bp.repository_id is null or bp.repository_id = r.id)
      and (bp.match_kind = 'DefaultBranch' or bp.ref_name = r.default_branch)
  );
```

```sql+sqlite
select
  r.name as repository_name,
  r.default_branch
from
  azuredevops_project as p
  join azuredevops_git_repository as r on r.project_id = p.id
where
  p.name = 'my-project'
  and not exists (
    select
      1
    from
      azuredevops_git_branch_policy as bp
    where
      bp.project = p.name
      and bp.type_display_name = 'Minimum number of reviewers'
      and bp.is_enabled = 1
      and bp.is_blocking = 1
      and (bp.repository_id is null or bp.repository_id = r.id)
      and (bp.match_kind = 'DefaultBranch' or bp.ref_name = r.default_branch)
  );
```
