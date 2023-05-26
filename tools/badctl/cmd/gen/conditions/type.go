package conditions

import (
	"go/types"
	"strings"

	"github.com/elliotchance/pie/v2"
)

func getTypeName(typeV types.Type) string {
	switch typeTyped := typeV.(type) {
	case *types.Named:
		return typeTyped.Obj().Name()
	// TODO ver el resto si al hacerlo me simplificaria algo
	default:
		return pie.Last(strings.Split(typeV.String(), "."))
	}
}

func getTypePkg(typeV types.Type) *types.Package {
	switch typeTyped := typeV.(type) {
	case *types.Named:
		return typeTyped.Obj().Pkg()
	// TODO ver el resto si al hacerlo me simplificaria algo
	default:
		return nil
	}
}
