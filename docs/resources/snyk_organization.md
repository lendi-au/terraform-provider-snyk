---
page_title: "snyk_organization Resource - terraform-provider-snyk"
subcategory: ""
description: |-
  Organization resource
---

# Resource `snyk_organization`

An individual organization within your Snyk Group.

## Example Usage

```terraform
resource "snyk_organization" "example" {
  name = "Example"
}
```

## Schema

### Required

- **name** (String, Required) The display name of the Organization. **Note: Currently, changing this will recreate the Org, destroying the existing one. This is to due with limitations within the Snyk API for changing display names. If you wish to update the name, do so in the UI first.**


