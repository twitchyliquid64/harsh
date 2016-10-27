package ast

func (n *IntegerLiteral) Exec(context *ExecContext) *Variant {
	return &Variant{
		Type: PRIMITIVE_TYPE_INT,
		Int:  n.Val,
	}
}

func (n *BoolLiteral) Exec(context *ExecContext) *Variant {
	return &Variant{
		Type: PRIMITIVE_TYPE_BOOL,
		Bool: n.Val,
	}
}

func (n *StringLiteral) Exec(context *ExecContext) *Variant {
	return &Variant{
		Type:   PRIMITIVE_TYPE_STRING,
		String: n.Str,
	}
}

func (n *ArrayLiteral) Exec(context *ExecContext) *Variant {
	sizeNode := n.Type.(ArrayType).Len.Exec(context)
	if sizeNode.Type != PRIMITIVE_TYPE_INT {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TYPE_ERR,
			CreatingNode: n,
			Text:         "Non-integer len used for array",
		})
		return &Variant{Type: PRIMITIVE_TYPE_UNDEFINED}
	}

	var values []*Variant = make([]*Variant, sizeNode.Int)
	if len(values) != len(n.Literal) && len(n.Literal) != 0 {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        BOUNDS_ERR,
			CreatingNode: n,
			Text:         "Literal used in array assignment does not match the size of the underlying array",
		})
		return &Variant{Type: PRIMITIVE_TYPE_UNDEFINED}
	}

	var i int
	for ; i < len(n.Literal); i++ {
		values[i] = n.Literal[i].Exec(context)
	}
	for ; i < len(values); i++ {
		values[i] = &Variant{Type: PRIMITIVE_TYPE_UNDEFINED}
	}

	return &Variant{
		Type:       COMPLEX_TYPE_ARRAY,
		VectorData: values,
	}
}

func (n *StatementList) Exec(context *ExecContext) *Variant {
	callingContext := (*context)
	newContext := callingContext
	newContext.IsFuncContext = false

	for _, node := range n.Stmts {
		v := node.Exec(&newContext)
		if v.IsReturn {
			for _, err := range newContext.Errors {
				context.Errors = append(context.Errors, err)
			}
			return v
		}
	}

	for _, err := range newContext.Errors {
		context.Errors = append(context.Errors, err)
	}
	return &Variant{
		Type: PRIMITIVE_TYPE_UNDEFINED,
	}
}

func (n *ReturnStmt) Exec(context *ExecContext) *Variant {
	v := n.Expr.Exec(context)
	temp := *v
	temp.IsReturn = true
	return &temp
}

func (n *BinaryOp) Exec(context *ExecContext) *Variant {
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
		case BINOP_EQUALITY:
			ret.Type = PRIMITIVE_TYPE_BOOL
			ret.Bool = l.Int == r.Int
		}
	} else if l.Type == PRIMITIVE_TYPE_STRING && r.Type == PRIMITIVE_TYPE_STRING {
		ret.Type = PRIMITIVE_TYPE_STRING
		switch n.Op {
		case BINOP_ADD:
			ret.String = l.String + r.String
		case BINOP_EQUALITY:
			ret.Type = PRIMITIVE_TYPE_BOOL
			ret.Bool = l.String == r.String
		default:
			context.Errors = append(context.Errors, ExecutionError{
				Class:        TYPE_ERR,
				CreatingNode: n,
				Text:         "Invalid operation for string operands: " + n.Op.String(),
			})
		}
	} else if l.Type == PRIMITIVE_TYPE_BOOL && r.Type == PRIMITIVE_TYPE_BOOL {
		ret.Type = PRIMITIVE_TYPE_BOOL
		switch n.Op {
		case BINOP_EQUALITY:
			ret.Bool = l.Bool && r.Bool
		case BINOP_LAND:
			ret.Bool = l.Bool && r.Bool
		case BINOP_LOR:
			ret.Bool = l.Bool || r.Bool
		default:
			context.Errors = append(context.Errors, ExecutionError{
				Class:        TYPE_ERR,
				CreatingNode: n,
				Text:         "Invalid operation for boolean operands: " + n.Op.String(),
			})
		}
	} else {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TYPE_ERR,
			CreatingNode: n,
			Text:         "Invalid types for operands: " + l.Type.String() + " and " + r.Type.String(),
		})

	}

	return &ret
}

func (n *VariableReference) Exec(context *ExecContext) *Variant {
	if v, ok := context.FunctionNamespace[n.Name]; ok {
		return v
	}
	if context.GlobalNamespace != nil {
		if v, ok := context.GlobalNamespace[n.Name]; ok {
			return v
		}
	}
	return &Variant{
		Type: PRIMITIVE_TYPE_UNDEFINED,
		VariableReferenceFailed: true,
	}
}

func (n *Assign) Exec(context *ExecContext) *Variant {
	variable := n.Variable.Exec(context)
	v := n.Value.Exec(context)
	if ident, ok := n.Variable.(*VariableReference); ok {
		if n.NewLocal || v.VariableReferenceFailed {
			context.FunctionNamespace.Save(ident.Name, v)
		} else {
			if _, ok := context.FunctionNamespace[ident.Name]; ok && context.IsFuncContext {
				context.FunctionNamespace.Save(ident.Name, v)
			} else if _, ok := context.GlobalNamespace[ident.Name]; ok {
				context.GlobalNamespace.Save(ident.Name, v)
			} else {
				if context.IsFuncContext {
					context.FunctionNamespace.Save(ident.Name, v)
				} else {
					context.GlobalNamespace.Save(ident.Name, v)
				}
			}
		}
	} else {
		newValue := *v
		newValue.IsReturn = false
		*variable = newValue
	}

	return &Variant{
		Type: PRIMITIVE_TYPE_UNDEFINED,
	}
}

func (n *IfStmt) Exec(context *ExecContext) *Variant {
	if n.Init != nil {
		n.Init.Exec(context)
	}

	conditionResult := n.Conditional.Exec(context)
	if conditionResult.Type == PRIMITIVE_TYPE_BOOL && conditionResult.Bool {
		return n.Code.Exec(context)
	} else if n.Else != nil {
		return n.Else.Exec(context)
	}

	return &Variant{
		Type: PRIMITIVE_TYPE_UNDEFINED,
	}
}

func (n *UnaryOp) Exec(context *ExecContext) *Variant {
	upper := n.Expr.Exec(context)
	if upper.Type == PRIMITIVE_TYPE_BOOL {
		switch n.Op {
		case UNOP_NOT:
			return &Variant{
				Type: PRIMITIVE_TYPE_BOOL,
				Bool: !upper.Bool,
			}
		default:
			context.Errors = append(context.Errors, ExecutionError{
				Class:        TYPE_ERR,
				CreatingNode: n,
				Text:         "Cannot perform unary operation on " + upper.Type.String(),
			})
		}
	} else {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TYPE_ERR,
			CreatingNode: n,
			Text:         "Cannot perform unary operations on type " + upper.Type.String(),
		})
	}
	return &Variant{
		Type: PRIMITIVE_TYPE_UNDEFINED,
	}
}

func (n *Subscript) Exec(context *ExecContext) *Variant {
	baseVar := n.Expr.Exec(context)
	subscript := n.Subscript.Exec(context)
	if baseVar.Type != COMPLEX_TYPE_ARRAY {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TYPE_ERR,
			CreatingNode: n,
			Text:         "Cannot perform subscript operation on type " + baseVar.Type.String(),
		})
		return &Variant{
			Type: PRIMITIVE_TYPE_UNDEFINED,
		}
	}
	if subscript.Type != PRIMITIVE_TYPE_INT {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TYPE_ERR,
			CreatingNode: n,
			Text:         "Cannot perform subscript operation on type " + baseVar.Type.String(),
		})
		return &Variant{
			Type: PRIMITIVE_TYPE_UNDEFINED,
		}
	}
	if int(subscript.Int) >= len(baseVar.VectorData) {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        BOUNDS_ERR,
			CreatingNode: n,
			Text:         "Subscript out of bounds",
		})
		return &Variant{
			Type: PRIMITIVE_TYPE_UNDEFINED,
		}
	}

	return baseVar.VectorData[subscript.Int]
}
