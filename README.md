# Snyk Terraform Provider

Provider for managing various aspects of Organizations within Snyk.

**Note: Requires a Business/Enterprise account, as that provides access to the API.**

Currently provides Terraform resources for:

- Organizations
- Integrations (yet to be released, but code on `main` - still need to write up docs)

## Currently Planned Functionality

- [ ] Projects (importing from integrations)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.16

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

```tf
terraform {
  required_providers {
    snyk = {
      source = "lendi-au/snyk"
      version = "<version>"
    }
  }
}

provider "snyk" {
    group_id = "GROUP_ID" # can also be provided in env as SNYK_API_GROUP
    api_key = "API_KEY" # can also be provided in env as SNYK_API_KEY, requires Group admin scope
}

resource "snyk_organization" "test" {
    name = "Test Organization"
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources within the configured Snyk group - requires an API key and Group ID to be set as environment variables.

```sh
$ make testacc
```
