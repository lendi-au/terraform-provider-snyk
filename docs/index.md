---
page_title: "Snyk Provider"
subcategory: ""
description: |-
  Provider for managing Snyk Groups and Organizations.
---

# Snyk Provider



## Example Usage

```terraform
provider "snyk" {
  group_id = "GROUP_ID"
  api_key = "API_KEY"
}
```

## Schema

- **group_id** (String, Required) Group ID for the Snyk Group you want to manage. Can also be provided as an environment variable `SNYK_API_GROUP`.

- **api_key** (String, Required) API key with Group Admin scope for the Snyk Group you want to manage. Can also be provided as an environment variable `SNYK_API_KEY`.
