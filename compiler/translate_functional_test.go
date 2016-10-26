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
			Globals: ns,
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
	if _, ok := p[0].(ast.NamedType); !ok {
		t.Error("First parameter not a NamedType:", reflect.TypeOf(p[0]))
	}
	if p1, ok := p[0].(ast.NamedType); !ok {
		t.Error("Parameter one is not named")
	} else {
		if p1.Name() != "inp1" {
			t.Error("First parameter is named incorrectly")
		}
	}
	if p2, ok := p[1].(ast.NamedType); !ok {
		t.Error("Parameter two is not named")
	} else {
		if p2.Name() != "inp2" {
			t.Error("Second parameter is named incorrectly")
		}
	}

	if p[0].ConcreteType() != ast.PRIMITIVE_TYPE_INT {
		t.Error("First parameter incorrect")
	}
	if _, ok := p[1].(ast.NamedType); !ok {
		t.Error("Second parameter not a NamedType:", reflect.TypeOf(p[0]))
	}
	if p[1].ConcreteType() != ast.PRIMITIVE_TYPE_INT {
		t.Error("Second parameter incorrect")
	}
	if len(context.Declarations[0].Results) != 1 {
		t.Error("Incorrect number of results")
	}
	if _, ok := context.Declarations[0].Results[0].(ast.TypeKindDescription); !ok {
		t.Error("Unexpected return type")
	}
	r := context.Declarations[0].Results
	if r[0] != ast.PRIMITIVE_TYPE_STRING {
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

	var v ast.Variant
	var ok bool
	if v, ok = context.Globals["testVar"]; !ok {
		t.Error("Global expected")
		t.FailNow()
	}
	if v.Type != ast.PRIMITIVE_TYPE_INT {
		t.Error("Integer type for variable expected, got ", v.Type.String())
	}
	if v.Int != 0 {
		t.Error("Default value not 0")
	}
}

func TestGlobalStringSavedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    var testVar string`, t)

	var v ast.Variant
	var ok bool
	if v, ok = context.Globals["testVar"]; !ok {
		t.Error("Global expected")
		t.FailNow()
	}
	if v.Type != ast.PRIMITIVE_TYPE_STRING {
		t.Error("String type for variable expected, got ", v.Type.String())
	}
	if v.String != "" {
		t.Error("Default value not \"\"")
	}
}

func TestAssignReturnsASTStructure(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

		var a int

    func test(){
      a = 3
    }`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.Assign); !ok {
		t.Error("Assign node expected")
	}

	if context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.Assign).NewLocal == true {
		t.Error("NewLocal flag truth expected")
	}

	if _, ok := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.Assign).Value.(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral node expected")
	}
	if context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.Assign).Value.(*ast.IntegerLiteral).Val != 3 {
		t.Error("Incorrect literal value, expected 3")
	}
}

func TestLocalDeclarationProducesAssignNodeSimple(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func test(){
			var testVar int
			var testString string
		}`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Code.(*ast.StatementList).Stmts) != 2 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	var s *ast.StatementList
	var ok bool
	if s, ok = context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.StatementList); !ok {
		t.Error("StatementList node expected")
	}

	if _, ok2 := s.Stmts[0].(*ast.Assign); !ok2 {
		t.Error("Assign Expected")
	}
	if s.Stmts[0].(*ast.Assign).NewLocal == false {
		t.Error("NewLocal flag truth expected")
	}
	if _, ok3 := s.Stmts[0].(*ast.Assign).Value.(*ast.IntegerLiteral); !ok3 {
		t.Error("IntegerLiteral Expected")
	}
	if s.Stmts[0].(*ast.Assign).Value.(*ast.IntegerLiteral).Val != 0 {
		t.Error("Default value expected")
	}

	if s, ok = context.Declarations[0].Code.(*ast.StatementList).Stmts[1].(*ast.StatementList); !ok {
		t.Error("StatementList node expected")
	}

	if _, ok := s.Stmts[0].(*ast.Assign); !ok {
		t.Error("Assign Expected")
	}
	if s.Stmts[0].(*ast.Assign).NewLocal == false {
		t.Error("NewLocal flag truth expected")
	}
	if _, ok := s.Stmts[0].(*ast.Assign).Value.(*ast.StringLiteral); !ok {
		t.Error("StringLiteral Expected")
	}
	if s.Stmts[0].(*ast.Assign).Value.(*ast.StringLiteral).Str != "" {
		t.Error("Default value expected")
	}
}

func TestArrayLocalDeclarationProducesCorrectAST(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func test(){
			var testVar [40]int
		}`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	var s *ast.StatementList
	var ok bool
	if s, ok = context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.StatementList); !ok {
		t.Error("StatementList node expected")
	}

	if _, ok2 := s.Stmts[0].(*ast.Assign); !ok2 {
		t.Error("Assign Expected")
	}
	if s.Stmts[0].(*ast.Assign).NewLocal == false {
		t.Error("NewLocal flag truth expected")
	}
	if _, ok3 := s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral); !ok3 {
		t.Error("ArrayLiteral Expected")
	}
	if s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind() != ast.COMPLEX_TYPE_ARRAY {
		t.Error("COMPLEX_TYPE_ARRAY expected, got", s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind())
	}
	lenNode := s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.(ast.ArrayType).Len
	if lenNode == nil {
		t.Error("Len node expected")
	}
	if lenLit, ok := lenNode.(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral len expected")
	} else {
		if lenLit.Val != 40 {
			t.Error("Expected len of 4")
		}
	}
}

func TestNestedArrayLocalDeclarationProducesCorrectAST(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func test(){
			var testVar [40][50]int
		}`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	var s *ast.StatementList
	var ok bool
	if s, ok = context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.StatementList); !ok {
		t.Error("StatementList node expected")
	}

	if _, ok2 := s.Stmts[0].(*ast.Assign); !ok2 {
		t.Error("Assign Expected")
	}
	if s.Stmts[0].(*ast.Assign).NewLocal == false {
		t.Error("NewLocal flag truth expected")
	}
	if _, ok3 := s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral); !ok3 {
		t.Error("ArrayLiteral Expected")
	}
	if s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind() != ast.COMPLEX_TYPE_ARRAY {
		t.Error("COMPLEX_TYPE_ARRAY expected, got", s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind())
	}
	subType := s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.(ast.ArrayType).SubType
	if subType == nil {
		t.Error("subType node expected")
	}
	if arraySubType, ok := subType.(ast.ArrayType); !ok {
		t.Error("ArrayType expected")
	} else {
		if arraySubType.ConcreteType() != ast.PRIMITIVE_TYPE_INT {
			t.Error("Expected underlying int type")
		}
	}
}

func TestLocalArrayInitialization(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func test(){
			var testVar [4]int = [4]int{1,2,3,4}
		}`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	var s *ast.StatementList
	var ok bool
	if s, ok = context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.StatementList); !ok {
		t.Error("StatementList node expected")
	}

	if _, ok2 := s.Stmts[0].(*ast.Assign); !ok2 {
		t.Error("Assign Expected")
	}
	if s.Stmts[0].(*ast.Assign).NewLocal == false {
		t.Error("NewLocal flag truth expected")
	}
	if _, ok3 := s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral); !ok3 {
		t.Error("ArrayLiteral Expected")
	}
	if s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind() != ast.COMPLEX_TYPE_ARRAY {
		t.Error("COMPLEX_TYPE_ARRAY expected, got", s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind())
	}
	literals := s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Literal
	if lit, ok := literals[0].(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral expected")
	} else {
		if lit.Val != 1 {
			t.Error("Incorrect value [0] != 1")
		}
	}
	if lit, ok := literals[3].(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral expected")
	} else {
		if lit.Val != 4 {
			t.Error("Incorrect value [3] != 4")
		}
	}
}

func TestNestedLocalArrayInitialization(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func test(){
			var testVar [2][2]int = [2][2]int{
				int{1,7},
				int{3,4},
			}
		}`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	var s *ast.StatementList
	var ok bool
	if s, ok = context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.StatementList); !ok {
		t.Error("StatementList node expected")
	}

	if _, ok2 := s.Stmts[0].(*ast.Assign); !ok2 {
		t.Error("Assign Expected")
	}
	if s.Stmts[0].(*ast.Assign).NewLocal == false {
		t.Error("NewLocal flag truth expected")
	}
	if _, ok3 := s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral); !ok3 {
		t.Error("ArrayLiteral Expected")
	}
	if s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind() != ast.COMPLEX_TYPE_ARRAY {
		t.Error("COMPLEX_TYPE_ARRAY expected, got", s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind())
	}
	literals := s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Literal
	if lit, ok := literals[0].(*ast.ArrayLiteral); !ok {
		t.Error("ArrayLiteral expected")
	} else {
		if lit.Type.ConcreteType() != ast.PRIMITIVE_TYPE_INT {
			t.Error("Expected underlying type to be PRIMITIVE_TYPE_INT")
		}
		if intLit, ok := lit.Literal[1].(*ast.IntegerLiteral); !ok {
			t.Error("Expected inner literal to be an integer")
		} else if intLit.Val != 7 {
			t.Error("Incorrect value for [0][1]")
		}
	}
}

func TestLocalDeclarationInitialization(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func test(){
			var testVar int = 4
		}`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	var s *ast.StatementList
	var ok bool
	if s, ok = context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.StatementList); !ok {
		t.Error("StatementList node expected")
	}

	if _, ok2 := s.Stmts[0].(*ast.Assign); !ok2 {
		t.Error("Assign Expected")
	}
	if s.Stmts[0].(*ast.Assign).NewLocal == false {
		t.Error("NewLocal flag truth expected")
	}
	if _, ok3 := s.Stmts[0].(*ast.Assign).Value.(*ast.IntegerLiteral); !ok3 {
		t.Error("IntegerLiteral Expected")
	}
	if s.Stmts[0].(*ast.Assign).Value.(*ast.IntegerLiteral).Val != 4 {
		t.Error("value 4 expected")
	}
}

func TestIfStatementASTStructureGeneratedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func test()bool{
			if true {
				return true
			} else {
				return false
			}
		}`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Identifier != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList, got ", reflect.TypeOf(context.Declarations[0].Code))
	}
	if _, ok := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.IfStmt); !ok {
		t.Error("Expected root node for declaration to be IfStmt")
		t.FailNow()
	}
	ifNode := context.Declarations[0].Code.(*ast.StatementList).Stmts[0].(*ast.IfStmt)
	if ifNode.Init != nil {
		t.Error("No init node expected")
	}
	if _, ok := ifNode.Conditional.(*ast.BoolLiteral); !ok {
		t.Error("Expected boolean literal condition")
	}
	if _, ok := ifNode.Code.(*ast.StatementList); !ok {
		t.Error("Expecting StatementList node for the If statment's code block")
	}
	if _, ok := ifNode.Else.(*ast.StatementList); !ok {
		t.Error("Expecting StatementList node for the If statment's 'else' code block")
	}
}
