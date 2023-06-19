// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	uuidmodel "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/uuidmodel"
	gorm "gorm.io/gorm"
	"time"
)

func UUIDModelId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[uuidmodel.UUIDModel] {
	return badorm.FieldCondition[uuidmodel.UUIDModel, badorm.UUID]{
		Expression: expr,
		Field:      "ID",
	}
}
func UUIDModelCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[uuidmodel.UUIDModel] {
	return badorm.FieldCondition[uuidmodel.UUIDModel, time.Time]{
		Expression: expr,
		Field:      "CreatedAt",
	}
}
func UUIDModelUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[uuidmodel.UUIDModel] {
	return badorm.FieldCondition[uuidmodel.UUIDModel, time.Time]{
		Expression: expr,
		Field:      "UpdatedAt",
	}
}
func UUIDModelDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[uuidmodel.UUIDModel] {
	return badorm.FieldCondition[uuidmodel.UUIDModel, gorm.DeletedAt]{
		Expression: expr,
		Field:      "DeletedAt",
	}
}
