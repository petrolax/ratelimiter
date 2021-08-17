package ratelimiter

import (
	"fmt"
	"sync"
	"testing"
)

func TestRateLimiter(t *testing.T) {
	intch := make(chan int)
	rt := NewRateLimiter(intch, 10, 100)

	go func() {
		var wg sync.WaitGroup

		for i := 0; i < 10000; i++ {
			wg.Add(1)
			go func(v int, wg *sync.WaitGroup) {
				intch <- v
				wg.Done()
			}(i, &wg)
		}
		wg.Wait()
		close(intch)
	}()

	fmt.Println(rt.Run())
}
