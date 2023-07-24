package multitype

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

var (
	ErrFieldTypeDoesNotMatch = errors.New("field type does not match attribute type")
	ErrParamNotValueOrField  = errors.New("parameter is neither a value nor a field of the attribute type")
)

func fieldTypeDoesNotMatchError(fieldType, attributeType reflect.Type, sqlOperator sql.Operator) error {
	return badorm.OperatorError(fmt.Errorf("%w; field type: %s, attribute type: %s",
		ErrFieldTypeDoesNotMatch,
		fieldType,
		attributeType,
	), sqlOperator)
}

func paramNotValueOrField[T any](value any, sqlOperator sql.Operator) error {
	return badorm.OperatorError(fmt.Errorf("%w; parameter type: %T, attribute type: %T",
		ErrParamNotValueOrField,
		value,
		*new(T),
	), sqlOperator)
}
