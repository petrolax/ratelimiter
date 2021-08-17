package ratelimiter

import (
	"sync"
	"sync/atomic"
	"time"
)

type RateLimiter struct {
	maxPerMinute   int32
	countPerMinute int32
	maxNow         int32
	countNow       int32
	num            chan int
	total          int64
}

func NewRateLimiter(data chan int, now, perMinute int32) *RateLimiter {
	if now == 0 || now < 0 {
		now = 30
	}

	if perMinute == 0 || perMinute < 0 {
		perMinute = 100
	}

	return &RateLimiter{
		num:          data,
		maxNow:       now,
		maxPerMinute: perMinute,
	}
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

			if atomic.LoadInt32(&rt.countNow) == rt.maxNow || rt.countPerMinute == rt.maxPerMinute {
				continue
			}

			wg.Add(1)
			atomic.AddInt32(&rt.countNow, int32(1))
			rt.countPerMinute++
			go rt.worker(val, &wg)
		case <-ticker.C:
			rt.countPerMinute = atomic.LoadInt32(&rt.countNow)
		}
	}

}

func (rt *RateLimiter) worker(value int, wg *sync.WaitGroup) {
	atomic.AddInt64(&rt.total, int64(value))
	atomic.AddInt32(&rt.countNow, int32(-1))
	wg.Done()
}
