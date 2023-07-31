// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func CountryId(v orm.UUID) orm.WhereCondition[models.Country] {
	return orm.WhereCondition[models.Country]{
		Field: "ID",
		Value: v,
	}
}
func CountryCreatedAt(v time.Time) orm.WhereCondition[models.Country] {
	return orm.WhereCondition[models.Country]{
		Field: "CreatedAt",
		Value: v,
	}
}
func CountryUpdatedAt(v time.Time) orm.WhereCondition[models.Country] {
	return orm.WhereCondition[models.Country]{
		Field: "UpdatedAt",
		Value: v,
	}
}
func CountryDeletedAt(v gorm.DeletedAt) orm.WhereCondition[models.Country] {
	return orm.WhereCondition[models.Country]{
		Field: "DeletedAt",
		Value: v,
	}
}
func CountryName(v string) orm.WhereCondition[models.Country] {
	return orm.WhereCondition[models.Country]{
		Field: "Name",
		Value: v,
	}
}
func CountryCapital(conditions ...orm.Condition[models.City]) orm.Condition[models.Country] {
	return orm.JoinCondition[models.Country, models.City]{
		Conditions: conditions,
		T1Field:    "ID",
		T2Field:    "CountryID",
	}
}
func CityCountry(conditions ...orm.Condition[models.Country]) orm.Condition[models.City] {
	return orm.JoinCondition[models.City, models.Country]{
		Conditions: conditions,
		T1Field:    "CountryID",
		T2Field:    "ID",
	}
}
