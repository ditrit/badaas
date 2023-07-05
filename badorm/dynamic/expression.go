package dynamic

import (
	"errors"
	"fmt"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/expressions"
)

// TODO ponerles mas informacion precisa a los errores
var (
	ErrFieldModelNotConcerned = errors.New("field's model is not concerned by the query so it can't be used in a expression")
	ErrJoinMustBeSelected     = errors.New("table is joined more than once, select which one you want to use")
)

const UndefinedJoinNumber = -1

func newDynamicExpression[T any](expression expressions.SQLExpression, field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	// TODO que pasa con los que solo aceptan cierto tipo, ver las de like por ejemplo
	// TODO que pasa si hay mas de uno, no se si me gusta mucho esto
	// creo que si hay mas de uno solo se tendria que aplicar en el primero,
	return &Expression[T]{
		SQLExpressions: []LiteralSQLExpression{
			{
				SQL:   expressions.ToSQL[expression],
				Field: field,
			},
		},
		JoinNumber: UndefinedJoinNumber,
	}
	// TODO soportar multivalue, no todos necesariamente dinamicos
}

// TODO doc
type Expression[T any] struct {
	// TODO hacer el cambio de nombre en el anterior tambien?
	SQLExpressions []LiteralSQLExpression
	JoinNumber     int
}

type LiteralSQLExpression struct {
	SQL   string
	Field badorm.IFieldIdentifier
}

func (expr Expression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

func (expr *Expression[T]) SelectJoin(joinNumber uint) badorm.DynamicExpression[T] {
	expr.JoinNumber = int(joinNumber)
	return expr
}

// verificar que en las condiciones anteriores alguien usó el field con el que se intenta comparar
// obtener de ahi cual es el nombre de la table a usar con ese field.
// TODO doc a ingles
func (expr Expression[T]) ToSQL(query *badorm.Query, columnName string) (string, []any, error) {
	exprString := columnName
	values := []any{}

	for _, sqlExpr := range expr.SQLExpressions {
		modelTables := query.GetTables(sqlExpr.Field.GetModelType())
		if modelTables == nil {
			return "", nil, ErrFieldModelNotConcerned
		}

		var modelTable badorm.Table

		if len(modelTables) == 1 {
			modelTable = modelTables[0]
		} else {
			if expr.JoinNumber == UndefinedJoinNumber {
				return "", nil, ErrJoinMustBeSelected
			}

			modelTable = modelTables[expr.JoinNumber]
		}

		exprString += fmt.Sprintf(
			" "+sqlExpr.SQL+" %s",
			sqlExpr.Field.ColumnSQL(query, modelTable),
		)
	}

	return exprString, values, nil
}

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
func Eq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newDynamicExpression(expressions.Eq, field)
}

// NotEqualTo
func NotEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newDynamicExpression(expressions.NotEq, field)
}

// LessThan
func Lt[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newDynamicExpression(expressions.Lt, field)
}

// LessThanOrEqualTo
func LtOrEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newDynamicExpression(expressions.LtOrEq, field)
}

// GreaterThan
func Gt[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	// TODO invertir orden de parametros para que quede igual que en badorm/expression
	return newDynamicExpression(expressions.Gt, field)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newDynamicExpression(expressions.GtOrEq, field)
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
func IsDistinct[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newDynamicExpression(expressions.IsDistinct, field)
}

// Not supported by: mysql
func IsNotDistinct[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newDynamicExpression(expressions.IsNotDistinct, field)
}
