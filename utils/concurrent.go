package utils

import (
	"sync"
)

func Concurrent[T any](array []T, execute func(T), maxRequests int) {
	if maxRequests < 1 {
		panic("Deadlock can't process items")
	}
	wg := sync.WaitGroup{}
	maxConcurrency := make(chan struct{}, maxRequests)

	for index, item := range array {
		wg.Add(1)
		maxConcurrency <- struct{}{}

		go func(i int, item T) {
			defer wg.Done()
			execute(item)

			<-maxConcurrency
		}(index, item)
	}

	wg.Wait()
}
