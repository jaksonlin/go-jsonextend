package util

import (
	"errors"
)

const (
	EMPTY_STACK = "Empty Stack"
)

var ErrorEodOfStack = errors.New("end of stack")

type Stack[T any] struct {
	s []T
}

func (s *Stack[T]) Length() int {
	return len(s.s)
}
func (s *Stack[T]) GetSlice() []T {
	return s.s
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{make([]T, 0)}
}

func (s *Stack[T]) Push(v T) {
	s.s = append(s.s, v)
}

func (s *Stack[T]) Pop() (T, error) {
	l := len(s.s)
	if l == 0 {
		return *new(T), ErrorEodOfStack
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (s *Stack[T]) Peek() (T, error) {
	l := len(s.s)
	if l == 0 {
		return *new(T), ErrorEodOfStack
	}
	return s.s[l-1], nil
}

func (s *Stack[T]) PushElements(v []T) {
	s.s = append(s.s, v...)
}
