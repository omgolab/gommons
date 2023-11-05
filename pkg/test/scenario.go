package gctest

type TestCase interface {
	NewScenario(name string) ToGivenStep
}

func (tc testCase) NewScenario(name string) ToGivenStep {
	tc.name = append(scenarioPrefix, []byte(name)...)
	return tc
}
