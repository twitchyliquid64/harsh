package ast

import "testing"

func TestIntLiteralExecReturnsCorrectValue(t *testing.T) {
	il := IntegerLiteral{
		Val: 493,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt return")
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
	if r.Type != PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool return")
	}
	if r.Bool != true {
		t.Error("Incorrect value")
	}
}

func TestNilLiteralExecReturnsCorrectValue(t *testing.T) {
	il := NilLiteral{}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined return")
	}
}

func TestEmptyArrayLiteralExecReturnsCorrectValue(t *testing.T) {
	il := ArrayLiteral{
		Type: ArrayType{
			SubType: PrimitiveTypeInt,
			Len: &IntegerLiteral{
				Val: 4,
			},
		},
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != ComplexTypeArray {
		t.Error("Expected ComplexTypeArray return")
	}
	if len(r.VectorData) != 4 {
		t.Error("Incorrect number of elements")
	}
	for i := 0; i < len(r.VectorData); i++ {
		if r.VectorData[i].Type != PrimitiveTypeUndefined {
			t.Error("Incorrect variant type at subscript", i)
		}
	}
}

func TestSimpleArrayLiteralExecReturnsCorrectValue(t *testing.T) {
	il := ArrayLiteral{
		Type: ArrayType{
			SubType: PrimitiveTypeInt,
			Len: &IntegerLiteral{
				Val: 3,
			},
		},
		Literal: []Node{
			&IntegerLiteral{Val: 33},
			&IntegerLiteral{Val: 88},
			&IntegerLiteral{Val: 232},
		},
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != ComplexTypeArray {
		t.Error("Expected ComplexTypeArray return")
	}
	if len(r.VectorData) != 3 {
		t.Error("Incorrect number of elements")
	}
	if r.VectorData[0].Int != 33 || r.VectorData[0].Type != PrimitiveTypeInt {
		t.Error("Incorrect type or value at subscript 0")
	}
	if r.VectorData[2].Int != 232 || r.VectorData[2].Type != PrimitiveTypeInt {
		t.Error("Incorrect type or value at subscript 2:", r.VectorData[2].Type.String(), r.VectorData[2].Int)
	}
}

func TestNestedArrayLiteralExecReturnsCorrectValue(t *testing.T) {
	il := ArrayLiteral{
		Type: ArrayType{
			SubType: ArrayType{
				SubType: PrimitiveTypeString,
				Len: &IntegerLiteral{
					Val: 1,
				},
			},
			Len: &IntegerLiteral{
				Val: 1,
			},
		},
		Literal: []Node{
			&ArrayLiteral{
				Type: ArrayType{
					SubType: PrimitiveTypeString,
					Len: &IntegerLiteral{
						Val: 1,
					},
				},
				Literal: []Node{
					&StringLiteral{
						Str: "LOLZ",
					},
				},
			},
		},
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != ComplexTypeArray {
		t.Error("Expected ComplexTypeArray return")
	}
	if len(r.VectorData) != 1 {
		t.Error("Incorrect number of elements")
	}
	if r.VectorData[0].Type != ComplexTypeArray {
		t.Error("Incorrect type", r.VectorData[0].Type.String())
	}
	if r.VectorData[0].VectorData[0].Type != PrimitiveTypeString {
		t.Error("Nested type incorrect")
	}
	if r.VectorData[0].VectorData[0].String != "LOLZ" {
		t.Error("Nested value incorrect")
	}
}

func TestStringLiteralExecReturnsCorrectValue(t *testing.T) {
	il := StringLiteral{
		Str: "Strdfjlkj_ fgklfjlgfjlkjlkj 'dd '",
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString return")
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
	if r.Type != PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt return")
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
		Op: BinOpAdd,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt return")
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
		Op: BinOpAdd,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined return")
	}
	if len(context.Errors) != 1 {
		t.Error("Errors expected")
		t.Fail()
	}
	if context.Errors[0].Class != TypeErr {
		t.Error("Expected error of type TypeErr")
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
		Op: BinOpAdd,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString return")
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
		Op: BinOpSub,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt return")
	}
	if r.Int != 489 {
		t.Error("Incorrect value")
	}
}

func TestBinaryOpModulusReturnsCorrectValue(t *testing.T) {
	il := BinaryOp{
		LHS: &IntegerLiteral{
			Val: 5,
		},
		RHS: &IntegerLiteral{
			Val: 4,
		},
		Op: BinOpMod,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt return")
	}
	if r.Int != 1 {
		t.Error("Incorrect value, got", r.Int)
	}
}

func TestBinaryOpDivisionReturnsCorrectValue(t *testing.T) {
	il := BinaryOp{
		LHS: &IntegerLiteral{
			Val: 6,
		},
		RHS: &IntegerLiteral{
			Val: 2,
		},
		Op: BinOpDiv,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt return")
	}
	if r.Int != 3 {
		t.Error("Incorrect value, got", r.Int)
	}
}

func TestBinaryOpIntEqualityReturnsCorrectValue(t *testing.T) {
	il := BinaryOp{
		LHS: &IntegerLiteral{
			Val: 6,
		},
		RHS: &IntegerLiteral{
			Val: 2,
		},
		Op: BinOpEquality,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool return")
	}
	if r.Bool != false {
		t.Error("Incorrect value, got", r.Bool)
	}
}

func TestBinaryOpBoolEqualityReturnsCorrectValue(t *testing.T) {
	il := BinaryOp{
		LHS: &BoolLiteral{
			Val: true,
		},
		RHS: &BoolLiteral{
			Val: false,
		},
		Op: BinOpEquality,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool return")
	}
	if r.Bool != false {
		t.Error("Incorrect value, got", r.Bool)
	}
	if len(context.Errors) != 0 {
		t.Error("Got errors, expected 0")
	}
}

func TestBinaryOpLogicalAndReturnsCorrectValue(t *testing.T) {
	il := BinaryOp{
		LHS: &BoolLiteral{
			Val: true,
		},
		RHS: &BoolLiteral{
			Val: false,
		},
		Op: BinOpLAnd,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool return")
	}
	if r.Bool != false {
		t.Error("Incorrect value, got", r.Int)
	}
}

func TestBinaryOpLogicalOrReturnsCorrectValue(t *testing.T) {
	il := BinaryOp{
		LHS: &BoolLiteral{
			Val: true,
		},
		RHS: &BoolLiteral{
			Val: false,
		},
		Op: BinOpLOr,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool return")
	}
	if r.Bool != true {
		t.Error("Incorrect value, got", r.Int)
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
		Op: BinOpMul,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt return")
	}
	if r.Int != 986 {
		t.Error("Incorrect value, got ", r.Int)
	}
}

func TestBinaryOpStringEquality(t *testing.T) {
	il := BinaryOp{
		LHS: &StringLiteral{
			Str: "ab",
		},
		RHS: &StringLiteral{
			Str: "cd",
		},
		Op: BinOpEquality,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool")
	}
	if len(context.Errors) != 0 {
		t.Error("Expected 0 errors, got", len(context.Errors))
	}
	if r.Bool != false {
		t.Error("Incorrect result")
	}
}

func TestBinaryOpInvalidStringOperationErrors(t *testing.T) {
	il := BinaryOp{
		LHS: &StringLiteral{
			Str: "ab",
		},
		RHS: &StringLiteral{
			Str: "cd",
		},
		Op: BinOpMod,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined")
	}
	if len(context.Errors) != 1 {
		t.Error("One error expected,", len(context.Errors))
	}
	if context.Errors[0].Class != TypeErr {
		t.Error("Incorrect error class")
	}
}

func TestBinaryOpInvalidBoolOperationErrors(t *testing.T) {
	il := BinaryOp{
		LHS: &BoolLiteral{
			Val: true,
		},
		RHS: &BoolLiteral{
			Val: true,
		},
		Op: BinOpMod,
	}
	context := ExecContext{}
	r := il.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined")
	}
	if len(context.Errors) != 1 {
		t.Error("One error expected,", len(context.Errors))
	}
	if context.Errors[0].Class != TypeErr {
		t.Error("Incorrect error class")
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
	if r.Type != PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt return, got " + r.Type.String())
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
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined return, got " + r.Type.String())
	}
}

func TestVariableReferenceReturnsFunctionParameter(t *testing.T) {
	vr := VariableReference{
		Name: "inInt",
	}
	context := ExecContext{
		IsFuncContext: true,
		FunctionNamespace: map[string]*Variant{
			"inInt": &Variant{
				Type: PrimitiveTypeInt,
				Int:  42,
			},
		},
	}

	r := vr.Exec(&context)
	if r.Type != PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt return, got " + r.Type.String())
	}
	if r.Int != 42 {
		t.Error("Incorrect value")
	}
}

func TestVariableReferenceReturnsGlobal(t *testing.T) {
	vr := VariableReference{
		Name: "glo",
	}
	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:   true,
		GlobalNamespace: ns,
	}
	context.GlobalNamespace.Save("glo", -42)

	r := vr.Exec(&context)
	if r.Type != PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt return, got " + r.Type.String())
	}
	if r.Int != -42 {
		t.Error("Incorrect value")
	}
}

func TestVariableReferenceReturnsParamFirst(t *testing.T) {
	vr := VariableReference{
		Name: "inInt",
	}
	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:   true,
		GlobalNamespace: ns,
		FunctionNamespace: map[string]*Variant{
			"inInt": &Variant{
				Type: PrimitiveTypeInt,
				Int:  1234,
			},
		},
	}
	context.GlobalNamespace.Save("inInt", -42)

	r := vr.Exec(&context)
	if r.Type != PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt return, got " + r.Type.String())
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
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined return, got " + r.Type.String())
	}
}

func TestLocalVariableWrite(t *testing.T) {
	ass := Assign{
		Variable: &VariableReference{
			Name: "testVar",
		},
		NewLocal: true,
		Value: &StringLiteral{
			Str: "abc",
		},
	}
	context := ExecContext{
		IsFuncContext:     true,
		FunctionNamespace: Namespace(map[string]*Variant{}),
		GlobalNamespace:   Namespace(map[string]*Variant{}),
	}

	ass.Exec(&context)
	v := context.FunctionNamespace["testVar"]
	if v.Type != PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString return, got " + v.Type.String())
	}
	if v.String != "abc" {
		t.Error("Incorrect value")
	}
	if _, ok := context.GlobalNamespace["testVar"]; ok {
		t.Error("Object should not be in global namespace")
	}
}

func TestLocalVariableWriteWhenAlreadyExists(t *testing.T) {
	ass := Assign{
		Variable: &VariableReference{
			Name: "testVar",
		},
		NewLocal: false,
		Value: &StringLiteral{
			Str: "abc",
		},
	}
	context := ExecContext{
		IsFuncContext: true,
		FunctionNamespace: Namespace(map[string]*Variant{
			"testVar": &Variant{Type: PrimitiveTypeString, String: "cba"},
		}),
		GlobalNamespace: Namespace(map[string]*Variant{}),
	}

	ass.Exec(&context)
	v := context.FunctionNamespace["testVar"]
	if v.Type != PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString return, got " + v.Type.String())
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
		Variable: &VariableReference{
			Name: "testVar",
		},
		Value: &StringLiteral{
			Str: "abc",
		},
	}
	context := ExecContext{
		IsFuncContext:     true,
		FunctionNamespace: Namespace(map[string]*Variant{}),
		GlobalNamespace:   Namespace(map[string]*Variant{}),
	}
	context.GlobalNamespace.Save("testVar", "asdsb")

	ass.Exec(&context)

	if _, ok := context.FunctionNamespace["testVar"]; ok {
		t.Error("Object should not be in function namespace")
	}

	v := context.GlobalNamespace["testVar"]
	if v.Type != PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString return, got " + v.Type.String())
	}
	if v.String != "abc" {
		t.Error("Incorrect value")
	}
}

func TestNewVariableAssignGoesToFunc(t *testing.T) {
	ass := Assign{
		Variable: &VariableReference{
			Name: "testVar",
		},
		Value: &StringLiteral{
			Str: "abc",
		},
	}
	context := ExecContext{
		IsFuncContext:     true,
		FunctionNamespace: Namespace(map[string]*Variant{}),
		GlobalNamespace:   Namespace(map[string]*Variant{}),
	}

	ass.Exec(&context)
	v := context.FunctionNamespace["testVar"]
	if v == nil {
		t.Error("Could not retrieve variable")
		t.FailNow()
	}
	if v.Type != PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString return, got " + v.Type.String())
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
		Variable: &VariableReference{
			Name: "testVar",
		},
		Value: &StringLiteral{
			Str: "abc",
		},
	}
	context := ExecContext{
		IsFuncContext:     false,
		FunctionNamespace: Namespace(map[string]*Variant{}),
		GlobalNamespace:   Namespace(map[string]*Variant{}),
	}

	ass.Exec(&context)
	v := context.GlobalNamespace["testVar"]
	if v.Type != PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString return, got " + v.Type.String())
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
			Variable: &VariableReference{
				Name: "testVar",
			},
			Value: &StringLiteral{
				Str: "main",
			},
		},
		Else: &Assign{
			Variable: &VariableReference{
				Name: "testVar",
			},
			Value: &StringLiteral{
				Str: "else",
			},
		},
	}
	context := ExecContext{
		IsFuncContext:     false,
		FunctionNamespace: Namespace(map[string]*Variant{}),
		GlobalNamespace:   Namespace(map[string]*Variant{}),
	}

	ifNode.Exec(&context)
	v := context.GlobalNamespace["testVar"]
	if v.Type != PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString return, got " + v.Type.String())
	}
	if v.String != "else" {
		t.Error("Incorrect value")
	}
}

func TestIfStatementTruthExecutesMainNotElse_alsoTestInit(t *testing.T) {
	ifNode := IfStmt{
		Conditional: &BoolLiteral{Val: true},
		Init: &Assign{
			Variable: &VariableReference{
				Name: "init",
			},
			Value: &StringLiteral{
				Str: "init has been run",
			},
		},
		Code: &Assign{
			Variable: &VariableReference{
				Name: "testVar",
			},
			Value: &StringLiteral{
				Str: "main",
			},
		},
		Else: &Assign{
			Variable: &VariableReference{
				Name: "testVar",
			},
			Value: &StringLiteral{
				Str: "else",
			},
		},
	}
	context := ExecContext{
		IsFuncContext:     false,
		FunctionNamespace: Namespace(map[string]*Variant{}),
		GlobalNamespace:   Namespace(map[string]*Variant{}),
	}

	ifNode.Exec(&context)
	v := context.GlobalNamespace["testVar"]
	if v.Type != PrimitiveTypeString {
		t.Error("Expected PrimitiveTypeString return, got " + v.Type.String())
	}
	if v.String != "main" {
		t.Error("Incorrect value")
	}

	if _, ok := context.GlobalNamespace["init"]; !ok {
		t.Error("Init AST node was not executed")
	}
}

func TestSubscriptLocalVariableReference(t *testing.T) {
	ass := &Assign{
		Variable: &Subscript{
			Subscript: &IntegerLiteral{
				Val: 0,
			},
			Expr: &VariableReference{
				Name: "testVar",
			},
		},
		Value: &IntegerLiteral{
			Val: 42,
		},
	}
	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:   true,
		GlobalNamespace: ns,
		FunctionNamespace: map[string]*Variant{
			"testVar": &Variant{
				Type: ComplexTypeArray,
				VectorData: []*Variant{
					&Variant{
						Type: PrimitiveTypeInt,
						Int:  11,
					},
				},
			},
		},
	}

	r := ass.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined return, got " + r.Type.String())
	}
	if v, ok := context.FunctionNamespace["testVar"]; !ok {
		t.Error("Could not read function variable")
	} else {
		if v.Type != ComplexTypeArray {
			t.Error("Type incorrect")
		}
		if v.VectorData[0].Type != PrimitiveTypeInt {
			t.Error("Element type incorrect")
		}
		if v.VectorData[0].Int != 42 {
			t.Error("Element value incorrect -", v.VectorData[0].Int)
		}
	}
}

func TestSubscriptGlobalVariableReference(t *testing.T) {
	ass := &Assign{
		Variable: &Subscript{
			Subscript: &IntegerLiteral{
				Val: 0,
			},
			Expr: &VariableReference{
				Name: "testVar",
			},
		},
		Value: &IntegerLiteral{
			Val: 42,
		},
	}
	ns := Namespace(map[string]*Variant{
		"testVar": &Variant{
			Type: ComplexTypeArray,
			VectorData: []*Variant{
				&Variant{
					Type: PrimitiveTypeInt,
					Int:  11,
				},
			},
		},
	})
	context := ExecContext{
		IsFuncContext:     true,
		GlobalNamespace:   ns,
		FunctionNamespace: map[string]*Variant{},
	}

	r := ass.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined return, got " + r.Type.String())
	}
	if v, ok := context.GlobalNamespace["testVar"]; !ok {
		t.Error("Could not read global variable")
	} else {
		if v.Type != ComplexTypeArray {
			t.Error("Type incorrect")
		}
		if v.VectorData[0].Type != PrimitiveTypeInt {
			t.Error("Element type incorrect")
		}
		if v.VectorData[0].Int != 42 {
			t.Error("Element value incorrect -", v.VectorData[0].Int)
		}
	}
}

func TestSubscriptVariableReferenceErrorsWhenNoneFound(t *testing.T) {
	ass := &Assign{
		Variable: &Subscript{
			Subscript: &IntegerLiteral{
				Val: 0,
			},
			Expr: &VariableReference{
				Name: "testVar",
			},
		},
		Value: &IntegerLiteral{
			Val: 42,
		},
	}
	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:     true,
		GlobalNamespace:   ns,
		FunctionNamespace: map[string]*Variant{},
	}

	r := ass.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined return, got " + r.Type.String())
	}
	if len(context.Errors) != 1 {
		t.Error("Expecting one error")
		t.Log(context.Errors[0].Error())
	}
	if context.Errors[0].Class != NotFoundErr {
		t.Error("Incorrect error class, got", context.Errors[0].Error())
	}
}

func TestUnaryNotOperation(t *testing.T) {
	op := &UnaryOp{
		Op: UnOpNot,
		Expr: &BoolLiteral{
			Val: true,
		},
	}

	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:     true,
		GlobalNamespace:   ns,
		FunctionNamespace: map[string]*Variant{},
	}

	r := op.Exec(&context)
	if r.Type != PrimitiveTypeBool {
		t.Error("Expected PrimitiveTypeBool")
	}
	if r.Bool != false {
		t.Error("Incorrect value")
	}
}

func TestUnaryBoolOperationWithIntErrors(t *testing.T) {
	op := &UnaryOp{
		Op: UnOpNot,
		Expr: &IntegerLiteral{
			Val: 6,
		},
	}

	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:     true,
		GlobalNamespace:   ns,
		FunctionNamespace: map[string]*Variant{},
	}

	r := op.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined")
	}
	if len(context.Errors) != 1 {
		t.Error("One error expected,", len(context.Errors))
	}
	if context.Errors[0].Class != TypeErr {
		t.Error("Incorrect error class")
	}
}

func TestArrayLiteralErrorsWhenNonIntSizeUsed(t *testing.T) {
	op := &ArrayLiteral{
		Type: ArrayType{
			Len: &StringLiteral{
				Str: "1",
			},
		},
	}

	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:     true,
		GlobalNamespace:   ns,
		FunctionNamespace: map[string]*Variant{},
	}

	r := op.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined")
	}
	if len(context.Errors) != 1 {
		t.Error("One error expected,", len(context.Errors))
	}
	if context.Errors[0].Class != TypeErr {
		t.Error("Incorrect error class")
	}
}

func TestArrayLiteralErrorsWhenLiteralLenMismatchToTypeLen(t *testing.T) {
	op := &ArrayLiteral{
		Type: ArrayType{
			Len: &IntegerLiteral{
				Val: 4,
			},
		},
		Literal: []Node{
			&IntegerLiteral{
				Val: 3,
			},
		},
	}

	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:     true,
		GlobalNamespace:   ns,
		FunctionNamespace: map[string]*Variant{},
	}

	r := op.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined")
	}
	if len(context.Errors) != 1 {
		t.Error("One error expected,", len(context.Errors))
	}
	if context.Errors[0].Class != BoundsErr {
		t.Error("Incorrect error class")
	}
}

func TestSubscriptErrorsWhenNonArrayBase(t *testing.T) {
	ass := &Assign{
		Variable: &Subscript{
			Subscript: &IntegerLiteral{
				Val: 0,
			},
			Expr: &IntegerLiteral{
				Val: 33,
			},
		},
		Value: &IntegerLiteral{
			Val: 42,
		},
	}
	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:     true,
		GlobalNamespace:   ns,
		FunctionNamespace: map[string]*Variant{},
	}

	ass.Exec(&context)
	if len(context.Errors) != 1 {
		t.Error("Expected 1 error, got", len(context.Errors))
	}
	if context.Errors[0].Class != TypeErr {
		t.Error("Incorrect error class")
	}
}

func TestSubscriptErrorsWithNonIntIndex(t *testing.T) {
	ass := &Assign{
		Variable: &Subscript{
			Subscript: &StringLiteral{
				Str: "1",
			},
			Expr: &VariableReference{
				Name: "testVar",
			},
		},
		Value: &IntegerLiteral{
			Val: 42,
		},
	}
	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:   true,
		GlobalNamespace: ns,
		FunctionNamespace: map[string]*Variant{
			"testVar": &Variant{
				Type: ComplexTypeArray,
				VectorData: []*Variant{
					&Variant{
						Type: PrimitiveTypeInt,
						Int:  11,
					},
				},
			},
		},
	}

	r := ass.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined return, got " + r.Type.String())
	}
	if len(context.Errors) != 1 {
		t.Error("Expected 1 error, got", len(context.Errors))
	}
	if context.Errors[0].Class != TypeErr {
		t.Error("Incorrect error class")
	}
}

func TestSubscriptErrorsWithOutOfBoundsSubscript(t *testing.T) {
	ass := &Assign{
		Variable: &Subscript{
			Subscript: &IntegerLiteral{
				Val: 1,
			},
			Expr: &VariableReference{
				Name: "testVar",
			},
		},
		Value: &IntegerLiteral{
			Val: 42,
		},
	}
	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:   true,
		GlobalNamespace: ns,
		FunctionNamespace: map[string]*Variant{
			"testVar": &Variant{
				Type: ComplexTypeArray,
				VectorData: []*Variant{
					&Variant{
						Type: PrimitiveTypeInt,
						Int:  11,
					},
				},
			},
		},
	}

	r := ass.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined return, got " + r.Type.String())
	}
	if len(context.Errors) != 1 {
		t.Error("Expected 1 error, got", len(context.Errors))
	}
	if context.Errors[0].Class != BoundsErr {
		t.Error("Incorrect error class")
	}
}

func TestStructLiteralWorksWhenGivenValue(t *testing.T) {
	op := &StructLiteral{
		Type: StructType{
			Fields: []NamedType{
				NamedType{
					Ident: "Lol",
					Type:  PrimitiveTypeInt,
				},
			},
		},
		Values: map[string]Node{
			"Lol": &IntegerLiteral{Val: 12},
		},
	}

	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:     true,
		GlobalNamespace:   ns,
		FunctionNamespace: map[string]*Variant{},
	}

	r := op.Exec(&context)
	if r.Type != ComplexTypeStruct {
		t.Error("Expected ComplexTypeStruct")
	}
	if len(context.Errors) != 0 {
		t.Error("0 errors expected,", len(context.Errors))
	}
	if lol, ok := r.NamedData["Lol"]; ok {
		if lol.Type != PrimitiveTypeInt {
			t.Error("Expected field Lol to be type PrimitiveTypeInt")
		}
		if lol.Int != 12 {
			t.Error("Expected integer value 12")
		}
	} else {
		t.Error("Field data missing in variant")
	}
}

func TestStructLiteralWorksWhenValueOmitted(t *testing.T) {
	op := &StructLiteral{
		Type: StructType{
			Fields: []NamedType{
				NamedType{
					Ident: "Lol",
					Type:  PrimitiveTypeInt,
				},
			},
		},
	}

	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:     true,
		GlobalNamespace:   ns,
		FunctionNamespace: map[string]*Variant{},
	}

	r := op.Exec(&context)
	if r.Type != ComplexTypeStruct {
		t.Error("Expected ComplexTypeStruct")
	}
	if len(context.Errors) != 0 {
		t.Error("0 errors expected,", len(context.Errors))
	}
	if lol, ok := r.NamedData["Lol"]; ok {
		if lol.Type != PrimitiveTypeInt {
			t.Error("Expected field Lol to be type PrimitiveTypeInt")
		}
		if lol.Int != 0 {
			t.Error("Expected integer value 0")
		}
	} else {
		t.Error("Field data missing in variant")
	}
}

func TestNamedSelectorWorksWithStructLiteral(t *testing.T) {
	s := &StructLiteral{
		Type: StructType{
			Fields: []NamedType{
				NamedType{
					Ident: "Lol",
					Type:  PrimitiveTypeInt,
				},
			},
		},
	}
	sel := &NamedSelector{
		Name: "Lol",
		Expr: s,
	}

	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:     true,
		GlobalNamespace:   ns,
		FunctionNamespace: map[string]*Variant{},
	}

	r := sel.Exec(&context)
	if r.Type != PrimitiveTypeInt {
		t.Error("Expected PrimitiveTypeInt")
	}
	if len(context.Errors) != 0 {
		t.Error("0 errors expected,", len(context.Errors))
	}
	if r.Int != 0 {
		t.Error("Expected default (0) value")
	}
}

func TestNamedSelectorErrorsIfNotFound(t *testing.T) {
	s := &StructLiteral{
		Type: StructType{
			Fields: []NamedType{
				NamedType{
					Ident: "Lol",
					Type:  PrimitiveTypeInt,
				},
			},
		},
	}
	sel := &NamedSelector{
		Name: "Bro",
		Expr: s,
	}

	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:     true,
		GlobalNamespace:   ns,
		FunctionNamespace: map[string]*Variant{},
	}

	r := sel.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined")
	}
	if len(context.Errors) != 1 {
		t.Error("1 errors expected, got", len(context.Errors))
	}
}

func TestNamedSelectorErrorsIfWrongUpstreamType(t *testing.T) {
	sel := &NamedSelector{
		Name: "Bro",
		Expr: &NilLiteral{},
	}

	ns := Namespace(map[string]*Variant{})
	context := ExecContext{
		IsFuncContext:     true,
		GlobalNamespace:   ns,
		FunctionNamespace: map[string]*Variant{},
	}

	r := sel.Exec(&context)
	if r.Type != PrimitiveTypeUndefined {
		t.Error("Expected PrimitiveTypeUndefined")
	}
	if len(context.Errors) != 1 {
		t.Error("1 errors expected, got", len(context.Errors))
	}
}
