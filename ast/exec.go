package ast

import "fmt"

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
func (n *IntegerLiteral) Exec(context *ExecContext) *Variant {
	return &Variant{
		Type: PrimitiveTypeInt,
		Int:  n.Val,
	}
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
func (n *BoolLiteral) Exec(context *ExecContext) *Variant {
	return &Variant{
		Type: PrimitiveTypeBool,
		Bool: n.Val,
	}
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
func (n *StringLiteral) Exec(context *ExecContext) *Variant {
	return &Variant{
		Type:   PrimitiveTypeString,
		String: n.Str,
	}
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
func (n *NilLiteral) Exec(context *ExecContext) *Variant {
	return &Variant{
		Type: PrimitiveTypeUndefined,
	}
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
func (n *ArrayLiteral) Exec(context *ExecContext) *Variant {
	sizeNode := n.Type.Len.Exec(context)
	if sizeNode.Type != PrimitiveTypeInt {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TypeErr,
			CreatingNode: n,
			Text:         "Non-integer len used for array",
		})
		return &Variant{Type: PrimitiveTypeUndefined}
	}

	var values = make([]*Variant, sizeNode.Int)
	if len(values) != len(n.Literal) && len(n.Literal) != 0 {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        BoundsErr,
			CreatingNode: n,
			Text:         "Literal used in array assignment does not match the size of the underlying array",
		})
		return &Variant{Type: PrimitiveTypeUndefined}
	}

	var i int
	for ; i < len(n.Literal); i++ {
		values[i] = n.Literal[i].Exec(context)
	}
	for ; i < len(values); i++ {
		values[i] = &Variant{Type: PrimitiveTypeUndefined}
	}

	return &Variant{
		Type:       ComplexTypeArray,
		VectorData: values,
	}
}

// Exec resolves the values for the literals specified (if any).
func (n *StructLiteral) Exec(context *ExecContext) *Variant {
	o := &Variant{
		Type:      ComplexTypeStruct,
		NamedData: map[string]*Variant{},
	}
	for _, field := range n.Type.Fields {
		if n.Values != nil && n.Values[field.Ident] != nil {
			o.NamedData[field.Ident] = n.Values[field.Ident].Exec(context)
		} else {
			var err error
			o.NamedData[field.Ident], err = DefaultVariantValue(field.Type)
			if err != nil {
				context.Errors = append(context.Errors, ExecutionError{
					Class:        InternalErr,
					CreatingNode: n,
					Text:         "Failed to create default value to populate field '" + field.Ident + "' with type: " + field.Type.String(),
				})
			}
		}
	}
	return o
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
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
		Type: PrimitiveTypeUndefined,
	}
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
func (n *ReturnStmt) Exec(context *ExecContext) *Variant {
	v := n.Expr.Exec(context)
	temp := *v
	temp.IsReturn = true
	return &temp
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
func (n *BinaryOp) Exec(context *ExecContext) *Variant {
	l := n.LHS.Exec(context)
	r := n.RHS.Exec(context)
	ret := Variant{
		Type: PrimitiveTypeUndefined,
	}
	if l.Type == PrimitiveTypeInt && r.Type == PrimitiveTypeInt {
		ret.Type = PrimitiveTypeInt
		switch n.Op {
		case BinOpAdd:
			ret.Int = l.Int + r.Int
		case BinOpSub:
			ret.Int = l.Int - r.Int
		case BinOpMul:
			ret.Int = l.Int * r.Int
		case BinOpDiv:
			ret.Int = l.Int / r.Int
		case BinOpMod:
			ret.Int = l.Int % r.Int
		case BinOpEquality:
			ret.Type = PrimitiveTypeBool
			ret.Bool = l.Int == r.Int
		case BinOpNotEquality:
			ret.Type = PrimitiveTypeBool
			ret.Bool = l.Int != r.Int
		}
	} else if l.Type == PrimitiveTypeString && r.Type == PrimitiveTypeString {
		ret.Type = PrimitiveTypeString
		switch n.Op {
		case BinOpAdd:
			ret.String = l.String + r.String
		case BinOpEquality:
			ret.Type = PrimitiveTypeBool
			ret.Bool = l.String == r.String
		case BinOpNotEquality:
			ret.Type = PrimitiveTypeBool
			ret.Bool = l.String != r.String
		default:
			ret.Type = PrimitiveTypeUndefined
			context.Errors = append(context.Errors, ExecutionError{
				Class:        TypeErr,
				CreatingNode: n,
				Text:         "Invalid operation for string operands: " + n.Op.String(),
			})
		}
	} else if l.Type == PrimitiveTypeBool && r.Type == PrimitiveTypeBool {
		ret.Type = PrimitiveTypeBool
		switch n.Op {
		case BinOpEquality:
			ret.Bool = l.Bool && r.Bool
		case BinOpNotEquality:
			ret.Bool = l.Bool != r.Bool
		case BinOpLAnd:
			ret.Bool = l.Bool && r.Bool
		case BinOpLOr:
			ret.Bool = l.Bool || r.Bool
		default:
			ret.Type = PrimitiveTypeUndefined
			context.Errors = append(context.Errors, ExecutionError{
				Class:        TypeErr,
				CreatingNode: n,
				Text:         "Invalid operation for boolean operands: " + n.Op.String(),
			})
		}
	} else {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TypeErr,
			CreatingNode: n,
			Text:         "Invalid types for operands: " + l.Type.String() + " and " + r.Type.String(),
		})

	}

	return &ret
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
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
		Type: PrimitiveTypeUndefined,
		VariableReferenceFailed: true,
	}
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
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
		Type: PrimitiveTypeUndefined,
	}
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
func (n *IfStmt) Exec(context *ExecContext) *Variant {
	if n.Init != nil {
		n.Init.Exec(context)
	}

	conditionResult := n.Conditional.Exec(context)
	if conditionResult.Type == PrimitiveTypeBool && conditionResult.Bool {
		return n.Code.Exec(context)
	} else if n.Else != nil {
		return n.Else.Exec(context)
	}

	return &Variant{
		Type: PrimitiveTypeUndefined,
	}
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
func (n *ForStmt) Exec(context *ExecContext) *Variant {
	if n.Init != nil {
		n.Init.Exec(context)
	}

	conditionResult := n.Conditional.Exec(context)
	for true {
		if conditionResult.Type != PrimitiveTypeBool {
			context.Errors = append(context.Errors, ExecutionError{
				Class:        TypeErr,
				CreatingNode: n,
				Text:         "Non-bool used as loop conditional: " + conditionResult.Type.String(),
			})
			return &Variant{
				Type: PrimitiveTypeUndefined,
			}
		}
		if conditionResult.Bool {
			n.Code.Exec(context)
			if n.PostIteration != nil {
				n.PostIteration.Exec(context)
			}
		} else {
			break
		}
		conditionResult = n.Conditional.Exec(context)
	}

	return &Variant{
		Type: PrimitiveTypeUndefined,
	}
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
func (n *UnaryOp) Exec(context *ExecContext) *Variant {
	upper := n.Expr.Exec(context)
	if upper.Type == PrimitiveTypeBool {
		switch n.Op {
		case UnOpNot:
			return &Variant{
				Type: PrimitiveTypeBool,
				Bool: !upper.Bool,
			}
		default:
			context.Errors = append(context.Errors, ExecutionError{
				Class:        TypeErr,
				CreatingNode: n,
				Text:         "Cannot perform boolean unary operation on " + upper.Type.String(),
			})
		}
	} else {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TypeErr,
			CreatingNode: n,
			Text:         "Cannot perform unary operations on type " + upper.Type.String(),
		})
	}
	return &Variant{
		Type: PrimitiveTypeUndefined,
	}
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
func (n *Subscript) Exec(context *ExecContext) *Variant {
	baseVar := n.Expr.Exec(context)
	subscript := n.Subscript.Exec(context)

	if baseVar.VariableReferenceFailed {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        NotFoundErr,
			CreatingNode: n,
			Text:         "Could not resolve a value/variable for base data of type " + baseVar.Type.String(),
		})
		return &Variant{
			Type: PrimitiveTypeUndefined,
		}
	}

	if baseVar.Type != ComplexTypeArray {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TypeErr,
			CreatingNode: n,
			Text:         "Cannot perform subscript operation on type " + baseVar.Type.String(),
		})
		return &Variant{
			Type: PrimitiveTypeUndefined,
		}
	}
	if subscript.Type != PrimitiveTypeInt {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TypeErr,
			CreatingNode: n,
			Text:         "Cannot perform subscript operation on type " + baseVar.Type.String(),
		})
		return &Variant{
			Type: PrimitiveTypeUndefined,
		}
	}
	if int(subscript.Int) >= len(baseVar.VectorData) {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        BoundsErr,
			CreatingNode: n,
			Text:         "Subscript out of bounds",
		})
		return &Variant{
			Type: PrimitiveTypeUndefined,
		}
	}

	return baseVar.VectorData[subscript.Int]
}

// Exec carries out node-specific logic, which may include evaluation of subnodes and primitive operations depending on the nodes type.
func (n *NamedSelector) Exec(context *ExecContext) *Variant {
	baseVar := n.Expr.Exec(context)

	if baseVar.Type.Kind() != ComplexTypeStruct {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TypeErr,
			CreatingNode: n,
			Text:         "Cannot perform named selection operation on type " + baseVar.Type.String(),
		})
		return &Variant{
			Type: PrimitiveTypeUndefined,
		}
	}

	if v, ok := baseVar.NamedData[n.Name]; ok {
		return v
	}

	context.Errors = append(context.Errors, ExecutionError{
		Class:        NotFoundErr,
		CreatingNode: n,
		Text:         "Cannot find named element " + n.Name,
	})
	return &Variant{
		Type: PrimitiveTypeUndefined,
	}
}

// Exec represents the invocation of the FunctionCall - with the function pointer and arguments resolved from the contained nodes.
func (n *FunctionCall) Exec(context *ExecContext) *Variant {
	fmt.Println(context.GlobalNamespace)

	functionPointer := n.Function.Exec(context)
	if functionPointer.Type.Kind() != ComplexTypeFunction {
		context.Errors = append(context.Errors, ExecutionError{
			Class:        TypeErr,
			CreatingNode: n,
			Text:         "Cannot call non-function type: " + functionPointer.Type.String(),
		})
		return &Variant{
			Type: PrimitiveTypeUndefined,
		}
	}

	fn := map[string]*Variant{}
	execContext := &ExecContext{
		IsFuncContext:     true,
		FunctionNamespace: fn,
		GlobalNamespace:   context.GlobalNamespace,
	}

	for i, paramNode := range functionPointer.Type.(FunctionType).Parameters {
		pn := n.Args[i].Exec(context)
		nt := paramNode.(NamedType)
		fn[nt.Ident] = pn
	}

	return functionPointer.Type.(FunctionType).Code.Exec(execContext)
}
