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
var (
	ErrFieldTypeDoesNotMatch = errors.New("field type does not match expression type")
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

func verifyFieldType[TAttribute, TField any]() badorm.DynamicExpression[TAttribute] {
	expressionType := reflect.TypeOf(*new(TAttribute))
	fieldType := reflect.TypeOf(*new(TField))

	if fieldType != expressionType &&
		!((isNullable(fieldType) && fieldType.Field(0).Type == expressionType) ||
			(isNullable(expressionType) && expressionType.Field(0).Type == fieldType)) {
		return badorm.NewInvalidExpression[TAttribute](ErrFieldTypeDoesNotMatch)
	}

	return nil
}

func newValueExpression[TAttribute, TField any](
	expression expressions.SQLExpression,
	field badorm.FieldIdentifier[TField],
) badorm.DynamicExpression[TAttribute] {
	invalidExpression := verifyFieldType[TAttribute, TField]()
	if invalidExpression != nil {
		return invalidExpression
	}

	// TODO que pasa con los que solo aceptan cierto tipo, ver las de like por ejemplo
	// TODO que pasa si hay mas de uno, no se si me gusta mucho esto
	// creo que si hay mas de uno solo se tendria que aplicar en el primero,
	// TODO soportar multivalue, no todos necesariamente dinamicos
	return &dynamic.ValueExpression[TAttribute]{
		SQLExpressions: []dynamic.LiteralSQLExpression{
			{
				SQL:   expressions.ToSQL[expression],
				Field: field,
			},
		},
		JoinNumber: dynamic.UndefinedJoinNumber,
	}
}

func newMultivalueExpression[TAttribute, TField any](
	sqlExpression expressions.SQLExpression,
	sqlConnector, sqlPrefix, sqlSuffix string,
	values ...any,
) badorm.DynamicExpression[TAttribute] {
	for _, value := range values {
		_, isT1 := value.(TAttribute)
		if isT1 {
			continue
		}

		_, isField := value.(badorm.FieldIdentifier[TField])
		if isField {
			invalidExpression := verifyFieldType[TAttribute, TField]()
			if invalidExpression != nil {
				return invalidExpression
			}

			continue
		}

		return badorm.NewInvalidExpression[TAttribute](ErrParamsNotValueOrField)
	}

	return &dynamic.MultivalueExpression[TAttribute]{
		Values:        values,
		SQLExpression: expressions.ToSQL[sqlExpression],
		SQLConnector:  sqlConnector,
		SQLPrefix:     sqlPrefix,
		SQLSuffix:     sqlSuffix,
		JoinNumber:    dynamic.UndefinedJoinNumber,
	}
}

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
func Eq[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicExpression[TAttribute] {
	return newValueExpression[TAttribute, TField](expressions.Eq, field)
}

// NotEqualTo
func NotEq[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicExpression[TAttribute] {
	return newValueExpression[TAttribute, TField](expressions.NotEq, field)
}

// LessThan
func Lt[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicExpression[TAttribute] {
	return newValueExpression[TAttribute, TField](expressions.Lt, field)
}

// LessThanOrEqualTo
func LtOrEq[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicExpression[TAttribute] {
	return newValueExpression[TAttribute, TField](expressions.LtOrEq, field)
}

// GreaterThan
func Gt[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicExpression[TAttribute] {
	// TODO invertir orden de parametros para que quede igual que en badorm/expression
	return newValueExpression[TAttribute, TField](expressions.Gt, field)
}

// GreaterThanOrEqualTo
func GtOrEq[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicExpression[TAttribute] {
	return newValueExpression[TAttribute, TField](expressions.GtOrEq, field)
}

// Comparison Predicates
// refs:
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to v1 < value < v2
func Between[TAttribute, TField any](v1 any, v2 any) badorm.DynamicExpression[TAttribute] {
	return newMultivalueExpression[TAttribute, TField](expressions.Between, "AND", "", "", v1, v2)
}

// Equivalent to NOT (field1 < value < field2)
func NotBetween[TAttribute, TField any](v1 any, v2 any) badorm.DynamicExpression[TAttribute] {
	return newMultivalueExpression[TAttribute, TField](expressions.NotBetween, "AND", "", "", v1, v2)
}

// Boolean Comparison Predicates

// Not supported by: mysql
func IsDistinct[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicExpression[TAttribute] {
	return newValueExpression[TAttribute, TField](expressions.IsDistinct, field)
}

// Not supported by: mysql
func IsNotDistinct[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicExpression[TAttribute] {
	return newValueExpression[TAttribute, TField](expressions.IsNotDistinct, field)
}
