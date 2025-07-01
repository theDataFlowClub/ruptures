package cost_test // _test suffix for package name when testing external functions

import (
	"errors"
	"math"
	"testing"

	// For the CostFunction interface
	"github.com/theDataFlowClub/ruptures/core/cost"       // The package being tested
	"github.com/theDataFlowClub/ruptures/core/exceptions" // For custom error types
	"github.com/theDataFlowClub/ruptures/core/types"      // For Matrix type
)

// Helper function to create a matrix (signal) for tests
func createMatrix(data [][]float64) types.Matrix {
	mat := make(types.Matrix, len(data))
	for i, row := range data {
		mat[i] = make([]float64, len(row))
		copy(mat[i], row)
	}
	return mat
}

func TestCostL2_Fit(t *testing.T) {
	testCases := []struct {
		name        string
		signal      types.Matrix
		expectError bool
	}{
		{
			name:        "ValidSignal_Univariate",
			signal:      createMatrix([][]float64{{1.0}, {2.0}, {3.0}}),
			expectError: false,
		},
		{
			name:        "ValidSignal_Multivariate",
			signal:      createMatrix([][]float64{{1.0, 10.0}, {2.0, 20.0}, {3.0, 30.0}}),
			expectError: false,
		},
		{
			name:        "EmptySignal",
			signal:      createMatrix([][]float64{}), // Empty outer slice
			expectError: true,                        // Should return ErrNotEnoughPoints from Fit
		},
		{
			name:        "NilSignal",
			signal:      nil,
			expectError: true, // Should return ErrNotEnoughPoints from Fit
		},
		{
			name:        "SignalWithEmptyRows",
			signal:      createMatrix([][]float64{{}, {}}), // Signal with 0 features per sample
			expectError: true,                              // Should return ErrNotEnoughPoints from Fit due to len(signal[0]) == 0 check
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l2Cost := cost.NewCostL2()
			err := l2Cost.Fit(tc.signal)

			if tc.expectError {
				if err == nil {
					t.Errorf("Fit() expected an error for signal %v, but got nil", tc.signal)
				} else if !errors.Is(err, exceptions.ErrNotEnoughPoints) {
					t.Errorf("Fit() got unexpected error: %v, want %v", err, exceptions.ErrNotEnoughPoints)
				}
			} else {
				if err != nil {
					t.Errorf("Fit() got unexpected error: %v, want nil", err)
				}
				if l2Cost.Signal == nil && len(tc.signal) > 0 {
					t.Errorf("Fit() failed to store signal: %v", l2Cost.Signal)
				}
			}
		})
	}
}

func TestCostL2_Error(t *testing.T) {
	// Sample signals for testing CostL2_Error
	univariateSignal := createMatrix([][]float64{{1.0}, {2.0}, {3.0}, {4.0}, {5.0}}) // Mean=3, Var=2 (for segment 0-5)
	multivariateSignal := createMatrix([][]float64{
		{1.0, 10.0}, {2.0, 20.0}, {3.0, 30.0}, {4.0, 40.0}, {5.0, 50.0},
	})

	testCases := []struct {
		name          string
		signal        types.Matrix
		minSize       int
		start         int
		end           int
		expectedCost  float64
		expectError   bool
		expectedError error
	}{
		// Valid Segments
		{
			name:         "ValidSegment_Univariate_Full",
			signal:       univariateSignal,
			minSize:      1,
			start:        0,
			end:          5,    // Segment: [1,2,3,4,5]
			expectedCost: 10.0, // (end-start) * var = 5 * 2 = 10
			expectError:  false,
		},
		{
			name:         "ValidSegment_Univariate_Partial",
			signal:       univariateSignal,
			minSize:      1,
			start:        0,
			end:          3,   // Segment: [1,2,3], Mean=2, Var=(1+0+1)/3 = 2/3
			expectedCost: 2.0, // (end-start) * var = 3 * 2/3 = 2
			expectError:  false,
		},
		{
			name:         "ValidSegment_Multivariate_Full",
			signal:       multivariateSignal,
			minSize:      1,
			start:        0,
			end:          5,
			expectedCost: 1010.0, // <-- CORREGIDO: 5 * (2 + 200) = 1010
			expectError:  false,
		},
		{
			name:         "ValidSegment_Multivariate_Partial",
			signal:       multivariateSignal,
			minSize:      1,
			start:        0,
			end:          3,     // Segments: [1,2,3] and [10,20,30]. Var each 2/3 and 200/3. Sum Var=202/3
			expectedCost: 202.0, // <-- CORREGIDO: 3 * (202/3) = 202
			expectError:  false,
		},
		{
			name:         "ValidSegment_SinglePoint",
			signal:       univariateSignal,
			minSize:      1,
			start:        2,
			end:          3,   // Segment: [3.0] -> Variance is 0
			expectedCost: 0.0, // (end-start) * 0 = 0
			expectError:  false,
		},

		// Invalid Segments due to minSize
		{
			name:          "InvalidSegment_TooShort_minSize2",
			signal:        univariateSignal,
			minSize:       2,
			start:         0,
			end:           1,   // Length 1, minSize 2
			expectedCost:  0.0, // Return zero on error
			expectError:   true,
			expectedError: exceptions.ErrNotEnoughPoints,
		},
		{
			name:          "InvalidSegment_TooShort_minSize3",
			signal:        univariateSignal,
			minSize:       3,
			start:         0,
			end:           2, // Length 2, minSize 3
			expectedCost:  0.0,
			expectError:   true,
			expectedError: exceptions.ErrNotEnoughPoints,
		},

		// Invalid Segments due to bounds (using the new specific error)
		{
			name:          "InvalidSegment_StartNegative",
			signal:        univariateSignal,
			minSize:       1,
			start:         -1,
			end:           2,
			expectedCost:  0.0,
			expectError:   true,
			expectedError: exceptions.ErrSegmentOutOfBounds, // <-- CORREGIDO
		},
		{
			name:          "InvalidSegment_EndTooLarge",
			signal:        univariateSignal,
			minSize:       1,
			start:         0,
			end:           len(univariateSignal) + 1, // End beyond signal length
			expectedCost:  0.0,
			expectError:   true,
			expectedError: exceptions.ErrSegmentOutOfBounds, // <-- CORREGIDO
		},
		{
			name:          "InvalidSegment_StartEqualsEnd",
			signal:        univariateSignal,
			minSize:       1,
			start:         2,
			end:           2, // Length 0
			expectedCost:  0.0,
			expectError:   true,
			expectedError: exceptions.ErrSegmentOutOfBounds, // <-- CORREGIDO (start >= end)
		},
		{
			name:          "InvalidSegment_StartGreaterThanEnd",
			signal:        univariateSignal,
			minSize:       1,
			start:         3,
			end:           2, // Negative length
			expectedCost:  0.0,
			expectError:   true,
			expectedError: exceptions.ErrSegmentOutOfBounds, // <-- CORREGIDO (start >= end)
		},
	}

	// Use a small epsilon for float comparisons due to potential precision issues.
	const floatTolerance = 1e-6

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l2Cost := cost.NewCostL2()
			// Fit the cost function with the signal
			err := l2Cost.Fit(tc.signal)
			if err != nil && !tc.expectError { // Only error if Fit should not fail
				t.Fatalf("Fit() failed unexpectedly: %v", err)
			}
			if err != nil && tc.expectError && errors.Is(err, exceptions.ErrNotEnoughPoints) {
				// If Fit() already errored with NotEnoughPoints and we expected it, skip Error() test
				return
			}

			l2Cost.MinSize = tc.minSize // AHORA CORRECTO

			resultCost, err := l2Cost.Error(tc.start, tc.end)

			if tc.expectError {
				if err == nil {
					t.Errorf("Error() expected an error, but got nil")
				} else if tc.expectedError != nil && !errors.Is(err, tc.expectedError) {
					t.Errorf("Error() got unexpected error: %v, want %v", err, tc.expectedError)
				}
				// Si expectedError es nil, cualquier error es válido (no debería ocurrir aquí con expectError=true)
			} else {
				if err != nil {
					t.Errorf("Error() got unexpected error: %v, want nil", err)
				}
				if math.Abs(resultCost-tc.expectedCost) > floatTolerance {
					t.Errorf("Error() = %f; want %f (diff: %f)", resultCost, tc.expectedCost, math.Abs(resultCost-tc.expectedCost))
				}
			}
		})
	}
}

func TestCostL2_Model(t *testing.T) {
	l2Cost := cost.NewCostL2()
	expectedModel := "l2"
	if l2Cost.Model() != expectedModel {
		t.Errorf("Model() = %s; want %s", l2Cost.Model(), expectedModel)
	}
}

// Define a small tolerance for float comparisons.
//const floatTolerance = 1e-9

// Helper function to create a matrix (signal) for tests

// CREADA AL INICIO PARA L2 TAMBIEN

func TestCostL1_Error(t *testing.T) {
	testCases := []struct {
		name          string
		signal        [][]float64 // Input signal for CostL1
		minSize       int         // minSize for the CostL1 instance
		start         int         // Start index of the segment to test (NEW FIELD)
		end           int         // End index of the segment to test (NEW FIELD)
		expectedCost  float64
		expectError   bool
		expectedError error // To check for specific error types
	}{
		{
			name:         "SimpleDataOddLength",
			signal:       [][]float64{{1.0}, {2.0}, {3.0}, {4.0}, {5.0}},
			minSize:      2,
			start:        0,
			end:          5, // Full segment
			expectedCost: 6.0,
			expectError:  false,
		},
		{
			name:         "SimpleDataEvenLength",
			signal:       [][]float64{{1.0}, {2.0}, {3.0}, {4.0}},
			minSize:      2,
			start:        0,
			end:          4, // Full segment
			expectedCost: 4.0,
			expectError:  false,
		},
		{
			name:         "NegativeNumbers",
			signal:       [][]float64{{-1.0}, {-2.0}, {-3.0}},
			minSize:      2,
			start:        0,
			end:          3, // Full segment
			expectedCost: 2.0,
			expectError:  false,
		},
		{
			name:         "MixedNumbers",
			signal:       [][]float64{{-10.0}, {0.0}, {5.0}, {-5.0}},
			minSize:      2,
			start:        0,
			end:          4, // Full segment
			expectedCost: 20.0,
			expectError:  false,
		},
		{
			name:         "SingleElement_CostZero", // Only works if minSize <= 1
			signal:       [][]float64{{7.0}},
			minSize:      1,
			start:        0,
			end:          1, // Full segment
			expectedCost: 0.0,
			expectError:  false,
		},
		{
			name:          "SingleElement_NotEnoughPoints", // Default L1 min_size is 2
			signal:        [][]float64{{7.0}},
			minSize:       2,
			start:         0,
			end:           1,
			expectedCost:  0.0,
			expectError:   true,
			expectedError: exceptions.ErrNotEnoughPoints,
		},
		{
			name:          "EmptySignal_FitError",
			signal:        [][]float64{}, // Empty outer slice
			minSize:       2,
			start:         0,
			end:           0,
			expectedCost:  0.0,
			expectError:   true,
			expectedError: exceptions.ErrNotEnoughPoints, // Fit should return this error
		},
		{
			name:          "EmptyInnerSlice_FitError",
			signal:        [][]float64{{}}, // Empty inner slice (feature)
			minSize:       2,
			start:         0,
			end:           1, // A segment of length 1, with 0 features
			expectedCost:  0.0,
			expectError:   true,
			expectedError: exceptions.ErrNotEnoughPoints, // Fit should return this error
		},
		{
			name: "Multivariate_TwoFeatures",
			signal: [][]float64{
				{1.0, 10.0},
				{2.0, 20.0},
				{3.0, 30.0},
			},
			minSize:      2,
			start:        0,
			end:          3,          // Full segment
			expectedCost: 2.0 + 20.0, // (1,2,3) -> Med=2, SumAbsDev=2; (10,20,30) -> Med=20, SumAbsDev=20. Total = 22.0
			expectError:  false,
		},
		{
			name:         "Zeroes_Multivariate",
			signal:       [][]float64{{0.0, 0.0}, {0.0, 0.0}},
			minSize:      2,
			start:        0,
			end:          2, // Full segment
			expectedCost: 0.0,
			expectError:  false,
		},
		// Invalid Segments due to bounds
		{
			name:          "InvalidSegment_StartNegative",
			signal:        [][]float64{{1.0}, {2.0}, {3.0}},
			minSize:       2,
			start:         -1, // This is the segment start
			end:           2,  // This is the segment end
			expectedCost:  0.0,
			expectError:   true,
			expectedError: exceptions.ErrSegmentOutOfBounds,
		},
		{
			name:          "InvalidSegment_EndTooLarge",
			signal:        [][]float64{{1.0}, {2.0}, {3.0}},
			minSize:       2,
			start:         0,
			end:           4, // End beyond signal length
			expectedCost:  0.0,
			expectError:   true,
			expectedError: exceptions.ErrSegmentOutOfBounds,
		},
		{
			name:          "InvalidSegment_StartEqualsEnd",
			signal:        [][]float64{{1.0}, {2.0}, {3.0}},
			minSize:       2,
			start:         2,
			end:           2, // Length 0
			expectedCost:  0.0,
			expectError:   true,
			expectedError: exceptions.ErrSegmentOutOfBounds,
		},
		{
			name:          "InvalidSegment_StartGreaterThanEnd",
			signal:        [][]float64{{1.0}, {2.0}, {3.0}},
			minSize:       2,
			start:         3,
			end:           2, // Negative length
			expectedCost:  0.0,
			expectError:   true,
			expectedError: exceptions.ErrSegmentOutOfBounds,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l1Cost := cost.NewCostL1()
			l1Cost.MinSize = tc.minSize // Set minSize for the test case

			// First, fit the signal. This is crucial.
			fitErr := l1Cost.Fit(createMatrix(tc.signal))

			// Check for errors during Fit, especially for empty/invalid signals
			if tc.name == "EmptySignal_FitError" || tc.name == "EmptyInnerSlice_FitError" {
				if fitErr == nil {
					t.Fatalf("CostL1.Fit() expected an error for '%s', but got nil", tc.name)
				}
				if tc.expectedError != nil && !errors.Is(fitErr, tc.expectedError) {
					t.Fatalf("CostL1.Fit() got unexpected error for '%s': %v, want %v", tc.name, fitErr, tc.expectedError)
				}
				// If Fit errored as expected, we don't proceed to call Error method
				return
			} else if fitErr != nil {
				// For other cases, if Fit unexpectedly failed, it's a fatal error
				t.Fatalf("CostL1.Fit() unexpectedly failed for '%s': %v", tc.name, fitErr)
			}

			// Now, calculate the error for the segment using the specified start and end.
			resultCost, err := l1Cost.Error(tc.start, tc.end) // Use tc.start and tc.end here!

			if tc.expectError {
				if err == nil {
					t.Errorf("CostL1.Error() for '%s' expected an error, but got nil", tc.name)
				} else if tc.expectedError != nil && !errors.Is(err, tc.expectedError) {
					t.Errorf("CostL1.Error() for '%s' got unexpected error: %v, want %v", tc.name, err, tc.expectedError)
				}
			} else {
				if err != nil {
					t.Errorf("CostL1.Error() for '%s' got unexpected error: %v, want nil", tc.name, err)
				}
				if math.Abs(resultCost-tc.expectedCost) > floatTolerance {
					t.Errorf("CostL1.Error() for '%s' = %f; want %f (diff: %f)", tc.name, resultCost, tc.expectedCost, math.Abs(resultCost-tc.expectedCost))
				}
			}
		})
	}
}

func TestCostL1_Model(t *testing.T) {
	l1Cost := cost.NewCostL1()
	expectedModel := "l1"
	if l1Cost.Model() != expectedModel {
		t.Errorf("Model() = %s; want %s", l1Cost.Model(), expectedModel)
	}
}
