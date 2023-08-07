==============================
Quickstart
==============================

Example
---------------------------

To quickly understand the features provided by Badaas-orm, you can head to the 
`example <https://github.com/ditrit/badaas-orm-example>`_, where you will find two different variations:

- `standalone/` where badaas-orm is used in the simplest possible way.
- `fx/` where badaas-orm is used within the :ref:`fx dependency injection system <badaas-orm/concepts:dependency injection>`

Refer to its README.md for running it.

Understand it
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

In this section we will see the steps carried out to develop this example.

**Standalone**

Once you have started your project with `go init`, you must add the dependency to BaDaaS and others:

.. code-block:: bash

    go get -u github.com/ditrit/badaas github.com/uber-go/zap gorm.io/gorm


In models.go the :ref:`models <badaas-orm/concepts:model>` are defined and 
in conditions/orm.go the file required to 
:ref:`generate the conditions <badaas-orm/concepts:conditions generation>` is created.

In main.go a main function is created with the configuration required to use the badaas-orm. 
First, we need to create a :ref:`gormDB <badaas-orm/concepts:gormDB>` that allows connection with the database:

.. code-block:: go

    gormDB, err := NewGormDBConnection()

After that, we have to call the :ref:`AutoMigrate <badaas-orm/concepts:auto migration>` 
method of the gormDB with the models you want to be persisted::

    err = gormDB.AutoMigrate(
      models.Product{},
      models.Company{},
      models.Seller{},
      models.Sale{},
    )

From here, we can start to use badaas-orm, getting the :ref:`CRUDService <badaas-orm/concepts:CRUDService>` 
and :ref:`CRUDRepository <badaas-orm/concepts:CRUDRepository>` of a model with the GetCRUD function:

.. code-block:: go

    crudProductService, crudProductRepository := orm.GetCRUD[models.Product, model.UUID](gormDB)

As you can see, we need to specify the type of the model and the kind 
of :ref:`id <badaas-orm/concepts:model ID>` this model uses.

Finally, you can use this service and repository to perform CRUD operations on your model:

.. code-block:: go

  CreateCRUDObjects(gormDB, crudProductRepository)
  QueryCRUDObjects(crudProductService)

This two functions are defined in `example.go`. 
In `QueryCRUDObjects` you can find a basic usage of the 
:ref:`compilable query system <badaas-orm/concepts:compilable query system>`.

**Fx**

Once you have started your project with `go init`, you must add the dependency to BaDaaS and others:

.. code-block:: bash

  go get -u github.com/ditrit/badaas github.com/uber-go/fx github.com/uber-go/zap gorm.io/gorm

In models.go the :ref:`models <badaas-orm/concepts:model>` are defined and 
in conditions/orm.go the file required to 
:ref:`generate the conditions <badaas-orm/concepts:conditions generation>` is created.

In main.go a main function is created with the configuration required to use the badaas-orm with fx. 
First, we will need to start your application with `fx`:

.. code-block:: go

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

There are some things you need to provide to the badaas-orm module:

- `NewGORMDBConnection` is the function that we need to create 
  a :ref:`gormDB <badaas-orm/concepts:gormDB>` that allows connection with the database.
- `GetModels` is a function that returns in a `orm.GetModelsResult` the list of models 
  you want to be persisted by the :ref:`auto migration <badaas-orm/concepts:auto migration>`.

After that, you can execute the auto-migration with `orm.AutoMigrate` 
and create :ref:`CRUDServices <badaas-orm/concepts:CRUDService>` 
to your models using `orm.GetCRUDServiceModule`.

Finally, we call the functions `CreateCRUDObjects` 
and `QueryCRUDObjects` where the CRUDServices are injected to create, 
read, update and delete the models easily. This two functions are defined in `example.go`. 
In `QueryCRUDObjects` you can find a basic usage of the :ref:`compilable query system <badaas-orm/concepts:compilable query system>`.