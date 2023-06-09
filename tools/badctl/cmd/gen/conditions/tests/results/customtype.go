// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	customtype "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/customtype"
	gorm "gorm.io/gorm"
	"time"
)

func CustomTypeId(v badorm.UUID) badorm.WhereCondition[customtype.CustomType] {
	return badorm.WhereCondition[customtype.CustomType]{
		Field: "ID",
		Value: v,
	}
}
func CustomTypeCreatedAt(v time.Time) badorm.WhereCondition[customtype.CustomType] {
	return badorm.WhereCondition[customtype.CustomType]{
		Field: "CreatedAt",
		Value: v,
	}
}
func CustomTypeUpdatedAt(v time.Time) badorm.WhereCondition[customtype.CustomType] {
	return badorm.WhereCondition[customtype.CustomType]{
		Field: "UpdatedAt",
		Value: v,
	}
}
func CustomTypeDeletedAt(v gorm.DeletedAt) badorm.WhereCondition[customtype.CustomType] {
	return badorm.WhereCondition[customtype.CustomType]{
		Field: "DeletedAt",
		Value: v,
	}
}
func CustomTypeCustom(v customtype.MultiString) badorm.WhereCondition[customtype.CustomType] {
	return badorm.WhereCondition[customtype.CustomType]{
		Field: "Custom",
		Value: v,
	}
}