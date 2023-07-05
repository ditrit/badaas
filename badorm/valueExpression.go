package badorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/ditrit/badaas/badorm/expressions"
)

// TODO ponerles mas informacion precisa a los errores
var (
	ErrFieldModelNotConcerned = errors.New("field's model is not concerned by the query so it can't be used in a expression")
	ErrJoinMustBeSelected     = errors.New("table is joined more than once, select which one you want to use")
)
var ErrValueCantBeNull = errors.New("value to compare can't be null")

const UndefinedJoinNumber = -1

// Expression that compares the value of the column against a fixed value
// If ExpressionsAndValues has multiple entries, comparisons will be nested
// Example (single): value = v1
// Example (multi): value LIKE v1 ESCAPE v2
type ValueExpression[T any] struct {
	SQLExpressions []SQLExpressionAndValue
	JoinNumber     int
}

type SQLExpressionAndValue struct {
	SQL   string
	Field IFieldIdentifier
	Value any
}

func (expr ValueExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

func (expr *ValueExpression[T]) SelectJoin(joinNumber uint) DynamicExpression[T] {
	expr.JoinNumber = int(joinNumber)
	return expr
}

func (expr *ValueExpression[T]) AddSQLExpression(value any, sqlExpression expressions.SQLExpression) ValueExpression[T] {
	expr.SQLExpressions = append(
		expr.SQLExpressions,
		SQLExpressionAndValue{
			Value: value,
			SQL:   expressions.ToSQL[sqlExpression],
		},
	)

	return *expr
}

// verificar que en las condiciones anteriores alguien usó el field con el que se intenta comparar
// obtener de ahi cual es el nombre de la table a usar con ese field.
// TODO doc a ingles
func (expr ValueExpression[T]) ToSQL(query *Query, columnName string) (string, []any, error) {
	exprString := columnName
	values := []any{}

	for _, sqlExpr := range expr.SQLExpressions {
		if sqlExpr.Field != nil {
			modelTable, err := getModelTable(query, sqlExpr.Field, expr.JoinNumber)
			if err != nil {
				return "", nil, err
			}

			exprString += fmt.Sprintf(
				" "+sqlExpr.SQL+" %s",
				sqlExpr.Field.ColumnSQL(query, modelTable),
			)
		} else {
			exprString += " " + sqlExpr.SQL + " ?"
			values = append(values, sqlExpr.Value)
		}
	}

	return exprString, values, nil
}

func NewValueExpression[T any](value any, sqlExpression expressions.SQLExpression) ValueExpression[T] {
	expr := &ValueExpression[T]{}

	return expr.AddSQLExpression(value, sqlExpression)
}

var nullableKinds = []reflect.Kind{
	reflect.Chan, reflect.Func,
	reflect.Map, reflect.Pointer,
	reflect.UnsafePointer, reflect.Interface,
	reflect.Slice,
}

func NewCantBeNullValueExpression[T any](value any, sqlExpression expressions.SQLExpression) Expression[T] {
	if value == nil || mapsToNull(value) {
		return NewInvalidExpression[T](ErrValueCantBeNull)
	}

	return NewValueExpression[T](value, sqlExpression)
}

func NewMustBePOSIXValueExpression[T string | sql.NullString](pattern string, sqlExpression expressions.SQLExpression) Expression[T] {
	_, err := regexp.CompilePOSIX(pattern)
	if err != nil {
		return NewInvalidExpression[T](err)
	}

	return NewValueExpression[T](pattern, sqlExpression)
}
