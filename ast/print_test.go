package ast

import "testing"

func TestPrintPrimitiveTypeString(t *testing.T) {
	pt := PrimitiveTypeInt
	if pt.String() != "int" {
		t.Error("PrimitiveType integer .String() is incorrect")
	}
}

func TestPrintPrimitiveTypeInteger(t *testing.T) {
	pt := PrimitiveTypeString

	if pt.String() != "string" {
		t.Error("PrimitiveType string .String() is incorrect")
	}
}

func TestPrintPrimitiveTypeBool(t *testing.T) {
	pt := PrimitiveTypeBool

	if pt.String() != "bool" {
		t.Error("PrimitiveType bool .String() is incorrect")
	}
}

func TestPrintPrimitiveTypeUndefined(t *testing.T) {
	pt := PrimitiveTypeUndefined

	if pt.String() != "undefined" {
		t.Error("PrimitiveType undef .String() is incorrect")
	}
}

func TestPrintArrayType(t *testing.T) {
	pt := ComplexTypeArray

	if pt.String() != "[?]" {
		t.Error("Expected [?], got", pt.String())
	}
}

func TestPrintBinops(t *testing.T) {
	b := BinOpAdd
	if b.String() != "+" {
		t.Error("BinOpAdd.String() incorrect")
	}
	b = BinOpDiv
	if b.String() != "/" {
		t.Error("BinOpDiv.String() incorrect")
	}
	b = BinOpEquality
	if b.String() != "==" {
		t.Error("BinOpEquality.String() incorrect")
	}
	b = BinOpLAnd
	if b.String() != "&&" {
		t.Error("BinOpLAnd.String() incorrect")
	}
	b = BinOpLOr
	if b.String() != "||" {
		t.Error("BinOpLOr.String() incorrect")
	}
	b = BinOpMod
	if b.String() != "%" {
		t.Error("BinOpMod.String() incorrect")
	}
	b = BinOpMul
	if b.String() != "*" {
		t.Error("BinOpMul.String() incorrect")
	}
	b = BinOpSub
	if b.String() != "-" {
		t.Error("BinOpSub.String() incorrect")
	}
	b = BinOpUnknown
	if b.String() != "UNK?" {
		t.Error("BinOpUnknown.String() incorrect:", b.String())
	}
}
