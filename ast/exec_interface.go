package ast

type ExecContext struct {
	IsFuncContext     bool
	FunctionNamespace Namespace
	GlobalNamespace   *Namespace
}

type Variant struct {
	Type     PrimitiveType
	Int      int64
	String   string
	IsReturn bool
}

type Namespace map[string]Variant

func (n *Namespace) Save(name string, v interface{}) {
	(*n)[name] = MakeVariant(v)
}
