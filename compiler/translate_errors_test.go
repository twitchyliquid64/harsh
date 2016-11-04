package compiler

import (
	goast "go/ast"
	"go/token"
	"reflect"
	"testing"

	"github.com/twitchyliquid64/harsh/ast"
)

func TestBasicLitKindUnknownProducesError(t *testing.T) {
	ns := ast.Namespace(map[string]*ast.Variant{})
	context := &Context{
		ConType: ContextAdhoc,
		Globals: ns,
		//Debug:   true,
	}
	node := translateGoNode(nil, context, reflect.ValueOf(goast.BasicLit{
		Kind: token.AND_NOT,
	}))
	if len(context.Errors) != 1 {
		t.Error("Error expected")
	}
	if context.Errors[0].Class != NotSupported{
		t.Error("Incorrect error class")
	}
	if node != nil {
		t.Error("Nil node expected")
	}
}

func TestMultipleReturnProducesError(t *testing.T) {
	ns := ast.Namespace(map[string]*ast.Variant{})
	context := &Context{
		ConType: ContextAdhoc,
		Globals: ns,
		//Debug:   true,
	}
	node := translateGoNode(nil, context, reflect.ValueOf(goast.ReturnStmt{
		Results: []goast.Expr{
			nil, nil,
		},
	}))
	if len(context.Errors) != 1 {
		t.Error("Error expected")
	}
	if context.Errors[0].Class != NotYetSupported {
		t.Error("Incorrect error class")
	}
	if node != nil {
		t.Error("Nil node expected")
	}
}

func TestImportsProducesError(t *testing.T) {
	ns := ast.Namespace(map[string]*ast.Variant{})
	context := &Context{
		ConType: ContextAdhoc,
		Globals: ns,
		//Debug:   true,
	}
	translateGoGenDecl(nil, context, &goast.GenDecl{
		Specs: []goast.Spec{
			&goast.ImportSpec{},
		},
	})
	if len(context.Errors) != 1 {
		t.Error("Error expected")
	}
	if context.Errors[0].Class != NotYetSupported {
		t.Error("Incorrect error class")
	}
}

func TestUnsupportedNodeProducesError(t *testing.T) {
	ns := ast.Namespace(map[string]*ast.Variant{})
	context := &Context{
		ConType: ContextAdhoc,
		Globals: ns,
		//Debug:   true,
	}
	node := translateGoNode(nil, context, reflect.ValueOf(goast.GoStmt{}))
	if len(context.Errors) != 1 {
		t.Error("Error expected")
	}
	if context.Errors[0].Class != NotSupported{
		t.Error("Incorrect error class")
	}
	if node != nil {
		t.Error("Nil node expected")
	}
}

func TestUnknownGlobalProducesError(t *testing.T) {
	ns := ast.Namespace(map[string]*ast.Variant{})
	context := &Context{
		ConType: ContextAdhoc,
		Globals: ns,
		//Debug:   true,
	}
	translateGoGenDecl(nil, context, &goast.GenDecl{
		Specs: []goast.Spec{
			&goast.ValueSpec{
				Type: &goast.Ident{
					Name: "unsupportedType",
				},
				Names: []*goast.Ident{
					&goast.Ident{
						Name: "comp",
					},
				},
			},
		},
	})
	if len(context.Errors) != 1 {
		t.Error("Error expected")
	}
	if context.Errors[0].Class != NotSupported{
		t.Error("Incorrect error class")
	}
}

func TestExecutionErrorStringCorrect(t *testing.T) {
	if (&ExecutionError{}).Error() != "0 execution errors" {
		t.Error("ExecutionError string incorrect")
	}
}
