package base

import (
	"github.com/theDataFlowClub/ruptures/core/utils"
)

// SumOfCosts calculates the total cost of a segmentation given a CostFunction and a list of breakpoints.
// It iterates through the segments defined by the breakpoints and sums the cost for each segment.
//
// The 'bkps' slice should represent the change point indices, with the last element
// conventionally being the total number of samples (n_samples) to define the end of the last segment.
// The function internally prepends a 0 to 'bkps' to represent the start of the first segment,
// and then uses `utils.Pairwise` to generate (start, end) pairs for each segment.
//
// Parameters:
//
//	cost: An implementation of the CostFunction interface, used to calculate the error for each segment.
//	bkps: A slice of integers representing the breakpoints. Example: [b1, b2, ..., bn_samples].
//	      If bkps is empty, the total cost is 0.0.
//
// Returns:
//
//	float64: The sum of the costs of all segments defined by the breakpoints.
//	error:   An error if any segment's cost calculation fails (e.g., due to ErrNotEnoughPoints).
func SumOfCosts(cost CostFunction, bkps []int) (float64, error) { // <--- ACTUALIZADO AQUÍ
	if len(bkps) == 0 {
		return 0.0, nil // <--- ACTUALIZADO AQUÍ
	}
	// Prepend 0 to the breakpoints to represent the start of the first segment.
	// Example: if bkps = [100, 200, 300], then breaks = [0, 100, 200, 300].
	breaks := append([]int{0}, bkps...)
	// Generate pairs (start, end) for each segment.
	// Example: [(0, 100), (100, 200), (200, 300)]
	pairs := utils.Pairwise(breaks)
	var sum float64
	// Sum the error for each segment defined by the pairs.
	for _, p := range pairs {
		segmentCost, err := cost.Error(p.First, p.Second) // <--- ACTUALIZADO AQUÍ
		if err != nil {
			return 0.0, err // <--- Maneja el error y lo propaga
		}
		sum += segmentCost
	}
	return sum, nil // <--- ACTUALIZADO AQUÍ
}
