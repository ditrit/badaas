package dynamic

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/expressions"
)

// TODO ponerles mas informacion precisa a los errores
var (
	ErrFieldModelNotConcerned = errors.New("field's model is not concerned by the query so it can't be used in a expression")
	ErrJoinMustBeSelected     = errors.New("table is joined more than once, select which one you want to use")
)

const UndefinedJoinNumber = -1

// TODO doc
type ValueExpression[T any] struct {
	// TODO hacer el cambio de nombre en el anterior tambien?
	SQLExpressions []LiteralSQLExpression
	JoinNumber     int
}

type LiteralSQLExpression struct {
	SQL   string
	Field badorm.IFieldIdentifier
}

func (expr ValueExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

func (expr *ValueExpression[T]) SelectJoin(joinNumber uint) badorm.DynamicExpression[T] {
	expr.JoinNumber = int(joinNumber)
	return expr
}

// verificar que en las condiciones anteriores alguien usó el field con el que se intenta comparar
// obtener de ahi cual es el nombre de la table a usar con ese field.
// TODO doc a ingles
func (expr ValueExpression[T]) ToSQL(query *badorm.Query, columnName string) (string, []any, error) {
	exprString := columnName
	values := []any{}

	for _, sqlExpr := range expr.SQLExpressions {
		modelTable, err := getModelTable(query, sqlExpr.Field, expr.JoinNumber)
		if err != nil {
			return "", nil, err
		}

		exprString += fmt.Sprintf(
			" "+sqlExpr.SQL+" %s",
			sqlExpr.Field.ColumnSQL(query, modelTable),
		)
	}

	return exprString, values, nil
}

func newValueExpression[T any](expression expressions.SQLExpression, field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	// TODO que pasa con los que solo aceptan cierto tipo, ver las de like por ejemplo
	// TODO que pasa si hay mas de uno, no se si me gusta mucho esto
	// creo que si hay mas de uno solo se tendria que aplicar en el primero,
	// TODO soportar multivalue, no todos necesariamente dinamicos
	return &ValueExpression[T]{
		SQLExpressions: []LiteralSQLExpression{
			{
				SQL:   expressions.ToSQL[expression],
				Field: field,
			},
		},
		JoinNumber: UndefinedJoinNumber,
	}
}

// Expression that compares the value of the column against multiple values
// Example: value IN (v1, v2, v3, ..., vN)
type MultivalueExpression[T any] struct {
	// TODO hacer el cambio de nombre en el anterior tambien?
	// TODO con esto podria reemplazar el SQLExpressionAndValue para que todos sean por dentro dynamics
	Values        []any  // the values to compare with
	SQLExpression string // the expression used to compare, example: IN
	SQLConnector  string // the connector between values, example: ', '
	SQLPrefix     string // something to put before the values, example: (
	SQLSuffix     string // something to put after the values, example: )
	JoinNumber    int
}

func (expr MultivalueExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

func (expr *MultivalueExpression[T]) SelectJoin(joinNumber uint) badorm.DynamicExpression[T] {
	expr.JoinNumber = int(joinNumber)
	return expr
}

func (expr MultivalueExpression[T]) ToSQL(query *badorm.Query, columnName string) (string, []any, error) {
	placeholderList := []string{}
	values := []any{}

	for _, value := range expr.Values {
		field, isField := value.(badorm.IFieldIdentifier)
		if isField {
			// if it is a field, add the field column to the query
			modelTable, err := getModelTable(query, field, expr.JoinNumber)
			if err != nil {
				return "", nil, err
			}

			placeholderList = append(placeholderList, field.ColumnSQL(query, modelTable))
		} else {
			// if it is not a field, it a value, ass the placeholder ? and the value to the list
			placeholderList = append(placeholderList, "?")
			values = append(values, value)
		}
	}

	placeholders := strings.Join(placeholderList, " "+expr.SQLConnector+" ")

	return fmt.Sprintf(
		"%s %s %s"+placeholders+"%s",
		columnName,
		expr.SQLExpression,
		expr.SQLPrefix,
		expr.SQLSuffix,
	), values, nil
}

func newMultivalueExpression[T any](sqlExpression expressions.SQLExpression, sqlConnector, sqlPrefix, sqlSuffix string, fields ...badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	// TODO no todos necesariamente dinamicos, pueden ser T o field[T]
	values := pie.Map(fields, func(field badorm.FieldIdentifier[T]) any {
		return field
	})

	return &MultivalueExpression[T]{
		Values:        values,
		SQLExpression: expressions.ToSQL[sqlExpression],
		SQLConnector:  sqlConnector,
		SQLPrefix:     sqlPrefix,
		SQLSuffix:     sqlSuffix,
		JoinNumber:    UndefinedJoinNumber,
	}
}

func getModelTable(query *badorm.Query, field badorm.IFieldIdentifier, joinNumber int) (badorm.Table, error) {
	modelTables := query.GetTables(field.GetModelType())
	if modelTables == nil {
		return badorm.Table{}, ErrFieldModelNotConcerned
	}

	if len(modelTables) == 1 {
		return modelTables[0], nil
	}

	if joinNumber == UndefinedJoinNumber {
		return badorm.Table{}, ErrJoinMustBeSelected
	}

	return modelTables[joinNumber], nil
}

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
func Eq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.Eq, field)
}

// NotEqualTo
func NotEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.NotEq, field)
}

// LessThan
func Lt[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.Lt, field)
}

// LessThanOrEqualTo
func LtOrEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.LtOrEq, field)
}

// GreaterThan
func Gt[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	// TODO invertir orden de parametros para que quede igual que en badorm/expression
	return newValueExpression(expressions.Gt, field)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.GtOrEq, field)
}

// Comparison Predicates
// refs:
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to v1 < value < v2
func Between[T any](field1 badorm.FieldIdentifier[T], field2 badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newMultivalueExpression(expressions.Between, "AND", "", "", field1, field2)
}

// Equivalent to NOT (v1 < value < v2)
func NotBetween[T any](field1 badorm.FieldIdentifier[T], field2 badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newMultivalueExpression(expressions.NotBetween, "AND", "", "", field1, field2)
}

// Boolean Comparison Predicates

// Not supported by: mysql
func IsDistinct[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.IsDistinct, field)
}

// Not supported by: mysql
func IsNotDistinct[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.IsNotDistinct, field)
}
