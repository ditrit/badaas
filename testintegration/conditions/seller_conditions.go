// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func SellerId(v badorm.UUID) badorm.WhereCondition[models.Seller] {
	return badorm.WhereCondition[models.Seller]{
		Field: "ID",
		Value: v,
	}
}
func SellerCreatedAt(v time.Time) badorm.WhereCondition[models.Seller] {
	return badorm.WhereCondition[models.Seller]{
		Field: "CreatedAt",
		Value: v,
	}
}
func SellerUpdatedAt(v time.Time) badorm.WhereCondition[models.Seller] {
	return badorm.WhereCondition[models.Seller]{
		Field: "UpdatedAt",
		Value: v,
	}
}
func SellerDeletedAt(v gorm.DeletedAt) badorm.WhereCondition[models.Seller] {
	return badorm.WhereCondition[models.Seller]{
		Field: "DeletedAt",
		Value: v,
	}
}
func SellerName(v string) badorm.WhereCondition[models.Seller] {
	return badorm.WhereCondition[models.Seller]{
		Field: "Name",
		Value: v,
	}
}
func SellerCompanyId(v *badorm.UUID) badorm.WhereCondition[models.Seller] {
	return badorm.WhereCondition[models.Seller]{
		Field: "CompanyID",
		Value: v,
	}
}
