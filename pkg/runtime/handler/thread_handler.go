package gcthreadhandler

import (
	gcenv "github.com/omar391/go-commons/pkg/env"
	gcruntime "github.com/omar391/go-commons/pkg/runtime"
	_ "go.uber.org/automaxprocs"
)

// The function initializes certain variables with default values and then
// monitors and reduces idle OS threads based on certain conditions.
func MonitorWith(intervalSec int) {
	threadMonitorEnabled := gcenv.Env[bool]("THREAD_MONITOR_ENABLED", true)
	timeoutSec := gcenv.Env[int]("OS_THREADS_TIMEOUT_SEC", intervalSec)
	rateLimit := gcenv.Env[int]("OS_THREADS_REDUCTION_RATE", 1)

	// monitor the idle os threads and adjust reduce them accordingly
	gcruntime.MonitorAndReduceIdleOSThreads(threadMonitorEnabled, timeoutSec, rateLimit)
}
