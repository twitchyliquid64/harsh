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
	return (3 + 1) * 3
}

func translateGoAST(fset *token.FileSet, inp *goast.File) (ast.Node, *Context) {
	ns := ast.Namespace(map[string]ast.Variant{})

	context := &Context{
		ConType: CONTEXT_ADHOC,
		Globals: ns,
		Debug:   true,
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

		case goast.Ident:
			return &ast.VariableReference{
				Name: v.Name,
			}

		case goast.AssignStmt:
			for _, l := range v.Lhs {
				if ident, ok := l.(*goast.Ident); ok {
					if _, ok := ident.Obj.Decl.(*goast.AssignStmt); ok { //new local variable
						return &ast.Assign{
							NewLocal:   true,
							Identifier: ident.Name,
							Value:      translateGoNode(fset, context, reflect.ValueOf(v.Rhs[0])),
						}
					} else if _, ok := ident.Obj.Decl.(*goast.ValueSpec); ok {
						return &ast.Assign{
							NewLocal:   false,
							Identifier: ident.Name,
							Value:      translateGoNode(fset, context, reflect.ValueOf(v.Rhs[0])),
						}
					}
					fmt.Println("Assignment object unknown: ", ident.Name, reflect.TypeOf(ident.Obj.Decl))
				}
			}

		case goast.BasicLit:
			if v.Kind == token.INT {
				v, _ := strconv.ParseInt(v.Value, 10, 64)
				return &ast.IntegerLiteral{
					Val: int64(v),
				}
			} else if v.Kind == token.STRING {
				s, _ := strconv.Unquote(v.Value)
				return &ast.StringLiteral{
					Str: s,
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
			translateGoGenDecl(fset, context, node)
		default:
			fmt.Println("Unknown ast.Decl: ", reflect.TypeOf(decl))
		}
	}
}

func translateGoGenDecl(fset *token.FileSet, context *Context, node *goast.GenDecl) declaration {
	for _, spec := range node.Specs {
		switch n := spec.(type) {
		case *goast.ImportSpec:
			if context.Debug {
				fmt.Println("IMPORT", n.Path)
			}
			fmt.Println("Imports not yet supported")
		case *goast.ValueSpec:
			if context.Debug {
				fmt.Println("GLOBAL: ", n.Type, n.Names, n.Values, reflect.TypeOf(n.Type))
			}
			//global initializer expressions currently ignored.
			for _, name := range n.Names {
				switch t := n.Type.(type) {
				case *goast.Ident:
					switch t.Name {
					case "int":
						context.Globals.Save(name.Name, 0)
					case "string":
						context.Globals.Save(name.Name, "")
					default:
						context.Globals.Save(name.Name, ast.Variant{Type: ast.PrimitiveType{Kind: ast.PRIMITIVE_TYPE_UNDEFINED}})
						fmt.Println("Unknown goast.Ident.Type: ", t.Name)
					}
				}
				//TODO: Save default value based on type
			}
		default:
			fmt.Println("Unknown GenDecl subspec: ", reflect.TypeOf(node.Specs[0]))
		}
	}

	return declaration{}
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
