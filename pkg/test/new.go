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

type testCfg[I, O any] struct {
	t                       *testing.T
	beforeAllTestsFn        StepFn
	afterAllTestsFn         StepFn
	commonBeforeEachTestsFn BeforeFn[I]
	commonAfterEachTestsFn  AfterFn[O]
	isParallel              bool
	activeTestCaseID        atomic.Int64
	wg                      sync.WaitGroup
	ch                      chan bool
}

type testCase[I, O any] struct {
	id   int64
	cfg  *testCfg[I, O]
	name []byte
	want O
	err  error
}

func NewTest[I, O any](t *testing.T, opts ...TestOption[I, O]) Test[I, O] {
	tc := &testCfg[I, O]{
		t:                t,
		activeTestCaseID: atomic.Int64{},
		ch:               make(chan bool),
	}

	// apply options
	for _, f := range opts {
		f(tc)
	}

	// execute beforeAllTestsFn if not nil
	if tc.beforeAllTestsFn != nil {
		err := tc.beforeAllTestsFn(t)
		assert.NoError(t, err)
	}

	// enable parallel exec if available
	if tc.isParallel {
		tc.t.Parallel()
	}

	// exec afterAllTestsFn
	go tc.execAfterAllTestFn()

	return tc
}

func (tc *testCfg[I, O]) execAfterAllTestFn() {
	// wait till minimum a test case's execution completes
	<-tc.ch

	// wait for the all wait group's count to be done
	tc.wg.Wait()

	// execute afterAllTestsFn if not nil
	if tc.afterAllTestsFn != nil {
		err := tc.afterAllTestsFn(tc.t)
		assert.NoError(tc.t, err)
	}
}
