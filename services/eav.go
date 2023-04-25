package services

import (
	"errors"
	"fmt"

	uuid "github.com/google/uuid"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrCantParseUUID = errors.New("can't parse uuid")
)

// EAV service provide handle EAV objects
type EAVService interface {
	GetEntity(entityTypeName string, id uuid.UUID) (*models.Entity, error)
	GetEntities(entityTypeName string, conditions map[string]string) ([]*models.Entity, error)
	CreateEntity(entityTypeName string, attributeValues map[string]any) (*models.Entity, error)
	UpdateEntity(entityTypeName string, entityID uuid.UUID, newValues map[string]any) (*models.Entity, error)
	DeleteEntity(entityTypeName string, entityID uuid.UUID) error
}

type eavServiceImpl struct {
	logger               *zap.Logger
	db                   *gorm.DB
	entityTypeRepository *repository.EntityTypeRepository
	entityRepository     *repository.EntityRepository
}

func NewEAVService(
	logger *zap.Logger,
	db *gorm.DB,
	entityTypeRepository *repository.EntityTypeRepository,
	entityRepository *repository.EntityRepository,
) EAVService {
	return &eavServiceImpl{
		logger:               logger,
		db:                   db,
		entityTypeRepository: entityTypeRepository,
		entityRepository:     entityRepository,
	}
}

// Get the Entity of type with name "entityTypeName" that has the "id"
func (eavService *eavServiceImpl) GetEntity(entityTypeName string, id uuid.UUID) (*models.Entity, error) {
	return eavService.entityRepository.Get(eavService.db, entityTypeName, id)
}

// Get entities of type with name "entityTypeName" that match all "conditions"
// entries in "conditions" that do not correspond to any attribute of "entityTypeName" are ignored
//
// "conditions" are in {"attributeName": expectedValue} format
// TODO relations join
func (eavService *eavServiceImpl) GetEntities(entityTypeName string, conditions map[string]string) ([]*models.Entity, error) {
	return ExecWithTransaction(
		eavService.db,
		func(tx *gorm.DB) ([]*models.Entity, error) {
			entityType, err := eavService.entityTypeRepository.GetByName(tx, entityTypeName)
			if err != nil {
				return nil, err
			}

			query := tx.Select("entities.*")

			// only entities that match the conditions
			for _, attribute := range entityType.Attributes {
				expectedValue, isPresent := conditions[attribute.Name]
				if isPresent {
					query = eavService.addValueCheckToQuery(query, attribute, expectedValue)
				}
			}

			// only entities of type entityType
			query = query.Where(
				"entities.entity_type_id = ?",
				entityType.ID,
			)

			// execute query
			var entities []*models.Entity
			query = query.Preload("Fields").Preload("Fields.Attribute").Preload("EntityType.Attributes").Preload("EntityType")
			err = query.Find(&entities).Error

			return entities, err
		},
	)
}

// Adds to the "query" the verification that the value for "attribute" is "expectedValue"
func (eavService *eavServiceImpl) addValueCheckToQuery(query *gorm.DB, attribute *models.Attribute, expectedValue string) *gorm.DB {
	var valToUseInQuery = string(attribute.ValueType) + "_val"
	if expectedValue == "null" { // TODO should be changed to be able to use nil
		valToUseInQuery = "is_null"
		expectedValue = "true"
	}

	return query.Joins(
		fmt.Sprintf(
			`JOIN attributes ON
				attributes.entity_type_id = entities.entity_type_id
				AND attributes.name = ?
			JOIN values ON
				values.attribute_id = attributes.id
				AND values.entity_id = entities.id
				AND values.%s = ?`,
			valToUseInQuery,
		),
		attribute.Name, expectedValue,
	)
}

// Creates a Entity of type "entityType" and its Values from the information provided in "attributeValues"
// not specified values are set with the default value (if exists) or a null value.
// entries in "attributeValues" that do not correspond to any attribute of "entityType" are ignored
//
// "attributeValues" are in {"attributeName": value} format
func (eavService *eavServiceImpl) CreateEntity(entityTypeName string, attributeValues map[string]any) (*models.Entity, error) {
	return ExecWithTransaction(
		eavService.db,
		func(tx *gorm.DB) (*models.Entity, error) {
			entityType, err := eavService.entityTypeRepository.GetByName(tx, entityTypeName)
			if err != nil {
				return nil, err
			}

			entity := models.NewEntity(entityType)
			for _, attribute := range entityType.Attributes {
				value, err := eavService.createValue(attribute, attributeValues)
				if err != nil {
					return nil, err
				}
				entity.Fields = append(entity.Fields, value)
			}

			return entity, eavService.entityRepository.Create(tx, entity)
		},
	)
}

// Creates a new Value for the "attribute" with the information provided in "attributesValues"
// or with the default value (if exists)
// or a null value if the value for "attribute" is not specified in "attributesValues"
//
// "attributesValues" is in {"attributeName": value} format
//
// returns error is the type of the value provided for "attribute"
// does not correspond with its ValueType
func (eavService *eavServiceImpl) createValue(attribute *models.Attribute, attributesValues map[string]any) (*models.Value, error) {
	attributeValue, isPresent := attributesValues[attribute.Name]
	if isPresent {
		value := &models.Value{Attribute: attribute, AttributeID: attribute.ID}
		err := eavService.updateValue(value, attributeValue)
		if err != nil {
			return nil, err
		}

		return value, nil
	}

	// attribute not present in params, set the default value (if exists) or a null value
	if attribute.Default {
		return attribute.GetNewDefaultValue()
	} else if attribute.Required {
		return nil, fmt.Errorf("field %s is missing and is required", attribute.Name)
	}
	return models.NewNullValue(attribute)
}

// Updates entity with type "entityTypeName" and "id" Values to the new values provided in "newValues"
// entries in "newValues" that do not correspond to any attribute of the EntityType are ignored
//
// "newValues" are in {"attributeName": newValue} format
func (eavService *eavServiceImpl) UpdateEntity(entityTypeName string, entityID uuid.UUID, newValues map[string]any) (*models.Entity, error) {
	return ExecWithTransaction(
		eavService.db,
		func(tx *gorm.DB) (*models.Entity, error) {
			entity, err := eavService.entityRepository.Get(tx, entityTypeName, entityID)
			if err != nil {
				return nil, err
			}

			for _, value := range entity.Fields {
				newValue, isPresent := newValues[value.Attribute.Name]
				if isPresent {
					err := eavService.updateValue(value, newValue)
					if err != nil {
						return nil, err
					}
				}
			}

			return entity, tx.Save(entity.Fields).Error
		},
	)
}

// Updates Value "value" to have the value "newValue"
//
// returns error is the type of the "newValue" does not correspond
// with the type expected for the "value"'s attribute
func (eavService *eavServiceImpl) updateValue(value *models.Value, newValue any) error {
	switch newValueTyped := newValue.(type) {
	case string:
		if value.Attribute.ValueType == models.RelationValueType {
			uuidVal, err := uuid.Parse(newValueTyped)
			if err != nil {
				return ErrCantParseUUID
			}

			// TODO verify that exists
			return value.SetRelationVal(uuidVal)
		}
		return value.SetStringVal(newValueTyped)
	case int:
		return value.SetIntVal(newValueTyped)
	case float64:
		if utils.IsAnInt(newValueTyped) && value.Attribute.ValueType == models.IntValueType {
			return value.SetIntVal(int(newValueTyped))
		}
		return value.SetFloatVal(newValueTyped)
	case bool:
		return value.SetBooleanVal(newValueTyped)
	case nil:
		return value.SetNull()
	default:
		return fmt.Errorf("unsupported type")
	}
}

// Deletes Entity with type "entityTypeName" and id "entityID" and its values
func (eavService *eavServiceImpl) DeleteEntity(entityTypeName string, entityID uuid.UUID) error {
	return ExecWithTransactionNoResponse(
		eavService.db,
		func(tx *gorm.DB) error {
			entity, err := eavService.entityRepository.Get(tx, entityTypeName, entityID)
			if err != nil {
				return err
			}

			err = tx.Where("entity_id = ?", entityID.String()).Delete(&models.Value{}).Error
			if err != nil {
				return err
			}

			return tx.Delete(entity).Error
		},
	)
}
