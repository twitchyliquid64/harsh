package compiler

import (
	"reflect"

	"github.com/twitchyliquid64/harsh/ast"
)

// TypecheckContext stores the persistent context of a recursive execution of Typecheck(), storing information such as type errors.
type TypecheckContext struct {
	Errors []TypeError
}

// TypeErrorKind represents an enum of error types which symbolise the kind of TypeError.
type TypeErrorKind int

const (
	// TypeerrorInternalErr represents a bug or an unreachable condition in the execution of TypeCheck.
	TypeerrorInternalErr TypeErrorKind = iota
	// TypeerrorIncompatibleTypesErr represents a combination of operands or operators which are invalid in respect to their types.
	TypeerrorIncompatibleTypesErr
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
func TypeEqual(l ast.TypeDecl, r ast.TypeDecl) bool {
	return l == r
}

// Typecheck is a recursive method that returns the effective type of the return value of the node, if it were executed.
// Any type errors are added to context.Errors.
func Typecheck(context *TypecheckContext, node ast.Node) ast.TypeDecl {
	switch n := (node).(type) {
	case *ast.StatementList:
		for _, subNode := range n.Stmts {
			Typecheck(context, subNode)
		}

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
		if !TypeEqual(l, r) {
			context.Errors = append(context.Errors, TypeError{
				Kind: TypeerrorIncompatibleTypesErr,
				Msg:  "Cannot perform binary operation " + n.Op.String() + " on operands with type " + l.String() + " and " + r.String(),
			})
			return ast.PrimitiveTypeUndefined
		}
		return l

	default:
		context.Errors = append(context.Errors, TypeError{
			Kind: TypeerrorInternalErr,
			Msg:  "Cannot typecheck node of type " + reflect.TypeOf(node).String(),
		})
	}
	return ast.PrimitiveTypeUndefined
}
