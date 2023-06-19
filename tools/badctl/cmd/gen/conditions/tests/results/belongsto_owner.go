// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	belongsto "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/belongsto"
	gorm "gorm.io/gorm"
	"time"
)

func OwnerId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[belongsto.Owner] {
	return badorm.FieldCondition[belongsto.Owner, badorm.UUID]{
		Expression: expr,
		Field:      "ID",
	}
}
func OwnerCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[belongsto.Owner] {
	return badorm.FieldCondition[belongsto.Owner, time.Time]{
		Expression: expr,
		Field:      "CreatedAt",
	}
}
func OwnerUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[belongsto.Owner] {
	return badorm.FieldCondition[belongsto.Owner, time.Time]{
		Expression: expr,
		Field:      "UpdatedAt",
	}
}
func OwnerDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[belongsto.Owner] {
	return badorm.FieldCondition[belongsto.Owner, gorm.DeletedAt]{
		Expression: expr,
		Field:      "DeletedAt",
	}
}
