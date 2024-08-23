package set

import "fmt"

var (
	NoSuchElementException = fmt.Errorf("no such elements in Set")
	IllegalStateException  = fmt.Errorf(
		"the next() method has not yet been called " +
			"or remove() has already been called since the last call to next()",
	)
)
