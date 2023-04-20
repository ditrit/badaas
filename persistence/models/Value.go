package models

import (
	"errors"
	"fmt"

	uuid "github.com/google/uuid"
)

// Describe the attribute value of an Entity
type Value struct {
	BaseModel
	IsNull      bool
	StringVal   string
	FloatVal    float64
	IntVal      int
	BoolVal     bool
	RelationVal uuid.UUID `gorm:"type:uuid;foreignKey:Entity;index:fk_relation_val_entity"`

	// GORM relations
	EntityID    uuid.UUID
	AttributeID uuid.UUID
	Attribute   *Attribute
}

var ErrAskingForWrongType = errors.New("attribute type doesn't match")

// Create a new null value
func NewNullValue(attr *Attribute) (*Value, error) {
	value := &Value{Attribute: attr, AttributeID: attr.ID}
	return value, value.SetNull()
}

// Create a new int value
func NewIntValue(attr *Attribute, i int) (*Value, error) {
	value := &Value{Attribute: attr, AttributeID: attr.ID}
	return value, value.SetIntVal(i)
}

// Create a new bool value
func NewBoolValue(attr *Attribute, b bool) (*Value, error) {
	value := &Value{Attribute: attr, AttributeID: attr.ID}
	return value, value.SetBooleanVal(b)
}

// Create a new float value
func NewFloatValue(attr *Attribute, f float64) (*Value, error) {
	value := &Value{Attribute: attr, AttributeID: attr.ID}
	return value, value.SetFloatVal(f)
}

// Create a new string value
func NewStringValue(attr *Attribute, s string) (*Value, error) {
	value := &Value{Attribute: attr, AttributeID: attr.ID}
	return value, value.SetStringVal(s)
}

// Create a new relation value.
// If et is nil, then the function returns an error
// If et is of the wrong types, then the function returns an error
func NewRelationValue(attr *Attribute, et *Entity) (*Value, error) {
	if et == nil {
		return nil, fmt.Errorf("can't create a new relation with a nil entity pointer")
	}
	if et.EntityType.ID != attr.RelationTargetEntityTypeID {
		return nil, fmt.Errorf(
			"can't create a relation with an entity of wrong EntityType (got the entityTypeID=%s, expected=%s)",
			et.EntityType.ID.String(), attr.RelationTargetEntityTypeID.String(),
		)
	}

	return NewRelationIDValue(attr, et.ID)
}

// Create a new relation value.
// If et is nil, then the function returns an error
// If et is of the wrong types, then the function returns an error
func NewRelationIDValue(attr *Attribute, uuidVal uuid.UUID) (*Value, error) {
	value := &Value{Attribute: attr, AttributeID: attr.ID}
	return value, value.SetRelationVal(uuidVal)
}

// Check if the Value is whole. eg, no fields are nil
func (v *Value) CheckWhole() error {
	if v.Attribute == nil {
		return fmt.Errorf("the Attribute pointer is nil in Value at %v", v)
	}

	return nil
}

// Return the underlying value as an interface
func (v *Value) Value() any {
	if v.IsNull {
		return nil
	}

	switch v.Attribute.ValueType {
	case StringValueType:
		return v.StringVal
	case IntValueType:
		return v.IntVal
	case FloatValueType:
		return v.FloatVal
	case BooleanValueType:
		return v.BoolVal
	case RelationValueType:
		return v.RelationVal
	default:
		panic(fmt.Errorf(
			"this Attribute.ValueType does not exists (got=%s)",
			v.Attribute.ValueType,
		))
	}
}

func (v *Value) SetNull() error {
	err := v.CheckWhole()
	if err != nil {
		return err
	}

	if v.Attribute.Required {
		return fmt.Errorf("can't set null value for a required attribute")
	}

	v.IsNull = true
	v.IntVal = 0
	v.FloatVal = 0.0
	v.StringVal = ""
	v.BoolVal = false
	v.RelationVal = uuid.Nil

	return nil
}

func (v *Value) SetStringVal(stringVal string) error {
	err := v.CheckWhole()
	if err != nil {
		return err
	}

	if v.Attribute.ValueType != StringValueType {
		return ErrAskingForWrongType
	}

	v.IsNull = false
	v.StringVal = stringVal

	return nil
}

func (v *Value) SetIntVal(intVal int) error {
	err := v.CheckWhole()
	if err != nil {
		return err
	}

	if v.Attribute.ValueType != IntValueType {
		return ErrAskingForWrongType
	}

	v.IsNull = false
	v.IntVal = intVal

	return nil
}

func (v *Value) SetFloatVal(floatVal float64) error {
	err := v.CheckWhole()
	if err != nil {
		return err
	}

	if v.Attribute.ValueType != FloatValueType {
		return ErrAskingForWrongType
	}

	v.IsNull = false
	v.FloatVal = floatVal

	return nil
}

func (v *Value) SetBooleanVal(boolVal bool) error {
	err := v.CheckWhole()
	if err != nil {
		return err
	}

	if v.Attribute.ValueType != BooleanValueType {
		return ErrAskingForWrongType
	}

	v.IsNull = false
	v.BoolVal = boolVal

	return nil
}

func (v *Value) SetRelationVal(relationVal uuid.UUID) error {
	err := v.CheckWhole()
	if err != nil {
		return err
	}

	if v.Attribute.ValueType != RelationValueType {
		return ErrAskingForWrongType
	}

	v.IsNull = false
	v.RelationVal = relationVal

	return nil
}

func (v Value) Equal(other Value) bool {
	return v.ID == other.ID &&
		v.AttributeID == other.AttributeID &&
		v.EntityID == other.EntityID &&
		v.Value() == other.Value()
}
