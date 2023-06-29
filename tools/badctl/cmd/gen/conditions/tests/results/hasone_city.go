// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	hasone "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/hasone"
	gorm "gorm.io/gorm"
	"time"
)

func CityId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[hasone.City] {
	return badorm.FieldCondition[hasone.City, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func CityCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[hasone.City] {
	return badorm.FieldCondition[hasone.City, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func CityUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[hasone.City] {
	return badorm.FieldCondition[hasone.City, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func CityDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[hasone.City] {
	return badorm.FieldCondition[hasone.City, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}
func CityCountry(conditions ...badorm.Condition[hasone.Country]) badorm.IJoinCondition[hasone.City] {
	return badorm.JoinCondition[hasone.City, hasone.Country]{
		Conditions:         conditions,
		RelationField:      "Country",
		T1Field:            "CountryID",
		T1PreloadCondition: CityPreloadAttributes,
		T2Field:            "ID",
	}
}

var CityPreloadCountry = CityCountry(CountryPreloadAttributes)
var cityCountryIdFieldID = badorm.FieldIdentifier{Field: "CountryID"}

func CityCountryId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[hasone.City] {
	return badorm.FieldCondition[hasone.City, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: cityCountryIdFieldID,
	}
}

var CityPreloadAttributes = badorm.NewPreloadCondition[hasone.City](cityCountryIdFieldID)
var CityPreloadRelations = []badorm.Condition[hasone.City]{CityPreloadCountry}
