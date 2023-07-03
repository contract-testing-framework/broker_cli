# README.md

`broker_cli`, command line interface for the contract testing broker.

## Cloning the repo

In your local *go* environment:

- Create a directory structure like this:
  `$GOPATH/src/github.com/contract-testing-framework`
- `cd` into `contract-testing-framework`
- `git clone` the repo
- `cd` into `broker_cli`

## Docs

### broker_cli publish

- The `publish` command pushes a local contract or spec to the broker. This automatically triggers contract/spec comparison if the broker already has a contract for the other participant in the integration.

```bash
flags:

-i --ignore-config  ingore .signetrc.yaml file if it exists

-p --path           the relative path to the contract or spec

-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)

-t -—type           the type of service contract (either 'consumer' or 'provider')

-n -—provider-name  canonical name of the provider service (only for —-type 'provider')

-v -—version        service version (required for --type 'consumer')
                    -—type=consumer: if flag not passed or passed without value, defaults to the git SHA of HEAD
                    -—type=provider: if the flag passed without value, defaults to git SHA

-b -—branch         git branch name (optional, defaults to current git branch)
```

- if a `.signetrc.yaml` file is present in the root directory, the broker_cli will read values for flags from it. Any flags which are explicitly passed on the command line will override the values in `.signetrc.yaml`.

- `.signetrc.yaml` supports these flags for consumers:
```yaml
type: consumer
path: ./data_test/cons-prov.json
broker-url: http://localhost:3000
```

- `.signetrc.yaml` supports these flags for providers:
```yaml
type: provider
path: ./data_test/api-spec.json
broker-url: http://localhost:3000
provider-name: user_service
```

### Publishing a Contract (in development)

`go run main.go publish --help` lists required arguments and flags

#### Publishing a Consumer Contract (with explicit flags)

```bash
broker_cli publish --path=./data_test/cons-prov.json --broker-url=http://localhost:3000 --type consumer
```

#### Publish a Provider Specification (yaml, with explicit flags)

```bash
broker_cli publish --path=./data_test/api-spec.yaml --broker-url=http://localhost:3000 --type provider --provider-name example-provider
```

## Release updated binaries

- build new binaries with `make build`
- create a new semantic version tag before committing: `git tag v0.1.4`
- commit changes
- push changes to github
- manually upload the binaries through the github releases page:
  - from the main `Code` tab, click on `Releases` in the right-hand sidebar
  - click `Draft a new release`
  - add the semantic version tag for the commit
  - upload binaries
  - click `set as latest release`
  - click `Publish release`

## Install binary executables

- go to the `Releases` page in the repo
- right click on the binary for your OS/Arch and copy the link address
- in the directory where you want to keep the binary, run `curl -sLO` followed by the link address.
  - ex.

  ```bash
  curl -sLO https://github.com/contract-testing-framework/broker_cli/releases/download/v0.1.4/broker_cli-darwin-arm64
  ```

- give the binary executable permissions: `chmod +x BINARY_FILE_NAME`
