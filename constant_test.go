package budget

import (
	"strconv"
	"testing"

	"github.com/xh3b4sd/tracer"
)

func Test_Budget_Constant_Errors(t *testing.T) {
	testCases := []struct {
		operation func() error
		matcher   func(error) bool
	}{
		// Case 0 tests error handling.
		{
			operation: func() error {
				return invalidConfigError
			},
			matcher: IsInvalidConfig,
		},
		// Case 1 tests error handling with masking.
		{
			operation: func() error {
				return tracer.Mask(invalidConfigError)
			},
			matcher: IsInvalidConfig,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			var b Interface
			{
				c := ConstantConfig{
					Budget:   1,
					Duration: 0,
				}

				b, err = NewConstant(c)
				if err != nil {
					t.Fatal(err)
				}
			}

			err = b.Execute(tc.operation)
			if !tc.matcher(err) {
				t.Fatalf("expected matcher to work")
			}
		})
	}
}

func Test_Budget_Constant_Retries(t *testing.T) {
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
		// Case 1 tests a budget of two which should result in two executions.
		{
			operation: func() error {
				executions++
				return nil
			},
			budget:     2,
			executions: 2,
		},
		// Case 2 tests a budget of three which should result in three executions.
		{
			operation: func() error {
				executions++
				return nil
			},
			budget:     3,
			executions: 3,
		},
		// Case 3 tests a budget of nine which should result in nine executions.
		{
			operation: func() error {
				executions++
				return nil
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
			if err != nil {
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
				return nil
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
				return nil
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
