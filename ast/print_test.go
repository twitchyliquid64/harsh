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
	pt := PrimitiveType{
		Kind: PRIMITIVE_TYPE_INT,
		Name: "dsfsdfds",
	}
	if pt.String() != "int{dsfsdfds}" {
		t.Error("PrimitiveType integer .String() is incorrect")
	}
}

func TestPrintPrimitiveTypeInteger(t *testing.T) {
	pt := PrimitiveType{
		Kind: PRIMITIVE_TYPE_STRING,
		Name: "dsfsdfds",
	}
	if pt.String() != "string{dsfsdfds}" {
		t.Error("PrimitiveType string .String() is incorrect")
	}
}

func TestPrintPrimitiveTypeUndefined(t *testing.T) {
	pt := PrimitiveType{
		Kind: PRIMITIVE_TYPE_UNDEFINED,
		Name: "dsfsdfds",
	}
	if pt.String() != "undefined{dsfsdfds}" {
		t.Error("PrimitiveType undef .String() is incorrect")
	}
}
