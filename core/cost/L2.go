// Package cost provides implementations of various cost functions used in change point detection.
// These cost functions quantify the "error" or "dissimilarity" within a given segment of a signal,
// enabling algorithms to identify optimal segmentations.
package cost

import (
	"errors"
	"fmt" // Agregado para fmt.Errorf

	"github.com/theDataFlowClub/ruptures/core/base"
	"github.com/theDataFlowClub/ruptures/core/exceptions" // Import custom error types
	"github.com/theDataFlowClub/ruptures/core/stat"       // <-- NUEVO: Importar el paquete stat
	"github.com/theDataFlowClub/ruptures/core/types"      // Import general types like Matrix
	// For mathematical operations if needed, though not directly in L2 var
)

// CostL2 represents the L2 (Least Squared Deviation) cost function.
// It calculates the sum of squared deviations from the mean for a given segment.
// This cost function is common for detecting changes in the mean of a signal.
//
// The L2 cost for a segment [start:end] is calculated as:
// Sum_{i=start}^{end-1} ||signal[i] - mean(signal[start:end])||^2
// This is equivalent to (end - start) * var(signal[start:end]), where var is the variance.
//
// CostL2 implements the base.CostFunction interface.
type CostL2 struct {
	Signal  types.Matrix // The signal on which the cost is calculated. Shape (n_samples, n_features).
	MinSize int          // The minimum required size for a segment to be valid.

	minSegmentSize int // También necesitará esto para la interfaz
}

// NewCostL2 creates and returns a new instance of CostL2.
// This constructor function helps in initializing the struct with default values.
func NewCostL2() *CostL2 {
	return &CostL2{
		MinSize: 1, // Default minimum segment size, consistent with Python.
	}
}

// Fit sets the parameters for the CostL2 instance.
// It receives the signal and stores it internally for subsequent error calculations.
//
// Parameters:
//
//	signal: The input signal as a types.Matrix (or types.Vector which gets converted).
//	        Expected shape is (n_samples, n_features) or (n_samples,) for univariate.
//
// Returns:
//
//	An error if the signal is invalid (e.g., nil or empty), otherwise nil.
func (c *CostL2) Fit(signal types.Matrix) error {
	if signal == nil || len(signal) == 0 || (len(signal) > 0 && len(signal[0]) == 0) {
		return exceptions.ErrNotEnoughPoints // Consider a more specific error like `exceptions.ErrEmptySignal` if appropriate.
	}

	c.Signal = signal
	return nil
}

// Error calculates the L2 cost for the segment [start:end].
// The cost is computed as (end - start) * variance of the segment.
// This function efficiently calculates the variance by summing squared differences
// from the mean of each feature over the segment.
//
// Parameters:
//
//	start: The starting index of the segment (inclusive).
//	end: The ending index of the segment (exclusive).
//
// Returns:
//
//	float64: The calculated L2 cost for the segment.
//	error:   An error if the segment length (end - start) is less than `c.MinSize`
//	         (specifically, exceptions.ErrNotEnoughPoints) or if indices are out of bounds.
func (c *CostL2) Error(start, end int) (float64, error) {
	if c.Signal == nil {
		return 0.0, errors.New("CostL2: signal not fitted, call Fit() first")
	}

	// Primero, verifica los límites del segmento antes de extraerlo.
	if start < 0 || end > len(c.Signal) || start >= end {
		return 0.0, exceptions.ErrSegmentOutOfBounds
	}

	segmentLen := end - start
	if segmentLen < c.MinSize {
		return 0.0, exceptions.ErrNotEnoughPoints
	}

	segment := c.Signal[start:end]

	nFeatures := len(segment[0])
	totalVarianceSum := 0.0 // Renombrado para mayor claridad

	for col := 0; col < nFeatures; col++ {
		// Extract current feature values for the segment
		featureValues := make([]float64, segmentLen)
		for row := 0; row < segmentLen; row++ {
			featureValues[row] = segment[row][col]
		}

		// Calculate variance using the new stat package
		varianceFeature, err := stat.Variance(featureValues) // <-- CAMBIO AQUÍ
		if err != nil {
			// This error should ideally not happen if segmentLen > 0, but good to check
			return 0.0, fmt.Errorf("error calculating variance for feature %d: %w", col, err) // Agregado fmt para el error
		}
		totalVarianceSum += varianceFeature
	}

	return totalVarianceSum * float64(segmentLen), nil
}

// Model returns the name of the cost function model, which is "l2".
func (c *CostL2) Model() string {
	return "l2"
}

// init function is called automatically when the package is initialized.
// It registers the CostL2 constructor with the cost factory.
func init() {
	// ¡Registra la función de costo L2!
	RegisterCostFunction("l2", func() base.CostFunction {
		return NewCostL2()
	})
}
