---
title: "Steampipe Table: azuredevops_serviceendpoint - Query Azure DevOps Service Endpoints using SQL"
description: "Allows users to query Service Endpoints in Azure DevOps, gaining visibility on the connections between Azure DevOps and external services."
---

# Table: azuredevops_serviceendpoint - Query Azure DevOps Service Endpoints using SQL

A Service Endpoint in Azure DevOps is a connection point to an external service. It is used to abstract the underlying connectivity to an external service, providing a stable interface for Azure DevOps while allowing the underlying connection details to change. Service Endpoints support services like Jenkins, Azure Resource Manager, and others.

## Table Usage Guide

The `azuredevops_serviceendpoint` table provides insights into Service Endpoints within Azure DevOps. As a DevOps engineer, you can use this table to explore details about Service Endpoints, including their types, authorization details, and associated project details. This table is beneficial for understanding the connections between your Azure DevOps environment and external services.

## Examples

### Basic info
Explore the readiness and sharing status of Azure DevOps service endpoints to understand their accessibility and ownership. This can help in identifying any potential issues or risks related to endpoint accessibility and management.

```sql
select
  id,
  name,
  is_ready,
  is_shared,
  owner,
  type
from
  azuredevops_serviceendpoint;
```

### List service endpoints which are not ready
Explore which Azure DevOps service endpoints are not ready for use. This can be beneficial in identifying potential issues in your deployment pipeline.

```sql
select
  id,
  name,
  is_ready,
  is_shared,
  owner,
  type
from
  azuredevops_serviceendpoint
where
  not is_ready;
```

### List shared service endpoints
Analyze the settings to understand which service endpoints in Azure DevOps are shared, providing insights into how resources are being utilized across different teams or projects. This can aid in resource optimization and improve collaboration efficiency.

```sql
select
  id,
  name,
  is_ready,
  is_shared,
  owner,
  type
from
  azuredevops_serviceendpoint
where
  is_shared;
```

### List service endpoints owned by library
Discover the segments that are owned by a specific entity in Azure DevOps. This is useful for understanding who has control over certain resources and can aid in resource management.

```sql
select
  id,
  name,
  is_ready,
  is_shared,
  owner,
  type
from
  azuredevops_serviceendpoint
where
  owner = 'Library';
```

### Get creator details of the service endpoints
Explore the creators of service endpoints, gaining insights into their unique identifiers and display names to understand who has set up each endpoint. This can be useful for auditing or tracking changes in your Azure DevOps environment.

```sql
select
  name,
  created_by ->> 'id' as id,
  created_by ->> 'displayName' as display_name,
  created_by ->> 'uniqueName' as unique_name,
  created_by ->> 'descriptor' as descriptor,
  created_by ->> 'url' as url,
  created_by ->> 'imageUrl' as image_url
from
  azuredevops_serviceendpoint;
```