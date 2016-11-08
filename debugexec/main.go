package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/twitchyliquid64/harsh/compiler"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("USAGE: ./debugexec <harsh file> <function name> [<parameter-name>=<parameter-value>...]")
		return
	}

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

	wereTypeErrors := false
	for _, f := range context.Declarations {
		c := &compiler.TypecheckContext{}

		compiler.Typecheck(c, f.Code)
		if len(c.Errors) > 0 {
			fmt.Println("Type errors in " + f.Identifier + ":")
			for i, e := range c.Errors {
				fmt.Printf("%02d: %s (%d)\r\n", i+1, e.Msg, e.Kind)
			}
			wereTypeErrors = true
		}
	}
	if wereTypeErrors {
		return
	}

	args := map[string]interface{}{}
	for i := 3; i < len(os.Args); i++ {
		spl := strings.Split(os.Args[i], "=")
		if len(spl) > 1 {
			intValue, err := strconv.Atoi(spl[1])
			if err == nil {
				args[spl[0]] = intValue
				continue
			}
			boolValue, err := strconv.ParseBool(spl[1])
			if err == nil {
				args[spl[0]] = boolValue
				continue
			}
			args[spl[0]] = spl[1]
		}
	}
	ret, err := context.CallFunc(os.Args[2], args)
	if err != nil {
		fmt.Println("Error:", err)
		if errs, ok := err.(compiler.ExecutionError); ok {
			for i, err := range errs.Errors {
				fmt.Printf("%02d: %s (%d)\r\n", i+1, err.Error(), err.Class)
				err.CreatingNode.Print(4)
			}
		}
	} else {
		fmt.Println("Return Value: ", ret)
	}
}
