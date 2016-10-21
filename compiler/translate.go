package compiler

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"reflect"
	"strconv"

	"github.com/twitchyliquid64/harsh/ast"
)

func testCrap(in int, bro, crap string) int {
	return 2*(1+2) - 4
}

func translateGoAST(fset *token.FileSet, inp *goast.File) (ast.Node, *Context) {
	context := &Context{
		ConType: CONTEXT_ADHOC,
	}
	return translateGoNode(fset, context, reflect.ValueOf(inp)), context
}

func translateGoNode(fset *token.FileSet, context *Context, t reflect.Value) ast.Node {
	if context.Debug {
		fmt.Println("translateGoNode(): ", t.Kind())
	}

	switch t.Kind() {
	case reflect.Ptr:
		return translateGoNode(fset, context, t.Elem())
	case reflect.Struct:
		s := t.Interface()
		switch v := s.(type) {
		case goast.File: //high level interface for a file.
			fileContext := context
			if context.ConType == CONTEXT_ADHOC {
				// do nothing - filecontext is the current context as it should be
			} else {
				fileContext = &Context{}
				context.ChildContexts = append(context.ChildContexts, fileContext)
			}
			fileContext.Name = v.Name.Name
			translateGoDecl(fset, fileContext, v.Decls)

		case goast.BlockStmt:
			sl := &ast.StatementList{}
			for _, stmt := range v.List {
				if n := translateGoNode(fset, context, reflect.ValueOf(stmt)); n != nil {
					sl.Stmts = append(sl.Stmts, n)
				}
			}
			return sl

		case goast.BinaryExpr:
			return &ast.BinaryOp{
				LHS: translateGoNode(fset, context, reflect.ValueOf(v.X)),
				RHS: translateGoNode(fset, context, reflect.ValueOf(v.Y)),
				Op:  translateGoBinop(v.Op),
			}

		case goast.ParenExpr:
			return translateGoNode(fset, context, reflect.ValueOf(v.X))

		case goast.AssignStmt:
			fmt.Println("Not implemented - ASSIGN: ", len(v.Lhs), len(v.Rhs))

		case goast.BasicLit:
			if v.Kind == token.INT {
				v, _ := strconv.ParseInt(v.Value, 10, 64)
				return &ast.IntegerLiteral{
					Val: int64(v),
				}
			}
			fmt.Println("Not implemented - BASICLIT: ", v.Value)

		case goast.ReturnStmt:
			if len(v.Results) > 0 { //only one return supported for now
				return &ast.ReturnStmt{
					Expr: translateGoNode(fset, context, reflect.ValueOf(v.Results[0])),
				}
			}
			fmt.Println("Not implemented - RETURN: ", len(v.Results))

		case goast.IfStmt:
			fmt.Println("Not implemented - IF: ") //whole lot of nodes (init, condition, else, main)

		case goast.SwitchStmt:
			fmt.Println("Not implemented - SWITCH: ", len(v.Body.List))

		default:
			fmt.Println("Got unknown struct type: ", t.Type())
		}

	}
	return nil
}

func translateGoBinop(tok token.Token) ast.BinOpType {
	switch tok {
	case token.ADD:
		return ast.BINOP_ADD
	case token.SUB:
		return ast.BINOP_SUB
	case token.MUL:
		return ast.BINOP_MUL
	default:
		fmt.Println("Unknown binop token.Token: ", reflect.TypeOf(tok))
		return ast.BINOP_UNK
	}
}

func translateType(typ *goast.Field) []ast.TypeDecl {
	var output []ast.TypeDecl
	switch node := typ.Type.(type) {
	case *goast.Ident:
		var kind ast.TypeKind
		if node.Name == "string" {
			kind = ast.PRIMITIVE_TYPE_STRING
		}
		if node.Name == "int" {
			kind = ast.PRIMITIVE_TYPE_INT
		}
		for _, name := range typ.Names {
			output = append(output, &ast.PrimitiveType{Kind: kind, Name: name.Name})
		}
		if len(typ.Names) == 0 {
			output = append(output, &ast.PrimitiveType{Kind: kind})
		}
	default:
		fmt.Println("translateType() unknown type: ", reflect.TypeOf(typ.Type))
		//goast.Print(nil, typ.Type)
	}
	return output
}

func translateGoDecl(fset *token.FileSet, context *Context, decls []goast.Decl) {
	for _, decl := range decls {
		switch node := decl.(type) {
		case *goast.FuncDecl:
			if context.Debug {
				fmt.Println("FUN DECL: ", node)
			}
			newDecl := translateGoFuncDecl(fset, context, node)
			context.Declarations = append(context.Declarations, newDecl)
		case *goast.GenDecl:
			if context.Debug {
				fmt.Println("GEN DECL: ", node)
			}
		default:
			fmt.Println("Unknown ast.Decl: ", reflect.TypeOf(decl))
		}
	}
}

func translateGoFuncDecl(fset *token.FileSet, context *Context, node *goast.FuncDecl) declaration {
	var returnVars []ast.TypeDecl
	var parameters []ast.TypeDecl

	if node.Type.Results != nil {
		for _, ret := range node.Type.Results.List {
			if t := translateType(ret); t != nil {
				for _, ret := range t {
					returnVars = append(returnVars, ret)
				}
			}
		}
	}
	if node.Type.Params != nil {
		for _, p := range node.Type.Params.List {
			if t := translateType(p); t != nil {
				for _, pm := range t {
					parameters = append(parameters, pm)
				}
			}
		}
	}

	return declaration{
		Identifier: node.Name.Name,
		Code:       translateGoNode(fset, context, reflect.ValueOf(node.Body)),
		Results:    returnVars,
		Parameters: parameters,
	}
}
