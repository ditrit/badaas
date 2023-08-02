package orm

// Comparison Operators
// ref: https://www.postgresql.org/docs/current/functions-comparison.html

// EqualTo
func Eq[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](value, "=")
}

// NotEqualTo
func NotEq[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](value, "<>")
}
