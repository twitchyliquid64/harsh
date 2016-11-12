package ast

// Node represents any component of the AST, such as literals, operations, variables, loops, etc.
type Node interface {
	Print(level int)
	Exec(context *ExecContext) *Variant
}

// StatementList represents a list of nodes to be executed sequentially.
type StatementList struct {
	Stmts []Node
}

// IntegerLiteral represents a literal whole number.
type IntegerLiteral struct {
	Val int64
}

// StringLiteral represents a literal string.
type StringLiteral struct {
	Str string
}

// BoolLiteral represents a literal boolean (true/false).
type BoolLiteral struct {
	Val bool
}

// ArrayLiteral represents a composite of literals which initialize a variable.
type ArrayLiteral struct {
	Type    ArrayType
	Literal []Node
}

// StructLiteral represents a composite of named values which initialize a variable of type struct.
type StructLiteral struct {
	Type   StructType
	Values map[string]Node
}

// NilLiteral symbolizes an invalid construct, or simply a null value.
type NilLiteral struct {
}

// ReturnStmt represents a short-circuit of linear StatementList execution, returning a value down to the function level.
type ReturnStmt struct {
	Expr Node
}

// NamedSelector represents the fetch of a named set of data from the upstream data structure.
type NamedSelector struct {
	Expr Node
	Name string
}

// IfStmt represents conditional branching, evaluating a condition then taking various actions based on the result.
type IfStmt struct {
	Conditional Node
	Code        Node
	Init        Node
	Else        Node
}

// VariableReference represents the fetching of a value at runtime from a variable. If possible the runtime type
// is inferred and stored in the structure for the sake of typechecking.
type VariableReference struct {
	Name string
	Type TypeKind
}

// BinaryOp represents a binary operation between two operands.
type BinaryOp struct {
	LHS Node
	RHS Node
	Op  BinOpType
}

// UnaryOp represents a unary operation (EG: NOT or !), done on a single operand.
type UnaryOp struct {
	Op   UnOpType
	Expr Node
}

// UnOpType encapsulates the valid operations for UnaryOp.
type UnOpType int

const (
	// UnOpNot symbolizes the boolean NOT operation.
	UnOpNot UnOpType = iota
)

// Subscript is a node representing the access of an index from a array/slice at runtime to retrieve a value.
type Subscript struct {
	Subscript Node
	Expr      Node
}

// BinOpType encapsulates valid operations for BinaryOp.
type BinOpType int

// Represents the possible binary operations.
const (
	BinOpAdd BinOpType = iota
	BinOpSub
	BinOpMul
	BinOpDiv
	BinOpMod

	BinOpLAnd
	BinOpLOr

	BinOpEquality
	BinOpUnknown
)

// Assign represents storing a value into a variable construct at runtime.
type Assign struct {
	Value    Node
	Variable Node
	NewLocal bool
}
