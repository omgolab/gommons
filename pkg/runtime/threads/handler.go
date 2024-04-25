package gcthreads

import (
	env "github.com/omgolab/go-commons/pkg/env"
	_ "go.uber.org/automaxprocs"
)

// The function check and initializes environment variables with default values and then
// monitors and reduces idle OS threads based on those variables.
func MonitorWith(intervalSec int) {
	threadMonitorEnabled := env.Env[bool]("THREAD_MONITOR_ENABLED", true)
	timeoutSec := env.Env[int]("OS_THREADS_TIMEOUT_SEC", intervalSec)
	rateLimit := env.Env[int]("OS_THREADS_REDUCTION_RATE", 1)

	// monitor the idle os threads and adjust reduce them accordingly
	if threadMonitorEnabled {
		MonitorAndReduceIdleOSThreads(timeoutSec, rateLimit)
	}
}
