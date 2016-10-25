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

type StringLiteral struct {
	Str string
}

type BoolLiteral struct {
	Val bool
}

type ArrayLiteral struct {
	Type    TypeKind
	Literal []Variant
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

type Assign struct {
	Value      Node
	Identifier string
	NewLocal   bool
}

// Types

type TypeKind interface {
	ConcreteType() TypeDecl
	Kind() TypeKindDescription
	String() string
}

type TypeKindDescription int

func (t TypeKindDescription) Kind() TypeKindDescription {
	return t
}
func (t TypeKindDescription) ConcreteType() TypeDecl {
	return t
}

const (
	PRIMITIVE_TYPE_INT TypeKindDescription = iota
	PRIMITIVE_TYPE_STRING
	PRIMITIVE_TYPE_BOOL
	COMPLEX_TYPE_ARRAY
	COMPLEX_TYPE_VECTOR
	PRIMITIVE_TYPE_UNDEFINED
)

type TypeDecl interface {
	String() string
	ConcreteType() TypeDecl
}

//A kind of named primitive variable
type NamedType struct {
	Type  TypeKind
	Ident string
}

func (p NamedType) ConcreteType() TypeDecl {
	return p.Type
}
func (p NamedType) Kind() TypeKindDescription {
	return p.Type.Kind()
}
func (p NamedType) Name() string {
	return p.Ident
}
func (p NamedType) SetName(n string) {
	p.Ident = n
}

type ArrayType struct {
	SubType TypeKind
	Len     Node
}

func (a ArrayType) String() string {
	return "[]" + a.ConcreteType().String()
}
func (a ArrayType) Kind() TypeKindDescription {
	return COMPLEX_TYPE_ARRAY
}
func (a ArrayType) ConcreteType() TypeDecl {
	return a.SubType
}
