package util

import (
	"errors"
)

const (
	EMPTY_STACK = "Empty Stack"
)

var ErrorEndOfStack = errors.New("end of stack")

type Stack[T any] struct {
	s []T
}

func (s *Stack[T]) Length() int {
	return len(s.s)
}
func (s *Stack[T]) GetSlice() []T {
	return s.s
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.s) == 0
}

func (s *Stack[T]) Clear() {
	s.s = make([]T, 0)
}

func (s *Stack[T]) Clone() *Stack[T] {
	return &Stack[T]{s.s}
}

func (s *Stack[T]) Get(i int) (T, error) {
	if i < 0 || i >= len(s.s) {
		return *new(T), ErrorEndOfStack
	}
	return s.s[i], nil
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
		return *new(T), ErrorEndOfStack
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (s *Stack[T]) Peek() (T, error) {
	l := len(s.s)
	if l == 0 {
		return *new(T), ErrorEndOfStack
	}
	return s.s[l-1], nil
}

func (s *Stack[T]) PushElements(v []T) {
	s.s = append(s.s, v...)
}
