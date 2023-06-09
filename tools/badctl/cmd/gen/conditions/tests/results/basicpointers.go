// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	basicpointers "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/basicpointers"
	gorm "gorm.io/gorm"
	"time"
)

func BasicPointersId(v badorm.UUID) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "ID",
		Value: v,
	}
}
func BasicPointersCreatedAt(v time.Time) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "CreatedAt",
		Value: v,
	}
}
func BasicPointersUpdatedAt(v time.Time) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "UpdatedAt",
		Value: v,
	}
}
func BasicPointersDeletedAt(v gorm.DeletedAt) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "DeletedAt",
		Value: v,
	}
}
func BasicPointersBool(v *bool) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "Bool",
		Value: v,
	}
}
func BasicPointersInt(v *int) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "Int",
		Value: v,
	}
}
func BasicPointersInt8(v *int8) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "Int8",
		Value: v,
	}
}
func BasicPointersInt16(v *int16) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "Int16",
		Value: v,
	}
}
func BasicPointersInt32(v *int32) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "Int32",
		Value: v,
	}
}
func BasicPointersInt64(v *int64) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "Int64",
		Value: v,
	}
}
func BasicPointersUInt(v *uint) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "UInt",
		Value: v,
	}
}
func BasicPointersUInt8(v *uint8) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "UInt8",
		Value: v,
	}
}
func BasicPointersUInt16(v *uint16) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "UInt16",
		Value: v,
	}
}
func BasicPointersUInt32(v *uint32) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "UInt32",
		Value: v,
	}
}
func BasicPointersUInt64(v *uint64) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "UInt64",
		Value: v,
	}
}
func BasicPointersUIntptr(v *uintptr) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "UIntptr",
		Value: v,
	}
}
func BasicPointersFloat32(v *float32) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "Float32",
		Value: v,
	}
}
func BasicPointersFloat64(v *float64) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "Float64",
		Value: v,
	}
}
func BasicPointersComplex64(v *complex64) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "Complex64",
		Value: v,
	}
}
func BasicPointersComplex128(v *complex128) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "Complex128",
		Value: v,
	}
}
func BasicPointersString(v *string) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "String",
		Value: v,
	}
}
func BasicPointersByte(v *uint8) badorm.WhereCondition[basicpointers.BasicPointers] {
	return badorm.WhereCondition[basicpointers.BasicPointers]{
		Field: "Byte",
		Value: v,
	}
}