# Contribute to the development of badaas

- [Contribute to the development of badaas](#contribute-to-the-development-of-badaas)
  - [Tests](#tests)
    - [Dependencies](#dependencies)
    - [Linting](#linting)
    - [Unit tests](#unit-tests)
    - [Integration tests](#integration-tests)
    - [Feature tests (or end to end tests)](#feature-tests-or-end-to-end-tests)
  - [Logger](#logger)
  - [Directory structure](#directory-structure)
  - [Git](#git)
    - [Branch naming policy](#branch-naming-policy)
    - [Default branch](#default-branch)
    - [How to release](#how-to-release)

## Tests

### Dependencies

Running tests have some dependencies as: `mockery`, `gotestsum`, etc.. Install them with `make install dependencies`.

### Linting

We use `golangci-lint` for linting our code. You can test it with `make lint`. The configuration file is in the default path (`.golangci.yml`). The file `.vscode.settings.json.template` is a template for your `.vscode/settings.json` that formats the code according to our configuration.

### Unit tests

We use the standard test suite in combination with [github.com/stretchr/testify](https://github.com/stretchr/testify) to do our unit testing. Mocks are generated using [mockery](https://github.com/vektra/mockery) a mock generator using this command `make test_generate_mocks`.

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

### Feature tests (or end to end tests)

We use docker to run a Badaas instance in combination with one node of CockroachDB.

Run:

```sh
make test_e2e
```

The feature files can be found in the `test_e2e/features` folder.

## Logger

We use uber's [zap](https://pkg.go.dev/go.uber.org/zap) to log stuff, please take `zap.Logger` as an argument for your services constructors. [fx](https://github.com/uber-go/fx) will provide your service with an instance.

## Directory structure

This is the directory structure we use for the project:

- `configuration/` *(Go code)*: Contains all the configuration holders. Please only use the interfaces, they are all mocked for easy testing.
- `controllers/` *(Go code)*: Contains all the http controllers, they handle http requests and consume services.
- `docs/`: Contains the documentation.
- `test_e2e/`: Contains all the feature and steps for e2e tests.
- `test_integration/`: Contains all the integration tests.
- `logger/` *(Go code)*: Contains the logger creation logic. Please don't call it from your own services and code, use the dependency injection system.
- `persistance/` *(Go code)*:
  - `gormdatabase/` *(Go code)*: Contains the logic to create a <https://gorm.io> database. Also contains a go package named `gormzap`: it is a compatibility layer between *gorm.io/gorm* and *github.com/uber-go/zap*.
  - `models/` *(Go code)*: Contains the models. (For a structure to me considered a valid model, it has to embed `badorm.UUIDModel` or `badorm.UIntModel`.).
    - `dto/` *(Go code)*: Contains the Data Transfer Objects. They are used mainly to decode json payloads.
  - `repository/` *(Go code)*: Contains the repository interfaces and implementations to manage queries to the database.
- `router/` *(Go code)*: Contains http router of badaas and the routes that can be added by the user.
  - `middlewares/` *(Go code)*: Contains the various http middlewares that we use.
- `docker/` : Contains the docker, docker-compose file and configuration files for different environments.
- `services/` *(Go code)*: Contains services.
  - `auth/protocols/`: Contains the implementations of authentication clients for different protocols.
    - `basicauth/` *(Go code)*: Handle the authentication using email/password.
    - `oidc/` *(Go code)*: Handle the authentication via Open-ID Connect.
  - `sessionservice/` *(Go code)*: Handle sessions and their lifecycle.
  - `userservice/` *(Go code)*: Handle users.
- `validators/` *(Go code)*: Contains validators such as an email validator.
- `badorm/` *(Go code)*: Contains the code of the orm used by badaas.

At the root of the project, you will find:

- The README.
- The changelog.
- The LICENSE file.

## Git

### Branch naming policy

`[BRANCH_TYPE]/[BRANCH_NAME]`

- `BRANCH_TYPE` is a prefix to describe the purpose of the branch.
  Accepted prefixes are:
  - `feature`, used for feature development
  - `bugfix`, used for bug fix
  - `improvement`, used for refactor
  - `library`, used for updating library
  - `prerelease`, used for preparing the branch for the release
  - `release`, used for releasing project
  - `hotfix`, used for applying a hotfix on main
  - `poc`, used for proof of concept
- `BRANCH_NAME` is managed by this regex: `[a-z0-9._-]` (`_` is used as space character).

### Default branch

The default branch is `main`. Direct commit on it is forbidden. The only way to update the application is through pull request.

Release tag are only done on the `main` branch.

### How to release

We use [Semantic Versioning](https://semver.org/spec/v2.0.0.html) as guideline for the version management.

Steps to release:

- Create a new branch labeled `release/vX.Y.Z` from the latest `main`.
- Improve the version number in `changelog.md` and `resources/api.go`.
- Verify the content of the `changelog.md`.
- Commit the modifications with the label `Release version X.Y.Z`.
- Create a pull request on github for this branch into `main`.
- Once the pull request validated and merged, tag the `main` branch with `vX.Y.Z`.
- After the tag is pushed, make the release on the tag in GitHub.
