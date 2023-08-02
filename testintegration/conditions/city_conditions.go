// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func CityId(operator orm.Operator[orm.UUID]) orm.WhereCondition[models.City] {
	return orm.FieldCondition[models.City, orm.UUID]{
		Field:    "ID",
		Operator: operator,
	}
}
func CityCreatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.City] {
	return orm.FieldCondition[models.City, time.Time]{
		Field:    "CreatedAt",
		Operator: operator,
	}
}
func CityUpdatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.City] {
	return orm.FieldCondition[models.City, time.Time]{
		Field:    "UpdatedAt",
		Operator: operator,
	}
}
func CityDeletedAt(operator orm.Operator[gorm.DeletedAt]) orm.WhereCondition[models.City] {
	return orm.FieldCondition[models.City, gorm.DeletedAt]{
		Field:    "DeletedAt",
		Operator: operator,
	}
}
func CityName(operator orm.Operator[string]) orm.WhereCondition[models.City] {
	return orm.FieldCondition[models.City, string]{
		Field:    "Name",
		Operator: operator,
	}
}
func CityCountryId(operator orm.Operator[orm.UUID]) orm.WhereCondition[models.City] {
	return orm.FieldCondition[models.City, orm.UUID]{
		Field:    "CountryID",
		Operator: operator,
	}
}
