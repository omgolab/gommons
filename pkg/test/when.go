package gctest

type ToWhenStep[I, O any] interface {
	When(name string) ToThenStep[I, O]
}

func (t test[I, O]) When(name string) ToThenStep[I, O] {
	t.tc.name = append(t.tc.name, append(whenPrefix, []byte(name)...)...)
	return t
}
