// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	uintmodel "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/uintmodel"
	gorm "gorm.io/gorm"
	"time"
)

func UintModelId(expr badorm.Expression[uint]) badorm.WhereCondition[uintmodel.UintModel] {
	return badorm.FieldCondition[uintmodel.UintModel, uint]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func UintModelCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[uintmodel.UintModel] {
	return badorm.FieldCondition[uintmodel.UintModel, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func UintModelUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[uintmodel.UintModel] {
	return badorm.FieldCondition[uintmodel.UintModel, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func UintModelDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[uintmodel.UintModel] {
	return badorm.FieldCondition[uintmodel.UintModel, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}
