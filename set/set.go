package set

import (
	atomicimpl "main/atomic_implements"
)

type SetImpl[T comparable] struct {
	Head *node[T]
	Tail *node[T]
}

func NewSetImpl[T comparable]() *SetImpl[T] {
	var headValue, tailValue T
	set := &SetImpl[T]{
		Head: newNode[T](headValue),
		Tail: newNode[T](tailValue),
	}
	set.Head.Next = atomicimpl.NewAtomicMarkableReference[node[T]](set.Tail, false)
	set.Tail.Next = atomicimpl.NewAtomicMarkableReference[node[T]](nil, false)
	return set
}

func (set *SetImpl[T]) Find(value T) (*node[T], *node[T]) {
	marked := []bool{false}

Loop:
	for {
		prev := set.Head
		curr := prev.Next.GetReference()
		for {
			next := curr.Next.Get(marked)
			for marked[0] {
				check := prev.Next.CompareAndSet(curr, next, false, false)
				if !check {
					continue Loop
				}

				curr = next
				next = curr.Next.Get(marked)
			}
			if curr == set.Tail || curr.Value == value {
				return prev, curr
			}

			prev = curr
			curr = next
		}
	}
}

func (set *SetImpl[T]) Add(value T) bool {
	for {
		prev, curr := set.Find(value)
		if curr.Value == value && curr != set.Tail {
			return false
		}

		newNode := newNode[T](value)
		newNode.Next = atomicimpl.NewAtomicMarkableReference[node[T]](set.Tail, false)
		if prev.Next.CompareAndSet(curr, newNode, false, false) {
			return true
		}
	}
}

func (set *SetImpl[T]) Remove(value T) bool {
	for {
		prev, curr := set.Find(value)
		if curr == set.Tail {
			return false
		}

		next := curr.Next.GetReference()
		if !curr.Next.AttemptMark(next, true) {
			continue
		}

		prev.Next.CompareAndSet(curr, next, false, false)
		return true
	}
}

func (set *SetImpl[T]) Contains(value T) bool {
	curr := set.Head
	for curr != set.Tail && curr.Value != value {
		curr = curr.Next.GetReference()
	}

	mark := curr.Next.IsMarked()
	return curr.Value == value && !mark
}

func (set *SetImpl[T]) IsEmpty() bool {
	for {
		prev, curr := set.Find(set.Tail.Value)

		prevMarked := prev.Next.IsMarked()
		currMarked := curr.Next.IsMarked()

		if curr != set.Tail {
			if currMarked {
				continue
			}
			return false
		}

		if prev != set.Head && prevMarked {
			continue
		}
		return prev == set.Head || prevMarked
	}
}

type SetIterator[T comparable] struct {
	curr             *node[T]
	lastNode         *node[T]
	callbackRemove   func(value T) bool
	isCurrentRemoved bool
}

func (set *SetImpl[T]) Iterator() SetIterator[T] {
	var snapShotHead, snapShotTail *node[T]

	for {
		snapShotHead, snapShotTail = set.getSnapShot()
		snapShotHeadCheck, snapShotTailCheck := set.getSnapShot()

		currentSnapShotNode := snapShotHead.Next.GetReference()
		currentSnapShotCheckNode := snapShotHeadCheck.Next.GetReference()
		for currentSnapShotNode != snapShotTail || currentSnapShotCheckNode != snapShotTailCheck {
			if currentSnapShotNode.Value != currentSnapShotCheckNode.Value {
				break
			}
			currentSnapShotNode = currentSnapShotNode.Next.GetReference()
			currentSnapShotCheckNode = currentSnapShotCheckNode.Next.GetReference()
		}
		if currentSnapShotNode == snapShotTail && currentSnapShotCheckNode == snapShotTailCheck {
			break
		}
	}

	return SetIterator[T]{
		curr:             snapShotHead,
		lastNode:         snapShotTail,
		callbackRemove:   set.Remove,
		isCurrentRemoved: false,
	}
}

func (iterator *SetIterator[T]) Next() (node[T], error) {
	var nilNode node[T]
	if iterator.curr == iterator.lastNode {
		return nilNode, NoSuchElementException
	}

	marked := []bool{false}
	curr := iterator.curr.Next.GetReference()
	next := curr.Next.Get(marked)
	for marked[0] {
		curr = next
		next = curr.Next.Get(marked)
	}

	if curr == iterator.lastNode {
		return nilNode, NoSuchElementException
	}

	iterator.curr = curr
	iterator.isCurrentRemoved = false
	return *curr, nil
}

func (iterator *SetIterator[T]) HasNext() bool {
	if iterator.curr == iterator.lastNode {
		return false
	}

	marked := []bool{false}
	curr := iterator.curr.Next.GetReference()
	next := curr.Next.Get(marked)
	for marked[0] {
		curr = next
		next = curr.Next.Get(marked)
	}

	if curr == iterator.lastNode {
		return false
	}
	return true
}

func (iterator *SetIterator[T]) Remove() error {
	if iterator.isCurrentRemoved {
		return IllegalStateException
	}
	iterator.callbackRemove(iterator.curr.Value)
	iterator.isCurrentRemoved = true
	return nil
}

func (set *SetImpl[T]) getSnapShot() (*node[T], *node[T]) {
	snapShotHead := newNode[T](set.Head.Value)

	currentSnapShotNode := snapShotHead
	currentNode := set.Head
	for currentNode != set.Tail {
		currentNode = currentNode.Next.GetReference()

		nodeCopy := newNode[T](currentNode.Value)
		currentSnapShotNode.Next = atomicimpl.NewAtomicMarkableReference[node[T]](nodeCopy, false)
		currentSnapShotNode = nodeCopy
	}
	currentSnapShotNode.Next = atomicimpl.NewAtomicMarkableReference[node[T]](nil, false)

	snapShotTail := currentSnapShotNode

	return snapShotHead, snapShotTail
}
