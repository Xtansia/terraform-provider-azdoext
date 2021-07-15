# terraform-provider-azdoext

[![Tests](https://github.com/Xtansia/terraform-provider-azdoext/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/Xtansia/terraform-provider-azdoext/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/Xtansia/terraform-provider-azdoext/branch/main/graph/badge.svg?token=UL1SU6ES0L)](https://codecov.io/gh/Xtansia/terraform-provider-azdoext)

Resources currently unavailable
in [microsoft/terraform-provider-azuredevops](https://github.com/microsoft/terraform-provider-azuredevops).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules). Please see the Go documentation for the most
up to date information about using Go modules.

To add a new dependency `github.com/author/dependency`:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (
see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin`
directory.

To generate or update documentation, run `go generate`.
