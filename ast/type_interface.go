package ast

// TypeKind is implemented by all Types which are represented in the AST.
// BaseType returns the underlying array/slice type if applicable, otherwise it returns the same value as Kind().
// Kind returns a value with represents the kind of value it is: ie int/string/slice/array.
type TypeKind interface {
	BaseType() TypeKind
	Kind() TypeKindDescription
	String() string
}

// TypeKindDescription represents at the top level the kind of primitive handled.
type TypeKindDescription int

// Kind returns the kind of value the type is.
func (t TypeKindDescription) Kind() TypeKindDescription {
	return t
}

// BaseType returns the same value as Kind() for TypeKindDescription.
func (t TypeKindDescription) BaseType() TypeKind {
	return t
}

// Represents the valid kinds in the type system.
const (
	PrimitiveTypeInt TypeKindDescription = iota
	PrimitiveTypeString
	PrimitiveTypeBool
	ComplexTypeArray
	ComplexTypeStruct
	ComplexTypeFunction
	PrimitiveTypeUndefined
	UnknownType //Used internally to signify the type could be valid but is currently unknown
)

// NamedType is a kind of named primitive variable, used mainly to represent named parameters.
type NamedType struct {
	Type  TypeKind
	Ident string
}

// BaseType returns the underlying type of the value.
func (p NamedType) BaseType() TypeKind {
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
	return "[]" + a.BaseType().String()
}

// Kind returns ComplexTypeArray.
func (a ArrayType) Kind() TypeKindDescription {
	return ComplexTypeArray
}

// BaseType returns the type of the underlying element primitive.
func (a ArrayType) BaseType() TypeKind {
	return a.SubType
}

// StructType represents a named set of fields contained within one structure.
type StructType struct {
	Fields []NamedType
}

func (a StructType) String() string {
	out := "struct{"
	for i, f := range a.Fields {
		out += f.String()
		if i+1 < len(a.Fields) {
			out += ", "
		}
	}
	return out + "}"
}

// Kind returns ComplexTypeArray.
func (a StructType) Kind() TypeKindDescription {
	return ComplexTypeStruct
}

// BaseType returns ComplexTypeStruct as there is no real base type.
func (a StructType) BaseType() TypeKind {
	return ComplexTypeStruct //no real base type
}

// FunctionType represents the parameters, return type and code node of a function.
type FunctionType struct {
	Parameters []TypeKind
	ReturnType TypeKind
	Code       Node
}

// Kind returns ComplexTypeFunction.
func (a FunctionType) Kind() TypeKindDescription {
	return ComplexTypeFunction
}

// BaseType returns ComplexTypeFunction as there is no real base type.
func (a FunctionType) BaseType() TypeKind {
	return ComplexTypeFunction //no real base type
}
