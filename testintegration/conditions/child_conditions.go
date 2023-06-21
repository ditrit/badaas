// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func ChildId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Child] {
	return badorm.FieldCondition[models.Child, badorm.UUID]{
		Expression: expr,
		Field:      "ID",
	}
}
func ChildCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Child] {
	return badorm.FieldCondition[models.Child, time.Time]{
		Expression: expr,
		Field:      "CreatedAt",
	}
}
func ChildUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Child] {
	return badorm.FieldCondition[models.Child, time.Time]{
		Expression: expr,
		Field:      "UpdatedAt",
	}
}
func ChildDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[models.Child] {
	return badorm.FieldCondition[models.Child, gorm.DeletedAt]{
		Expression: expr,
		Field:      "DeletedAt",
	}
}
func ChildParent1(conditions ...badorm.Condition[models.Parent1]) badorm.Condition[models.Child] {
	return badorm.JoinCondition[models.Child, models.Parent1]{
		Conditions:    conditions,
		RelationField: "Parent1",
		T1Field:       "Parent1ID",
		T2Field:       "ID",
	}
}
func ChildParent1Id(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Child] {
	return badorm.FieldCondition[models.Child, badorm.UUID]{
		Expression: expr,
		Field:      "Parent1ID",
	}
}
func ChildParent2(conditions ...badorm.Condition[models.Parent2]) badorm.Condition[models.Child] {
	return badorm.JoinCondition[models.Child, models.Parent2]{
		Conditions:    conditions,
		RelationField: "Parent2",
		T1Field:       "Parent2ID",
		T2Field:       "ID",
	}
}
func ChildParent2Id(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Child] {
	return badorm.FieldCondition[models.Child, badorm.UUID]{
		Expression: expr,
		Field:      "Parent2ID",
	}
}
