package ast

func (n *IntegerLiteral) Exec(context *ExecContext) Variant {
	return Variant{
		Type: PRIMITIVE_TYPE_INT,
		Int:  n.Val,
	}
}

func (n *BoolLiteral) Exec(context *ExecContext) Variant {
	return Variant{
		Type: PRIMITIVE_TYPE_BOOL,
		Bool: n.Val,
	}
}

func (n *StringLiteral) Exec(context *ExecContext) Variant {
	return Variant{
		Type:   PRIMITIVE_TYPE_STRING,
		String: n.Str,
	}
}

func (n *ArrayLiteral) Exec(context *ExecContext) Variant {
	return Variant{
		Type: n.Type,
	}
}

func (n *StatementList) Exec(context *ExecContext) Variant {
	callingContext := (*context)
	newContext := callingContext
	newContext.IsFuncContext = false

	for _, node := range n.Stmts {
		v := node.Exec(&newContext)
		if v.IsReturn {
			return v
		}
	}

	return Variant{
		Type: PRIMITIVE_TYPE_UNDEFINED,
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
		Type: PRIMITIVE_TYPE_UNDEFINED,
	}
	if l.Type == PRIMITIVE_TYPE_INT && r.Type == PRIMITIVE_TYPE_INT {
		ret.Type = PRIMITIVE_TYPE_INT
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
	} else if l.Type == PRIMITIVE_TYPE_STRING && r.Type == PRIMITIVE_TYPE_STRING {
		ret.Type = PRIMITIVE_TYPE_STRING
		switch n.Op {
		case BINOP_ADD:
			ret.String = l.String + r.String
			//TODO: Add default case which adds an error to the context
		}
	} else {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TYPE_ERR,
			CreatingNode: n,
			Text:         "Invalid types for operands: " + l.Type.String() + " and " + r.Type.String(),
		})
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
		Type: PRIMITIVE_TYPE_UNDEFINED,
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
		Type: PRIMITIVE_TYPE_UNDEFINED,
	}
}

func (n *IfStmt) Exec(context *ExecContext) Variant {
	if n.Init != nil {
		n.Init.Exec(context)
	}

	conditionResult := n.Conditional.Exec(context)
	if conditionResult.Type == PRIMITIVE_TYPE_BOOL && conditionResult.Bool {
		return n.Code.Exec(context)
	} else if n.Else != nil {
		return n.Else.Exec(context)
	}

	return Variant{
		Type: PRIMITIVE_TYPE_UNDEFINED,
	}
}
