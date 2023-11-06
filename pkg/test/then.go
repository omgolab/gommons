package gctest

type ToThenStep[I, O any] interface {
	Then(name string, want O) ToErrorStep[I, O]
}

func (tc testCase[I, O]) Then(name string, want O) ToErrorStep[I, O] {
	tc.name = append(tc.name, append(thenPrefix, []byte(name)...)...)
	tc.want = want
	return tc
}
