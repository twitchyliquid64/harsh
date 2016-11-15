package compiler

import (
	"reflect"

	"github.com/twitchyliquid64/harsh/ast"
)

// TypecheckContext stores the persistent context of a recursive execution of Typecheck(), storing information such as type errors.
type TypecheckContext struct {
	Errors     []TypeError
	ReturnType ast.TypeKind
}

// TypeErrorKind represents an enum of error types which symbolise the kind of TypeError.
type TypeErrorKind int

const (
	// TypeerrorInternalErr represents a bug or an unreachable condition in the execution of TypeCheck.
	TypeerrorInternalErr TypeErrorKind = iota
	// TypeErrorIncompatibleTypesErr represents a combination of operands or operators which are invalid in respect to their types.
	TypeErrorIncompatibleTypesErr
	// TypeErrorNotFoundErr represents a situation where a named sub element is used which does not exist.
	TypeErrorNotFoundErr
)

// TypeError represents an error in the AST found during Typecheck().
type TypeError struct {
	Msg  string
	Kind TypeErrorKind
}

func (t TypeError) Error() string {
	return t.Msg
}

func checkStructsEqual(l ast.StructType, r ast.StructType) bool {
	// left join with left struct -- if equal continue
	for _, lfield := range l.Fields {
		for _, rfield := range r.Fields {
			if rfield.Ident == lfield.Ident {
				if !TypeEqual(lfield.Type, rfield.Type) {
					return false
				}
			}
		}
	}
	//left join with right struct
	for _, lfield := range r.Fields {
		for _, rfield := range l.Fields {
			if rfield.Ident == lfield.Ident {
				if !TypeEqual(lfield.Type, rfield.Type) {
					return false
				}
			}
		}
	}
	return true
}

// TypeEqual returns true if the given types are equivalent and can be operated without promotion.
func TypeEqual(l ast.TypeKind, r ast.TypeKind) bool {
	if l.Kind() == ast.ComplexTypeStruct && r.Kind() == ast.ComplexTypeStruct {
		return checkStructsEqual(l.(ast.StructType), r.(ast.StructType))
	}
	if l.Kind() == ast.ComplexTypeArray && r.Kind() == ast.ComplexTypeArray {
		return TypeEqual(l.(ast.ArrayType).SubType, r.(ast.ArrayType).SubType)
	}
	//TODO(twitchyliquid64): Check for function and check return/args individually

	return l == r
}

// Typecheck is a recursive method that returns the effective type of the return value of the node, if it were executed.
// Any type errors are added to context.Errors.
func Typecheck(context *TypecheckContext, node ast.Node) ast.TypeKind {
	switch n := (node).(type) {
	case *ast.StatementList:
		for _, subNode := range n.Stmts {
			Typecheck(context, subNode)
		}

	case *ast.FunctionCall:
		funcNodeType := Typecheck(context, n.Function)
		if funcNodeType.Kind() != ast.ComplexTypeFunction {
			context.Errors = append(context.Errors, TypeError{
				Kind: TypeErrorIncompatibleTypesErr,
				Msg:  "Cannot perform function invocation on type " + funcNodeType.String(),
			})
			return ast.UnknownType
		}
		if len(funcNodeType.(ast.FunctionType).Parameters) != len(n.Args) {
			context.Errors = append(context.Errors, TypeError{
				Kind: TypeErrorIncompatibleTypesErr,
				Msg:  "Cannot perform function invocation - incorrect number of parameters",
			})
			return ast.UnknownType
		}
		for i, param := range funcNodeType.(ast.FunctionType).Parameters {
			paramType := Typecheck(context, n.Args[i])
			if !TypeEqual(paramType, param) {
				context.Errors = append(context.Errors, TypeError{
					Kind: TypeErrorIncompatibleTypesErr,
					Msg:  "Parameter type mismatch: parameter has type " + param.String() + " but was called with " + paramType.String(),
				})
			}
		}
		return funcNodeType.(ast.FunctionType).ReturnType

	case *ast.VariableReference:
		if n.Type == nil {
			context.Errors = append(context.Errors, TypeError{
				Kind: TypeerrorInternalErr,
				Msg:  "VariableReference.Type should never be nil",
			})
			return ast.UnknownType
		}
		return n.Type
	case *ast.NilLiteral:
	case *ast.StringLiteral:
		return ast.PrimitiveTypeString
	case *ast.IntegerLiteral:
		return ast.PrimitiveTypeInt
	case *ast.BoolLiteral:
		return ast.PrimitiveTypeBool
	case *ast.BinaryOp:
		l := Typecheck(context, n.LHS)
		r := Typecheck(context, n.RHS)
		if l == ast.UnknownType || r == ast.UnknownType {
			return ast.UnknownType
		}
		if !TypeEqual(l, r) {
			context.Errors = append(context.Errors, TypeError{
				Kind: TypeErrorIncompatibleTypesErr,
				Msg:  "Cannot perform binary operation " + n.Op.String() + " on operands with type " + l.String() + " and " + r.String(),
			})
			return ast.UnknownType
		}
		return l

	case *ast.StructLiteral:
		for _, field := range n.Type.Fields {
			if v, ok := n.Values[field.Ident]; ok {
				vType := Typecheck(context, v)
				if !TypeEqual(vType, field.Type) {
					context.Errors = append(context.Errors, TypeError{
						Kind: TypeErrorIncompatibleTypesErr,
						Msg:  "Invalid struct literal - cannot have value of type " + vType.String() + " when the field is typed " + field.Type.String(),
					})
					return ast.UnknownType
				}
			}
		}
		return n.Type

	case *ast.ArrayLiteral:
		for _, element := range n.Literal {
			eType := Typecheck(context, element)
			if !TypeEqual(eType, n.Type.SubType) {
				context.Errors = append(context.Errors, TypeError{
					Kind: TypeErrorIncompatibleTypesErr,
					Msg:  "Invalid array literal - cannot have value of type " + eType.String() + " when the array contains elements of type " + n.Type.SubType.String(),
				})
				return ast.UnknownType
			}
		}
		return n.Type

	case *ast.Assign:
		l := Typecheck(context, n.Value)
		r := Typecheck(context, n.Variable)
		if l == ast.UnknownType || r == ast.UnknownType {
			return ast.UnknownType
		}
		if !TypeEqual(l, r) {
			context.Errors = append(context.Errors, TypeError{
				Kind: TypeErrorIncompatibleTypesErr,
				Msg:  "Cannot perform assignment to " + r.String() + " with type " + l.String(),
			})
			return ast.UnknownType
		}
		return l

	case *ast.ReturnStmt:
		if context.ReturnType != nil { //return type is known, test it
			v := Typecheck(context, n.Expr)
			if v == ast.UnknownType {
				return ast.UnknownType //cant compare to unknown
			}
			if !TypeEqual(v, context.ReturnType) {
				context.Errors = append(context.Errors, TypeError{
					Kind: TypeErrorIncompatibleTypesErr,
					Msg:  "Returned value does not match return type " + context.ReturnType.String() + ". Upstream value is " + Typecheck(context, n.Expr).String(),
				})
				return ast.UnknownType
			}
		}
		return Typecheck(context, n.Expr)

	case *ast.Subscript:
		sub := Typecheck(context, n.Subscript)
		if sub != ast.UnknownType && sub != ast.PrimitiveTypeInt {
			context.Errors = append(context.Errors, TypeError{
				Kind: TypeErrorIncompatibleTypesErr,
				Msg:  "Cannot subscript with non-integer index - got type: " + sub.String(),
			})
			return ast.UnknownType
		}
		RHS := Typecheck(context, n.Expr)
		if RHS.Kind() != ast.ComplexTypeArray { //TODO: support others
			context.Errors = append(context.Errors, TypeError{
				Kind: TypeErrorIncompatibleTypesErr,
				Msg:  "Cannot subscript non-array type " + RHS.String(),
			})
			return ast.UnknownType
		}
		return RHS.BaseType()

	case *ast.NamedSelector:
		up := Typecheck(context, n.Expr)
		if up.Kind() != ast.ComplexTypeStruct {
			context.Errors = append(context.Errors, TypeError{
				Kind: TypeErrorIncompatibleTypesErr,
				Msg:  "Cannot select non-struct type " + up.String(),
			})
			return ast.UnknownType
		}
		for _, field := range up.(ast.StructType).Fields {
			if field.Name() == n.Name {
				return field.BaseType()
			}
		}
		context.Errors = append(context.Errors, TypeError{
			Kind: TypeErrorNotFoundErr,
			Msg:  "Cannot find sub-element " + n.Name,
		})
		return ast.UnknownType

	case nil:
		context.Errors = append(context.Errors, TypeError{
			Kind: TypeerrorInternalErr,
			Msg:  "Cannot typecheck nil node - translate error?",
		})

	default:
		context.Errors = append(context.Errors, TypeError{
			Kind: TypeerrorInternalErr,
			Msg:  "Cannot typecheck node of type " + reflect.TypeOf(node).String(),
		})
	}
	return ast.UnknownType
}
