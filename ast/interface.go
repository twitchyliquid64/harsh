package ast

type Node interface {
	Print(level int)
}

type StatementList struct {
	Stmts []Node
}

type IntegerLiteral struct {
	Val int64
}

type ReturnStmt struct {
	Expr Node
}

// Types

const (
	PRIMITIVE_TYPE_INT    = 0
	PRIMITIVE_TYPE_STRING = 1
)

type TypeDecl interface {
	String() string
}

type PrimitiveType struct {
	Kind int
	Name string
}
