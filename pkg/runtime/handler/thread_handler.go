package gcthreadhandler

import (
	gcenv "github.com/omar391/go-commons/pkg/env"
	gcruntime "github.com/omar391/go-commons/pkg/runtime"
	_ "go.uber.org/automaxprocs"
)

// The function initializes certain variables with default values and then
// monitors and reduces idle OS threads based on certain conditions.
func InitWith(tOut, maxLimit int) {
	threadMonitorEnabled := gcenv.Env[bool]("THREAD_MONITOR_ENABLED", true)
	timeoutSec := gcenv.Env[int]("OS_THREADS_TIMEOUT_SEC", tOut)
	maxLimitation := gcenv.Env[int]("MAX_OS_THREADS_LIMITATION", maxLimit)

	// monitor the idle os threads and adjust reduce them accordingly
	gcruntime.MonitorAndReduceIdleOSThreads(threadMonitorEnabled, timeoutSec, maxLimitation)
}
