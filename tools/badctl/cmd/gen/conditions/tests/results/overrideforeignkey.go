// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	overrideforeignkey "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/overrideforeignkey"
	gorm "gorm.io/gorm"
	"time"
)

func BicycleId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[overrideforeignkey.Bicycle] {
	return badorm.FieldCondition[overrideforeignkey.Bicycle, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func BicycleCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[overrideforeignkey.Bicycle] {
	return badorm.FieldCondition[overrideforeignkey.Bicycle, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func BicycleUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[overrideforeignkey.Bicycle] {
	return badorm.FieldCondition[overrideforeignkey.Bicycle, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func BicycleDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[overrideforeignkey.Bicycle] {
	return badorm.FieldCondition[overrideforeignkey.Bicycle, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}
func BicycleOwner(conditions ...badorm.Condition[overrideforeignkey.Person]) badorm.Condition[overrideforeignkey.Bicycle] {
	return badorm.JoinCondition[overrideforeignkey.Bicycle, overrideforeignkey.Person]{
		Conditions:    conditions,
		RelationField: "Owner",
		T1Field:       "OwnerSomethingID",
		T2Field:       "ID",
	}
}
func BicycleOwnerSomethingId(expr badorm.Expression[string]) badorm.WhereCondition[overrideforeignkey.Bicycle] {
	return badorm.FieldCondition[overrideforeignkey.Bicycle, string]{
		Expression:      expr,
		FieldIdentifier: badorm.FieldIdentifier{Field: "OwnerSomethingID"},
	}
}
