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
	if r.Type != ast.PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined")
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
	if v.Type != ast.PrimitiveTypeUndefined {
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
	if r.Type != ast.PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt")
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
	if r.Type != ast.PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString")
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
	if r.Type != ast.PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString")
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
	if r.Type != ast.PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt")
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
	if r.Type != ast.PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt")
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
	if r.Type != ast.PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt")
	}
	if r.Int != 13 {
		t.Error("Expected value 4")
	}
}

func TestGlobalReadCorrectly(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

		var testVar string

    func testFetch()string{
      return testVar
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}
	c.Globals.Save("testVar", "testData121")

	r, err := c.CallFunc("testFetch", nil)
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}
	if r.Type != ast.PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString, got " + r.Type.String())
	}
	if r.String != "testData121" {
		t.Error("Expected value testData121, got '" + r.String + "'")
	}
}

func TestNewLocalWriteCorrectly(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func testFetch()string{
			testVar := "abc"
      return testVar
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("testFetch", nil)
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}

	if r.Type != ast.PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString, got " + r.Type.String())
	}
	if r.String != "abc" {
		t.Error("Expected value abc, got '" + r.String + "'")
	}
}

func TestGlobalWriteCorrectly(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

		var testVar string
		var test2 string

    func testFetch()string{
			testVar = "abbc"
			test2 = "1234"
			testVar = test2 + testVar
      return testVar
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("testFetch", nil)
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}

	if r.Type != ast.PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString, got " + r.Type.String())
	}
	if r.String != "1234abbc" {
		t.Error("Expected value abc, got '" + r.String + "'")
	}
}

func TestAssignScopePrecedenceCorrectness(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

		var testVar string
		var crap string

    func testFetch()string{
			testVar = "abbc"
			crap = "abc"
			crap := "123"
      return testVar + crap
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("testFetch", nil)
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}

	if r.Type != ast.PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString, got " + r.Type.String())
	}
	if r.String != "abbc123" {
		t.Error("Expected value abc, got '" + r.String + "'")
	}
}

func TestLocalDeclarationDefaultsInt(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func test()int{
			var t1 int
			var t2 int = 44
      return t1 + t2
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("test", nil)
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}

	if r.Type != ast.PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt, got " + r.Type.String())
	}
	if r.Int != 44 {
		t.Error("Expected value 44, got ", r.Int)
	}
}

func TestLocalDeclarationDefaultsString(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func test()string{
			var t1 string
			var t2 string = "abc"
      return t1 + t2
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("test", nil)
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}

	if r.Type != ast.PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString, got " + r.Type.String())
	}
	if r.String != "abc" {
		t.Error("Expected value abc, got ", r.Int)
	}
}

func TestBasicBoolType(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func TestBoolParam(a_input bool)bool{
      return a_input
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("TestBoolParam", map[string]interface{}{
		"a_input": true,
	})
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}
	if r.Type != ast.PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool")
	}
	if r.Bool != true {
		t.Error("Expected value false")
	}
}

func TestIfStatementReturnsCorrectly(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

		var didRunInit bool

    func TestIF(a_input bool)bool{
      if didRunInit = true; a_input {
				return true
			} else {
				return false
			}
    }

		func getInitVar()bool{
			return didRunInit
		}
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("TestIF", map[string]interface{}{
		"a_input": true,
	})
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}
	if r.Type != ast.PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool, got ", r.Type.String())
	}
	if r.Bool != true {
		t.Error("Expected value true")
	}

	r, err = c.CallFunc("getInitVar", map[string]interface{}{})
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}
	if r.Type != ast.PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool, got ", r.Type.String())
	}
	if r.Bool != true {
		t.Error("Expected value true")
	}
}

func TestArrayLocalDeclarationAssignment(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func Test()int{
      var testVar [4]int
			testVar[1] = 4
			return testVar[1]
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("Test", map[string]interface{}{})
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}
	if r.Type != ast.PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt, got ", r.Type.String())
	}
	if r.Int != 4 {
		t.Error("Expected value 4")
	}
}

func TestArrayLocalInitializationWorksAndReadable(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func Test()int{
      var testVar [4]int = [4]int{1,2,3,4}
			return testVar[1]
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("Test", map[string]interface{}{})
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}
	if r.Type != ast.PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt, got ", r.Type.String())
	}
	if r.Int != 2 {
		t.Error("Expected value 2")
	}
}

func TestArrayShorthandAssignmentWorks(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func Test(inVal int)int{
      testVar := [4]int{1,8,3,inVal}
			return testVar[3]
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, err := c.CallFunc("Test", map[string]interface{}{
		"inVal": 123,
	})
	if err != nil {
		t.Error("CallFunc(): Error")
		t.Error(err)
		t.FailNow()
	}
	if r.Type != ast.PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt, got ", r.Type.String())
	}
	if r.Int != 123 {
		t.Error("Expected value 123")
	}
}

func TestArrayOutOfBoundsWritesError(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func Test()int{
      testVar := [4]int{1,8,3,123}
			return testVar[30]
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, er := c.CallFunc("Test", map[string]interface{}{})
	if er == nil {
		t.Error("Errors expected")
	} else {
		errs := er.(ExecutionError)
		if len(errs.Errors) != 1 {
			t.Error("Expected one error")
		}
		if errs.Errors[0].Class != ast.BoundsErr {
			t.Error("Expected error ast.BoundsErr")
		}
	}
	if r.Type != ast.PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined, got ", r.Type.String())
	}
}

func TestBooleanOperationsAreCorrect(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func Test()bool{
			return true && false || false && (true == false)
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, er := c.CallFunc("Test", map[string]interface{}{})
	if er != nil {
		t.Error("Errors when executing")
	}
	if r.Type != ast.PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool")
	}
	if r.Bool != false {
		t.Error("Incorrect value")
	}
}

func TestStringEquality(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func Test()int{
			if "a" == "b" {
				return 1
			}
			if "ab" == "ab" {
				return 2
			}
			return 0
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, er := c.CallFunc("Test", map[string]interface{}{})
	if er != nil {
		t.Error("Errors when executing")
	}
	if r.Type != ast.PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt")
	}
	if r.Int != 2 {
		t.Error("Incorrect value")
	}
}

func TestLayeredArrayAccessAndAssignWorks(t *testing.T) {
	c, err := ParseLiteral("main.go", `package test

		func Test() string {
			var testVar [2][2]string = [2][2]string{
				[2]string{
					"a",
					"b",
				},
				[2]string{
					"c",
					"d",
				},
			}
			testVar[0][1] = "CRAPLOL"
			return testVar[0][1] + testVar[1][1]
		}`)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, er := c.CallFunc("Test", map[string]interface{}{})
	if er != nil {
		t.Error("Errors when executing")
	}
	if r.Type != ast.PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString")
	}
	if r.String != "CRAPLOLd" {
		t.Error("Incorrect value")
	}
}

func TestUnaryNotBooleanReturnsCorrect(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

    func Test()int{
			return !false
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, er := c.CallFunc("Test", map[string]interface{}{})
	if er != nil {
		t.Error("Errors when executing")
	}
	if r.Type != ast.PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool")
	}
	if r.Bool != true {
		t.Error("Incorrect value")
	}
}

func TestMultiFunctionDeclAndCall(t *testing.T) {
	c, err := ParseLiteral("test.go", `
    package test

		func testSub2(in2, in1 int) bool {
			return (2*in2) == in1
		}

		func testSub(in1 int) bool {
			return testSub2(in1, 8)
		}

    func Test() bool{
			return testSub(4)
    }
    `)

	if err != nil {
		t.Error("ParseLiteral(): Error")
		t.Error(err)
		t.FailNow()
	}

	r, er := c.CallFunc("Test", map[string]interface{}{})
	if er != nil {
		t.Error("Errors when executing")
	}
	if r.Type != ast.PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool")
	}
	if r.Bool != true {
		t.Error("Incorrect value")
	}
}
