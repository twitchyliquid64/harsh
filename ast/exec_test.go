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
