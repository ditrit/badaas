// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	belongsto "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/belongsto"
	gorm "gorm.io/gorm"
	"time"
)

func OwnedId(v badorm.UUID) badorm.WhereCondition[belongsto.Owned] {
	return badorm.WhereCondition[belongsto.Owned]{
		Field: "ID",
		Value: v,
	}
}
func OwnedCreatedAt(v time.Time) badorm.WhereCondition[belongsto.Owned] {
	return badorm.WhereCondition[belongsto.Owned]{
		Field: "CreatedAt",
		Value: v,
	}
}
func OwnedUpdatedAt(v time.Time) badorm.WhereCondition[belongsto.Owned] {
	return badorm.WhereCondition[belongsto.Owned]{
		Field: "UpdatedAt",
		Value: v,
	}
}
func OwnedDeletedAt(v gorm.DeletedAt) badorm.WhereCondition[belongsto.Owned] {
	return badorm.WhereCondition[belongsto.Owned]{
		Field: "DeletedAt",
		Value: v,
	}
}
func OwnedOwner(conditions ...badorm.Condition[belongsto.Owner]) badorm.Condition[belongsto.Owned] {
	return badorm.JoinCondition[belongsto.Owned, belongsto.Owner]{
		Conditions: conditions,
		T1Field:    "OwnerID",
		T2Field:    "ID",
	}
}
func OwnedOwnerId(v badorm.UUID) badorm.WhereCondition[belongsto.Owned] {
	return badorm.WhereCondition[belongsto.Owned]{
		Field: "OwnerID",
		Value: v,
	}
}
