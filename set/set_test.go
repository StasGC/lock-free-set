package set

import (
	"sync"
	"testing"
	"time"
)

func TestSetImpl_Add(t *testing.T) {
	threadsCount := 5
	set := NewSetImpl[int]()

	var values []int
	var sharedValue int = 20
	for i := 1; i < threadsCount+1; i++ {
		values = append(values, i)
	}

	wg := sync.WaitGroup{}
	for _, value := range values {
		wg.Add(1)
		go func(currValue int) {
			defer wg.Done()
			time.Sleep(time.Microsecond)

			set.Add(currValue)
			set.Add(sharedValue)
		}(value)
	}

	wg.Wait()

	if set.IsEmpty() {
		t.Fatalf("set should not be empty")
	}

	for _, value := range values {
		if !set.Contains(value) {
			t.Fatalf("set did not contains value: %v", value)
		}
	}

	if !set.Contains(sharedValue) {
		t.Fatalf("set did not contains shared value: %v", sharedValue)
	}
	set.Remove(sharedValue)
	if set.Contains(sharedValue) {
		t.Fatalf("set contains shared value, but should not contains value: %v", sharedValue)
	}

}

func TestSetImpl_Remove(t *testing.T) {
	threadsCount := 5
	set := NewSetImpl[int]()

	var values []int
	var sharedValue int = 20
	for i := 1; i < threadsCount+1; i++ {
		values = append(values, i)
		set.Add(i)
	}
	set.Add(sharedValue)

	if set.IsEmpty() {
		t.Fatalf("set should not be empty")
	}

	wg := sync.WaitGroup{}
	for _, value := range values {
		wg.Add(1)
		go func(currValue int) {
			defer wg.Done()
			time.Sleep(time.Microsecond)

			set.Remove(currValue)
			set.Remove(sharedValue)
		}(value)
	}

	wg.Wait()

	if !set.IsEmpty() {
		t.Fatalf("set should be empty")
	}
}

func TestSetImpl_Iterator(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7}
	valuesToRemove := []int{1, 2, 4, 7}
	iterValues := map[int]bool{}

	set := NewSetImpl[int]()
	for _, value := range values {
		set.Add(value)
	}

	iterator := set.Iterator()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, val := range valuesToRemove {
			time.Sleep(time.Millisecond)
			set.Remove(val)
		}
	}()

	for {
		node, err := iterator.Next()
		if err != nil {
			break
		}

		iterValues[node.Value] = true
	}

	wg.Wait()

	for _, val := range values {
		if _, ok := iterValues[val]; !ok {
			t.Fatalf("value not in iterator: %v", val)
		}
	}
}

func TestSetImpl_Iterator_Remove(t *testing.T) {
	set := NewSetImpl[int]()
	set.Add(10)

	iterator := set.Iterator()
	_, err := iterator.Next()
	if err != nil {
		t.Fatalf("there shouldn't be an error: %v", err)
	}

	err = iterator.Remove()
	if err != nil {
		t.Fatalf("there shouldn't be an error: %v", err)
	}

	err = iterator.Remove()
	if err != IllegalStateException {
		t.Fatalf("error shold be equal to IllegalStateException. Error: %v", err)
	}

	if !set.IsEmpty() {
		t.Fatalf("iterator remove failed, set should be empty")
	}
}
