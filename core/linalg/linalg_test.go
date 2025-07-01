package linalg_test

import (
	"math"
	"strings" // Para strings.Contains
	"testing"

	"github.com/theDataFlowClub/ruptures/core/linalg" // El paquete a probar
	"github.com/theDataFlowClub/ruptures/core/types"  // Para el tipo Matrix
)

// Define una pequeña tolerancia para comparaciones de punto flotante.
const floatTolerance = 1e-9

// Función auxiliar para comparar dos slices de float64 para igualdad aproximada.
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

// Función auxiliar para comparar dos types.Matrix para igualdad aproximada.
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

// ### **TestPdistSqEuclidean**

func TestPdistSqEuclidean(t *testing.T) {
	testCases := []struct {
		name          string
		matrix        types.Matrix // Usamos types.Matrix aquí para ser explícitos
		expectedDists []float64
		expectError   bool
		expectedError string
	}{
		{
			name:          "Simple2D_3x2",
			matrix:        types.Matrix{{0.0, 0.0}, {1.0, 1.0}, {0.0, 1.0}},
			expectedDists: []float64{2.0, 1.0, 1.0}, // (0,0)-(1,1)=2, (0,0)-(0,1)=1, (1,1)-(0,1)=1
			expectError:   false,
		},
		{
			name:          "Simple1D_3x1",
			matrix:        types.Matrix{{0.0}, {1.0}, {2.0}},
			expectedDists: []float64{1.0, 4.0, 1.0}, // (0)-(1)=1, (0)-(2)=4, (1)-(2)=1
			expectError:   false,
		},
		{
			name:          "EmptyMatrix",
			matrix:        types.Matrix{},
			expectedDists: []float64{},
			expectError:   false,
		},
		{
			name:          "SingleRowMatrix",
			matrix:        types.Matrix{{10.0, 20.0}},
			expectedDists: []float64{}, // pdist para N=1 devuelve vacío
			expectError:   false,
		},
		{
			name:          "MatrixWithInconsistentRows", // Renamed for clarity, it means variable features
			matrix:        types.Matrix{{1.0, 1.0}, {1.0}},
			expectedDists: nil,
			expectError:   true,
			expectedError: "inconsistent feature dimension",
		},
		{
			name:          "MatrixWithEmptyFirstRow",
			matrix:        types.Matrix{{}, {1.0, 2.0}},
			expectedDists: nil,
			expectError:   true,
			expectedError: "input matrix has zero features",
		},
		{
			name:          "MatrixWithZeroFeaturesAllRows",
			matrix:        types.Matrix{{}, {}, {}},
			expectedDists: nil,
			expectError:   true,
			expectedError: "input matrix has zero features",
		},
		{
			name:          "DuplicateRows",
			matrix:        types.Matrix{{1.0, 2.0}, {1.0, 2.0}, {3.0, 4.0}},
			expectedDists: []float64{0.0, 8.0, 8.0}, // (1,2)-(1,2)=0, (1,2)-(3,4)=8, (1,2)-(3,4)=8
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := linalg.PdistSqEuclidean(tc.matrix)

			if tc.expectError {
				if err == nil {
					t.Errorf("PdistSqEuclidean() expected an error containing %q, but got nil", tc.expectedError)
				} else if !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("PdistSqEuclidean() got unexpected error: %v, want error containing %q", err, tc.expectedError)
				}
			} else {
				if err != nil {
					t.Errorf("PdistSqEuclidean() got unexpected error: %v, want nil", err)
				}
				if !compareFloatSlices(result, tc.expectedDists, floatTolerance) {
					t.Errorf("PdistSqEuclidean() = %v; want %v", result, tc.expectedDists)
				}
			}
		})
	}
}

//### **TestSquareform**

func TestSquareform(t *testing.T) {
	testCases := []struct {
		name           string
		dists          []float64
		n              int // Número de observaciones originales
		expectedMatrix types.Matrix
		expectError    bool
		expectedError  string
	}{
		{
			name:           "3x3Matrix",
			dists:          []float64{2.0, 1.0, 1.0}, // Desde PdistSqEuclidean para matriz 3x2
			n:              3,
			expectedMatrix: types.Matrix{{0.0, 2.0, 1.0}, {2.0, 0.0, 1.0}, {1.0, 1.0, 0.0}},
			expectError:    false,
		},
		{
			name:           "2x2Matrix",
			dists:          []float64{5.0},
			n:              2,
			expectedMatrix: types.Matrix{{0.0, 5.0}, {5.0, 0.0}},
			expectError:    false,
		},
		{
			name:           "EmptyDists_n0",
			dists:          []float64{},
			n:              0,
			expectedMatrix: types.Matrix{},
			expectError:    false,
		},
		{
			name:           "EmptyDists_n1",
			dists:          []float64{},
			n:              1, // Squareform para n=1 significa una matriz 1x1 con 0.
			expectedMatrix: types.Matrix{{0.0}},
			expectError:    false,
		},
		{
			name:           "InconsistentDimensions",
			dists:          []float64{1.0, 2.0}, // Se esperan 3 elementos para n=3 (3*2/2 = 3)
			n:              3,
			expectedMatrix: nil,
			expectError:    true,
			expectedError:  "inconsistent dimensions",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := linalg.Squareform(tc.dists, tc.n)

			if tc.expectError {
				if err == nil {
					t.Errorf("Squareform() expected an error containing %q, but got nil", tc.expectedError)
				} else if !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("Squareform() got unexpected error: %v, want error containing %q", err, tc.expectedError)
				}
			} else {
				if err != nil {
					t.Errorf("Squareform() got unexpected error: %v, want nil", err)
				}
				if !compareMatrices(result, tc.expectedMatrix, floatTolerance) {
					t.Errorf("Squareform() = %v; want %v", result, tc.expectedMatrix)
				}
			}
		})
	}
}

// ### **TestClipSlice**

func TestClipSlice(t *testing.T) {
	testCases := []struct {
		name            string
		data            []float64
		min             float64
		max             float64
		expectedClipped []float64
	}{
		{
			name:            "NoClipping",
			data:            []float64{1.0, 2.0, 3.0},
			min:             0.0,
			max:             4.0,
			expectedClipped: []float64{1.0, 2.0, 3.0},
		},
		{
			name:            "ClipBelowMin",
			data:            []float64{-1.0, 0.5, 2.0},
			min:             0.0,
			max:             10.0,
			expectedClipped: []float64{0.0, 0.5, 2.0},
		},
		{
			name:            "ClipAboveMax",
			data:            []float64{5.0, 10.5, 12.0},
			min:             0.0,
			max:             10.0,
			expectedClipped: []float64{5.0, 10.0, 10.0},
		},
		{
			name:            "ClipBothEnds",
			data:            []float64{-5.0, 2.0, 15.0},
			min:             0.0,
			max:             10.0,
			expectedClipped: []float64{0.0, 2.0, 10.0},
		},
		{
			name:            "EmptySlice",
			data:            []float64{},
			min:             0.0,
			max:             10.0,
			expectedClipped: []float64{},
		},
		{
			name:            "MinEqualsMax",
			data:            []float64{-1.0, 5.0, 10.0},
			min:             5.0,
			max:             5.0,
			expectedClipped: []float64{5.0, 5.0, 5.0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Hacer una copia para asegurar que el original no sea modificado por ejecuciones de prueba anteriores
			dataCopy := make([]float64, len(tc.data))
			copy(dataCopy, tc.data)

			result := linalg.ClipSlice(dataCopy, tc.min, tc.max)
			if !compareFloatSlices(result, tc.expectedClipped, floatTolerance) {
				t.Errorf("ClipSlice() = %v; want %v", result, tc.expectedClipped)
			}
		})
	}
}

// ### **TestDiagonalSum**

func TestDiagonalSum(t *testing.T) {
	testCases := []struct {
		name          string
		matrix        types.Matrix
		expectedSum   float64
		expectError   bool
		expectedError string
	}{
		{
			name:        "SimpleSquareMatrix",
			matrix:      types.Matrix{{1.0, 2.0}, {3.0, 4.0}},
			expectedSum: 5.0, // 1 + 4
			expectError: false,
		},
		{
			name:        "LargerSquareMatrix",
			matrix:      types.Matrix{{1.0, 0.0, 0.0}, {0.0, 2.0, 0.0}, {0.0, 0.0, 3.0}},
			expectedSum: 6.0, // 1 + 2 + 3
			expectError: false,
		},
		{
			name:        "SingleElementMatrix",
			matrix:      types.Matrix{{7.0}},
			expectedSum: 7.0,
			expectError: false,
		},
		{
			name:          "NonSquareMatrix_Rectangular",
			matrix:        types.Matrix{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}},
			expectedSum:   0.0,
			expectError:   true,
			expectedError: "matrix is not square for diagonal sum",
		},
		{
			name:          "NonSquareMatrix_Ragged",
			matrix:        types.Matrix{{1.0, 2.0}, {3.0}}, // Arreglo irregular
			expectedSum:   0.0,
			expectError:   true,
			expectedError: "matrix is not square for diagonal sum",
		},
		{
			name:          "EmptyMatrix",
			matrix:        types.Matrix{},
			expectedSum:   0.0,
			expectError:   true,
			expectedError: "empty matrix for diagonal sum",
		},
		{
			name:          "MatrixWithEmptyRows", // Matriz inválida para esta función
			matrix:        types.Matrix{{1.0}, {}},
			expectedSum:   0.0,
			expectError:   true,
			expectedError: "matrix is not square for diagonal sum", // La captura porque len(row) != nRows
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := linalg.DiagonalSum(tc.matrix)

			if tc.expectError {
				if err == nil {
					t.Errorf("DiagonalSum() expected an error containing %q, but got nil", tc.expectedError)
				} else if !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("DiagonalSum() got unexpected error: %v, want error containing %q", err, tc.expectedError)
				}
			} else {
				if err != nil {
					t.Errorf("DiagonalSum() got unexpected error: %v, want nil", err)
				}
				if math.Abs(result-tc.expectedSum) > floatTolerance {
					t.Errorf("DiagonalSum() = %f; want %f (diff: %f)", result, tc.expectedSum, math.Abs(result-tc.expectedSum))
				}
			}
		})
	}
}

// ### **TestSum**

func TestSum(t *testing.T) {
	testCases := []struct {
		name          string
		matrix        types.Matrix // Usamos types.Matrix aquí
		expectedSum   float64
		expectError   bool
		expectedError string
	}{
		{
			name:        "SimpleMatrix",
			matrix:      types.Matrix{{1.0, 2.0}, {3.0, 4.0}},
			expectedSum: 10.0, // 1+2+3+4
			expectError: false,
		},
		{
			name:        "MatrixWithNegativeValues",
			matrix:      types.Matrix{{-1.0, 0.0}, {1.0, -2.0}},
			expectedSum: -2.0, // -1+0+1-2
			expectError: false,
		},
		{
			name:        "SingleElementMatrix",
			matrix:      types.Matrix{{5.0}},
			expectedSum: 5.0,
			expectError: false,
		},
		{
			name:          "EmptyMatrix",
			matrix:        types.Matrix{},
			expectedSum:   0.0, // El valor de retorno es 0.0 si hay error
			expectError:   true,
			expectedError: "empty matrix for sum",
		},
		{
			name:          "MatrixWithEmptyRowInMiddle", // Test para tu nueva validación de fila vacía
			matrix:        types.Matrix{{1.0, 2.0}, {}, {3.0, 4.0}},
			expectedSum:   0.0,
			expectError:   true,
			expectedError: "matrix contains empty row at index 1", // Ajusta el mensaje esperado
		},
		{
			name:        "MatrixWithZeroes",
			matrix:      types.Matrix{{0.0, 0.0}, {0.0, 0.0}},
			expectedSum: 0.0,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sum, err := linalg.Sum(tc.matrix)

			if tc.expectError {
				if err == nil {
					t.Errorf("Sum() expected an error containing %q, but got nil", tc.expectedError)
				} else if !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("Sum() got unexpected error: %v, want error containing %q", err, tc.expectedError)
				}
			} else {
				if err != nil {
					t.Errorf("Sum() got unexpected error: %v", err)
				}
				if math.Abs(sum-tc.expectedSum) > floatTolerance { // Usa tolerancia para floats
					t.Errorf("expected sum %v, got %v (diff: %f)", tc.expectedSum, sum, math.Abs(sum-tc.expectedSum))
				}
			}
		})
	}
}

// ... (tus tests existentes para PdistSqEuclidean, DiagonalSum, Sum, ClipSlice, etc.)

func TestDot(t *testing.T) {
	testCases := []struct {
		name        string
		vec1        types.Vector
		vec2        types.Vector
		expected    float64
		expectError bool
		errorMsg    string
	}{
		{"ValidDotProduct", types.Vector{1, 2, 3}, types.Vector{4, 5, 6}, 32.0, false, ""},
		{"ZeroVectors", types.Vector{0, 0, 0}, types.Vector{1, 2, 3}, 0.0, false, ""},
		{"NegativeValues", types.Vector{-1, -2}, types.Vector{3, 4}, -11.0, false, ""},
		{"MismatchedLengths", types.Vector{1, 2}, types.Vector{3}, 0.0, true, "input vectors must have the same length"},
		{"EmptyVectors", types.Vector{}, types.Vector{}, 0.0, false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := linalg.Dot(tc.vec1, tc.vec2)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				} else if !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("Got error %v, want error containing %q", err, tc.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got %v", err)
				}
				if math.Abs(result-tc.expected) > floatTolerance {
					t.Errorf("Dot product mismatch. Got %f, want %f", result, tc.expected)
				}
			}
		})
	}
}

func TestSquaredEuclideanDistance(t *testing.T) {
	testCases := []struct {
		name        string
		vec1        types.Vector
		vec2        types.Vector
		expected    float64
		expectError bool
		errorMsg    string
	}{
		{"SimpleDistance", types.Vector{0, 0}, types.Vector{3, 4}, 25.0, false, ""},  // (3-0)^2 + (4-0)^2 = 9 + 16 = 25
		{"NegativeValues", types.Vector{-1, -1}, types.Vector{1, 1}, 8.0, false, ""}, // (1-(-1))^2 + (1-(-1))^2 = 2^2 + 2^2 = 4 + 4 = 8
		{"IdenticalVectors", types.Vector{5, 5}, types.Vector{5, 5}, 0.0, false, ""},
		{"MismatchedLengths", types.Vector{1, 2}, types.Vector{3}, 0.0, true, "input vectors must have the same dimension"},
		{"EmptyVectors", types.Vector{}, types.Vector{}, 0.0, false, ""}, // Distance between empty vectors is 0
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := linalg.SquaredEuclideanDistance(tc.vec1, tc.vec2)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				} else if !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("Got error %v, want error containing %q", err, tc.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got %v", err)
				}
				if math.Abs(result-tc.expected) > floatTolerance {
					t.Errorf("Squared Euclidean distance mismatch. Got %f, want %f", result, tc.expected)
				}
			}
		})
	}
}

func TestVectorNorm(t *testing.T) {
	testCases := []struct {
		name        string
		vec         types.Vector
		expected    float64
		expectError bool
		errorMsg    string // No error expected for now, but good to keep
	}{
		{"SimpleVector", types.Vector{3, 4}, 5.0, false, ""}, // sqrt(3^2 + 4^2) = sqrt(9+16) = sqrt(25) = 5
		{"ZeroVector", types.Vector{0, 0, 0}, 0.0, false, ""},
		{"NegativeValues", types.Vector{-3, -4}, 5.0, false, ""}, // sqrt((-3)^2 + (-4)^2) = sqrt(9+16) = 5
		{"SingleElement", types.Vector{7}, 7.0, false, ""},
		{"EmptyVector", types.Vector{}, 0.0, false, ""}, // Norm of empty vector is 0
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := linalg.VectorNorm(tc.vec)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				} else if !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("Got error %v, want error containing %q", err, tc.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got %v", err)
				}
				if math.Abs(result-tc.expected) > floatTolerance {
					t.Errorf("Vector norm mismatch. Got %f, want %f", result, tc.expected)
				}
			}
		})
	}
}
