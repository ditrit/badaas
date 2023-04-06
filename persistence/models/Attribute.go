package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	RelationValueType ValueTypeT = "relation"
	BooleanValueType  ValueTypeT = "bool"
	StringValueType   ValueTypeT = "string"
	IntValueType      ValueTypeT = "int"
	FloatValueType    ValueTypeT = "float"
)

// Describe the type of an attribut
type ValueTypeT string

// Describe the attribut of a en EntityType
type Attribute struct {
	BaseModel
	Name     string
	Unique   bool
	Required bool

	// there is a default value
	Default bool

	// Default values
	DefaultInt    int
	DefaultBool   bool
	DefaultString string
	DefaultFloat  float64

	// the type the values of this attr are. Can be "int", "float", "string", "bool", "relation"
	ValueType          ValueTypeT
	TargetEntityTypeId uuid.UUID // name of the EntityType

	// GORM relations
	EntityTypeId uuid.UUID
}

var ErrNoDefaultValueSet = errors.New("no default value found")

// Get a new value with the default value associated with the attribut
func (a *Attribute) GetNewDefaultValue() (*Value, error) {
	if !a.Default {
		return nil, ErrNoDefaultValueSet
	}
	switch a.ValueType {
	case StringValueType:
		v, err := NewStringValue(a, a.DefaultString)
		if err != nil {
			return nil, err
		}
		return v, nil
	case IntValueType:
		v, err := NewIntValue(a, a.DefaultInt)
		if err != nil {
			return nil, err
		}
		return v, nil
	case FloatValueType:
		v, err := NewFloatValue(a, a.DefaultFloat)
		if err != nil {
			return nil, err
		}
		return v, nil
	case BooleanValueType:
		v, err := NewBoolValue(a, a.DefaultBool)
		if err != nil {
			return nil, err
		}
		return v, nil
	case RelationValueType:
		return nil, fmt.Errorf("can't provide default value for relations")
	default:
		panic("hmmm we are not supposed to be here")
	}
}
