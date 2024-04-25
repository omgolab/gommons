package gctest

type ToErrorStep[I, O any] interface {
	ReturnsError(err error) ToFinalStep[I, O]
	Returns(want O) ToFinalStep[I, O]
}

func (t test[I, O]) ReturnsError(err error) ToFinalStep[I, O] {
	t.tc.name = append(t.tc.name, returnsErrors...)
	t.tc.err = err
	return &t
}

func (t test[I, O]) Returns(want O) ToFinalStep[I, O] {
	t.tc.name = append(t.tc.name, returnsNoErrors...)
	t.tc.err = nil
	t.tc.want = want
	return &t
}
