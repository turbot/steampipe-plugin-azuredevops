# Table: azuredevops_pipeline

Azure Pipelines automatically builds and tests code projects. It supports all major languages and project types and combines continuous integration, continuous delivery, and continuous testing to build, test, and deliver your code to any destination.

## Examples

### Basic info

```sql
select
  id,
  name,
  configuration_type,
  project_id,
  folder,
  revision
from
  azuredevops_pipeline;
```

### List yaml based pipelines

```sql
select
  id,
  name,
  configuration_type,
  project_id,
  folder,
  revision
from
  azuredevops_pipeline
where
  configuration_type = 'yaml';
```

### List pipelines associated with a particular project

```sql
select
  l.id as pipeline_id,
  l.name as pipeline_name,
  configuration_type,
  project_id,
  folder,
  l.revision
from
  azuredevops_pipeline as l,
  azuredevops_project as p
where
  l.project_id = p.id
  and p.name = 'private_project';
```
