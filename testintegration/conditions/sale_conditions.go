// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"reflect"
	"time"
)

var saleType = reflect.TypeOf(*new(models.Sale))
var SaleIdField = badorm.FieldIdentifier{
	Field:     "ID",
	ModelType: saleType,
	Type:      reflect.TypeOf(*new(badorm.UUID)),
}

func SaleId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Sale] {
	return badorm.FieldCondition[models.Sale, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: SaleIdField,
	}
}

var SaleCreatedAtField = badorm.FieldIdentifier{
	Field:     "CreatedAt",
	ModelType: saleType,
	Type:      reflect.TypeOf(*new(time.Time)),
}

func SaleCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Sale] {
	return badorm.FieldCondition[models.Sale, time.Time]{
		Expression:      expr,
		FieldIdentifier: SaleCreatedAtField,
	}
}

var SaleUpdatedAtField = badorm.FieldIdentifier{
	Field:     "UpdatedAt",
	ModelType: saleType,
	Type:      reflect.TypeOf(*new(time.Time)),
}

func SaleUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Sale] {
	return badorm.FieldCondition[models.Sale, time.Time]{
		Expression:      expr,
		FieldIdentifier: SaleUpdatedAtField,
	}
}

var SaleDeletedAtField = badorm.FieldIdentifier{
	Field:     "DeletedAt",
	ModelType: saleType,
	Type:      reflect.TypeOf(*new(gorm.DeletedAt)),
}

func SaleDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[models.Sale] {
	return badorm.FieldCondition[models.Sale, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: SaleDeletedAtField,
	}
}

var SaleCodeField = badorm.FieldIdentifier{
	Field:     "Code",
	ModelType: saleType,
	Type:      reflect.TypeOf(*new(int)),
}

func SaleCode(expr badorm.Expression[int]) badorm.WhereCondition[models.Sale] {
	return badorm.FieldCondition[models.Sale, int]{
		Expression:      expr,
		FieldIdentifier: SaleCodeField,
	}
}

var SaleDescriptionField = badorm.FieldIdentifier{
	Field:     "Description",
	ModelType: saleType,
	Type:      reflect.TypeOf(*new(string)),
}

func SaleDescription(expr badorm.Expression[string]) badorm.WhereCondition[models.Sale] {
	return badorm.FieldCondition[models.Sale, string]{
		Expression:      expr,
		FieldIdentifier: SaleDescriptionField,
	}
}
func SaleProduct(conditions ...badorm.Condition[models.Product]) badorm.IJoinCondition[models.Sale] {
	return badorm.JoinCondition[models.Sale, models.Product]{
		Conditions:         conditions,
		RelationField:      "Product",
		T1Field:            "ProductID",
		T1PreloadCondition: SalePreloadAttributes,
		T2Field:            "ID",
	}
}

var SalePreloadProduct = SaleProduct(ProductPreloadAttributes)
var SaleProductIdField = badorm.FieldIdentifier{
	Field:     "ProductID",
	ModelType: saleType,
	Type:      reflect.TypeOf(*new(badorm.UUID)),
}

func SaleProductId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Sale] {
	return badorm.FieldCondition[models.Sale, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: SaleProductIdField,
	}
}
func SaleSeller(conditions ...badorm.Condition[models.Seller]) badorm.IJoinCondition[models.Sale] {
	return badorm.JoinCondition[models.Sale, models.Seller]{
		Conditions:         conditions,
		RelationField:      "Seller",
		T1Field:            "SellerID",
		T1PreloadCondition: SalePreloadAttributes,
		T2Field:            "ID",
	}
}

var SalePreloadSeller = SaleSeller(SellerPreloadAttributes)
var SaleSellerIdField = badorm.FieldIdentifier{
	Field:     "SellerID",
	ModelType: saleType,
	Type:      reflect.TypeOf(*new(badorm.UUID)),
}

func SaleSellerId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Sale] {
	return badorm.FieldCondition[models.Sale, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: SaleSellerIdField,
	}
}

var SalePreloadAttributes = badorm.NewPreloadCondition[models.Sale](SaleIdField, SaleCreatedAtField, SaleUpdatedAtField, SaleDeletedAtField, SaleCodeField, SaleDescriptionField, SaleProductIdField, SaleSellerIdField)
var SalePreloadRelations = []badorm.Condition[models.Sale]{SalePreloadProduct, SalePreloadSeller}
