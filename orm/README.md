# BaDaaS ORM: Backend and Distribution ORM (Object Relational Mapping) <!-- omit in toc -->

BaDaaS ORM is the BaDaaS component that allows for easy persistence and querying of objects. It is built on top of gorm and adds for each entity a service and a repository that allows complex queries without any extra effort.

BaDaaS ORM can be used both within a BaDaaS application and as a stand-alone application.

- [Quickstart](#quickstart)
  - [Stand-alone Example](#stand-alone-example)
  - [BaDaaS Example](#badaas-example)
  - [Step-by-step instructions](#step-by-step-instructions)
- [Provided functionalities](#provided-functionalities)
  - [Base models](#base-models)
  - [CRUDServiceModule](#crudservicemodule)

## Quickstart

### Stand-alone Example

To quickly understand the features provided by badaas-orm, you can head to the [example](https://github.com/ditrit/badaas-orm-example). This example will help you to see how to use badaas-orm and as a template to start your own project.

### BaDaaS Example

If you are interested in using badaas-orm within a BaDaaS application you can consult the [example](https://github.com/ditrit/badaas-example) in which besides using the services and repositories provided by badaas-orm, BaDaaS adds a controller that allows the query of objects via an http api.

### Step-by-step instructions

Once you have started your project with `go init`, you must add the dependency to BaDaaS:

```bash
go get -u github.com/ditrit/badaas
```

In order to use badaas-orm you will also need to use the following libraries:

```bash
go get -u github.com/uber-go/fx github.com/uber-go/zap gorm.io/gorm
```

First of all, you will need to start your application with `fx`:

```go
func main() {
  fx.New(
    // connect to db
    fx.Provide(NewGormDBConnection),
    // activate badaas-orm
    fx.Provide(GetModels),
    orm.AutoMigrate,

    // create crud services for models
    orm.GetCRUDServiceModule[models.Company](),
    orm.GetCRUDServiceModule[models.Product](),
    orm.GetCRUDServiceModule[models.Seller](),
    orm.GetCRUDServiceModule[models.Sale](),

    // start example data
    fx.Provide(CreateCRUDObjects),
    fx.Invoke(QueryCRUDObjects),
  ).Run()
}
```

There are some things you need to provide to the badaas-orm module:

- `NewGORMDBConnection` is the function that establish the connection to the database where you data will be saved.
- `GetModels` is the function that returns a `orm.GetModelsResult`, to tell badaas-orm which are the models you want to be auto-migrated.

After that, you can execute the auto-migration with `orm.AutoMigrate` and create the CRUD services to your models using `orm.GetCRUDServiceModule`.

Finally, you can call your application functions as `CreateCRUDObjects` and `QueryCRUDObjects` where created CRUDServices can be injected to create, read, update and delete your models easily.

## Provided functionalities

### Base models

badaas-orm gives you two types of base models for your classes: `orm.UUIDModel` and `orm.UIntModel`.

To use them, simply embed the desired model in any of your classes:

```go
type MyClass struct {
  orm.UUIDModel

  // your code here
}
```

Once done your class will be considered a **BaDaaS Model**.

The difference between them is the type they will use as primary key: a random uuid and an auto incremental uint respectively. Both provide date created, edited and deleted (<https://gorm.io/docs/delete.html#Soft-Delete>).

### CRUDServiceModule

`CRUDServiceModule` provides you a CRUDService and a CRUDRepository for your badaas Model. After calling it as, for example, `orm.GetCRUDServiceModule[models.Company](),` the following can be used by dependency injection:

- `crudCompanyService orm.CRUDService[models.Company, orm.UUID]`
- `crudCompanyRepository orm.CRUDRepository[models.Company, orm.UUID]`

These classes will allow you to perform queries using the compilable query system generated with badaas-cli. For details on how to do this visit [badaas-cli docs](github.com/ditrit/badaas-cli/README.md).
