package ast

type ExecContext struct {
	IsFuncContext     bool
	FunctionNamespace Namespace
	GlobalNamespace   Namespace
	Errors            []ExecutionError
}

type Namespace map[string]*Variant

// Save constructs a Variant from v and saves it. If a *variant is given, a shallow copy is performed.
func (n Namespace) Save(name string, v interface{}) {
	n[name] = MakeVariant(v) //makes a copy
}

func (n *Namespace) Names() []string {
	var o []string
	for name, _ := range *n {
		o = append(o, name)
	}
	return o
}
