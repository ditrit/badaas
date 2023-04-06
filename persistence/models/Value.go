package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Describe the attribute value of an Entity
type Value struct {
	BaseModel
	IsNull      bool
	StringVal   string
	FloatVal    float64
	IntVal      int
	BoolVal     bool
	RelationVal uuid.UUID

	StringifiedValue string

	// GORM relations
	EntityId    uuid.UUID
	AttributeId uuid.UUID
	Attribute   *Attribute
}

var (
	ErrValueIsNull        = errors.New("you can't get the value from a null Value")
	ErrAskingForWrongType = errors.New("you can't get this type of value, the attribute type doesn't match")
)

// Create a new null value
func NewNullValue(attr *Attribute) (*Value, error) {
	val := new(Value)
	if attr.Required {
		return nil, fmt.Errorf("can't create new null value for a required attribute")
	}
	val.IsNull = true
	val.Attribute = attr
	return val, nil
}

// Create a new int value
func NewIntValue(attr *Attribute, i int) (*Value, error) {
	val := new(Value)
	if attr.ValueType != IntValueType {
		return nil, fmt.Errorf("can't create a new int value with a %s attribute", attr.ValueType)
	}
	val.IsNull = false
	val.IntVal = i
	val.Attribute = attr
	return val, nil
}

// Create a new bool value
func NewBoolValue(attr *Attribute, b bool) (*Value, error) {
	val := new(Value)
	if attr.ValueType != BooleanValueType {
		return nil, fmt.Errorf("can't create a new bool value with a %s attribute", attr.ValueType)
	}
	val.IsNull = false
	val.BoolVal = b
	val.Attribute = attr
	return val, nil
}

// Create a new float value
func NewFloatValue(attr *Attribute, f float64) (*Value, error) {
	val := new(Value)
	if attr.ValueType != FloatValueType {
		return nil, fmt.Errorf("can't create a new float value with a %s attribute", attr.ValueType)
	}
	val.IsNull = false
	val.FloatVal = f
	val.Attribute = attr
	return val, nil
}

// Create a new string value
func NewStringValue(attr *Attribute, s string) (*Value, error) {
	val := new(Value)
	if attr.ValueType != StringValueType {
		return nil, fmt.Errorf("can't create a new string value with a %s attribute", attr.ValueType)
	}
	val.IsNull = false
	val.StringVal = s
	val.Attribute = attr
	return val, nil
}

// Create a new relation value.
// If et is nil, then the function return an error
// If et is of the wrong types
func NewRelationValue(attr *Attribute, et *Entity) (*Value, error) {
	val := new(Value)
	if attr.ValueType != RelationValueType {
		return nil, fmt.Errorf("can't create a new relation value with a %s attribute", attr.ValueType)
	}
	if et == nil {
		return nil, fmt.Errorf("can't create a new relation with a nill entity pointer")
	}
	if et.EntityType.ID != attr.TargetEntityTypeId {
		return nil, fmt.Errorf(
			"can't create a relation with an entity of wrong EntityType. (got the entityid=%d, expected=%d)",
			et.EntityType.ID, attr.TargetEntityTypeId,
		)
	}
	val.IsNull = false
	val.RelationVal = et.ID
	val.Attribute = attr
	return val, nil
}

// Check if the Value is whole. eg, no fields are nil
func (v *Value) CheckWhole() error {
	if v.Attribute == nil {
		return fmt.Errorf("the Attribute pointer is nil in Value at %v", v)
	}
	return nil
}

// Return the string value
// If the Value is null, it return ErrValueIsNull
// If the Value not of the requested type, it return ErrAskingForWrongType
// If the Value.Attribute == nil, it panic
func (v *Value) GetStringVal() (string, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}

	if v.IsNull {
		if v.Attribute.Default {
			return v.Attribute.DefaultString, nil
		}
		return "", ErrValueIsNull
	}
	if v.Attribute.ValueType != StringValueType {
		return "", ErrAskingForWrongType
	}
	return v.StringVal, nil
}

// Return the float value
// If the Value is null, it return ErrValueIsNull
// If the Value not of the requested type, it return ErrAskingForWrongType
// If the Value.Attribute == nil, it panic
func (v *Value) GetFloatVal() (float64, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}

	if v.IsNull {
		if v.Attribute.Default {
			return v.Attribute.DefaultFloat, nil
		}
		return 0.0, ErrValueIsNull
	}
	if v.Attribute.ValueType != FloatValueType {
		return 0.0, ErrAskingForWrongType
	}
	return v.FloatVal, nil
}

// Return the int value
// If the Value is null, it return ErrValueIsNull
// If the Value not of the requested type, it return ErrAskingForWrongType
// If the Value.Attribute == nil, it panic
func (v *Value) GetIntVal() (int, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}

	if v.IsNull {
		if v.Attribute.Default {
			return v.Attribute.DefaultInt, nil
		}
		return 0, ErrValueIsNull
	}
	if v.Attribute.ValueType != IntValueType {
		return 0, ErrAskingForWrongType
	}
	return v.IntVal, nil
}

// Return the bool value
// If the Value is null, it return ErrValueIsNull
// If the Value not of the requested type, it return ErrAskingForWrongType
// If the Value.Attribute == nil, it panic
func (v *Value) GetBoolVal() (bool, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}

	if v.IsNull {
		if v.Attribute.Default {
			return v.Attribute.DefaultBool, nil
		}
		return false, ErrValueIsNull
	}
	if v.Attribute.ValueType != BooleanValueType {
		return false, ErrAskingForWrongType
	}
	return v.BoolVal, nil
}

// Return the Relation value as a *Entity
// If the Value is null, it return nil
// If the Value not of the requested type, it return ErrAskingForWrongType
// If the Value.Attribute == nil, it panic
func (v *Value) GetComputedRelationVal(db *gorm.DB) (*Entity, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}
	if v.Attribute.ValueType != RelationValueType {
		return nil, ErrAskingForWrongType
	}

	if v.IsNull {
		return nil, nil
	}
	var et Entity
	err = db.Preload("Fields").Preload("Fields.attribute").Preload("EntityType.Attributes").Preload("EntityType").First(&et, v.RelationVal).Error
	if err != nil {
		return nil, err
	}
	return &et, nil
}

// Return the Relation value as an uint (returns the ID)
// If the Value is null, it return nil
// If the Value not of the requested type, it return ErrAskingForWrongType
// If the Value.Attribute == nil, it panic
func (v *Value) GetRelationVal() (uuid.UUID, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}
	if v.Attribute.ValueType != RelationValueType {
		return uuid.Nil, ErrAskingForWrongType
	}

	if v.IsNull {
		return uuid.Nil, fmt.Errorf("the relation is null")
	}
	return v.RelationVal, nil
}

// Return the underlying value as an interface
func (v *Value) Value() any {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}
	switch v.Attribute.ValueType {
	case StringValueType:
		if v.IsNull {
			if v.Attribute.Default {
				return v.Attribute.DefaultString
			}
			return nil
		}
		return v.StringVal
	case IntValueType:
		if v.IsNull {
			if v.Attribute.Default {
				return v.Attribute.DefaultInt
			}
			return nil
		}
		return v.IntVal
	case FloatValueType:
		if v.IsNull {
			if v.Attribute.Default {
				return v.Attribute.DefaultFloat
			}
			return nil
		}
		return v.FloatVal
	case BooleanValueType:
		if v.IsNull {
			if v.Attribute.Default {
				return v.Attribute.DefaultBool
			}
			return nil
		}
		return v.BoolVal
	case RelationValueType:
		if v.IsNull {
			return nil
		}
		return v.RelationVal
	default:
		panic(fmt.Errorf(
			"this Attribute.ValueType does not exists (got=%s)",
			v.Attribute.ValueType,
		))
	}
}

// When Value isNull, it is impossible to build a Key/Value pair
var ErrCantBuildKVPairForNullValue = errors.New("can't build key/value pair from null value")

// Build a key/value pair to be included in a JSON
// If the value hold an int=8 with an attribute named "voila" then the string returned will be `"voila":8`
func (v *Value) BuildJsonKVPair() (string, error) {
	err := v.CheckWhole()
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(v.Value())
	if err != nil {
		return "", fmt.Errorf("an error happened while trying to marshall the %q attr: (%w)", v.Attribute.Name, err)
	}
	return fmt.Sprintf("%q:%s", v.Attribute.Name, bytes), nil
}
