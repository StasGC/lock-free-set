package atomicimpl

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type testStruct struct {
	intValue    int
	stringValue string
}

func TestAtomicMarkableReference_CompareAndSet(t *testing.T) {
	var threadsCount int64 = 5
	var successCompareAndSetTimes int64 = 0
	var expectedNewValueTimes int64 = 0

	currentStruct := testStruct{
		intValue:    1,
		stringValue: "first",
	}
	newStruct := testStruct{
		intValue:    2,
		stringValue: "second",
	}

	markableReference := NewAtomicMarkableReference[testStruct](&currentStruct, false)
	wg := sync.WaitGroup{}
	for i := int64(0); i < threadsCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			time.Sleep(time.Microsecond)

			if markableReference.CompareAndSet(&currentStruct, &newStruct, false, false) {
				atomic.AddInt64(&successCompareAndSetTimes, 1)
			} else {
				if *markableReference.GetReference() == newStruct {
					atomic.AddInt64(&expectedNewValueTimes, 1)
				}
			}
		}()
	}

	wg.Wait()
	if successCompareAndSetTimes != 1 {
		t.Fatalf("wrong success compareAndSet times: %v", successCompareAndSetTimes)
	}
	if expectedNewValueTimes != threadsCount-1 {
		t.Fatalf("wrong expected new value times: %v", successCompareAndSetTimes)
	}
}

func TestAtomicMarkableReference_AttemptMark(t *testing.T) {
	var threadsCount int64 = 5
	var successAttemptMarkTimes int64 = 0
	var expectedSeenValueTimes int64 = 0

	currentStruct := testStruct{
		intValue:    1,
		stringValue: "first",
	}

	markableReference := NewAtomicMarkableReference[testStruct](&currentStruct, false)
	wg := sync.WaitGroup{}
	for i := int64(0); i < threadsCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			time.Sleep(time.Microsecond)

			if markableReference.AttemptMark(&currentStruct, true) {
				atomic.AddInt64(&successAttemptMarkTimes, 1)
			}
			if *markableReference.GetReference() == currentStruct {
				atomic.AddInt64(&expectedSeenValueTimes, 1)
			}
		}()
	}

	wg.Wait()
	if successAttemptMarkTimes != 1 {
		t.Fatalf("wrong success compareAndSet times: %v", successAttemptMarkTimes)
	}
	if expectedSeenValueTimes != threadsCount {
		t.Fatalf("wrong expected new value times: %v", expectedSeenValueTimes)
	}
}
