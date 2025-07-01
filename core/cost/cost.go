package cost

import "github.com/theDataFlowClub/ruptures/core/types"

// CostFunction defines the interface for all cost functions used in change point detection.

// CostFunction defines the interface for all cost functions used in change point detection.
type CostFunction interface {
	// Fit prepares the cost function with the given signal.
	Fit(signal types.Matrix) error

	// Error calculates the cost for a segment [start, end).
	Error(start, end int) (float64, error)

	// Model returns the name of the cost function model (e.g., "l1", "l2", "rbf").
	Model() string

	// MinSize returns the minimum required length of a segment for this cost function.
	// This method ensures that all cost function implementations provide a minimum segment size.
	MinSize() int // <--- ¡ASEGÚRATE DE QUE ESTA LÍNEA ESTÉ AQUÍ!
}
