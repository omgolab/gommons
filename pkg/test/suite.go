package gctest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Test[I, O any] interface {
	Suite(name string) ToGivenStep[I, O]
	ExecuteAll() // call this via defer
}

// we are not using pointer receiver so we can copy the
// - "set name sofar" and create new duplicates

func (tt test[I, O]) Suite(name string) ToGivenStep[I, O] {
	return test[I, O]{
		tc: testCase[I, O]{
			name: append(scenarioPrefix, []byte(name)...),
		},
		rt: tt.rt,
	}
}

func (tt test[I, O]) ExecuteAll() {
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
