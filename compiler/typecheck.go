package compiler

import (
	"reflect"

	"github.com/twitchyliquid64/harsh/ast"
)

type TypecheckContext struct {
	Errors []TypeError
}

type TypeErrorKind int

const (
	TYPEERROR_INTERNAL_ERR TypeErrorKind = iota
	TYPEERROR_INCOMPATIBLE_TYPES_ERR
)

type TypeError struct {
	Msg  string
	Kind TypeErrorKind
}

func (t TypeError) Error() string {
	return t.Msg
}

func TypeEqual(l ast.TypeDecl, r ast.TypeDecl) bool {
	return l == r
}

func Typecheck(context *TypecheckContext, node ast.Node) ast.TypeDecl {
	switch n := (node).(type) {
	case *ast.StatementList:
		for _, subNode := range n.Stmts {
			Typecheck(context, subNode)
		}

	case *ast.NilLiteral:
	case *ast.StringLiteral:
		return ast.PRIMITIVE_TYPE_STRING
	case *ast.IntegerLiteral:
		return ast.PRIMITIVE_TYPE_INT
	case *ast.BoolLiteral:
		return ast.PRIMITIVE_TYPE_BOOL
	case *ast.BinaryOp:
		l := Typecheck(context, n.LHS)
		r := Typecheck(context, n.RHS)
		if !TypeEqual(l, r) {
			context.Errors = append(context.Errors, TypeError{
				Kind: TYPEERROR_INCOMPATIBLE_TYPES_ERR,
				Msg:  "Cannot perform binary operation " + n.Op.String() + " on operands with type " + l.String() + " and " + r.String(),
			})
			return ast.PRIMITIVE_TYPE_UNDEFINED
		}
		return l

	default:
		context.Errors = append(context.Errors, TypeError{
			Kind: TYPEERROR_INTERNAL_ERR,
			Msg:  "Cannot typecheck node of type " + reflect.TypeOf(node).String(),
		})
	}
	return ast.PRIMITIVE_TYPE_UNDEFINED
}
