package gctest

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// we may create a test case like as follows:
// NewTest ->
//
//	Suite -> Given -> When -> Then -> Error -> FinalStep
var (
	scenarioPrefix  = []byte("Suite: ")
	givenPrefix     = []byte("\tGiven: ")
	whenPrefix      = []byte(" When: ")
	thenPrefix      = []byte(" Then: ")
	returnsErrors   = []byte(" And returns errors.")
	returnsNoErrors = []byte(" And returns NO errors.")
)

type BeforeFn[I any] func(t *testing.T) (I, error)
type ExecFn[I, O any] func(t *testing.T, arg I) (O, error)
type AfterFn[O any] func(t *testing.T, arg O) error
type StepFn func(t *testing.T) error

type rootTest[I, O any] struct {
	t                       *testing.T
	beforeAllTestsFn        StepFn
	afterAllTestsFn         StepFn
	commonBeforeEachTestsFn BeforeFn[I]
	commonAfterEachTestsFn  AfterFn[O]
	activeTestCaseID        atomic.Int64
	wg                      sync.WaitGroup
	ch                      chan bool
	tcs                     []testCase[I, O]
	isParallel              bool // is root test being parallel
}

type testCase[I, O any] struct {
	id         int64
	name       []byte
	want       O
	err        error
	isParallel bool // is this test case being parallel
	execFn     ExecFn[I, O]
}

type test[I, O any] struct {
	rt *rootTest[I, O]
	tc testCase[I, O]
}

func NewTest[I, O any](t *testing.T, opts ...TestOption[I, O]) Test[I, O] {
	rt := &rootTest[I, O]{
		t:                t,
		activeTestCaseID: atomic.Int64{},
		ch:               make(chan bool),
	}

	// apply options
	for _, f := range opts {
		f(rt)
	}

	// execute beforeAllTestsFn if not nil
	if rt.beforeAllTestsFn != nil {
		err := rt.beforeAllTestsFn(t)
		assert.NoError(t, err)
	}

	// enable root test parallel exec if available
	if rt.isParallel {
		rt.t.Parallel()
	}

	tt := test[I, O]{
		rt: rt,
	}

	// queue all test run
	go tt.executeAll()

	return tt
}

func (tt test[I, O]) executeAll() {
	// wait till minimum a test case's execution completes
	<-tt.rt.ch

	// wait for the all wait group's count to be done
	tt.rt.wg.Wait()

	// execute all actual tests
	for _, tc := range tt.rt.tcs {
		tt.rt.t.Run(string(tc.name), func(t *testing.T) {
			// 1. check if parallel
			if tc.isParallel {
				t.Parallel()
			}

			// 2. skip check
			if tc.id != tt.rt.activeTestCaseID.Load() {
				t.Skip()
				return
			}

			// 3. exec common before each test fn
			var err error
			var a I
			if tt.rt.commonBeforeEachTestsFn != nil {
				a, err = tt.rt.commonBeforeEachTestsFn(t)
			}
			assert.NoError(t, err)

			// 4. exec execFn test
			b, err := tc.execFn(t, a)
			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, tc.err, err)
			}

			// 5. check deep equality
			assert.Equal(t, tc.want, b)

			// 6. exec common after each test fn
			if tt.rt.commonAfterEachTestsFn != nil {
				err = tt.rt.commonAfterEachTestsFn(t, b)
				assert.NoError(t, err)
			}
		})
	}

	// execute afterAllTestsFn if not nil
	if tt.rt.afterAllTestsFn != nil {
		err := tt.rt.afterAllTestsFn(tt.rt.t)
		assert.NoError(tt.rt.t, err)
	}
}
