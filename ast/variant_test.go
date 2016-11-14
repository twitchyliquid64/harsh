package ast

import "testing"

func TestMakeVariantString(t *testing.T) {
	v := MakeVariant("abc")
	if v.Type != PrimitiveTypeString {
		t.Error("Expected string type")
	}
	if v.String != "abc" {
		t.Error("Expected \"abc\"")
	}
}

func TestMakeVariantInt(t *testing.T) {
	v := MakeVariant(123)
	if v.Type != PrimitiveTypeInt {
		t.Error("Expected int type")
	}
	if v.Int != 123 {
		t.Error("Expected 123")
	}
}

func TestMakeVariantBool(t *testing.T) {
	v := MakeVariant(true)
	if v.Type != PrimitiveTypeBool {
		t.Error("Expected bool type")
	}
	if v.Bool != true {
		t.Error("Expected true")
	}
}

func TestMakeVariantInt64(t *testing.T) {
	v := MakeVariant(int64(-53))
	if v.Type != PrimitiveTypeInt {
		t.Error("Expected int type")
	}
	if v.Int != -53 {
		t.Error("Expected 123")
	}
}

type WierdUnknownTypeMock int

func TestMakeVariantUndefined(t *testing.T) {
	v := MakeVariant(WierdUnknownTypeMock(1))
	if v.Type != PrimitiveTypeUndefined {
		t.Error("Expected undefined type")
	}
}

func TestDefaultVariantValueInt(t *testing.T) {
	v, err := DefaultVariantValue(PrimitiveTypeInt)
	if v.Type != PrimitiveTypeInt {
		t.Error("Expected int type")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDefaultVariantValueString(t *testing.T) {
	v, err := DefaultVariantValue(PrimitiveTypeString)
	if v.Type != PrimitiveTypeString {
		t.Error("Expected string type")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDefaultVariantValueBool(t *testing.T) {
	v, err := DefaultVariantValue(PrimitiveTypeBool)
	if v.Type != PrimitiveTypeBool {
		t.Error("Expected bool type")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDefaultVariantValueUndefined(t *testing.T) {
	v, err := DefaultVariantValue(PrimitiveTypeUndefined)
	if v.Type != PrimitiveTypeUndefined {
		t.Error("Expected undefined type")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDefaultVariantValueBasicArray(t *testing.T) {
	at := ArrayType{
		SubType: PrimitiveTypeInt,
		Len: &IntegerLiteral{
			Val: 20,
		},
	}
	v, err := DefaultVariantValue(at)
	if v.Type.Kind() != ComplexTypeArray {
		t.Error("Expected array type")
	}
	if len(v.VectorData) != 20 {
		t.Error("Expected len 20")
	}
	if v.VectorData[4].Type != PrimitiveTypeInt {
		t.Error("Expected value at index to be type int")
	}
	if v.VectorData[4].Int != 0 {
		t.Error("Expected value at index to be default (0)")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDefaultVariantValueArrayErrorsIfNotIntegerLen(t *testing.T) {
	at := ArrayType{
		SubType: PrimitiveTypeInt,
		Len:     &StringLiteral{},
	}
	_, err := DefaultVariantValue(at)
	if err == nil {
		t.Error("Expected error")
	}
	if err.Error() != "Resolved length of array was not an integer" {
		t.Error("Expected message \"Resolved length of array was not an integer\"")
	}
}

func TestDefaultVariantValueArrayErrorsIfLenNotResolvable(t *testing.T) {
	at := ArrayType{
		SubType: PrimitiveTypeInt,
		Len: &VariableReference{
			Name: "br",
		},
	}
	_, err := DefaultVariantValue(at)
	if err == nil {
		t.Error("Expected error")
		t.FailNow()
	}
	if err.Error() != "Could not statically resolve the length of the given array" {
		t.Error("Unexpected error, got", err)
	}
}

func TestDefaultVariantValueComplexArray(t *testing.T) {
	at := ArrayType{
		SubType: ArrayType{
			SubType: PrimitiveTypeBool,
			Len: &IntegerLiteral{
				Val: 2,
			},
		},
		Len: &IntegerLiteral{
			Val: 20,
		},
	}
	v, err := DefaultVariantValue(at)
	if v.Type.Kind() != ComplexTypeArray {
		t.Error("Expected array type")
	}
	if len(v.VectorData) != 20 {
		t.Error("Expected len 20")
	}
	if v.VectorData[4].Type.Kind() != ComplexTypeArray {
		t.Error("Expected value at index to be type ComplexTypeArray")
	}
	if len(v.VectorData[4].VectorData) != 2 {
		t.Error("Expected subarray len to be 2")
	}
	if v.VectorData[4].VectorData[0].Type != PrimitiveTypeBool {
		t.Error("Expected subarray subtype to be bool")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDefaultVariantValueBasicStruct(t *testing.T) {
	at := StructType{
		Fields: []NamedType{
			NamedType{
				Type:  PrimitiveTypeInt,
				Ident: "IntType",
			},
			NamedType{
				Type:  PrimitiveTypeBool,
				Ident: "BoolType",
			},
		},
	}
	v, err := DefaultVariantValue(at)
	if v.Type.Kind() != ComplexTypeStruct {
		t.Error("Expected struct type")
	}
	if len(v.NamedData) != 2 {
		t.Error("Expected len 2")
	}

	if v.NamedData["IntType"].Type != PrimitiveTypeInt {
		t.Error("Expected IntType == type int")
	}
	if v.NamedData["BoolType"].Type != PrimitiveTypeBool {
		t.Error("Expected BoolType == type bool")
	}

	if err != nil {
		t.Error(err)
	}
}
