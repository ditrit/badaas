==============================
Functionalities
==============================

Base models
-----------------------

BaDORM gives you two types of base models for your classes: `badorm.UUIDModel` and `badorm.UIntModel`.

To use them, simply embed the desired model in any of your classes::

    type MyClass struct {
      badorm.UUIDModel

      // your code here
    }

Once done your class will be considered a **BaDORM Model**.

The difference between them is the type they will use as primary key: 
a random uuid and an auto incremental uint respectively. 
Both provide date created, edited and deleted (<https://gorm.io/docs/delete.html#Soft-Delete>).

CRUDServiceModule
-----------------------

`CRUDServiceModule` provides you a CRUDService and a CRUDRepository for your BaDORM Model. 
After calling it as, for example, `badorm.GetCRUDServiceModule[models.Company](),` 
the following can be used by dependency injection:

- `crudCompanyService badorm.CRUDService[models.Company, uuid.UUID]`
- `crudCompanyRepository badorm.CRUDRepository[models.Company, uuid.UUID]`

These classes will allow you to perform queries using the compilable query system generated with BaDctl. 
For details on how to do this visit :doc:`/badctl/commands`

CRUDUnsafeServiceModule
-----------------------

`CRUDUnsafeServiceModule` provides you a CRUDUnsafeService and a CRUDUnsafeRepository for your BaDORM Model. 
After calling it as, for example, `badorm.GetCRUDUnsafeServiceModule[models.Company](),` 
the following can be used by dependency injection:

- `crudCompanyService badorm.CRUDUnsafeService[models.Company, uuid.UUID]`
- `crudCompanyRepository badorm.CRUDUnsafeRepository[models.Company, uuid.UUID]`

These classes will allow you to perform queries using maps as conditions. 
**Its direct use is not recommended**, since using the compilable query system we can make 
sure that the query is correct at compile time, while here errors will happen at runtime in 
case your condition map is not well structured. 
This functionality is used internally by BaDaaS to provide an http api for queries.