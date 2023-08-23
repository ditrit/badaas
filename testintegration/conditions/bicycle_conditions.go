// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func BicycleId(v orm.UUID) orm.WhereCondition[models.Bicycle] {
	return orm.WhereCondition[models.Bicycle]{
		Field: "ID",
		Value: v,
	}
}
func BicycleCreatedAt(v time.Time) orm.WhereCondition[models.Bicycle] {
	return orm.WhereCondition[models.Bicycle]{
		Field: "CreatedAt",
		Value: v,
	}
}
func BicycleUpdatedAt(v time.Time) orm.WhereCondition[models.Bicycle] {
	return orm.WhereCondition[models.Bicycle]{
		Field: "UpdatedAt",
		Value: v,
	}
}
func BicycleDeletedAt(v gorm.DeletedAt) orm.WhereCondition[models.Bicycle] {
	return orm.WhereCondition[models.Bicycle]{
		Field: "DeletedAt",
		Value: v,
	}
}
func BicycleName(v string) orm.WhereCondition[models.Bicycle] {
	return orm.WhereCondition[models.Bicycle]{
		Field: "Name",
		Value: v,
	}
}
func BicycleOwner(conditions ...orm.Condition[models.Person]) orm.Condition[models.Bicycle] {
	return orm.JoinCondition[models.Bicycle, models.Person]{
		Conditions: conditions,
		T1Field:    "OwnerName",
		T2Field:    "Name",
	}
}
func BicycleOwnerName(v string) orm.WhereCondition[models.Bicycle] {
	return orm.WhereCondition[models.Bicycle]{
		Field: "OwnerName",
		Value: v,
	}
}
