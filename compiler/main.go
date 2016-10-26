package compiler

import (
	//goast "go/ast"
	"errors"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"

	"github.com/twitchyliquid64/harsh/ast"
)

var ErrFuncNotFound = errors.New("No function found")

func ParseFile(fname string) (ast.Node, *Context, error) {
	fset := token.NewFileSet()

	goAstFile, err := parser.ParseFile(fset, fname, nil, parser.AllErrors)
	if err != nil {
		return nil, nil, err
	}
	rootNode, context := translateGoAST(fset, goAstFile)

	//goast.Print(fset, goAstFile)
	return rootNode, context, nil
}

func ParseLiteral(fname, inCode string) (context *Context, err error) {
	fset := token.NewFileSet()
	goAst, err := parser.ParseFile(fset, fname, inCode, 0)
	if err != nil {
		return nil, err
	}
	ns := ast.Namespace(map[string]*ast.Variant{})
	context = &Context{
		ConType: CONTEXT_ADHOC,
		Globals: ns,
	}

	if err != nil {
		return nil, err
	}
	a := translateGoNode(fset, context, reflect.ValueOf(goAst))
	if a != nil {
		return nil, err
	}
	return context, nil
}

type ExecutionError struct {
	Errors []ast.ExecutionError
}

func (e ExecutionError) Error() string {
	return strconv.Itoa(len(e.Errors)) + " execution errors"
}

func (c *Context) CallFunc(name string, args map[string]interface{}) (*ast.Variant, error) {
	for _, decl := range c.Declarations {
		if decl.Identifier == name {
			execContext := &ast.ExecContext{
				IsFuncContext:     true,
				FunctionNamespace: map[string]*ast.Variant{},
				GlobalNamespace:   c.Globals,
			}
			if args != nil {
				for name, arg := range args {
					execContext.FunctionNamespace[name] = ast.MakeVariant(arg)
				}
			}

			retValue := decl.Code.Exec(execContext)
			if len(execContext.Errors) == 0 {
				return retValue, nil
			} else {
				return retValue, ExecutionError{Errors: execContext.Errors}
			}
		}
	}
	return &ast.Variant{
		Type: ast.PRIMITIVE_TYPE_UNDEFINED,
	}, ErrFuncNotFound
}
