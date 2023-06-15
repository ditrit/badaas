package conditions

import (
	"errors"
	"go/types"

	"github.com/dave/jennifer/jen"

	"github.com/ditrit/badaas/tools/badctl/cmd/cmderrors"
)

type JenParam struct {
	statement    *jen.Statement
	internalType *jen.Statement
}

func NewJenParam() *JenParam {
	return &JenParam{
		statement: jen.Id("exprs").Op("...").Qual(
			badORMPath, badORMExpression,
		),
		internalType: &jen.Statement{},
	}
}

func (param JenParam) Statement() *jen.Statement {
	return param.statement.Types(param.internalType)
}

func (param JenParam) GenericType() *jen.Statement {
	return param.internalType
}

func (param JenParam) ToBasicKind(basicType *types.Basic) {
	switch basicType.Kind() {
	case types.Bool:
		param.internalType.Bool()
	case types.Int:
		param.internalType.Int()
	case types.Int8:
		param.internalType.Int8()
	case types.Int16:
		param.internalType.Int16()
	case types.Int32:
		param.internalType.Int32()
	case types.Int64:
		param.internalType.Int64()
	case types.Uint:
		param.internalType.Uint()
	case types.Uint8:
		param.internalType.Uint8()
	case types.Uint16:
		param.internalType.Uint16()
	case types.Uint32:
		param.internalType.Uint32()
	case types.Uint64:
		param.internalType.Uint64()
	case types.Uintptr:
		param.internalType.Uintptr()
	case types.Float32:
		param.internalType.Float32()
	case types.Float64:
		param.internalType.Float64()
	case types.Complex64:
		param.internalType.Complex64()
	case types.Complex128:
		param.internalType.Complex128()
	case types.String:
		param.internalType.String()
	case types.Invalid, types.UnsafePointer,
		types.UntypedBool, types.UntypedInt,
		types.UntypedRune, types.UntypedFloat,
		types.UntypedComplex, types.UntypedString,
		types.UntypedNil:
		cmderrors.FailErr(errors.New("unreachable! untyped types can't be inside a struct"))
	}
}

func (param JenParam) ToPointer() {
	param.internalType.Op("*")
}

func (param JenParam) ToSlice() {
	param.internalType.Index()
}

func (param JenParam) ToCustomType(destPkg string, typeV Type) {
	param.internalType.Qual(
		getRelativePackagePath(destPkg, typeV),
		typeV.Name(),
	)
}