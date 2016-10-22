package ast

import "testing"

func TestMakeVariantString(t *testing.T) {
	v := MakeVariant("abc")
	if v.Type.Kind != PRIMITIVE_TYPE_STRING {
		t.Error("Expected string type")
	}
	if v.String != "abc" {
		t.Error("Expected \"abc\"")
	}
}

func TestMakeVariantInt(t *testing.T) {
	v := MakeVariant(123)
	if v.Type.Kind != PRIMITIVE_TYPE_INT {
		t.Error("Expected int type")
	}
	if v.Int != 123 {
		t.Error("Expected 123")
	}
}

func TestMakeVariantInt64(t *testing.T) {
	v := MakeVariant(int64(-53))
	if v.Type.Kind != PRIMITIVE_TYPE_INT {
		t.Error("Expected int type")
	}
	if v.Int != -53 {
		t.Error("Expected 123")
	}
}

type WierdUnknownTypeMock int

func TestMakeVariantUndefined(t *testing.T) {
	v := MakeVariant(WierdUnknownTypeMock(1))
	if v.Type.Kind != PRIMITIVE_TYPE_UNDEFINED {
		t.Error("Expected undefined type")
	}
}
