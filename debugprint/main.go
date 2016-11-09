package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/twitchyliquid64/harsh/compiler"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE: ./debugprint <harsh file>")
		return
	}

	_, context, err := compiler.ParseFile(os.Args[1])
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Print("Globals: {")
	for i, name := range context.Globals.Names() {
		v := context.Globals[name]
		fmt.Print(name, " ", v.Type.String())
		if i+1 < len(context.Globals.Names()) {
			fmt.Print(", ")
		}
	}
	fmt.Println("}")

	for _, decl := range context.Declarations {
		fmt.Println("FUNCTION: ", decl.Identifier)
		fmt.Print("  Params: {")
		for i, param := range decl.Parameters {
			fmt.Print(param.String())
			if i+1 < len(decl.Parameters) {
				fmt.Print(", ")
			}
		}
		fmt.Println("}")
		if decl.Code != nil {
			decl.Code.Print(2)
		}
		for _, ret := range decl.Results {
			fmt.Println("  -", ret.String(), "(return)")
		}

		c := &compiler.TypecheckContext{}
		if len(decl.Results) == 1 {
			c.ReturnType = decl.Results[0]
		}
		compiler.Typecheck(c, decl.Code)
		if len(c.Errors) > 0 {
			fmt.Println("  Type errors:")
			for i, e := range c.Errors {
				fmt.Printf("   %02d: %s (%d)\r\n", i+1, e.Msg, e.Kind)
			}
		}
	}

	if len(context.Errors) > 0 {
		fmt.Println("Translate Errors:")
		for i, err := range context.Errors {
			fmt.Printf("%02d: %s (%s)\r\n", i+1, err.Text, err.Pos.String())
		}
	}

	if len(os.Args) > 2 && os.Args[2] == "--goast" {
		fset := token.NewFileSet() // positions are relative to fset
		f, err := parser.ParseFile(fset, os.Args[1], nil, 0)
		if err != nil {
			fmt.Println(err)
		}
		ast.Print(fset, f)
	}
}
