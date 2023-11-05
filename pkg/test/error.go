package gctest

type ToErrorStep interface {
	ReturnsError() ToFinalStep
	ReturnsNoError() ToFinalStep
}

func (tc *testCase) ReturnsError() ToFinalStep {
	tc.name = append(tc.name, returnsErrors...)
	return tc
}

func (tc *testCase) ReturnsNoError() ToFinalStep {
	tc.name = append(tc.name, returnsNoErrors...)
	return tc
}
