package gctest

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var once sync.Once

type ToFinalStep[I, O any] interface {
	TestCaseName() string
	TestCaseFn(execFn ExecFn[I, O]) func(t *testing.T)
}

func (tc *testCase[I, O]) TestCaseName() string {
	return string(tc.name)
}

func (tc *testCase[I, O]) TestCaseFn(execFn ExecFn[I, O]) func(t *testing.T) {
	// update wg count
	tc.cfg.wg.Add(1)

	// return the exec fn
	return func(t *testing.T) {
		defer release(&tc.cfg.wg, tc.cfg.ch)

		// 1. skip check
		if tc.id != tc.cfg.activeTestCaseID.Load() {
			t.Skip()
			return
		}

		// 2. exec common before each test fn
		var err error
		var a I
		if tc.cfg.commonBeforeEachTestsFn != nil {
			a, err = tc.cfg.commonBeforeEachTestsFn(t)
		}
		assert.NoError(t, err)

		// 3. exec execFn test
		b, err := execFn(t, a)
		if tc.err == nil {
			assert.NoError(t, err)
		} else {
			assert.ErrorIs(t, tc.err, err)
		}

		// 4. check deep equality
		assert.Equal(t, tc.want, b)

		// 5. exec common after each test fn
		if tc.cfg.commonAfterEachTestsFn != nil {
			err = tc.cfg.commonAfterEachTestsFn(t, b)
			assert.NoError(t, err)
		}
	}
}

func release(wg *sync.WaitGroup, ch chan bool) {
	// release the wait group
	wg.Done()

	// release the waiting channel
	once.Do(func() {
		close(ch)
	})
}
