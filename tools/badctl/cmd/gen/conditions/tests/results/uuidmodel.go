// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	uuidmodel "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/uuidmodel"
	gorm "gorm.io/gorm"
	"time"
)

func UUIDModelId(exprs ...badorm.Expression[badorm.UUID]) badorm.FieldCondition[uuidmodel.UUIDModel, badorm.UUID] {
	return badorm.FieldCondition[uuidmodel.UUIDModel, badorm.UUID]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func UUIDModelCreatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[uuidmodel.UUIDModel, time.Time] {
	return badorm.FieldCondition[uuidmodel.UUIDModel, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func UUIDModelUpdatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[uuidmodel.UUIDModel, time.Time] {
	return badorm.FieldCondition[uuidmodel.UUIDModel, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func UUIDModelDeletedAt(exprs ...badorm.Expression[gorm.DeletedAt]) badorm.FieldCondition[uuidmodel.UUIDModel, gorm.DeletedAt] {
	return badorm.FieldCondition[uuidmodel.UUIDModel, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
