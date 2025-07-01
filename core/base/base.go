// Package base defines the core interfaces and fundamental components
// for change point detection estimators and cost functions within the ruptures library.
// These interfaces establish clear contracts for algorithms and cost models,
// ensuring a consistent API and facilitating modularity and testability.
package base

import "github.com/theDataFlowClub/ruptures/core/types"

// Estimator is the base interface for all change point detection algorithms.
// Any algorithm implementing this interface must provide methods to:
//   - Fit: Prepare the estimator with the input signal.
//   - Predict: Compute the change points based on a given penalty.
//   - FitPredict: A convenience method that combines Fit and Predict in one call.
//
// Implementations should handle specific algorithm logic (e.g., PELT, BinSeg, DynP)
// and manage internal state required for prediction.
type Estimator interface {
	// Fit trains the estimator on the provided signal.
	// It typically performs initial computations or pre-processes the data
	// required for the prediction step.
	// Returns an error if the signal is invalid or fitting fails.
	Fit(signal types.Matrix) error
	// Predict computes the change points given a penalty value.
	// The penalty influences the number of detected change points;
	// higher penalties generally result in fewer breakpoints.
	// Returns a slice of breakpoint indices or an error if prediction fails.
	Predict(penalty float64) ([]int, error)
	// FitPredict is a convenience method that first fits the estimator to the signal
	// and then predicts the change points based on the provided penalty.
	// This method is useful for a streamlined workflow.
	// Returns a slice of breakpoint indices or an error if the operation fails.
	FitPredict(signal types.Matrix, penalty float64) ([]int, error)
}

// CostFunction is the base interface for all segment cost functions.
// Cost functions quantify the "cost" or "error" within a given segment of a signal.
// They are crucial for change point detection algorithms to evaluate potential segmentations.
// Any implementation must provide methods to:
//   - Fit: Prepare the cost function with the input signal (e.g., precompute sums, matrices).
//   - Error: Calculate the cost for a specific segment.
//   - Model: Return the name or type of the cost model (e.g., "l2", "rbf").
type CostFunction interface {
	// Fit prepares the cost function by processing the input signal.
	// This method is typically called once before computing segment costs,
	// allowing for pre-computation of necessary statistics (e.g., cumulative sums, Gram matrices).
	// Returns an error if fitting fails.
	Fit(signal types.Matrix) error
	// Error calculates the cost (or error) for a segment spanning from 'start' to 'end' indices (inclusive of start, exclusive of end).
	// The cost represents how well the data within this segment conforms to a specific model (e.g., constant mean, linear trend).
	// Returns the calculated cost as a float64 and an error if the segment is invalid (e.g., too short).
	Error(start, end int) (float64, error) // <--- ACTUALIZADO AQUÍ
	// Model returns a string identifier for the cost function (e.g., "l2", "rbf", "linear").
	// This can be useful for logging, debugging, or configuring algorithms based on the cost model.
	Model() string
}

// createSignal es una función de ayuda para convertir un slice de float64 en types.Matrix.
func createSignal(data []float64, dims int) types.Matrix {
	signal := make(types.Matrix, len(data)/dims)
	for i := 0; i < len(data)/dims; i++ {
		signal[i] = make([]float64, dims)
		copy(signal[i], data[i*dims:(i+1)*dims])
	}
	return signal
}
