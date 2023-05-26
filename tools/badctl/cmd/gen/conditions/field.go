package conditions

import (
	"errors"
	"go/types"

	"github.com/ettle/strcase"
)

type Field struct {
	Name     string
	Type     Type
	Embedded bool
	Tags     GormTags
}

// Get the name of the column where the data for a field will be saved
func (field Field) getColumnName() string {
	columnTag, isPresent := field.Tags[columnTagName]
	if isPresent {
		// field has a gorm column tag, so the name of the column will be that tag
		return columnTag
	}

	// column name generated automatically by gorm
	// TODO support https://gorm.io/docs/conventions.html#NamingStrategy
	return strcase.ToSnake(field.Name)
}

// Get name of the attribute of the object that is a foreign key to the field's object
func (field Field) getFKAttribute() string {
	foreignKeyTag, isPresent := field.Tags[foreignKeyTagName]
	if isPresent {
		// field has a foreign key tag, so the name will be that tag
		return foreignKeyTag
	}

	// gorm default
	return field.Name + "ID"
}

// Get name of the attribute of the field's object that is references by the foreign key
func (field Field) getFKReferencesAttribute() string {
	referencesTag, isPresent := field.Tags[referencesTagName]
	// TODO testear cuando hay redefinicion en la inversa
	if isPresent {
		// field has a references tag, so the name will be that tag
		return referencesTag
	}

	// gorm default
	return "ID"
}

// Get name of the attribute of field's object that is a foreign key to the object
func (field Field) getRelatedTypeFKAttribute(structName string) string {
	// TODO testear cuando hay redefinicion
	foreignKeyTag, isPresent := field.Tags[foreignKeyTagName]
	if isPresent {
		// field has a foreign key tag, so the name will that tag
		return foreignKeyTag
	}

	// gorm default
	return structName + "ID"
}

// Get field's type full string (pkg + name)
func (field Field) TypeString() string {
	return field.Type.String()
}

// Get field's type name
func (field Field) TypeName() string {
	return field.Type.Name()
}

// Create a new field with the same name and tags but a different type
func (field Field) ChangeType(newType types.Type) Field {
	return Field{
		Name: field.Name,
		Type: Type{newType},
		Tags: field.Tags,
	}
}

// Get fields of a BaDORM model
// Returns error is objectType is not a BaDORM model
func getFields(objectType Type, prefix string) ([]Field, error) {
	// The underlying type has to be a struct and a BaDORM Model
	// (ignore const, var, func, etc.)
	structType, err := objectType.BadORMModelStruct()
	if err != nil {
		return nil, err
	}

	return getStructFields(structType, prefix)
}

// Get fields of a struct
// Returns errors if the struct has not fields
func getStructFields(structType *types.Struct, prefix string) ([]Field, error) {
	numFields := structType.NumFields()
	if numFields == 0 {
		return nil, errors.New("struct has 0 fields")
	}

	fields := []Field{}

	// Iterate over struct fields
	for i := 0; i < numFields; i++ {
		fieldObject := structType.Field(i)
		gormTags := getGormTags(structType.Tag(i))
		fields = append(fields, Field{
			Name:     prefix + fieldObject.Name(),
			Type:     Type{fieldObject.Type()},
			Embedded: fieldObject.Embedded() || gormTags.hasEmbedded(),
			Tags:     gormTags,
		})
	}

	return fields, nil
}
