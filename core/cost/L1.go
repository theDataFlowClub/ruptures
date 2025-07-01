// Package cost provides implementations of various cost functions used in change point detection.
package cost

import (
	"errors" // For generic error (e.g., signal not fitted)
	"fmt"
	"math" // For math.Abs

	"github.com/theDataFlowClub/ruptures/core/base"       // For CostFunction interface
	"github.com/theDataFlowClub/ruptures/core/exceptions" // For custom error types
	"github.com/theDataFlowClub/ruptures/core/stat"       // <-- NUEVO: Importar el paquete stat
	"github.com/theDataFlowClub/ruptures/core/types"      // For Matrix type
	// "sort" // Ya no se necesita aquí porque Median está en stat
)

// CostL1 represents the L1 (Least Absolute Deviation) cost function.
// It calculates the sum of absolute deviations from the median for a given segment.
// This cost function is more robust to outliers compared to CostL2 (least squares).
//
// The L1 cost for a segment [start:end] is calculated as:
// Sum_{i=start}^{end-1} ||signal[i] - median(signal[start:end])||_1
// Where ||.||_1 is the L1 norm (sum of absolute differences).
//
// CostL1 implements the base.CostFunction interface.
type CostL1 struct {
	Signal  types.Matrix // The signal on which the cost is calculated. Shape (n_samples, n_features).
	MinSize int          // The minimum required size for a segment to be valid. Default is 2.

	minSegmentSize int // También necesitará esto para la interfaz

}

// NewCostL1 creates and returns a new instance of CostL1.
// This constructor function helps in initializing the struct with default values.
func NewCostL1() *CostL1 {
	return &CostL1{
		MinSize: 2, // Default minimum segment size for L1, as in Python.
	}
}

// Fit sets the parameters for the CostL1 instance.
// It receives the signal and stores it internally for subsequent error calculations.
//
// Parameters:
//
//	signal: The input signal as a types.Matrix.
//	        Expected shape is (n_samples, n_features).
//
// Returns:
//
//	An error if the signal is invalid (e.g., nil or empty), otherwise nil.
func (c *CostL1) Fit(signal types.Matrix) error {
	if signal == nil || len(signal) == 0 || (len(signal) > 0 && len(signal[0]) == 0) {
		return exceptions.ErrNotEnoughPoints // Or a more specific "ErrEmptySignal"
	}
	c.Signal = signal
	return nil
}

// Error calculates the L1 cost for the segment [start:end].
// The cost is computed as the sum of absolute deviations from the median for each feature.
//
// Parameters:
//
//	start: The starting index of the segment (inclusive).
//	end: The ending index of the segment (exclusive).
//
// Returns:
//
//	float64: The calculated L1 cost for the segment.
//	error:   An error if the segment is too short (exceptions.ErrNotEnoughPoints)
//	         or if indices are out of bounds (exceptions.ErrSegmentOutOfBounds).
func (c *CostL1) Error(start, end int) (float64, error) {
	if c.Signal == nil {
		return 0.0, errors.New("CostL1: signal not fitted, call Fit() first")
	}

	segmentLen := end - start

	// Check bounds and invalid segment definition (start >= end) first
	if start < 0 || end > len(c.Signal) || start >= end {
		return 0.0, exceptions.ErrSegmentOutOfBounds
	}

	// Check minimum size required for calculation (specifically for L1, min_size=2)
	if segmentLen < c.MinSize {
		return 0.0, exceptions.ErrNotEnoughPoints
	}

	segment := c.Signal[start:end]
	nFeatures := len(segment[0])
	totalAbsoluteDeviation := 0.0

	// Calculate median for each feature (column) and sum absolute deviations
	for col := 0; col < nFeatures; col++ {
		// Extract current feature values for the segment
		featureValues := make([]float64, segmentLen)
		for row := 0; row < segmentLen; row++ {
			featureValues[row] = segment[row][col]
		}

		// Calculate median using the new stat package
		medianFeature, err := stat.Median(featureValues) // <-- CAMBIO AQUÍ
		if err != nil {
			// This error should ideally not happen if segmentLen > 0, but good to check
			return 0.0, fmt.Errorf("error calculating median for feature %d: %w", col, err) // Agregado fmt para el error
		}

		// Sum absolute deviations from the median for the current feature
		for row := 0; row < segmentLen; row++ {
			totalAbsoluteDeviation += math.Abs(featureValues[row] - medianFeature)
		}
	}

	return totalAbsoluteDeviation, nil
}

// Model returns the name of the cost function model, which is "l1".
func (c *CostL1) Model() string {
	return "l1"
}

func init() {
	// ¡Registra la función de costo L1!
	RegisterCostFunction("l1", func() base.CostFunction {
		return NewCostL1()
	})
}
