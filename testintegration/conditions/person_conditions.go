// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	"reflect"
	"time"
)

var personType = reflect.TypeOf(*new(models.Person))
var PersonIdField = orm.FieldIdentifier[orm.UUID]{
	Field:     "ID",
	ModelType: personType,
}

func PersonId(operator orm.Operator[orm.UUID]) orm.WhereCondition[models.Person] {
	return orm.FieldCondition[models.Person, orm.UUID]{
		FieldIdentifier: PersonIdField,
		Operator:        operator,
	}
}

var PersonCreatedAtField = orm.FieldIdentifier[time.Time]{
	Field:     "CreatedAt",
	ModelType: personType,
}

func PersonCreatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Person] {
	return orm.FieldCondition[models.Person, time.Time]{
		FieldIdentifier: PersonCreatedAtField,
		Operator:        operator,
	}
}

var PersonUpdatedAtField = orm.FieldIdentifier[time.Time]{
	Field:     "UpdatedAt",
	ModelType: personType,
}

func PersonUpdatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Person] {
	return orm.FieldCondition[models.Person, time.Time]{
		FieldIdentifier: PersonUpdatedAtField,
		Operator:        operator,
	}
}

var PersonDeletedAtField = orm.FieldIdentifier[time.Time]{
	Field:     "DeletedAt",
	ModelType: personType,
}

func PersonDeletedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Person] {
	return orm.FieldCondition[models.Person, time.Time]{
		FieldIdentifier: PersonDeletedAtField,
		Operator:        operator,
	}
}

var PersonNameField = orm.FieldIdentifier[string]{
	Field:     "Name",
	ModelType: personType,
}

func PersonName(operator orm.Operator[string]) orm.WhereCondition[models.Person] {
	return orm.FieldCondition[models.Person, string]{
		FieldIdentifier: PersonNameField,
		Operator:        operator,
	}
}

var PersonPreloadAttributes = orm.NewPreloadCondition[models.Person](PersonIdField, PersonCreatedAtField, PersonUpdatedAtField, PersonDeletedAtField, PersonNameField)
