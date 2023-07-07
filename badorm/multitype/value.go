package multitype

import (
	"database/sql"
	"reflect"

	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	badormSQL "github.com/ditrit/badaas/badorm/sql"
)

var nullableTypes = []reflect.Type{
	reflect.TypeOf(sql.NullBool{}),
	reflect.TypeOf(sql.NullByte{}),
	reflect.TypeOf(sql.NullFloat64{}),
	reflect.TypeOf(sql.NullInt16{}),
	reflect.TypeOf(sql.NullInt32{}),
	reflect.TypeOf(sql.NullInt64{}),
	reflect.TypeOf(sql.NullString{}),
	reflect.TypeOf(sql.NullTime{}),
	reflect.TypeOf(gorm.DeletedAt{}),
}

func isNullable(fieldType reflect.Type) bool {
	return pie.Contains(nullableTypes, fieldType)
}

func verifyFieldType[TAttribute, TField any](sqlOperator badormSQL.Operator) badorm.DynamicOperator[TAttribute] {
	attributeType := reflect.TypeOf(*new(TAttribute))
	fieldType := reflect.TypeOf(*new(TField))

	if fieldType != attributeType &&
		!((isNullable(fieldType) && fieldType.Field(0).Type == attributeType) ||
			(isNullable(attributeType) && attributeType.Field(0).Type == fieldType)) {
		return badorm.NewInvalidOperator[TAttribute](
			fieldTypeDoesNotMatchError(fieldType, attributeType, sqlOperator),
		)
	}

	return nil
}

func NewValueOperator[TAttribute, TField any](
	sqlOperator badormSQL.Operator,
	field badorm.FieldIdentifier[TField],
) badorm.DynamicOperator[TAttribute] {
	invalidOperator := verifyFieldType[TAttribute, TField](sqlOperator)
	if invalidOperator != nil {
		return invalidOperator
	}

	return &badorm.ValueOperator[TAttribute]{
		Operations: []badorm.Operation{
			{
				SQLOperator: sqlOperator,
				Value:       field,
			},
		},
		JoinNumber: badorm.UndefinedJoinNumber,
	}
}