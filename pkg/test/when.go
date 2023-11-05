package gctest

type ToWhenStep interface {
	When(name string) ToThenStep
}

func (tc testCase) When(name string) ToThenStep {
	tc.name = append(tc.name, append(whenPrefix, []byte(name)...)...)
	return &tc
}
