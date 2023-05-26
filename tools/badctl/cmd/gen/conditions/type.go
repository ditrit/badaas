package conditions

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/elliotchance/pie/v2"
)

// TODO me gustaria que esten en una clase
func getTypeName(typeV types.Type) string {
	switch typeTyped := typeV.(type) {
	case *types.Named:
		return typeTyped.Obj().Name()
	default:
		return pie.Last(strings.Split(typeV.String(), "."))
	}
}

func getTypePkg(typeV types.Type) *types.Package {
	switch typeTyped := typeV.(type) {
	case *types.Named:
		return typeTyped.Obj().Pkg()
	default:
		return nil
	}
}

func getBadORMModelStruct(typeV types.Type) (*types.Struct, error) {
	structType, ok := typeV.Underlying().(*types.Struct)
	if !ok || !isBadORMModel(structType) {
		return nil, fmt.Errorf("type %s is not a BaDORM Model", typeV.String())
	}

	return structType, nil
}

func isBadORMModel(structType *types.Struct) bool {
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)

		if field.Embedded() && pie.Contains(badORMModels, field.Name()) {
			return true
		}
	}

	return false
}

func hasFK(typeV types.Type, field Field) (bool, error) {
	objectFields, err := getFields(
		typeV,
		// TODO testear esto si esta bien aca
		field.Tags.getEmbeddedPrefix(),
	)
	if err != nil {
		return false, err
	}
	return pie.Any(objectFields, func(otherField Field) bool {
		return otherField.Name == field.getFKAttribute()
	}), nil
}
