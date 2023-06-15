// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	overrideforeignkey "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/overrideforeignkey"
	gorm "gorm.io/gorm"
	"time"
)

func BicycleId(exprs ...badorm.Expression[badorm.UUID]) badorm.FieldCondition[overrideforeignkey.Bicycle, badorm.UUID] {
	return badorm.FieldCondition[overrideforeignkey.Bicycle, badorm.UUID]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func BicycleCreatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[overrideforeignkey.Bicycle, time.Time] {
	return badorm.FieldCondition[overrideforeignkey.Bicycle, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func BicycleUpdatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[overrideforeignkey.Bicycle, time.Time] {
	return badorm.FieldCondition[overrideforeignkey.Bicycle, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func BicycleDeletedAt(exprs ...badorm.Expression[gorm.DeletedAt]) badorm.FieldCondition[overrideforeignkey.Bicycle, gorm.DeletedAt] {
	return badorm.FieldCondition[overrideforeignkey.Bicycle, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func BicycleOwner(conditions ...badorm.Condition[overrideforeignkey.Person]) badorm.Condition[overrideforeignkey.Bicycle] {
	return badorm.JoinCondition[overrideforeignkey.Bicycle, overrideforeignkey.Person]{
		Conditions: conditions,
		T1Field:    "OwnerSomethingID",
		T2Field:    "ID",
	}
}
func BicycleOwnerSomethingId(exprs ...badorm.Expression[string]) badorm.FieldCondition[overrideforeignkey.Bicycle, string] {
	return badorm.FieldCondition[overrideforeignkey.Bicycle, string]{
		Expressions: exprs,
		Field:       "OwnerSomethingID",
	}
}