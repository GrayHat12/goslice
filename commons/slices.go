package commons

import (
	"context"
	"fmt"
	"sync"
)

func Find[T any](slice *[]T, predicate func(*T, int, *[]T) bool) *T {
	for index, element := range *slice {
		if predicate(&element, index, slice) {
			return &element
		}
	}
	return nil
}

func FindIndex[T any](slice *[]T, predicate func(*T, int, *[]T) bool) int {
	for index, element := range *slice {
		if predicate(&element, index, slice) {
			return index
		}
	}
	return -1
}

type FilterCallback[T any] func(context.Context, *T, int, *[]T)

func ForEachWithoutWait[T any](wg *sync.WaitGroup, ctx context.Context, slice *[]T, callback FilterCallback[T]) {
	callbackWrapper := func(waitGroup *sync.WaitGroup, _callback func()) {
		defer waitGroup.Done()
		_callback()
	}
	for index, element := range *slice {
		wg.Add(1)
		go callbackWrapper(wg, func() {
			callback(ctx, &element, index, slice)
		})
	}
}

func ForEach[T any](ctx context.Context, slice *[]T, callback FilterCallback[T]) {
	wg := sync.WaitGroup{}
	ForEachWithoutWait(&wg, ctx, slice, callback)
	wg.Wait()
}

func Every[T any](slice *[]T, predicate func(context.Context, *T, int, *[]T) bool) bool {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	responses := make(chan bool, len(*slice))
	// defer close(responses)
	ForEachWithoutWait(wg, ctx, slice, func(_ctx context.Context, item *T, index int, array *[]T) {
		select {
		case <-ctx.Done():
			fmt.Printf("ctx done for %+v\n", item)
			return
		default:
			fmt.Printf("pushing chan\n")
			responses <- predicate(_ctx, item, index, array)
		}
	})
	returnVal := true
	responseIndex := 0
	for response := range responses {
		fmt.Printf("got chan %+v\n", response)
		if response == false {
			cancel()
			returnVal = false
			close(responses)
			break
		}
		responseIndex += 1
		if responseIndex == len(*slice) {
			close(responses)
			break
		}
	}
	cancel()
	wg.Wait()
	return returnVal
}

func Reduce[T any](slice *[]T, callback func(*T, *T, int, *[]T) *T, initialValue *T) *T {
	accumulatedValue := initialValue
	for index, element := range *slice {
		accumulatedValue = callback(accumulatedValue, &element, index, slice)
	}
	return accumulatedValue
}

func Some[T any](slice *[]T, predicate func(context.Context, *T, int, *[]T) bool) bool {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	responses := make(chan bool, len(*slice))
	// defer close(responses)
	ForEachWithoutWait(wg, ctx, slice, func(_ctx context.Context, item *T, index int, array *[]T) {
		select {
		case <-ctx.Done():
			fmt.Printf("ctx done for %+v\n", item)
			return
		default:
			fmt.Printf("pushing chan\n")
			responses <- predicate(_ctx, item, index, array)
		}
	})
	returnVal := false
	responseIndex := 0
	for response := range responses {
		fmt.Printf("got chan %+v\n", response)
		if response == true {
			cancel()
			returnVal = true
			close(responses)
			break
		}
		responseIndex += 1
		if responseIndex == len(*slice) {
			close(responses)
			break
		}
	}
	cancel()
	wg.Wait()
	return returnVal
}
