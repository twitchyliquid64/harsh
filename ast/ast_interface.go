package ast

type Node interface {
	Print(level int)
	Exec(context *ExecContext) Variant
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

type BinaryOp struct {
	LHS Node
	RHS Node
	Op  BinOpType
}

// BinaryOp Ops
type BinOpType int

const (
	BINOP_ADD BinOpType = iota
	BINOP_SUB
	BINOP_MUL
	BINOP_DIV
	BINOP_MOD
	BINOP_UNK
)

// Types

type TypeKind int

const (
	PRIMITIVE_TYPE_INT TypeKind = iota
	PRIMITIVE_TYPE_STRING
	PRIMITIVE_TYPE_UNDEFINED
)

type TypeDecl interface {
	String() string
}

//A kind of named primitive variable
type PrimitiveType struct {
	Kind TypeKind
	Name string
}
