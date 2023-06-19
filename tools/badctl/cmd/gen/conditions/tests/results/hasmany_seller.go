// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	hasmany "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/hasmany"
	gorm "gorm.io/gorm"
	"time"
)

func SellerId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[hasmany.Seller] {
	return badorm.FieldCondition[hasmany.Seller, badorm.UUID]{
		Expression: expr,
		Field:      "ID",
	}
}
func SellerCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[hasmany.Seller] {
	return badorm.FieldCondition[hasmany.Seller, time.Time]{
		Expression: expr,
		Field:      "CreatedAt",
	}
}
func SellerUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[hasmany.Seller] {
	return badorm.FieldCondition[hasmany.Seller, time.Time]{
		Expression: expr,
		Field:      "UpdatedAt",
	}
}
func SellerDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[hasmany.Seller] {
	return badorm.FieldCondition[hasmany.Seller, gorm.DeletedAt]{
		Expression: expr,
		Field:      "DeletedAt",
	}
}
func SellerCompanyId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[hasmany.Seller] {
	return badorm.FieldCondition[hasmany.Seller, badorm.UUID]{
		Expression: expr,
		Field:      "CompanyID",
	}
}
