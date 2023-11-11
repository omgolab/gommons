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

	// exec afterAllTestsFn
	// go rt.execAfterAllTestFn()

	return test[I, O]{
		rt: rt,
	}
}
