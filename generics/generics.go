package generics

import "fmt"

type Parser[T any] func(string) (T, string, error)

type Result[T any] struct {
	Val T
	Err error
}

func AcceptGeneric[T any](val T)Result[T]{
	return Result[T]{
		Val: val,
		Err: fmt.Errorf(""),
	}
}
