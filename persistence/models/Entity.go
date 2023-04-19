package models

import (
	"encoding/json"
	"fmt"

	uuid "github.com/google/uuid"

	"github.com/ditrit/badaas/utils"
)

// Describe an instance of an EntityType
type Entity struct {
	BaseModel
	Fields []*Value

	// GORM relations
	EntityTypeID uuid.UUID
	EntityType   *EntityType
}

// Encode the entity to json
// use the [encoding/json.Marshaler] interface
func (e *Entity) MarshalJSON() ([]byte, error) {
	dto := make(map[string]any)
	dto["id"] = e.ID
	dto["type"] = e.EntityType.Name
	dto["createdAt"] = e.CreatedAt
	dto["updatedAt"] = e.UpdatedAt
	dto["attrs"] = e.encodeAttributes()

	return json.Marshal(dto)
}

// return the attribute in a json encoded string
func (e *Entity) encodeAttributes() map[string]any {
	pairs := make(map[string]any, len(e.Fields))
	for _, field := range e.Fields {
		if field.IsNull {
			continue
		}
		pairs[field.Attribute.Name] = field.Value()
	}

	return pairs
}

func (e *Entity) GetValue(attrName string) (interface{}, error) {
	value := utils.FindFirst(e.Fields,
		func(v *Value) bool {
			return v.Attribute.Name == attrName
		},
	)
	if value == nil {
		return nil, fmt.Errorf("value for attribute %s not found", attrName)
	}

	return (*value).Value(), nil
}

func (Entity) TableName() string {
	return "entities"
}

func (e Entity) Equal(other Entity) bool {
	return e.ID == other.ID &&
		e.EntityType.ID == other.EntityType.ID
}
