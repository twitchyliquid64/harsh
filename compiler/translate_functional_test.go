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
		ns := ast.Namespace(map[string]*ast.Variant{})
		context = &Context{
			ConType: ContextAdhoc,
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
	if context.Declarations[0].Ident != "testLiteralReturn" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt); !ok {
		t.Error("ReturnStmt node expected")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral node expected")
	}
	if context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.IntegerLiteral).Val != 3 {
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
	if context.Declarations[0].Ident != "testParamReturn" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt); !ok {
		t.Error("ReturnStmt node expected")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.VariableReference); !ok {
		t.Error("VariableReference node expected")
	}
	if context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.VariableReference).Name != "in" {
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
	if context.Declarations[0].Ident != "TestBasicArithmetic" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt); !ok {
		t.Error("ReturnStmt node expected")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.BinaryOp); !ok {
		t.Error("BinaryOp node expected")
	}
	op := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.BinaryOp)
	if op.Op != ast.BinOpAdd {
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
	if context.Declarations[0].Ident != "TestComplexArithmetic" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt); !ok {
		t.Error("ReturnStmt node expected")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.BinaryOp); !ok {
		t.Error("BinaryOp node expected")
	}
	op := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.BinaryOp)
	if op.Op != ast.BinOpSub {
		t.Error("Subtraction operation expected")
	}
	if _, ok := op.LHS.(*ast.BinaryOp); !ok {
		t.Error("BinaryOp node expected")
	}
	if op.LHS.(*ast.BinaryOp).Op != ast.BinOpMul {
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
	if op.RHS.(*ast.BinaryOp).Op != ast.BinOpAdd {
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
	if context.Declarations[0].Ident != "testParamResultsDefinition" {
		t.Error("Unexpected declaration name")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Parameters) != 2 {
		t.Error("Expected 2 parameters")
	}
	p := context.Declarations[0].Type.(ast.FunctionType).Parameters
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

	if p[0].BaseType() != ast.PrimitiveTypeInt {
		t.Error("First parameter incorrect")
	}
	if _, ok := p[1].(ast.NamedType); !ok {
		t.Error("Second parameter not a NamedType:", reflect.TypeOf(p[0]))
	}
	if p[1].BaseType() != ast.PrimitiveTypeInt {
		t.Error("Second parameter incorrect")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).ReturnType.(ast.TypeKindDescription); !ok {
		t.Error("Unexpected return type")
	}
	r := context.Declarations[0].Type.(ast.FunctionType).ReturnType
	if r != ast.PrimitiveTypeString {
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
		ConType: ContextFile,
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
	if context.ChildContexts[0].Declarations[0].Ident != "testCrap" || context.ChildContexts[0].Declarations[1].Ident != "testCrap2" {
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
	if context.ChildContexts[1].Declarations[0].Ident != "test2Crap" || context.ChildContexts[1].Declarations[1].Ident != "test2Crap2" {
		t.Error("Declarations in child context named incorrectly")
	}
}

func TestGlobalIntSavedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    var testVar int`, t)

	var v *ast.Variant
	var ok bool
	if v, ok = context.Globals["testVar"]; !ok {
		t.Error("Global expected")
		t.FailNow()
	}
	if v.Type != ast.PrimitiveTypeInt {
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

	var v *ast.Variant
	var ok bool
	if v, ok = context.Globals["testVar"]; !ok {
		t.Error("Global expected")
		t.FailNow()
	}
	if v.Type != ast.PrimitiveTypeString {
		t.Error("String type for variable expected, got ", v.Type.String())
	}
	if v.String != "" {
		t.Error("Default value not \"\"")
	}
}

func TestGlobalIntArraySavedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    var testVar [3]int`, t)

	var v *ast.Variant
	var ok bool
	if v, ok = context.Globals["testVar"]; !ok {
		t.Error("Global expected")
		t.FailNow()
	}
	if v.Type.BaseType().Kind() != ast.PrimitiveTypeInt {
		t.Error("Int type for variable concrete type expected, got ", v.Type.BaseType().String())
	}
	if v.Type.Kind() != ast.ComplexTypeArray {
		t.Error("Complex array variable kind expected, got ", v.Type.Kind().String())
	}
}

func TestGlobalIntArrayArraySavedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    var testVar [3][9]int`, t)

	var v *ast.Variant
	var ok bool
	if v, ok = context.Globals["testVar"]; !ok {
		t.Error("Global expected")
		t.FailNow()
	}
	if v.Type.BaseType().Kind() != ast.ComplexTypeArray {
		t.Error("base type complex array expected, got ", v.Type.BaseType().String())
	}
	if v.Type.Kind() != ast.ComplexTypeArray {
		t.Error("Complex array variable kind expected, got ", v.Type.Kind().String())
	}
	if v.Type.BaseType().(ast.ArrayType).BaseType().Kind() != ast.PrimitiveTypeInt {
		t.Error("second level based type was expected to be int, got", v.Type.BaseType().(ast.ArrayType).BaseType().Kind().String())
	}
}

func TestGlobalStructSavedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    var testVar struct{
			Bro string
			Cuzz int
		}`, t)

	var v *ast.Variant
	var ok bool
	if v, ok = context.Globals["testVar"]; !ok {
		t.Error("Global expected")
		t.FailNow()
	}
	if v.Type.Kind() != ast.ComplexTypeStruct {
		t.Error("ComplexTypeStruct type for variable kind expected, got ", v.Type.String())
	}
	if len(v.Type.(ast.StructType).Fields) != 2 {
		t.Error("Expected 2 fields")
	}
	if v.Type.(ast.StructType).Fields[0].Ident != "Bro" {
		t.Error("First field named incorrectly")
	}
	if v.Type.(ast.StructType).Fields[0].Type != ast.PrimitiveTypeString {
		t.Error("First field typed incorrectly")
	}
	if v.Type.(ast.StructType).Fields[1].Ident != "Cuzz" {
		t.Error("Second field named incorrectly")
	}
	if v.Type.(ast.StructType).Fields[1].Type != ast.PrimitiveTypeInt {
		t.Error("Second field typed incorrectly")
	}
}

func TestGlobalStructWithArraySavedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    var testVar struct{
			MajorLazorArrayBro [3][3]int
			Crap 							 [1111111]string
		}`, t)

	var v *ast.Variant
	var ok bool
	if v, ok = context.Globals["testVar"]; !ok {
		t.Error("Global expected")
		t.FailNow()
	}
	if v.Type.Kind() != ast.ComplexTypeStruct {
		t.Error("ComplexTypeStruct type for variable kind expected, got ", v.Type.String())
	}
	if _, ok := v.Type.(ast.StructType); !ok {
		t.Error("Variable is not struct type, got " + reflect.TypeOf(v.Type).String())
		t.FailNow()
	}
	if len(v.Type.(ast.StructType).Fields) != 2 {
		t.Error("Expected 2 fields")
	}
	if v.Type.(ast.StructType).Fields[0].Ident != "MajorLazorArrayBro" {
		t.Error("First field named incorrectly")
	}
	if v.Type.(ast.StructType).Fields[0].Type.Kind() != ast.ComplexTypeArray {
		t.Error("First field typed incorrectly")
	}
	if v.Type.(ast.StructType).Fields[0].Type.(ast.ArrayType).SubType.Kind() != ast.ComplexTypeArray {
		t.Error("First field sub array type incorrect, got " + v.Type.(ast.StructType).Fields[0].Type.(ast.ArrayType).SubType.String())
	}
	if v.Type.(ast.StructType).Fields[1].Ident != "Crap" {
		t.Error("Second field named incorrectly")
	}
	if v.Type.(ast.StructType).Fields[1].Type.Kind() != ast.ComplexTypeArray {
		t.Error("Second field typed incorrectly")
	}
	if v.Type.(ast.StructType).Fields[1].Type.(ast.ArrayType).SubType.Kind() != ast.PrimitiveTypeString {
		t.Error("Second field sub type incorrect, got " + v.Type.(ast.StructType).Fields[0].Type.(ast.ArrayType).SubType.String())
	}
}

func TestGlobalStructWithStructSavedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    var testVar struct{
			Crap struct{
				Brah string
			}
		}`, t)

	var v *ast.Variant
	var ok bool
	if v, ok = context.Globals["testVar"]; !ok {
		t.Error("Global expected")
		t.FailNow()
	}
	if v.Type.Kind() != ast.ComplexTypeStruct {
		t.Error("ComplexTypeStruct type for variable kind expected, got ", v.Type.String())
	}
	if _, ok := v.Type.(ast.StructType); !ok {
		t.Error("Variable is not struct type, got " + reflect.TypeOf(v.Type).String())
		t.FailNow()
	}
	if len(v.Type.(ast.StructType).Fields) != 1 {
		t.Error("Expected 1 field")
	}
	if v.Type.(ast.StructType).Fields[0].Ident != "Crap" {
		t.Error("First field named incorrectly")
	}
	if v.Type.(ast.StructType).Fields[0].Type.Kind() != ast.ComplexTypeStruct {
		t.Error("First field typed incorrectly")
	}
	subStruct := v.Type.(ast.StructType).Fields[0].Type.(ast.StructType)
	if len(subStruct.Fields) != 1 {
		t.Error("Expected 1 field in subStruct")
	}
	if subStruct.Fields[0].Ident != "Brah" {
		t.Error("Expected name 'Brah' in sub struct first field identifier")
	}
	if subStruct.Fields[0].Type != ast.PrimitiveTypeString {
		t.Error("Expected type string in field of sub struct")
	}
}

func TestSelectorExprTranslatesCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    var testVar struct{
			Crap struct{
				Brah string
			}
		}

		func test() string {
			return testVar.Crap
		}`, t)

	if context.Declarations[1].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[1].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[1].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	nodes := context.Declarations[1].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts

	if _, ok := nodes[0].(*ast.ReturnStmt).Expr.(*ast.NamedSelector); !ok {
		t.Error("Expected node to be NamedSelector, got", reflect.TypeOf(nodes[0]))
		t.FailNow()
	}
	if nodes[0].(*ast.ReturnStmt).Expr.(*ast.NamedSelector).Name != "Crap" {
		t.Error("Expected selector to be brah")
	}
	if nodes[0].(*ast.ReturnStmt).Expr.(*ast.NamedSelector).Expr.(*ast.VariableReference).Name != "testVar" {
		t.Error("Expected variable to be testVar")
	}
}

func TestGlobalArrayWithStructSavedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    var testVar [43]struct{
			Brah string
		}`, t)

	var v *ast.Variant
	var ok bool
	if v, ok = context.Globals["testVar"]; !ok {
		t.Error("Global expected")
		t.FailNow()
	}
	if v.Type.Kind() != ast.ComplexTypeArray {
		t.Error("ComplexTypeArray type for variable kind expected, got ", v.Type.String())
	}
	if _, ok := v.Type.(ast.ArrayType); !ok {
		t.Error("Variable is not array type, got " + reflect.TypeOf(v.Type).String())
		t.FailNow()
	}
	if v.Type.(ast.ArrayType).SubType.Kind() != ast.ComplexTypeStruct {
		t.Error("Expected struct subtype")
		t.FailNow()
	}
	if len(v.Type.(ast.ArrayType).SubType.(ast.StructType).Fields) == 1 && v.Type.(ast.ArrayType).SubType.(ast.StructType).Fields[0].Ident != "Brah" {
		t.Error("Expected 1 field with name Brah")
	}
}

func TestLocalStructSavedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func test(){
			var testVar struct{
				Field string
			}
		}`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	nodes := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts
	if _, ok := nodes[0].(*ast.StatementList).Stmts[0].(*ast.Assign); !ok {
		t.Error("Assign node expected, got " + reflect.TypeOf(nodes[0].(*ast.StatementList).Stmts[0]).String())
	}
	assign := nodes[0].(*ast.StatementList).Stmts[0].(*ast.Assign)
	if assign.Variable.(*ast.VariableReference).Name != "testVar" {
		t.Error("Incorrect variable name")
	}
	if assign.Variable.(*ast.VariableReference).Type.Kind() != ast.ComplexTypeStruct {
		t.Error("Expected variable base type to be complex struct")
	}
	if len(assign.Value.(*ast.StructLiteral).Type.Fields) != 1 {
		t.Error("Expected 1 field")
	}
	if assign.Value.(*ast.StructLiteral).Type.Fields[0].Ident != "Field" {
		t.Error("Field named incorrectly")
	}
	if assign.Value.(*ast.StructLiteral).Type.Fields[0].Type != ast.PrimitiveTypeString {
		t.Error("Field type string expected")
	}
	if len(assign.Value.(*ast.StructLiteral).Values) != 0 {
		t.Error("Expected no values")
	}
}

func TestLocalStructDeclaredThenInitializedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

		func test() {
			var testVar struct {
				Field string
			} = struct {
				Field string
			}{
				Field: "abc",
			}
		}`, t)
	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	nodes := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts
	if _, ok := nodes[0].(*ast.StatementList).Stmts[0].(*ast.Assign); !ok {
		t.Error("Assign node expected, got " + reflect.TypeOf(nodes[0].(*ast.StatementList).Stmts[0]).String())
	}
	assign := nodes[0].(*ast.StatementList).Stmts[0].(*ast.Assign)
	if assign.Variable.(*ast.VariableReference).Name != "testVar" {
		t.Error("Incorrect variable name")
	}
	if assign.Variable.(*ast.VariableReference).Type.Kind() != ast.ComplexTypeStruct {
		t.Error("Expected variable base type to be complex struct")
	}
	if len(assign.Value.(*ast.StructLiteral).Type.Fields) != 1 {
		t.Error("Expected 1 field")
	}
	if assign.Value.(*ast.StructLiteral).Type.Fields[0].Ident != "Field" {
		t.Error("Field named incorrectly")
	}
	if assign.Value.(*ast.StructLiteral).Type.Fields[0].Type != ast.PrimitiveTypeString {
		t.Error("Field type string expected")
	}
	if len(assign.Value.(*ast.StructLiteral).Values) != 1 {
		t.Error("Expected 1 value")
	}
	if assign.Value.(*ast.StructLiteral).Values["Field"].(*ast.StringLiteral).Str != "abc" {
		t.Error("Literal value incorrect")
	}
}

func TestLocalStructInitializedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func test(){
			testVar := struct{
				Field string
			}{
				Field: "abc",
			}
		}`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	nodes := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts
	if _, ok := nodes[0].(*ast.Assign); !ok {
		t.Error("Assign node expected, got " + reflect.TypeOf(nodes[0].(*ast.StatementList).Stmts[0]).String())
	}
	assign := nodes[0].(*ast.Assign)
	if assign.Variable.(*ast.VariableReference).Name != "testVar" {
		t.Error("Incorrect variable name")
	}
	if assign.Variable.(*ast.VariableReference).Type.Kind() != ast.ComplexTypeStruct {
		t.Error("Expected variable base type to be complex struct")
	}
	if len(assign.Value.(*ast.StructLiteral).Type.Fields) != 1 {
		t.Error("Expected 1 field")
	}
	if assign.Value.(*ast.StructLiteral).Type.Fields[0].Ident != "Field" {
		t.Error("Field named incorrectly")
	}
	if assign.Value.(*ast.StructLiteral).Type.Fields[0].Type != ast.PrimitiveTypeString {
		t.Error("Field type string expected")
	}
	if len(assign.Value.(*ast.StructLiteral).Values) != 1 {
		t.Error("Expected 1 value")
	}
	if assign.Value.(*ast.StructLiteral).Values["Field"].(*ast.StringLiteral).Str != "abc" {
		t.Error("Literal value incorrect")
	}
}

func TestGlobalIntArrayInvalidLenErrors(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

		var aa int
    var testVar [aa+1]int`, t)

	if len(context.Errors) != 1 {
		t.Error("Expected one error, got ", len(context.Errors))
	}
	if context.Errors[0].Class != NotStatic {
		t.Error("Expected error class to be NotStatic")
	}
}

func TestAssignReturnsASTStructure(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func test(){
      a = 3
    }

		var a int`, t)

	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.Assign); !ok {
		t.Error("Assign node expected")
	}

	if context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.Assign).NewLocal == true {
		t.Error("NewLocal flag truth expected")
	}

	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.Assign).Value.(*ast.IntegerLiteral); !ok {
		t.Error("IntegerLiteral node expected")
	}
	if context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.Assign).Value.(*ast.IntegerLiteral).Val != 3 {
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
	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 2 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	var s *ast.StatementList
	var ok bool
	if s, ok = context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.StatementList); !ok {
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

	if s, ok = context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[1].(*ast.StatementList); !ok {
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
	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	var s *ast.StatementList
	var ok bool
	if s, ok = context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.StatementList); !ok {
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
	if s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind() != ast.ComplexTypeArray {
		t.Error("ComplexTypeArray expected, got", s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind())
	}
	lenNode := s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Len
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
	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	var s *ast.StatementList
	var ok bool
	if s, ok = context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.StatementList); !ok {
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
	if s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind() != ast.ComplexTypeArray {
		t.Error("ComplexTypeArray expected, got", s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind())
	}
	subType := s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.SubType
	if subType == nil {
		t.Error("subType node expected")
	}
	if arraySubType, ok := subType.(ast.ArrayType); !ok {
		t.Error("ArrayType expected")
	} else {
		if arraySubType.BaseType() != ast.PrimitiveTypeInt {
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
	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	var s *ast.StatementList
	var ok bool
	if s, ok = context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.StatementList); !ok {
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
	if s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind() != ast.ComplexTypeArray {
		t.Error("ComplexTypeArray expected, got", s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind())
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
	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	var s *ast.StatementList
	var ok bool
	if s, ok = context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.StatementList); !ok {
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
	if s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind() != ast.ComplexTypeArray {
		t.Error("ComplexTypeArray expected, got", s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Type.Kind())
	}
	literals := s.Stmts[0].(*ast.Assign).Value.(*ast.ArrayLiteral).Literal
	if lit, ok := literals[0].(*ast.ArrayLiteral); !ok {
		t.Error("ArrayLiteral expected")
	} else {
		if lit.Type.BaseType() != ast.PrimitiveTypeInt {
			t.Error("Expected underlying type to be PrimitiveTypeInt")
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
	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts) != 1 {
		t.Error("Unexpected number of declarations in root node StatementList")
	}
	var s *ast.StatementList
	var ok bool
	if s, ok = context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.StatementList); !ok {
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
	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList, got ", reflect.TypeOf(context.Declarations[0].Type.(ast.FunctionType).Code))
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.IfStmt); !ok {
		t.Error("Expected root node for declaration to be IfStmt")
		t.FailNow()
	}
	ifNode := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.IfStmt)
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

func TestUnaryNotASTStructureGeneratedCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func test(k bool)bool{
			return !k
		}`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList, got ", reflect.TypeOf(context.Declarations[0].Type.(ast.FunctionType).Code))
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt); !ok {
		t.Error("Expected next node for declaration to be ReturnStmt")
		t.FailNow()
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.UnaryOp); !ok {
		t.Error("Expected next node for declaration to be UnaryOp")
		t.FailNow()
	}
}

func TestFunctionCallTranslatesCorrectly(t *testing.T) {
	_, context := setupTestGetAST(nil, `
		package test

		func test() int {
			mate(7, 4)
		}

		var mate int`, t)

	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList, got ", reflect.TypeOf(context.Declarations[0].Type.(ast.FunctionType).Code))
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.FunctionCall); !ok {
		t.Error("Expected next node for declaration to be FunctionCall")
		t.FailNow()
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.FunctionCall).Function.(*ast.VariableReference); !ok {
		t.Error("Expected function call to be VariableReference")
	}
	if len(context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.FunctionCall).Args) != 2 {
		t.Error("Expected 2 arguments")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.FunctionCall).Args[0].(*ast.IntegerLiteral); !ok {
		t.Error("Expected argument 1 to be type IntegerLiteral")
		t.FailNow()
	}
	if context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.FunctionCall).Args[0].(*ast.IntegerLiteral).Val != 7 {
		t.Error("Expected Argument 1 to be equal to 7")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.FunctionCall).Args[1].(*ast.IntegerLiteral); !ok {
		t.Error("Expected argument 2 to be type IntegerLiteral")
		t.FailNow()
	}
	if context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.FunctionCall).Args[1].(*ast.IntegerLiteral).Val != 4 {
		t.Error("Expected Argument 2 to be equal to 7")
	}
}

func TestEmptyReturnProducesCorrectASTStructure(t *testing.T) {
	_, context := setupTestGetAST(nil, `
    package test

    func test(){
			return
		}`, t)

	if len(context.Declarations) != 1 {
		t.Error("Unexpected number of declarations")
	}
	if context.Declarations[0].Ident != "test" {
		t.Error("Unexpected declaration name")
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList); !ok {
		t.Error("Expected root node for declaration to be StatementList, got ", reflect.TypeOf(context.Declarations[0].Type.(ast.FunctionType).Code))
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt); !ok {
		t.Error("Expected next node for declaration to be ReturnStmt")
		t.FailNow()
	}
	if _, ok := context.Declarations[0].Type.(ast.FunctionType).Code.(*ast.StatementList).Stmts[0].(*ast.ReturnStmt).Expr.(*ast.NilLiteral); !ok {
		t.Error("Expected next node for declaration to be NilLiteral")
		t.FailNow()
	}
}
