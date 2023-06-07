// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	overrideforeignkey "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/overrideforeignkey"
	gorm "gorm.io/gorm"
	"time"
)

func BicycleId(v badorm.UUID) badorm.WhereCondition[overrideforeignkey.Bicycle] {
	return badorm.WhereCondition[overrideforeignkey.Bicycle]{
		Field: "id",
		Value: v,
	}
}
func BicycleCreatedAt(v time.Time) badorm.WhereCondition[overrideforeignkey.Bicycle] {
	return badorm.WhereCondition[overrideforeignkey.Bicycle]{
		Field: "created_at",
		Value: v,
	}
}
func BicycleUpdatedAt(v time.Time) badorm.WhereCondition[overrideforeignkey.Bicycle] {
	return badorm.WhereCondition[overrideforeignkey.Bicycle]{
		Field: "updated_at",
		Value: v,
	}
}
func BicycleDeletedAt(v gorm.DeletedAt) badorm.WhereCondition[overrideforeignkey.Bicycle] {
	return badorm.WhereCondition[overrideforeignkey.Bicycle]{
		Field: "deleted_at",
		Value: v,
	}
}
func BicycleOwner(conditions ...badorm.Condition[overrideforeignkey.Person]) badorm.Condition[overrideforeignkey.Bicycle] {
	return badorm.JoinCondition[overrideforeignkey.Bicycle, overrideforeignkey.Person]{
		Conditions: conditions,
		T1Field:    "owner_something_id",
		T2Field:    "id",
	}
}
func BicycleOwnerSomethingId(v string) badorm.WhereCondition[overrideforeignkey.Bicycle] {
	return badorm.WhereCondition[overrideforeignkey.Bicycle]{
		Field: "owner_something_id",
		Value: v,
	}
}
