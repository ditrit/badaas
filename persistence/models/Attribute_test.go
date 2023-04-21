package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultValueReturnsErrorIfNotDefault(t *testing.T) {
	attr := &Attribute{
		Default: false,
	}

	_, err := attr.GetNewDefaultValue()
	assert.ErrorIs(t, err, ErrNoDefaultValueSet)
}

func TestNewDefaultValueReturnsDefaultIntForInt(t *testing.T) {
	attr := &Attribute{
		ValueType:  IntValueType,
		Default:    true,
		DefaultInt: 1,
	}

	value, err := attr.GetNewDefaultValue()
	assert.Nil(t, err)
	assert.Equal(t, 1, value.Value())
}

func TestNewDefaultValueReturnsDefaultFloatForFloat(t *testing.T) {
	attr := &Attribute{
		ValueType:    FloatValueType,
		Default:      true,
		DefaultFloat: 1.0,
	}

	value, err := attr.GetNewDefaultValue()
	assert.Nil(t, err)
	assert.Equal(t, 1.0, value.Value())
}

func TestNewDefaultValueReturnsDefaultStringForString(t *testing.T) {
	attr := &Attribute{
		ValueType:     StringValueType,
		Default:       true,
		DefaultString: "salut",
	}

	value, err := attr.GetNewDefaultValue()
	assert.Nil(t, err)
	assert.Equal(t, "salut", value.Value())
}

func TestNewDefaultValueReturnsDefaultBoolForBool(t *testing.T) {
	attr := &Attribute{
		ValueType:   BooleanValueType,
		Default:     true,
		DefaultBool: true,
	}

	value, err := attr.GetNewDefaultValue()
	assert.Nil(t, err)
	assert.Equal(t, true, value.Value())
}

func TestNewDefaultValueReturnsErrorForRelation(t *testing.T) {
	attr := &Attribute{
		Default:   true,
		ValueType: RelationValueType,
	}

	_, err := attr.GetNewDefaultValue()
	assert.Error(t, err, "can't provide default value for relations")
}

func TestNewDefaultValueReturnsErrorForUnsupportedType(t *testing.T) {
	attr := &Attribute{
		Default:   true,
		ValueType: "something else",
	}

	_, err := attr.GetNewDefaultValue()
	assert.Error(t, err, "unsupported ValueType")
}
