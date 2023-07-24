// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"reflect"
	"time"
)

var parent2Type = reflect.TypeOf(*new(models.Parent2))
var Parent2IdField = badorm.FieldIdentifier[badorm.UUID]{
	Field:     "ID",
	ModelType: parent2Type,
}

func Parent2Id(operator badorm.Operator[badorm.UUID]) badorm.WhereCondition[models.Parent2] {
	return badorm.FieldCondition[models.Parent2, badorm.UUID]{
		FieldIdentifier: Parent2IdField,
		Operator:        operator,
	}
}

var Parent2CreatedAtField = badorm.FieldIdentifier[time.Time]{
	Field:     "CreatedAt",
	ModelType: parent2Type,
}

func Parent2CreatedAt(operator badorm.Operator[time.Time]) badorm.WhereCondition[models.Parent2] {
	return badorm.FieldCondition[models.Parent2, time.Time]{
		FieldIdentifier: Parent2CreatedAtField,
		Operator:        operator,
	}
}

var Parent2UpdatedAtField = badorm.FieldIdentifier[time.Time]{
	Field:     "UpdatedAt",
	ModelType: parent2Type,
}

func Parent2UpdatedAt(operator badorm.Operator[time.Time]) badorm.WhereCondition[models.Parent2] {
	return badorm.FieldCondition[models.Parent2, time.Time]{
		FieldIdentifier: Parent2UpdatedAtField,
		Operator:        operator,
	}
}

var Parent2DeletedAtField = badorm.FieldIdentifier[gorm.DeletedAt]{
	Field:     "DeletedAt",
	ModelType: parent2Type,
}

func Parent2DeletedAt(operator badorm.Operator[gorm.DeletedAt]) badorm.WhereCondition[models.Parent2] {
	return badorm.FieldCondition[models.Parent2, gorm.DeletedAt]{
		FieldIdentifier: Parent2DeletedAtField,
		Operator:        operator,
	}
}
func Parent2ParentParent(conditions ...badorm.Condition[models.ParentParent]) badorm.IJoinCondition[models.Parent2] {
	return badorm.JoinCondition[models.Parent2, models.ParentParent]{
		Conditions:         conditions,
		RelationField:      "ParentParent",
		T1Field:            "ParentParentID",
		T1PreloadCondition: Parent2PreloadAttributes,
		T2Field:            "ID",
	}
}

var Parent2PreloadParentParent = Parent2ParentParent(ParentParentPreloadAttributes)
var Parent2ParentParentIdField = badorm.FieldIdentifier[badorm.UUID]{
	Field:     "ParentParentID",
	ModelType: parent2Type,
}

func Parent2ParentParentId(operator badorm.Operator[badorm.UUID]) badorm.WhereCondition[models.Parent2] {
	return badorm.FieldCondition[models.Parent2, badorm.UUID]{
		FieldIdentifier: Parent2ParentParentIdField,
		Operator:        operator,
	}
}

var Parent2PreloadAttributes = badorm.NewPreloadCondition[models.Parent2](Parent2IdField, Parent2CreatedAtField, Parent2UpdatedAtField, Parent2DeletedAtField, Parent2ParentParentIdField)
var Parent2PreloadRelations = []badorm.Condition[models.Parent2]{Parent2PreloadParentParent}
