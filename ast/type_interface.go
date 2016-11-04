package ast

// TypeKind is implemented by all Types which are represented in the AST.
// ConcreteType returns the underlying array/slice type if applicable, otherwise it returns the same value as Kind().
// Kind returns a value with represents the kind of value it is: ie int/string/slice/array.
type TypeKind interface {
	ConcreteType() TypeDecl
	Kind() TypeKindDescription
	String() string
}

// TypeKindDescription represents at the top level the kind of primitive handled.
type TypeKindDescription int

// Kind returns the kind of value the type is.
func (t TypeKindDescription) Kind() TypeKindDescription {
	return t
}

// ConcreteType returns the same value as Kind() for TypeKindDescription.
func (t TypeKindDescription) ConcreteType() TypeDecl {
	return t
}

// Represents the valid kinds in the type system.
const (
	PrimitiveTypeInt TypeKindDescription = iota
	PrimitiveTypeString
	PrimitiveTypeBool
	ComplexTypeArray
	PrimitiveTypeUndefined
	UnknownType //Used internally to signify the type could be valid but is currently unknown
)

// TypeDecl is a subset of TypeKind, but otherwise has the same meaning and semantics.
type TypeDecl interface {
	String() string
	ConcreteType() TypeDecl
}

// NamedType is a kind of named primitive variable, used mainly to represent named parameters.
type NamedType struct {
	Type  TypeKind
	Ident string
}

// ConcreteType returns the underlying type of the value.
func (p NamedType) ConcreteType() TypeDecl {
	return p.Type
}

// Kind returns the kind of value the type is.
func (p NamedType) Kind() TypeKindDescription {
	return p.Type.Kind()
}

// Name returns the name of the value.
func (p NamedType) Name() string {
	return p.Ident
}

// SetName sets the name of the value.
func (p NamedType) SetName(n string) {
	p.Ident = n
}

// ArrayType represents an array built on primitive of type SubType, with an array length represented by the execution of Len.
type ArrayType struct {
	SubType TypeKind
	Len     Node
}

func (a ArrayType) String() string {
	return "[]" + a.ConcreteType().String()
}

// Kind returns ComplexTypeArray.
func (a ArrayType) Kind() TypeKindDescription {
	return ComplexTypeArray
}

// ConcreteType returns the type of the underlying element primitive.
func (a ArrayType) ConcreteType() TypeDecl {
	return a.SubType
}
