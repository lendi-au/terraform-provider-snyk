# Snyk Terraform Provider

Provider for managing various aspects of Organizations within Snyk.

**Note: Requires a Business/Enterprise account, as that provides access to the API.**

*This currently has very limited functionality, limited to creating Organizations within a group.*

## Currently Planned Functionality

- [ ] Integrations
- [ ] Project imports from cloud integrations

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

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

To override any installed provider versions with the local development copy, add the following block to `~/.terraformrc`:

```tf
provider_installation {

  dev_overrides {
    "lendi-au/snyk" = "<FOLDER_CONTAINING_BINARY>"
  }

  direct {}

}
```


To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources within the configured Snyk group.

```sh
$ make testacc
```
