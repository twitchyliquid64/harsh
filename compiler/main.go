package compiler

import (
	//goast "go/ast"
	"go/parser"
	"go/token"

	"github.com/twitchyliquid64/harsh/ast"
)

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
