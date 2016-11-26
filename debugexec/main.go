package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/twitchyliquid64/harsh/ast"
	"github.com/twitchyliquid64/harsh/compiler"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("USAGE: ./debugexec <harsh file> <function name> [<parameter-name>=<parameter-value>...]")
		return
	}

	// Parse & translate the file
	_, context, err := compiler.ParseFile(os.Args[1])
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	if len(context.Errors) > 0 {
		fmt.Println("Parse Errors:")
		for i, err := range context.Errors {
			fmt.Printf("%02d: %s (%s)\r\n", i+1, err.Text, err.Pos.String())
		}
		return
	}

	// Type check all the functional declarations
	wereTypeErrors := false
	for _, f := range context.Declarations {
		if _, ok := f.Type.(ast.FunctionType); !ok {
			continue
		}
		c := &compiler.TypecheckContext{}
		c.ReturnType = f.Type.(ast.FunctionType).ReturnType

		compiler.Typecheck(c, f.Type.(ast.FunctionType).Code)
		if len(c.Errors) > 0 {
			fmt.Println("Type errors in " + f.Ident + ":")
			for i, e := range c.Errors {
				fmt.Printf("%02d: %s (%d)\r\n", i+1, e.Msg, e.Kind)
			}
			wereTypeErrors = true
		}
	}
	if wereTypeErrors { //abort execution if there were type errors
		return
	}

	//Find the function we will execute - so we can match input parameters
	var funcDecl ast.NamedType
	for _, f := range context.Declarations {
		if _, ok := f.Type.(ast.FunctionType); ok && f.Ident == os.Args[2] {
			funcDecl = f
			break
		}
	}
	if funcDecl.Ident == "" {
		fmt.Println("A function with the name '" + os.Args[2] + "' does not exist")
		return
	}

	args := map[string]interface{}{}
	for i := 3; i < len(os.Args); i++ { //read in any args
		spl := strings.Split(os.Args[i], "=")
		if len(spl) > 1 {
			switch findFuncType(funcDecl.Type.(ast.FunctionType), spl[0]) {
			case "int":
				intValue, e := strconv.Atoi(spl[1])
				if e != nil {
					fmt.Println("Failed converting parameter '" + spl[0] + "' to int: " + e.Error())
					fmt.Println("Defaulting to 0.")
				}
				args[spl[0]] = intValue
			case "bool":
				boolValue, e := strconv.ParseBool(spl[1])
				if e != nil {
					fmt.Println("Failed converting parameter '" + spl[0] + "' to bool: " + e.Error())
					fmt.Println("Defaulting to false.")
				}
				args[spl[0]] = boolValue
			case "string":
				args[spl[0]] = spl[1]
			default:
				fmt.Println("Failed aligning input parameter '" + spl[0] + "' with function parameter definition.")
			}
		}
	}

	ret, err := context.CallFunc(os.Args[2], args)
	if err != nil {
		fmt.Println("Error:", err)
		if errs, ok := err.(compiler.ExecutionError); ok {
			for i, err := range errs.Errors {
				fmt.Printf("%02d: %s (%d)\r\n", i+1, err.Error(), err.Class)
				err.CreatingNode.Print(4, &ast.PrintContext{Output: os.Stdout, Color: true})
			}
		}
	} else {
		fmt.Println("Return Value: ", ret)
	}
}

func findFuncType(function ast.FunctionType, paramName string) string {
	for _, p := range function.Parameters {
		if p.(ast.NamedType).Ident == paramName {
			return p.(ast.NamedType).Type.String()
		}
	}
	return "?"
}
