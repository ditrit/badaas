// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	belongsto "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/belongsto"
	gorm "gorm.io/gorm"
	"time"
)

func OwnedId(exprs ...badorm.Expression[badorm.UUID]) badorm.FieldCondition[belongsto.Owned, badorm.UUID] {
	return badorm.FieldCondition[belongsto.Owned, badorm.UUID]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func OwnedCreatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[belongsto.Owned, time.Time] {
	return badorm.FieldCondition[belongsto.Owned, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func OwnedUpdatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[belongsto.Owned, time.Time] {
	return badorm.FieldCondition[belongsto.Owned, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func OwnedDeletedAt(exprs ...badorm.Expression[gorm.DeletedAt]) badorm.FieldCondition[belongsto.Owned, gorm.DeletedAt] {
	return badorm.FieldCondition[belongsto.Owned, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func OwnedOwner(conditions ...badorm.Condition[belongsto.Owner]) badorm.Condition[belongsto.Owned] {
	return badorm.JoinCondition[belongsto.Owned, belongsto.Owner]{
		Conditions: conditions,
		T1Field:    "OwnerID",
		T2Field:    "ID",
	}
}
func OwnedOwnerId(exprs ...badorm.Expression[badorm.UUID]) badorm.FieldCondition[belongsto.Owned, badorm.UUID] {
	return badorm.FieldCondition[belongsto.Owned, badorm.UUID]{
		Expressions: exprs,
		Field:       "OwnerID",
	}
}