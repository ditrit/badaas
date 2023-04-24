package eavservice

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
	ErrIDDontMatchEntityType = errors.New("this object doesn't belong to this type")
	ErrCantParseUUID         = errors.New("can't parse uuid")
)

// EAV service provide handle EAV objects
type EAVService interface {
	GetEntityTypeByName(name string) (*models.EntityType, error)
	GetEntitiesWithParams(ett *models.EntityType, params map[string]string) []*models.Entity
	DeleteEntity(et *models.Entity) error
	GetEntity(ett *models.EntityType, id uuid.UUID) (*models.Entity, error)
	CreateEntity(ett *models.EntityType, attrs map[string]any) (*models.Entity, error)
	UpdateEntity(et *models.Entity, attrs map[string]any) error
}

type eavServiceImpl struct {
	logger           *zap.Logger
	db               *gorm.DB
	entityRepository *repository.EntityRepository
}

func NewEAVService(
	logger *zap.Logger,
	db *gorm.DB,
	entityRepository *repository.EntityRepository,
) EAVService {
	return &eavServiceImpl{
		logger:           logger,
		db:               db,
		entityRepository: entityRepository,
	}
}

// Get EntityType by name (string)
func (eavService *eavServiceImpl) GetEntityTypeByName(name string) (*models.EntityType, error) {
	var ett models.EntityType
	err := eavService.db.Preload("Attributes").First(&ett, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("EntityType named %q not found", name)
		}
	}
	return &ett, nil
}

// Get entities of type entityType that match all conditions given in params
// params are in {"attributeName": expectedValue} format
// TODO relations join
func (eavService *eavServiceImpl) GetEntitiesWithParams(entityType *models.EntityType, params map[string]string) []*models.Entity {
	query := eavService.db.Select("entities.*")

	// only entities that match the conditions in params
	for _, attribute := range entityType.Attributes {
		expectedValue, isPresent := params[attribute.Name]
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
	query.Find(&entities)

	return entities
}

// Adds to the query the verification that the value for the attribute is expectedValue
func (eavService *eavServiceImpl) addValueCheckToQuery(query *gorm.DB, attribute *models.Attribute, expectedValue string) *gorm.DB {
	var valToUseInQuery string
	if expectedValue != "null" { // TODO should be changed to be able to use nil
		switch attribute.ValueType {
		case models.StringValueType:
			valToUseInQuery = "string_val"
		case models.IntValueType:
			valToUseInQuery = "int_val"
		case models.FloatValueType:
			valToUseInQuery = "float_val"
		case models.BooleanValueType:
			valToUseInQuery = "bool_val"
		case models.RelationValueType:
			valToUseInQuery = "relation_val"
		}
	} else {
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

// Delete an entity and its values
func (eavService *eavServiceImpl) DeleteEntity(et *models.Entity) error {
	return eavService.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("entity_id = ?", et.ID.String()).Delete(&models.Value{}).Error
		if err != nil {
			return err
		}

		return tx.Delete(et).Error
	})
}

func (eavService *eavServiceImpl) GetEntity(ett *models.EntityType, id uuid.UUID) (*models.Entity, error) {
	var et models.Entity
	err := eavService.db.Preload("Fields").Preload("Fields.Attribute").Preload("EntityType.Attributes").Preload("EntityType").First(&et, id).Error
	if err != nil {
		return nil, err
	}
	if ett.ID != et.EntityTypeID {
		return nil, ErrIDDontMatchEntityType
	}
	return &et, nil
}

// Create a brand new entity
func (eavService *eavServiceImpl) CreateEntity(entityType *models.EntityType, params map[string]any) (*models.Entity, error) {
	entity := models.NewEntity(entityType)
	for _, attribute := range entityType.Attributes {
		value, err := eavService.createValueFromParams(attribute, params)
		if err != nil {
			return nil, err
		}
		entity.Fields = append(entity.Fields, value)
	}

	return entity, eavService.entityRepository.Create(entity)
}

func (eavService *eavServiceImpl) createValueFromParams(attribute *models.Attribute, params map[string]any) (*models.Value, error) {
	attributeValue, isPresent := params[attribute.Name]
	if isPresent {
		switch attributeValueTyped := attributeValue.(type) {
		case string:
			if attribute.ValueType == models.RelationValueType {
				uuidVal, err := uuid.Parse(attributeValueTyped)
				if err != nil {
					return nil, ErrCantParseUUID
				}
				// TODO verify that exists
				return models.NewRelationIDValue(attribute, uuidVal)
			}
			return models.NewStringValue(attribute, attributeValueTyped)
		case int:
			return models.NewIntValue(attribute, attributeValueTyped)
		case float64:
			if utils.IsAnInt(attributeValueTyped) && attribute.ValueType == models.IntValueType {
				return models.NewIntValue(attribute, int(attributeValueTyped))
			}
			return models.NewFloatValue(attribute, attributeValueTyped)
		case bool:
			return models.NewBoolValue(attribute, attributeValueTyped)
		case nil:
			return models.NewNullValue(attribute)
		default:
			return nil, fmt.Errorf("unsupported type")
		}
	}

	if attribute.Default {
		return attribute.GetNewDefaultValue()
	} else if attribute.Required {
		return nil, fmt.Errorf("field %s is missing and is required", attribute.Name)
	}
	return models.NewNullValue(attribute)
}
func (eavService *eavServiceImpl) UpdateEntity(et *models.Entity, attrs map[string]any) error {
	for _, value := range et.Fields {
		attribute := value.Attribute
		for k, v := range attrs {
			if k == attribute.Name {
				switch t := v.(type) {
				case string:
					if attribute.ValueType == models.RelationValueType {
						uuidVal, err := uuid.Parse(t)
						if err != nil {
							return ErrCantParseUUID
						}

						// TODO verify that exists
						err = value.SetRelationVal(uuidVal)
						if err != nil {
							return err
						}
					} else {
						err := value.SetStringVal(t)
						if err != nil {
							return err
						}
					}
				case int:
					err := value.SetIntVal(t)
					if err != nil {
						return err
					}
				case float64:
					if utils.IsAnInt(t) && attribute.ValueType == models.IntValueType {
						err := value.SetIntVal(int(t))
						if err != nil {
							return err
						}
					} else {
						err := value.SetFloatVal(t)
						if err != nil {
							return err
						}
					}
				case bool:
					err := value.SetBooleanVal(t)
					if err != nil {
						return err
					}
				case nil:
					err := value.SetNull()
					if err != nil {
						return err
					}

				default:
					return fmt.Errorf("unsupported type")
				}
			}
		}
	}

	return eavService.db.Save(et.Fields).Error
}
