package services

import (
	"fmt"

	"github.com/ditrit/badaas/persistence/models"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

type CRUDService[T models.Tabler, ID any] interface {
	GetEntity(id ID) (*T, error)
	GetEntities(conditions map[string]any) ([]*T, error)
	CreateEntity(attributeValues map[string]any) (*T, error)
	UpdateEntity(entityID ID, newValues map[string]any) (*T, error)
	DeleteEntity(entityID ID) error
}

// check interface compliance
var _ CRUDService[models.User, uuid.UUID] = (*crudServiceImpl[models.User, uuid.UUID])(nil)

// Implementation of the Generic CRUD Repository
type crudServiceImpl[T models.Tabler, ID any] struct {
	CRUDService[T, ID]
	logger *zap.Logger
	db     *gorm.DB
}

func NewCRUDService[T models.Tabler](
	logger *zap.Logger,
	db *gorm.DB,
) CRUDService[T, uuid.UUID] { // TODO ver este UUID hardcodeado aca
	return &crudServiceImpl[T, uuid.UUID]{
		logger: logger,
		db:     db,
	}
}

// TODO todo el codigo duplicado con CRUD repository
// Get the Entity of type with name "entityTypeName" that has the "id"
func (service *crudServiceImpl[T, ID]) GetEntity(id uuid.UUID) (*T, error) {
	var entity T
	err := service.db.First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// Get entities of type with name "entityTypeName" that match all "conditions"
//
// "conditions" is in {"attributeName": expectedValue} format
// in case of join "conditions" can have the format:
//
//	{"relationAttributeName": {"attributeName": expectedValue}}
func (service *crudServiceImpl[T, ID]) GetEntities(conditions map[string]any) ([]*T, error) {
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

	query := service.db.Where(thisEntityConditions)

	var entity T
	// only entities that match the conditions
	for joinAttributeName, joinConditions := range joinConditions {
		// TODO +s no necesariamente, hay plurales que cambian
		err := service.addJoinToQuery(query, entity.TableName(), joinAttributeName, joinConditions)
		if err != nil {
			return nil, err
		}
	}

	// execute query
	var entities []*T
	err := query.Find(&entities).Error

	return entities, err
}

// Adds a join to the "query" by the "attributeName" that must be relation type
// then, adds the verification that the values for the joined entity are "expectedValues"

// "expectedValues" is in {"attributeName": expectedValue} format
func (service *crudServiceImpl[T, ID]) addJoinToQuery(
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
	tableName := pluralize.Plural(joinAttributeName)
	tableWithSuffix := tableName + "_" + previousTable
	stringQuery := fmt.Sprintf(
		`JOIN %[1]s %[2]s ON
			%[2]s.id = %[3]s.%[4]s_id
		`,
		tableName,
		tableWithSuffix,
		previousTable,
		joinAttributeName,
		// TODO ver que pase si attributeName no existe como tabla
	)

	for _, attributeName := range maps.Keys(thisEntityConditions) {
		stringQuery += fmt.Sprintf(
			`AND %[1]s.%[2]s = ?
			`,
			tableWithSuffix, attributeName,
		)
	}

	query.Joins(stringQuery, maps.Values(thisEntityConditions)...)

	// only entities that match the conditions
	// TODO codigo duplicado
	for joinAttributeName, joinConditions := range joinConditions {
		err := service.addJoinToQuery(query, tableWithSuffix, joinAttributeName, joinConditions)
		if err != nil {
			return err
		}
	}

	return nil
}

// Creates a Entity of type "entityType" and its Values from the information provided in "attributeValues"
// not specified values are set with the default value (if exists) or a null value.
// entries in "attributeValues" that do not correspond to any attribute of "entityType" are ignored
//
// "attributeValues" are in {"attributeName": value} format
func (service *crudServiceImpl[T, ID]) CreateEntity(attributeValues map[string]any) (*T, error) {
	var entity T
	// TODO ver si dejo esto o desencodeo el json directo en la entidad
	// TODO testear lo de que se le pueden agregar relaciones a esto
	err := mapstructure.Decode(attributeValues, &entity)
	if err != nil {
		return nil, err // TODO ver que errores puede haber aca
	}

	err = service.db.Create(&entity).Error
	if err != nil {
		return nil, err
	}

	// TODO eliminar esto
	// err := service.db.Model(&entity).Create(attributeValues).Error
	// if err != nil {
	// 	return nil, err
	// }
	// entity.ID = attributeValues["id"]

	return &entity, nil
}

// Updates entity with type "entityTypeName" and "id" Values to the new values provided in "newValues"
// entries in "newValues" that do not correspond to any attribute of the EntityType are ignored
//
// "newValues" are in {"attributeName": newValue} format
func (service *crudServiceImpl[T, ID]) UpdateEntity(entityID uuid.UUID, newValues map[string]any) (*T, error) {
	return nil, nil
}

// Deletes Entity with type "entityTypeName" and id "entityID" and its values
func (service *crudServiceImpl[T, ID]) DeleteEntity(entityID uuid.UUID) error {
	return nil
}
