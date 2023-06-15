// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	hasone "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/hasone"
	gorm "gorm.io/gorm"
	"time"
)

func CityId(exprs ...badorm.Expression[badorm.UUID]) badorm.FieldCondition[hasone.City, badorm.UUID] {
	return badorm.FieldCondition[hasone.City, badorm.UUID]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func CityCreatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[hasone.City, time.Time] {
	return badorm.FieldCondition[hasone.City, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func CityUpdatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[hasone.City, time.Time] {
	return badorm.FieldCondition[hasone.City, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func CityDeletedAt(exprs ...badorm.Expression[gorm.DeletedAt]) badorm.FieldCondition[hasone.City, gorm.DeletedAt] {
	return badorm.FieldCondition[hasone.City, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func CityCountryId(exprs ...badorm.Expression[badorm.UUID]) badorm.FieldCondition[hasone.City, badorm.UUID] {
	return badorm.FieldCondition[hasone.City, badorm.UUID]{
		Expressions: exprs,
		Field:       "CountryID",
	}
}