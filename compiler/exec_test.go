package compiler

import (
	"testing"

	"github.com/twitchyliquid64/harsh/ast"
)

func TestBasicCallFuncReturnsUndefined(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func TestNoReturn(){

    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("TestNoReturn", nil)
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}
	if r.Type.Kind != ast.PRIMITIVE_TYPE_UNDEFINED {
		t.Error("Expected PRIMITIVE_TYPE_UNDEFINED")
	}
}

func TestBasicCallFuncBadNameFails(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func TestBadName(){

    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	v, err := c.CallFunc("TestNoReturn", nil)
	if err != ErrFuncNotFound {
		t.Error("CallFunc(): Unexpected Error")
		t.Error(err)
		t.FailNow()
	}
	if v.Type.Kind != ast.PRIMITIVE_TYPE_UNDEFINED {
		t.Error("Expected undefined result")
	}
}

func TestBasicCallFuncReturnsLiteral(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func Test(){
      return 1
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("Test", nil)
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}
	if r.Type.Kind != ast.PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT")
	}
	if r.Int != 1 {
		t.Error("Expected value 1")
	}
}

func TestBasicCallFuncReturnsArithmeticResult(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func TestArithmetic(){
      return (2*3) - 1 + 2
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("TestArithmetic", nil)
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}
	if r.Type.Kind != ast.PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT")
	}
	if r.Int != 7 {
		t.Error("Expected value 1")
	}
}
