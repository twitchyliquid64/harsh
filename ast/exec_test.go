package ast

import "testing"

func TestIntLiteralExecReturnsCorrectValue(t *testing.T) {
	il := IntegerLiteral{
		Val: 493,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT return")
	}
	if r.Int != 493 {
		t.Error("Incorrect value")
	}
}

func TestBoolLiteralExecReturnsCorrectValue(t *testing.T) {
	il := BoolLiteral{
		Val: true,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_BOOL {
		t.Error("Expected PRIMITIVE_TYPE_BOOL return")
	}
	if r.Bool != true {
		t.Error("Incorrect value")
	}
}

func TestStringLiteralExecReturnsCorrectValue(t *testing.T) {
	il := StringLiteral{
		Str: "Strdfjlkj_ fgklfjlgfjlkjlkj 'dd '",
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_STRING {
		t.Error("Expected PRIMITIVE_TYPE_STRING return")
	}
	if r.String != "Strdfjlkj_ fgklfjlgfjlkjlkj 'dd '" {
		t.Error("Incorrect value")
	}
}

func TestReturnReturnsUpstreamValue(t *testing.T) {
	il := ReturnStmt{
		Expr: &IntegerLiteral{
			Val: 493,
		},
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT return")
	}
	if r.Int != 493 {
		t.Error("Incorrect value")
	}
}

func TestBinaryOpAdditionReturnsCorrectValue(t *testing.T) {
	il := BinaryOp{
		LHS: &IntegerLiteral{
			Val: 493,
		},
		RHS: &IntegerLiteral{
			Val: 4,
		},
		Op: BINOP_ADD,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT return")
	}
	if r.Int != 497 {
		t.Error("Incorrect value")
	}
}

func TestBinaryOpInvalidOperandsErrors(t *testing.T) {
	il := BinaryOp{
		LHS: &IntegerLiteral{
			Val: 493,
		},
		RHS: &StringLiteral{
			Str: "4",
		},
		Op: BINOP_ADD,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_UNDEFINED {
		t.Error("Expected PRIMITIVE_TYPE_UNDEFINED return")
	}
	if len(context.Errors) != 1 {
		t.Error("Errors expected")
		t.Fail()
	}
	if context.Errors[0].Class != TYPE_ERR {
		t.Error("Expected error of type TYPE_ERR")
	}
	if context.Errors[0].Text != "Invalid types for operands: int and string" {
		t.Error("Got unexpected error text: " + context.Errors[0].Error())
	}
}

func TestBinaryOpConcatenationReturnsCorrectValue(t *testing.T) {
	il := BinaryOp{
		LHS: &StringLiteral{
			Str: "ab",
		},
		RHS: &StringLiteral{
			Str: "cd",
		},
		Op: BINOP_ADD,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_STRING {
		t.Error("Expected PRIMITIVE_TYPE_STRING return")
	}
	if r.String != "abcd" {
		t.Error("Incorrect value")
	}
}

func TestBinaryOpSubtractionReturnsCorrectValue(t *testing.T) {
	il := BinaryOp{
		LHS: &IntegerLiteral{
			Val: 493,
		},
		RHS: &IntegerLiteral{
			Val: 4,
		},
		Op: BINOP_SUB,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT return")
	}
	if r.Int != 489 {
		t.Error("Incorrect value")
	}
}

func TestBinaryOpMultiplicationReturnsCorrectValue(t *testing.T) {
	il := BinaryOp{
		LHS: &IntegerLiteral{
			Val: 493,
		},
		RHS: &IntegerLiteral{
			Val: 2,
		},
		Op: BINOP_MUL,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT return")
	}
	if r.Int != 986 {
		t.Error("Incorrect value, got ", r.Int)
	}
}

func TestStatementListReturnsShortCircuitValue(t *testing.T) {
	n := StatementList{
		Stmts: []Node{
			&IntegerLiteral{
				Val: 11,
			},
			&ReturnStmt{
				Expr: &IntegerLiteral{
					Val: 44,
				},
			},
			&IntegerLiteral{
				Val: 1,
			},
		},
	}

	context := ExecContext{
		IsFuncContext: true,
	}
	r := n.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT return, got " + r.Type.String())
	}
	if r.Int != 44 {
		t.Error("Incorrect value")
	}
}

func TestEmptyStatementListReturnsNoValue(t *testing.T) {
	n := StatementList{}

	context := ExecContext{
		IsFuncContext: true,
	}
	r := n.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_UNDEFINED {
		t.Error("Expected PRIMITIVE_TYPE_UNDEFINED return, got " + r.Type.String())
	}
}

func TestVariableReferenceReturnsFunctionParameter(t *testing.T) {
	vr := VariableReference{
		Name: "inInt",
	}
	context := ExecContext{
		IsFuncContext: true,
		FunctionNamespace: map[string]Variant{
			"inInt": Variant{
				Type: PrimitiveType{
					Kind: PRIMITIVE_TYPE_INT,
				},
				Int: 42,
			},
		},
	}

	r := vr.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT return, got " + r.Type.String())
	}
	if r.Int != 42 {
		t.Error("Incorrect value")
	}
}

func TestVariableReferenceReturnsGlobal(t *testing.T) {
	vr := VariableReference{
		Name: "glo",
	}
	ns := Namespace(map[string]Variant{})
	context := ExecContext{
		IsFuncContext:   true,
		GlobalNamespace: ns,
	}
	context.GlobalNamespace.Save("glo", -42)

	r := vr.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT return, got " + r.Type.String())
	}
	if r.Int != -42 {
		t.Error("Incorrect value")
	}
}

func TestVariableReferenceReturnsParamFirst(t *testing.T) {
	vr := VariableReference{
		Name: "inInt",
	}
	ns := Namespace(map[string]Variant{})
	context := ExecContext{
		IsFuncContext:   true,
		GlobalNamespace: ns,
		FunctionNamespace: map[string]Variant{
			"inInt": Variant{
				Type: PrimitiveType{
					Kind: PRIMITIVE_TYPE_INT,
				},
				Int: 1234,
			},
		},
	}
	context.GlobalNamespace.Save("inInt", -42)

	r := vr.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_INT {
		t.Error("Expected PRIMITIVE_TYPE_INT return, got " + r.Type.String())
	}
	if r.Int != 1234 {
		t.Error("Incorrect value")
	}
}

func TestVariableReferenceReturnsUndefinedWhenNoMatch(t *testing.T) {
	vr := VariableReference{
		Name: "inInt",
	}
	context := ExecContext{
		IsFuncContext: true,
	}

	r := vr.Exec(&context)
	if r.Type.Kind != PRIMITIVE_TYPE_UNDEFINED {
		t.Error("Expected PRIMITIVE_TYPE_UNDEFINED return, got " + r.Type.String())
	}
}

func TestLocalVariableWrite(t *testing.T) {
	ass := Assign{
		Identifier: "testVar",
		NewLocal:   true,
		Value: &StringLiteral{
			Str: "abc",
		},
	}
	context := ExecContext{
		IsFuncContext:     true,
		FunctionNamespace: Namespace(map[string]Variant{}),
		GlobalNamespace:   Namespace(map[string]Variant{}),
	}

	ass.Exec(&context)
	v := context.FunctionNamespace["testVar"]
	if v.Type.Kind != PRIMITIVE_TYPE_STRING {
		t.Error("Expected PRIMITIVE_TYPE_STRING return, got " + v.Type.String())
	}
	if v.String != "abc" {
		t.Error("Incorrect value")
	}
	if _, ok := context.GlobalNamespace["testVar"]; ok {
		t.Error("Object should not be in global namespace")
	}
}

func TestGlobalVariableWrite(t *testing.T) {
	ass := Assign{
		Identifier: "testVar",
		Value: &StringLiteral{
			Str: "abc",
		},
	}
	context := ExecContext{
		IsFuncContext:     true,
		FunctionNamespace: Namespace(map[string]Variant{}),
		GlobalNamespace:   Namespace(map[string]Variant{}),
	}
	context.GlobalNamespace.Save("testVar", "asdsb")

	ass.Exec(&context)

	if _, ok := context.FunctionNamespace["testVar"]; ok {
		t.Error("Object should not be in function namespace")
	}

	v := context.GlobalNamespace["testVar"]
	if v.Type.Kind != PRIMITIVE_TYPE_STRING {
		t.Error("Expected PRIMITIVE_TYPE_STRING return, got " + v.Type.String())
	}
	if v.String != "abc" {
		t.Error("Incorrect value")
	}
}

func TestNewVariableAssignGoesToFunc(t *testing.T) {
	ass := Assign{
		Identifier: "testVar",
		Value: &StringLiteral{
			Str: "abc",
		},
	}
	context := ExecContext{
		IsFuncContext:     true,
		FunctionNamespace: Namespace(map[string]Variant{}),
		GlobalNamespace:   Namespace(map[string]Variant{}),
	}

	ass.Exec(&context)
	v := context.FunctionNamespace["testVar"]
	if v.Type.Kind != PRIMITIVE_TYPE_STRING {
		t.Error("Expected PRIMITIVE_TYPE_STRING return, got " + v.Type.String())
	}
	if v.String != "abc" {
		t.Error("Incorrect value")
	}
	if _, ok := context.GlobalNamespace["testVar"]; ok {
		t.Error("Object should not be in global namespace")
	}
}

func TestNewVariableAssignGoesToGlobalWhenNotFuncContext(t *testing.T) {
	ass := Assign{
		Identifier: "testVar",
		Value: &StringLiteral{
			Str: "abc",
		},
	}
	context := ExecContext{
		IsFuncContext:     false,
		FunctionNamespace: Namespace(map[string]Variant{}),
		GlobalNamespace:   Namespace(map[string]Variant{}),
	}

	ass.Exec(&context)
	v := context.GlobalNamespace["testVar"]
	if v.Type.Kind != PRIMITIVE_TYPE_STRING {
		t.Error("Expected PRIMITIVE_TYPE_STRING return, got " + v.Type.String())
	}
	if v.String != "abc" {
		t.Error("Incorrect value")
	}
	if _, ok := context.FunctionNamespace["testVar"]; ok {
		t.Error("Object should not be in functional namespace")
	}
}

func TestIfStatementTruthExecutesElseNotMain(t *testing.T) {
	ifNode := IfStmt{
		Conditional: &BoolLiteral{},
		Code: &Assign{
			Identifier: "testVar",
			Value: &StringLiteral{
				Str: "main",
			},
		},
		Else: &Assign{
			Identifier: "testVar",
			Value: &StringLiteral{
				Str: "else",
			},
		},
	}
	context := ExecContext{
		IsFuncContext:     false,
		FunctionNamespace: Namespace(map[string]Variant{}),
		GlobalNamespace:   Namespace(map[string]Variant{}),
	}

	ifNode.Exec(&context)
	v := context.GlobalNamespace["testVar"]
	if v.Type.Kind != PRIMITIVE_TYPE_STRING {
		t.Error("Expected PRIMITIVE_TYPE_STRING return, got " + v.Type.String())
	}
	if v.String != "else" {
		t.Error("Incorrect value")
	}
}

func TestIfStatementTruthExecutesMainNotElse_alsoTestInit(t *testing.T) {
	ifNode := IfStmt{
		Conditional: &BoolLiteral{Val: true},
		Init: &Assign{
			Identifier: "init",
			Value: &StringLiteral{
				Str: "init has been run",
			},
		},
		Code: &Assign{
			Identifier: "testVar",
			Value: &StringLiteral{
				Str: "main",
			},
		},
		Else: &Assign{
			Identifier: "testVar",
			Value: &StringLiteral{
				Str: "else",
			},
		},
	}
	context := ExecContext{
		IsFuncContext:     false,
		FunctionNamespace: Namespace(map[string]Variant{}),
		GlobalNamespace:   Namespace(map[string]Variant{}),
	}

	ifNode.Exec(&context)
	v := context.GlobalNamespace["testVar"]
	if v.Type.Kind != PRIMITIVE_TYPE_STRING {
		t.Error("Expected PRIMITIVE_TYPE_STRING return, got " + v.Type.String())
	}
	if v.String != "main" {
		t.Error("Incorrect value")
	}

	if _, ok := context.GlobalNamespace["init"]; !ok {
		t.Error("Init AST node was not executed")
	}
}
