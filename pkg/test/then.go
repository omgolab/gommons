package gctest

type ToThenStep[I, O any] interface {
	Then(name string) ToErrorStep[I, O]
}

func (t test[I, O]) Then(name string) ToErrorStep[I, O] {
	t.tc.name = append(t.tc.name, append(thenPrefix, []byte(name)...)...)
	return t
}
