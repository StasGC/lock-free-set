package set

type Set[T comparable] interface {
	Add(value T) bool
	Remove(value T) bool
	Contains(value T) bool
	IsEmpty() bool
	Iterator() SetIterator[T]
}

type Iterator[K comparable] interface {
	Next() (node[K], error)
	HasNext() bool
	Remove() error
}
