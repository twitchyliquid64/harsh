package compiler

import (
	"go/token"

	"github.com/twitchyliquid64/harsh/ast"
)

// Context represents a parsed Go module.
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
	Results    []ast.TypeKind
	Parameters []ast.TypeKind
}

type conType int

const (
	// ContextAdhoc represents a Parse/translation pass done on arbitrarily structured code - IE: not broken up into modules.
	ContextAdhoc = 0
	// ContextFile represents a Parse/translation done on a file or set of files, where nested contexts are used to represent each layer.
	ContextFile = 1
)

type translateErrClass int

const (
	// NotSupported indicates that the semantic/syntatic feature represented in the go AST is not supported by harsh.
	NotSupported translateErrClass = iota
	// NotYetSupported indicates a feature cannot yet be translated, but further development may support the feature (submit a PR!)
	NotYetSupported
	// TypeErrorFound indicates translate type typechecking failed.
	TypeErrorFound
	// NotStatic indicates that an expression was not able to be resolved at compiletime, when such a condition must exist.
	NotStatic
)

// TranslateError represents any issue translating the go.ast format into harsh's native AST format.
type TranslateError struct {
	Pos   token.Position
	Class translateErrClass
	Text  string
}
