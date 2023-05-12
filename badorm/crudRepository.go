package badorm

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/badorm/pagination"
	"github.com/ditrit/badaas/configuration"
	"github.com/gertd/go-pluralize"
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
	Get(tx *gorm.DB, conditions map[string]any) (*T, error)
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

func (repository *CRUDRepositoryImpl[T, ID]) Get(tx *gorm.DB, conditions map[string]any) (*T, error) {
	entity, err := repository.GetOptional(tx, conditions)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, ErrObjectNotFound
	}

	return entity, nil
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

	var entity T
	// only entities that match the conditions
	for joinAttributeName, joinConditions := range joinConditions {
		schemaName, err := schema.Parse(entity, &sync.Map{}, tx.NamingStrategy)
		if err != nil {
			return nil, err
		}
		err = repository.addJoinToQuery(
			query,
			schemaName.Table,
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

func (repository *CRUDRepositoryImpl[T, ID]) GetAll(tx *gorm.DB) ([]*T, error) {
	return repository.GetMultiple(tx, map[string]any{})
}

// Adds a join to the "query" by the "attributeName" that must be relation type
// then, adds the verification that the values for the joined entity are "expectedValues"

// "expectedValues" is in {"attributeName": expectedValue} format
func (repository *CRUDRepositoryImpl[T, ID]) addJoinToQuery(
	query *gorm.DB, previousTable, joinAttributeName string, conditions map[string]any,
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

	pluralize := pluralize.NewClient()
	// TODO poder no hacer esto, en caso de que hayan difinido otra cosa
	tableName := pluralize.Plural(joinAttributeName)
	// TODO creo que deberia ser al revez
	tableWithSuffix := tableName + "_" + previousTable
	stringQuery := fmt.Sprintf(
		`JOIN %[1]s %[2]s ON
			%[2]s.id = %[3]s.%[4]s_id
			AND %[2]s.deleted_at IS NULL
		`,
		tableName,
		tableWithSuffix,
		previousTable,
		joinAttributeName,
		// TODO que pasa si el atributo no existe
		// TODO ver que pase si attributeName no existe como tabla
	)

	conditionsValues := []any{}
	for attributeName, conditionValue := range thisEntityConditions {
		stringQuery += fmt.Sprintf(
			`AND %[1]s.%[2]s = ?
			`,
			tableWithSuffix, attributeName,
		)
		conditionsValues = append(conditionsValues, conditionValue)
	}

	query.Joins(stringQuery, conditionsValues...)

	// only entities that match the conditions
	// TODO codigo duplicado
	for joinAttributeName, joinConditions := range joinConditions {
		err := repository.addJoinToQuery(query, tableWithSuffix, joinAttributeName, joinConditions)
		if err != nil {
			return err
		}
	}

	return nil
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
