==============================
Advanced query
==============================

Dynamic operators
--------------------------------

In :doc:`/badorm/query` we have seen how to use the operators 
to make comparisons between the attributes of a model and static values such as a string, 
a number, etc. But if we want to make comparisons between two or more attributes of 
the same type we need to use the dynamic operators. 
These, instead of a dynamic value, receive a FieldIdentifier, that is, 
an object that identifies the attribute with which the operation is to be performed.

These identifiers are also generated during the generation of conditions and 
their name of these FieldIdentifiers will be <Model><Attribute>Field where 
<Model> is the model type and <Attribute> is the attribute name.

For example we query all YourModels that has the same value in its String attribute that 
its related Related's String attribute.

.. code-block:: go

    type Related struct {
        badorm.UUIDModel

        String string
    }

    type YourModel struct {
        badorm.UUIDModel

        String string

        Related   Related
        RelatedID badorm.UUID
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelRelated(
            conditions.RelatedString(
                dynamic.Eq(conditions.YourModelStringField),
            ),
        ),
    )

**Attention**, when using dynamic operators the verification that the FieldIdentifier 
is concerned by the query is performed at run time, returning an error otherwise. 
For example:

.. code-block:: go

    type Related struct {
        badorm.UUIDModel

        String string
    }

    type YourModel struct {
        badorm.UUIDModel

        String string

        Related   Related
        RelatedID badorm.UUID
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelString(
            dynamic.Eq(conditions.RelatedStringField),
        ),
    )

will respond badorm.ErrFieldModelNotConcerned in err.

All operators supported by BaDORM that receive any value are available in their dynamic version at
<https://pkg.go.dev/github.com/ditrit/badaas/badorm/dynamic>. 
In turn, there are dynamic versions of the operators specific to each database that can be found at
<https://pkg.go.dev/github.com/ditrit/badaas/badorm/dynamic/mysqldynamic> and 
<https://pkg.go.dev/github.com/ditrit/badaas/badorm/dynamic/sqlserverdynamic>.

Select join
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

In case the attribute to be used by the dynamic operator is present more 
than once in the query, it will be necessary to select the join to be used, 
to avoid getting the error badorm.ErrJoinMustBeSelected. 
To do this, you must use the SelectJoin method, as in the following example:

.. code-block:: go

    type ParentParent struct {
        badorm.UUIDModel
    }

    type Parent1 struct {
        badorm.UUIDModel

        ParentParent   ParentParent
        ParentParentID badorm.UUID
    }

    type Parent2 struct {
        badorm.UUIDModel

        ParentParent   ParentParent
        ParentParentID badorm.UUID
    }

    type Child struct {
        badorm.UUIDModel

        Parent1   Parent1
        Parent1ID badorm.UUID

        Parent2   Parent2
        Parent2ID badorm.UUID
    }

    models, err := ts.crudChildService.Query(
        conditions.ChildParent1(
            conditions.Parent1ParentParent(),
        ),
        conditions.ChildParent2(
            conditions.Parent2ParentParent(),
        ),
        conditions.ChildName(
            // for the value 0 (conditions.ParentParentNameField),
            // choose the first (0) join (made by conditions.ChildParent1())
            dynamic.Eq(conditions.ParentParentNameField).SelectJoin(0, 0),
        ),
    )

Multitype operators
----------------------------

To make as many checks as possible at compile time to avoid run-time errors,
the dynamic operators only accept FieldIdentifiers that are of the same type as the condition attribute. 
But there are cases in which we want to make this limitation more flexible, 
either because we want to compare attributes of :ref:`nullable <badorm/concepts:nullable types>` 
and non-nullable type or because the operator accepts multiple values and we want 
to combine both static and dynamic values, as in the following example:

.. code-block:: go

    type YourModel struct {
        badorm.UUIDModel

        Int int
        NullableInt sql.NullInt32
    }

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelInt(
            multitype.Between(1, conditions.YourModelNullableIntField),
        ),
    )

In case the type of any of the operator parameters is not related to the type of the condition's attribute, 
err will be multitype.ErrFieldTypeDoesNotMatch.

All operators supported by BaDORM that receive any value are available in their multitype version at
<https://pkg.go.dev/github.com/ditrit/badaas/badorm/multitype>. 
In turn, there are multitype versions of the operators specific to each database that can be found at
<https://pkg.go.dev/github.com/ditrit/badaas/badorm/multitype/mysqlmultitype> and 
<https://pkg.go.dev/github.com/ditrit/badaas/badorm/multitype/sqlservermultitype>.

Unsafe operators
--------------------------------

With multitype operators we add more flexibility to the operators at the cost of more validations 
to be performed at runtime. 
However, BaDORM will validate that the types of the values to be used inside the operator 
are the same or related. 

In case you want to avoid this validation, unsafe operators should be used. 
Although their use is not recommended, this can be useful when the database 
used allows operations between different types or when attributes of different 
types map at the same time in the database (see <https://gorm.io/docs/data_types.html>).

If it is neither of these two cases, the use of an unsafe operator will result in 
an error in the execution of the query that depends on the database used.

All operators supported by BaDORM that receive any value are available in their unsafe version at
<https://pkg.go.dev/github.com/ditrit/badaas/badorm/unsafe>. 
In turn, there are unsafe versions of the operators specific to each database that can be found at
<https://pkg.go.dev/github.com/ditrit/badaas/badorm/unsafe/mysqlunsafe> and 
<https://pkg.go.dev/github.com/ditrit/badaas/badorm/unsafe/sqlserverunsafe>.

Unsafe conditions (raw SQL)
--------------------------------

In case you need to use operators that are not supported by BaDORM 
(please create an issue in our repository if you think we have forgotten any), 
you can always run raw SQL with unsafe.NewCondition, as in the following example:

.. code-block:: go

    yourModels, err := ts.crudYourModelService.Query(
        conditions.YourModelString(
            unsafe.NewCondition[models.YourModel]("%s.name = NULL"),
        ),
    )

As you can see in the example, "%s" can be used in the raw SQL to be replaced 
by the table name of the model to which the condition belongs.

Of course, its use is not recommended because it can generate errors in the execution 
of the query that will depend on the database used.