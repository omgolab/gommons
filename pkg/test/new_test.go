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

	tt.Suite("Test Suite-Given-When-Then pattern without only").Only().
		Given("a sample given").When("a sample when").Then("a sample then").Returns(nil).Exec(func(t *testing.T, arg any) (any, error) {
		fmt.Print("hello test 1")
		return nil, nil
	})
	// t.Run("---------a", fn)

	tt.Suite("Test Suite-Given-When-Then pattern").
		Given("a sample given").When("a sample when").Then("a sample then").Returns(nil).Exec(func(t *testing.T, arg any) (any, error) {
		fmt.Print("hello test 2")
		return nil, nil
	})
	// t.Run("--------b", fn)
}
