package ast

type errClass int

const (
	TYPE_ERR errClass = iota
	INVALID_AST
)

type ExecutionError struct {
	Class        errClass
	CreatingNode Node
	Text         string
}

func (e ExecutionError) Error() string {
	return e.Text
}
