package gctest

type TestOption func(tc *testCfg)

func WithEnvVars(envVars map[string]string) TestOption {
	return func(tc *testCfg) {
		for k, v := range envVars {
			tc.t.Setenv(k, v)
		}
	}
}

func WithParallel() TestOption {
	return func(tc *testCfg) {
		tc.isParallel = true
	}
}

func WithBeforeAllTestsFn(fn BeforeFn[any]) TestOption {
	return func(tc *testCfg) {
		tc.beforeAllTestsFn = fn
	}
}

func WithAfterAllTestsFn(fn AfterFn[any]) TestOption {
	return func(tc *testCfg) {
		tc.afterAllTestsFn = fn
	}
}

func WithCommonBeforeEachTestsFn(fn BeforeFn[any]) TestOption {
	return func(tc *testCfg) {
		tc.commonBeforeEachTestsFn = fn
	}
}

func WithCommonAfterEachTestsFn(fn AfterFn[any]) TestOption {
	return func(tc *testCfg) {
		tc.commonAfterEachTestsFn = fn
	}
}
