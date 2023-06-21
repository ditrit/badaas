// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func CityId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.City] {
	return badorm.FieldCondition[models.City, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func CityCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.City] {
	return badorm.FieldCondition[models.City, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func CityUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.City] {
	return badorm.FieldCondition[models.City, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func CityDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[models.City] {
	return badorm.FieldCondition[models.City, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}

var cityNameFieldID = badorm.FieldIdentifier{Field: "Name"}

func CityName(expr badorm.Expression[string]) badorm.WhereCondition[models.City] {
	return badorm.FieldCondition[models.City, string]{
		Expression:      expr,
		FieldIdentifier: cityNameFieldID,
	}
}
func CityCountry(conditions ...badorm.Condition[models.Country]) badorm.Condition[models.City] {
	return badorm.JoinCondition[models.City, models.Country]{
		Conditions:    conditions,
		RelationField: "Country",
		T1Field:       "CountryID",
		T2Field:       "ID",
	}
}

var CityPreloadCountry = CityCountry(CountryPreloadAttributes)
var cityCountryIdFieldID = badorm.FieldIdentifier{Field: "CountryID"}

func CityCountryId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.City] {
	return badorm.FieldCondition[models.City, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: cityCountryIdFieldID,
	}
}

var CityPreloadAttributes = badorm.NewPreloadCondition[models.City](cityNameFieldID, cityCountryIdFieldID)
var CityPreloadRelations = []badorm.Condition[models.City]{CityPreloadCountry}
