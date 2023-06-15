// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func BicycleId(exprs ...badorm.Expression[badorm.UUID]) badorm.FieldCondition[models.Bicycle, badorm.UUID] {
	return badorm.FieldCondition[models.Bicycle, badorm.UUID]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func BicycleCreatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[models.Bicycle, time.Time] {
	return badorm.FieldCondition[models.Bicycle, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func BicycleUpdatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[models.Bicycle, time.Time] {
	return badorm.FieldCondition[models.Bicycle, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func BicycleDeletedAt(exprs ...badorm.Expression[gorm.DeletedAt]) badorm.FieldCondition[models.Bicycle, gorm.DeletedAt] {
	return badorm.FieldCondition[models.Bicycle, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func BicycleName(exprs ...badorm.Expression[string]) badorm.FieldCondition[models.Bicycle, string] {
	return badorm.FieldCondition[models.Bicycle, string]{
		Expressions: exprs,
		Field:       "Name",
	}
}
func BicycleOwner(conditions ...badorm.Condition[models.Person]) badorm.Condition[models.Bicycle] {
	return badorm.JoinCondition[models.Bicycle, models.Person]{
		Conditions: conditions,
		T1Field:    "OwnerName",
		T2Field:    "Name",
	}
}
func BicycleOwnerName(exprs ...badorm.Expression[string]) badorm.FieldCondition[models.Bicycle, string] {
	return badorm.FieldCondition[models.Bicycle, string]{
		Expressions: exprs,
		Field:       "OwnerName",
	}
}