// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func PhoneId(exprs ...badorm.Expression[uint]) badorm.FieldCondition[models.Phone, uint] {
	return badorm.FieldCondition[models.Phone, uint]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func PhoneCreatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[models.Phone, time.Time] {
	return badorm.FieldCondition[models.Phone, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func PhoneUpdatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[models.Phone, time.Time] {
	return badorm.FieldCondition[models.Phone, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func PhoneDeletedAt(exprs ...badorm.Expression[gorm.DeletedAt]) badorm.FieldCondition[models.Phone, gorm.DeletedAt] {
	return badorm.FieldCondition[models.Phone, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func PhoneName(exprs ...badorm.Expression[string]) badorm.FieldCondition[models.Phone, string] {
	return badorm.FieldCondition[models.Phone, string]{
		Expressions: exprs,
		Field:       "Name",
	}
}
func PhoneBrand(conditions ...badorm.Condition[models.Brand]) badorm.Condition[models.Phone] {
	return badorm.JoinCondition[models.Phone, models.Brand]{
		Conditions: conditions,
		T1Field:    "BrandID",
		T2Field:    "ID",
	}
}
func PhoneBrandId(exprs ...badorm.Expression[uint]) badorm.FieldCondition[models.Phone, uint] {
	return badorm.FieldCondition[models.Phone, uint]{
		Expressions: exprs,
		Field:       "BrandID",
	}
}