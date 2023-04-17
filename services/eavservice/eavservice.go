package eavservice

import (
	"errors"
	"fmt"

	uuid "github.com/google/uuid"

	"github.com/ditrit/badaas/persistence/models"
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
	CreateEntity(ett *models.EntityType, attrs map[string]interface{}) (*models.Entity, error)
	UpdateEntity(et *models.Entity, attrs map[string]interface{}) error
}

type eavServiceImpl struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewEAVService(
	logger *zap.Logger,
	db *gorm.DB,
) EAVService {
	return &eavServiceImpl{
		logger: logger,
		db:     db,
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

func (eavService *eavServiceImpl) GetEntitiesWithParams(ett *models.EntityType, params map[string]string) []*models.Entity {
	var entities []*models.Entity

	// TODO check pivot and cross join
	// TODO distinct
	// query := eavService.db.Select("entities.id").Joins(
	// 	"JOIN attributes ON attributes.entity_type_id = ?",
	// 	ett.ID,
	// ).Joins(
	// 	`JOIN values ON
	// 		values.entity_id = entities.id
	// 		AND values.attribute_id = attributes.id`,
	// )

	// eavService.db.Raw(
	// 	`SELECT entities.*
	// 	FROM entities
	// 	JOIN attributes
	// 		ON attributes.entity_type_id = entities.entity_type_id
	// 	JOIN values ON
	// 		values.entity_id = entities.id AND
	// 		values.attribute_id = attributes.id
	// 	// TODO pivot table
	// 	WHERE
	// 		entities.entity_type_id = ?`,
	// 	ett.ID,
	// ).Preload("Fields").Preload("Fields.Attribute").Preload("EntityType.Attributes").Preload("EntityType").Find(&entities)
	// TODO hacer version sin gorm para no tener que hacer todas estas queries al pedo.
	// TODO tambien hacer una version que no usa gorm para evitar las queries multiples hechas para buscar el resto de objetos.

	// multiple joins version
	// TODO filter by null
	// TODO relations

	query := eavService.db.Select("entities.*")

	for _, attr := range ett.Attributes {
		v, present := params[attr.Name]
		if present {
			var valToUse string

			// TODO check that v is a valid value for that type to avoid error in db
			switch attr.ValueType {
			case models.StringValueType:
				valToUse = "string_val"
			case models.IntValueType:
				valToUse = "int_val"
			case models.FloatValueType:
				valToUse = "float_val"
			case models.BooleanValueType:
				valToUse = "bool_val"
			}

			query = query.Joins(
				fmt.Sprintf(`
					JOIN attributes ON
						attributes.entity_type_id = entities.entity_type_id
						AND attributes.name = ?
					JOIN values ON
						values.attribute_id = attributes.id
						AND values.entity_id = entities.id
						AND values.%s = ?`,
					valToUse,
				),
				attr.Name, v,
			)
		}
	}

	query = query.Where(
		"entities.entity_type_id = ?",
		ett.ID,
	)
	query = query.Preload("Fields").Preload("Fields.Attribute").Preload("EntityType.Attributes").Preload("EntityType")
	query.Find(&entities)

	return entities

	// version con un solo join
	// eavService.db.Joins(
	// 	"JOIN values ON values.entity_id = entities.id",
	// ).Joins(
	// 	`JOIN attributes ON
	// 		values.attribute_id = attributes.id
	// 		AND attributes.name IN ?`,
	// 	maps.Keys(params),
	// ).Where(
	// 	`entities.entity_type_id = ?
	// 	AND values.value = ?`,
	// 	ett.ID,
	// 	params["no se que poner aca"],
	// ).Find(&entities)

	// version all in memory
	// eavService.db.Where(
	// 	"entity_type_id = ?",
	// 	ett.ID,
	// ).Preload("Fields").Preload("Fields.Attribute").Preload("EntityType.Attributes").Preload("EntityType").Find(&entities)

	// resultSet := make([]*models.Entity, 0, len(entities))
	// var keep bool
	// for _, et := range entities {
	// 	keep = true
	// 	for _, value := range et.Fields {
	// 		for k, v := range params {
	// 			if k == value.Attribute.Name {
	// 				switch value.Attribute.ValueType {
	// 				case models.StringValueType:
	// 					if v != value.StringVal {
	// 						keep = false
	// 					}
	// 				case models.IntValueType:
	// 					intVal, err := strconv.Atoi(v)
	// 					if err != nil {
	// 						break
	// 					}
	// 					if intVal != value.IntVal {
	// 						keep = false
	// 					}
	// 				case models.FloatValueType:
	// 					floatVal, err := strconv.ParseFloat(v, 64)
	// 					if err != nil {
	// 						break
	// 					}
	// 					if floatVal != value.FloatVal {
	// 						keep = false
	// 					}
	// 				case models.RelationValueType:
	// 					uuidVal, err := uuid.Parse(v)
	// 					if err != nil {
	// 						keep = false
	// 					}
	// 					if uuidVal != value.RelationVal {
	// 						keep = false
	// 					}
	// 				case models.BooleanValueType:
	// 					boolVal, err := strconv.ParseBool(v)
	// 					if err != nil {
	// 						break
	// 					}
	// 					if boolVal != value.BoolVal {
	// 						keep = false
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// 	if keep {
	// 		resultSet = append(resultSet, et)
	// 	}
	// }
	// return resultSet
}

// Delete an entity
func (eavService *eavServiceImpl) DeleteEntity(et *models.Entity) error {
	for _, v := range et.Fields {
		err := eavService.db.Delete(v).Error
		if err != nil {
			return err
		}
	}
	return eavService.db.Delete(et).Error
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
func (eavService *eavServiceImpl) CreateEntity(ett *models.EntityType, attrs map[string]interface{}) (*models.Entity, error) {
	var et models.Entity
	for _, a := range ett.Attributes {
		present := false

		var value *models.Value
		var err error

		for k, v := range attrs {
			if k == a.Name {
				present = true

				switch t := v.(type) {
				case string:
					if a.ValueType == models.RelationValueType {
						uuidVal, err := uuid.Parse(t)
						if err != nil {
							return nil, ErrCantParseUUID
						}
						value, err = models.NewRelationIDValue(a, uuidVal)
						if err != nil {
							return nil, err
						}
					} else {
						value, err = models.NewStringValue(a, t)
						if err != nil {
							return nil, err
						}
					}
				case int:
					value, err = models.NewIntValue(a, t)
					if err != nil {
						return nil, err
					}
				case float64:
					if utils.IsAnInt(t) && a.ValueType == models.IntValueType {
						value, err = models.NewIntValue(a, int(t))
						if err != nil {
							return nil, err
						}
					} else {
						value, err = models.NewFloatValue(a, t)
						if err != nil {
							return nil, err
						}
					}
				case bool:
					value, err = models.NewBoolValue(a, t)
					if err != nil {
						return nil, err
					}
				// TODO is this really necessary?
				case nil:
					value, err = models.NewNullValue(a)
					if err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("unsupported type")
				}
			}
		}

		if !present {
			if a.Required {
				return nil, fmt.Errorf("field %q is missing and is required", a.Name)
			}

			if a.Default {
				value, err = a.GetNewDefaultValue()
				if err != nil {
					return nil, err
				}
			} else {
				// TODO is this really necessary?
				value, err = models.NewNullValue(a)
				if err != nil {
					return nil, err
				}
			}
		}

		et.Fields = append(et.Fields, value)
	}

	et.EntityType = ett
	return &et, eavService.db.Create(&et).Error
}

func (eavService *eavServiceImpl) UpdateEntity(et *models.Entity, attrs map[string]interface{}) error {
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
				// TODO is this really necessary?
				case nil:
					if attribute.Required {
						return fmt.Errorf("can't set a required variable to null")
					}
					value.SetNull()

				default:
					return fmt.Errorf("unsupported type")
				}
			}
		}
		eavService.db.Save(value)
	}

	return nil
}
