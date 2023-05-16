# BADAAS: Backend And Distribution As A Service

Badaas enables the effortless construction of ***distributed, resilient, highly available and secure applications by design***, while ensuring very simple deployment and management (NoOps).

> **Warning**
> BaDaaS is still under development. Each of its components can have a different state of evolution that you can consult in [Features and components](#features-and-components)

- [BADAAS: Backend And Distribution As A Service](#badaas-backend-and-distribution-as-a-service)
  - [Features and components](#features-and-components)
  - [Quickstart](#quickstart)
    - [Example](#example)
    - [Step-by-step instructions](#step-by-step-instructions)
    - [Provided functionalities](#provided-functionalities)
      - [InfoControllerModule](#infocontrollermodule)
      - [AuthControllerModule](#authcontrollermodule)
      - [EAVControllerModule](#eavcontrollermodule)
    - [Configuration](#configuration)
  - [Contributing](#contributing)
  - [License](#license)

## Features and components

Badaas provides several key features, each provided by a component that can be used independently and has a different state of evolution:

- **Authentication**(unstable): Badaas can authenticate users using its internal authentication scheme or externally by using protocols such as OIDC, SAML, Oauth2...
- **Authorization**(wip_unstable): On resource access, Badaas will check if the user is authorized using a RBAC model.
- **Distribution**(todo): Badaas is built to run in clusters by default. Communications between nodes are TLS encrypted using [shoset](https://github.com/ditrit/shoset).
- **Persistence**(wip_unstable): Applicative objects are persisted as well as user files. Those resources are shared across the clusters to increase resiliency. To achieve this, BaDaaS uses the [BaDorm](https://github.com/ditrit/badaas/badorm) component.
- **Querying Resources**(unstable): Resources are accessible via a REST API.
- **Posix compliant**(stable): Badaas strives towards being a good unix citizen and respecting commonly accepted norms. (see [Configuration](#configuration))
- **Advanced logs management**(todo): Badaas provides an interface to interact with the logs produced by the clusters. Logs are formatted in json by default.

## Quickstart

### Example

To quickly get badaas up and running, you can head to the [example](https://github.com/ditrit/badaas-example). This example will help you to see how to use badaas and as a template to start your own project.

### Step-by-step instructions

Once you have started your project with `go init`, you must add the dependency to badaas. To use badaas, your project must also use [`fx`](https://github.com/uber-go/fx) and [`verdeter`](https://github.com/ditrit/verdeter):

<!-- TODO remove commit when badaas as a library has a first tagged version -->
```bash
go get -u github.com/ditrit/badaas@83b120f0853bce9dccb32fd27e858aa0fd71d0e6 github.com/uber-go/fx github.com/ditrit/verdeter
```

Then, your application must be defined as a `verdeter command` and you have to call the configuration of this command:

```go
var command = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
  Use:   "badaas",
  Short: "Backend and Distribution as a Service",
  Run:   runCommandFunc,
})

func main() {
  badaas.ConfigCommandParameters(command)
  command.Execute()
}
```

Then, in the Run function of your command, you must use `fx` and start the badaas functions:

```go
func runCommandFunc(cmd *cobra.Command, args []string) {
  fx.New(
    badaas.BadaasModule,

    // Here you can add the functionalities provided by badaas
    // Here you can start the rest of the modules that your project uses.
  ).Run()
}
```

You are free to choose which badaas functionalities you wish to use. To add them, you must initialise the corresponding module:

```go
func runCommandFunc(cmd *cobra.Command, args []string) {
  fx.New(
    badaas.BadaasModule,

    fx.Provide(NewAPIVersion),
    // add routes provided by badaas
    badaasControllers.InfoControllerModule,
    badaasControllers.AuthControllerModule,
    badaasControllers.EAVControllerModule,
    // Here you can start the rest of the modules that your project uses.
  ).Run()
}

func NewAPIVersion() *semver.Version {
  return semver.MustParse("0.0.0-unreleased")
}
```

Once you have defined the functionalities of your project (an http api for example), you can generate everything you need to run your application using `badctl`.

For installing it, use:

<!-- TODO remove commit when badctl has a first tagged version -->
```bash
go install github.com/ditrit/badaas/tools/badctl@cbd4c9e035709de25df59ec17e4b302b3a7b9931
```

Then generate files to make this project work with `cockroach` as database:

```bash
badctl gen --db_provider cockroachdb
```

For more information about `badctl` refer to [badctl docs](https://github.com/ditrit/badaas/tools/badctl/README.md).

Finally, you can run the api with:

```bash
badctl run
```

The api will be available at <http://localhost:8000>.

### Provided functionalities

#### InfoControllerModule

`InfoControllerModule` adds the path `/info`, where the api version will be answered. To set the version we want to be responded to we must provide the version using fx:

```go
func runCommandFunc(cmd *cobra.Command, args []string) {
  fx.New(
    badaas.BadaasModule,

    // provide api version
    fx.Provide(NewAPIVersion),
    // add /info route provided by badaas
    badaasControllers.InfoControllerModule,
  ).Run()
}

func NewAPIVersion() *semver.Version {
  return semver.MustParse("0.0.0-unreleased")
}
```

#### AuthControllerModule

`AuthControllerModule` adds `/login` and `/logout`, which allow us to add authentication to our application in a simple way:

```go
func runCommandFunc(cmd *cobra.Command, args []string) {
  fx.New(
    badaas.BadaasModule,

    // add /login and /logout routes provided by badaas
    badaasControllers.AuthControllerModule,
  ).Run()
}
```

#### EAVControllerModule

`EAVControllerModule` adds `/objects/{type}` and `/objects/{type}/{id}`, where `{type}` is any defined type and `{id}` is any uuid. These routes allow us to create, read, update and remove objects. For more information on how to use them, see the [miniblog example](https://github.com/ditrit/badaas-example).

```go
func runCommandFunc(cmd *cobra.Command, args []string) {
  fx.New(
    badaas.BadaasModule,

    // add /login and /logout routes provided by badaas
    badaasControllers.EAVControllerModule,
  ).Run()
}
```

### Configuration

Badaas use [verdeter](https://github.com/ditrit/verdeter) to manage it's configuration, so Badaas is POSIX compliant by default.

Badgen automatically generates a default configuration in `badaas/config/badaas.yml`, but you are free to modify it if you need to.

This can be done using environment variables, configuration files or CLI flags.
CLI flags take priority on the environment variables and the environment variables take priority on the content of the configuration file.

As an example we will define the `database.port` configuration key using the 3 methods:

- Using a CLI flag: `--database.port=1222`
- Using an environment variable: `export BADAAS_DATABASE_PORT=1222` (*dots are replaced by underscores*)
- Using a config file (in YAML here):

    ```yml
    # /etc/badaas/badaas.yml
    database:
      port: 1222
    ```

The config file can be placed at `/etc/badaas/badaas.yml` or `$HOME/.config/badaas/badaas.yml` or in the same folder as the badaas binary `./badaas.yml`.

If needed, the location can be overridden using the config key `config_path`.

***For a full overview of the configuration keys: please head to the [configuration documentation](./configuration.md).***

## Contributing

See [this section](./CONTRIBUTING.md).

## License

Badaas is Licensed under the [Mozilla Public License Version 2.0](./LICENSE).
