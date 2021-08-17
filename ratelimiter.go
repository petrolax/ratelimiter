package ratelimiter

import (
	"sync"
	"sync/atomic"
	"time"
)

type RateLimiter struct {
	maxPerMinut   int32
	countPerMinut int32
	maxNow        int32
	countNow      int32
	num           chan int
	total         int64
	mu            sync.RWMutex
}

func NewRateLimiter(data chan int) *RateLimiter {
	rt := &RateLimiter{
		num:         data,
		maxNow:      30,
		maxPerMinut: 100,
	}

	return rt
}

func (rt *RateLimiter) Run() (int64, string) {
	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case val, ok := <-rt.num:
			if !ok {
				wg.Wait()
				return rt.total, "RateLimiter was stop"
			}

			if atomic.LoadInt32(&rt.countNow) == rt.maxNow || rt.countPerMinut == rt.maxPerMinut {
				continue
			}

			// if rt.countPerMinut == rt.maxPerMinut {
			// 	continue
			// }

			wg.Add(1)
			atomic.AddInt32(&rt.countNow, int32(1))
			rt.countPerMinut++
			go rt.worker(val, &wg)
		case <-ticker.C:
			rt.countPerMinut = atomic.LoadInt32(&rt.countNow)
		}
	}

}

func (rt *RateLimiter) worker(value int, wg *sync.WaitGroup) {
	atomic.AddInt64(&rt.total, int64(value))
	atomic.AddInt32(&rt.countNow, int32(-1))
	wg.Done()
}
