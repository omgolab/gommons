package gctest

type Test[I, O any] interface {
	Suite(name string) ToGivenStep[I, O]
}

// we are not using pointer receiver so we can copy the
// - "set name sofar" and create new duplicates

func (t *testCfg[I, O]) Suite(name string) ToGivenStep[I, O] {
	return testCase[I, O]{
		cfg:  t,
		name: append(scenarioPrefix, []byte(name)...),
	}
}
