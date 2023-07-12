# README.md

The command line interface for the Signet contract testing framework.

# Installation

- requires node and npm

```bash
npm install -g signet-cli
```

# Docs

Every signet command supports `--help` flag, for example:
`signet publish --help`

## `signet deploy`

- The `signet deploy` command deploys the Signet broker to the users AWS account using AWS ECS with AWS Fargate. The Signet broker can be torn down by calling `signet deploy` again with the `-d` flag.
- requies local docker engine to be running, and for AWS credentials to be available on a docker ecs context (this can be accomplished with `docker create context ecs [context name]`)

```bash
flags:

-c --ecs-context    the name of the local docker ecs context with AWS credentials

-s -—silent         (bool) silence docker's status updates as it provisions AWS infrastructure

-d --destroy        (bool) causes the Signet broker to be torn down from AWS instead of deployed
```

- `.signetrc.yaml` supports these flags for `signet proxy`:
```yaml
deploy:
  ecs-context: myecscontext
  silent: true
```

#### Deploying the Signet Broker (with explicit flags)
```bash
signet deploy --ecs-context myecscontext
```

## `signet proxy`

- The `signet proxy` command is used to automatically generate a consumer contract by recording requests and responses generated during unit and service tests. `signet proxy` starts up a server that acts as a transparent proxy between the consumer service under test and the mock or stub of the provider service. `signet proxy` captures the requests and responses between the two, and automatically generates a valid consumer contract. `signet proxy` uses Mountebank to record the messages, and then transforms Mountebank's output into a Pact-complient consumer contract.

```bash
flags:

-o --port           the port that signet proxy should run on

-t --target         the URL of the running provider stub or mock

-p --path           the relative path and filename that the consumer contract will be written to

-n -—name           the canonical name of the consumer service

-m --provider-name  the canonical name of the provider service that the mock or stub represents

-i --ignore-config  ingore .signetrc.yaml file if it exists
```
- `.signetrc.yaml` supports these flags for `signet proxy`:
```yaml
proxy:
  path: ./contracts/cons-prov.json
  port: 3004
  target: http://localhost:3002
  name: service_1
  provider-name: user_service
```
#### Using Signet Proxy (with explicit flags)
```bash
signet proxy --port 3005 --target http://localhost:3002 --path ./contracts/contract.json --name service_1 --provider-name user_service
```

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

-v --version        the version of the provider service (defaults to git SHA of HEAD)

-b --branch         version control branch (passing --branch without a value will default to git branch of HEAD, also defaults if --version)

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

## `signet deploy-guard`
- The `deploy-guard` command checks whether a service version can be safely deployed to an environment without introducing any breakages with other services in that environemnt. The `deploy-guard` command will fail (with an exit code of 1) if any of the following conditions are NOT met: 

1. The service is compatible with all of its consumers which are deployed in the environment
2. All of the service's providers are deployed in the environment
3. All of the service's providers are compatible with the service.

- If any of these are not true, the service version cannot be safely deployed to the environemnt, because doing so would either break the service or break one of its consumers.
	
```bash
flags:

-n --name 					the name of the service

-v --version        the version of the service (defaults to git SHA of HEAD)

-e --environment		the name of the environment that the service is deployed to (ex. production)

-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted (ex. http://localhost:3000)

-i --ignore-config  ingore .signetrc.yaml file if it exists
```
- `.signetrc.yaml` supports these flags for `deploy-guard`:
```yaml
broker-url: http://localhost:3000

deploy-guard:
  name: user_service
```
#### Check if it is safe to deploy a new version of a service
```bash
signet deploy-guard --broker-url=http://localhost:3000 --name=example-provider --version=version1 --environment=production
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
