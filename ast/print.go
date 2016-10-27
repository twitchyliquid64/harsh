package ast

import (
	"fmt"
	"reflect"
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
	switch n := node.Type.(type) {
	case ArrayType:
		printLeveled("len{", level)
		n.Len.Print(level + 1)
		printLeveled("}"+node.Type.String(), level)
		return
	default:
		printLeveled("ERR unexpected node type: "+reflect.TypeOf(node.Type).Name(), level)
	}

	printLeveled(node.Type.String(), level)
}

func (node *BoolLiteral) Print(level int) {
	printLeveled(strconv.FormatBool(node.Val)+" bool", level)
}

func (node *BinaryOp) Print(level int) {
	printLeveled(node.Op.String()+" {", level)
	if node.LHS != nil {
		node.LHS.Print(level + 1)
	} else {
		printLeveled("NIL", level+1)
	}
	if node.RHS != nil {
		node.RHS.Print(level + 1)
	} else {
		printLeveled("NIL", level+1)
	}
	printLeveled("}", level)
}

func (node *UnaryOp) Print(level int) {
	printLeveled(node.Op.String()+" {", level)
	if node.Expr != nil {
		node.Expr.Print(level + 1)
	} else {
		printLeveled("NIL", level+1)
	}
	printLeveled("}", level)
}

func (node *Subscript) Print(level int) {
	printLeveled("subscript {", level)
	printLeveled("Index {", level+1)
	if node.Subscript != nil {
		node.Subscript.Print(level + 2)
	} else {
		printLeveled("NIL", level+2)
	}
	printLeveled("}", level+1)
	printLeveled("Expr {", level+1)
	if node.Expr != nil {
		node.Expr.Print(level + 2)
	} else {
		printLeveled("NIL", level+2)
	}
	printLeveled("}", level+1)
	printLeveled("}", level)
}

func (node *Assign) Print(level int) {
	if _, ok := node.Variable.(*VariableReference); ok {
		printLeveled(node.Variable.(*VariableReference).Name+" <= {", level)
	} else {
		printLeveled("assign {", level)
		node.Variable.Print(level + 1)
	}
	if node.Value == nil {
		printLeveled("NIL", level+1)
	} else {
		node.Value.Print(level + 1)
	}
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

func (op *UnOpType) String() string {
	switch *op {
	case UNOP_NOT:
		return "!"
	}
	return "UNOP?"
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
	case BINOP_EQUALITY:
		return "=="
	case BINOP_LAND:
		return "&&"
	case BINOP_LOR:
		return "||"
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
