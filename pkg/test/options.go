package gctest

type TestOption[I, O any] func(tc *testCfg[I, O])

func WithEnvVars[I, O any](envVars map[string]string) TestOption[I, O] {
	return func(tc *testCfg[I, O]) {
		for k, v := range envVars {
			tc.t.Setenv(k, v)
		}
	}
}

func WithParallel[I, O any]() TestOption[I, O] {
	return func(tc *testCfg[I, O]) {
		tc.isParallel = true
	}
}

func WithBeforeAllTestsFn[I, O any](fn StepFn) TestOption[I, O] {
	return func(tc *testCfg[I, O]) {
		tc.beforeAllTestsFn = fn
	}
}

func WithAfterAllTestsFn[I, O any](fn StepFn) TestOption[I, O] {
	return func(tc *testCfg[I, O]) {
		tc.afterAllTestsFn = fn
	}
}

func WithCommonBeforeEachTestsFn[I, O any](fn BeforeFn[I]) TestOption[I, O] {
	return func(tc *testCfg[I, O]) {
		tc.commonBeforeEachTestsFn = fn
	}
}

func WithCommonAfterEachTestsFn[I, O any](fn AfterFn[O]) TestOption[I, O] {
	return func(tc *testCfg[I, O]) {
		tc.commonAfterEachTestsFn = fn
	}
}
