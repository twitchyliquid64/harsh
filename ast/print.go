package ast

import (
	"fmt"
	"strconv"
)

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *StatementList) Print(level int, printContext *PrintContext) {
	openSection("", level, printContext)
	for _, n := range node.Stmts {
		if n == nil {
			outputNil(level+1, printContext)
		} else {
			n.Print(level+1, printContext)
		}
	}
	closeSection(level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *NilLiteral) Print(level int, printContext *PrintContext) {
	if printContext.Color {
		outputLeveled(red()+"nil"+resetColor(), level, printContext)
	} else {
		outputNil(level, printContext)
	}
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *ReturnStmt) Print(level int, printContext *PrintContext) {
	openSection("return", level, printContext)
	if node.Expr == nil {
		outputNil(level+1, printContext)
	} else {
		node.Expr.Print(level+1, printContext)
	}
	closeSection(level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *NamedSelector) Print(level int, printContext *PrintContext) {
	outputLeveled(ifColor(yellow(), printContext)+"."+ifColor(resetColor(), printContext)+
		outputBaseSource(node.Name, printContext)+
		ifColor(blue(), printContext)+" {", level, printContext)
	if node.Expr == nil {
		outputNil(level+1, printContext)
	} else {
		node.Expr.Print(level+1, printContext)
	}
	closeSection(level, printContext)
}

// Print writes a description of the function call to standard output.
func (node *FunctionCall) Print(level int, printContext *PrintContext) {
	openSection("INVOCATION", level, printContext)
	if node.Function == nil {
		outputNil(level+1, printContext)
	} else {
		node.Function.Print(level+1, printContext)
	}
	openSection("args", level+1, printContext)
	for _, arg := range node.Args {
		if arg == nil {
			outputNil(level+2, printContext)
		} else {
			arg.Print(level+2, printContext)
		}
	}
	closeSection(level+1, printContext)
	closeSection(level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *IntegerLiteral) Print(level int, printContext *PrintContext) {
	outputLeveled(outputBaseSource(strconv.FormatInt(node.Val, 10), printContext)+outputType(" int64", printContext), level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *StringLiteral) Print(level int, printContext *PrintContext) {
	outputLeveled(outputBaseSource(strconv.Quote(node.Str), printContext)+outputType(" string", printContext), level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *ArrayLiteral) Print(level int, printContext *PrintContext) {
	openSection("array", level, printContext)
	openSection("len", level+1, printContext)
	node.Type.Len.Print(level+2, printContext)
	closeSection(level+1, printContext)
	if len(node.Literal) > 0 {
		openSection("values", level+1, printContext)
		for _, v := range node.Literal {
			if v == nil {
				outputNil(level+2, printContext)
			} else {
				v.Print(level+2, printContext)
			}
		}
		closeSection(level+1, printContext)
	}
	outputLeveled(ifColor(blue(), printContext)+"} "+ifColor(resetColor(), printContext)+
		outputType("<"+node.Type.String()+">", printContext), level, printContext)
}

// Print writes a description of the struct to standard output, at the specified indentation level.
func (node *StructLiteral) Print(level int, printContext *PrintContext) {
	openSection("struct", level, printContext)
	for _, field := range node.Type.Fields {
		if node.Values != nil && node.Values[field.Ident] != nil {
			openSection("Field: '"+field.Ident+"'", level+1, printContext)
			node.Values[field.Ident].Print(level+2, printContext)
			outputLeveled("} <"+field.Type.String()+">", level+1, printContext)
		} else {
			outputLeveled(ifColor(yellow(), printContext)+"Field: '"+
				outputBaseSource(field.Ident, printContext)+
				ifColor(yellow(), printContext)+"' "+
				ifColor(red(), printContext)+"(nil) "+
				outputType("<"+field.Type.String()+">", printContext), level+1, printContext)
		}
	}
	closeSection(level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *BoolLiteral) Print(level int, printContext *PrintContext) {
	outputLeveled(outputBaseSource(strconv.FormatBool(node.Val), printContext)+outputType(" bool", printContext), level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *BinaryOp) Print(level int, printContext *PrintContext) {
	openSection(node.Op.String(), level, printContext)
	if node.LHS != nil {
		node.LHS.Print(level+1, printContext)
	} else {
		outputNil(level+1, printContext)
	}
	if node.RHS != nil {
		node.RHS.Print(level+1, printContext)
	} else {
		outputNil(level+1, printContext)
	}
	closeSection(level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *UnaryOp) Print(level int, printContext *PrintContext) {
	openSection(node.Op.String(), level, printContext)
	if node.Expr != nil {
		node.Expr.Print(level+1, printContext)
	} else {
		outputNil(level+1, printContext)
	}
	closeSection(level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *Subscript) Print(level int, printContext *PrintContext) {
	openSection("subscript", level, printContext)
	openSection("index", level+1, printContext)
	if node.Subscript != nil {
		node.Subscript.Print(level+2, printContext)
	} else {
		outputNil(level+2, printContext)
	}
	closeSection(level+1, printContext)
	openSection("expr", level+1, printContext)
	if node.Expr != nil {
		node.Expr.Print(level+2, printContext)
	} else {
		outputNil(level+2, printContext)
	}
	closeSection(level+1, printContext)
	closeSection(level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *Assign) Print(level int, printContext *PrintContext) {
	if _, ok := node.Variable.(*VariableReference); ok {
		typeStr := "<?>"
		if node.Variable.(*VariableReference).Type != nil {
			typeStr = "<" + node.Variable.(*VariableReference).Type.String() + ">"
		}
		outputLeveled(outputBaseSource(node.Variable.(*VariableReference).Name, printContext)+
			" "+outputType(typeStr, printContext)+
			ifColor(blue(), printContext)+" <= {"+ifColor(resetColor(), printContext), level, printContext)
	} else {
		openSection("assign", level, printContext)
		node.Variable.Print(level+1, printContext)
	}
	if node.Value == nil {
		outputNil(level+1, printContext)
	} else {
		node.Value.Print(level+1, printContext)
	}
	closeSection(level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *VariableReference) Print(level int, printContext *PrintContext) {
	outputLeveled("{"+outputBaseSource(node.Name, printContext)+"} "+outputType("("+node.Type.String()+")", printContext), level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *IfStmt) Print(level int, printContext *PrintContext) {
	openSection("if", level, printContext)
	if node.Init != nil {
		openSection("init", level+2, printContext)
		node.Init.Print(level+3, printContext)
		closeSection(level+2, printContext)
	}
	openSection("condition", level+2, printContext)
	if node.Conditional != nil {
		node.Conditional.Print(level+3, printContext)
	} else {
		outputNil(level+3, printContext)
	}
	closeSection(level+2, printContext)
	openSection("code", level+2, printContext)
	if node.Code != nil {
		node.Code.Print(level+3, printContext)
	} else {
		outputNil(level+3, printContext)
	}
	closeSection(level+2, printContext)
	if node.Else != nil {
		openSection("else", level+2, printContext)
		node.Else.Print(level+3, printContext)
		closeSection(level+2, printContext)
	}
	closeSection(level, printContext)
}

// Print writes a description of the node to standard output, at the specified indentation level.
func (node *ForStmt) Print(level int, printContext *PrintContext) {
	openSection("for", level, printContext)
	if node.Init != nil {
		openSection("init", level+2, printContext)
		node.Init.Print(level+3, printContext)
		closeSection(level+2, printContext)
	}
	openSection("condition", level+2, printContext)
	if node.Conditional != nil {
		node.Conditional.Print(level+3, printContext)
	} else {
		outputNil(level+3, printContext)
	}
	closeSection(level+2, printContext)
	openSection("code", level+2, printContext)
	if node.Code != nil {
		node.Code.Print(level+3, printContext)
	} else {
		outputNil(level+3, printContext)
	}
	closeSection(level+2, printContext)
	if node.PostIteration != nil {
		openSection("PostIteration", level+2, printContext)
		node.PostIteration.Print(level+3, printContext)
		closeSection(level+2, printContext)
	}
	closeSection(level, printContext)
}

func openSection(sectionName string, level int, printContext *PrintContext) {
	joiner := " {"
	if sectionName == "" {
		joiner = "{"
	}
	if printContext.Color {
		outputLeveled(yellow()+sectionName+blue()+joiner+resetColor(), level, printContext)

	} else {
		outputLeveled(sectionName+joiner, level, printContext)
	}
}
func closeSection(level int, printContext *PrintContext) {
	outputLeveled(blue()+"}"+resetColor(), level, printContext)
}

func outputNil(level int, printContext *PrintContext) {
	if printContext.Color {
		outputLeveled(red()+"NIL"+resetColor(), level, printContext)
	} else {
		outputLeveled("NIL", level, printContext)
	}
}

func outputBaseSource(str string, printContext *PrintContext) string {
	return str
}

func outputType(typeStr string, printContext *PrintContext) string {
	if printContext.Color {
		return cyan() + typeStr + resetColor()
	}
	return typeStr
}

func blue() string {
	return string([]byte{27, '[', '3', '4', 'm'})
}
func green() string {
	return string([]byte{27, '[', '3', '2', 'm'})
}
func red() string {
	return string([]byte{27, '[', '3', '1', 'm'})
}
func yellow() string {
	return string([]byte{27, '[', '3', '3', 'm'})
}
func magenta() string {
	return string([]byte{27, '[', '3', '5', 'm'})
}
func cyan() string {
	return string([]byte{27, '[', '3', '6', 'm'})
}
func white() string {
	return string([]byte{27, '[', '3', '7', 'm'})
}
func resetColor() string {
	return string([]byte{27, '[', '0', 'm'})
}
func ifColor(in string, ps *PrintContext) string {
	if ps.Color {
		return in
	}
	return ""
}

func outputLeveled(str string, level int, printContext *PrintContext) {
	for i := 0; i < level; i++ {
		fmt.Fprint(printContext.Output, " ")
	}
	fmt.Fprintln(printContext.Output, str)
}

func (op *UnOpType) String() string {
	switch *op {
	case UnOpNot:
		return "!"
	}
	return "UNOP?"
}

func (op *BinOpType) String() string {
	switch *op {
	case BinOpAdd:
		return "+"
	case BinOpSub:
		return "-"
	case BinOpMul:
		return "*"
	case BinOpDiv:
		return "/"
	case BinOpUnknown:
		return "UNK?"
	case BinOpMod:
		return "%"
	case BinOpEquality:
		return "=="
	case BinOpLAnd:
		return "&&"
	case BinOpLOr:
		return "||"
	case BinOpNotEquality:
		return "!="
	}
	return "BINOP?"
}

//Type system

func (t NamedType) String() string {
	return t.Ident + " " + t.Type.String()
}

func (t FunctionType) String() string {
	paramList := ""
	for i, p := range t.Parameters {
		paramList += p.String()
		if i+1 < len(t.Parameters) {
			paramList += ", "
		}
	}
	return "(" + paramList + ")" + t.ReturnType.String()
}

func (tk TypeKindDescription) String() string {
	switch tk {
	case PrimitiveTypeInt:
		return "int"
	case PrimitiveTypeString:
		return "string"
	case PrimitiveTypeBool:
		return "bool"
	case ComplexTypeArray:
		return "[?]"
	case PrimitiveTypeUndefined:
		return "undefined"
	}
	return "?"
}
