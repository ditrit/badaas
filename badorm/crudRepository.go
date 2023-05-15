package badorm

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/badorm/pagination"
	"github.com/ditrit/badaas/configuration"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type BadaasID interface {
	int | uuid.UUID
}

// Generic CRUD Repository
type CRUDRepository[T any, ID BadaasID] interface {
	// create
	Create(tx *gorm.DB, entity *T) error
	// read
	GetByID(tx *gorm.DB, id ID) (*T, error)
	Get(tx *gorm.DB, conditions map[string]any) (T, error)
	GetOptional(tx *gorm.DB, conditions map[string]any) (*T, error)
	GetMultiple(tx *gorm.DB, conditions map[string]any) ([]*T, error)
	GetAll(tx *gorm.DB) ([]*T, error)
	Find(tx *gorm.DB, filters squirrel.Sqlizer, pagination pagination.Paginator, sort SortOption) (*pagination.Page[T], error)
	// update
	Save(tx *gorm.DB, entity *T) error
	// delete
	Delete(tx *gorm.DB, entity *T) error
}

var (
	ErrMoreThanOneObjectFound = errors.New("found more that one object that meet the requested conditions")
	ErrObjectNotFound         = errors.New("no object exists that meets the requested conditions")
	ErrObjectsNotRelated      = func(typeName, attributeName string) error {
		return fmt.Errorf("%[1]s has not attribute named %[2]s or %[2]sID", typeName, attributeName)
	}
	ErrModelNotRegistered = func(typeName, attributeName string) error {
		return fmt.Errorf("%[1]s has an attribute named %[2]s or %[2]sID but %[2]s is not registered as model (use AddModel)", typeName, attributeName)
	}
)

// Implementation of the Generic CRUD Repository
type CRUDRepositoryImpl[T any, ID BadaasID] struct {
	CRUDRepository[T, ID]
	logger                  *zap.Logger
	paginationConfiguration configuration.PaginationConfiguration
}

// Constructor of the Generic CRUD Repository
func NewCRUDRepository[T any, ID BadaasID](
	logger *zap.Logger,
	paginationConfiguration configuration.PaginationConfiguration,
) CRUDRepository[T, ID] {
	return &CRUDRepositoryImpl[T, ID]{
		logger:                  logger,
		paginationConfiguration: paginationConfiguration,
	}
}

// Create an entity of a Model
func (repository *CRUDRepositoryImpl[T, ID]) Create(tx *gorm.DB, entity *T) error {
	err := tx.Create(entity).Error
	if err != nil {
		return err
	}

	return nil
}

// Delete an entity of a Model
func (repository *CRUDRepositoryImpl[T, ID]) Delete(tx *gorm.DB, entity *T) error {
	err := tx.Delete(entity).Error
	if err != nil {
		return err
	}

	return nil
}

// Save an entity of a Model
func (repository *CRUDRepositoryImpl[T, ID]) Save(tx *gorm.DB, entity *T) error {
	err := tx.Save(entity).Error
	if err != nil {
		return err
	}

	return nil
}

// Get an entity of a Model By ID
func (repository *CRUDRepositoryImpl[T, ID]) GetByID(tx *gorm.DB, id ID) (*T, error) {
	var entity T
	err := tx.First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (repository *CRUDRepositoryImpl[T, ID]) Get(tx *gorm.DB, conditions map[string]any) (T, error) {
	entity, err := repository.GetOptional(tx, conditions)
	var nilValue T
	if err != nil {
		return nilValue, err
	}

	if entity == nil {
		return nilValue, ErrObjectNotFound
	}

	return *entity, nil
}

func (repository *CRUDRepositoryImpl[T, ID]) GetOptional(tx *gorm.DB, conditions map[string]any) (*T, error) {
	entities, err := repository.GetMultiple(tx, conditions)
	if err != nil {
		return nil, err
	}

	if len(entities) > 1 {
		return nil, ErrMoreThanOneObjectFound
	} else if len(entities) == 1 {
		return entities[0], nil
	}

	return nil, nil
}

// Get all entities of a Model
func (repository *CRUDRepositoryImpl[T, ID]) GetMultiple(tx *gorm.DB, conditions map[string]any) ([]*T, error) {
	thisEntityConditions := map[string]any{}
	joinConditions := map[string]map[string]any{}

	// only entities that match the conditions
	for attributeName, expectedValue := range conditions {
		switch typedExpectedValue := expectedValue.(type) {
		case float64, bool, string, nil:
			thisEntityConditions[attributeName] = expectedValue
		case map[string]any:
			joinConditions[attributeName] = typedExpectedValue
		default:
			return nil, fmt.Errorf("unsupported type")
		}
	}

	query := tx.Where(thisEntityConditions)

	entity := new(T)
	// only entities that match the conditions
	for joinAttributeName, joinConditions := range joinConditions {
		tableName, err := getTableName(tx, entity)
		if err != nil {
			return nil, err
		}

		err = repository.addJoinToQuery(
			query,
			entity,
			tableName,
			joinAttributeName,
			joinConditions,
		)
		if err != nil {
			return nil, err
		}
	}

	// execute query
	var entities []*T
	err := query.Find(&entities).Error

	return entities, err
}

func getTableName(db *gorm.DB, entity any) (string, error) {
	schemaName, err := schema.Parse(entity, &sync.Map{}, db.NamingStrategy)
	if err != nil {
		return "", err
	}

	return schemaName.Table, nil
}

func (repository *CRUDRepositoryImpl[T, ID]) GetAll(tx *gorm.DB) ([]*T, error) {
	return repository.GetMultiple(tx, map[string]any{})
}

// Adds a join to the "query" by the "attributeName" that must be relation type
// then, adds the verification that the values for the joined entity are "expectedValues"

// "expectedValues" is in {"attributeName": expectedValue} format
// TODO support ManyToMany relations
// previousEntity is pointer
func (repository *CRUDRepositoryImpl[T, ID]) addJoinToQuery(
	query *gorm.DB, previousEntity any,
	previousTableName, joinAttributeName string,
	conditions map[string]any,
) error {
	// TODO codigo duplicado
	thisEntityConditions := map[string]any{}
	joinConditions := map[string]map[string]any{}

	for attributeName, expectedValue := range conditions {
		switch typedExpectedValue := expectedValue.(type) {
		case float64, bool, string, nil:
			thisEntityConditions[attributeName] = expectedValue
		case map[string]any:
			joinConditions[attributeName] = typedExpectedValue
		default:
			return fmt.Errorf("unsupported type")
		}
	}

	relatedObject, relationIDIsInPreviousTable, err := getRelatedObject(
		previousEntity,
		joinAttributeName,
	)
	if err != nil {
		return err
	}

	joinTableName, err := getTableName(query, relatedObject)
	if err != nil {
		return err
	}

	tableWithSuffix := joinTableName + "_" + previousTableName

	var stringQuery string
	if relationIDIsInPreviousTable {
		stringQuery = fmt.Sprintf(
			`JOIN %[1]s %[2]s ON
				%[2]s.id = %[3]s.%[4]s_id
				AND %[2]s.deleted_at IS NULL
			`,
			joinTableName,
			tableWithSuffix,
			previousTableName,
			joinAttributeName,
		)
	} else {
		// TODO foreignKey can be redefined (https://gorm.io/docs/has_one.html#Override-References)
		previousAttribute := reflect.TypeOf(previousEntity).Elem().Name()
		stringQuery = fmt.Sprintf(
			`JOIN %[1]s %[2]s ON
				%[2]s.%[4]s_id = %[3]s.id
				AND %[2]s.deleted_at IS NULL
			`,
			joinTableName,
			tableWithSuffix,
			previousTableName,
			previousAttribute,
		)
	}

	conditionsValues := []any{}
	for attributeName, conditionValue := range thisEntityConditions {
		stringQuery += fmt.Sprintf(
			`AND %[1]s.%[2]s = ?
			`,
			tableWithSuffix, attributeName,
			// TODO que pasa si el atributo no existe
		)
		conditionsValues = append(conditionsValues, conditionValue)
	}

	query.Joins(stringQuery, conditionsValues...)

	// TODO codigo duplicado
	for joinAttributeName, joinConditions := range joinConditions {
		err := repository.addJoinToQuery(
			query,
			relatedObject,
			tableWithSuffix,
			joinAttributeName,
			joinConditions,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// entity can be a pointer of not, now only works with pointer
func getRelatedObject(entity any, relationName string) (any, bool, error) {
	entityType := reflect.TypeOf(entity)

	// entityType will be a pointer if the relation can be nullable
	if entityType.Kind() == reflect.Pointer {
		entityType = entityType.Elem()
	}

	field, isPresent := entityType.FieldByName(relationName)
	if !isPresent {
		// some gorm relations dont have a direct relation in the model, only the id
		return getRelatedObjectByID(entityType, relationName)
	}

	_, isIDPresent := entityType.FieldByName(relationName + "ID")

	return createObject(field.Type), isIDPresent, nil
}

func getRelatedObjectByID(entityType reflect.Type, relationName string) (any, bool, error) {
	_, isPresent := entityType.FieldByName(relationName + "ID")
	if !isPresent {
		return nil, false, ErrObjectsNotRelated(entityType.Name(), relationName)
	}

	fieldType, isPresent := modelsMapping[relationName]
	if !isPresent {
		return nil, false, ErrModelNotRegistered(entityType.Name(), relationName)
	}

	return createObject(fieldType), true, nil
}

func createObject(entityType reflect.Type) any {
	return reflect.New(entityType).Elem().Interface()
}

// Find entities of a Model
func (repository *CRUDRepositoryImpl[T, ID]) Find(
	tx *gorm.DB,
	filters squirrel.Sqlizer,
	page pagination.Paginator,
	sortOption SortOption,
) (*pagination.Page[T], error) {
	var instances []*T
	whereClause, values, err := filters.ToSql()
	if err != nil {
		return nil, err
	}

	if page != nil {
		tx = tx.
			Offset(
				int((page.Offset() - 1) * page.Limit()),
			).
			Limit(
				int(page.Limit()),
			)
	} else {
		page = pagination.NewPaginator(0, repository.paginationConfiguration.GetMaxElemPerPage())
	}

	if sortOption != nil {
		tx = tx.Order(buildClauseFromSortOption(sortOption))
	}

	err = tx.Where(whereClause, values...).Find(&instances).Error
	if err != nil {
		return nil, err
	}

	// Get Count
	nbElem, err := repository.count(tx, whereClause, values)
	if err != nil {
		return nil, err
	}

	return pagination.NewPage(instances, page.Offset(), page.Limit(), nbElem), nil
}

// Count the number of record that match the where clause with the provided values on the db
func (repository *CRUDRepositoryImpl[T, ID]) count(tx *gorm.DB, whereClause string, values []interface{}) (uint, error) {
	var entity *T
	var count int64
	err := tx.Model(entity).Where(whereClause, values).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return uint(count), nil
}

// Build a gorm order clause from a SortOption
func buildClauseFromSortOption(sortOption SortOption) clause.OrderByColumn {
	return clause.OrderByColumn{Column: clause.Column{Name: sortOption.Column()}, Desc: sortOption.Desc()}
}
