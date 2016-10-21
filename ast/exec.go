package ast

func (n *IntegerLiteral) Exec(context *ExecContext) Variant {
	return Variant{
		Type: PrimitiveType{
			Kind: PRIMITIVE_TYPE_INT,
		},
		Int: n.Val,
	}
}

func (n *StatementList) Exec(context *ExecContext) Variant {
	callingContext := (*context)
	newContext := callingContext
	newContext.IsFuncContext = false

	for _, node := range n.Stmts {
		v := node.Exec(&newContext)
		if v.IsReturn && context.IsFuncContext {
			v.IsReturn = false
			return v
		}
	}

	return Variant{
		Type: PrimitiveType{
			Kind: PRIMITIVE_TYPE_UNDEFINED,
		},
	}
}

func (n *ReturnStmt) Exec(context *ExecContext) Variant {
	v := n.Expr.Exec(context)
	v.IsReturn = true
	return v
}

func (n *BinaryOp) Exec(context *ExecContext) Variant {
	l := n.LHS.Exec(context)
	r := n.RHS.Exec(context)
	ret := Variant{
		Type: PrimitiveType{
			Kind: PRIMITIVE_TYPE_UNDEFINED,
		},
	}
	if l.Type.Kind == PRIMITIVE_TYPE_INT && r.Type.Kind == PRIMITIVE_TYPE_INT {
		ret.Type = PrimitiveType{
			Kind: PRIMITIVE_TYPE_INT,
		}
		switch n.Op {
		case BINOP_ADD:
			ret.Int = l.Int + r.Int
		case BINOP_SUB:
			ret.Int = l.Int - r.Int
		case BINOP_MUL:
			ret.Int = l.Int * r.Int
		case BINOP_DIV:
			ret.Int = l.Int / r.Int
		}
	}

	return ret
}
