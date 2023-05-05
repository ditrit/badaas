# BADAAS: Backend And Distribution As A Service

Badaas enables the effortless construction of ***distributed, resilient, highly available and secure applications by design***, while ensuring very simple deployment and management (NoOps).

Badaas provides several key features:

- **Authentication**: Badaas can authenticate users using its internal authentication scheme or externally by using protocols such as OIDC, SAML, Oauth2...
- **Authorization**: On resource access, Badaas will check if the user is authorized using a RBAC model.
- **Distribution**: Badaas is built to run in clusters by default. Communications between nodes are TLS encrypted using [shoset](https://github.com/ditrit/shoset).
- **Persistence**: Applicative objects are persisted as well as user files. Those resources are shared across the clusters to increase resiliency.
- **Querying Resources**: Resources are accessible via a REST API.
- **Posix compliant**: Badaas strives towards being a good unix citizen and respecting commonly accepted norms. (see [Configuration](#configuration))
- **Advanced logs management**: Badaas provides an interface to interact with the logs produced by the clusters. Logs are formatted in json by default.

- [BADAAS: Backend And Distribution As A Service](#badaas-backend-and-distribution-as-a-service)
  - [Quickstart](#quickstart)
    - [Example](#example)
    - [Step-by-step instructions](#step-by-step-instructions)
  - [Configuration](#configuration)
  - [Contributing](#contributing)
  - [License](#license)

## Quickstart

### Example

To quickly get badaas up and running, you can head to the [miniblog example](https://github.com/ditrit/badaas-example). This example will help you to see how to use badaas and as a template to start your own project

### Step-by-step instructions

Once you have started your project with `go init`, you must add the dependency to badaas. To use badaas, your project must also use [`fx`](https://github.com/uber-go/fx) and [`verdeter`](https://github.com/ditrit/verdeter):

<!-- TODO remove commit when badaas as a library has a first tagged version -->
```bash
go get -u github.com/ditrit/badaas@dbd7e55
go get -u github.com/uber-go/fx
go get -u github.com/ditrit/verdeter
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

    // add routes provided by badaas
    fx.Invoke(router.AddInfoRoutes),
    fx.Invoke(router.AddLoginRoutes),
    fx.Invoke(router.AddCRUDRoutes),
    // Here you can start the rest of the modules that your project uses.
  ).Run()
}
```

Once you have defined the functionalities of your project (an http api for example), you can generate everything you need to run your application using `badctl`.

For installing it, use:

<!-- TODO remove commit when badctl has a first tagged version -->
```bash
go install github.com/ditrit/badaas/tools/badctl@bef1116
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
