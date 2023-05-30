// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	uuid "github.com/google/uuid"
	"time"
)

func CityId(v uuid.UUID) badorm.WhereCondition[models.City] {
	return badorm.WhereCondition[models.City]{
		Field: "id",
		Value: v,
	}
}
func CityCreatedAt(v time.Time) badorm.WhereCondition[models.City] {
	return badorm.WhereCondition[models.City]{
		Field: "created_at",
		Value: v,
	}
}
func CityUpdatedAt(v time.Time) badorm.WhereCondition[models.City] {
	return badorm.WhereCondition[models.City]{
		Field: "updated_at",
		Value: v,
	}
}
func CityName(v string) badorm.WhereCondition[models.City] {
	return badorm.WhereCondition[models.City]{
		Field: "name",
		Value: v,
	}
}
func CityCountryId(v uuid.UUID) badorm.WhereCondition[models.City] {
	return badorm.WhereCondition[models.City]{
		Field: "country_id",
		Value: v,
	}
}