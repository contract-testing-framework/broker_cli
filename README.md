# README.md

The command line interface for the Signet contract testing framework.

## Cloning the repo

In your local *go* environment:

- Create a directory structure like this:
  `$GOPATH/src/github.com/contract-testing-framework`
- `cd` into `contract-testing-framework`
- `git clone` the repo
- `cd` into `broker_cli`

# Docs
------------

## `signet publish`

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
broker-url: http://localhost:3000

publish:
  type: consumer
  path: ./data_test/cons-prov.json
```

- `.signetrc.yaml` supports these flags for providers:
```yaml
broker-url: http://localhost:3000

publish:
  type: provider
  path: ./data_test/api-spec.json
  provider-name: user_service
```

#### Publishing a Contract (development)

`go run main.go publish --help` lists required flags

#### Publishing a Consumer Contract (binary - with explicit flags)

```bash
signet publish --path=./data_test/cons-prov.json --broker-url=http://localhost:3000 --type consumer
```

#### Publish a Provider Specification (binary - yaml, with explicit flags)

```bash
signet publish --path=./data_test/api-spec.yaml --broker-url=http://localhost:3000 --type provider --provider-name example-provider
```

## `signet register-env`

- The `register-env` command informs the Signet broker about a deployment environment. This is a prerequisite for allowing to the Signet broker to keep track of which versions of each service are
currently deployed to which environment.

```bash
flags:

-n --name           the name of the deployment environment being registered (ex. production)

-i --ignore-config  ingore .signetrc.yaml file if it exists

-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)
```

# Other Development Details
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
  curl -sLO https://github.com/contract-testing-framework/broker_cli/releases/download/v0.1.4/signet-darwin-arm64
  ```

- give the binary executable permissions: `chmod +x BINARY_FILE_NAME`
