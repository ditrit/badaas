// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	uuid "github.com/google/uuid"
	pq "github.com/lib/pq"
	"time"
)

func ProductId(v uuid.UUID) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "id",
		Value: v,
	}
}
func ProductCreatedAt(v time.Time) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "created_at",
		Value: v,
	}
}
func ProductUpdatedAt(v time.Time) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "updated_at",
		Value: v,
	}
}
func ProductString(v string) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "string_something_else",
		Value: v,
	}
}
func ProductInt(v int) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "int",
		Value: v,
	}
}
func ProductIntPointer(v *int) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "int_pointer",
		Value: v,
	}
}
func ProductFloat(v float64) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "float",
		Value: v,
	}
}
func ProductBool(v bool) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "bool",
		Value: v,
	}
}
func ProductByteArray(v []uint8) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "byte_array",
		Value: v,
	}
}
func ProductMultiString(v models.MultiString) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "multi_string",
		Value: v,
	}
}
func ProductStringArray(v pq.StringArray) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "string_array",
		Value: v,
	}
}
func ProductEmbeddedInt(v int) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "embedded_int",
		Value: v,
	}
}
func ProductGormEmbeddedInt(v int) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field: "gorm_embedded_int",
		Value: v,
	}
}
