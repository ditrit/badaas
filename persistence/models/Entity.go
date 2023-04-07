package models

import (
	"encoding/json"
	"fmt"

	"github.com/ditrit/badaas/services/eavservice/utils"
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
	var attrID uuid.UUID
	for _, a := range e.EntityType.Attributes {
		if a.Name == attrName {
			attrID = a.ID
			break
		}
	}
	if uuid.Nil == attrID {
		return nil, fmt.Errorf("attr not found: got=%s", attrName)
	}
	for _, v := range e.Fields {
		if v.AttributeID == attrID {
			return v.Value(), nil
		}
	}
	return nil, fmt.Errorf("value not found")
}

// Encode the Entity to json
func (e *Entity) EncodeToJSON() []byte {
	var pairs []string
	pairs = append(pairs,
		fmt.Sprintf("%q: %d", "id", e.ID),
		fmt.Sprintf("%q: %q", "type", e.EntityType.Name),
		fmt.Sprintf("%q: %s", "attrs", e.encodeAttributesold()),
	)
	return []byte(utils.BuildJSONFromStrings(pairs))
}

// return the attribute in a json encoded string
func (e *Entity) encodeAttributesold() string {
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
