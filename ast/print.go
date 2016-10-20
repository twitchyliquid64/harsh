package ast

import (
	"fmt"
	"strconv"
)

func (node *StatementList) Print(level int) {
	printLeveled("nodes{", level)
	for _, n := range node.Stmts {
		if n == nil {
			printLeveled("NIL", level+1)
		} else {
			n.Print(level + 1)
		}
	}
	printLeveled("}", level)
}

func (node *ReturnStmt) Print(level int) {
	printLeveled("return{", level)
	if node.Expr == nil {
		printLeveled("NIL", level+1)
	} else {
		node.Expr.Print(level + 1)
	}
	printLeveled("}", level)
}

func (node *IntegerLiteral) Print(level int) {
	printLeveled(strconv.FormatInt(node.Val, 10)+" int64", level)
}

func printLeveled(str string, level int) {
	for i := 0; i < level; i++ {
		fmt.Print(" ")
	}
	fmt.Println(str)
}

//TypeDecl
func (t *PrimitiveType) String() string {
	if t.Kind == PRIMITIVE_TYPE_INT {
		return "int{" + t.Name + "}"
	}
	if t.Kind == PRIMITIVE_TYPE_STRING {
		return "string{" + t.Name + "}"
	}
	return "?{" + t.Name + "}"
}
