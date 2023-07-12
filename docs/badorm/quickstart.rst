==============================
Quickstart
==============================

BaDaaS Example
---------------------------

If you are interested in using BaDORM within a BaDaaS application you can 
consult the `badaas example <https://github.com/ditrit/badaas-example>`_. 
in which besides using the services and repositories provided by BaDorm, 
BaDaaS adds a controller that allows the query of objects via an http api.

Stand-alone Example
---------------------------

To quickly understand the features provided by BaDORM, you can head to the 
`example <https://github.com/ditrit/badorm-example>`_. 
This example will help you to see how to use BaDorm and as a template to start your own project.

Step-by-step instructions
---------------------------

Once you have started your project with `go init`, you must add the dependency to BaDaaS::

    go get -u github.com/ditrit/badaas github.com/uber-go/fx github.com/uber-go/zap

First of all, you will need to start your application with `fx`::

    func main() {
      fx.New(
        fx.Provide(NewLogger),

        // DB modules
        fx.Provide(NewGORMDBConnection),
        // activate BaDORM
        badorm.BaDORMModule,
        // start example data
        badorm.GetCRUDServiceModule[models.Company](),
        badorm.GetCRUDServiceModule[models.Product](),
        badorm.GetCRUDServiceModule[models.Seller](),
        badorm.GetCRUDServiceModule[models.Sale](),
        fx.Provide(CreateCRUDObjects),
        fx.Invoke(QueryCRUDObjects),
      ).Run()
    }

There are some things you need to provide to the BaDORM module:

- `NewLogger` is the function that provides a zap logger to the BaDorm components.
- `NewGORMDBConnection` if the function that establish the connection to the 
    database where you data will be saved.

After that, you can start the `badorm.BaDORMModule` and crete the CRUD 
services to your models using `badorm.GetCRUDServiceModule`.

Finally, you can call your application functions as `CreateCRUDObjects` 
and `QueryCRUDObjects` where created  CRUDServices can be injected to create, 
read, update and delete your models easily.