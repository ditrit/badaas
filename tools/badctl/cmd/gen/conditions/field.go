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

func (field Field) getColumnName() string {
	columnTag, isPresent := field.Tags[columnTagName]
	if isPresent {
		return columnTag
	}

	return strcase.ToSnake(field.Name)
}

func (field Field) getFKAttribute() string {
	foreignKeyTag, isPresent := field.Tags[foreignKeyTagName]
	if isPresent {
		return foreignKeyTag
	}

	return field.Name + "ID"
}

func (field Field) getFKReferencesAttribute() string {
	referencesTag, isPresent := field.Tags[referencesTagName]
	// TODO testear cuando hay redefinicion en la inversa
	if isPresent {
		return referencesTag
	}

	return "ID"
}

func (field Field) getRelatedTypeFKAttribute(structName string) string {
	// TODO testear cuando hay redefinicion
	foreignKeyTag, isPresent := field.Tags[foreignKeyTagName]
	if isPresent {
		return foreignKeyTag
	}

	return structName + "ID"
}

func (field Field) TypeString() string {
	return field.Type.String()
}

func (field Field) TypeName() string {
	return field.Type.Name()
}

func (field Field) ChangeType(newType types.Type) Field {
	return Field{
		Name: field.Name,
		Type: Type{newType},
		Tags: field.Tags,
	}
}

func getFields(objectType Type, prefix string) ([]Field, error) {
	// The underlying type has to be a struct and a BaDORM Model
	// (ignore const, var, func, etc.)
	structType, err := objectType.BadORMModelStruct()
	if err != nil {
		return nil, err
	}

	return getStructFields(structType, prefix)
}

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
