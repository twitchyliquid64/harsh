package ast

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
	UNKNOWN_TYPE //Used internally to signify the type could be valid but is currently unknown
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
