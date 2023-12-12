## v0.2.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/install/steampipe.sh), as a [Postgres FDW](https://steampipe.io/install/postgres.sh), as a [SQLite extension](https://steampipe.io/install/sqlite.sh) and as a standalone [exporter](https://steampipe.io/install/export.sh).
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension.
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-azuredevops/blob/main/docs/LICENSE).

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server enacapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#19](https://github.com/turbot/steampipe-plugin-azuredevops/pull/19))

## v0.1.1 [2023-10-05]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#12](https://github.com/turbot/steampipe-plugin-azuredevops/pull/12))

## v0.1.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#10](https://github.com/turbot/steampipe-plugin-azuredevops/pull/10))
- Recompiled plugin with Go version `1.21`. ([#10](https://github.com/turbot/steampipe-plugin-azuredevops/pull/10))

## v0.0.1 [2023-05-23]

_What's new?_

- New tables added
  - [azuredevops_build](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_build)
  - [azuredevops_build_definition](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_build_definition)
  - [azuredevops_dashboard](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_dashboard)
  - [azuredevops_git_repository](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_git_repository)
  - [azuredevops_git_repository_branch](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_git_repository_branch)
  - [azuredevops_group](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_group)
  - [azuredevops_pipeline](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_pipeline)
  - [azuredevops_project](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_project)
  - [azuredevops_release](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_release)
  - [azuredevops_serviceendpoint](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_serviceendpoint)
  - [azuredevops_team](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_team)
  - [azuredevops_team_member](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_team_member)
  - [azuredevops_user](https://hub.steampipe.io/plugins/turbot/azuredevops/tables/azuredevops_user)
