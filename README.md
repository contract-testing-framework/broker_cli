A command line interface for the contract testing broker.

# cloning the repo

In your local go environemnt:
- create a directory structure like this:
  `$GOPATH/src/github.com/contract-testing-framework`
- `cd` into `contract-testing-framework`
- `git clone` the repo
- `cd` into `broker_cli`


# Docs

### broker_cli publish
- The `publish` command pushes a local contract to the contract broker. This automatically triggers contract verification if the broker
already has a contract for the other microservice in the integration.

args:

`publish [path to contract file] [broker-url]`

flags:

-t -—type           enum(’consumer’, ‘provider’)

-v -—version        (optional, defaults to SHA of git commit)

-b -—branch         (optional)

-n -—provider-name  (only for -—type ‘provider’)

# Publishing a Contract (in development)

`go run main.go publish --help` lists required arguments and flags

Publish an example provider specification (yaml):
`go run main.go publish ./data_test/api-spec.yaml http://localhost:3000/api/contracts --type provider --provider-name example-provider`