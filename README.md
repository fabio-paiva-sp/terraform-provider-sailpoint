# SailPoint Terraform Provider

## Roadmap

### High priority resources 

- [ ] Managed Clusters - https://developer.sailpoint.com/docs/api/v2025/managed-clusters
- [ ] Sources - https://developer.sailpoint.com/docs/api/v2025/sources
- [ ] Managed Clients - https://developer.sailpoint.com/docs/api/v2025/managed-clients
- [ ] Identity Profiles - https://developer.sailpoint.com/docs/api/v2025/identity-profiles
- [ ] Access Profiles - https://developer.sailpoint.com/docs/api/v2025/access-profiles
- [ ] Roles - https://developer.sailpoint.com/docs/api/v2025/roles
- [ ] Accounts (Delimited file) - https://developer.sailpoint.com/docs/api/v2025/accounts
- [ ] Transforms - https://developer.sailpoint.com/docs/api/v2025/transforms
- [ ] Workflows - https://developer.sailpoint.com/docs/api/v2025/workflows
- [ ] Custom Forms - https://developer.sailpoint.com/docs/api/v2025/custom-forms
- [ ] Launchers - https://developer.sailpoint.com/docs/api/v2025/get-launchers

#### Data lookup

- [ ] Identities https://developer.sailpoint.com/docs/api/v2025/identities
- [ ] Accounts (Delimited file) - https://developer.sailpoint.com/docs/api/v2025/accounts
- [ ] Entitlements - https://developer.sailpoint.com/docs/api/v2025/entitlements
- [ ] Access Profiles - https://developer.sailpoint.com/docs/api/v2025/access-profiles
- [ ] Roles - https://developer.sailpoint.com/docs/api/v2025/roles
- [ ] Connectors - https://developer.sailpoint.com/docs/api/v2025/connectors

### Beta resources

- [ ] Applications - https://developer.sailpoint.com/docs/api/beta/apps

### Experimental resources

- [ ] Governance Groups - https://developer.sailpoint.com/docs/api/v2025/governance-groups
- [ ] Triggers - https://developer.sailpoint.com/docs/api/v2025/triggers
- [ ] Custom User Levels - https://developer.sailpoint.com/docs/api/v2025/custom-user-levels

### Medium priority resources

- [ ] Connector Rule Management - https://developer.sailpoint.com/docs/api/v2025/connector-rule-management
- [ ] Certification Campaign Filters - https://developer.sailpoint.com/docs/api/v2025/certification-campaign-filters
- [ ] Access Request segment - https://developer.sailpoint.com/docs/api/v2025/segments
- [ ] Service Desk Integrations - https://developer.sailpoint.com/docs/api/v2025/service-desk-integration
- [ ] Custom Connectors - https://developer.sailpoint.com/docs/api/v2025/connectors
- [ ] Connector Customizers - https://developer.sailpoint.com/docs/api/v2025/connector-customizers

### Low priority resources

- [ ] Branding - https://developer.sailpoint.com/docs/api/v2025/branding

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
