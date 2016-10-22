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

func TestBasicCallFuncReturnsIntLiteral(t *testing.T) {
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

func TestBasicCallFuncReturnsStringLiteral(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func Test(){
      return "bantz"
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
	if r.Type.Kind != ast.PRIMITIVE_TYPE_STRING {
		t.Error("Expected PRIMITIVE_TYPE_STRING")
	}
	if r.String != "bantz" {
		t.Error("Expected value \"bantz\"")
	}
}

func TestStringConcatCallFuncReturnsStringLiteral(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func Test(){
      return "bantz" + " :D"
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
	if r.Type.Kind != ast.PRIMITIVE_TYPE_STRING {
		t.Error("Expected PRIMITIVE_TYPE_STRING")
	}
	if r.String != "bantz :D" {
		t.Error("Expected value \"bantz :D\"")
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
		t.Error("Expected value 7")
	}
}

func TestBasicCallFuncReturnsParameters(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func TestIntParam(a_input int){
      return a_input
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("TestIntParam", map[string]interface{}{
		"a_input": 4,
	})
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}
	if r.Type.Kind != ast.PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT")
	}
	if r.Int != 4 {
		t.Error("Expected value 4")
	}
}

func TestBasicCallFuncReturnsArithmeticFromParam(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func TestIntParam(a_input int){
      return (a_input * 3) + 1
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("TestIntParam", map[string]interface{}{
		"a_input": 4,
	})
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}
	if r.Type.Kind != ast.PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT")
	}
	if r.Int != 13 {
		t.Error("Expected value 4")
	}
}
