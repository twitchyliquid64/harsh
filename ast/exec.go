package ast

func (n *IntegerLiteral) Exec(context *ExecContext) Variant {
	return Variant{
		Type: PrimitiveType{
			Kind: PRIMITIVE_TYPE_INT,
		},
		Int: n.Val,
	}
}

func (n *StringLiteral) Exec(context *ExecContext) Variant {
	return Variant{
		Type: PrimitiveType{
			Kind: PRIMITIVE_TYPE_STRING,
		},
		String: n.Str,
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
		case BINOP_MOD:
			ret.Int = l.Int % r.Int
		}
	} else if l.Type.Kind == PRIMITIVE_TYPE_STRING && r.Type.Kind == PRIMITIVE_TYPE_STRING {
		ret.Type = PrimitiveType{
			Kind: PRIMITIVE_TYPE_STRING,
		}
		switch n.Op {
		case BINOP_ADD:
			ret.String = l.String + r.String
			//TODO: Add default case which adds an error to the context
		}
	}

	return ret
}

func (n *VariableReference) Exec(context *ExecContext) Variant {
	if v, ok := context.FunctionNamespace[n.Name]; ok {
		return v
	}
	if context.GlobalNamespace != nil {
		if v, ok := context.GlobalNamespace[n.Name]; ok {
			return v
		}
	}
	return Variant{
		Type: PrimitiveType{
			Kind: PRIMITIVE_TYPE_UNDEFINED,
		},
	}
}

func (n *Assign) Exec(context *ExecContext) Variant {
	if n.NewLocal {
		context.FunctionNamespace.Save(n.Identifier, n.Value.Exec(context))
	} else {
		if _, ok := context.FunctionNamespace[n.Identifier]; ok && context.IsFuncContext {
			context.FunctionNamespace.Save(n.Identifier, n.Value.Exec(context))
		} else if _, ok := context.GlobalNamespace[n.Identifier]; ok {
			context.GlobalNamespace.Save(n.Identifier, n.Value.Exec(context))
		} else {
			if context.IsFuncContext {
				context.FunctionNamespace.Save(n.Identifier, n.Value.Exec(context))
			} else {
				context.GlobalNamespace.Save(n.Identifier, n.Value.Exec(context))
			}
		}
	}

	return Variant{
		Type: PrimitiveType{
			Kind: PRIMITIVE_TYPE_UNDEFINED,
		},
	}
}
