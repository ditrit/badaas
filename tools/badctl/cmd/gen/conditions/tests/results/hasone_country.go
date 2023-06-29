// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	hasone "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/hasone"
	gorm "gorm.io/gorm"
	"time"
)

func CountryId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[hasone.Country] {
	return badorm.FieldCondition[hasone.Country, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func CountryCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[hasone.Country] {
	return badorm.FieldCondition[hasone.Country, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func CountryUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[hasone.Country] {
	return badorm.FieldCondition[hasone.Country, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func CountryDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[hasone.Country] {
	return badorm.FieldCondition[hasone.Country, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}
func CountryCapital(conditions ...badorm.Condition[hasone.City]) badorm.IJoinCondition[hasone.Country] {
	return badorm.JoinCondition[hasone.Country, hasone.City]{
		Conditions:         conditions,
		RelationField:      "Capital",
		T1Field:            "ID",
		T1PreloadCondition: CountryPreloadAttributes,
		T2Field:            "CountryID",
	}
}

var CountryPreloadCapital = CountryCapital(CityPreloadAttributes)
var CountryPreloadAttributes = badorm.NewPreloadCondition[hasone.Country]()
var CountryPreloadRelations = []badorm.Condition[hasone.Country]{CountryPreloadCapital}
