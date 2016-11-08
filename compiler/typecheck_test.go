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
