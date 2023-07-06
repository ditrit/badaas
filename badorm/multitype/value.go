package multitype

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	badormSQL "github.com/ditrit/badaas/badorm/sql"
)

// TODO ponerles mas informacion precisa a los errores
var (
	ErrFieldTypeDoesNotMatch = errors.New("field type does not match operator type")
	ErrParamsNotValueOrField = errors.New("parameter is neither a value or a field")
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

func verifyFieldType[TAttribute, TField any]() badorm.DynamicOperator[TAttribute] {
	attributeType := reflect.TypeOf(*new(TAttribute))
	fieldType := reflect.TypeOf(*new(TField))

	if fieldType != attributeType &&
		!((isNullable(fieldType) && fieldType.Field(0).Type == attributeType) ||
			(isNullable(attributeType) && attributeType.Field(0).Type == fieldType)) {
		return badorm.NewInvalidOperator[TAttribute](ErrFieldTypeDoesNotMatch)
	}

	return nil
}

func NewValueOperator[TAttribute, TField any](
	sqlOperator badormSQL.Operator,
	field badorm.FieldIdentifier[TField],
) badorm.DynamicOperator[TAttribute] {
	invalidOperator := verifyFieldType[TAttribute, TField]()
	if invalidOperator != nil {
		return invalidOperator
	}

	// TODO que pasa con los que solo aceptan cierto tipo, ver las de like por ejemplo
	// TODO que pasa si hay mas de uno, no se si me gusta mucho esto
	// creo que si hay mas de uno solo se tendria que aplicar en el primero,
	// TODO soportar multivalue, no todos necesariamente dinamicos
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
