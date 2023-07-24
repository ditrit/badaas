==============================
CRUD Operations
==============================

CRUDServices and CRUDRepositories
--------------------------------------

CRUD operations are made to your models via the CRUDServices and CRUDRepositories. 
The difference between them is that a CRUDService will execute this operations within a transaction 
while the CRUDRepository will be executed within a transaction received by parameter, 
thus allowing defining services that perform multiple CRUD operations within the same transaction.

Create, Save and Delete methods are just hooks to the gorm's corresponding methods. 
For details visit 
<https://gorm.io/docs/create.html>, <https://gorm.io/docs/update.html> and <https://gorm.io/docs/delete.html>. 
On the other hand, read (query) operations are provided by BaDORM via its 
:ref:`compilable query system <badorm/concepts:compilable query system>` 
(see how in :doc:`/badorm/query`).

Each pair of CRUDService and CRUDRepository corresponds to a model. To create them you must use 
the `badorm.GetCRUD[<model>, <modelID>](gormDB)` where 
`<model>` is the type of your :ref:`model <badorm/concepts:model>`, 
`<modelID>` is the type of your :ref:`model's id <badorm/concepts:model id>` 
and `gormDB` is the :ref:`GormDB <badorm/concepts:GormDB>` object.

When using BaDORM with `fx` as :ref:`dependency injector <badorm/concepts:Dependency injection>` you 
will need to provide to fx `badorm.GetCRUDServiceModule[<model>]()` 
where `<model>` is the type of your :ref:`model <badorm/concepts:model>`. 
After that the following can be used by dependency injection:

- `crudYourModelService badorm.CRUDService[<model>, <modelID>]`
- `crudYourModelRepository badorm.CRUDRepository[<model>, <modelID>]`

For example:

.. code-block:: go


    type YourModel struct {
        badorm.UUIDModel
    }

    func main() {
        fx.New(
            // activate BaDORM
            fx.Provide(NewGormDBConnection),
            fx.Provide(GetModels),
            badorm.BaDORMModule,

            badorm.GetCRUDServiceModule[YourModel](),
            fx.Invoke(QueryCRUDObjects),
        ).Run()
    }

    func QueryCRUDObjects(crudYourModelService badorm.CRUDService[YourModel, badorm.UUID]) {
        // use crudYourModelService
    }