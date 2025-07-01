package linalg

import (
	"errors"
	"fmt"
	"math"

	"github.com/theDataFlowClub/ruptures/core/types" // Import types.Matrix
)

// PdistSqEuclidean calculates the pairwise squared Euclidean distances between
// rows of a matrix. It returns a condensed distance matrix (1D array).
//
// Equivalent to scipy.spatial.distance.pdist(matrix, metric="sqeuclidean").
//
// Parameters:
//
//	matrix: The input matrix (n_samples, n_features).
//
// Returns:
//
//	[]float64: A 1D array of squared Euclidean distances. The order is
//	           (row_0, row_1), (row_0, row_2), ..., (row_0, row_n-1),
//	           (row_1, row_2), ..., (row_1, row_n-1), ..., (row_n-2, row_n-1).
//	error:     An error if the input matrix is invalid (e.g., empty or inconsistent dimensions).
func PdistSqEuclidean(matrix types.Matrix) ([]float64, error) {
	nSamples := len(matrix)
	if nSamples == 0 || nSamples == 1 {
		return []float64{}, nil // No distances for empty or single-row matrix
	}
	// Determine nFeatures from the first row.
	// Check if the first row exists and is not empty.
	if len(matrix[0]) == 0 {
		return nil, errors.New("input matrix has zero features (first row is empty)")
	}
	nFeatures := len(matrix[0])

	// Calculate the size of the condensed distance matrix
	numPairs := nSamples * (nSamples - 1) / 2
	distances := make([]float64, numPairs)
	k := 0 // Index for the distances array

	for i := 0; i < nSamples; i++ {
		// Ensure consistent feature dimension for matrix[i]
		if len(matrix[i]) != nFeatures {
			return nil, fmt.Errorf("inconsistent feature dimension at row %d: got %d, want %d", i, len(matrix[i]), nFeatures)
		}
		for j := i + 1; j < nSamples; j++ {
			// Ensure consistent feature dimension for matrix[j] BEFORE accessing its elements
			if len(matrix[j]) != nFeatures {
				return nil, fmt.Errorf("inconsistent feature dimension at row %d: got %d, want %d", j, len(matrix[j]), nFeatures)
			}

			sumSqDiff := 0.0
			for f := 0; f < nFeatures; f++ {
				// This line (line 55 in your output) is now safe because matrix[j] is guaranteed to have nFeatures elements
				diff := matrix[i][f] - matrix[j][f]
				sumSqDiff += diff * diff
			}
			distances[k] = sumSqDiff
			k++
		}
	}
	return distances, nil
}

// Squareform converts a 1D condensed distance matrix into a 2D square symmetric matrix.
//
// Equivalent to scipy.spatial.distance.squareform(distances).
//
// Parameters:
//
//	distances: The 1D condensed distance matrix (from PdistSqEuclidean).
//	n:         The number of original observations (samples).
//
// Returns:
//
//	types.Matrix: The 2D square symmetric distance matrix.
//	error:        An error if the input dimensions are inconsistent.
func Squareform(distances []float64, n int) (types.Matrix, error) {
	// The number of elements in condensed form should be n * (n - 1) / 2
	expectedLen := n * (n - 1) / 2
	if len(distances) != expectedLen {
		return nil, fmt.Errorf("inconsistent dimensions: len(distances) = %d, expected %d for n = %d", len(distances), expectedLen, n)
	}

	// Create a square matrix of size n x n
	matrix := make(types.Matrix, n)
	for i := range matrix {
		matrix[i] = make([]float64, n)
	}

	k := 0 // Index for the distances array
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			matrix[i][j] = distances[k]
			matrix[j][i] = distances[k] // Symmetric
			k++
		}
	}
	return matrix, nil
}

// ClipSlice clips the values in a slice to be within the specified min and max.
//
// Equivalent to numpy.clip(data, min, max, out=data).
//
// Parameters:
//
//	data: The input slice of float64 values.
//	min:  The minimum value.
//	max:  The maximum value.
//
// Returns:
//
//	[]float64: The clipped slice (modified in place).
func ClipSlice(data []float64, min, max float64) []float64 {
	for i := range data {
		if data[i] < min { // <-- Aquí está el problema si min > 0
			data[i] = min
		} else if data[i] > max {
			data[i] = max
		}
	}
	return data
}

// DiagonalSum calculates the sum of the elements on the main diagonal of a square matrix.
//
// Parameters:
//
//	matrix: The input square matrix.
//
// Returns:
//
//	float64: The sum of the diagonal elements.
//	error:   An error if the matrix is not square or empty.
func DiagonalSum(matrix types.Matrix) (float64, error) {
	nRows := len(matrix)
	if nRows == 0 {
		return 0.0, errors.New("empty matrix for diagonal sum")
	}

	sum := 0.0
	for i := 0; i < nRows; i++ {
		if len(matrix[i]) != nRows {
			return 0.0, errors.New("matrix is not square for diagonal sum")
		}
		sum += matrix[i][i]
	}
	return sum, nil
}

// Sum calculates the sum of all elements in a matrix.
//
// Parameters:
//
//	matrix: The input matrix.
//
// Returns:
//
//	float64: The sum of all elements.
//	error:   An error if the matrix is empty.
//
// Sum computes the sum of all elements in a matrix.
func Sum(matrix [][]float64) (float64, error) {
	if len(matrix) == 0 {
		return 0.0, errors.New("empty matrix for sum")
	}
	var sum float64
	for i, row := range matrix {
		if len(row) == 0 { // maneja explícitamente el caso de filas vacías dentro de la matriz,
			return 0.0, fmt.Errorf("matrix contains empty row at index %d", i)
		}
		for _, val := range row {
			sum += val
		}
	}
	return sum, nil
}

// core/linalg/linalg.go (Adiciones o verificaciones)

// Dot calcula el producto punto de dos vectores.
// Retorna un error si los vectores no tienen la misma longitud.
func Dot(x, y types.Vector) (float64, error) { // <--- CAMBIO AQUÍ: Añadir ", error"
	if len(x) != len(y) {
		return 0, errors.New("linalg.Dot: input vectors must have the same length")
	}
	res := 0.0
	for i := range x {
		res += x[i] * y[i]
	}
	return res, nil // <--- CAMBIO AQUÍ: Añadir ", nil"
}

// SquaredEuclideanDistance calcula la distancia euclidiana al cuadrado entre dos vectores.
// Retorna un error si los vectores no tienen la misma longitud.
func SquaredEuclideanDistance(x, y types.Vector) (float64, error) { // <--- CAMBIO AQUÍ
	if len(x) != len(y) {
		return 0, errors.New("linalg.SquaredEuclideanDistance: input vectors must have the same dimension")
	}
	sumSq := 0.0
	for i := range x {
		diff := x[i] - y[i]
		sumSq += diff * diff
	}
	return sumSq, nil // <--- CAMBIO AQUÍ
}

// VectorNorm calcula la norma euclidiana (L2 norm) de un vector.
// La norma es la raíz cuadrada de la suma de los cuadrados de sus elementos.
// Retorna un error si el vector es nil (aunque un vector vacío tiene norma 0).
func VectorNorm(x types.Vector) (float64, error) { // <--- CAMBIO AQUÍ
	if x == nil { // Considerar si un slice nil debe dar error o norma 0
		return 0.0, errors.New("linalg.VectorNorm: input vector is nil")
	}
	if len(x) == 0 {
		return 0.0, nil
	}
	sumSq := 0.0
	for _, val := range x {
		sumSq += val * val
	}
	return math.Sqrt(sumSq), nil // <--- CAMBIO AQUÍ
}
