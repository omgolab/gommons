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
// https://github.com/golang/go/issues/14592 possible pitfall:
// https://github.com/golang/go/issues/14592#issuecomment-693186098
func MonitorAndReduceIdleOSThreads(threadMonitorEnabled bool, timeoutSec, maxLimitation int) {
	if !threadMonitorEnabled {
		return
	}

	// The default value is 60 seconds and minimum value is 5 seconds
	if timeoutSec <= 5 {
		timeoutSec = 5
	}

	// The default value is 1000 threads and minimum value is 1 thread
	if maxLimitation <= 0 {
		maxLimitation = 1
	}

	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Duration(timeoutSec) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		go func() {
			mThreadNum, _ := runtime.ThreadCreateProfile(nil)
			reduce := mThreadNum - maxLimitation
			if reduce > 0 {
				wg.Add(reduce)
				for i := 0; i < reduce; i++ {
					go func() {
						runtime.LockOSThread()
						wg.Done()
					}()
				}
				wg.Wait()
			}
		}()
	}
}
