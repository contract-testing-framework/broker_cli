# README.md

The command line interface for the Signet contract testing framework.

# Installation

## Install npm pkg (full featured, requires node and npm)

```bash
npm install -g test_signet_cli
```

## Install only the Signet CLI golang binary (does not support provider verification)

MacOS arm64
```bash
curl -sLO https://github.com/contract-testing-framework/broker_cli/releases/download/v0.3.0/signet-darwin-arm64 \
&& mv signet-darwin-arm64 signet \
&& chmod +x signet \
&& export PATH=$PATH:$(pwd)
```

MacOS amd64 (x86-64)
```bash
curl -sLO https://github.com/contract-testing-framework/broker_cli/releases/download/v0.3.0/signet-darwin-amd64 \
&& mv signet-darwin-amd64 signet \
&& chmod +x signet \
&& export PATH=$PATH:$(pwd)
```

Linux amd64 (x86-64)
```bash
curl -sLO https://github.com/contract-testing-framework/broker_cli/releases/download/v0.3.0/signet-linux-amd64 \
&& mv signet-linux-amd64 signet \
&& chmod +x signet \
&& export PATH=$PATH:$(pwd)
```

# Docs

Every signet command supports `--help` flag, for example:
`signet publish --help`

## `signet publish`

- The `publish` command pushes a local contract or API spec to the broker. This automatically triggers contract/spec comparison if the broker already has a contract or API spec for the other participant in the integration.

- When publishing a consumer contract, it required to pass a `--version`. This is used to inform the Signet broker of which versions of the consumer service the consumer contract is tested against.

- When publishing a provider spec, `--version` and `--branch` flags are ignored. Versions of a provider service are proven to correctly implement an API spec with the `signet test` command. A passing `signet test` will inform the Signet broker of which versions of the provider service are tested against the API spec.

```bash
flags:

-p --path           the relative path to the contract or spec

-t -—type           the type of service contract (either 'consumer' or 'provider')

-n -—name           canonical name of the provider service (only for —-type 'provider')

-v -—version        service version (only for --type 'consumer', if flag not passed or passed without value, defaults to the git SHA of HEAD)

-b -—branch         git branch name (optional, only for --type 'consumer', defaults to git branch of HEAD)

-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)

-i --ignore-config  ingore .signetrc.yaml file if it exists
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
  name: user_service
```

#### Publishing a Consumer Contract (with explicit flags)

```bash
signet publish --path=./data_test/cons-prov.json --broker-url=http://localhost:3000 --type consumer
```

#### Publish a Provider Specification (with explicit flags)

```bash
signet publish --path=./data_test/api-spec.yaml --broker-url=http://localhost:3000 --type provider --name=example-provider
```

## `signet test`
- The `signet test` command determines if a provider service correctly implements an API spec. First, it fetches the most recently updated API spec from the Signet broker. Then, it leverages an open source tool (dredd) to parse the API spec, generate mock requests and expected responses, and execute those interactions against the provider service. If the tests are successful, `signet test` notifies the Signet broker that this version of the provider service is verified -- it is proven to implement the API spec through testing. If any tests fail, an analysis of the failing tests is logged.

- Before running `signet test`, the provider service must be running, and an API spec for that service must be published to the Signet broker.

```bash
flags:

-n --name           the name of the provider service

-v --version        the version of the provider service (required, passing --version without a value will default to git SHA of HEAD)

-b --branch         version control branch (passing --branch without a value will default to git branch of HEAD)

-s --provider-url   the URL where the provider service is running

-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)

-i --ignore-config  ingore .signetrc.yaml file if it exists (optional)
```

- `.signetrc.yaml` supports these flags for `signet test`:
```yaml
broker-url: http://localhost:3000

test:
  name: user_service
  provider-url: http://localhost:3002
```

#### Test a provider service against an API spec (with explicit flags)
```bash
signet test --broker-url=http://localhost:3000 --provider-url=http://localhost:3002 --name=example-provider --version=version1 --branch
```

## `signet update-deployment`

- The `update-deployment` command informs the Signet broker of which service versions are currently deployed in an environment. If broker does not already know about the `--environment`, it will create it.

```bash
flags:

-n --name           the name of the service

-v --version        the version of the service (required)

-e --environment    the name of the environment that the service is deployed to (ex. production)

-d --delete         the presence of this flag inidicates that the service is no longer deployed to the environment

-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)

-i --ignore-config  ingore .signetrc.yaml file if it exists
```
- `.signetrc.yaml` supports these flags for `update-deployment`:
```yaml
broker-url: http://localhost:3000

update-deployment:
  name: user_service
```

#### Notify the Signet broker of a deployment
```bash
signet update-deployment --broker-url=http://localhost:3000 --name=example-provider --version=version1 --environment=production
```

# Development Details
## Cloning the repo

In your local `go` environment:

- Create a directory structure like this:
  `$GOPATH/src/github.com/contract-testing-framework`
- `cd` into `contract-testing-framework`
- `git clone` the repo
- `cd` into `broker_cli`

## Run the CLI in development
`go run main.go [cmd]`

## Run the test suite
`make test`


## Releasing a new version of the Signet CLI (as an NPM package)

#### In the Signet CLI project
- commit changes locally
- add an annotated tag with a new semantic version (check the 'Releases' page in the GitHub repo for the last version number)
    `git tag -a v0.3.10 -m "updated changes"`
- push the commit and tag to GitHub
    `git push origin v0.3.10`
- use goreleaser to publish binaries to 'Releases' on GitHub (this automatically builds new binaries before releasing)
    `goreleaser release --clean`

#### In the cli_npm_pkg project
- delete the outdated binaries
    `rm -rf dist/`

#### In the Signet CLI project
- copy the new binaries over to wherever your local cli_npm_pkg root directory is
    `cp -r dist <relative path to your cli_npm_pkg root>`

#### In the cli_npm_pkg project
- open up package.json, and change the `"version"` to the new semantic version
    `"version": 0.3.10`
- publish the updated npm package
    `npm publish`