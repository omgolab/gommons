package gctest

type ToErrorStep[I, O any] interface {
	ReturnsError(err error) ToFinalStep[I, O]
	ReturnsNoError() ToFinalStep[I, O]
}

func (tc testCase[I, O]) ReturnsError(err error) ToFinalStep[I, O] {
	tc.name = append(tc.name, returnsErrors...)
	tc.err = err
	return &tc
}

func (tc testCase[I, O]) ReturnsNoError() ToFinalStep[I, O] {
	tc.name = append(tc.name, returnsNoErrors...)
	tc.err = nil
	return &tc
}
