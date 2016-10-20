package compiler

import "github.com/twitchyliquid64/harsh/ast"

type Context struct {
	Name          string
	ConType       conType
	Declarations  []declaration
	ChildContexts []*Context
	Debug         bool
}

type declaration struct {
	Identifier string
	Code       ast.Node
	Results    []ast.TypeDecl
	Parameters []ast.TypeDecl
}

type conType int

const (
	CONTEXT_ADHOC = 0
	CONTEXT_FILE  = 1
)
