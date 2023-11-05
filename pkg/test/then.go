package gctest

type ToThenStep interface {
	Then(name string) ToErrorStep
}

func (tc *testCase) Then(name string) ToErrorStep {
	tc.name = append(tc.name, append(thenPrefix, []byte(name)...)...)
	return tc
}
