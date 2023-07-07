// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	hasmany "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/hasmany"
	gorm "gorm.io/gorm"
	"reflect"
	"time"
)

var sellerType = reflect.TypeOf(*new(hasmany.Seller))
var SellerIdField = badorm.FieldIdentifier[badorm.UUID]{
	Field:     "ID",
	ModelType: sellerType,
}

func SellerId(operator badorm.Operator[badorm.UUID]) badorm.WhereCondition[hasmany.Seller] {
	return badorm.FieldCondition[hasmany.Seller, badorm.UUID]{
		FieldIdentifier: SellerIdField,
		Operator:        operator,
	}
}

var SellerCreatedAtField = badorm.FieldIdentifier[time.Time]{
	Field:     "CreatedAt",
	ModelType: sellerType,
}

func SellerCreatedAt(operator badorm.Operator[time.Time]) badorm.WhereCondition[hasmany.Seller] {
	return badorm.FieldCondition[hasmany.Seller, time.Time]{
		FieldIdentifier: SellerCreatedAtField,
		Operator:        operator,
	}
}

var SellerUpdatedAtField = badorm.FieldIdentifier[time.Time]{
	Field:     "UpdatedAt",
	ModelType: sellerType,
}

func SellerUpdatedAt(operator badorm.Operator[time.Time]) badorm.WhereCondition[hasmany.Seller] {
	return badorm.FieldCondition[hasmany.Seller, time.Time]{
		FieldIdentifier: SellerUpdatedAtField,
		Operator:        operator,
	}
}

var SellerDeletedAtField = badorm.FieldIdentifier[gorm.DeletedAt]{
	Field:     "DeletedAt",
	ModelType: sellerType,
}

func SellerDeletedAt(operator badorm.Operator[gorm.DeletedAt]) badorm.WhereCondition[hasmany.Seller] {
	return badorm.FieldCondition[hasmany.Seller, gorm.DeletedAt]{
		FieldIdentifier: SellerDeletedAtField,
		Operator:        operator,
	}
}
func SellerCompany(conditions ...badorm.Condition[hasmany.Company]) badorm.IJoinCondition[hasmany.Seller] {
	return badorm.JoinCondition[hasmany.Seller, hasmany.Company]{
		Conditions:         conditions,
		RelationField:      "Company",
		T1Field:            "CompanyID",
		T1PreloadCondition: SellerPreloadAttributes,
		T2Field:            "ID",
	}
}

var SellerPreloadCompany = SellerCompany(CompanyPreloadAttributes)
var SellerCompanyIdField = badorm.FieldIdentifier[badorm.UUID]{
	Field:     "CompanyID",
	ModelType: sellerType,
}

func SellerCompanyId(operator badorm.Operator[badorm.UUID]) badorm.WhereCondition[hasmany.Seller] {
	return badorm.FieldCondition[hasmany.Seller, badorm.UUID]{
		FieldIdentifier: SellerCompanyIdField,
		Operator:        operator,
	}
}

var SellerPreloadAttributes = badorm.NewPreloadCondition[hasmany.Seller](SellerIdField, SellerCreatedAtField, SellerUpdatedAtField, SellerDeletedAtField, SellerCompanyIdField)
var SellerPreloadRelations = []badorm.Condition[hasmany.Seller]{SellerPreloadCompany}
