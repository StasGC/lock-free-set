package atomicimpl

import (
	"sync/atomic"
	"unsafe"
)

type AtomicMarkableReference[T any] struct {
	p *unsafe.Pointer
}

type Pair[T any] struct {
	reference *T
	mark      bool
}

func NewAtomicMarkableReference[T any](reference *T, mark bool) *AtomicMarkableReference[T] {
	p := unsafe.Pointer(&Pair[T]{reference: reference, mark: mark})
	return &AtomicMarkableReference[T]{p: &p}
}

func (sr *AtomicMarkableReference[T]) GetReference() *T {
	item := (*Pair[T])(atomic.LoadPointer(sr.p))
	return item.reference
}

func (sr *AtomicMarkableReference[T]) IsMarked() bool {
	item := (*Pair[T])(atomic.LoadPointer(sr.p))
	return item.mark
}

func (sr *AtomicMarkableReference[T]) Get(markHolder []bool) *T {
	item := (*Pair[T])(atomic.LoadPointer(sr.p))

	markHolder[0] = item.mark
	return item.reference
}

func (sr *AtomicMarkableReference[T]) CompareAndSet(
	expectedReference *T,
	newReference *T,
	expectedMark bool,
	newMark bool,
) bool {
	expectedPointer := atomic.LoadPointer(sr.p)
	current := (*Pair[T])(expectedPointer)
	if current.reference == expectedReference && current.mark == expectedMark {
		return atomic.CompareAndSwapPointer(
			sr.p,
			expectedPointer,
			unsafe.Pointer(&Pair[T]{reference: newReference, mark: newMark}),
		)
	}
	return false
}

func (sr *AtomicMarkableReference[T]) AttemptMark(expectedReference *T, newMark bool) bool {
	expectedPointer := atomic.LoadPointer(sr.p)
	current := (*Pair[T])(expectedPointer)
	if current.reference == expectedReference && current.mark != newMark {
		return atomic.CompareAndSwapPointer(
			sr.p,
			expectedPointer,
			unsafe.Pointer(&Pair[T]{reference: expectedReference, mark: newMark}),
		)
	}
	return false
}

func (sr *AtomicMarkableReference[T]) Set(newReference *T, newMark bool) {
	expectedPointer := atomic.LoadPointer(sr.p)
	current := (*Pair[T])(expectedPointer)
	if newReference != current.reference || newMark != current.mark {
		p := unsafe.Pointer(&Pair[T]{reference: newReference, mark: newMark})
		sr.p = &p
	}
}
