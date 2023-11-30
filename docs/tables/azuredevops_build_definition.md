---
title: "Steampipe Table: azuredevops_build_definition - Query Azure DevOps Build Definitions using SQL"
description: "Allows users to query Azure DevOps Build Definitions, enabling the retrieval of details on build pipelines, including their tasks, triggers, and repository information."
---

# Table: azuredevops_build_definition - Query Azure DevOps Build Definitions using SQL

Azure DevOps Build Definitions is a feature within Azure DevOps that allows you to define the steps required to build your code. These definitions include tasks, triggers, and repository information. It provides a structured way to automate the process of transforming your code into a finished product.

## Table Usage Guide

The `azuredevops_build_definition` table provides insights into build definitions within Azure DevOps. As a DevOps engineer, explore build definition details through this table, including tasks, triggers, and associated repository information. Utilize it to uncover information about build pipelines, such as their configuration, the tasks they perform, and the triggers that initiate them.

## Examples

### Basic info
Analyze the settings to understand the quality and creation date of different build definitions in your Azure DevOps project. This can help you assess the elements within your project and make necessary improvements.

```sql
select
  id,
  name,
  quality,
  project_id,
  repository_id,
  created_date
from
  azuredevops_build_definition;
```

### List build definitions where badge is enabled
Explore the build definitions in Azure DevOps that have the badge feature enabled. This can be useful for understanding which projects are actively promoting their build status.

```sql
select
  id,
  name,
  quality,
  project_id,
  repository_id,
  created_date
from
  azuredevops_build_definition
where
  badge_enabled;
```

### List build definitions where builds can be queued
Explore which build definitions in Azure DevOps are set up to allow queuing. This is useful for identifying areas where build processes can be scheduled and managed effectively.

```sql
select
  id,
  name,
  quality,
  project_id,
  repository_id,
  created_date
from
  azuredevops_build_definition
where
  queue_status = 'enabled';
```

### List build definitions of a particular project
Discover the build definitions associated with a specific project. This can be useful for auditing project configurations and understanding the quality of different builds within the project.

```sql
select
  id,
  name,
  quality,
  project_id,
  repository_id,
  created_date
from
  azuredevops_build_definition
where
  project ->> 'name' = 'private_project';
```

### List build definitions without repository
Discover the segments that have build definitions without an associated repository in Azure DevOps. This is useful to identify potential misconfigurations or orphaned build definitions that may need attention.

```sql
select
  id,
  name,
  quality,
  project_id,
  repository_id,
  created_date
from
  azuredevops_build_definition
where
  repository_id is null;
```

### Get the author details of a particular build definition
Explore the authorship information related to a specific build definition in Azure DevOps to understand who created or modified it, which can be crucial for tracking changes and maintaining accountability.

```sql
select
  authored_by ->> 'id' as author_id,
  authored_by ->> 'displayName' as display_name,
  authored_by ->> 'uniqueName' as unique_name,
  authored_by ->> 'descriptor' as descriptor,
  authored_by ->> 'url' as url,
  authored_by ->> 'imageUrl' as image_url
from
  azuredevops_build_definition
where
  name = 'private_project';
```