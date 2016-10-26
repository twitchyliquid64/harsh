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

func TestPrintBinops(t *testing.T) {
	b := BINOP_ADD
	if b.String() != "+" {
		t.Error("BINOP_ADD.String() incorrect")
	}
	b = BINOP_DIV
	if b.String() != "/" {
		t.Error("BINOP_DIV.String() incorrect")
	}
	b = BINOP_EQUALITY
	if b.String() != "==" {
		t.Error("BINOP_EQUALITY.String() incorrect")
	}
	b = BINOP_LAND
	if b.String() != "&&" {
		t.Error("BINOP_LAND.String() incorrect")
	}
	b = BINOP_LOR
	if b.String() != "||" {
		t.Error("BINOP_LOR.String() incorrect")
	}
	b = BINOP_MOD
	if b.String() != "%" {
		t.Error("BINOP_MOD.String() incorrect")
	}
	b = BINOP_MUL
	if b.String() != "*" {
		t.Error("BINOP_MUL.String() incorrect")
	}
	b = BINOP_SUB
	if b.String() != "-" {
		t.Error("BINOP_SUB.String() incorrect")
	}
	b = BINOP_UNK
	if b.String() != "UNK?" {
		t.Error("BINOP_UNK.String() incorrect:", b.String())
	}
}
