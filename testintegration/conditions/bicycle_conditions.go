// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func BicycleId(exprs ...orm.Expression[orm.UUID]) orm.FieldCondition[models.Bicycle, orm.UUID] {
	return orm.FieldCondition[models.Bicycle, orm.UUID]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func BicycleCreatedAt(exprs ...orm.Expression[time.Time]) orm.FieldCondition[models.Bicycle, time.Time] {
	return orm.FieldCondition[models.Bicycle, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func BicycleUpdatedAt(exprs ...orm.Expression[time.Time]) orm.FieldCondition[models.Bicycle, time.Time] {
	return orm.FieldCondition[models.Bicycle, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func BicycleDeletedAt(exprs ...orm.Expression[gorm.DeletedAt]) orm.FieldCondition[models.Bicycle, gorm.DeletedAt] {
	return orm.FieldCondition[models.Bicycle, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func BicycleName(exprs ...orm.Expression[string]) orm.FieldCondition[models.Bicycle, string] {
	return orm.FieldCondition[models.Bicycle, string]{
		Expressions: exprs,
		Field:       "Name",
	}
}
func BicycleOwner(conditions ...orm.Condition[models.Person]) orm.Condition[models.Bicycle] {
	return orm.JoinCondition[models.Bicycle, models.Person]{
		Conditions: conditions,
		T1Field:    "OwnerName",
		T2Field:    "Name",
	}
}
func BicycleOwnerName(exprs ...orm.Expression[string]) orm.FieldCondition[models.Bicycle, string] {
	return orm.FieldCondition[models.Bicycle, string]{
		Expressions: exprs,
		Field:       "OwnerName",
	}
}
