package ast

import "testing"

func TestPrintPrimitiveTypeString(t *testing.T) {
	pt := PRIMITIVE_TYPE_INT
	if pt.String() != "int" {
		t.Error("PrimitiveType integer .String() is incorrect")
	}
}

func TestPrintPrimitiveTypeInteger(t *testing.T) {
	pt := PRIMITIVE_TYPE_STRING

	if pt.String() != "string" {
		t.Error("PrimitiveType string .String() is incorrect")
	}
}

func TestPrintPrimitiveTypeBool(t *testing.T) {
	pt := PRIMITIVE_TYPE_BOOL

	if pt.String() != "bool" {
		t.Error("PrimitiveType bool .String() is incorrect")
	}
}

func TestPrintPrimitiveTypeUndefined(t *testing.T) {
	pt := PRIMITIVE_TYPE_UNDEFINED

	if pt.String() != "undefined" {
		t.Error("PrimitiveType undef .String() is incorrect")
	}
}

func TestPrintArrayType(t *testing.T) {
	pt := COMPLEX_TYPE_ARRAY

	if pt.String() != "[?]" {
		t.Error("Expected [?], got", pt.String())
	}
}
