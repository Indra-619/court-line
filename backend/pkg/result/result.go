package result

import "errors"

// Result is a generic container for a value or an error
type Result[T any] struct {
	Value T
	Error error
}

// Success creates a successful result
func Success[T any](value T) Result[T] {
	return Result[T]{
		Value: value,
		Error: nil,
	}
}

// Failure creates a failed result
func Failure[T any](err error) Result[T] {
	return Result[T]{
		Error: err,
	}
}

// Errs creates a failed result with a string error
func Errs[T any](msg string) Result[T] {
	return Result[T]{
		Error: errors.New(msg),
	}
}

// IsSuccess returns true if the result is successful
func (r Result[T]) IsSuccess() bool {
	return r.Error == nil
}

// IsFailure returns true if the result is a failure
func (r Result[T]) IsFailure() bool {
	return r.Error != nil
}

// Unwrap returns the value or panics if there is an error
func (r Result[T]) Unwrap() T {
	if r.Error != nil {
		panic(r.Error)
	}
	return r.Value
}
