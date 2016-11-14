package ast

type errClass int

// Represents possible classes of execution errors.
const (
	TypeErr errClass = iota
	BoundsErr
	NotFoundErr
	InvalidAst
	InternalErr
	NotImplementedErr
)

// ExecutionError encapsulates errors encountered while executing the AST at runtime.
type ExecutionError struct {
	Class        errClass
	CreatingNode Node
	Text         string
}

func (e ExecutionError) Error() string {
	return e.Text
}
