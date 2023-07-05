package multitype

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/dynamic"
	"github.com/ditrit/badaas/badorm/expressions"
)

// TODO ponerles mas informacion precisa a los errores
var ErrFieldTypeDoesNotMatch = errors.New("field type does not match expression type")

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

func newMultitypeValueExpression[T1 any, T2 any](expression expressions.SQLExpression, field badorm.FieldIdentifier[T2]) badorm.DynamicExpression[T1] {
	expressionType := reflect.TypeOf(*new(T1))
	fieldType := reflect.TypeOf(*new(T2))

	if fieldType != expressionType &&
		!((isNullable(fieldType) && fieldType.Field(0).Type == expressionType) ||
			(isNullable(expressionType) && expressionType.Field(0).Type == fieldType)) {
		return badorm.NewInvalidExpression[T1](ErrFieldTypeDoesNotMatch)
	}

	// TODO que pasa con los que solo aceptan cierto tipo, ver las de like por ejemplo
	// TODO que pasa si hay mas de uno, no se si me gusta mucho esto
	// creo que si hay mas de uno solo se tendria que aplicar en el primero,
	return &dynamic.ValueExpression[T1]{
		SQLExpressions: []dynamic.LiteralSQLExpression{
			{
				SQL:   expressions.ToSQL[expression],
				Field: field,
			},
		},
		JoinNumber: dynamic.UndefinedJoinNumber,
	}
	// TODO soportar multivalue, no todos necesariamente dinamicos
}

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
func Eq[T1 any, T2 any](field badorm.FieldIdentifier[T2]) badorm.DynamicExpression[T1] {
	return newMultitypeValueExpression[T1, T2](expressions.Eq, field)
}

// NotEqualTo
func NotEq[T1 any, T2 any](field badorm.FieldIdentifier[T2]) badorm.DynamicExpression[T1] {
	return newMultitypeValueExpression[T1, T2](expressions.NotEq, field)
}

// LessThan
func Lt[T1 any, T2 any](field badorm.FieldIdentifier[T2]) badorm.DynamicExpression[T1] {
	return newMultitypeValueExpression[T1, T2](expressions.Lt, field)
}

// LessThanOrEqualTo
func LtOrEq[T1 any, T2 any](field badorm.FieldIdentifier[T2]) badorm.DynamicExpression[T1] {
	return newMultitypeValueExpression[T1, T2](expressions.LtOrEq, field)
}

// GreaterThan
func Gt[T1 any, T2 any](field badorm.FieldIdentifier[T2]) badorm.DynamicExpression[T1] {
	// TODO invertir orden de parametros para que quede igual que en badorm/expression
	return newMultitypeValueExpression[T1, T2](expressions.Gt, field)
}

// GreaterThanOrEqualTo
func GtOrEq[T1 any, T2 any](field badorm.FieldIdentifier[T2]) badorm.DynamicExpression[T1] {
	return newMultitypeValueExpression[T1, T2](expressions.GtOrEq, field)
}

// Comparison Predicates
// refs:
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to v1 < value < v2
// func Between[T any](v1 T, v2 T) MultivalueExpression[T] {
// 	return NewMultivalueExpression("BETWEEN", "AND", "", "", v1, v2)
// }

// // Equivalent to NOT (v1 < value < v2)
// func NotBetween[T any](v1 T, v2 T) MultivalueExpression[T] {
// 	return NewMultivalueExpression("NOT BETWEEN", "AND", "", "", v1, v2)
// }

// Boolean Comparison Predicates

// Not supported by: mysql
func IsDistinct[T1 any, T2 any](field badorm.FieldIdentifier[T2]) badorm.DynamicExpression[T1] {
	return newMultitypeValueExpression[T1, T2](expressions.IsDistinct, field)
}

// Not supported by: mysql
func IsNotDistinct[T1 any, T2 any](field badorm.FieldIdentifier[T2]) badorm.DynamicExpression[T1] {
	return newMultitypeValueExpression[T1, T2](expressions.IsNotDistinct, field)
}
