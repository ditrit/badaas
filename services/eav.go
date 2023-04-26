package services

import (
	"errors"
	"fmt"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrCantParseUUID = errors.New("can't parse uuid")
)

// EAV service provide handle EAV objects
type EAVService interface {
	GetEntity(entityTypeName string, id uuid.UUID) (*models.Entity, error)
	GetEntities(entityTypeName string, conditions map[string]any) ([]*models.Entity, error)
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
func (eavService *eavServiceImpl) GetEntities(entityTypeName string, conditions map[string]any) ([]*models.Entity, error) {
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
					query, err = eavService.addValueCheckToQuery(query, attribute.Name, expectedValue)
					if err != nil {
						return nil, err
					}
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
func (eavService *eavServiceImpl) addValueCheckToQuery(query *gorm.DB, attributeName string, expectedValue any) (*gorm.DB, error) {
	stringQuery := fmt.Sprintf(
		`JOIN attributes attributes_%[1]s ON
			attributes_%[1]s.entity_type_id = entities.entity_type_id
			AND attributes_%[1]s.name = ?
		JOIN values values_%[1]s ON
			values_%[1]s.attribute_id = attributes_%[1]s.id
			AND values_%[1]s.entity_id = entities.id
		`,
		attributeName,
	)
	switch expectedValueTyped := expectedValue.(type) {
	case float64:
		stringQuery += fmt.Sprintf(
			"AND ((%s) OR (%s))",
			getQueryCheckValueOfType(attributeName, models.IntValueType),
			getQueryCheckValueOfType(attributeName, models.FloatValueType),
		)
	case bool:
		stringQuery += "AND " + getQueryCheckValueOfType(attributeName, models.BooleanValueType)
	case string:
		_, err := uuid.Parse(expectedValueTyped)
		if err == nil {
			stringQuery += "AND " + getQueryCheckValueOfType(attributeName, models.RelationValueType)
		} else {
			stringQuery += "AND " + getQueryCheckValueOfType(attributeName, models.StringValueType)
		}
	case nil:
		stringQuery += fmt.Sprintf(
			"AND values_%s.is_null = true",
			attributeName,
		)
	default:
		return nil, fmt.Errorf("unsupported type")
	}

	query = query.Joins(stringQuery, attributeName, expectedValue, expectedValue)

	return query, nil
}

func getQueryCheckValueOfType(attributeName string, valueType models.ValueTypeT) string {
	return fmt.Sprintf(
		"attributes_%[1]s.value_type = '%[2]s' AND values_%[1]s.%[2]s_val = ?",
		attributeName, string(valueType),
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
				value, err := eavService.createValue(tx, attribute, attributeValues)
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
func (eavService *eavServiceImpl) createValue(tx *gorm.DB, attribute *models.Attribute, attributesValues map[string]any) (*models.Value, error) {
	attributeValue, isPresent := attributesValues[attribute.Name]
	if isPresent {
		value := &models.Value{Attribute: attribute, AttributeID: attribute.ID}
		err := eavService.updateValue(tx, value, attributeValue)
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
					err := eavService.updateValue(tx, value, newValue)
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
func (eavService *eavServiceImpl) updateValue(tx *gorm.DB, value *models.Value, newValue any) error {
	switch newValueTyped := newValue.(type) {
	case string:
		if value.Attribute.ValueType == models.RelationValueType {
			return eavService.updateRelationValue(tx, value, newValueTyped)
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

// Updates Value "value" to point to "newID"
//
// returns error if entity with "newID" does not exist
// or its type is not the same as the one pointed by the attribute's RelationTargetEntityTypeID
func (eavService *eavServiceImpl) updateRelationValue(tx *gorm.DB, value *models.Value, newID string) error {
	uuidVal, err := uuid.Parse(newID)
	if err != nil {
		return ErrCantParseUUID
	}

	relationEntityType, err := eavService.entityTypeRepository.Get(tx, value.Attribute.RelationTargetEntityTypeID)
	if err != nil {
		return err
	}

	relationEntity, err := eavService.entityRepository.Get(tx, relationEntityType.Name, uuidVal)
	if err != nil {
		return err
	}

	return value.SetRelationVal(relationEntity)
}

// Deletes Entity with type "entityTypeName" and id "entityID" and its values
func (eavService *eavServiceImpl) DeleteEntity(entityTypeName string, entityID uuid.UUID) error {
	return eavService.db.Transaction(func(tx *gorm.DB) error {
		entity, err := eavService.entityRepository.Get(tx, entityTypeName, entityID)
		if err != nil {
			return err
		}

		err = tx.Where("entity_id = ?", entityID.String()).Delete(&models.Value{}).Error
		if err != nil {
			return err
		}

		return tx.Delete(entity).Error
	})
}
