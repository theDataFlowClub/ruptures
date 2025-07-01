package stat_test

import (
	"math"
	"testing"

	"github.com/theDataFlowClub/ruptures/core/stat" // The package being tested
)

// Define a small tolerance for float comparisons.
const floatTolerance = 1e-9

// --- Test functions for Mean ---

func TestMean(t *testing.T) {
	testCases := []struct {
		name         string
		data         []float64
		expectedMean float64
		expectError  bool
	}{
		{
			name:         "PositiveNumbers",
			data:         []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			expectedMean: 3.0,
			expectError:  false,
		},
		{
			name:         "NegativeNumbers",
			data:         []float64{-1.0, -2.0, -3.0},
			expectedMean: -2.0,
			expectError:  false,
		},
		{
			name:         "MixedNumbers",
			data:         []float64{-1.0, 0.0, 1.0, 2.0},
			expectedMean: 0.5,
			expectError:  false,
		},
		{
			name:         "SingleElement",
			data:         []float64{7.0},
			expectedMean: 7.0,
			expectError:  false,
		},
		{
			name:         "EmptySlice",
			data:         []float64{},
			expectedMean: 0.0, // Error case, but return value is 0.0
			expectError:  true,
		},
		{
			name:         "Zeroes",
			data:         []float64{0.0, 0.0, 0.0},
			expectedMean: 0.0,
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := stat.Mean(tc.data)

			if tc.expectError {
				if err == nil {
					t.Errorf("Mean() expected an error, but got nil")
				}
				// We're expecting a specific error message, checking for it.
				expectedErrStr := "empty slice for mean calculation"
				if err != nil && err.Error() != expectedErrStr {
					t.Errorf("Mean() got unexpected error: %v, want %s", err, expectedErrStr)
				}
			} else {
				if err != nil {
					t.Errorf("Mean() got unexpected error: %v, want nil", err)
				}
				if math.Abs(result-tc.expectedMean) > floatTolerance {
					t.Errorf("Mean() = %f; want %f (diff: %f)", result, tc.expectedMean, math.Abs(result-tc.expectedMean))
				}
			}
		})
	}
}

// --- Test functions for Median ---

func TestMedian(t *testing.T) {
	testCases := []struct {
		name           string
		data           []float64
		expectedMedian float64
		expectError    bool
	}{
		{
			name:           "OddNumberOfElements",
			data:           []float64{1.0, 3.0, 2.0},
			expectedMedian: 2.0,
			expectError:    false,
		},
		{
			name:           "EvenNumberOfElements",
			data:           []float64{1.0, 2.0, 3.0, 4.0},
			expectedMedian: 2.5,
			expectError:    false,
		},
		{
			name:           "NegativeNumbers",
			data:           []float64{-5.0, -1.0, -3.0},
			expectedMedian: -3.0,
			expectError:    false,
		},
		{
			name:           "MixedNumbers",
			data:           []float64{-10.0, 0.0, 5.0, -5.0},
			expectedMedian: -2.5,
			expectError:    false,
		},
		{
			name:           "SingleElement",
			data:           []float64{42.0},
			expectedMedian: 42.0,
			expectError:    false,
		},
		{
			name:           "EmptySlice",
			data:           []float64{},
			expectedMedian: 0.0, // Error case, but return value is 0.0
			expectError:    true,
		},
		{
			name:           "DuplicateValues",
			data:           []float64{1.0, 2.0, 2.0, 3.0},
			expectedMedian: 2.0,
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := stat.Median(tc.data)

			if tc.expectError {
				if err == nil {
					t.Errorf("Median() expected an error, but got nil")
				}
				expectedErrStr := "empty slice for median calculation"
				if err != nil && err.Error() != expectedErrStr {
					t.Errorf("Median() got unexpected error: %v, want %s", err, expectedErrStr)
				}
			} else {
				if err != nil {
					t.Errorf("Median() got unexpected error: %v, want nil", err)
				}
				if math.Abs(result-tc.expectedMedian) > floatTolerance {
					t.Errorf("Median() = %f; want %f (diff: %f)", result, tc.expectedMedian, math.Abs(result-tc.expectedMedian))
				}
			}
		})
	}
}

// --- Test functions for Variance ---

func TestVariance(t *testing.T) {
	testCases := []struct {
		name             string
		data             []float64
		expectedVariance float64
		expectError      bool
	}{
		{
			name:             "SmallPositiveNumbers",
			data:             []float64{1.0, 2.0, 3.0, 4.0, 5.0}, // Mean=3, SumSqDiff=10, N=5, Var=10/5=2
			expectedVariance: 2.0,
			expectError:      false,
		},
		{
			name:             "NegativeNumbers",
			data:             []float64{-1.0, -2.0, -3.0}, // Mean=-2, SumSqDiff=2, N=3, Var=2/3
			expectedVariance: 0.6666666666666666,
			expectError:      false,
		},
		{
			name:             "ZeroVariance",
			data:             []float64{5.0, 5.0, 5.0}, // Mean=5, SumSqDiff=0, N=3, Var=0
			expectedVariance: 0.0,
			expectError:      false,
		},
		{
			name:             "SingleElement",
			data:             []float64{10.0}, // Mean=10, SumSqDiff=0, N=1, Var=0
			expectedVariance: 0.0,
			expectError:      false,
		},
		{
			name:             "EmptySlice",
			data:             []float64{},
			expectedVariance: 0.0,
			expectError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := stat.Variance(tc.data)

			if tc.expectError {
				if err == nil {
					t.Errorf("Variance() expected an error, but got nil")
				}
				expectedErrStr := "not enough points for variance calculation"
				if err != nil && err.Error() != expectedErrStr {
					t.Errorf("Variance() got unexpected error: %v, want %s", err, expectedErrStr)
				}
			} else {
				if err != nil {
					t.Errorf("Variance() got unexpected error: %v, want nil", err)
				}
				if math.Abs(result-tc.expectedVariance) > floatTolerance {
					t.Errorf("Variance() = %f; want %f (diff: %f)", result, tc.expectedVariance, math.Abs(result-tc.expectedVariance))
				}
			}
		})
	}
}
