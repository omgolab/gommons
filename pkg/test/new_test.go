package gctest_test

import (
	"fmt"
	"testing"

	gctest "github.com/omar391/go-commons/pkg/test"
	"github.com/tj/assert"
)

// Returns a TestCase object
func TestReturnsATestCaseObject(t *testing.T) {
	testCase := gctest.NewTest[any, any](t)
	assert.NotNil(t, testCase)
}

// Executes beforeAllTestsFn if not nil
func TestExecutesBeforeAllTestsFnIfNotNil(t *testing.T) {
	var executed bool
	beforeAllTestsFn := func(t *testing.T) error {
		executed = true
		return nil
	}
	_ = gctest.NewTest[any, any](t, gctest.WithBeforeAllTestsFn[any, any](beforeAllTestsFn))
	assert.True(t, executed)
}

func TestBehavioralTestPattern(t *testing.T) {

	tt := gctest.NewTest[any, any](t)

	tc1 := tt.Suite("Test Suite-Given-When-Then pattern without only").
		Given("a sample given").When("a sample when").Then("a sample then", nil).ReturnsNoError()

	t.Run("a", tc1.TestCaseFn(func(t *testing.T, _ any) (any, error) {
		t.Parallel()
		fmt.Print("hello test")
		return nil, nil
	}))

	tc2 := tt.Suite("Test Suite-Given-When-Then pattern").Only().
	Given("a sample given").When("a sample when").Then("a sample then", nil).ReturnsNoError()
	t.Run("a", tc2.TestCaseFn(func(t *testing.T, _ any) (any, error) {
		fmt.Print("--")
		return nil, nil
	}))
}
