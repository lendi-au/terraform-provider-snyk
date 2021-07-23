---
page_title: "snyk_organization Data Source - terraform-provider-snyk"
subcategory: ""
description: |-
  Data source for individual Snyk organizations
---

# Data Source `snyk_organization`

Data source for fetching information about Snyk organizations.

(Not too much real-world usage, was just a testing ground for schema work. Will probably be removed in later releases.)

## Example Usage

```terraform
data "snyk_organization" "example" {
  id = "ORG_ID"
}
```

## Schema

### Required

- **id** (String, Required) Sample attribute.

## Outputs

- **name** (String) Display name of Organization.


