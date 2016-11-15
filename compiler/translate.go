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
		Debug:   true,
	}
	return translateGoNode(fset, context, reflect.ValueOf(inp)), context
}

func translateGoNode(fset *token.FileSet, context *Context, t reflect.Value) ast.Node {
	if context.Debug {
		fmt.Println("translateGoNode(): ", t.Kind(), t.Type().String())
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
			if v.Obj != nil {
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
				case *goast.FuncDecl:
					nt := translateGoFuncDecl(fset, context, n)
					t = nt.Type //unwrap Named node - all we want is the function type
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
			}
			return &ast.VariableReference{
				Name: v.Name,
				Type: ast.PrimitiveTypeUndefined,
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
				} else if _, ok := l.(*goast.SelectorExpr); ok {
					return &ast.Assign{
						NewLocal: false,
						Variable: translateGoNode(fset, context, reflect.ValueOf(l)),
						Value:    translateGoNode(fset, context, reflect.ValueOf(v.Rhs[0])),
					}
				}
				context.Errors = append(context.Errors, TranslateError{
					Class: NotSupported,
					Pos:   fset.Position(v.Pos()),
					Text:  "Assignment LHS unknown: " + reflect.TypeOf(l).Name(),
				})
			}
		case goast.UnaryExpr:
			if v.Op == token.NOT {
				return &ast.UnaryOp{
					Op:   ast.UnOpNot,
					Expr: translateGoNode(fset, context, reflect.ValueOf(v.X)),
				}
			}

		case goast.SelectorExpr:
			return &ast.NamedSelector{
				Name: v.Sel.Name,
				Expr: translateGoNode(fset, context, reflect.ValueOf(v.X)),
			}

		case goast.DeclStmt:
			switch d := v.Decl.(type) {
			case *goast.GenDecl:
				ln := ast.StatementList{}
				for _, spec := range d.Specs {
					if s, ok := spec.(*goast.ValueSpec); ok {
						for i, ident := range s.Names {
							assignNode := defaultValue(convertTypeToTypeKind(fset, s.Type, context), context)
							if i < len(s.Values) {
								assignNode = translateGoNode(fset, context, reflect.ValueOf(s.Values[i]))
							}
							ln.Stmts = append(ln.Stmts, &ast.Assign{
								NewLocal: true,
								Variable: translateGoNode(fset, context, reflect.ValueOf(ident)),
								Value:    assignNode,
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

		case goast.ExprStmt:
			if function, ok := v.X.(*goast.CallExpr); ok {
				return translateGoNode(fset, context, reflect.ValueOf(function))
			}

		case goast.CallExpr:
			var args []ast.Node
			for _, astArg := range v.Args {
				args = append(args, translateGoNode(fset, context, reflect.ValueOf(astArg)))
			}
			return &ast.FunctionCall{
				Function: translateGoNode(fset, context, reflect.ValueOf(v.Fun)),
				Args:     args,
			}

		case goast.CompositeLit: //composite literal: <type>{<values>...}
			subTypeOfComposite := convertTypeToTypeKind(fset, v.Type, context)
			orderedLiterals := []ast.Node{}
			namedLiterals := map[string]ast.Node{}

			// collect values
			for _, n := range v.Elts {
				switch valueNode := n.(type) {
				case *goast.KeyValueExpr:
					if _, ok := valueNode.Key.(*goast.Ident); !ok {
						context.Errors = append(context.Errors, TranslateError{
							Class: NotSupported,
							Pos:   fset.Position(v.Pos()),
							Text:  "Literal in composite with non-deterministic key is not supported: " + reflect.TypeOf(valueNode.Key).String(),
						})
					}
					namedLiterals[valueNode.Key.(*goast.Ident).Name] = translateGoNode(fset, context, reflect.ValueOf(valueNode.Value))
				default:
					orderedLiterals = append(orderedLiterals, translateGoNode(fset, context, reflect.ValueOf(n)))
				}
			}

			if _, ok := v.Type.(*goast.StructType); ok {
				if len(orderedLiterals) > 0 {
					context.Errors = append(context.Errors, TranslateError{
						Class: NotSupported,
						Pos:   fset.Position(v.Pos()),
						Text:  "Cannot have unnamed literals in struct composite literal",
					})
				}
				return &ast.StructLiteral{
					Type:   subTypeOfComposite.(ast.StructType),
					Values: namedLiterals,
				}
			}
			if len(namedLiterals) > 0 {
				context.Errors = append(context.Errors, TranslateError{
					Class: NotSupported,
					Pos:   fset.Position(v.Pos()),
					Text:  "Cannot have key-value pairs for non-struct composite literal of subtype: " + reflect.TypeOf(v.Type).String(),
				})
			}
			return &ast.ArrayLiteral{
				Type: ast.ArrayType{
					SubType: subTypeOfComposite,
					Len: &ast.IntegerLiteral{
						Val: int64(len(orderedLiterals)),
					},
				},
				Literal: orderedLiterals,
			}

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

func defaultValue(k ast.TypeKind, context *Context) ast.Node {
	if k == ast.PrimitiveTypeInt {
		return &ast.IntegerLiteral{}
	}
	if k == ast.PrimitiveTypeString {
		return &ast.StringLiteral{}
	}
	if a, ok := k.(ast.ArrayType); ok {
		return &ast.ArrayLiteral{
			Type:    a,
			Literal: nil,
		}
	}
	if st, ok := k.(ast.StructType); ok {
		return &ast.StructLiteral{
			Type:   st,
			Values: nil,
		}
	}
	context.Errors = append(context.Errors, TranslateError{
		Class: InternalErr,
		Text:  "Could not generate a default value for type: " + reflect.TypeOf(k).String(),
	})
	return &ast.NilLiteral{}
}

func convertTypeToTypeKind(fset *token.FileSet, t goast.Expr, context *Context) ast.TypeKind {
	if context.Debug {
		fmt.Println("convertTypeToTypeKind(): ", reflect.TypeOf(t))
	}
	//TODO: Refactor this mess to use a type switch
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
	} else if node, ok := t.(*goast.StructType); ok {
		structRet := ast.StructType{}
		if context.Debug {
			fmt.Println("Got struct", node.Fields.List)
		}
		for _, field := range node.Fields.List {
			ft := translateType(fset, field, context)
			if len(ft) != 1 {
				context.Errors = append(context.Errors, TranslateError{
					Class: InternalErr,
					Pos:   fset.Position(t.Pos()),
					Text:  "Struct field resolves to more than one TypeKind",
				})
				return ast.PrimitiveTypeUndefined
			}
			structRet.Fields = append(structRet.Fields, ft[0].(ast.NamedType))
		}
		return structRet
	}

	context.Errors = append(context.Errors, TranslateError{
		Class: NotSupported,
		Pos:   fset.Position(t.Pos()),
		Text:  "Cannot convert go/ast node to TypeKind: " + reflect.TypeOf(t).String(),
	})
	return ast.PrimitiveTypeUndefined
}

func translateType(fset *token.FileSet, typ *goast.Field, context *Context) []ast.TypeKind {
	if context.Debug {
		fmt.Println("translateType(): ", reflect.TypeOf(typ.Type))
	}
	kind := convertTypeToTypeKind(fset, typ.Type, context)
	var output []ast.TypeKind

	for _, name := range typ.Names {
		output = append(output, ast.NamedType{Type: kind, Ident: name.Name})
	}
	if len(typ.Names) == 0 {
		output = append(output, kind)
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
			newDecl := translateGoGenDecl(fset, context, node)
			context.Declarations = append(context.Declarations, newDecl)
		default:
			fmt.Println("Unknown ast.Decl: ", reflect.TypeOf(decl))
		}
	}
}

func translateGoGenDecl(fset *token.FileSet, context *Context, node *goast.GenDecl) ast.NamedType {
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
						return ast.NamedType{
							Ident: name.Name,
							Type:  ast.PrimitiveTypeInt,
						}
					case "bool":
						context.Globals.Save(name.Name, false)
						return ast.NamedType{
							Ident: name.Name,
							Type:  ast.PrimitiveTypeBool,
						}
					case "string":
						context.Globals.Save(name.Name, "")
						return ast.NamedType{
							Ident: name.Name,
							Type:  ast.PrimitiveTypeString,
						}
					default:
						context.Globals.Save(name.Name, ast.PrimitiveTypeUndefined)
						context.Errors = append(context.Errors, TranslateError{
							Class: NotSupported,
							Pos:   fset.Position(spec.Pos()),
							Text:  "Unknown global type: " + t.Name,
						})
					}
				case *goast.StructType:
					tk := convertTypeToTypeKind(fset, t, context)
					v, err := ast.DefaultVariantValue(tk)
					if err != nil {
						context.Errors = append(context.Errors, TranslateError{
							Class: NotStatic,
							Pos:   fset.Position(spec.Pos()),
							Text:  "Could not calculated default value for global: " + err.Error(),
						})
					} else {
						if context.Debug {
							fmt.Println(v)
						}
						context.Globals.Save(name.Name, v)
						return ast.NamedType{
							Ident: name.Name,
							Type:  tk,
						}
					}
				case *goast.ArrayType:
					tk := convertTypeToTypeKind(fset, t, context)
					v, err := ast.DefaultVariantValue(tk)
					if err != nil {
						context.Errors = append(context.Errors, TranslateError{
							Class: NotStatic,
							Pos:   fset.Position(spec.Pos()),
							Text:  "Could not calculated default value for global: " + err.Error(),
						})
					} else {
						context.Globals.Save(name.Name, v)
						return ast.NamedType{
							Ident: name.Name,
							Type:  tk,
						}
					}
				default:
					context.Errors = append(context.Errors, TranslateError{
						Class: NotYetSupported,
						Pos:   fset.Position(node.Pos()),
						Text:  "Cannot understand type of global " + name.Name + " - " + reflect.TypeOf(n.Type).String(),
					})
				}
			}
		default:
			fmt.Println("Unknown GenDecl subspec: ", reflect.TypeOf(node.Specs[0]))
		}
	}

	return ast.NamedType{
		Ident: "",
		Type:  ast.PrimitiveTypeUndefined,
	}
}

func translateGoFuncDecl(fset *token.FileSet, context *Context, node *goast.FuncDecl) ast.NamedType {
	var returnType ast.TypeKind = ast.PrimitiveTypeUndefined
	var parameters []ast.TypeKind

	if node.Type.Results != nil {
		if len(node.Type.Results.List) == 1 {
			if t := translateType(fset, node.Type.Results.List[0], context); t != nil {
				if len(t) == 1 {
					returnType = t[0]
				}
			}
		} else if len(node.Type.Results.List) > 1 {
			context.Errors = append(context.Errors, TranslateError{
				Class: NotYetSupported,
				Pos:   fset.Position(node.Pos()),
				Text:  "Functions with more than one result are not supported",
			})
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

	return ast.NamedType{
		Ident: node.Name.Name,
		Type: ast.FunctionType{
			Parameters: parameters,
			Code:       translateGoNode(fset, context, reflect.ValueOf(node.Body)),
			ReturnType: returnType,
		},
	}
}
