package gctest

import (
	"time"
)

type ToGivenStep interface {
	CoreGivenStep
	Only() CoreGivenStep
}

type CoreGivenStep interface {
	Given(name string) ToWhenStep
}

func (tc testCase) Given(name string) ToWhenStep {
	tc.name = append(tc.name, append(givenPrefix, []byte(name)...)...)
	return tc
}

func (tc testCase) Only() CoreGivenStep {
	if tc.cfg.activeTestCaseID.Load() != 0 {
		panic("only() can only be used once per test")
	}

	// create and update an atomic test case id
	tc.id = time.Now().UnixNano()
	tc.cfg.activeTestCaseID.Store(tc.id)

	return tc
}
