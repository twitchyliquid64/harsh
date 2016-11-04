package compiler

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"reflect"
	"strconv"

	"github.com/twitchyliquid64/harsh/ast"
)

func translateGoAST(fset *token.FileSet, inp *goast.File) (ast.Node, *Context) {
	ns := ast.Namespace(map[string]*ast.Variant{})

	context := &Context{
		ConType: ContextAdhoc,
		Globals: ns,
		//Debug:   true,
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
			if context.ConType == ContextAdhoc {
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
			if v.Name == "true" || v.Name == "false" {
				b, _ := strconv.ParseBool(v.Name) //TODO: Process error`
				return &ast.BoolLiteral{
					Val: b,
				}
			}
			var t ast.TypeKind = ast.UnknownType
			switch n := v.Obj.Decl.(type) {
			case *goast.ValueSpec:
				t = convertTypeToTypeKind(fset, n.Type, context)
			case *goast.AssignStmt:
				//try inferring type by typechecking the RHS of the assignment.
				assignRHSNode := translateGoNode(fset, context, reflect.ValueOf(n.Rhs[0]))
				tc := &TypecheckContext{}
				t = Typecheck(tc, assignRHSNode)
				if len(tc.Errors) > 0 {
					context.Errors = append(context.Errors, TranslateError{
						Class: TypeErrorFound,
						Pos:   fset.Position(v.Pos()),
						Text:  "Could not typecheck RHS of assignment to " + v.Name,
					})
				}
			case *goast.Field:
				t = convertTypeToTypeKind(fset, n.Type, context)
			default:
				context.Errors = append(context.Errors, TranslateError{
					Class: NotSupported,
					Pos:   fset.Position(v.Pos()),
					Text:  "ast.Ident.Obj.Decl type unknown: " + v.Name + " (" + reflect.TypeOf(n).Name() + ")",
				})
			}
			return &ast.VariableReference{
				Name: v.Name,
				Type: t,
			}

		case goast.AssignStmt:
			for _, l := range v.Lhs {
				if ident, ok := l.(*goast.Ident); ok {
					if _, ok := ident.Obj.Decl.(*goast.AssignStmt); ok { //new local variable
						return &ast.Assign{
							NewLocal: true,
							Variable: translateGoNode(fset, context, reflect.ValueOf(l)),
							Value:    translateGoNode(fset, context, reflect.ValueOf(v.Rhs[0])),
						}
					} else if _, ok := ident.Obj.Decl.(*goast.ValueSpec); ok {
						return &ast.Assign{
							NewLocal: false,
							Variable: translateGoNode(fset, context, reflect.ValueOf(l)),
							Value:    translateGoNode(fset, context, reflect.ValueOf(v.Rhs[0])),
						}
					}
					context.Errors = append(context.Errors, TranslateError{
						Class: NotSupported,
						Pos:   fset.Position(v.Pos()),
						Text:  "Assignment object unknown: " + ident.Name + " (" + reflect.TypeOf(ident.Obj.Decl).Name() + ")",
					})
				} else if _, ok := l.(*goast.IndexExpr); ok {
					return &ast.Assign{
						NewLocal: false,
						Variable: translateGoNode(fset, context, reflect.ValueOf(l)),
						Value:    translateGoNode(fset, context, reflect.ValueOf(v.Rhs[0])),
					}
				} else {
					context.Errors = append(context.Errors, TranslateError{
						Class: NotSupported,
						Pos:   fset.Position(v.Pos()),
						Text:  "Assignment LHS unknown: " + reflect.TypeOf(l).Name(),
					})
				}
			}
		case goast.UnaryExpr:
			if v.Op == token.NOT {
				return &ast.UnaryOp{
					Op:   ast.UnOpNot,
					Expr: translateGoNode(fset, context, reflect.ValueOf(v.X)),
				}
			}

		case goast.DeclStmt:
			switch d := v.Decl.(type) {
			case *goast.GenDecl:
				ln := ast.StatementList{}
				for _, spec := range d.Specs {
					if s, ok := spec.(*goast.ValueSpec); ok {
						for i := range s.Names {
							assignNode := defaultValue(convertTypeToTypeKind(fset, s.Type, context))
							if i < len(s.Values) {
								assignNode = translateGoNode(fset, context, reflect.ValueOf(s.Values[i]))
							}
							ln.Stmts = append(ln.Stmts, &ast.Assign{
								NewLocal: true,
								Variable: &ast.VariableReference{
									Name: s.Names[i].Name,
								},
								Value: assignNode,
							})
						}
					} else {
						context.Errors = append(context.Errors, TranslateError{
							Class: NotSupported,
							Pos:   fset.Position(v.Pos()),
							Text:  "Spec in Declaration unknown: (" + reflect.TypeOf(spec).Name() + ")",
						})
					}
				}
				return &ln
			}

		case goast.CompositeLit: //composite literal: <type>{<values>...}
			compLit := &ast.ArrayLiteral{
				Type: convertTypeToTypeKind(fset, v.Type, context),
			}
			for _, n := range v.Elts {
				compLit.Literal = append(compLit.Literal, translateGoNode(fset, context, reflect.ValueOf(n)))
			}
			return compLit

		case goast.IndexExpr:
			return &ast.Subscript{
				Expr:      translateGoNode(fset, context, reflect.ValueOf(v.X)),
				Subscript: translateGoNode(fset, context, reflect.ValueOf(v.Index)),
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
			} else {
				context.Errors = append(context.Errors, TranslateError{
					Class: NotSupported,
					Pos:   fset.Position(v.Pos()),
					Text:  "BasicLit Kind is not recognised: " + v.Kind.String(),
				})
			}

		case goast.ReturnStmt:
			if len(v.Results) == 1 { //only one return supported for now
				return &ast.ReturnStmt{
					Expr: translateGoNode(fset, context, reflect.ValueOf(v.Results[0])),
				}
			} else if len(v.Results) == 0 { //TODO: make a undefined node and return it
				return &ast.ReturnStmt{
					Expr: &ast.NilLiteral{},
				}
			} else {
				context.Errors = append(context.Errors, TranslateError{
					Class: NotYetSupported,
					Pos:   fset.Position(v.Pos()),
					Text:  "Returning multiple values is not supported.",
				})
			}

		case goast.IfStmt:
			return &ast.IfStmt{
				Init:        translateGoNode(fset, context, reflect.ValueOf(v.Init)),
				Code:        translateGoNode(fset, context, reflect.ValueOf(v.Body)),
				Else:        translateGoNode(fset, context, reflect.ValueOf(v.Else)),
				Conditional: translateGoNode(fset, context, reflect.ValueOf(v.Cond)),
			}

		case goast.SwitchStmt:
			fmt.Println("Not implemented - SWITCH: ", len(v.Body.List))

		default:
			context.Errors = append(context.Errors, TranslateError{
				Class: NotSupported,
				Text:  "Translation of go/ast node not supported: " + t.Type().Name(),
			})
		}

	}
	return nil
}

func translateGoBinop(tok token.Token) ast.BinOpType {
	switch tok {
	case token.ADD:
		return ast.BinOpAdd
	case token.SUB:
		return ast.BinOpSub
	case token.MUL:
		return ast.BinOpMul
	case token.REM:
		return ast.BinOpMod
	case token.QUO:
		return ast.BinOpDiv
	case token.LAND:
		return ast.BinOpLAnd
	case token.LOR:
		return ast.BinOpLOr
	case token.EQL:
		return ast.BinOpEquality
	default:
		fmt.Println("Unknown binop token.Token: ", tok.String())
		return ast.BinOpUnknown
	}
}

func defaultValue(k ast.TypeKind) ast.Node {
	if k == ast.PrimitiveTypeInt {
		return &ast.IntegerLiteral{}
	}
	if k == ast.PrimitiveTypeString {
		return &ast.StringLiteral{}
	}
	if _, ok := k.(ast.ArrayType); ok {
		return &ast.ArrayLiteral{
			Type:    k,
			Literal: nil,
		}
	}
	return &ast.NilLiteral{}
}

func convertTypeToTypeKind(fset *token.FileSet, t goast.Expr, context *Context) ast.TypeKind {
	if node, ok := t.(*goast.Ident); ok {
		if node.Name == "string" {
			return ast.PrimitiveTypeString
		}
		if node.Name == "int" {
			return ast.PrimitiveTypeInt
		}
		if node.Name == "bool" {
			return ast.PrimitiveTypeBool
		}
		context.Errors = append(context.Errors, TranslateError{
			Class: NotSupported,
			Text:  "Cannot convert go/ast.Ident to TypeKind: " + node.Name,
		})
	} else if node, ok := t.(*goast.ArrayType); ok {
		childTypeKind := convertTypeToTypeKind(fset, node.Elt, context)
		if childTypeKind == ast.PrimitiveTypeUndefined {
			context.Errors = append(context.Errors, TranslateError{
				Class: NotSupported,
				Pos:   fset.Position(node.Pos()),
				Text:  "Array uses unsupported type: " + reflect.TypeOf(node.Elt).String(),
			})
		} else { //build an array type based on it
			var lenNode = node.Len
			if lenNode == nil {
				context.Errors = append(context.Errors, TranslateError{
					Class: NotSupported,
					Pos:   fset.Position(node.Pos()),
					Text:  "Slices are not yet supported",
				})
			} else {
				return ast.ArrayType{
					SubType: childTypeKind,
					Len:     translateGoNode(fset, context, reflect.ValueOf(node.Len)),
				}
			}
		}
	} else {
		context.Errors = append(context.Errors, TranslateError{
			Class: NotSupported,
			Pos:   fset.Position(t.Pos()),
			Text:  "Cannot convert go/ast node to TypeKind: " + reflect.TypeOf(t).String(),
		})
	}
	return ast.PrimitiveTypeUndefined
}

func translateType(fset *token.FileSet, typ *goast.Field, context *Context) []ast.TypeDecl {
	var output []ast.TypeDecl
	switch typ.Type.(type) {
	case *goast.Ident:
		kind := convertTypeToTypeKind(fset, typ.Type, context)
		for _, name := range typ.Names {
			output = append(output, ast.NamedType{Type: kind, Ident: name.Name})
		}
		if len(typ.Names) == 0 {
			output = append(output, kind)
		}
	default:
		context.Errors = append(context.Errors, TranslateError{
			Class: NotSupported,
			Pos:   fset.Position(typ.Type.Pos()),
			Text:  "Translate type encounted unexpected type: " + reflect.TypeOf(typ.Type).String(),
		}) //goast.Print(nil, typ.Type)
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
			context.Errors = append(context.Errors, TranslateError{
				Class: NotYetSupported,
				Pos:   fset.Position(node.Pos()),
				Text:  "Import statements are not yet supported",
			})
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
					case "bool":
						context.Globals.Save(name.Name, false)
					case "string":
						context.Globals.Save(name.Name, "")
					default:
						context.Globals.Save(name.Name, ast.PrimitiveTypeUndefined)
						context.Errors = append(context.Errors, TranslateError{
							Class: NotSupported,
							Pos:   fset.Position(spec.Pos()),
							Text:  "Unknown global type: " + t.Name,
						})
					}
				}
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
			if t := translateType(fset, ret, context); t != nil {
				for _, ret := range t {
					returnVars = append(returnVars, ret)
				}
			}
		}
	}
	if node.Type.Params != nil {
		for _, p := range node.Type.Params.List {
			if t := translateType(fset, p, context); t != nil {
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
