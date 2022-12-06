package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMarshall(t *testing.T) {
	id := uuid.MustParse("11e1d4b8-701d-47cc-852f-6d36922bcc75")
	ett := &EntityType{
		BaseModel: BaseModel{ID: id},
		Name:      "bird",
		Attributs: []*Attribut{
			&Attribut{
				BaseModel: BaseModel{
					ID: id,
				},
				Name:          "color",
				DefaultString: "red",
				ValueType:     StringValueType,
				EntityTypeId:  id,
			}},
	}

	et := &Entity{
		BaseModel: BaseModel{
			ID: id,
		},
		Fields: []*Value{&Value{
			BaseModel: BaseModel{
				ID: id,
			},

			IsNull:     false,
			StringVal:  "blue",
			EntityId:   id,
			AttributId: id,
			Attribut:   ett.Attributs[0],
		}},
		EntityTypeId: id,
		EntityType:   ett,
	}

	p, _ := et.MarshalJSON()
	assert.JSONEq(t, `{"attrs":{"color":"blue"},"id":"11e1d4b8-701d-47cc-852f-6d36922bcc75","type":"bird", "createdAt":"0001-01-01T00:00:00Z", "updatedAt":"0001-01-01T00:00:00Z"}`, string(p))
}
