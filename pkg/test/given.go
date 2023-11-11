package gctest

import "time"

type ToGivenStep[I, O any] interface {
	CoreGivenStep[I, O]
	Only() CoreGivenStep[I, O]
}

type CoreGivenStep[I, O any] interface {
	Given(name string) ToWhenStep[I, O]
}

func (t test[I, O]) Given(name string) ToWhenStep[I, O] {
	t.tc.name = append(t.tc.name, append(givenPrefix, []byte(name)...)...)

	//assign an ID if empty
	if t.tc.id == 0 {
		t.tc.id = time.Now().UnixNano()
	}

	return t
}

func (t test[I, O]) Only() CoreGivenStep[I, O] {
	if t.rt.activeTestCaseID.Load() != 0 {
		panic("only() can only be used once per test")
	}

	// create and update an atomic test case id
	t.tc.id = time.Now().UnixNano()
	t.rt.activeTestCaseID.Store(t.tc.id)

	return t
}
