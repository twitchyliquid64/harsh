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
	printLeveled(strconv.Quote(node.Str)+" string", level)
}

func (node *ArrayLiteral) Print(level int) {
	if node.Type.(ArrayType).Len != nil {
		printLeveled("len{", level)
		node.Type.(ArrayType).Len.Print(level + 1)
		printLeveled("}"+node.Type.String(), level)
		return
	}
	printLeveled(node.Type.String(), level)
}

func (node *BoolLiteral) Print(level int) {
	printLeveled(strconv.FormatBool(node.Val)+" bool", level)
}

func (node *BinaryOp) Print(level int) {
	printLeveled(node.Op.String()+" {", level)
	node.LHS.Print(level + 1)
	node.RHS.Print(level + 1)
	printLeveled("}", level)
}

func (node *Assign) Print(level int) {
	printLeveled(node.Identifier+" <= {", level)
	node.Value.Print(level + 1)
	printLeveled("}", level)
}

func (node *VariableReference) Print(level int) {
	printLeveled("{"+node.Name+"}", level)
}

func (node *IfStmt) Print(level int) {
	printLeveled("if {", level)
	if node.Init != nil {
		printLeveled("init {", level+2)
		node.Init.Print(level + 3)
		printLeveled("}", level+2)
	}
	printLeveled("condition {", level+2)
	node.Conditional.Print(level + 3)
	printLeveled("}", level+2)
	printLeveled("code {", level+2)
	node.Code.Print(level + 3)
	printLeveled("}", level+2)
	if node.Else != nil {
		printLeveled("else {", level+2)
		node.Else.Print(level + 3)
		printLeveled("}", level+2)
	}
	printLeveled("}", level)
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

//Type system

func (t NamedType) String() string {
	return t.Type.String() + "{" + t.Ident + "}"
}

func (tk TypeKindDescription) String() string {
	switch tk {
	case PRIMITIVE_TYPE_INT:
		return "int"
	case PRIMITIVE_TYPE_STRING:
		return "string"
	case PRIMITIVE_TYPE_BOOL:
		return "bool"
	case COMPLEX_TYPE_ARRAY:
		return "[?]"
	case PRIMITIVE_TYPE_UNDEFINED:
		return "undefined"
	}
	return "?"
}
