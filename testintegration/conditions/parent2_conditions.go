// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func Parent2Id(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Parent2] {
	return badorm.FieldCondition[models.Parent2, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func Parent2CreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Parent2] {
	return badorm.FieldCondition[models.Parent2, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func Parent2UpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Parent2] {
	return badorm.FieldCondition[models.Parent2, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func Parent2DeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[models.Parent2] {
	return badorm.FieldCondition[models.Parent2, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}
func Parent2ParentParent(conditions ...badorm.Condition[models.ParentParent]) badorm.Condition[models.Parent2] {
	return badorm.JoinCondition[models.Parent2, models.ParentParent]{
		Conditions:    conditions,
		RelationField: "ParentParent",
		T1Field:       "ParentParentID",
		T2Field:       "ID",
	}
}
func Parent2ParentParentId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Parent2] {
	return badorm.FieldCondition[models.Parent2, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.FieldIdentifier{Field: "ParentParentID"},
	}
}
