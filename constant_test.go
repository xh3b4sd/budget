package budget

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/xh3b4sd/tracer"
)

func Test_Budget_Constant_Errors(t *testing.T) {
	var testError = &tracer.Error{
		Kind: "testError",
	}

	var executions int

	testCases := []struct {
		operation func() error
	}{
		// Case 0 tests error handling.
		{
			operation: func() error {
				executions++
				return testError
			},
		},
		// Case 1 tests error handling with masking.
		{
			operation: func() error {
				executions++
				return tracer.Mask(testError)
			},
		},
	}

	var err error

	// Note that the budget implementation is reused across all test cases in
	// order to ensure the reusability of a single budget instance. This is a
	// feature we want to ensure. Using up the configured budget of a given
	// budget instance should only happen in isolation and not affect
	// consecutive calls of the same instance.
	var b Interface
	{
		c := ConstantConfig{
			Budget:   5,
			Duration: 0,
		}

		b, err = NewConstant(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			// Note that this counter has to be reset for each test in order to
			// lead to accurate results.
			executions = 0

			err = b.Execute(tc.operation)
			if !errors.Is(err, testError) {
				t.Fatalf("expected test error")
			}

			// We expect five executions because we hard coded the errors
			// returned while the budget is five. That means we should see five
			// attempts of executing the operation.
			if executions != 5 {
				t.Fatalf("expected %#v got %#v", 5, executions)
			}
		})
	}
}

func Test_Budget_Constant_Retries(t *testing.T) {
	var testError = &tracer.Error{
		Kind: "testError",
	}

	var executions int

	testCases := []struct {
		operation  func() error
		budget     int
		executions int
	}{
		// Case 0 tests a budget of one which should result in one execution.
		{
			operation: func() error {
				executions++
				return nil
			},
			budget:     1,
			executions: 1,
		},
		// Case 1 tests a budget of two which should result in one execution.
		{
			operation: func() error {
				executions++
				return nil
			},
			budget:     2,
			executions: 1,
		},
		// Case 3 tests a budget of nine which should result in one execution.
		{
			operation: func() error {
				executions++
				return nil
			},
			budget:     9,
			executions: 1,
		},
		// Case 4 tests a budget of one which should result in one execution due
		// to the returned error.
		{
			operation: func() error {
				executions++
				return testError
			},
			budget:     1,
			executions: 1,
		},
		// Case 5 tests a budget of two which should result in two executions
		// due to the returned error.
		{
			operation: func() error {
				executions++
				return testError
			},
			budget:     2,
			executions: 2,
		},
		// Case 6 tests a budget of nine which should result in nine executions
		// due to the returned error.
		{
			operation: func() error {
				executions++
				return tracer.Mask(testError)
			},
			budget:     9,
			executions: 9,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			var b Interface
			{
				c := ConstantConfig{
					Budget:   tc.budget,
					Duration: 0,
				}

				b, err = NewConstant(c)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Note that this counter has to be reset for each test in order to
			// lead to accurate results.
			executions = 0

			err = b.Execute(tc.operation)
			if errors.Is(err, testError) {
				// Simply fall through in case we get the test error because in
				// case we get it we produced it purposefully in order to check
				// the exection results in any given situation.
			} else if err != nil {
				t.Fatal(err)
			}

			if executions != tc.executions {
				t.Fatalf("expected %#v got %#v", tc.executions, executions)
			}
		})
	}
}

func Test_Budget_Constant_Stop(t *testing.T) {
	var executions int

	testCases := []struct {
		operation  func() error
		budget     int
		executions int
	}{
		// Case 0 tests the execution stop at a budget of three which should
		// result in two executions.
		{
			operation: func() error {
				executions++
				if executions == 2 {
					return Stop{}
				}
				return fmt.Errorf("test error")
			},
			budget:     3,
			executions: 2,
		},
		// Case 1 tests the execution stop at a budget of eight which should
		// result in four executions.
		{
			operation: func() error {
				executions++
				if executions == 4 {
					return Stop{}
				}
				return fmt.Errorf("test error")
			},
			budget:     8,
			executions: 4,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			var b Interface
			{
				c := ConstantConfig{
					Budget:   tc.budget,
					Duration: 0,
				}

				b, err = NewConstant(c)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Note that this counter has to be reset for each test in order to
			// lead to accurate results.
			executions = 0

			err = b.Execute(tc.operation)
			if err != nil {
				t.Fatal(err)
			}

			if executions != tc.executions {
				t.Fatalf("expected %#v got %#v", tc.executions, executions)
			}
		})
	}
}
