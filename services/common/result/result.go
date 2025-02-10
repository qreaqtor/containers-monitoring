package result

type Result[T any] struct {
	Error error
	Value T
}

func NewValueResult[T any](value T) Result[T] {
	return Result[T]{Value: value}
}

func NewErrorResult[T any](err error) Result[T] {
	return Result[T]{Error: err}
}
