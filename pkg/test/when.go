package gctest

type ToWhenStep[I, O any] interface {
	When(name string) ToThenStep[I, O]
}

func (tc testCase[I, O]) When(name string) ToThenStep[I, O] {
	tc.name = append(tc.name, append(whenPrefix, []byte(name)...)...)
	return tc
}
