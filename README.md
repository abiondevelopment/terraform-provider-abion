# Terraform Provider: Abion

The Abion provider supports resources that performs DNS updates and data sources via the Abion API.  

## Requirements

* [Terraform](https://www.terraform.io/downloads)
* [Go](https://go.dev/doc/install) (1.22)
* [GNU Make](https://www.gnu.org/software/make/)
* [golangci-lint](https://golangci-lint.run/usage/install/#local-installation) (optional)

## Using the provider

Official documentation on how to use this provider can be found on the
[Terraform Registry](https://registry.terraform.io/providers/abion/abion/latest/docs).

## Developing the Abion Provider
The provided `GNUmakefile` defines commands generally useful during development,
like for trigger a Golang build and install the Abion Provider, running acceptance tests, generating documentation, code formatting and linting.

`git clone` this repository and `cd` into its directory
```shell
make 
```
will run linting, formatting, generate docs build and install the provider

### Building

```shell
make install 
```
will trigger the Golang build and install the provider in your `$GOBIN` folder

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

### Testing

In order to test the Abion provider, run

```shell
ABION_API_KEY=<api key> ABION_API_HOST=<host> make testacc
```
to run the full suite of acceptance tests

It's important to note that acceptance tests (`make testacc`) will actually call the Abion API and update zone records. You need a valid API KEY, access to the existing zone and also, your IP address must be whitelisted by Abion to be able to access the Abion API

### Generating documentation

The Abion provider uses [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs/)
to generate documentation and store it in the `docs/` directory.
Once a release is cut, the Terraform Registry will download the documentation from `docs/`
and associate it with the release version. Read more about how this works on the
[official page](https://www.terraform.io/registry/providers/docs).

Use 
```shell
make generate
``` 
to ensure the documentation is re-generated with latest changes.

### Code formatting and linting
```shell
make fmt
make lint
````

### Using a development build

If [running acceptance tests](#Testing) isn't enough, it's possible to set up a local terraform configuration
to use a development builds of the provider. This can be achieved by leveraging the Terraform CLI
[configuration file development overrides](https://www.terraform.io/cli/config/config-file#development-overrides-for-provider-developers).

First, use `make install` to place a fresh development build of the provider in your
[`${GOBIN}`](https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies)
(defaults to `${GOPATH}/bin` or `${HOME}/go/bin` if `${GOPATH}` is not set). Repeat
this every time you make changes to the provider locally.

Then, setup your environment following [these instructions](https://www.terraform.io/plugin/debugging#terraform-cli-development-overrides)
to make your local terraform use your local build.

Example of a `.terraformrc` file:
```
provider_installation {

  dev_overrides {
      "abion/abion" = "<path_to_/go/bin>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```