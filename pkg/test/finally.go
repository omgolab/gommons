package gctest

import (
	"fmt"
	"sync"
	"testing"
)

var once sync.Once

type ToFinalStep[I, O any] interface {
	Exec(execFn ExecFn[I, O]) (string, func(t *testing.T))
	ParallelExec(execFn ExecFn[I, O]) (string, func(t *testing.T))
}

func (tt *test[I, O]) Exec(execFn ExecFn[I, O]) (string, func(t *testing.T)) {
	// update wg count
	tt.rt.wg.Add(1)

	// return id and the exec fn which will track the use of "Only" in other tests
	return fmt.Sprint(tt.tc.id), func(t *testing.T) {
		defer release(&tt.rt.wg, tt.rt.ch)
		t.Helper()
		tt.tc.execFn = execFn
		tt.rt.tcs = append(tt.rt.tcs, tt.tc)
	}
}

func (tt *test[I, O]) ParallelExec(execFn ExecFn[I, O]) (string, func(t *testing.T)) {
	tt.tc.isParallel = true
	return tt.Exec(execFn)
}

func release(wg *sync.WaitGroup, ch chan bool) {
	// release the wait group
	wg.Done()

	// release the waiting channel
	once.Do(func() {
		close(ch)
	})
}
