package models

import (
	"encoding/json"
	"errors"
	"fmt"

	uuid "github.com/google/uuid"

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
	RelationVal uuid.UUID `gorm:"type:uuid;foreignKey:Entity;index:fk_relation_val_entity"`

	// GORM relations
	EntityID    uuid.UUID
	AttributeID uuid.UUID
	Attribute   *Attribute
}

var (
	ErrValueIsNull        = errors.New("you can't get the value from a null Value")
	ErrAskingForWrongType = errors.New("you can't get this type of value, the attribute type doesn't match")
)

// Create a new null value
func NewNullValue(attr *Attribute) (*Value, error) {
	if attr.Required {
		return nil, fmt.Errorf("can't create new null value for a required attribute")
	}

	return &Value{IsNull: true, Attribute: attr, AttributeID: attr.ID}, nil
}

// Create a new int value
func NewIntValue(attr *Attribute, i int) (*Value, error) {
	if attr.ValueType != IntValueType {
		return nil, fmt.Errorf("can't create a new int value with a %s attribute", attr.ValueType)
	}

	return &Value{IntVal: i, Attribute: attr, AttributeID: attr.ID}, nil
}

// Create a new bool value
func NewBoolValue(attr *Attribute, b bool) (*Value, error) {
	if attr.ValueType != BooleanValueType {
		return nil, fmt.Errorf("can't create a new bool value with a %s attribute", attr.ValueType)
	}

	return &Value{BoolVal: b, Attribute: attr, AttributeID: attr.ID}, nil
}

// Create a new float value
func NewFloatValue(attr *Attribute, f float64) (*Value, error) {
	if attr.ValueType != FloatValueType {
		return nil, fmt.Errorf("can't create a new float value with a %s attribute", attr.ValueType)
	}

	return &Value{FloatVal: f, Attribute: attr, AttributeID: attr.ID}, nil
}

// Create a new string value
func NewStringValue(attr *Attribute, s string) (*Value, error) {
	if attr.ValueType != StringValueType {
		return nil, fmt.Errorf("can't create a new string value with a %s attribute", attr.ValueType)
	}

	return &Value{StringVal: s, Attribute: attr, AttributeID: attr.ID}, nil
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
			"can't create a relation with an entity of wrong EntityType. (got the entityid=%d, expected=%d)",
			et.EntityType.ID, attr.RelationTargetEntityTypeID,
		)
	}

	return NewRelationIDValue(attr, et.ID)
}

// Create a new relation value.
// If et is nil, then the function returns an error
// If et is of the wrong types, then the function returns an error
func NewRelationIDValue(attr *Attribute, uuidVal uuid.UUID) (*Value, error) {
	if attr.ValueType != RelationValueType {
		return nil, fmt.Errorf("can't create a new relation value with a %s attribute", attr.ValueType)
	}

	return &Value{RelationVal: uuidVal, Attribute: attr, AttributeID: attr.ID}, nil
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
		return "", err
	}

	if v.Attribute.ValueType != StringValueType {
		return "", ErrAskingForWrongType
	}

	if v.IsNull {
		if v.Attribute.Default {
			return v.Attribute.DefaultString, nil
		}
		return "", ErrValueIsNull
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
		return 0.0, err
	}

	if v.Attribute.ValueType != FloatValueType {
		return 0.0, ErrAskingForWrongType
	}

	if v.IsNull {
		if v.Attribute.Default {
			return v.Attribute.DefaultFloat, nil
		}
		return 0.0, ErrValueIsNull
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
		return 0, err
	}

	if v.Attribute.ValueType != IntValueType {
		return 0, ErrAskingForWrongType
	}

	if v.IsNull {
		if v.Attribute.Default {
			return v.Attribute.DefaultInt, nil
		}
		return 0, ErrValueIsNull
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
		return false, err
	}

	if v.Attribute.ValueType != BooleanValueType {
		return false, ErrAskingForWrongType
	}

	if v.IsNull {
		if v.Attribute.Default {
			return v.Attribute.DefaultBool, nil
		}
		return false, ErrValueIsNull
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
		return nil, err
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
		return uuid.Nil, err
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

func (v *Value) SetNull() {
	v.IsNull = true
	v.IntVal = 0
	v.FloatVal = 0.0
	v.StringVal = ""
	v.BoolVal = false
	v.RelationVal = uuid.Nil
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

var ErrCantBuildKVPairForNullValue = errors.New("can't build key/value pair from null value") // When Value isNull, it is impossible to build a Key/Value pair

// Build a key/value pair to be included in a JSON
// If the value hold an int=8 with an attribute named "voila" then the string returned will be `"voila":8`
func (v *Value) BuildJSONKVPair() (string, error) {
	err := v.CheckWhole()
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(v.Value())
	if err != nil {
		return "", fmt.Errorf("an error happened while trying to marshall the %q attr: (%w)", v.Attribute.Name, err)
	}

	return fmt.Sprintf("%q:%s", v.Attribute.Name, bytes), nil
}

func (v Value) Equal(other Value) bool {
	return v.ID == other.ID &&
		v.AttributeID == other.AttributeID &&
		v.EntityID == other.EntityID &&
		v.Value() == other.Value()
}
