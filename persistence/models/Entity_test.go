package models

import (
	"testing"

	uuid "github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestMarshall(t *testing.T) {
	id := uuid.MustParse("11e1d4b8-701d-47cc-852f-6d36922bcc75")
	ett := &EntityType{
		BaseModel: BaseModel{ID: id},
		Name:      "bird",
		Attributes: []*Attribute{
			{
				BaseModel: BaseModel{
					ID: id,
				},
				Name:          "color",
				DefaultString: "red",
				ValueType:     StringValueType,
				EntityTypeID:  id,
			}},
	}

	et := &Entity{
		BaseModel: BaseModel{
			ID: id,
		},
		Fields: []*Value{{
			BaseModel: BaseModel{
				ID: id,
			},

			IsNull:      false,
			StringVal:   "blue",
			EntityID:    id,
			AttributeID: id,
			Attribute:   ett.Attributes[0],
		}},
		EntityTypeID: id,
		EntityType:   ett,
	}

	p, _ := et.MarshalJSON()
	assert.JSONEq(t, `{"attrs":{"color":"blue"},"id":"11e1d4b8-701d-47cc-852f-6d36922bcc75","type":"bird", "createdAt":"0001-01-01T00:00:00Z", "updatedAt":"0001-01-01T00:00:00Z"}`, string(p))
}
