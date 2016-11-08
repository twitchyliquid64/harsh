package ast

import "errors"

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
	case TypeKind:
		return &Variant{
			Type: v,
		}
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

// DefaultVariantValue returns a valid *Variant setup with the given type, and the appropriate default values.
func DefaultVariantValue(t TypeKind) (*Variant, error) {
	ret := &Variant{
		Type: t,
	}

	switch t.Kind() {
	//default values are fine
	case PrimitiveTypeInt:
	case PrimitiveTypeString:
	case PrimitiveTypeUndefined:
	case PrimitiveTypeBool:
	case ComplexTypeArray:
		context := &ExecContext{}
		arrayLen := 0
		lenEval := t.(ArrayType).Len.Exec(context)

		if len(context.Errors) == 0 && lenEval.Type == PrimitiveTypeInt {
			arrayLen = int(lenEval.Int)
			ret.VectorData = make([]*Variant, arrayLen)
			for i := 0; i < arrayLen; i++ {
				v, e := DefaultVariantValue(t.BaseType())
				if e != nil {
					return ret, errors.New("Array basetype error: " + e.Error())
				}
				ret.VectorData[i] = v
			}
		} else if len(context.Errors) != 0 {
			return ret, errors.New("Could not statically resolve the length of the given array")
		} else {
			return ret, errors.New("Resolved length of array was not an integer")
		}
	}
	return ret, nil
}
