# Developing

This document provides the information you need to know before developing code for a pull request.

## Environment

- Install [go](https://go.dev/doc/install) >= v1.20
- Install project dependencies: `go get`
- Install [docker](https://docs.docker.com/engine/install/) and [compose plugin](https://docs.docker.com/compose/install/)

## Directory structure

This is the directory structure we use for the project:

- `badorm/`: Contains the code of the BaDORM component.
- `configuration/`: Contains all the configuration holders. Please only use the interfaces, they are all mocked for easy testing.
- `controllers/`: Contains all the http controllers, they handle http requests and consume services.
- `docker/` : Contains the docker, docker-compose and configuration files for different environments.
- `docs/`: Contains the documentation showed for readthedocs.io.
- `httperrors/`: Contains the http errors that can be responded by the http api. Should be moved to `controller/` when services stop using them.
- `logger/`: Contains the logger creation logic. Please don't call it from your own services and code, use the dependency injection system.
- `mocks/`: Contains the mocks generated with `mockery`.
- `persistance/`:
  - `gormdatabase/`: Contains the logic to create a <https://gorm.io> database. Also contains a go package named `gormzap`: it is a compatibility layer between *gorm.io/gorm* and *github.com/uber-go/zap*.
  - `models/`: Contains the models.
    - `dto/`: Contains the Data Transfer Objects. They are used mainly to decode json payloads.
  - `repository/`: Contains the repository interfaces and implementations to manage queries to the database.
- `router/`: Contains http router of badaas and the routes that can be added by the user.
  - `middlewares/`: Contains the various http middlewares that we use.
- `services/`: Contains services.
  - `auth/protocols/`: Contains the implementations of authentication clients for different protocols.
    - `basicauth/`: Handle the authentication using email/password.
    - `oidc/`: Handle the authentication via Open-ID Connect.
- `test_e2e/`: Contains all the feature and steps for e2e tests.
- `testintegration/`: Contains all the integration tests.
- `tools/`: Contains the go tools necessary to use BaDaaS.
  - `badctl`: Contains the command line tool that makes it possible to configure and run BaDaaS applications easily.
- `utils/`: Contains functions that can be util all around the project, as managing data structures, time, etc.

At the root of the project, you will find:

- The README.
- The changelog.
- The LICENSE file.

## Tests

### Dependencies

Running tests have some dependencies as: `mockery`, `gotestsum`, etc.. Install them with `make install_dependencies`.

### Linting

We use `golangci-lint` for linting our code. You can test it with `make lint`. The configuration file is in the default path (`.golangci.yml`). The file `.vscode.settings.json.template` is a template for your `.vscode/settings.json` that formats the code according to our configuration.

### Unit tests

We use the standard test suite in combination with [github.com/stretchr/testify](https://github.com/stretchr/testify) to do our unit testing. Mocks are generated using [mockery](https://github.com/vektra/mockery) a mock generator using the command `make test_generate_mocks`.

To run them, please run:

```sh
make -k test_unit
```

### Integration tests

Integration tests have a database and the dependency injection system. BaDaaS and BaDORM are tested on multiple databases (those supported by gorm, which is the base of BaDORM). By default, the database used will be postgresql:

```sh
make test_integration
```

To run the tests on another database you can use: `make test_integration_postgresql`, `make test_integration_cockroachdb`, `make test_integration_mysql`, `make test_integration_sqlite`, `make test_integration_sqlserver`. All of them will be verified by our continuous integration system.

### Feature tests (end to end tests)

We use docker to run a Badaas instance in combination with one node of CockroachDB.

Run:

```sh
make test_e2e
```

The feature files can be found in the `test_e2e/features` folder.

## Use of Third-party code

Third-party code must include licenses.
