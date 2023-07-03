// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	"reflect"
	"time"

	gorm "gorm.io/gorm"

	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
)

func SellerId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func SellerCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func SellerUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func SellerDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}

// TODO generacion automatica
var SellerNameField = badorm.FieldIdentifier{Field: "Name", Type: reflect.String, ModelType: reflect.TypeOf(models.Seller{})}

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
var sellerCompanyIdFieldID = badorm.FieldIdentifier{Field: "CompanyID"}

func SellerCompanyId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: sellerCompanyIdFieldID,
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
var sellerUniversityIdFieldID = badorm.FieldIdentifier{Field: "UniversityID"}

func SellerUniversityId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Seller] {
	return badorm.FieldCondition[models.Seller, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: sellerUniversityIdFieldID,
	}
}

var SellerPreloadAttributes = badorm.NewPreloadCondition[models.Seller](SellerNameField, sellerCompanyIdFieldID, sellerUniversityIdFieldID)
var SellerPreloadRelations = []badorm.Condition[models.Seller]{SellerPreloadCompany, SellerPreloadUniversity}
