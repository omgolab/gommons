package gctest

type Test[I, O any] interface {
	Suite(name string) ToGivenStep[I, O]
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
