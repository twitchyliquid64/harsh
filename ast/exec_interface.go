package ast

// ExecContext is a structure passed to AST nodes during execution to contain namespaces or contextualise behaviour.
type ExecContext struct {
	IsFuncContext     bool
	FunctionNamespace Namespace
	GlobalNamespace   Namespace
	Errors            []ExecutionError
}

// Namespace represents a mapping of (variable) names to values.
type Namespace map[string]*Variant

// Save constructs a Variant from v and saves it. If a *variant is given, a shallow copy is performed.
func (n Namespace) Save(name string, v interface{}) {
	n[name] = MakeVariant(v) //makes a copy
}

// Names returns a list of all the names (variables) in the namespace.
func (n *Namespace) Names() []string {
	var o []string
	for name := range *n {
		o = append(o, name)
	}
	return o
}
