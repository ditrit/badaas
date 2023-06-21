package badorm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
)

const DeletedAtField = "DeletedAt"

var (
	IDColumnID        = ColumnIdentifier{Field: "ID"}
	CreatedAtColumnID = ColumnIdentifier{Field: "CreatedAt"}
	UpdatedAtColumnID = ColumnIdentifier{Field: "UpdatedAt"}
	DeletedAtColumnID = ColumnIdentifier{Field: DeletedAtField}
)

var ErrEmptyConditions = errors.New("condition must have at least one inner condition")

type Condition[T any] interface {
	// Applies the condition to the "query"
	// using the "tableName" as name for the table holding
	// the data for object of type T
	ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error)

	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T],
	// since if no method receives by parameter a type T,
	// any other Condition[T2] would also be considered a Condition[T].
	interfaceVerificationMethod(T)
}

// Conditions that can be used in a where clause
// (or in a on of a join)
type WhereCondition[T any] interface {
	Condition[T]

	// Get the sql string and values to use in the query
	GetSQL(query *gorm.DB, tableName string) (string, []any, error)

	// Returns true if the DeletedAt column if affected by the condition
	// If no condition affects the DeletedAt, the verification that it's null will be added automatically
	affectsDeletedAt() bool
}

// Condition that contains a internal condition.
// Example: NOT (internal condition)
type ContainerCondition[T any] struct {
	ConnectionCondition WhereCondition[T]
	Prefix              string
}

//nolint:unused // see inside
func (condition ContainerCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition ContainerCondition[T]) ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error) {
	return applyWhereCondition[T](condition, query, tableName)
}

func (condition ContainerCondition[T]) GetSQL(query *gorm.DB, tableName string) (string, []any, error) {
	sqlString, values, err := condition.ConnectionCondition.GetSQL(query, tableName)
	if err != nil {
		return "", nil, err
	}

	sqlString = condition.Prefix + " (" + sqlString + ")"

	return sqlString, values, nil
}

//nolint:unused // is used
func (condition ContainerCondition[T]) affectsDeletedAt() bool {
	return condition.ConnectionCondition.affectsDeletedAt()
}

// Condition that contains a internal condition.
// Example: NOT (internal condition)
func NewContainerCondition[T any](prefix string, conditions ...WhereCondition[T]) WhereCondition[T] {
	if len(conditions) == 0 {
		return NewInvalidCondition[T](ErrEmptyConditions)
	}

	return ContainerCondition[T]{
		Prefix:              prefix,
		ConnectionCondition: And(conditions...),
	}
}

// Condition that connects multiple conditions.
// Example: condition1 AND condition2
type ConnectionCondition[T any] struct {
	Connector  string
	Conditions []WhereCondition[T]
}

//nolint:unused // see inside
func (condition ConnectionCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition ConnectionCondition[T]) ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error) {
	return applyWhereCondition[T](condition, query, tableName)
}

func (condition ConnectionCondition[T]) GetSQL(query *gorm.DB, tableName string) (string, []any, error) {
	sqlStrings := []string{}
	values := []any{}

	for _, internalCondition := range condition.Conditions {
		internalSQLString, exprValues, err := internalCondition.GetSQL(query, tableName)
		if err != nil {
			return "", nil, err
		}

		sqlStrings = append(sqlStrings, internalSQLString)

		values = append(values, exprValues...)
	}

	return strings.Join(sqlStrings, " "+condition.Connector+" "), values, nil
}

//nolint:unused // is used
func (condition ConnectionCondition[T]) affectsDeletedAt() bool {
	return pie.Any(condition.Conditions, func(internalCondition WhereCondition[T]) bool {
		return internalCondition.affectsDeletedAt()
	})
}

// Condition that connects multiple conditions.
// Example: condition1 AND condition2
func NewConnectionCondition[T any](connector string, conditions ...WhereCondition[T]) WhereCondition[T] {
	return ConnectionCondition[T]{
		Connector:  connector,
		Conditions: conditions,
	}
}

// TODO usar tambien en las conditions
// poner en variables y reutilizar
type ColumnIdentifier struct {
	Column       string
	Field        string
	ColumnPrefix string
}

func (columnID ColumnIdentifier) ColumnName(db *gorm.DB, tableName string) string {
	// TODO codigo repetido
	columnName := columnID.Column
	if columnName == "" {
		columnName = db.NamingStrategy.ColumnName(tableName, columnID.Field)
	}

	// add column prefix and table name once we know the column name
	return columnID.ColumnPrefix + columnName
}

// TODO doc
type PreloadCondition[T any] struct {
	Columns []ColumnIdentifier
}

//nolint:unused // see inside
func (condition PreloadCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition PreloadCondition[T]) ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error) {
	for _, columnID := range condition.Columns {
		columnName := columnID.ColumnName(query, tableName)

		// Remove first table name as GORM only adds it from the second join
		_, attributePrefix, _ := strings.Cut(tableName, "__")

		query.Statement.Selects = append(
			query.Statement.Selects,
			fmt.Sprintf(
				"%[1]s.%[2]s AS \"%[3]s__%[2]s\"", // name used by gorm to load the fields inside the models
				tableName,
				columnName,
				attributePrefix,
			),
		)
	}

	return query, nil
}

// TODO doc
func NewPreloadCondition[T any](columns ...ColumnIdentifier) PreloadCondition[T] {
	return PreloadCondition[T]{
		Columns: append(
			columns,
			// base model fields
			IDColumnID,
			CreatedAtColumnID,
			UpdatedAtColumnID,
			DeletedAtColumnID,
		),
	}
}

// Condition that verifies the value of a field,
// using the Expression
type FieldCondition[TObject any, TAtribute any] struct {
	Field        string
	Column       string
	ColumnPrefix string
	Expression   Expression[TAtribute]
}

//nolint:unused // see inside
func (condition FieldCondition[TObject, TAtribute]) interfaceVerificationMethod(_ TObject) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Returns a gorm Where condition that can be used
// to filter that the Field as a value of Value
func (condition FieldCondition[TObject, TAtribute]) ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error) {
	return applyWhereCondition[TObject](condition, query, tableName)
}

func applyWhereCondition[T any](condition WhereCondition[T], query *gorm.DB, tableName string) (*gorm.DB, error) {
	sql, values, err := condition.GetSQL(query, tableName)
	if err != nil {
		return nil, err
	}

	if condition.affectsDeletedAt() {
		query = query.Unscoped()
	}

	return query.Where(
		sql,
		values...,
	), nil
}

//nolint:unused // is used
func (condition FieldCondition[TObject, TAtribute]) affectsDeletedAt() bool {
	return condition.Field == DeletedAtField
}

func (condition FieldCondition[TObject, TAtribute]) GetSQL(query *gorm.DB, tableName string) (string, []any, error) {
	// TODO codigo repetido
	columnName := condition.Column
	if columnName == "" {
		columnName = query.NamingStrategy.ColumnName(tableName, condition.Field)
	}

	// add column prefix and table name once we know the column name
	columnName = tableName + "." + condition.ColumnPrefix + columnName

	return condition.Expression.ToSQL(columnName)
}

// Condition that joins with other table
type JoinCondition[T1 any, T2 any] struct {
	T1Field         string
	T2Field         string
	ConnectionField string
	Conditions      []Condition[T2]
}

//nolint:unused // see inside
func (condition JoinCondition[T1, T2]) interfaceVerificationMethod(_ T1) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Applies a join between the tables of T1 and T2
// previousTableName is the name of the table of T1
// It also applies the nested conditions
func (condition JoinCondition[T1, T2]) ApplyTo(query *gorm.DB, previousTableName string) (*gorm.DB, error) {
	// get the name of the table for T2
	toBeJoinedTableName, err := getTableName(query, *new(T2))
	if err != nil {
		return nil, err
	}

	// add a suffix to avoid tables with the same name when joining
	// the same table more than once
	nextTableAlias := previousTableName + "__" + condition.ConnectionField

	whereConditions, joinConditions, preloadCondition := divideConditionsByType(condition.Conditions)

	// apply WhereConditions to join in "on" clause
	connectionCondition := And(whereConditions...)

	onQuery, onValues, err := connectionCondition.GetSQL(query, nextTableAlias)
	if err != nil {
		return nil, err
	}

	// get the sql to do the join with T2
	// if it's only a preload use a left join
	// TODO una condicion para ver que la relacion sea null (ademas de hacerle is null al fk)
	// TODO no me termina de convencer que para el preload hay que hacer el join si o si
	isLeftJoin := len(whereConditions) == 0 && preloadCondition != nil
	joinQuery := condition.getSQLJoin(
		query,
		toBeJoinedTableName,
		nextTableAlias,
		previousTableName,
		isLeftJoin,
	)

	if onQuery != "" {
		joinQuery += " AND " + onQuery
	}

	if !connectionCondition.affectsDeletedAt() {
		joinQuery += fmt.Sprintf(
			" AND %s.deleted_at IS NULL",
			nextTableAlias,
		)
	}

	// add the join to the query
	query = query.Joins(joinQuery, onValues...)

	// apply preload condition
	if preloadCondition != nil {
		query, err = preloadCondition.ApplyTo(query, nextTableAlias)
		if err != nil {
			return nil, err
		}
	}

	// apply nested joins
	for _, joinCondition := range joinConditions {
		query, err = joinCondition.ApplyTo(query, nextTableAlias)
		if err != nil {
			return nil, err
		}
	}

	return query, nil
}

// Returns the SQL string to do a join between T1 and T2
// taking into account that the ID attribute necessary to do it
// can be either in T1's or T2's table.
func (condition JoinCondition[T1, T2]) getSQLJoin(
	query *gorm.DB,
	toBeJoinedTableName, nextTableAlias, previousTableName string,
	isLeftJoin bool,
) string {
	joinString := "INNER JOIN"
	if isLeftJoin {
		joinString = "LEFT JOIN"
	}

	return fmt.Sprintf(
		`%[6]s %[1]s %[2]s ON %[2]s.%[3]s = %[4]s.%[5]s
		`,
		toBeJoinedTableName,
		nextTableAlias,
		query.NamingStrategy.ColumnName(nextTableAlias, condition.T2Field),
		previousTableName,
		query.NamingStrategy.ColumnName(previousTableName, condition.T1Field),
		joinString,
	)
}

// Divides a list of conditions by its type: WhereConditions and JoinConditions
func divideConditionsByType[T any](
	conditions []Condition[T],
) (whereConditions []WhereCondition[T], joinConditions []Condition[T], preloadCondition *PreloadCondition[T]) {
	for _, condition := range conditions {
		whereCondition, ok := condition.(WhereCondition[T])
		if ok {
			whereConditions = append(whereConditions, whereCondition)
		} else {
			possiblePreloadCondition, ok := condition.(PreloadCondition[T])
			if ok {
				preloadCondition = &possiblePreloadCondition
			} else {
				joinConditions = append(joinConditions, condition)
			}
		}
	}

	return
}

// Condition that can be used to express conditions that are not supported (yet?) by BaDORM
// Example: table1.columnX = table2.columnY
type UnsafeCondition[T any] struct {
	SQLCondition string
	Values       []any
}

//nolint:unused // see inside
func (condition UnsafeCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition UnsafeCondition[T]) ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error) {
	return applyWhereCondition[T](condition, query, tableName)
}

func (condition UnsafeCondition[T]) GetSQL(_ *gorm.DB, tableName string) (string, []any, error) {
	return fmt.Sprintf(
		condition.SQLCondition,
		tableName,
	), condition.Values, nil
}

//nolint:unused // is used
func (condition UnsafeCondition[T]) affectsDeletedAt() bool {
	return false
}

// Condition that can be used to express conditions that are not supported (yet?) by BaDORM
// Example: table1.columnX = table2.columnY
func NewUnsafeCondition[T any](condition string, values []any) UnsafeCondition[T] {
	return UnsafeCondition[T]{
		SQLCondition: condition,
		Values:       values,
	}
}

// Condition used to returns an error when the query is executed
type InvalidCondition[T any] struct {
	Err error
}

//nolint:unused // see inside
func (condition InvalidCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition InvalidCondition[T]) ApplyTo(_ *gorm.DB, _ string) (*gorm.DB, error) {
	return nil, condition.Err
}

func (condition InvalidCondition[T]) GetSQL(_ *gorm.DB, _ string) (string, []any, error) {
	return "", nil, condition.Err
}

//nolint:unused // is used
func (condition InvalidCondition[T]) affectsDeletedAt() bool {
	return false
}

// Condition used to returns an error when the query is executed
func NewInvalidCondition[T any](err error) InvalidCondition[T] {
	return InvalidCondition[T]{
		Err: err,
	}
}

// Logical Operators
// ref:
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-logical.html
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/logical-operators.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/logical-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

func And[T any](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition("AND", conditions...)
}

func Or[T any](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition("OR", conditions...)
}

func Not[T any](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewContainerCondition("NOT", conditions...)
}
