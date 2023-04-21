package models

import (
	"encoding/json"

	uuid "github.com/google/uuid"
)

// Describe an instance of an EntityType
type Entity struct {
	BaseModel
	Fields []*Value

	// GORM relations
	EntityTypeID uuid.UUID
	EntityType   *EntityType
}

func NewEntity(entityType *EntityType) *Entity {
	return &Entity{
		EntityType:   entityType,
		EntityTypeID: entityType.ID,
	}
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

func (e Entity) Equal(other Entity) bool {
	return e.ID == other.ID &&
		e.EntityType.ID == other.EntityType.ID
}

func (Entity) TableName() string {
	return "entities"
}
