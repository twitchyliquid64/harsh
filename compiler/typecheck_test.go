package compiler

import (
	"testing"

	"github.com/twitchyliquid64/harsh/ast"
)

func TestTypecheckStatementList(t *testing.T) {
	node := &ast.StatementList{
		Stmts: []ast.Node{
			&ast.NilLiteral{},
		},
	}
	c := &TypecheckContext{}
	Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
}

func TestTypecheckString(t *testing.T) {
	node := &ast.StringLiteral{}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if ty != ast.PrimitiveTypeString {
		t.Error("Expected string type")
	}
}

func TestTypecheckInt(t *testing.T) {
	node := &ast.IntegerLiteral{}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if ty != ast.PrimitiveTypeInt {
		t.Error("Expected integer type")
	}
}

func TestTypecheckBool(t *testing.T) {
	node := &ast.BoolLiteral{}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if ty != ast.PrimitiveTypeBool {
		t.Error("Expected bool type")
	}
}

func TestTypecheckBinaryOpMatches(t *testing.T) {
	node := &ast.BinaryOp{
		LHS: &ast.IntegerLiteral{},
		RHS: &ast.IntegerLiteral{},
		Op:  ast.BinOpAdd,
	}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if ty != ast.PrimitiveTypeInt {
		t.Error("Expected int type")
	}
}

func TestTypecheckBinaryOpEquality(t *testing.T) {
	node := &ast.BinaryOp{
		LHS: &ast.IntegerLiteral{},
		RHS: &ast.IntegerLiteral{},
		Op:  ast.BinOpEquality,
	}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if ty != ast.PrimitiveTypeBool {
		t.Error("Expected bool type")
	}
}

func TestTypecheckBinaryOpLAnd(t *testing.T) {
	node := &ast.BinaryOp{
		LHS: &ast.IntegerLiteral{},
		RHS: &ast.IntegerLiteral{},
		Op:  ast.BinOpLAnd,
	}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if ty != ast.PrimitiveTypeBool {
		t.Error("Expected bool type")
	}
}

func TestTypecheckBinaryOpLOr(t *testing.T) {
	node := &ast.BinaryOp{
		LHS: &ast.IntegerLiteral{},
		RHS: &ast.IntegerLiteral{},
		Op:  ast.BinOpLOr,
	}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if ty != ast.PrimitiveTypeBool {
		t.Error("Expected bool type")
	}
}

func TestTypecheckBinaryOpMismatchCausesError(t *testing.T) {
	node := &ast.BinaryOp{
		LHS: &ast.IntegerLiteral{},
		RHS: &ast.BoolLiteral{},
		Op:  ast.BinOpAdd,
	}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 1 {
		t.Error("Type errors expected")
	}
	if c.Errors[0].Kind != TypeErrorIncompatibleTypesErr {
		t.Error("Expected incompatible types error")
	}
	if c.Errors[0].Msg != "Cannot perform binary operation + on operands with type int and bool" {
		t.Error("Incorrect error message, got: " + c.Errors[0].Msg)
	}
	if ty != ast.UnknownType {
		t.Error("Expected UnknownType")
	}
}

func TestBinaryOpWithOperandWithUnknownTypeReturnsUnknownType(t *testing.T) {
	node := &ast.BinaryOp{
		LHS: &ast.VariableReference{Name: "aa", Type: ast.UnknownType},
		RHS: &ast.IntegerLiteral{},
		Op:  ast.BinOpAdd,
	}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if ty != ast.UnknownType {
		t.Error("Expected UnknownType")
	}
}

func TestTypecheckVariabeReferenceWithKnownType(t *testing.T) {
	node := &ast.BinaryOp{
		LHS: &ast.VariableReference{Name: "aa", Type: ast.PrimitiveTypeInt},
		RHS: &ast.IntegerLiteral{},
		Op:  ast.BinOpAdd,
	}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if ty != ast.PrimitiveTypeInt {
		t.Error("Expected int type")
	}
}

func TestTypecheckVariabeReferenceWithUnknownTypeErrors(t *testing.T) {
	node := &ast.BinaryOp{
		LHS: &ast.VariableReference{Name: "aa"},
		RHS: &ast.IntegerLiteral{},
		Op:  ast.BinOpAdd,
	}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 1 {
		t.Error("1 Type error expected")
	}
	if c.Errors[0].Kind != TypeerrorInternalErr {
		t.Error("Expected error of type TypeerrorInternalErr")
	}
	if ty != ast.UnknownType {
		t.Error("Expected unknown type")
	}
}

func TestTypecheckAssignment(t *testing.T) {
	node := &ast.Assign{
		Variable: &ast.VariableReference{Name: "aa", Type: ast.PrimitiveTypeInt},
		Value:    &ast.IntegerLiteral{},
	}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if ty != ast.PrimitiveTypeInt {
		t.Error("Expected int type")
	}
}

func TestTypecheckAssignmentErrorsIfMismatchType(t *testing.T) {
	node := &ast.Assign{
		Variable: &ast.VariableReference{Name: "aa", Type: ast.PrimitiveTypeInt},
		Value:    &ast.StringLiteral{},
	}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 1 {
		t.Error("1 Type error expected")
	}
	if c.Errors[0].Kind != TypeErrorIncompatibleTypesErr {
		t.Error("Expected TypeErrorIncompatibleTypesErr")
	}
	if ty != ast.UnknownType {
		t.Error("Expected unknown type")
	}
}

func TestTypecheckAssignmentReturnsUnknownIfAnyUnknown(t *testing.T) {
	node := &ast.Assign{
		Variable: &ast.VariableReference{Name: "aa", Type: ast.UnknownType},
		Value:    &ast.StringLiteral{},
	}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("0 Type errors expected")
	}
	if ty != ast.UnknownType {
		t.Error("Expected unknown type")
	}
}

func TestTypecheckReturnWorksIfUnknown(t *testing.T) {
	node := &ast.ReturnStmt{
		Expr: &ast.VariableReference{Name: "aa", Type: ast.PrimitiveTypeInt},
	}
	c := &TypecheckContext{}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("0 Type errors expected")
	}
	if ty != ast.PrimitiveTypeInt {
		t.Error("Expected int type, got " + ty.String())
	}
}

func TestTypecheckReturnWorksIfKnown(t *testing.T) {
	node := &ast.ReturnStmt{
		Expr: &ast.VariableReference{Name: "aa", Type: ast.PrimitiveTypeInt},
	}
	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeInt}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("0 Type errors expected")
	}
	if ty != ast.PrimitiveTypeInt {
		t.Error("Expected int type, got " + ty.String())
	}
}

func TestTypecheckReturnErrorsIfTypeMismatch(t *testing.T) {
	node := &ast.ReturnStmt{
		Expr: &ast.VariableReference{Name: "aa", Type: ast.PrimitiveTypeInt},
	}
	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, node)
	if len(c.Errors) != 1 {
		t.Error("1 Type error expected")
	}
	if c.Errors[0].Kind != TypeErrorIncompatibleTypesErr {
		t.Error("Expected TypeErrorIncompatibleTypesErr")
	}
	if ty != ast.UnknownType {
		t.Error("Expected unknown type")
	}
}

func TestTypecheckSubscriptErrorsIfIndexNonInteger(t *testing.T) {
	node := &ast.Subscript{
		Subscript: &ast.VariableReference{Name: "aa", Type: ast.PrimitiveTypeString},
		Expr: &ast.VariableReference{Name: "aa",
			Type: ast.ArrayType{
				SubType: ast.PrimitiveTypeInt,
			},
		},
	}
	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, node)
	if len(c.Errors) != 1 {
		t.Error("1 Type error expected")
	}
	if c.Errors[0].Kind != TypeErrorIncompatibleTypesErr {
		t.Error("Expected TypeErrorIncompatibleTypesErr")
	}
	if ty != ast.UnknownType {
		t.Error("Expected unknown type")
	}
}

func TestTypecheckSubscriptErrorsIfExprNonArray(t *testing.T) {
	node := &ast.Subscript{
		Subscript: &ast.VariableReference{Name: "aa", Type: ast.PrimitiveTypeInt},
		Expr: &ast.VariableReference{Name: "aa",
			Type: ast.PrimitiveTypeInt,
		},
	}
	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, node)
	if len(c.Errors) != 1 {
		t.Error("1 Type error expected")
	}
	if c.Errors[0].Kind != TypeErrorIncompatibleTypesErr {
		t.Error("Expected TypeErrorIncompatibleTypesErr")
	}
	if ty != ast.UnknownType {
		t.Error("Expected unknown type")
	}
}

func TestTypecheckSubscriptWorks(t *testing.T) {
	node := &ast.Subscript{
		Subscript: &ast.VariableReference{Name: "aa", Type: ast.PrimitiveTypeInt},
		Expr: &ast.VariableReference{Name: "aa",
			Type: ast.ArrayType{
				SubType: ast.PrimitiveTypeBool,
			},
		},
	}
	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("0 Type errors expected")
	}
	if ty != ast.PrimitiveTypeBool {
		t.Error("Expected bool type")
	}
}

func TestTypecheckStructLiteralWorks(t *testing.T) {
	node := &ast.StructLiteral{
		Type: ast.StructType{
			Fields: []ast.NamedType{
				ast.NamedType{
					Ident: "Abc",
					Type:  ast.PrimitiveTypeInt,
				},
			},
		},
		Values: map[string]ast.Node{
			"Abc": &ast.IntegerLiteral{
				Val: 1234,
			},
		},
	}
	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("0 Type errors expected")
	}
	if ty.Kind() != ast.ComplexTypeStruct {
		t.Error("Expected struct type")
	}
}

func TestTypecheckStructLiteralErrorsOnFieldMismatch(t *testing.T) {
	node := &ast.StructLiteral{
		Type: ast.StructType{
			Fields: []ast.NamedType{
				ast.NamedType{
					Ident: "Abc",
					Type:  ast.PrimitiveTypeInt,
				},
			},
		},
		Values: map[string]ast.Node{
			"Abc": &ast.StringLiteral{
				Str: "1234",
			},
		},
	}
	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, node)
	if len(c.Errors) != 1 {
		t.Error("1 Type errors expected")
	}
	if ty.Kind() != ast.UnknownType {
		t.Error("Expected unknown type")
	}
}

func TestTypeEqualStructsReturnsTrue(t *testing.T) {
	l := ast.StructType{
		Fields: []ast.NamedType{
			ast.NamedType{
				Ident: "Abc",
				Type:  ast.PrimitiveTypeInt,
			},
		},
	}
	r := ast.StructType{
		Fields: []ast.NamedType{
			ast.NamedType{
				Ident: "Abc",
				Type:  ast.PrimitiveTypeInt,
			},
		},
	}
	if !TypeEqual(l, r) {
		t.Error("Expected types to be equal")
	}
}

func TestTypeEqualStructDifferentFieldTypesReturnsFalse(t *testing.T) {
	l := ast.StructType{
		Fields: []ast.NamedType{
			ast.NamedType{
				Ident: "Abc",
				Type:  ast.PrimitiveTypeInt,
			},
		},
	}
	r := ast.StructType{
		Fields: []ast.NamedType{
			ast.NamedType{
				Ident: "Abc",
				Type:  ast.PrimitiveTypeString,
			},
		},
	}
	if TypeEqual(l, r) {
		t.Error("Expected types to be different")
	}
}

func TestTypeEqualStructExtraFieldsReturnsFalse1(t *testing.T) {
	l := ast.StructType{
		Fields: []ast.NamedType{
			ast.NamedType{
				Ident: "Abc",
				Type:  ast.PrimitiveTypeInt,
			},
		},
	}
	r := ast.StructType{
		Fields: []ast.NamedType{
			ast.NamedType{
				Ident: "Abc",
				Type:  ast.PrimitiveTypeString,
			},
			ast.NamedType{
				Ident: "Abcdr",
				Type:  ast.PrimitiveTypeString,
			},
		},
	}
	if TypeEqual(l, r) {
		t.Error("Expected types to be different")
	}
}

func TestTypeEqualStructExtraFieldsReturnsFalse2(t *testing.T) {
	l := ast.StructType{
		Fields: []ast.NamedType{
			ast.NamedType{
				Ident: "Abc",
				Type:  ast.PrimitiveTypeInt,
			},
			ast.NamedType{
				Ident: "Abcdr",
				Type:  ast.PrimitiveTypeString,
			},
		},
	}
	r := ast.StructType{
		Fields: []ast.NamedType{
			ast.NamedType{
				Ident: "Abc",
				Type:  ast.PrimitiveTypeString,
			},
		},
	}
	if TypeEqual(l, r) {
		t.Error("Expected types to be different")
	}
}

func TestTypecheckArrayLiteralWorks(t *testing.T) {
	node := &ast.ArrayLiteral{
		Type: ast.ArrayType{
			SubType: ast.PrimitiveTypeInt,
		},
		Literal: []ast.Node{
			&ast.IntegerLiteral{Val: 23},
		},
	}
	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("0 Type errors expected")
	}
	if ty.Kind() != ast.ComplexTypeArray {
		t.Error("Expected array type")
	}
}

func TestTypecheckArrayLiteralErrorsOnFieldMismatch(t *testing.T) {
	node := &ast.ArrayLiteral{
		Type: ast.ArrayType{
			SubType: ast.PrimitiveTypeInt,
		},
		Literal: []ast.Node{
			&ast.StringLiteral{
				Str: "1234",
			},
		},
	}
	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, node)
	if len(c.Errors) != 1 {
		t.Error("1 Type errors expected")
	}
	if ty.Kind() != ast.UnknownType {
		t.Error("Expected unknown type")
	}
}

func TestTypeEqualArrayReturnsTrue(t *testing.T) {
	l := ast.ArrayType{
		SubType: ast.PrimitiveTypeBool,
	}
	r := ast.ArrayType{
		SubType: ast.PrimitiveTypeBool,
	}
	if !TypeEqual(l, r) {
		t.Error("Expected types to be equal")
	}
}

func TestTypeEqualArrayArrayReturnsTrue(t *testing.T) {
	l := ast.ArrayType{
		SubType: ast.ArrayType{
			SubType: ast.PrimitiveTypeBool,
		},
	}
	r := ast.ArrayType{
		SubType: ast.ArrayType{
			SubType: ast.PrimitiveTypeBool,
		},
	}
	if !TypeEqual(l, r) {
		t.Error("Expected types to be equal")
	}
}

func TestTypecheckNamedSelectorWorks(t *testing.T) {
	s := &ast.StructLiteral{
		Type: ast.StructType{
			Fields: []ast.NamedType{
				ast.NamedType{
					Ident: "Abc",
					Type:  ast.PrimitiveTypeInt,
				},
			},
		},
		Values: map[string]ast.Node{
			"Abc": &ast.IntegerLiteral{
				Val: 1234,
			},
		},
	}
	sel := &ast.NamedSelector{
		Name: "Abc",
		Expr: s,
	}

	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, sel)
	if len(c.Errors) != 0 {
		t.Error("0 Type errors expected")
	}
	if ty != ast.PrimitiveTypeInt {
		t.Error("Expected int type")
	}
}

func TestTypecheckNamedSelectorErrorsIfNotExists(t *testing.T) {
	s := &ast.StructLiteral{
		Type: ast.StructType{
			Fields: []ast.NamedType{
				ast.NamedType{
					Ident: "Abc",
					Type:  ast.PrimitiveTypeInt,
				},
			},
		},
		Values: map[string]ast.Node{
			"Abc": &ast.IntegerLiteral{
				Val: 1234,
			},
		},
	}
	sel := &ast.NamedSelector{
		Name: "doesnt exist bro",
		Expr: s,
	}

	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, sel)
	if len(c.Errors) != 1 {
		t.Error("1 Type error expected")
	}
	if c.Errors[0].Kind != TypeErrorNotFoundErr {
		t.Error("Expected TypeErrorNotFoundErr")
	}
	if ty != ast.UnknownType {
		t.Error("Expected UnknownType type")
	}
}

func TestTypecheckNamedSelectorErrorsIfWrongUpstreamType(t *testing.T) {
	sel := &ast.NamedSelector{
		Name: "doesnt exist bro",
		Expr: &ast.IntegerLiteral{
			Val: 1234,
		},
	}

	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, sel)
	if len(c.Errors) != 1 {
		t.Error("1 Type error expected")
	}
	if c.Errors[0].Kind != TypeErrorIncompatibleTypesErr {
		t.Error("Expected TypeErrorIncompatibleTypesErr")
	}
	if ty != ast.UnknownType {
		t.Error("Expected UnknownType type")
	}
}

func TestTypecheckIfStmtWorks(t *testing.T) {
	s := &ast.IfStmt{
		Conditional: &ast.BoolLiteral{},
		Code:        &ast.BoolLiteral{},
		Init:        &ast.BoolLiteral{},
		Else:        &ast.BoolLiteral{},
	}

	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, s)
	if len(c.Errors) != 0 {
		t.Error("0 Type errors expected")
	}
	if ty != ast.UnknownType {
		t.Error("Expected UnknownType")
	}
}

func TestTypecheckIfStmtErrorsOnNonBooleanConditional(t *testing.T) {
	s := &ast.IfStmt{
		Conditional: &ast.StringLiteral{},
		Code:        &ast.BoolLiteral{},
		Init:        &ast.BoolLiteral{},
		Else:        &ast.BoolLiteral{},
	}

	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	Typecheck(c, s)
	if len(c.Errors) != 1 {
		t.Error("1 Type errors expected")
	}
	if c.Errors[0].Kind != TypeErrorIncompatibleTypesErr {
		t.Error("Expected TypeErrorIncompatibleTypesErr")
	}
}

func TestTypecheckForStmtWorks(t *testing.T) {
	s := &ast.ForStmt{
		Conditional:   &ast.BoolLiteral{},
		Code:          &ast.BoolLiteral{},
		Init:          &ast.BoolLiteral{},
		PostIteration: &ast.BoolLiteral{},
	}

	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	ty := Typecheck(c, s)
	if len(c.Errors) != 0 {
		t.Error("0 Type errors expected")
	}
	if ty != ast.UnknownType {
		t.Error("Expected UnknownType")
	}
}

func TestTypecheckForStmtErrorsIfNotBooleanConditional(t *testing.T) {
	s := &ast.ForStmt{
		Conditional: &ast.IntegerLiteral{},
	}

	c := &TypecheckContext{ReturnType: ast.PrimitiveTypeBool}
	Typecheck(c, s)
	if len(c.Errors) != 1 {
		t.Error("1 Type errors expected")
	}
	if c.Errors[0].Kind != TypeErrorIncompatibleTypesErr {
		t.Error("Expected TypeErrorIncompatibleTypesErr")
	}
}

func TestTypeEqualFuncTypeReturnsTrue(t *testing.T) {
	l := ast.FunctionType{
		ReturnType: ast.PrimitiveTypeInt,
	}
	r := ast.FunctionType{
		ReturnType: ast.PrimitiveTypeInt,
	}
	if !TypeEqual(l, r) {
		t.Error("Expected types to be equal")
	}
}

func TestTypeEqualFuncTypeWithDiffReturnReturnsFalse(t *testing.T) {
	l := ast.FunctionType{
		ReturnType: ast.PrimitiveTypeInt,
	}
	r := ast.FunctionType{
		ReturnType: ast.PrimitiveTypeBool,
	}
	if TypeEqual(l, r) {
		t.Error("Expected types to not be equal")
	}
}

func TestTypeEqualFuncTypeWithDiffNumberParamsReturnsFalse(t *testing.T) {
	l := ast.FunctionType{
		ReturnType: ast.PrimitiveTypeInt,
		Parameters: []ast.TypeKind{ast.PrimitiveTypeString},
	}
	r := ast.FunctionType{
		ReturnType: ast.PrimitiveTypeInt,
	}
	if TypeEqual(l, r) {
		t.Error("Expected types to not be equal")
	}
}

func TestTypeEqualFuncTypeWithDiffParamsReturnsFalse(t *testing.T) {
	l := ast.FunctionType{
		ReturnType: ast.PrimitiveTypeInt,
		Parameters: []ast.TypeKind{ast.PrimitiveTypeString},
	}
	r := ast.FunctionType{
		ReturnType: ast.PrimitiveTypeInt,
		Parameters: []ast.TypeKind{ast.PrimitiveTypeInt},
	}
	if TypeEqual(l, r) {
		t.Error("Expected types to not be equal")
	}
}
