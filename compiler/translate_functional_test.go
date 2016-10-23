package compiler

import (
	goast "go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"

	"github.com/twitchyliquid64/harsh/ast"
)

func setupTestGetAST(context *Context, inCode string, t *testing.T) (ast.Node, *Context) {
	fset := token.NewFileSet()
	goAst, err := parser.ParseFile(fset, "testfile.go", inCode, 0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if context == nil {
		ns := ast.Namespace(map[string]ast.Variant{})
		context = &Context{
			ConType: CONTEXT_ADHOC,
			Globals: &ns,
		}
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	a := translateGoNode(fset, context, reflect.ValueOf(goAst))
	if a != nil {
		t.Error("translateGoNode() did not return nil")
		t.FailNow()
	}
	return a, context
}

func TestLiteralReturnASTStructure(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func testLiteralReturn()int{
      return 3
    }`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "testLiteralReturn" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt); !ok {
		t.Error("ReturnStmt node expected")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral node expected")
	}
	if context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.IntegerLiteral).Val != 3 {
		t.Error("Incorrect literal value, expected 3")
	}
}

func TestFuncParamReturnASTStructure(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func testParamReturn(in int)int{
      return in
    }`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "testParamReturn" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt); !ok {
		t.Error("ReturnStmt node expected")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.VariableReference); !ok {
		t.Error("VariableReference node expected")
	}
	if context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.VariableReference).Name != "in" {
		t.Error("Incorrect ident value, expected \"in\"")
	}
}

func TestBasicArithmeticASTStructureCorrectness(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func TestBasicArithmetic()int{
      return 3+1
    }`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "TestBasicArithmetic" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt); !ok {
		t.Error("ReturnStmt node expected")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.BinaryOp); !ok {
		t.Error("BinaryOp node expected")
	}
	op := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.BinaryOp)
	if op.Op != ast.BINOP_ADD {
		t.Error("Addition operation expected")
	}
	if _, ok := op.LHS.(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral node expected")
	}
	if _, ok := op.RHS.(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral node expected")
	}
	if op.RHS.(*ast.IntegerLiteral).Val != 1 || op.LHS.(*ast.IntegerLiteral).Val != 3 {
		t.Error("Values incorrect")
	}
}

func TestComplexArithmeticASTStructureCorrectness(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func TestComplexArithmetic()int{
      return 2*(1+2)-4
    }`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "TestComplexArithmetic" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt); !ok {
		t.Error("ReturnStmt node expected")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.BinaryOp); !ok {
		t.Error("BinaryOp node expected")
	}
	op := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.BinaryOp)
	if op.Op != ast.BINOP_SUB {
		t.Error("Subtraction operation expected")
	}
	if _, ok := op.LHS.(*ast.BinaryOp); !ok {
		t.Error("BinaryOp node expected")
	}
	if op.LHS.(*ast.BinaryOp).Op != ast.BINOP_MUL {
		t.Error("Multiplication operation expected")
	}
	if _, ok := op.RHS.(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral node expected")
	}
	if op.RHS.(*ast.IntegerLiteral).Val != 4 {
		t.Error("Value incorrect")
	}

	op = op.LHS.(*ast.BinaryOp)
	if _, ok := op.LHS.(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral node expected")
	}
	if op.LHS.(*ast.IntegerLiteral).Val != 2 {
		t.Error("Values incorrect")
	}

	if _, ok := op.RHS.(*ast.BinaryOp); !ok {
		t.Error("BinaryOp node expected")
	}
	if op.RHS.(*ast.BinaryOp).Op != ast.BINOP_ADD {
		t.Error("Multiplication operation expected")
	}
	op = op.RHS.(*ast.BinaryOp)
	if _, ok := op.LHS.(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral node expected")
	}
	if _, ok := op.RHS.(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral node expected")
	}
	if op.RHS.(*ast.IntegerLiteral).Val != 2 || op.LHS.(*ast.IntegerLiteral).Val != 1 {
		t.Error("Values incorrect")
	}
}

func TestFunctionParamsAndResultsAreTypedAndNamedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func testParamResultsDefinition(inp1 int, inp2 int)string{
      return 3
    }`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "testParamResultsDefinition" {
		t.Error("Unexpected declaration name")
	}
	if len(context.Declarations[0].Parameters) != 2 {
		t.Error("Expected 2 parameters")
	}
	p := context.Declarations[0].Parameters
	if _, ok := p[0].(*ast.PrimitiveType); !ok {
		t.Error("First parameter not a primitive")
	}
	if p[0].(*ast.PrimitiveType).Kind != ast.PRIMITIVE_TYPE_INT || p[0].(*ast.PrimitiveType).Name != "inp1" {
		t.Error("First parameter incorrect")
	}
	if _, ok := p[1].(*ast.PrimitiveType); !ok {
		t.Error("Second parameter not a primitive")
	}
	if p[1].(*ast.PrimitiveType).Kind != ast.PRIMITIVE_TYPE_INT || p[1].(*ast.PrimitiveType).Name != "inp2" {
		t.Error("Second parameter incorrect")
	}
	if len(context.Declarations[0].Results) != 1 {
		t.Error("Incorrect number of results")
	}
	if _, ok := context.Declarations[0].Results[0].(*ast.PrimitiveType); !ok {
		t.Error("Unexpected return type")
	}
	r := context.Declarations[0].Results
	if r[0].(*ast.PrimitiveType).Kind != ast.PRIMITIVE_TYPE_STRING || r[0].(*ast.PrimitiveType).Name != "" {
		t.Error("Return incorrect")
	}
}

func TestContextMetadataCorrectness(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func testCrap(){
      return 3
    }`, t)
	if context.Name != "test" {
		t.Error("Context name incorrect")
	}
	if len(context.ChildContexts) != 0 {
		t.Error("Unexpected child contexts")
	}
}

func TestFileContextCorrectness(t *testing.T) {
	_, context := setupTestGetAST(&Context{
		ConType: CONTEXT_FILE,
	}, `
    package test

    func testCrap(){
      return 3
    }
    func testCrap2(){
      return 6
    }
    `, t)

	if len(context.ChildContexts) != 1 {
		t.Error("1 child context (one for each file) expected, got ", len(context.ChildContexts))
		goast.Print(nil, context)
	}

	if context.ChildContexts[0].Name != "test" {
		t.Error("Incorrect child context name")
	}
	if len(context.ChildContexts[0].Declarations) != 2 {
		t.Error("Incorrect number of declarations in child context")
	}
	if context.ChildContexts[0].Declarations[0].Identifier != "testCrap" || context.ChildContexts[0].Declarations[1].Identifier != "testCrap2" {
		t.Error("Declarations in child context named incorrectly")
	}

	_, context = setupTestGetAST(context, `
    package testSecond

    func test2Crap(){
      return 3
    }
    func test2Crap2(){
      return 6
    }
    `, t)

	if len(context.ChildContexts) != 2 {
		t.Error("1 child context (one for each file) expected, got ", len(context.ChildContexts))
		goast.Print(nil, context)
	}

	if context.ChildContexts[1].Name != "testSecond" {
		t.Error("Incorrect child context name")
	}
	if len(context.ChildContexts[1].Declarations) != 2 {
		t.Error("Incorrect number of declarations in child context")
	}
	if context.ChildContexts[1].Declarations[0].Identifier != "test2Crap" || context.ChildContexts[1].Declarations[1].Identifier != "test2Crap2" {
		t.Error("Declarations in child context named incorrectly")
	}
}

func TestGlobalIntSavedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    var testVar int`, t)
	if context.Name != "test" {
		t.Error("Context name incorrect")
	}

	var v ast.Variant
	var ok bool
	if v, ok = (*context.Globals)["testVar"]; !ok {
		t.Error("Global expected")
		t.FailNow()
	}
	if v.Type.Kind != ast.PRIMITIVE_TYPE_INT {
		t.Error("Integer type for variable expected, got ", v.Type.Kind.String())
	}
	if v.Int != 0 {
		t.Error("Default value not 0")
	}
}

func TestGlobalStringSavedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    var testVar string`, t)
	if context.Name != "test" {
		t.Error("Context name incorrect")
	}

	var v ast.Variant
	var ok bool
	if v, ok = (*context.Globals)["testVar"]; !ok {
		t.Error("Global expected")
		t.FailNow()
	}
	if v.Type.Kind != ast.PRIMITIVE_TYPE_STRING {
		t.Error("String type for variable expected, got ", v.Type.Kind.String())
	}
	if v.String != "" {
		t.Error("Default value not \"\"")
	}
}
