package models

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRelationValueNeedsEntity(t *testing.T) {
	brandEtt := &EntityType{
		Name: "brand",
	}

	carEtt := &EntityType{
		Name: "car",
	}
	brandAttr := NewRelationAttribute(carEtt, "brand", false, true, brandEtt)
	carEtt.Attributes = []*Attribute{brandAttr}

	_, err := NewRelationValue(brandAttr, nil)
	assert.Error(t, err, "can't create a new relation with a nil entity pointer")
}

func TestRelationValueEntityHasToBeTheTargetOfTheAttribute(t *testing.T) {
	brandEttID := uuid.New()
	carEttID := uuid.New()

	brandEtt := &EntityType{
		BaseModel: BaseModel{
			ID: brandEttID,
		},
		Name: "brand",
	}

	carEtt := &EntityType{
		BaseModel: BaseModel{
			ID: carEttID,
		},
		Name: "car",
	}
	brandAttr := NewRelationAttribute(carEtt, "brand", false, true, brandEtt)
	carEtt.Attributes = []*Attribute{brandAttr}

	car := &Entity{
		EntityType: carEtt,
	}

	_, err := NewRelationValue(brandAttr, car)
	assert.Error(
		t, err,
		fmt.Sprintf(
			"can't create a relation with an entity of wrong EntityType (got the entityTypeID=%s, expected=%s)",
			carEttID.String(),
			brandEttID.String(),
		),
	)
}

func TestRelationValueWithCorrespondingEntity(t *testing.T) {
	brandEtt := &EntityType{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},
		Name: "brand",
	}

	carEtt := &EntityType{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},
		Name: "car",
	}
	brandAttr := NewRelationAttribute(carEtt, "brand", false, true, brandEtt)
	carEtt.Attributes = []*Attribute{brandAttr}

	brand := &Entity{
		EntityType: brandEtt,
	}

	value, err := NewRelationValue(brandAttr, brand)
	assert.Nil(t, err)

	assert.Equal(t, value.Value(), brand.ID)
}

func TestNewNullValueRespondErrorIfAttributeIsRequired(t *testing.T) {
	requiredAttr := &Attribute{
		Name:      "attr",
		ValueType: StringValueType,
		Required:  true,
	}

	_, err := NewNullValue(requiredAttr)
	assert.Error(t, err, "can't create new null value for a required attribute")
}

func TestNewNullValueWorksIfAttributeIsNotRequired(t *testing.T) {
	notRequiredAttr := &Attribute{
		Name:      "attr",
		ValueType: StringValueType,
		Required:  false,
	}

	value, err := NewNullValue(notRequiredAttr)
	assert.Nil(t, err)
	assert.Nil(t, value.Value())
}

func TestNewIntValueRespondErrorIsAttributeIfNotIntType(t *testing.T) {
	stringAttr := &Attribute{
		Name:      "attr",
		ValueType: StringValueType,
	}

	_, err := NewIntValue(stringAttr, 5)
	assert.ErrorIs(t, err, ErrAskingForWrongType)
}

func TestNewIntValueWorks(t *testing.T) {
	intAttr := &Attribute{
		Name:      "attr",
		ValueType: IntValueType,
	}

	value, err := NewIntValue(intAttr, 5)
	assert.Nil(t, err)
	assert.Equal(t, value.Value(), 5)
}

func TestNewFloatValueRespondErrorIsAttributeIfNotFloatType(t *testing.T) {
	stringAttr := &Attribute{
		Name:      "attr",
		ValueType: StringValueType,
	}

	_, err := NewFloatValue(stringAttr, 5.5)
	assert.ErrorIs(t, err, ErrAskingForWrongType)
}

func TestNewFloatValueWorks(t *testing.T) {
	attr := &Attribute{
		Name:      "attr",
		ValueType: FloatValueType,
	}

	value, err := NewFloatValue(attr, 5.5)
	assert.Nil(t, err)
	assert.Equal(t, value.Value(), 5.5)
}

func TestNewBoolValueRespondErrorIsAttributeIfNotBoolType(t *testing.T) {
	stringAttr := &Attribute{
		Name:      "attr",
		ValueType: StringValueType,
	}

	_, err := NewBoolValue(stringAttr, true)
	assert.ErrorIs(t, err, ErrAskingForWrongType)
}

func TestNewBoolValueWorks(t *testing.T) {
	attr := &Attribute{
		Name:      "attr",
		ValueType: BooleanValueType,
	}

	value, err := NewBoolValue(attr, true)
	assert.Nil(t, err)
	assert.Equal(t, value.Value(), true)
}

func TestNewStringValueRespondErrorIsAttributeIfNotStringType(t *testing.T) {
	boolAttr := &Attribute{
		Name:      "attr",
		ValueType: BooleanValueType,
	}

	_, err := NewStringValue(boolAttr, "salut")
	assert.ErrorIs(t, err, ErrAskingForWrongType)
}

func TestNewStringValueWorks(t *testing.T) {
	attr := &Attribute{
		Name:      "attr",
		ValueType: StringValueType,
	}

	value, err := NewStringValue(attr, "salut")
	assert.Nil(t, err)
	assert.Equal(t, value.Value(), "salut")
}

func TestNewRelationValueRespondErrorIsAttributeIfNotRelationType(t *testing.T) {
	brandEtt := &EntityType{
		Name: "brand",
	}
	brand := &Entity{
		EntityType: brandEtt,
	}

	boolAttr := &Attribute{
		Name:      "attr",
		ValueType: BooleanValueType,
	}

	_, err := NewRelationValue(boolAttr, brand)
	assert.ErrorIs(t, err, ErrAskingForWrongType)
}

func TestSetNullWorks(t *testing.T) {
	attr := &Attribute{
		Name:      "attr",
		ValueType: StringValueType,
	}

	value, err := NewStringValue(attr, "salut")
	assert.Nil(t, err)
	value.SetNull()
	assert.Nil(t, value.Value())
}
