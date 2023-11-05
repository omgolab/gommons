package gctest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ToFinalStep interface {
	TestCaseName() string
	TestCaseFn(execFn ExecFn[any, any]) func(t *testing.T)
}

func (tc *testCase) TestCaseName() string {
	return string(tc.name)
}

func (tc *testCase) TestCaseFn(execFn ExecFn[any, any]) func(t *testing.T) {
	return func(t *testing.T) {
		// 1. skip check
		if tc.id != tc.cfg.activeTestCaseID.Load() {
			t.Skip()
			return
		}

		// 2. exec common before each test fn
		a, err := tc.cfg.commonBeforeEachTestsFn(t)
		assert.NoError(t, err)

		// 3. exec execFn test
		b, err := execFn(t, a)
		assert.NoError(t, err)

		// 4. exec common after each test fn
		err = tc.cfg.commonAfterEachTestsFn(t, b)
		assert.NoError(t, err)
	}
}
