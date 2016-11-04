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
	if c.Errors[0].Kind != TypeerrorIncompatibleTypesErr {
		t.Error("Expected incompatible types error")
	}
	if c.Errors[0].Msg != "Cannot perform binary operation + on operands with type int and bool" {
		t.Error("Incorrect error message, got: " + c.Errors[0].Msg)
	}
	if ty != ast.PrimitiveTypeUndefined {
		t.Error("Expected undefined type")
	}
}
