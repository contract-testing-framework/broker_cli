A command line interface for the contract testing broker.

# Cloning the repo

In your local go environemnt:
- create a directory structure like this:
  `$GOPATH/src/github.com/contract-testing-framework`
- `cd` into `contract-testing-framework`
- `git clone` the repo
- `cd` into `broker_cli`

# Building a new version of the binary executable

goreleaser --snapshot

# Documentation

## Publishing a Contract (in development)

`go run main.go publish --help` lists required arguments and flags

Publish an example provider specification (yaml):
`go run main.go publish ./data_test/api-spec.yaml http://localhost:3000/api/contracts --type provider --provider-name example-provider`