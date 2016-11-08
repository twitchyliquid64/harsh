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
)

// TypeError represents an error in the AST found during Typecheck().
type TypeError struct {
	Msg  string
	Kind TypeErrorKind
}

func (t TypeError) Error() string {
	return t.Msg
}

// TypeEqual returns true if the given types are equivalent and can be operated without promotion.
func TypeEqual(l ast.TypeKind, r ast.TypeKind) bool {
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

	default:
		context.Errors = append(context.Errors, TypeError{
			Kind: TypeerrorInternalErr,
			Msg:  "Cannot typecheck node of type " + reflect.TypeOf(node).String(),
		})
	}
	return ast.UnknownType
}
