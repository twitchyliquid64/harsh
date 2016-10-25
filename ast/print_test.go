package ast

import "testing"

func TestPrintTypeKindIntegerString(t *testing.T) {
	p := PRIMITIVE_TYPE_INT
	if p.String() != "int" {
		t.FailNow()
	}
}

func TestPrintTypeKindStringString(t *testing.T) {
	p := PRIMITIVE_TYPE_STRING
	if p.String() != "string" {
		t.FailNow()
	}
}

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
