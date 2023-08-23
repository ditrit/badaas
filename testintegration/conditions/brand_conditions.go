// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	condition "github.com/ditrit/badaas/orm/condition"
	model "github.com/ditrit/badaas/orm/model"
	operator "github.com/ditrit/badaas/orm/operator"
	query "github.com/ditrit/badaas/orm/query"
	models "github.com/ditrit/badaas/testintegration/models"
	"reflect"
	"time"
)

var brandType = reflect.TypeOf(*new(models.Brand))
var BrandIdField = query.FieldIdentifier[model.UIntID]{
	Field:     "ID",
	ModelType: brandType,
}

func BrandId(operator operator.Operator[model.UIntID]) condition.WhereCondition[models.Brand] {
	return condition.NewFieldCondition[models.Brand, model.UIntID](BrandIdField, operator)
}

var BrandCreatedAtField = query.FieldIdentifier[time.Time]{
	Field:     "CreatedAt",
	ModelType: brandType,
}

func BrandCreatedAt(operator operator.Operator[time.Time]) condition.WhereCondition[models.Brand] {
	return condition.NewFieldCondition[models.Brand, time.Time](BrandCreatedAtField, operator)
}

var BrandUpdatedAtField = query.FieldIdentifier[time.Time]{
	Field:     "UpdatedAt",
	ModelType: brandType,
}

func BrandUpdatedAt(operator operator.Operator[time.Time]) condition.WhereCondition[models.Brand] {
	return condition.NewFieldCondition[models.Brand, time.Time](BrandUpdatedAtField, operator)
}

var BrandDeletedAtField = query.FieldIdentifier[time.Time]{
	Field:     "DeletedAt",
	ModelType: brandType,
}

func BrandDeletedAt(operator operator.Operator[time.Time]) condition.WhereCondition[models.Brand] {
	return condition.NewFieldCondition[models.Brand, time.Time](BrandDeletedAtField, operator)
}

var BrandNameField = query.FieldIdentifier[string]{
	Field:     "Name",
	ModelType: brandType,
}

func BrandName(operator operator.Operator[string]) condition.WhereCondition[models.Brand] {
	return condition.NewFieldCondition[models.Brand, string](BrandNameField, operator)
}

var BrandPreloadAttributes = condition.NewPreloadCondition[models.Brand](BrandIdField, BrandCreatedAtField, BrandUpdatedAtField, BrandDeletedAtField, BrandNameField)
