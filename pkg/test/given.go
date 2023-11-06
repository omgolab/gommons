package gctest

import "time"

type ToGivenStep[I, O any] interface {
	CoreGivenStep[I, O]
	Only() CoreGivenStep[I, O]
}

type CoreGivenStep[I, O any] interface {
	Given(name string) ToWhenStep[I, O]
}

func (tc testCase[I, O]) Given(name string) ToWhenStep[I, O] {
	tc.name = append(tc.name, append(givenPrefix, []byte(name)...)...)
	return tc
}

func (tc testCase[I, O]) Only() CoreGivenStep[I, O] {
	if tc.cfg.activeTestCaseID.Load() != 0 {
		panic("only() can only be used once per test")
	}

	// create and update an atomic test case id
	tc.id = time.Now().UnixNano()
	tc.cfg.activeTestCaseID.Store(tc.id)

	return tc
}
