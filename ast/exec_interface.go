package ast

type ExecContext struct {
	IsFuncContext bool
}

type Variant struct {
	Type     PrimitiveType
	Int      int64
	String   string
	IsReturn bool
}
