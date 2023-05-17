# BaDorm: Backend and Distribution ORM (Object Relational Mapping)

BaDorm is the BaDaaS component that allows for easy persistence and querying of objects. It is built on top of gorm and adds for each entity a service and a repository that allows complex queries without any extra effort.

BaDorm can be used both within a BaDaaS application and as a stand-alone application.

- [BaDorm: Backend and Distribution ORM (Object Relational Mapping)](#badorm-backend-and-distribution-orm-object-relational-mapping)
  - [Quickstart](#quickstart)
    - [Stand-alone Example](#stand-alone-example)
    - [BaDaaS Example](#badaas-example)
    - [Step-by-step instructions](#step-by-step-instructions)
  - [Provided functionalities](#provided-functionalities)
    - [CRUDServiceModule](#crudservicemodule)

## Quickstart

### Stand-alone Example

To quickly understand the features provided by BaDorm, you can head to the [example](https://github.com/ditrit/badorm-example). This example will help you to see how to use BaDorm and as a template to start your own project.

### BaDaaS Example

If you are interested in using BaDorm within a BaDaaS application you can consult the [example](https://github.com/ditrit/badaas-example) in which besides using the services and repositories provided by BaDorm, BaDaaS adds a controller that allows the query of objects via an http api.

### Step-by-step instructions

Once you have started your project with `go init`, you must add the dependency to BaDaaS:

<!-- TODO remove commit when badaas as a library has a first tagged version -->
```bash
go get -u github.com/ditrit/badaas@7fae89e
```

In order to use BaDorm you will also need to use the following libraries:

```bash
go get -u github.com/uber-go/fx github.com/uber-go/zap gorm.io/gorm
```

First of all, you will need to start your application with `fx`:

```go
func main() {
  fx.New(
    fx.Provide(NewLogger),

    // DB modules
    fx.Provide(NewGORMDBConnection),
    // activate BaDORM
    badorm.BaDORMModule,
    // start example data
    badorm.GetCRUDServiceModule[models.Company, uuid.UUID](),
    badorm.GetCRUDServiceModule[models.Product, uuid.UUID](),
    badorm.GetCRUDServiceModule[models.Seller, uuid.UUID](),
    badorm.GetCRUDServiceModule[models.Sale, uuid.UUID](),
    fx.Provide(CreateCRUDObjects),
    fx.Invoke(QueryCRUDObjects),
  ).Run()
}
```

There are some things you need to provide to the BaDorm module:

- `NewLogger` is the function that provides a zap logger to the BaDorm components.
- `NewGORMDBConnection` if the function that establish the connection to the database where you data will be saved.

After that, you can start the `badorm.BaDORMModule` and crete the CRUD services to your models using `badorm.GetCRUDServiceModule`.

Finally, you can call your application functions as `CreateCRUDObjects` and `QueryCRUDObjects` where created  CRUDServices can be injected to create, read, update and delete your models easily.

## Provided functionalities

### CRUDServiceModule

`CRUDServiceModule` provides you a CRUDService, a CRUDRepository for your model class and registers it. After calling it as, for example, `badorm.GetCRUDServiceModule[models.Company, uuid.UUID](),` the following can be used by dependency injection:

- `crudCompanyService badorm.CRUDService[models.Company, uuid.UUID]`
- `crudCompanyRepository badorm.CRUDRepository[models.Company, uuid.UUID]`
