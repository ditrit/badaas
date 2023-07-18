==============================
Query
==============================

Query methods
------------------------

In CRUDRepository you will find different methods that will 
allow you to perform queries on the model to which that repository belongs:

.. TODO ver si se mantienen estos nombres

- GetByID: will allow you to obtain a model by its id.
- Get: will allow you to obtain the model that meets the conditions received by parameter.
- Query: will allow you to obtain the models that meet the conditions received by parameter.

Compilable query system
------------------------

The set of conditions that are received by the read operations of the CRUDService 
and CRUDRepository form the BaDORM compilable query system. 
It is so named because the conditions will verify at compile time that the query to be executed is correct.

These conditions are objects of type badorm.Condition that contain the 
necessary information to perform the queries in a safe way. 
They are generated from the definition of your models using badctl.

Conditions generation
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The generation of conditions is done with badctl. For this, we need to install badctl:

.. code-block:: bash

    go install github.com/ditrit/badaas/tools/badctl

Then, inside our project we will have to create a package called conditions 
(or another name if you wish) and inside it a file with the following content:

.. code-block:: go

    package conditions

    //go:generate badctl gen conditions ../models_path_1 ../models_path_2

where ../models_path_1 ../models_path_2 are the relative paths between the package conditions 
and the packages containing the definition of your models (can be only one).

Now, from the root of your project you can execute:

.. code-block:: bash

  go generate ./...

and the conditions for each of your models will be created in the conditions package.

Use of the conditions
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

After generating the conditions you will have the following conditions:

- One condition for each attribute of each of your models. 
  The name of these conditions will be <Model><Attribute> where 
  <Model> is the model type and <Attribute> is the attribute name. 
  These conditions are of type WhereCondition.
- One condition for each relationship with another model that each of your models has. 
  The name of these conditions will be <Model><Relation> where 
  <Model> is the model type and <Relation> is the name of the attribute that creates the relation. 
  These conditions are of type JoinCondition because using them will 
  mean performing a join within the executed query.

Then, combining these conditions, the Connection Conditions (badorm.And, badorm.Or, badorm.Not) 
and the Operators (badorm.Eq, badorm.Lt, etc.) you will be able to make all 
the queries you need in a safe way.

Examples
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**Filter by an attribute**

In this example we query all YourModel that has "a_string" in the Attribute attribute.

.. code-block:: go

    type YourModel struct {
        badorm.UUIDModel

        Attribute string
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelAttribute(badorm.Eq("a_string")),
    )

**Filter by an attribute of a related model**

In this example we query all YourModels whose related Related has "a_string" in its Attribute attribute.

.. code-block:: go

    type Related struct {
        badorm.UUIDModel

        Attribute string
    }

    type YourModel struct {
        badorm.UUIDModel

        Related   Related
        RelatedID badorm.UUID
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelRelated(
            conditions.RelatedAttribute(badorm.Eq("a_string")),
        ),
    )

**Multiple conditions**

In this example we query all YourModels that has a 4 in the IntAttribute attribute and 
whose related Related has "a_string" in its Attribute attribute.

.. code-block:: go

    type Related struct {
        badorm.UUIDModel

        Attribute string
    }

    type YourModel struct {
        badorm.UUIDModel

        IntAttribute int

        Related   Related
        RelatedID badorm.UUID
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelIntAttribute(badorm.Eq(4)),
        conditions.YourModelRelated(
            conditions.RelatedAttribute(badorm.Eq("a_string")),
        ),
    )

Operators
------------------------

Below you will find the complete list of available operators:

- badorm.Eq(value): EqualTo
- badorm.EqOrIsNull(value): if value is not NULL returns a Eq operator but if value is NULL returns a IsNull operator
- badorm.NotEq(value): NotEqualTo
- badorm.NotEqOrIsNotNull(value): if value is not NULL returns a NotEq operator but if value is NULL returns a IsNotNull operator
- badorm.Lt(value): LessThan
- badorm.LtOrEq(value): LessThanOrEqualTo
- badorm.Gt(value): GreaterThan
- badorm.GtOrEq(value): GreaterThanOrEqualTo
- badorm.IsNull()
- badorm.IsNotNull()
- badorm.Between(v1, v2): Equivalent to v1 < attribute < v2
- badorm.NotBetween(v1, v2): Equivalent to NOT (v1 < attribute < v2)
- badorm.IsTrue() (Not supported by: sqlserver)
- badorm.IsNotTrue() (Not supported by: sqlserver)
- badorm.IsFalse() (Not supported by: sqlserver)
- badorm.IsNotFalse() (Not supported by: sqlserver)
- badorm.IsUnknown() (Not supported by: sqlserver, sqlite)
- badorm.IsNotUnknown() (Not supported by: sqlserver, sqlite)
- badorm.IsDistinct(value) (Not supported by: mysql)
- badorm.IsNotDistinct(value) (Not supported by: mysql)
- badorm.Like(pattern)
- badorm.Like(pattern).Escape(escape)
- badorm.ArrayIn(values)
- badorm.ArrayNotIn(values)

In addition to these, BaDORM gives the possibility to use operators 
that are only supported by a certain database (outside the standard). 
These operators can be found in <https://pkg.go.dev/github.com/ditrit/badaas/badorm/mysql>, 
<https://pkg.go.dev/github.com/ditrit/badaas/badorm/sqlserver>, 
<https://pkg.go.dev/github.com/ditrit/badaas/badorm/psql> 
and <https://pkg.go.dev/github.com/ditrit/badaas/badorm/sqlite>.