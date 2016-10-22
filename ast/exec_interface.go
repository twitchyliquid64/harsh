package ast

type ExecContext struct {
	IsFuncContext     bool
	FunctionNamespace map[string]Variant
}

type Variant struct {
	Type     PrimitiveType
	Int      int64
	String   string
	IsReturn bool
}
