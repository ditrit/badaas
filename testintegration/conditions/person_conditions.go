// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	"time"
)

func PersonId(operator orm.Operator[orm.UUID]) orm.WhereCondition[models.Person] {
	return orm.FieldCondition[models.Person, orm.UUID]{
		FieldIdentifier: orm.IDFieldID,
		Operator:        operator,
	}
}
func PersonCreatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Person] {
	return orm.FieldCondition[models.Person, time.Time]{
		FieldIdentifier: orm.CreatedAtFieldID,
		Operator:        operator,
	}
}
func PersonUpdatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Person] {
	return orm.FieldCondition[models.Person, time.Time]{
		FieldIdentifier: orm.UpdatedAtFieldID,
		Operator:        operator,
	}
}
func PersonDeletedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Person] {
	return orm.FieldCondition[models.Person, time.Time]{
		FieldIdentifier: orm.DeletedAtFieldID,
		Operator:        operator,
	}
}

var personNameFieldID = orm.FieldIdentifier{Field: "Name"}

func PersonName(operator orm.Operator[string]) orm.WhereCondition[models.Person] {
	return orm.FieldCondition[models.Person, string]{
		FieldIdentifier: personNameFieldID,
		Operator:        operator,
	}
}

var PersonPreloadAttributes = orm.NewPreloadCondition[models.Person](personNameFieldID)
