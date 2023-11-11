package gctest

type ToErrorStep[I, O any] interface {
	ReturnsError(err error) ToFinalStep[I, O]
	ReturnsNoError() ToFinalStep[I, O]
}

func (t test[I, O]) ReturnsError(err error) ToFinalStep[I, O] {
	t.tc.name = append(t.tc.name, returnsErrors...)
	t.tc.err = err
	return &t
}

func (t test[I, O]) ReturnsNoError() ToFinalStep[I, O] {
	t.tc.name = append(t.tc.name, returnsNoErrors...)
	t.tc.err = nil
	return &t
}
