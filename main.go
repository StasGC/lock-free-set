package main

import (
	"fmt"
	"log"
	"main/set"
	"time"
)

func main() {
	lockFreeSet := set.NewSetImpl[int]()

	lockFreeSet.Add(10)
	lockFreeSet.Remove(10)
	if !lockFreeSet.IsEmpty() {
		log.Fatal("set is not empty")
	}

	UsageLockFreeSet[int](lockFreeSet)

	iterator := lockFreeSet.Iterator()
	for i := 0; i < 5; i++ {
		value := i
		go func() {
			time.Sleep(time.Millisecond)
			lockFreeSet.Add(value)
		}()
	}

	time.Sleep(time.Millisecond)

	for j := 0; ; j++ {
		val, err := iterator.Next()
		if err != nil {
			break
		}
		fmt.Println(val.Value)
	}
}

func UsageLockFreeSet[T comparable](lockFreeSet set.Set[T]) {
	// some usage of lockFreeSet...
	fmt.Println(lockFreeSet)
}
