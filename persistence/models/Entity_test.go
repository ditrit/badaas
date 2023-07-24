package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/badorm"
)

func TestMarshall(t *testing.T) {
	id := badorm.UUID(uuid.MustParse("11e1d4b8-701d-47cc-852f-6d36922bcc75"))
	ett := &EntityType{
		UUIDModel: badorm.UUIDModel{ID: id},
		Name:      "bird",
		Attributes: []*Attribute{
			{
				UUIDModel: badorm.UUIDModel{
					ID: id,
				},
				Name:          "color",
				DefaultString: "red",
				ValueType:     StringValueType,
				EntityTypeID:  id,
			},
			{
				UUIDModel: badorm.UUIDModel{
					ID: id,
				},
				Name:         "heigh",
				ValueType:    IntValueType,
				EntityTypeID: id,
			},
		},
	}

	et := &Entity{
		UUIDModel: badorm.UUIDModel{
			ID: id,
		},
		Fields: []*Value{
			{
				UUIDModel: badorm.UUIDModel{
					ID: id,
				},

				IsNull:      false,
				StringVal:   "blue",
				EntityID:    id,
				AttributeID: id,
				Attribute:   ett.Attributes[0],
			},
			{
				UUIDModel: badorm.UUIDModel{
					ID: id,
				},

				IsNull:      true,
				EntityID:    id,
				AttributeID: id,
				Attribute:   ett.Attributes[1],
			},
		},
		EntityTypeID: id,
		EntityType:   ett,
	}

	p, _ := et.MarshalJSON()
	assert.JSONEq(
		t,
		`{"attrs":{"color":"blue"},"id":"11e1d4b8-701d-47cc-852f-6d36922bcc75","type":"bird", "createdAt":"0001-01-01T00:00:00Z", "updatedAt":"0001-01-01T00:00:00Z"}`, //nolint:lll
		string(p),
	)
}
