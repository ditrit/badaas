// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	basicpointers "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/basicpointers"
	gorm "gorm.io/gorm"
	"time"
)

func BasicPointersId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func BasicPointersCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func BasicPointersUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func BasicPointersDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}

var basicPointersBoolFieldID = badorm.FieldIdentifier{Field: "Bool"}

func BasicPointersBool(expr badorm.Expression[bool]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, bool]{
		Expression:      expr,
		FieldIdentifier: basicPointersBoolFieldID,
	}
}

var basicPointersIntFieldID = badorm.FieldIdentifier{Field: "Int"}

func BasicPointersInt(expr badorm.Expression[int]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, int]{
		Expression:      expr,
		FieldIdentifier: basicPointersIntFieldID,
	}
}

var basicPointersInt8FieldID = badorm.FieldIdentifier{Field: "Int8"}

func BasicPointersInt8(expr badorm.Expression[int8]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, int8]{
		Expression:      expr,
		FieldIdentifier: basicPointersInt8FieldID,
	}
}

var basicPointersInt16FieldID = badorm.FieldIdentifier{Field: "Int16"}

func BasicPointersInt16(expr badorm.Expression[int16]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, int16]{
		Expression:      expr,
		FieldIdentifier: basicPointersInt16FieldID,
	}
}

var basicPointersInt32FieldID = badorm.FieldIdentifier{Field: "Int32"}

func BasicPointersInt32(expr badorm.Expression[int32]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, int32]{
		Expression:      expr,
		FieldIdentifier: basicPointersInt32FieldID,
	}
}

var basicPointersInt64FieldID = badorm.FieldIdentifier{Field: "Int64"}

func BasicPointersInt64(expr badorm.Expression[int64]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, int64]{
		Expression:      expr,
		FieldIdentifier: basicPointersInt64FieldID,
	}
}

var basicPointersUIntFieldID = badorm.FieldIdentifier{Field: "UInt"}

func BasicPointersUInt(expr badorm.Expression[uint]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, uint]{
		Expression:      expr,
		FieldIdentifier: basicPointersUIntFieldID,
	}
}

var basicPointersUInt8FieldID = badorm.FieldIdentifier{Field: "UInt8"}

func BasicPointersUInt8(expr badorm.Expression[uint8]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, uint8]{
		Expression:      expr,
		FieldIdentifier: basicPointersUInt8FieldID,
	}
}

var basicPointersUInt16FieldID = badorm.FieldIdentifier{Field: "UInt16"}

func BasicPointersUInt16(expr badorm.Expression[uint16]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, uint16]{
		Expression:      expr,
		FieldIdentifier: basicPointersUInt16FieldID,
	}
}

var basicPointersUInt32FieldID = badorm.FieldIdentifier{Field: "UInt32"}

func BasicPointersUInt32(expr badorm.Expression[uint32]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, uint32]{
		Expression:      expr,
		FieldIdentifier: basicPointersUInt32FieldID,
	}
}

var basicPointersUInt64FieldID = badorm.FieldIdentifier{Field: "UInt64"}

func BasicPointersUInt64(expr badorm.Expression[uint64]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, uint64]{
		Expression:      expr,
		FieldIdentifier: basicPointersUInt64FieldID,
	}
}

var basicPointersUIntptrFieldID = badorm.FieldIdentifier{Field: "UIntptr"}

func BasicPointersUIntptr(expr badorm.Expression[uintptr]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, uintptr]{
		Expression:      expr,
		FieldIdentifier: basicPointersUIntptrFieldID,
	}
}

var basicPointersFloat32FieldID = badorm.FieldIdentifier{Field: "Float32"}

func BasicPointersFloat32(expr badorm.Expression[float32]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, float32]{
		Expression:      expr,
		FieldIdentifier: basicPointersFloat32FieldID,
	}
}

var basicPointersFloat64FieldID = badorm.FieldIdentifier{Field: "Float64"}

func BasicPointersFloat64(expr badorm.Expression[float64]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, float64]{
		Expression:      expr,
		FieldIdentifier: basicPointersFloat64FieldID,
	}
}

var basicPointersComplex64FieldID = badorm.FieldIdentifier{Field: "Complex64"}

func BasicPointersComplex64(expr badorm.Expression[complex64]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, complex64]{
		Expression:      expr,
		FieldIdentifier: basicPointersComplex64FieldID,
	}
}

var basicPointersComplex128FieldID = badorm.FieldIdentifier{Field: "Complex128"}

func BasicPointersComplex128(expr badorm.Expression[complex128]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, complex128]{
		Expression:      expr,
		FieldIdentifier: basicPointersComplex128FieldID,
	}
}

var basicPointersStringFieldID = badorm.FieldIdentifier{Field: "String"}

func BasicPointersString(expr badorm.Expression[string]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, string]{
		Expression:      expr,
		FieldIdentifier: basicPointersStringFieldID,
	}
}

var basicPointersByteFieldID = badorm.FieldIdentifier{Field: "Byte"}

func BasicPointersByte(expr badorm.Expression[uint8]) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.FieldCondition[basicpointers.BasicPointers, uint8]{
		Expression:      expr,
		FieldIdentifier: basicPointersByteFieldID,
	}
}

var BasicPointersPreloadAttributes = badorm.NewPreloadCondition[basicpointers.BasicPointers](basicPointersBoolFieldID, basicPointersIntFieldID, basicPointersInt8FieldID, basicPointersInt16FieldID, basicPointersInt32FieldID, basicPointersInt64FieldID, basicPointersUIntFieldID, basicPointersUInt8FieldID, basicPointersUInt16FieldID, basicPointersUInt32FieldID, basicPointersUInt64FieldID, basicPointersUIntptrFieldID, basicPointersFloat32FieldID, basicPointersFloat64FieldID, basicPointersComplex64FieldID, basicPointersComplex128FieldID, basicPointersStringFieldID, basicPointersByteFieldID)
var BasicPointersPreloadRelations = []badorm.Condition[basicpointers.BasicPointers]{}
