package models

import (
	"encoding/json"
	"fmt"

	"github.com/ditrit/badaas/utils"
	"github.com/google/uuid"
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

// Encode the Entity to json
func (e *Entity) EncodeToJSON() []byte {
	var pairs []string
	pairs = append(pairs,
		fmt.Sprintf("%q: %d", "id", e.ID),
		fmt.Sprintf("%q: %q", "type", e.EntityType.Name),
		fmt.Sprintf("%q: %s", "attrs", e.encodeAttributesToJSON()),
	)

	return []byte(utils.BuildJSONFromStrings(pairs))
}

// return the attribute in a json encoded string
func (e *Entity) encodeAttributesToJSON() string {
	var pairs []string
	for _, f := range e.Fields {
		if f.IsNull {
			continue
		}
		pair, _ := f.BuildJSONKVPair()
		pairs = append(pairs, pair)
	}

	return utils.BuildJSONFromStrings(pairs)
}

func (Entity) TableName() string {
	return "entities"
}

func (e Entity) Equal(other Entity) bool {
	return e.ID == other.ID &&
		e.EntityTypeID == other.EntityTypeID
}
