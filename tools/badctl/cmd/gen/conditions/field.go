package conditions

import (
	"errors"
	"go/types"
	"regexp"

	"github.com/ettle/strcase"
)

type Field struct {
	Name     string
	Type     types.Type
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
	return getTypeName(field.Type)
}

func (field Field) TypePkg() *types.Package {
	return getTypePkg(field.Type)
}

func (field Field) ChangeType(newType types.Type) Field {
	return Field{
		Name: field.Name,
		Type: newType,
		Tags: field.Tags,
	}
}

var scanMethod = regexp.MustCompile(`func \(\*.*\)\.Scan\([a-zA-Z0-9_-]* interface\{\}\) error$`)
var valueMethod = regexp.MustCompile(`func \(.*\)\.Value\(\) \(database/sql/driver\.Value\, error\)$`)

func (field Field) IsGormCustomType() bool {
	typeNamed, isNamedType := field.Type.(*types.Named)
	if !isNamedType {
		return false
	}

	hasScanMethod := false
	hasValueMethod := false
	for i := 0; i < typeNamed.NumMethods(); i++ {
		methodSignature := typeNamed.Method(i).String()

		if !hasScanMethod && scanMethod.MatchString(methodSignature) {
			hasScanMethod = true
		} else if !hasValueMethod && valueMethod.MatchString(methodSignature) {
			hasValueMethod = true
		}
	}

	return hasScanMethod && hasValueMethod
}

func getFields(objectType types.Type, prefix string) ([]Field, error) {
	// The underlying type has to be a struct and a BaDORM Model
	// (ignore const, var, func, etc.)
	structType, err := getBadORMModelStruct(objectType)
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
			Type:     fieldObject.Type(),
			Embedded: fieldObject.Embedded() || gormTags.hasEmbedded(),
			Tags:     gormTags,
		})
	}

	return fields, nil
}
