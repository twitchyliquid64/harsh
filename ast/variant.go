package ast

func MakeVariant(in interface{}) Variant {
	switch v := in.(type) {
	case Variant:
		return v
	case int:
		return Variant{
			Type: PrimitiveType{
				Kind: PRIMITIVE_TYPE_INT,
			},
			Int: int64(v),
		}
	case int64:
		return Variant{
			Type: PrimitiveType{
				Kind: PRIMITIVE_TYPE_INT,
			},
			Int: v,
		}
	case string:
		return Variant{
			Type: PrimitiveType{
				Kind: PRIMITIVE_TYPE_STRING,
			},
			String: v,
		}
	}

	return Variant{
		Type: PrimitiveType{
			Kind: PRIMITIVE_TYPE_UNDEFINED,
		},
	}
}
