package ast

type Node interface {
	Print(level int)
	Exec(context *ExecContext) *Variant
}

type StatementList struct {
	Stmts []Node
}

type IntegerLiteral struct {
	Val int64
}

type StringLiteral struct {
	Str string
}

type BoolLiteral struct {
	Val bool
}

type ArrayLiteral struct {
	Type    TypeKind
	Literal []Node
}

type NilLiteral struct {
}

type ReturnStmt struct {
	Expr Node
}

type IfStmt struct {
	Conditional Node
	Code        Node
	Init        Node
	Else        Node
}

type VariableReference struct {
	Name string
	Type TypeKind
}

type BinaryOp struct {
	LHS Node
	RHS Node
	Op  BinOpType
}

type UnaryOp struct {
	Op   UnOpType
	Expr Node
}

type Subscript struct {
	Subscript Node
	Expr      Node
}

type UnOpType int

const (
	UNOP_NOT UnOpType = iota
)

// BinaryOp Ops
type BinOpType int

const (
	BINOP_ADD BinOpType = iota
	BINOP_SUB
	BINOP_MUL
	BINOP_DIV
	BINOP_MOD

	BINOP_LAND
	BINOP_LOR

	BINOP_EQUALITY
	BINOP_UNK
)

type Assign struct {
	Value    Node
	Variable Node
	NewLocal bool
}
