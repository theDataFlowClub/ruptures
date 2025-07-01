package base_test

import (
	"errors" // Import for checking error types
	"testing"

	"github.com/theDataFlowClub/ruptures/core/base"
	"github.com/theDataFlowClub/ruptures/core/exceptions" // Import custom exceptions
	"github.com/theDataFlowClub/ruptures/core/types"
)

// MockCostFunction is a simple mock implementation of the base.CostFunction interface
// for testing purposes. It returns a predefined cost and can simulate an error.
type MockCostFunction struct {
	costPerSegment float64
	simulateError  bool // New field to simulate an error return
}

func (m *MockCostFunction) Fit(signal types.Matrix) error {
	// No-op for this mock
	return nil
}

// Error now returns a float64 and an error, matching the updated interface.
func (m *MockCostFunction) Error(start, end int) (float64, error) {
	if m.simulateError {
		return 0.0, exceptions.ErrNotEnoughPoints // Simulate an error
	}
	return m.costPerSegment, nil // Return fixed cost and no error
}

func (m *MockCostFunction) Model() string {
	return "mock_cost"
}

func TestSumOfCosts(t *testing.T) {
	testCases := []struct {
		name          string
		bkps          []int
		mockCostValue float64
		simulateErr   bool
		expectedSum   float64
		expectError   bool
		expectedErr   error
	}{
		{
			name:          "NoBreakpoints_EmptySlice",
			bkps:          []int{},
			mockCostValue: 10.0,
			simulateErr:   false,
			expectedSum:   0.0,
			expectError:   false,
		},
		{
			name:          "OneSegment",
			bkps:          []int{100},
			mockCostValue: 10.0,
			simulateErr:   false,
			expectedSum:   10.0, // 1 segment * 10.0 cost/segment
			expectError:   false,
		},
		{
			name:          "TwoSegments",
			bkps:          []int{50, 100},
			mockCostValue: 10.0,
			simulateErr:   false,
			expectedSum:   20.0, // 2 segments * 10.0 cost/segment
			expectError:   false,
		},
		{
			name:          "MultipleSegments",
			bkps:          []int{25, 50, 75, 100},
			mockCostValue: 10.0,
			simulateErr:   false,
			expectedSum:   40.0, // 4 segments * 10.0 cost/segment
			expectError:   false,
		},
		{
			name:          "PropagateError_FirstSegment",
			bkps:          []int{50, 100},
			mockCostValue: 0.0, // Cost value doesn't matter if error is simulated
			simulateErr:   true,
			expectedSum:   0.0, // Expected sum when an error occurs
			expectError:   true,
			expectedErr:   exceptions.ErrNotEnoughPoints,
		},
		{
			name:          "PropagateError_LaterSegment",
			bkps:          []int{25, 50, 75, 100},
			mockCostValue: 0.0,
			simulateErr:   true, // If one segment fails, the whole sum fails
			expectedSum:   0.0,
			expectError:   true,
			expectedErr:   exceptions.ErrNotEnoughPoints,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock cost function for the current test case
			mockCost := &MockCostFunction{
				costPerSegment: tc.mockCostValue,
				simulateError:  tc.simulateErr,
			}

			resultSum, err := base.SumOfCosts(mockCost, tc.bkps)

			if tc.expectError {
				if err == nil {
					t.Errorf("SumOfCosts() expected an error, but got nil")
				} else if !errors.Is(err, tc.expectedErr) {
					t.Errorf("SumOfCosts() got unexpected error: %v, want: %v", err, tc.expectedErr)
				}
			} else {
				if err != nil {
					t.Errorf("SumOfCosts() got unexpected error: %v, want nil", err)
				}
				if resultSum != tc.expectedSum {
					t.Errorf("SumOfCosts() = %f; want %f", resultSum, tc.expectedSum)
				}
			}
		})
	}
}
