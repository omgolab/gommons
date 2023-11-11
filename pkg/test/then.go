package gctest

type ToThenStep[I, O any] interface {
	Then(name string, want O) ToErrorStep[I, O]
}

func (t test[I, O]) Then(name string, want O) ToErrorStep[I, O] {
	t.tc.name = append(t.tc.name, append(thenPrefix, []byte(name)...)...)
	t.tc.want = want
	return t
}
