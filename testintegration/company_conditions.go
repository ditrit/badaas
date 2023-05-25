// Code generated by badctl v0.0.0, DO NOT EDIT.
package testintegration

import (
	badorm "github.com/ditrit/badaas/badorm"
	uuid "github.com/google/uuid"
	"time"
)

func CompanyIdCondition(v uuid.UUID) badorm.WhereCondition[Company] {
	return badorm.WhereCondition[Company]{
		Field: "id",
		Value: v,
	}
}
func CompanyCreatedAtCondition(v time.Time) badorm.WhereCondition[Company] {
	return badorm.WhereCondition[Company]{
		Field: "created_at",
		Value: v,
	}
}
func CompanyUpdatedAtCondition(v time.Time) badorm.WhereCondition[Company] {
	return badorm.WhereCondition[Company]{
		Field: "updated_at",
		Value: v,
	}
}
func CompanyNameCondition(v string) badorm.WhereCondition[Company] {
	return badorm.WhereCondition[Company]{
		Field: "name",
		Value: v,
	}
}
func SellerCompanyCondition(conditions ...badorm.Condition[Company]) badorm.Condition[Seller] {
	return badorm.JoinCondition[Seller, Company]{
		Conditions: conditions,
		T1Field:    "company_id",
		T2Field:    "id",
	}
}
