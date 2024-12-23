package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/GrayHat12/goslice/commons"
	"github.com/GrayHat12/goslice/outofplace"
)

func testroutine(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("forcing cancel")
			return
		default:
			fmt.Println("still running")
			time.Sleep(4 * time.Second)
		}
	}
}

func main() {
	// ctx, cancel := context.WithCancel(context.Background())
	// wg := &sync.WaitGroup{}
	// wg.Add(1)
	// go testroutine(wg, ctx)
	// go func() {
	// 	time.Sleep(6 * time.Second)
	// 	fmt.Println("cancelling")
	// 	cancel()
	// }()
	// wg.Wait()
	// fmt.Println("Finish")
	list := []int{2, 4, 10, 24, 0}
	allEven := commons.Every(&list, func(_ context.Context, item *int, index int, _ *[]int) bool {
		return (*item)%2 == 0
	})
	fmt.Printf("allEven=%+v\n", allEven)
	fmt.Println("==========================================")
	allOdd := commons.Every(&list, func(_ context.Context, item *int, index int, _ *[]int) bool {
		return (*item)%2 != 0
	})
	fmt.Printf("allOdd=%+v\n", allOdd)
	fmt.Println("==========================================")
	reduced := 0
	reduced = *commons.Reduce(&list, func(sum *int, val *int, _ int, _ *[]int) *int {
		out := (*sum) + (*val)
		return &out
	}, &reduced)
	fmt.Printf("reduced=%+v\n", reduced)
	mapped := outofplace.Map(list, func(v *int, _ int, _ *[]int) *string {
		a := fmt.Sprintf("o%d", *v)
		return &a
	})
	fmt.Printf("mapped=%+v\n", mapped)
}
