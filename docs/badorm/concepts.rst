==============================
Concepts
==============================

Model
------------------------------

A model is any object (struct) of go that you want to persist 
in the database and on which you can perform queries. 
For this, the struct must have an embedded BaDORM base model.

BaDORM base model
-----------------------------

It is a struct that when embedded allows your structures to become BadORM models, 
adding attributes ID, CreatedAt, UpdatedAt and DeletedAt attributes and the possibility to persist, 
create conditions and perform queries on these structures.

Model ID
-----------------------------

The id is a unique identifier needed to persist a model in the database. 
It can be an auto-incremental uint or a badorm.uuid, depending on the base model used.

Auto Migration
----------------------------------------------------------

To persist the models it is necessary to migrate the database, 
so that the structure of the tables corresponds to the definition of the model. 
This migration is performed by gorm through the gormDB. 
For details visit `gorm docs <https://gorm.io/docs/migration.html>`_

GormDB
-----------------------------

GormDB is a gorm.DB object that allows communication with the database.

Condition
-----------------------------

Conditions are the basis of the BaDORM query system, every query is composed of a set of conditions. 
Conditions belong to a particular model and there are 4 different types: 
WhereConditions, ConnectionConditions, JoinConditions and PreloadConditions.

WhereCondition
-----------------------------

Type of condition that allows filters to be made on the model to which they belong 
and an attribute of this model. These filters are performed through operators.

ConnectionCondition
-----------------------------

Type of condition that allows the use of logical operators 
(and, or, or, not) between WhereConditions.

JoinCondition
-----------------------------

Condition type that allows to navigate relationships between models, 
which will result in a join in the executed query 
(don't worry, if you don't know what a join is, 
you don't need to understand the queries that badorm executes).

PreloadCondition
-----------------------------

Type of condition that allows retrieving information from a model as a result of the database (preload). 
This information can be all its attributes and/or another model that is related to it.

Operator
-----------------------------

Concept similar to database operators, 
which allow different operations to be performed on an attribute of a model, 
such as comparisons, predicates, pattern matching, etc.

Operators can be classified as static, dynamic, multitype and unsafe.

Static operator
-----------------------------

Static operators are those that perform operations on an attribute and static values, 
such as a boolean value, an integer, etc.

Dynamic operator
-----------------------------

Dynamic operators are those that perform operations between an attribute and other attributes, 
either from the same model or from a different model, as long as the type of these attributes is the same.

Multitype operator
-----------------------------

Multitype operators are those that can perform operations between an attribute 
and static values or other attributes at the same time and, in addition, 
these values and attributes can be of a type related to the type of the attribute.

Unsafe operator
-----------------------------

Unsafe operators are those that can perform operations between an attribute and 
any type of value or attribute.

Database specific operator
-----------------------------

Since not all SQL databases support the same set of operators, 
there are operators that only work for a specific database.

CRUDService
-----------------------------

A CrudService is a service that allows us to perform CRUD (create, read, update and delete) 
operations on a specific model, executing all the necessary operations within a transaction. 
Internally they use the CRUDRepository of that model.

CRUDRepository
-----------------------------

A CRUDRepository is an object that allows us to perform CRUD operations (create, read, update, delete) 
on a model but, unlike services, its internal operations are performed within a transaction received 
by parameter. 
This is useful to be able to define services that perform multiple CRUD 
operations within the same transaction.

Compilable query system
-----------------------------

The set of conditions that are received by the read operations of the CRUDService 
and CRUDRepository form the BaDORM compilable query system. 
It is so named because the conditions will verify at compile time that the query to be executed is correct.

Conditions generation
----------------------------

Conditions are the basis of the compilable query system. 
They are generated for each model and attribute and can then be used. 
Their generation is done with BaDctl.

Dependency injection
-----------------------------------

Dependency injection is a programming technique in which an object or function 
receives other objects or functions that it depends on. BaDORM is compatible with 
`uber fx <https://uber-go.github.io/fx/>`_ to inject the CRUDServices and 
CRUDRepositories in your objects and functions.

Relation getter
-----------------------------------

Relationships between objects can be loaded from the database using PreloadConditions. 
In order to safely navigate the relations in the loaded model BaDORM provides methods 
called "relation getters".