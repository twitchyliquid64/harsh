package ast

type Variant struct {
	Type       TypeKind
	Int        int64
	String     string
	Bool       bool
	IsReturn   bool
	VectorData []Variant
}

func MakeVariant(in interface{}) Variant {
	switch v := in.(type) {
	case Variant:
		return v
	case int:
		return Variant{
			Type: PRIMITIVE_TYPE_INT,
			Int:  int64(v),
		}
	case int64:
		return Variant{
			Type: PRIMITIVE_TYPE_INT,
			Int:  v,
		}
	case bool:
		return Variant{
			Type: PRIMITIVE_TYPE_BOOL,
			Bool: v,
		}
	case string:
		return Variant{
			Type:   PRIMITIVE_TYPE_STRING,
			String: v,
		}
	}

	return Variant{
		Type: PRIMITIVE_TYPE_UNDEFINED,
	}
}
