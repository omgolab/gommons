package gtime

import (
	"time"
)

// Goal: Track and wait for a certain time with a starting time
type Waiter struct {
	// start time
	start time.Time
	// goal time
	goal int64
}

// Create a new Waiter
func NewWaiter(goalDurationInMillis int64) *Waiter {
	return &Waiter{
		start: time.Now(),
		goal:  goalDurationInMillis,
	}
}

// Wait for the remaining time of goal time
func (w *Waiter) Wait(overheadMillis ...int64) int64 {
	// sum up the extra durations
	var overhead time.Duration
	for _, d := range overheadMillis {
		overhead += time.Duration(d * int64(time.Millisecond))
	}
	st := time.Since(w.start)
	rt := time.Duration(w.goal*int64(time.Millisecond)) - st - overhead
	time.Sleep(rt)

	return rt.Milliseconds()
}
