package badorm

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// TODO creo que estos podrian ser todos tipos privados
type Table struct {
	Name    string
	Alias   string
	Initial bool
}

// Returns true if the Table is the initial table in a query
func (t Table) IsInitial() bool {
	return t.Initial
}

// Returns the related Table corresponding to the model
func (t Table) DeliverTable(query *query, model Model, relationName string) (Table, error) {
	// get the name of the table for the model
	tableName, err := getTableName(query.gormDB, model)
	if err != nil {
		return Table{}, err
	}

	// add a suffix to avoid tables with the same name when joining
	// the same table more than once
	tableAlias := relationName
	if !t.IsInitial() {
		tableAlias = t.Alias + "__" + relationName
	}

	return Table{
		Name:    tableName,
		Alias:   tableAlias,
		Initial: false,
	}, nil
}

type FieldIdentifier struct {
	Column       string
	Field        string
	ColumnPrefix string
	Type         reflect.Type
	ModelType    reflect.Type
}

// Returns the name of the column in which the field is saved in the table
func (columnID FieldIdentifier) ColumnName(query *query, table Table) string {
	columnName := columnID.Column
	if columnName == "" {
		columnName = query.ColumnName(table, columnID.Field)
	}

	// add column prefix and table name once we know the column name
	return columnID.ColumnPrefix + columnName
}

// Returns the SQL to get the value of the field in the table
func (columnID FieldIdentifier) ColumnSQL(query *query, table Table) string {
	return table.Alias + "." + columnID.ColumnName(query, table)
}

type query struct {
	gormDB          *gorm.DB
	concernedModels map[reflect.Type]Table
}

func (query *query) AddSelect(table Table, fieldID FieldIdentifier) {
	columnName := fieldID.ColumnName(query, table)

	query.gormDB.Statement.Selects = append(
		query.gormDB.Statement.Selects,
		fmt.Sprintf(
			"%[1]s.%[2]s AS \"%[1]s__%[2]s\"", // name used by gorm to load the fields inside the models
			table.Alias,
			columnName,
		),
	)
}

func (query *query) Preload(preloadQuery string, args ...interface{}) {
	query.gormDB = query.gormDB.Preload(preloadQuery, args...)
}

func (query *query) Unscoped() {
	query.gormDB = query.gormDB.Unscoped()
}

func (query *query) Where(whereQuery interface{}, args ...interface{}) {
	query.gormDB = query.gormDB.Where(whereQuery, args...)
}

func (query *query) Joins(joinQuery string, args ...interface{}) {
	query.gormDB = query.gormDB.Joins(joinQuery, args...)
}

func (query *query) Find(dest interface{}, conds ...interface{}) error {
	query.gormDB = query.gormDB.Find(dest, conds...)

	return query.gormDB.Error
}

func (query *query) AddConcernedModel(model Model, table Table) {
	// TODO que pasa si ya estaba
	query.concernedModels[reflect.TypeOf(model)] = table
}

// TODO ver esta, porque no estoy usando los fields aca y que pasa si hay fk override y todo eso
func (query query) ColumnName(table Table, fieldName string) string {
	return query.gormDB.NamingStrategy.ColumnName(table.Name, fieldName)
}

func NewQuery[T Model](db *gorm.DB, conditions []Condition[T]) (*query, error) {
	model := *new(T)

	initialTableName, err := getTableName(db, model)
	if err != nil {
		return nil, err
	}

	initialTable := Table{
		Name:    initialTableName,
		Alias:   initialTableName,
		Initial: true,
	}

	query := &query{
		gormDB:          db.Select(initialTableName + ".*"),
		concernedModels: map[reflect.Type]Table{},
	}
	query.AddConcernedModel(model, initialTable)

	for _, condition := range conditions {
		err = condition.ApplyTo(query, initialTable)
		if err != nil {
			return nil, err
		}
	}

	return query, nil
}
