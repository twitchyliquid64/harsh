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

// ErrFuncNotFound is returned if Context.CallFunc() is called with a function that is not known in the context.
var ErrFuncNotFound = errors.New("No function found")

// ParseFile parses a Go source file and returns the root node, a AST Context, and any translation / parse errors.
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

// ParseLiteral takes a string of Go code, returning an AST context and any translation/parse errors.
func ParseLiteral(fname, inCode string) (context *Context, err error) {
	fset := token.NewFileSet()
	goAst, err := parser.ParseFile(fset, fname, inCode, 0)
	if err != nil {
		return nil, err
	}
	ns := ast.Namespace(map[string]*ast.Variant{})
	context = &Context{
		ConType: ContextAdhoc,
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

// ExecutionError represents a problem that arose when executing a function.
type ExecutionError struct {
	Errors []ast.ExecutionError
}

func (e ExecutionError) Error() string {
	return strconv.Itoa(len(e.Errors)) + " execution errors"
}

// CallFunc executes the named function in Context, with args, and returning a value. If the function does not exist
// or execution raises an error, an error is returned.
func (c *Context) CallFunc(name string, args map[string]interface{}) (*ast.Variant, error) {
	for _, decl := range c.Declarations {
		if decl.Ident == name {
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

			if _, ok := decl.Type.(ast.FunctionType); !ok {
				return &ast.Variant{Type: ast.PrimitiveTypeUndefined}, errors.New("Declaration is not a function")
			}

			retValue := decl.Type.(ast.FunctionType).Code.Exec(execContext)
			if len(execContext.Errors) == 0 {
				return retValue, nil
			}
			return retValue, ExecutionError{Errors: execContext.Errors}
		}
	}
	return &ast.Variant{
		Type: ast.PrimitiveTypeUndefined,
	}, ErrFuncNotFound
}
