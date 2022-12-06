package eavservice

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/eavservice/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrIdDontMatchEntityType = errors.New("this object doesn't belong to this type")
	ErrCantParseUUID         = errors.New("can't parse uuid")
)

// EAV service provide handle EAV objets
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
	err := eavService.db.Preload("Attributs").First(&ett, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(" EntityType named %q not found", name)
		}
	}
	return &ett, nil
}

func (eavService *eavServiceImpl) GetEntitiesWithParams(ett *models.EntityType, params map[string]string) []*models.Entity {
	var ets []*models.Entity
	eavService.db.Where("entity_type_id = ?", ett.ID).Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").Find(&ets)
	resultSet := make([]*models.Entity, 0, len(ets))
	var keep bool
	for _, et := range ets {
		keep = true

		for _, value := range et.Fields {
			for k, v := range params {
				if k == value.Attribut.Name {
					switch value.Attribut.ValueType {
					case models.StringValueType:
						if v != value.StringVal {
							keep = false
						}
					case models.IntValueType:
						intVal, err := strconv.Atoi(v)
						if err != nil {
							break
						}
						if intVal != value.IntVal {
							keep = false
						}
					case models.FloatValueType:
						floatVal, err := strconv.ParseFloat(v, 64)
						if err != nil {
							break
						}
						if floatVal != value.FloatVal {
							keep = false
						}
					case models.RelationValueType:
						uuidVal, err := uuid.Parse(v)
						if err != nil {
							keep = false
						}
						if uuidVal != value.RelationVal {
							keep = false
						}
					case models.BooleanValueType:
						boolVal, err := strconv.ParseBool(v)
						if err != nil {
							break
						}
						if boolVal != value.BoolVal {
							keep = false
						}
					}
				}
			}
		}
		if keep {
			resultSet = append(resultSet, et)
		}
	}
	return resultSet
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
	err := eavService.db.Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").First(&et, id).Error
	if err != nil {
		return nil, err
	}
	if ett.ID != et.EntityTypeId {
		return nil, ErrIdDontMatchEntityType
	}
	return &et, nil
}

// Create a brand new entity
func (eavService *eavServiceImpl) CreateEntity(ett *models.EntityType, attrs map[string]interface{}) (*models.Entity, error) {
	var et models.Entity
	for _, a := range ett.Attributs {
		present := false
		var value models.Value
		for k, v := range attrs {
			if k == a.Name {
				present = true
				switch t := v.(type) {
				case string:
					if a.ValueType == models.RelationValueType {
						uuidVal, err := uuid.Parse(v.(string))
						if err != nil {
							return nil, ErrCantParseUUID
						}
						value = models.Value{RelationVal: uuidVal}
					} else if a.ValueType != models.StringValueType {
						return nil, fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
					}
					value = models.Value{StringVal: v.(string)}

				case float64:
					if a.ValueType != models.IntValueType &&
						a.ValueType != models.FloatValueType &&
						a.ValueType != models.RelationValueType {
						return nil, fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
					}
					if utils.IsAnInt(v.(float64)) {
						value = models.Value{IntVal: int(v.(float64))}
					} else {
						// is a float
						value = models.Value{FloatVal: v.(float64)}
					}

				case bool:
					if a.ValueType != models.BooleanValueType {
						return nil, fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
					}
					value = models.Value{BoolVal: v.(bool)}

				case nil:
					if a.Required {
						return nil, fmt.Errorf("can't have a null field with a required attribut")
					}
					value = models.Value{IsNull: true}

				default:
					panic("mmmh you just discovered a new json type (https://go.dev/blog/json#generic-json-with-interface)")
				}
			}
		}
		if a.Required && !present {
			return nil, fmt.Errorf("field %q is missing and is required", a.Name)
		}
		if !present {
			if a.Default {
				v, err := a.GetNewDefaultValue()
				if err != nil {
					return nil, err
				} else {
					value = *v
				}
			} else {
				value = models.Value{IsNull: true}
			}
		}
		value.Attribut = a
		et.Fields = append(et.Fields, &value)
	}
	et.EntityType = ett
	return &et, eavService.db.Create(&et).Error
}

func (eavService *eavServiceImpl) UpdateEntity(et *models.Entity, attrs map[string]interface{}) error {
	for _, a := range et.EntityType.Attributs {
		for _, value := range et.Fields {
			if a.ID == value.AttributId {
				for k, v := range attrs {
					if k == a.Name {
						switch t := v.(type) {
						case string:
							if a.ValueType == models.RelationValueType {
								uuidVal, err := uuid.Parse(v.(string))
								if err != nil {
									return ErrCantParseUUID
								}
								value.RelationVal = uuidVal
							} else if a.ValueType != models.StringValueType {
								return fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
							}
							value.StringVal = v.(string)
						case float64:
							if a.ValueType != models.IntValueType &&
								a.ValueType != models.FloatValueType &&
								a.ValueType != models.RelationValueType {
								return fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
							}
							if utils.IsAnInt(v.(float64)) {
								// is an int

								value.IntVal = int(v.(float64))
							} else {
								// is a float
								value.FloatVal = v.(float64)
							}

						case bool:
							if a.ValueType != models.BooleanValueType {
								return fmt.Errorf("types dont match (expected=%v got=%T)", a.ValueType, t)
							}
							value.BoolVal = v.(bool)

						case nil:
							if a.Required {
								return fmt.Errorf("can't set a required variable to null (expected=%v got=%T)", a.ValueType, t)
							}
							value.IsNull = true
							value.IntVal = 0
							value.FloatVal = 0.0
							value.StringVal = ""
							value.BoolVal = false
							value.RelationVal = uuid.Nil

						default:
							panic("mmmh you just discovered a new json type (https://go.dev/blog/json#generic-json-with-interface)")
						}
					}
				}
				value.Attribut = a
				eavService.db.Save(value)
			}
		}
	}
	return nil
}
