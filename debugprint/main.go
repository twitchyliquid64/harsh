package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	myast "github.com/twitchyliquid64/harsh/ast"
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

	for _, decl := range context.Declarations {
		if fType, ok := decl.Type.(myast.FunctionType); ok {
			fmt.Println("FUNCTION: ", decl.String())
			if fType.Code != nil {
				fType.Code.Print(2, &myast.PrintContext{Output: os.Stdout, Color: true})
			}
			fmt.Println("  -", fType.ReturnType.String(), "(return)")

			c := &compiler.TypecheckContext{}
			c.ReturnType = fType.ReturnType
			compiler.Typecheck(c, fType.Code)
			if len(c.Errors) > 0 {
				fmt.Println("  Type errors:")
				for i, e := range c.Errors {
					fmt.Printf("   %02d: %s (%d)\r\n", i+1, e.Msg, e.Kind)
				}
			}
		} else {
			fmt.Println("DECLARATION: ", decl.String())
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
