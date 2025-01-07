# Terraform Provider: Abion

The Terraform Abion Provider supports resources that performs DNS record updates and data sources for DNS records via the Abion API.
You **must** have an Abion account to retrieve an API key. Contact [Abion](https://abion.com/) for help how to create an account and API key. To access 
the Abion API you also need to get your IP addresses whitelisted by Abion. The Terraform Abion Provider operates on record level so the actual zone
must exist and the API Key/account need to have access to the zone.  

## Requirements

* [Terraform](https://www.terraform.io/downloads)
* [Go](https://go.dev/doc/install) (1.22)
* [GNU Make](https://www.gnu.org/software/make/)
* [golangci-lint](https://golangci-lint.run/welcome/install#local-installation) (optional)

## Limitations

* You are only allowed to make 100 updates of a zone per day. 
* The zone must exist (and managed by Abion) and the API Key/account needs access to the zone. 
* The supported record types: 
  * `A`
  * `AAAA`
  * `CNAME`
  * `MX`
  * `TXT`
  * `NS`
  * `SRV` 
  * `PTR`


## Using the provider

Official documentation on how to use this provider can be found on the
[Terraform Registry](https://registry.terraform.io/providers/abiondevelopment/abion/latest/docs).

## Developing the Abion Provider

The provided `GNUmakefile` defines commands generally useful during development, like for trigger a Golang build and install the Abion Provider, 
running acceptance tests, generating documentation, code formatting and linting.

`git clone` this repository and `cd` into its directory

The default `make` command will execute linting, formatting, generate docs, build and install the Abion provider
```shell
make 
```

### Building

The `make install` target will trigger the Golang build and install the provider in your `$GOBIN` folder

```shell
make install 
```

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to the Abion provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

### Testing

In order to test the Terraform Abion Provider, run to run the full suite of acceptance tests by executing:

```shell
ABION_API_KEY=<api key> ABION_API_HOST=<host> make testacc
```
Configure the ABION_API_KEY and (optional) ABION_API_HOST environment variables. 

**NOTE!**  
It's important to understand that acceptance tests (`make testacc`) will actually call the Abion API and update zone records. You need a valid API KEY, access to the existing zone and also, your IP address must be whitelisted by Abion to be able to access the Abion API

### Generating documentation

The Terraform Abion Provider uses [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs/)
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

### Using a development build in Terraform

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
      "abiondevelopment/abion" = "<path_to_/go/bin>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

## Releasing

The release process is automated via GitHub Actions, and it's defined in the Workflow
[release.yml](./.github/workflows/release.yml).

Each release is cut by pushing a [semantically versioned](https://semver.org/) tag (e.g. `v.0.1.0`) to the main branch.
