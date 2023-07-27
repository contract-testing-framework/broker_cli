
`signet-cli` is the command line interface for the Signet contract testing framework. The functionality and documentation for the cli lives here. The [signet-broker README](https://github.com/signet-framework/signet-broker) has example use-cases, and more information about the Signet broker.
&nbsp;  

# Installation

- requires node and npm

Install:
```bash
npm install -g signet-cli
```

Uninstall:
```bash
uninstall-signet-cli
```
&nbsp;  
# `signet` commands
Syntax:
```bash
signet [command] [--flags]
```
&nbsp;  

Signet-cli checks for a `.signetrc.yaml` file in the current working directory. All required flags, and most optional flags can be set in the config file instead of being passed on the command line.

Config file syntax:
```yaml
global-flag: string

any-signet-command:
  flag-for-command: string
```
&nbsp;  
## `signet deploy`

- The `deploy` command deploys the Signet broker to a new ECS Fargate cluster on AWS. Using the golang SDK and a CloudFormation template, `deploy` provisions all of the necessary AWS infrastructure for self-hosting the Signet framework. When the broker is successfully deployed, the command outputs the public URL for the load balancer serving as the entrypoint to the application. 

- `deploy` requires the `aws cli` to be installed, and an IAM user with the `AdministratorAccess` permission set to be logged in as the `default` CLI profile.

```bash
signet deploy
```
&nbsp;  
## `signet undeploy`
- The `undeploy` command tears down all of the cloud infrastructure created by `signet deploy`
```bash
signet undeploy
```
&nbsp;  
## `signet proxy`

- The `proxy` command is used to automatically generate a consumer contract by recording requests and responses generated during unit and service tests. `proxy` starts up a server that acts as a transparent proxy between the consumer service under test and the mock or stub of the provider service. `proxy` captures the requests and responses between the two services, and automatically generates a valid consumer contract. `proxy` uses an open source tool (mountebank) to record the requests and responses, and then transforms the recorded messages into a Pact-complient consumer contract.

```bash
signet proxy


flags:

-o --port           the port that signet proxy should run on

-t --target         the URL of the running provider stub or mock

-p --path           the relative path and filename that the consumer contract will be written to

-n -—name           the canonical name of the consumer service

-m --provider-name  the canonical name of the provider service that the mock or stub represents

-i --ignore-config  ingore .signetrc.yaml file if it exists (optional)
```
- `.signetrc.yaml` supports these flags for `signet proxy`:
```yaml
broker-url: http://localhost:3000

proxy:
  path: ./contracts/cons-prov.json
  port: 3004
  target: http://localhost:3002
  name: service_1
  provider-name: user_service
```
&nbsp;  
## `signet publish`
- The `publish` command pushes a local contract or API spec to the broker. This automatically triggers contract/spec comparison if the broker already has a contract or API spec for the other participant in the integration.

- When publishing a consumer contract, it required to pass a `--version`. This informs the Signet broker of which versions of the consumer service the consumer contract is tested against.

- When publishing a provider API spec, `--version` and `--branch` flags are ignored. This is becuase a provider spec is not generated from unit tests (like a consumer contract), and is not guarenteed to be correctly implemented by a provider at the time the spec is published. Versions of a provider service are proven to correctly implement an API spec with the `signet test` command. A passing `signet test` will inform the Signet broker of which versions of the provider service are tested against the API spec.

```bash
signet publish


flags:

-p --path           the relative path to the contract or API spec

-t -—type           the type of service contract (either 'consumer' or 'provider')

-n -—name           canonical name of the provider service (only for —-type 'provider')

-v -—version        service version (only for --type 'consumer', defaults to the git SHA of HEAD if no value is provided)

-b -—branch         git branch name (optional, only for --type 'consumer', defaults to git branch of HEAD if no value is provided)

-u --broker-url     the scheme, domain, and port where the Signet broker is being hosted

-i --ignore-config  ingore .signetrc.yaml file if it exists (optional)
```

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
&nbsp;  
## `signet test`
- The `test` command determines if a provider service correctly implements an API spec. First, it fetches the latest API spec from the Signet broker. Then, it leverages an open source tool (dredd) to parse the API spec, generate mock requests and expected responses, and execute those interactions against the provider service. If the tests are successful, `test` notifies the Signet broker that this version of the provider service is verified -- it is proven to implement the API spec through testing. If any tests fail, an analysis of the failing tests is logged.

- Before running `test`, the provider service must be running, and an API spec for that service must be published to the Signet broker.

```bash
signet test


flags:

-n --name           the name of the provider service

-v --version        the version of the provider service (defaults to git SHA of HEAD if no value is provided)

-b --branch         git branch (optional, defaults to git branch of HEAD if '--branch' is passed with no value, or if '--version' defaulted to git SHA)

-s --provider-url   the URL where the provider service is running

-u --broker-url     the scheme, domain, and port where the Signet broker is being hosted

-i --ignore-config  ingore .signetrc.yaml file if it exists (optional)
```

- `.signetrc.yaml` supports these flags for `signet test`:
```yaml
broker-url: http://localhost:3000

test:
  name: user_service
  provider-url: http://localhost:3002
```
&nbsp;  
## `signet register-env`

- The `register-env` command informs the Signet broker about a new deployment environment. 

```bash
signet register-env


flags:

-e --environment    the name of the deployment environment being registered (ex. production)

-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted

-i --ignore-config  ingore .signetrc.yaml file if it exists (optional)
```
- `.signetrc.yaml` supports these flags for `update-deployment`:
```yaml
broker-url: http://localhost:3000

register-env:
  environment: production
```
&nbsp;  
## `signet update-deployment`

- The `update-deployment` command informs the Signet broker of which service versions are currently deployed in an environment. If broker does not already know about the `--environment`, it will create it.

```bash
signet update-deployment


flags:

-n --name           the name of the service

-v --version        the version of the service (defaults to git SHA of HEAD if no value is provided)

-e --environment    the name of the environment that the service is deployed to (ex. production)

-d --delete         the presence of this flag indicates that the service is no longer deployed to the environment (optional)

-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted

-i --ignore-config  ingore .signetrc.yaml file if it exists (optional)
```
- `.signetrc.yaml` supports these flags for `update-deployment`:
```yaml
broker-url: http://localhost:3000

update-deployment:
  name: user_service
  environment: production
```
&nbsp;  
## `signet deploy-guard`
- The `deploy-guard` command checks whether a service version can be safely deployed to an environment without introducing any breakages with other services in that environemnt. The `deploy-guard` command will fail (with an exit code of 1) if ANY of the following conditions are NOT met: 

1. The service is compatible with all of its consumers which are deployed in the environment
2. All of the service's providers are deployed in the environment
3. All of the service's providers are compatible with the service.

- If any of these are not true, the service version cannot be safely deployed to the environemnt, because doing so would either break the service or break one of its consumers. `deploy-guard` allows a CI/CD pipeline to automatically gate a deployment if it will lead to unintended breakages.
	
```bash
signet deploy-guard


flags:

-n --name           the name of the service

-v --version        the version of the service (defaults to git SHA of HEAD if no value is provided)

-e --environment    the name of the environment that the service is deployed to (ex. production)

-u --broker-url     the scheme, domain, and port where the Signet Broker is being hosted

-i --ignore-config  ingore .signetrc.yaml file if it exists (optional)
```
- `.signetrc.yaml` supports these flags for `deploy-guard`:
```yaml
broker-url: http://localhost:3000

deploy-guard:
  name: user_service
```