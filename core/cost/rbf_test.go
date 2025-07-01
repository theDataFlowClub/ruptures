package cost_test

import (
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/theDataFlowClub/ruptures/core/cost"
	"github.com/theDataFlowClub/ruptures/core/exceptions"
	"github.com/theDataFlowClub/ruptures/core/types"
)

// Define una pequeña tolerancia para comparaciones de punto flotante.
const floatTolerance = 1e-6 // Mayor tolerancia para RBF debido a las exponenciales

// Helper function to compare two float64 slices for approximate equality
func compareFloatSlices(a, b []float64, tolerance float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > tolerance {
			return false
		}
	}
	return true
}

// Helper function to compare two types.Matrix for approximate equality
func compareMatrices(a, b types.Matrix, tolerance float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !compareFloatSlices(a[i], b[i], tolerance) {
			return false
		}
	}
	return true
}

// ### **TestNewCostRbf**

func TestNewCostRbf(t *testing.T) {
	// Test with nil gamma
	c1 := cost.NewCostRbf(nil)
	if c1 == nil {
		t.Fatal("NewCostRbf returned nil")
	}
	if c1.MinSize() != 1 {
		t.Errorf("Expected MinSize 1, got %d", c1.MinSize())
	}
	if c1.Gamma != nil {
		t.Errorf("Expected Gamma to be nil, got %v", *c1.Gamma)
	}
	if c1.Model() != "rbf" {
		t.Errorf("Expected Model 'rbf', got %s", c1.Model())
	}

	// Test with a specific gamma
	g := 0.5
	c2 := cost.NewCostRbf(&g)
	if c2 == nil {
		t.Fatal("NewCostRbf returned nil with gamma")
	}
	if c2.Gamma == nil || *c2.Gamma != 0.5 {
		t.Errorf("Expected Gamma 0.5, got %v", c2.Gamma)
	}
}

// ### TestCostRbfFit

func TestCostRbfFit(t *testing.T) {
	testCases := []struct {
		name        string
		signal      types.Matrix
		expectError bool
		expectedErr error
	}{
		{
			name:   "ValidSignal",
			signal: types.Matrix{{1.0, 2.0}, {3.0, 4.0}, {5.0, 6.0}},
		},
		{
			name:        "EmptySignal",
			signal:      types.Matrix{},
			expectError: true,
			expectedErr: exceptions.ErrNotEnoughPoints,
		},
		{
			name:        "SignalWithEmptyRows",
			signal:      types.Matrix{{1.0}, {}}, // Should be caught by internal linalg functions
			expectError: true,
			expectedErr: errors.New("inconsistent feature dimension"), // Error from PdistSqEuclidean
		},
		{
			name:        "NilSignal",
			signal:      nil,
			expectError: true,
			expectedErr: exceptions.ErrNotEnoughPoints,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rbfCost := cost.NewCostRbf(nil) // Start with nil gamma
			err := rbfCost.Fit(tc.signal)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				} else if !errors.Is(err, tc.expectedErr) && (tc.expectedErr.Error() != "" && !strings.Contains(err.Error(), tc.expectedErr.Error())) {
					t.Errorf("Got unexpected error: %v, want error %v (containing %q)", err, tc.expectedErr, tc.expectedErr.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got %v", err)
				}
				if rbfCost.Signal == nil || len(rbfCost.Signal) == 0 {
					t.Errorf("Signal was not set after Fit()")
				}
				// After fitting, if gamma was nil, it should have been calculated.
				if rbfCost.Gamma == nil {
					t.Error("Gamma was not set after Fit() with nil initial gamma")
				}
				// Also check that gram matrix was calculated and cached
				// Utiliza el método exportado para acceder al estado cacheado.
				if rbfCost.GetCachedGramForTest() == nil && len(tc.signal) > 1 {
					t.Errorf("Gram matrix was not computed and cached after Fit()")
				}
			}
		})
	}
}

// ### TestCostRbfGetGram

func TestCostRbfGetGram(t *testing.T) {
	// Test case: Basic signal, check gram matrix values and gamma heuristic
	t.Run("BasicSignalWithNilGamma", func(t *testing.T) {
		signal := types.Matrix{{0.0}, {1.0}} // Simple 1D signal
		rbfCost := cost.NewCostRbf(nil)      // gamma is nil

		err := rbfCost.Fit(signal) // Fit the signal
		if err != nil {
			t.Fatalf("Fit failed: %v", err)
		}

		gram, err := rbfCost.GetGram() // This should trigger calculation
		if err != nil {
			t.Fatalf("GetGram failed: %v", err)
		}

		// Expected calculations for {0}, {1}:
		// PdistSqEuclidean: [1.0] (distance between 0 and 1 is 1 squared)
		// Median: 1.0 (for [1.0])
		// Gamma: 1/1.0 = 1.0
		// Clipped: [1.0] (no clipping for 1.0 within 1e-2, 1e2)
		// squareform(-K): types.Matrix{{0, -1}, {-1, 0}}
		// np.exp(squareform(-K)): types.Matrix{{exp(0), exp(-1)}, {exp(-1), exp(0)}}
		// Expected Gram Matrix: {{1.0, 0.36787944}, {0.36787944, 1.0}} (approx)
		expectedGram := types.Matrix{{1.0, math.Exp(-1.0)}, {math.Exp(-1.0), 1.0}}

		if !compareMatrices(gram, expectedGram, floatTolerance) {
			t.Errorf("Calculated Gram matrix mismatch.\nGot:\n%v\nWant:\n%v", gram, expectedGram)
		}
		if rbfCost.Gamma == nil || *rbfCost.Gamma != 1.0 {
			t.Errorf("Expected gamma to be 1.0 (median heuristic), got %v", rbfCost.Gamma)
		}

		// ... (el resto del test para lazy loading)
		cachedGramForModification := rbfCost.GetCachedGramForTest()
		if len(cachedGramForModification) > 0 && len(cachedGramForModification[0]) > 0 {
			cachedGramForModification[0][0] = 99.0
		} else {
			t.Fatalf("Gram matrix unexpectedly empty for modification test")
		}

		cachedGram, err := rbfCost.GetGram()
		if err != nil {
			t.Fatalf("GetGram failed on second call: %v", err)
		}
		if cachedGram[0][0] != 99.0 {
			t.Errorf("Gram matrix was recomputed, expected cached value (99.0), got %f", cachedGram[0][0])
		}
	})

	// Test case: Basic signal, with predefined gamma
	t.Run("BasicSignalWithDefinedGamma", func(t *testing.T) {
		signal := types.Matrix{{0.0}, {1.0}}
		g := 2.0
		rbfCost := cost.NewCostRbf(&g) // predefined gamma

		err := rbfCost.Fit(signal)
		if err != nil {
			t.Fatalf("Fit failed: %v", err)
		}

		gram, err := rbfCost.GetGram()
		if err != nil {
			t.Fatalf("GetGram failed: %v", err)
		}

		// Expected calculations for {0}, {1} with gamma=2.0:
		// PdistSqEuclidean: [1.0]
		// Gamma scaling: 1.0 * 2.0 = 2.0
		// Clipped: [2.0]
		// np.exp(squareform(-K)): types.Matrix{{exp(0), exp(-2)}, {exp(-2), exp(0)}}
		// Expected Gram Matrix: {{1.0, 0.13533528}, {0.13533528, 1.0}} (approx)
		expectedGram := types.Matrix{{1.0, math.Exp(-2.0)}, {math.Exp(-2.0), 1.0}}

		if !compareMatrices(gram, expectedGram, floatTolerance) {
			t.Errorf("Calculated Gram matrix mismatch with predefined gamma.\nGot:\n%v\nWant:\n%v", gram, expectedGram)
		}
		if rbfCost.Gamma == nil || *rbfCost.Gamma != 2.0 {
			t.Errorf("Expected gamma to be 2.0 (predefined), got %v", rbfCost.Gamma)
		}
	})

	// Test case: GetGram without fitting signal
	t.Run("GetGramWithoutFit", func(t *testing.T) {
		rbfCost := cost.NewCostRbf(nil) // No signal fitted yet
		_, err := rbfCost.GetGram()
		if err == nil {
			t.Error("Expected error when calling GetGram before Fit(), got nil")
		}
		if err != nil && !strings.Contains(err.Error(), "signal not fitted or empty") {
			t.Errorf("Expected 'signal not fitted' error, got: %v", err)
		}
	})

	// Test with signal leading to medianDistSq = 0 (e.g., all identical points)
	t.Run("IdenticalPointsMedianHeuristic", func(t *testing.T) {
		signal := types.Matrix{{1.0, 1.0}, {1.0, 1.0}, {1.0, 1.0}}
		rbfCost := cost.NewCostRbf(nil) // gamma is nil

		err := rbfCost.Fit(signal)
		if err != nil {
			t.Fatalf("Fit failed: %v", err)
		}

		gram, err := rbfCost.GetGram()
		if err != nil {
			t.Fatalf("GetGram failed: %v", err)
		}

		// All distances are 0.0. Median is 0.0. Gamma should default to 1.0.
		// BUT, the C implementation clips gamma * squared_distance to 0.01 before exp.
		// So exp(-0.01) is expected for off-diagonal (where distSq is 0).
		expectedValForZeroDist := math.Exp(-0.01) // This is 0.9900498337491681

		// --- CORRECCIÓN AQUÍ ---
		expectedGram := types.Matrix{
			{1.0, expectedValForZeroDist, expectedValForZeroDist},
			{expectedValForZeroDist, 1.0, expectedValForZeroDist},
			{expectedValForZeroDist, expectedValForZeroDist, 1.0},
		}

		if !compareMatrices(gram, expectedGram, floatTolerance) {
			t.Errorf("Calculated Gram matrix mismatch for identical points.\nGot:\n%v\nWant:\n%v", gram, expectedGram)
		}
		if rbfCost.Gamma == nil || *rbfCost.Gamma != 1.0 {
			t.Errorf("Expected gamma to be 1.0 (median heuristic for zero median), got %v", rbfCost.Gamma)
		}
	})
}

// ### TestCostRbfError

func TestCostRbfError(t *testing.T) {
	// Signal y Gram Matrix precalculados para facilitar la verificación
	// signal: {{0.0, 0.0}, {1.0, 1.0}, {0.0, 1.0}, {1.0, 0.0}}
	// Para simplificar, asumimos gamma=1.0 para este ejemplo.
	// Distancias cuadradas (sqeuclidean) entre 4 puntos:
	// (0,1): 2.0
	// (0,2): 1.0
	// (0,3): 1.0
	// (1,2): 1.0
	// (1,3): 1.0
	// (2,3): 2.0
	// K (condensed): [2.0, 1.0, 1.0, 1.0, 1.0, 2.0]
	// If gamma = 1.0, then values for exp(-K) are:
	// exp(-2) = 0.135335
	// exp(-1) = 0.367879
	// Gram Matrix (approx) for gamma=1.0 (diagonal 1.0):
	//     P0       P1       P2       P3
	// P0 [[1.0,   0.135,   0.368,   0.368],
	// P1  [0.135,  1.0,    0.368,   0.368],
	// P2  [0.368,  0.368,   1.0,    0.135],
	// P3  [0.368,  0.368,   0.135,   1.0  ]]
	predefinedSignal := types.Matrix{{0.0, 0.0}, {1.0, 1.0}, {0.0, 1.0}, {1.0, 0.0}}
	gammaVal := 1.0
	predefinedCostRbf := cost.NewCostRbf(&gammaVal)
	if err := predefinedCostRbf.Fit(predefinedSignal); err != nil {
		t.Fatalf("Failed to fit predefined signal for error test: %v", err)
	}
	// Force gram matrix calculation once to get actual values
	_, err := predefinedCostRbf.GetGram()
	if err != nil {
		t.Fatalf("Failed to get Gram matrix for error test: %v", err)
	}
	// Access cached Gram matrix via the exported helper method
	precomputedGram := predefinedCostRbf.GetCachedGramForTest()
	//precomputedGram := predefinedCostRbf._gram // Access cached Gram

	testCases := []struct {
		name         string
		start        int
		end          int
		expectedCost float64
		expectError  bool
		expectedErr  error
	}{
		{
			name:  "ValidSegment_0_2", // Segment [0:2], points P0, P1
			start: 0,
			end:   2,
			// subGram P0,P1: {{1.0, 0.135}, {0.135, 1.0}}
			// diagSum = 1.0 + 1.0 = 2.0
			// totalSum = 1.0 + 0.135 + 0.135 + 1.0 = 2.27
			// len = 2
			// cost = 2.0 - (2.27 / 2) = 2.0 - 1.135 = 0.865
			expectedCost: 0.8646647167633873, // Actual value from np.diagonal(sub_gram).sum() - sub_gram.sum() / (end-start)
			expectError:  false,
		},
		{
			name:  "ValidSegment_1_4", // Segment [1:4], points P1, P2, P3
			start: 1,
			end:   4,
			// subGram P1,P2,P3:
			// {{1.0,    0.368,  0.368},
			//  {0.368,   1.0,   0.135},
			//  {0.368,  0.135,   1.0   }}
			// diagSum = 1.0 + 1.0 + 1.0 = 3.0
			// totalSum = 1.0+0.368+0.368 + 0.368+1.0+0.135 + 0.368+0.135+1.0 = 4.742
			// len = 3
			// cost = 3.0 - (4.742 / 3) = 3.0 - 1.58066 = 1.41934
			expectedCost: 1.419270556280335, // Actual value
			expectError:  false,
		},
		{
			name:         "ValidSegment_Length1", // Renombrado para claridad
			start:        0,
			end:          1,
			expectedCost: 0.0, // Costo de un segmento de longitud 1
			expectError:  false,
		},
		{
			name:         "SegmentTooShort_Length0",
			start:        0,
			end:          0, // length 0
			expectedCost: 0.0,
			expectError:  true,
			expectedErr:  exceptions.ErrNotEnoughPoints,
		},
		{
			name:         "InvalidSegmentBounds_EndTooLarge",
			start:        0,
			end:          5, // out of bounds for signal len 4
			expectedCost: 0.0,
			expectError:  true,
			expectedErr:  errors.New("sub-segment indices out of bounds"),
		},
		{
			name:         "InvalidSegmentBounds_StartTooLarge",
			start:        4,
			end:          4, // length 0
			expectedCost: 0.0,
			expectError:  true,
			// ¡Cambiar a este error!
			expectedErr: exceptions.ErrNotEnoughPoints, // <--- CAMBIO AQUÍ
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rbfCost := cost.NewCostRbf(&gammaVal) // Use predefined gamma
			rbfCost.Fit(predefinedSignal)         // Fit the signal (already done above, but good practice)

			//rbfCost._gram = precomputedGram       // Ensure cached gram is used for consistent testing
			// En TestCostRbfError
			// rbfCost._gram = precomputedGram // Esto es lo que da error
			// Cambia a:
			rbfCost.SetCachedGramForTest(precomputedGram)

			costVal, err := rbfCost.Error(tc.start, tc.end)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error for %s, but got nil", tc.name)
				} else if !errors.Is(err, tc.expectedErr) && (tc.expectedErr.Error() != "" && !strings.Contains(err.Error(), tc.expectedErr.Error())) {
					t.Errorf("For %s, got unexpected error: %v, want error %v (containing %q)", tc.name, err, tc.expectedErr, tc.expectedErr.Error())
				}
			} else {
				if err != nil {
					t.Errorf("For %s, did not expect an error, but got %v", tc.name, err)
				}
				if math.Abs(costVal-tc.expectedCost) > floatTolerance {
					t.Errorf("For %s, cost mismatch.\nGot: %f\nWant: %f\nDifference: %f", tc.name, costVal, tc.expectedCost, math.Abs(costVal-tc.expectedCost))
				}
			}
		})
	}
}
