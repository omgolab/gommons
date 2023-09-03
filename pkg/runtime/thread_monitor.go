package gcruntime

import (
	"runtime"
	"sync"
	"time"
)

// The function monitors and reduces the number of
// idle OS threads in order to stay within a maximum limitation. This is a
// temporary solution to reduce the number of M idle OS threads.
// Open issue:
// https://github.com/golang/go/issues/14592
// https://github.com/golang/go/issues/20395
// possible pitfall:
// https://github.com/golang/go/issues/14592#issuecomment-693186098
func MonitorAndReduceIdleOSThreads(threadMonitorEnabled bool, timeoutSec, rateLimit int) {
	if !threadMonitorEnabled {
		return
	}

	// The default value is 60 seconds and minimum value is 5 seconds
	if timeoutSec <= 5 {
		timeoutSec = 5
	}

	// The minimum rate limit is 1 thread
	if rateLimit <= 0 {
		rateLimit = 1
	}

	go func() {
		var wg sync.WaitGroup
		ticker := time.NewTicker(time.Duration(timeoutSec) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			mThreadNum, _ := runtime.ThreadCreateProfile(nil)
			if mThreadNum <= rateLimit {
				return
			}

			if rateLimit > 0 {
				wg.Add(rateLimit)
				for i := 0; i < rateLimit; i++ {
					go func() {
						runtime.LockOSThread()
						defer wg.Done()
						// this cause system exit rather than OS thread exit
						// syscall.Syscall(syscall.SYS_EXIT, 0, 0, 0)
					}()
				}
				wg.Wait()
			}
		}
	}()
}
