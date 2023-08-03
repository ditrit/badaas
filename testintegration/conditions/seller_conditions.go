// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	"time"
)

func SellerId(operator orm.Operator[orm.UUID]) orm.WhereCondition[models.Seller] {
	return orm.FieldCondition[models.Seller, orm.UUID]{
		FieldIdentifier: orm.IDFieldID,
		Operator:        operator,
	}
}
func SellerCreatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Seller] {
	return orm.FieldCondition[models.Seller, time.Time]{
		FieldIdentifier: orm.CreatedAtFieldID,
		Operator:        operator,
	}
}
func SellerUpdatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Seller] {
	return orm.FieldCondition[models.Seller, time.Time]{
		FieldIdentifier: orm.UpdatedAtFieldID,
		Operator:        operator,
	}
}
func SellerDeletedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Seller] {
	return orm.FieldCondition[models.Seller, time.Time]{
		FieldIdentifier: orm.DeletedAtFieldID,
		Operator:        operator,
	}
}

var sellerNameFieldID = orm.FieldIdentifier{Field: "Name"}

func SellerName(operator orm.Operator[string]) orm.WhereCondition[models.Seller] {
	return orm.FieldCondition[models.Seller, string]{
		FieldIdentifier: sellerNameFieldID,
		Operator:        operator,
	}
}
func SellerCompany(conditions ...orm.Condition[models.Company]) orm.IJoinCondition[models.Seller] {
	return orm.JoinCondition[models.Seller, models.Company]{
		Conditions:         conditions,
		RelationField:      "Company",
		T1Field:            "CompanyID",
		T1PreloadCondition: SellerPreloadAttributes,
		T2Field:            "ID",
	}
}

var SellerPreloadCompany = SellerCompany(CompanyPreloadAttributes)
var sellerCompanyIdFieldID = orm.FieldIdentifier{Field: "CompanyID"}

func SellerCompanyId(operator orm.Operator[orm.UUID]) orm.WhereCondition[models.Seller] {
	return orm.FieldCondition[models.Seller, orm.UUID]{
		FieldIdentifier: sellerCompanyIdFieldID,
		Operator:        operator,
	}
}

var SellerPreloadAttributes = orm.NewPreloadCondition[models.Seller](sellerNameFieldID, sellerCompanyIdFieldID)
var SellerPreloadRelations = []orm.Condition[models.Seller]{SellerPreloadCompany}
