package badorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"

	badormSQL "github.com/ditrit/badaas/badorm/sql"
)

// TODO ponerles mas informacion precisa a los errores
var (
	ErrFieldModelNotConcerned = errors.New("field's model is not concerned by the query so it can't be used in a operator")
	ErrJoinMustBeSelected     = errors.New("table is joined more than once, select which one you want to use")
	ErrValueCantBeNull        = errors.New("value to compare can't be null")
)

const UndefinedJoinNumber = -1

// Operator that compares the value of the column against a fixed value
// If SQLOperators has multiple entries, comparisons will be nested
// Example (single): value = v1
// Example (multi): value LIKE v1 ESCAPE v2
type ValueOperator[T any] struct {
	Operations []Operation
	// TODO join deberia estar en cada operator
	JoinNumber int
}

type Operation struct {
	SQLOperator badormSQL.Operator
	Value       any
}

func (expr ValueOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

func (expr *ValueOperator[T]) SelectJoin(joinNumber uint) DynamicOperator[T] {
	expr.JoinNumber = int(joinNumber)
	return expr
}

func (expr *ValueOperator[T]) AddOperation(value any, sqlOperator badormSQL.Operator) ValueOperator[T] {
	expr.Operations = append(
		expr.Operations,
		Operation{
			Value:       value,
			SQLOperator: sqlOperator,
		},
	)

	return *expr
}

// verificar que en las condiciones anteriores alguien usó el field con el que se intenta comparar
// obtener de ahi cual es el nombre de la table a usar con ese field.
// TODO doc a ingles
func (expr ValueOperator[T]) ToSQL(query *Query, columnName string) (string, []any, error) {
	operationString := columnName
	values := []any{}

	for _, operation := range expr.Operations {
		field, isField := operation.Value.(IFieldIdentifier)
		if isField {
			modelTable, err := getModelTable(query, field, expr.JoinNumber)
			if err != nil {
				return "", nil, err
			}

			operationString += fmt.Sprintf(
				" %s %s",
				operation.SQLOperator,
				field.ColumnSQL(query, modelTable),
			)
		} else {
			operationString += " " + operation.SQLOperator.String() + " ?"
			values = append(values, operation.Value)
		}
	}

	return operationString, values, nil
}

func NewValueOperator[T any](value any, sqlOperator badormSQL.Operator) ValueOperator[T] {
	expr := &ValueOperator[T]{}

	return expr.AddOperation(value, sqlOperator)
}

var nullableKinds = []reflect.Kind{
	reflect.Chan, reflect.Func,
	reflect.Map, reflect.Pointer,
	reflect.UnsafePointer, reflect.Interface,
	reflect.Slice,
}

func NewCantBeNullValueOperator[T any](value any, sqlOperator badormSQL.Operator) Operator[T] {
	if value == nil || mapsToNull(value) {
		return NewInvalidOperator[T](ErrValueCantBeNull)
	}

	return NewValueOperator[T](value, sqlOperator)
}

func NewMustBePOSIXValueOperator[T string | sql.NullString](pattern string, sqlOperator badormSQL.Operator) Operator[T] {
	_, err := regexp.CompilePOSIX(pattern)
	if err != nil {
		return NewInvalidOperator[T](err)
	}

	return NewValueOperator[T](pattern, sqlOperator)
}
