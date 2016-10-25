package ast

type ExecContext struct {
	IsFuncContext     bool
	FunctionNamespace Namespace
	GlobalNamespace   Namespace
	Errors            []ExecutionError
}

type Variant struct {
	Type     PrimitiveType
	Int      int64
	String   string
	Bool     bool
	IsReturn bool
}

type Namespace map[string]Variant

func (n Namespace) Save(name string, v interface{}) {
	n[name] = MakeVariant(v)
}

func (n *Namespace) Names() []string {
	var o []string
	for name, _ := range *n {
		o = append(o, name)
	}
	return o
}
