// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func Parent1Id(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Parent1] {
	return badorm.FieldCondition[models.Parent1, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func Parent1CreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Parent1] {
	return badorm.FieldCondition[models.Parent1, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func Parent1UpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Parent1] {
	return badorm.FieldCondition[models.Parent1, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func Parent1DeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[models.Parent1] {
	return badorm.FieldCondition[models.Parent1, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}
func Parent1ParentParent(conditions ...badorm.Condition[models.ParentParent]) badorm.Condition[models.Parent1] {
	return badorm.JoinCondition[models.Parent1, models.ParentParent]{
		Conditions:    conditions,
		RelationField: "ParentParent",
		T1Field:       "ParentParentID",
		T2Field:       "ID",
	}
}

var Parent1PreloadParentParent = Parent1ParentParent(ParentParentPreloadAttributes)
var parent1ParentParentIdFieldID = badorm.FieldIdentifier{Field: "ParentParentID"}

func Parent1ParentParentId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Parent1] {
	return badorm.FieldCondition[models.Parent1, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: parent1ParentParentIdFieldID,
	}
}

var Parent1PreloadAttributes = badorm.NewPreloadCondition[models.Parent1](parent1ParentParentIdFieldID)
var Parent1PreloadRelations = []badorm.Condition[models.Parent1]{Parent1PreloadParentParent}
