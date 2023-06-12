// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func SellerId(exprs ...orm.Expression[orm.UUID]) orm.FieldCondition[models.Seller, orm.UUID] {
	return orm.FieldCondition[models.Seller, orm.UUID]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func SellerCreatedAt(exprs ...orm.Expression[time.Time]) orm.FieldCondition[models.Seller, time.Time] {
	return orm.FieldCondition[models.Seller, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func SellerUpdatedAt(exprs ...orm.Expression[time.Time]) orm.FieldCondition[models.Seller, time.Time] {
	return orm.FieldCondition[models.Seller, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func SellerDeletedAt(exprs ...orm.Expression[gorm.DeletedAt]) orm.FieldCondition[models.Seller, gorm.DeletedAt] {
	return orm.FieldCondition[models.Seller, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func SellerName(exprs ...orm.Expression[string]) orm.FieldCondition[models.Seller, string] {
	return orm.FieldCondition[models.Seller, string]{
		Expressions: exprs,
		Field:       "Name",
	}
}
func SellerCompanyId(exprs ...orm.Expression[*orm.UUID]) orm.FieldCondition[models.Seller, *orm.UUID] {
	return orm.FieldCondition[models.Seller, *orm.UUID]{
		Expressions: exprs,
		Field:       "CompanyID",
	}
}
