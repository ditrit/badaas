// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	basicslices "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/basicslices"
	gorm "gorm.io/gorm"
	"time"
)

func BasicSlicesId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func BasicSlicesCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func BasicSlicesUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func BasicSlicesDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}

var basicSlicesBoolFieldID = badorm.FieldIdentifier{Field: "Bool"}

func BasicSlicesBool(expr badorm.Expression[[]bool]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []bool]{
		Expression:      expr,
		FieldIdentifier: basicSlicesBoolFieldID,
	}
}

var basicSlicesIntFieldID = badorm.FieldIdentifier{Field: "Int"}

func BasicSlicesInt(expr badorm.Expression[[]int]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []int]{
		Expression:      expr,
		FieldIdentifier: basicSlicesIntFieldID,
	}
}

var basicSlicesInt8FieldID = badorm.FieldIdentifier{Field: "Int8"}

func BasicSlicesInt8(expr badorm.Expression[[]int8]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []int8]{
		Expression:      expr,
		FieldIdentifier: basicSlicesInt8FieldID,
	}
}

var basicSlicesInt16FieldID = badorm.FieldIdentifier{Field: "Int16"}

func BasicSlicesInt16(expr badorm.Expression[[]int16]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []int16]{
		Expression:      expr,
		FieldIdentifier: basicSlicesInt16FieldID,
	}
}

var basicSlicesInt32FieldID = badorm.FieldIdentifier{Field: "Int32"}

func BasicSlicesInt32(expr badorm.Expression[[]int32]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []int32]{
		Expression:      expr,
		FieldIdentifier: basicSlicesInt32FieldID,
	}
}

var basicSlicesInt64FieldID = badorm.FieldIdentifier{Field: "Int64"}

func BasicSlicesInt64(expr badorm.Expression[[]int64]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []int64]{
		Expression:      expr,
		FieldIdentifier: basicSlicesInt64FieldID,
	}
}

var basicSlicesUIntFieldID = badorm.FieldIdentifier{Field: "UInt"}

func BasicSlicesUInt(expr badorm.Expression[[]uint]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []uint]{
		Expression:      expr,
		FieldIdentifier: basicSlicesUIntFieldID,
	}
}

var basicSlicesUInt8FieldID = badorm.FieldIdentifier{Field: "UInt8"}

func BasicSlicesUInt8(expr badorm.Expression[[]uint8]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []uint8]{
		Expression:      expr,
		FieldIdentifier: basicSlicesUInt8FieldID,
	}
}

var basicSlicesUInt16FieldID = badorm.FieldIdentifier{Field: "UInt16"}

func BasicSlicesUInt16(expr badorm.Expression[[]uint16]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []uint16]{
		Expression:      expr,
		FieldIdentifier: basicSlicesUInt16FieldID,
	}
}

var basicSlicesUInt32FieldID = badorm.FieldIdentifier{Field: "UInt32"}

func BasicSlicesUInt32(expr badorm.Expression[[]uint32]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []uint32]{
		Expression:      expr,
		FieldIdentifier: basicSlicesUInt32FieldID,
	}
}

var basicSlicesUInt64FieldID = badorm.FieldIdentifier{Field: "UInt64"}

func BasicSlicesUInt64(expr badorm.Expression[[]uint64]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []uint64]{
		Expression:      expr,
		FieldIdentifier: basicSlicesUInt64FieldID,
	}
}

var basicSlicesUIntptrFieldID = badorm.FieldIdentifier{Field: "UIntptr"}

func BasicSlicesUIntptr(expr badorm.Expression[[]uintptr]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []uintptr]{
		Expression:      expr,
		FieldIdentifier: basicSlicesUIntptrFieldID,
	}
}

var basicSlicesFloat32FieldID = badorm.FieldIdentifier{Field: "Float32"}

func BasicSlicesFloat32(expr badorm.Expression[[]float32]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []float32]{
		Expression:      expr,
		FieldIdentifier: basicSlicesFloat32FieldID,
	}
}

var basicSlicesFloat64FieldID = badorm.FieldIdentifier{Field: "Float64"}

func BasicSlicesFloat64(expr badorm.Expression[[]float64]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []float64]{
		Expression:      expr,
		FieldIdentifier: basicSlicesFloat64FieldID,
	}
}

var basicSlicesComplex64FieldID = badorm.FieldIdentifier{Field: "Complex64"}

func BasicSlicesComplex64(expr badorm.Expression[[]complex64]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []complex64]{
		Expression:      expr,
		FieldIdentifier: basicSlicesComplex64FieldID,
	}
}

var basicSlicesComplex128FieldID = badorm.FieldIdentifier{Field: "Complex128"}

func BasicSlicesComplex128(expr badorm.Expression[[]complex128]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []complex128]{
		Expression:      expr,
		FieldIdentifier: basicSlicesComplex128FieldID,
	}
}

var basicSlicesStringFieldID = badorm.FieldIdentifier{Field: "String"}

func BasicSlicesString(expr badorm.Expression[[]string]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []string]{
		Expression:      expr,
		FieldIdentifier: basicSlicesStringFieldID,
	}
}

var basicSlicesByteFieldID = badorm.FieldIdentifier{Field: "Byte"}

func BasicSlicesByte(expr badorm.Expression[[]uint8]) badorm.WhereCondition[basicslices.BasicSlices] {
	return badorm.FieldCondition[basicslices.BasicSlices, []uint8]{
		Expression:      expr,
		FieldIdentifier: basicSlicesByteFieldID,
	}
}

var BasicSlicesPreloadAttributes = badorm.NewPreloadCondition[basicslices.BasicSlices](basicSlicesBoolFieldID, basicSlicesIntFieldID, basicSlicesInt8FieldID, basicSlicesInt16FieldID, basicSlicesInt32FieldID, basicSlicesInt64FieldID, basicSlicesUIntFieldID, basicSlicesUInt8FieldID, basicSlicesUInt16FieldID, basicSlicesUInt32FieldID, basicSlicesUInt64FieldID, basicSlicesUIntptrFieldID, basicSlicesFloat32FieldID, basicSlicesFloat64FieldID, basicSlicesComplex64FieldID, basicSlicesComplex128FieldID, basicSlicesStringFieldID, basicSlicesByteFieldID)
var BasicSlicesPreloadRelations = []badorm.Condition[basicslices.BasicSlices]{}
