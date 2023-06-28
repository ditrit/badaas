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
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func ChildCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Child] {
	return badorm.FieldCondition[models.Child, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func ChildUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Child] {
	return badorm.FieldCondition[models.Child, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func ChildDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[models.Child] {
	return badorm.FieldCondition[models.Child, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}
func ChildParent1(conditions ...badorm.Condition[models.Parent1]) badorm.IJoinCondition[models.Child] {
	return badorm.JoinCondition[models.Child, models.Parent1]{
		Conditions:         conditions,
		RelationField:      "Parent1",
		T1Field:            "Parent1ID",
		T1PreloadCondition: ChildPreloadAttributes,
		T2Field:            "ID",
	}
}

var ChildPreloadParent1 = ChildParent1(Parent1PreloadAttributes)
var childParent1IdFieldID = badorm.FieldIdentifier{Field: "Parent1ID"}

func ChildParent1Id(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Child] {
	return badorm.FieldCondition[models.Child, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: childParent1IdFieldID,
	}
}
func ChildParent2(conditions ...badorm.Condition[models.Parent2]) badorm.IJoinCondition[models.Child] {
	return badorm.JoinCondition[models.Child, models.Parent2]{
		Conditions:         conditions,
		RelationField:      "Parent2",
		T1Field:            "Parent2ID",
		T1PreloadCondition: ChildPreloadAttributes,
		T2Field:            "ID",
	}
}

var ChildPreloadParent2 = ChildParent2(Parent2PreloadAttributes)
var childParent2IdFieldID = badorm.FieldIdentifier{Field: "Parent2ID"}

func ChildParent2Id(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Child] {
	return badorm.FieldCondition[models.Child, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: childParent2IdFieldID,
	}
}

var ChildPreloadAttributes = badorm.NewPreloadCondition[models.Child](childParent1IdFieldID, childParent2IdFieldID)
var ChildPreloadRelations = []badorm.Condition[models.Child]{ChildPreloadParent1, ChildPreloadParent2}
