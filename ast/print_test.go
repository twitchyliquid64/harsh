package ast

import "testing"

func TestTypeKindIntegerString(t *testing.T) {
	p := PRIMITIVE_TYPE_INT
	if p.String() != "int" {
		t.FailNow()
	}
}

func TestTypeKindStringString(t *testing.T) {
	p := PRIMITIVE_TYPE_STRING
	if p.String() != "string" {
		t.FailNow()
	}
}

func TestPrimitiveTypeString(t *testing.T) {
	pt := PrimitiveType{
		Kind: PRIMITIVE_TYPE_INT,
		Name: "dsfsdfds",
	}
	if pt.String() != "int{dsfsdfds}" {
		t.Error("PrimitiveType integer .String() is incorrect")
	}
}

func TestPrimitiveTypeInteger(t *testing.T) {
	pt := PrimitiveType{
		Kind: PRIMITIVE_TYPE_STRING,
		Name: "dsfsdfds",
	}
	if pt.String() != "string{dsfsdfds}" {
		t.Error("PrimitiveType string .String() is incorrect")
	}
}
