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
	type_ := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if type_ != ast.PRIMITIVE_TYPE_STRING {
		t.Error("Expected string type")
	}
}

func TestTypecheckInt(t *testing.T) {
	node := &ast.IntegerLiteral{}
	c := &TypecheckContext{}
	type_ := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if type_ != ast.PRIMITIVE_TYPE_INT {
		t.Error("Expected integer type")
	}
}

func TestTypecheckBool(t *testing.T) {
	node := &ast.BoolLiteral{}
	c := &TypecheckContext{}
	type_ := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if type_ != ast.PRIMITIVE_TYPE_BOOL {
		t.Error("Expected bool type")
	}
}

func TestTypecheckBinaryOpMatches(t *testing.T) {
	node := &ast.BinaryOp{
		LHS: &ast.IntegerLiteral{},
		RHS: &ast.IntegerLiteral{},
		Op:  ast.BINOP_ADD,
	}
	c := &TypecheckContext{}
	type_ := Typecheck(c, node)
	if len(c.Errors) != 0 {
		t.Error("Type errors not expected")
	}
	if type_ != ast.PRIMITIVE_TYPE_INT {
		t.Error("Expected int type")
	}
}

func TestTypecheckBinaryOpMismatchCausesError(t *testing.T) {
	node := &ast.BinaryOp{
		LHS: &ast.IntegerLiteral{},
		RHS: &ast.BoolLiteral{},
		Op:  ast.BINOP_ADD,
	}
	c := &TypecheckContext{}
	type_ := Typecheck(c, node)
	if len(c.Errors) != 1 {
		t.Error("Type errors expected")
	}
	if c.Errors[0].Kind != TYPEERROR_INCOMPATIBLE_TYPES_ERR {
		t.Error("Expected incompatible types error")
	}
	if type_ != ast.PRIMITIVE_TYPE_UNDEFINED {
		t.Error("Expected undefined type")
	}
}
