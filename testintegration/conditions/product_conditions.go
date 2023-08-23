// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func ProductId(v orm.UUID) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		Field: "ID",
		Value: v,
	}
}
func ProductCreatedAt(v time.Time) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		Field: "CreatedAt",
		Value: v,
	}
}
func ProductUpdatedAt(v time.Time) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		Field: "UpdatedAt",
		Value: v,
	}
}
func ProductDeletedAt(v gorm.DeletedAt) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		Field: "DeletedAt",
		Value: v,
	}
}
func ProductString(v string) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		Column: "string_something_else",
		Value:  v,
	}
}
func ProductInt(v int) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		Field: "Int",
		Value: v,
	}
}
func ProductIntPointer(v int) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		Field: "IntPointer",
		Value: v,
	}
}
func ProductFloat(v float64) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		Field: "Float",
		Value: v,
	}
}
func ProductBool(v bool) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		Field: "Bool",
		Value: v,
	}
}
func ProductByteArray(v []uint8) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		Field: "ByteArray",
		Value: v,
	}
}
func ProductMultiString(v models.MultiString) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		Field: "MultiString",
		Value: v,
	}
}
func ProductEmbeddedInt(v int) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		Field: "EmbeddedInt",
		Value: v,
	}
}
func ProductGormEmbeddedInt(v int) orm.WhereCondition[models.Product] {
	return orm.WhereCondition[models.Product]{
		ColumnPrefix: "gorm_embedded_",
		Field:        "Int",
		Value:        v,
	}
}
