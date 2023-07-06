package multitype

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
func Eq[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicOperator[TAttribute] {
	return newValueOperator[TAttribute, TField](sql.Eq, field)
}

// NotEqualTo
func NotEq[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicOperator[TAttribute] {
	return newValueOperator[TAttribute, TField](sql.NotEq, field)
}

// LessThan
func Lt[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicOperator[TAttribute] {
	return newValueOperator[TAttribute, TField](sql.Lt, field)
}

// LessThanOrEqualTo
func LtOrEq[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicOperator[TAttribute] {
	return newValueOperator[TAttribute, TField](sql.LtOrEq, field)
}

// GreaterThan
func Gt[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicOperator[TAttribute] {
	// TODO invertir orden de parametros para que quede igual que en badorm/expression
	return newValueOperator[TAttribute, TField](sql.Gt, field)
}

// GreaterThanOrEqualTo
func GtOrEq[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicOperator[TAttribute] {
	return newValueOperator[TAttribute, TField](sql.GtOrEq, field)
}

// Comparison Predicates
// refs:
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to v1 < value < v2
func Between[TAttribute, TField any](v1, v2 any) badorm.DynamicOperator[TAttribute] {
	return newMultivalueOperator[TAttribute, TField](sql.Between, "AND", "", "", v1, v2)
}

// Equivalent to NOT (v1 < value < v2)
func NotBetween[TAttribute, TField any](v1, v2 any) badorm.DynamicOperator[TAttribute] {
	return newMultivalueOperator[TAttribute, TField](sql.NotBetween, "AND", "", "", v1, v2)
}

// Boolean Comparison Predicates

// Not supported by: mysql
func IsDistinct[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicOperator[TAttribute] {
	return newValueOperator[TAttribute, TField](sql.IsDistinct, field)
}

// Not supported by: mysql
func IsNotDistinct[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicOperator[TAttribute] {
	return newValueOperator[TAttribute, TField](sql.IsNotDistinct, field)
}
