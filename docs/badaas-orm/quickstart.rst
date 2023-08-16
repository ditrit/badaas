==============================
Quickstart
==============================

To integrate badaas-orm into your project, you can head to the 
`quickstart <https://github.com/ditrit/badaas-orm-quickstart>`_, where you will find three different variations:

1. Standalone (not using any other dependency) in `standalone/`
2. Using uber fx for :ref:`dependency injection <badaas-orm/concepts:dependency injection>` in `fx/`
3. Inside a badaas application in `badaas/`

Refer to its README.md for running it.

Understand it
---------------------------------

In this section we will see the steps carried out to develop this quickstart.

Standalone
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Once you have started your project with `go init`, you must add the dependency to BaDaaS:

.. code-block:: bash

    go get -u github.com/ditrit/badaas


In models.go the :ref:`models <badaas-orm/concepts:model>` are defined and 
in conditions/orm.go the file required to 
:ref:`generate the conditions <badaas-orm/concepts:conditions generation>` is created.

In main.go a main function is created with the configuration required to use the badaas-orm. 
First, we need to create a :ref:`gorm.DB <badaas-orm/concepts:GormDB>` that allows connection with the database:

.. code-block:: go

    gormDB, err := NewDBConnection()

After that, we have to call the :ref:`AutoMigrate <badaas-orm/concepts:auto migration>` 
method of the gormDB with the models you want to be persisted::

    err = gormDB.AutoMigrate(
      models.MyModel{},
    )

From here, we can start to use badaas-orm, getting the :ref:`CRUDService <badaas-orm/concepts:CRUDService>` 
and :ref:`CRUDRepository <badaas-orm/concepts:CRUDRepository>` of a model with the GetCRUD function:

.. code-block:: go

    crudMyModelService, crudMyModelRepository := orm.GetCRUD[models.MyModel, model.UUID](gormDB)

As you can see, we need to specify the type of the model and the kind 
of :ref:`id <badaas-orm/concepts:model ID>` this model uses.

Finally, you can use this service and repository to perform CRUD operations on your model:

.. code-block:: go

  Run(crudMyModelService, crudMyModelRepository)

This function is defined in `example.go`. 

Fx
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

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
        fx.Provide(NewZapLogger),
        // connect to db
        fx.Provide(NewDBConnection),
        fx.Provide(GetModels),
        orm.AutoMigrate,

        // logger for fx
        fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
          return &fxevent.ZapLogger{Logger: logger}
        }),

        // create crud services for models
        orm.GetCRUDServiceModule[models.MyModel](),

        // run your code
        fx.Invoke(Run),
      ).Run()
    }

There are some things you need to provide to the badaas-orm module:

- `NewZapLogger` (optional) in this case we will use the zap logger instead of the gorm logger, 
  so we have to provide it and then use it as a logger for fx. 
  For more information visit :doc:`logger`.
- `NewDBConnection` is the function that we need to create 
  a :ref:`gorm.DB <badaas-orm/concepts:GormDB>` that allows connection with the database.
- `GetModels` is a function that returns in a `orm.GetModelsResult` the list of models 
  you want to be persisted by the :ref:`auto migration <badaas-orm/concepts:auto migration>`.

After that, you can execute the auto-migration with `orm.AutoMigrate` 
and create :ref:`CRUDServices <badaas-orm/concepts:CRUDService>` 
to your models using `orm.GetCRUDServiceModule`.

Finally, we call the function `Run` where the CRUDServices and CRUDRepositories are injected, 
allowing to perform CRUD operations on your models. 
This function is defined in `example.go`.

Badaas
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Once you have started your project with `go init`, you must add the dependency to BaDaaS and others:

.. code-block:: bash

  go get -u github.com/ditrit/badaas github.com/uber-go/fx github.com/uber-go/zap

In models.go the :ref:`models <badaas-orm/concepts:model>` are defined and 
in conditions/orm.go the file required to 
:ref:`generate the conditions <badaas-orm/concepts:conditions generation>` is created.

In main.go a main function is created with the configuration required to use the badaas-orm 
services and repositories inside a badaas application: 

.. code-block:: go

  func main() {
    badaas.BaDaaS.AddModules(
      orm.AutoMigrate,
      // create crud services for models
      orm.GetCRUDServiceModule[models.MyModel](),
    ).Provide(
      GetModels,
    ).Invoke(
      // run your code
      Run,
    ).Start()
  }

You need to provide to the badaas application `orm.AutoMigrate` and 
`GetModels` for running the :ref:`auto migration <badaas-orm/concepts:auto migration>`.

After that, you can create :ref:`CRUDServices <badaas-orm/concepts:CRUDService>` 
to your models using `orm.GetCRUDServiceModule`.

Finally, we call the function `Run` where the CRUDServices and CRUDRepositories are injected, 
allowing to perform CRUD operations on your models. 
This function is defined in `example.go`.


For more details about badaas visit :doc:`/index`.

Use it
----------------------

Now that you know how to integrate badaas-orm into your project, 
you can learn how to use it by following the :doc:`tutorial`.