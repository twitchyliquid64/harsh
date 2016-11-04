package ast

// Variant represents a value at runtime.
type Variant struct {
	Type                    TypeKind
	Int                     int64
	String                  string
	Bool                    bool
	IsReturn                bool
	VariableReferenceFailed bool
	VectorData              []*Variant
}

// MakeVariant takes a value of type *Variant or a go primitive (int/int64/bool/string) and constructs a *Variant.
func MakeVariant(in interface{}) *Variant {
	switch v := in.(type) {
	case *Variant:
		temp := *v
		temp.IsReturn = false
		temp.VariableReferenceFailed = false
		return &temp
	case int:
		return &Variant{
			Type: PrimitiveTypeInt,
			Int:  int64(v),
		}
	case int64:
		return &Variant{
			Type: PrimitiveTypeInt,
			Int:  v,
		}
	case bool:
		return &Variant{
			Type: PrimitiveTypeBool,
			Bool: v,
		}
	case string:
		return &Variant{
			Type:   PrimitiveTypeString,
			String: v,
		}
	}

	return &Variant{
		Type: PrimitiveTypeUndefined,
	}
}
