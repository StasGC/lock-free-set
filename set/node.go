package set

import (
	atomicimpl "main/atomic_implements"
)

type node[T comparable] struct {
	Value T
	Next  *atomicimpl.AtomicMarkableReference[node[T]]
}

func newNode[T comparable](value T) *node[T] {
	return &node[T]{Value: value}
}
