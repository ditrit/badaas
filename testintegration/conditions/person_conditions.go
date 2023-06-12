// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func PersonId(exprs ...orm.Expression[orm.UUID]) orm.FieldCondition[models.Person, orm.UUID] {
	return orm.FieldCondition[models.Person, orm.UUID]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func PersonCreatedAt(exprs ...orm.Expression[time.Time]) orm.FieldCondition[models.Person, time.Time] {
	return orm.FieldCondition[models.Person, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func PersonUpdatedAt(exprs ...orm.Expression[time.Time]) orm.FieldCondition[models.Person, time.Time] {
	return orm.FieldCondition[models.Person, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func PersonDeletedAt(exprs ...orm.Expression[gorm.DeletedAt]) orm.FieldCondition[models.Person, gorm.DeletedAt] {
	return orm.FieldCondition[models.Person, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func PersonName(exprs ...orm.Expression[string]) orm.FieldCondition[models.Person, string] {
	return orm.FieldCondition[models.Person, string]{
		Expressions: exprs,
		Field:       "Name",
	}
}
