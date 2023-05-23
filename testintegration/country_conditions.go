// Code generated by badctl v0.0.0, DO NOT EDIT.
package testintegration

import (
	badorm "github.com/ditrit/badaas/badorm"
	uuid "github.com/google/uuid"
	"time"
)

func CountryIdCondition(v uuid.UUID) badorm.WhereCondition[Country] {
	return badorm.WhereCondition[Country]{
		Field: "id",
		Value: v,
	}
}
func CountryCreatedAtCondition(v time.Time) badorm.WhereCondition[Country] {
	return badorm.WhereCondition[Country]{
		Field: "created_at",
		Value: v,
	}
}
func CountryUpdatedAtCondition(v time.Time) badorm.WhereCondition[Country] {
	return badorm.WhereCondition[Country]{
		Field: "updated_at",
		Value: v,
	}
}
func CountryNameCondition(v string) badorm.WhereCondition[Country] {
	return badorm.WhereCondition[Country]{
		Field: "name",
		Value: v,
	}
}
func CountryCapitalCondition(conditions ...badorm.Condition[City]) badorm.Condition[Country] {
	return badorm.JoinCondition[Country, City]{
		Conditions: conditions,
		Field:      "capital",
	}
}
func CityCountryCondition(conditions ...badorm.Condition[Country]) badorm.Condition[City] {
	return badorm.JoinCondition[City, Country]{
		Conditions: conditions,
		Field:      "country",
	}
}
