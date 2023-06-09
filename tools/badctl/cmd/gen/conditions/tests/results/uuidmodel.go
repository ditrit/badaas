// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	uuidmodel "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/uuidmodel"
	gorm "gorm.io/gorm"
	"time"
)

func UUIDModelId(v badorm.UUID) badorm.WhereCondition[uuidmodel.UUIDModel] {
	return badorm.WhereCondition[uuidmodel.UUIDModel]{
		Field: "ID",
		Value: v,
	}
}
func UUIDModelCreatedAt(v time.Time) badorm.WhereCondition[uuidmodel.UUIDModel] {
	return badorm.WhereCondition[uuidmodel.UUIDModel]{
		Field: "CreatedAt",
		Value: v,
	}
}
func UUIDModelUpdatedAt(v time.Time) badorm.WhereCondition[uuidmodel.UUIDModel] {
	return badorm.WhereCondition[uuidmodel.UUIDModel]{
		Field: "UpdatedAt",
		Value: v,
	}
}
func UUIDModelDeletedAt(v gorm.DeletedAt) badorm.WhereCondition[uuidmodel.UUIDModel] {
	return badorm.WhereCondition[uuidmodel.UUIDModel]{
		Field: "DeletedAt",
		Value: v,
	}
}