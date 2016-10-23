package main

import (
	"fmt"
	"os"

	"github.com/twitchyliquid64/harsh/compiler"
)

func main() {
	fmt.Println("harsh: shiz.go")
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
		fmt.Print(name, v.Type.Kind.String())
		if i+1 < len(context.Globals.Names()) {
			fmt.Print(", ")
		}
	}
	fmt.Println("}")

	for _, decl := range context.Declarations {
		fmt.Println("FUNCTION: ", decl.Identifier)
		fmt.Print("Params: {")
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
			fmt.Println("-", ret.String(), "(return)")
		}
	}
}
