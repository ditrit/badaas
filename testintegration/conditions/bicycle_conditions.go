// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func BicycleId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Bicycle] {
	return badorm.FieldCondition[models.Bicycle, badorm.UUID]{
		Expression: expr,
		Field:      "ID",
	}
}
func BicycleCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Bicycle] {
	return badorm.FieldCondition[models.Bicycle, time.Time]{
		Expression: expr,
		Field:      "CreatedAt",
	}
}
func BicycleUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Bicycle] {
	return badorm.FieldCondition[models.Bicycle, time.Time]{
		Expression: expr,
		Field:      "UpdatedAt",
	}
}
func BicycleDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[models.Bicycle] {
	return badorm.FieldCondition[models.Bicycle, gorm.DeletedAt]{
		Expression: expr,
		Field:      "DeletedAt",
	}
}
func BicycleName(expr badorm.Expression[string]) badorm.WhereCondition[models.Bicycle] {
	return badorm.FieldCondition[models.Bicycle, string]{
		Expression: expr,
		Field:      "Name",
	}
}
func BicycleOwner(conditions ...badorm.Condition[models.Person]) badorm.Condition[models.Bicycle] {
	return badorm.JoinCondition[models.Bicycle, models.Person]{
		Conditions: conditions,
		T1Field:    "OwnerName",
		T2Field:    "Name",
	}
}
func BicycleOwnerName(expr badorm.Expression[string]) badorm.WhereCondition[models.Bicycle] {
	return badorm.FieldCondition[models.Bicycle, string]{
		Expression: expr,
		Field:      "OwnerName",
	}
}
