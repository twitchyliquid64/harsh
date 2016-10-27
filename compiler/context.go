package compiler

import (
	"go/token"

	"github.com/twitchyliquid64/harsh/ast"
)

type Context struct {
	Name          string
	ConType       conType
	Declarations  []declaration
	ChildContexts []*Context
	Debug         bool
	Globals       ast.Namespace
	Errors        []TranslateError
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

type translateErrClass int

const (
	NOT_SUPPORTED translateErrClass = iota
	NOT_YET_SUPPORTED
)

type TranslateError struct {
	Pos   token.Position
	Class translateErrClass
	Text  string
}
