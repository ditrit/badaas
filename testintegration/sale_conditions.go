// Code generated by badctl v0.0.0, DO NOT EDIT.
package testintegration

import (
	badorm "github.com/ditrit/badaas/badorm"
	uuid "github.com/google/uuid"
	"time"
)

func SaleIdCondition(v uuid.UUID) badorm.WhereCondition[Sale] {
	return badorm.WhereCondition[Sale]{
		Field: "id",
		Value: v,
	}
}
func SaleCreatedAtCondition(v time.Time) badorm.WhereCondition[Sale] {
	return badorm.WhereCondition[Sale]{
		Field: "created_at",
		Value: v,
	}
}
func SaleUpdatedAtCondition(v time.Time) badorm.WhereCondition[Sale] {
	return badorm.WhereCondition[Sale]{
		Field: "updated_at",
		Value: v,
	}
}
func SaleCodeCondition(v int) badorm.WhereCondition[Sale] {
	return badorm.WhereCondition[Sale]{
		Field: "code",
		Value: v,
	}
}
func SaleDescriptionCondition(v string) badorm.WhereCondition[Sale] {
	return badorm.WhereCondition[Sale]{
		Field: "description",
		Value: v,
	}
}
func SaleProductCondition(conditions ...badorm.Condition[Product]) badorm.Condition[Sale] {
	return badorm.JoinCondition[Sale, Product]{
		Conditions: conditions,
		Field:      "product",
	}
}
func SaleProductIdCondition(v uuid.UUID) badorm.WhereCondition[Sale] {
	return badorm.WhereCondition[Sale]{
		Field: "product_id",
		Value: v,
	}
}
func SaleSellerCondition(conditions ...badorm.Condition[Seller]) badorm.Condition[Sale] {
	return badorm.JoinCondition[Sale, Seller]{
		Conditions: conditions,
		Field:      "seller",
	}
}
func SaleSellerIdCondition(v *uuid.UUID) badorm.WhereCondition[Sale] {
	return badorm.WhereCondition[Sale]{
		Field: "seller_id",
		Value: v,
	}
}