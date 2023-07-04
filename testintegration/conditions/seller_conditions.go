// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"reflect"
	"time"
)

var sellerType = reflect.TypeOf(*new(models.Seller))
var SellerIdField = badorm.FieldIdentifier[badorm.UUID]{
	Field:     "ID",
	ModelType: sellerType,
}

func SellerId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: SellerIdField,
	}
}

var SellerCreatedAtField = badorm.FieldIdentifier[time.Time]{
	Field:     "CreatedAt",
	ModelType: sellerType,
}

func SellerCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, time.Time]{
		Expression:      expr,
		FieldIdentifier: SellerCreatedAtField,
	}
}

var SellerUpdatedAtField = badorm.FieldIdentifier[time.Time]{
	Field:     "UpdatedAt",
	ModelType: sellerType,
}

func SellerUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, time.Time]{
		Expression:      expr,
		FieldIdentifier: SellerUpdatedAtField,
	}
}

var SellerDeletedAtField = badorm.FieldIdentifier[gorm.DeletedAt]{
	Field:     "DeletedAt",
	ModelType: sellerType,
}

func SellerDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: SellerDeletedAtField,
	}
}

var SellerNameField = badorm.FieldIdentifier[string]{
	Field:     "Name",
	ModelType: sellerType,
}

func SellerName(expr badorm.Expression[string]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, string]{
		Expression:      expr,
		FieldIdentifier: SellerNameField,
	}
}
func SellerCompany(conditions ...badorm.Condition[models.Company]) badorm.IJoinCondition[models.Seller] {
	return badorm.JoinCondition[models.Seller, models.Company]{
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

func SellerCompanyId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: SellerCompanyIdField,
	}
}
func SellerUniversity(conditions ...badorm.Condition[models.University]) badorm.IJoinCondition[models.Seller] {
	return badorm.JoinCondition[models.Seller, models.University]{
		Conditions:         conditions,
		RelationField:      "University",
		T1Field:            "UniversityID",
		T1PreloadCondition: SellerPreloadAttributes,
		T2Field:            "ID",
	}
}

var SellerPreloadUniversity = SellerUniversity(UniversityPreloadAttributes)
var SellerUniversityIdField = badorm.FieldIdentifier[badorm.UUID]{
	Field:     "UniversityID",
	ModelType: sellerType,
}

func SellerUniversityId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: SellerUniversityIdField,
	}
}

var SellerPreloadAttributes = badorm.NewPreloadCondition[models.Seller](SellerIdField, SellerCreatedAtField, SellerUpdatedAtField, SellerDeletedAtField, SellerNameField, SellerCompanyIdField, SellerUniversityIdField)
var SellerPreloadRelations = []badorm.Condition[models.Seller]{SellerPreloadCompany, SellerPreloadUniversity}
