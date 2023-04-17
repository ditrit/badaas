package models

import (
	"errors"
	"fmt"

	uuid "github.com/google/uuid"
)

// Describe the type of an attribute
type ValueTypeT string

const (
	RelationValueType ValueTypeT = "relation"
	BooleanValueType  ValueTypeT = "bool"
	StringValueType   ValueTypeT = "string"
	IntValueType      ValueTypeT = "int"
	FloatValueType    ValueTypeT = "float"
)

// Describe the attribute of a en EntityType
type Attribute struct {
	BaseModel
	Name     string
	Unique   bool
	Required bool

	Default bool // if there is a default value

	// Default values, only if Default == true
	DefaultInt    int
	DefaultBool   bool
	DefaultString string
	DefaultFloat  float64

	ValueType ValueTypeT // the type the values of this attr are. Can be "int", "float", "string", "bool", "relation"
	// id of the EntityType to which a RelationValueType points to. Only if ValueType == RelationValueType
	RelationTargetEntityTypeID uuid.UUID `gorm:"foreignKey:EntityType"`

	// GORM relations
	EntityTypeID uuid.UUID
}

var ErrNoDefaultValueSet = errors.New("no default value found")

func NewRelationAttribute(entityType *EntityType, name string,
	unique bool, required bool,
	relationTargetEntityType *EntityType) *Attribute {
	return &Attribute{
		EntityTypeID:               entityType.ID,
		Name:                       name,
		ValueType:                  RelationValueType,
		Required:                   required,
		Unique:                     unique,
		RelationTargetEntityTypeID: relationTargetEntityType.ID,
	}
}

// Get a new value with the default value associated with the attribute
func (a *Attribute) GetNewDefaultValue() (*Value, error) {
	if !a.Default {
		return nil, ErrNoDefaultValueSet
	}

	switch a.ValueType {
	case StringValueType:
		return NewStringValue(a, a.DefaultString)
	case IntValueType:
		return NewIntValue(a, a.DefaultInt)
	case FloatValueType:
		return NewFloatValue(a, a.DefaultFloat)
	case BooleanValueType:
		return NewBoolValue(a, a.DefaultBool)
	case RelationValueType:
		return nil, fmt.Errorf("can't provide default value for relations")
	default:
		return nil, fmt.Errorf("unsupported ValueType")
	}
}
