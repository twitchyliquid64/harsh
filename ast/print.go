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

func (node *StringLiteral) Print(level int) {
	printLeveled(node.Str+" string", level)
}

func (node *BinaryOp) Print(level int) {
	printLeveled(node.Op.String()+" {", level)
	node.LHS.Print(level + 1)
	node.RHS.Print(level + 1)
	printLeveled("}", level)
}

func (node *VariableReference) Print(level int) {
	printLeveled("{"+node.Name+"}", level)
}

func (op *BinOpType) String() string {
	switch *op {
	case BINOP_ADD:
		return "+"
	case BINOP_SUB:
		return "-"
	case BINOP_MUL:
		return "*"
	case BINOP_DIV:
		return "/"
	case BINOP_UNK:
		return "UNK?"
	case BINOP_MOD:
		return "%"
	}
	return "BINOP?"
}

func printLeveled(str string, level int) {
	for i := 0; i < level; i++ {
		fmt.Print(" ")
	}
	fmt.Println(str)
}

//TypeDecl
func (t *PrimitiveType) String() string {
	return t.Kind.String() + "{" + t.Name + "}"
}

func (tk *TypeKind) String() string {
	switch *tk {
	case PRIMITIVE_TYPE_INT:
		return "int"
	case PRIMITIVE_TYPE_STRING:
		return "string"
	case PRIMITIVE_TYPE_UNDEFINED:
		return "undefined"
	}
	return "?"
}
