// Package stat provides common statistical functions for numerical data.
// These functions are designed to be general-purpose and reusable across
// different parts of the ruptures library, such as cost functions
// and potentially other analysis modules.
package stat

import (
	"errors"
	"sort"
)

// Median calculates the median of a slice of float64 values.
//
// Parameters:
//
//	data: The slice of float64 for which to calculate the median.
//
// Returns:
//
//	float64: The calculated median.
//	error:   An error if the input slice is empty (errors.New("empty slice for median calculation")).
func Median(data []float64) (float64, error) {
	if len(data) == 0 {
		return 0.0, errors.New("empty slice for median calculation")
	}

	// Create a copy to avoid modifying the original slice passed by reference.
	sortedData := make([]float64, len(data))
	copy(sortedData, data)
	sort.Float64s(sortedData)

	n := len(sortedData)
	if n%2 == 1 {
		// Odd number of elements, return the middle one
		return sortedData[n/2], nil
	} else {
		// Even number of elements, return the average of the two middle ones
		mid1 := sortedData[n/2-1]
		mid2 := sortedData[n/2]
		return (mid1 + mid2) / 2.0, nil
	}
}

// Mean calculates the arithmetic mean of a slice of float64 values.
//
// Parameters:
//
//	data: The slice of float64 for which to calculate the mean.
//
// Returns:
//
//	float64: The calculated mean.
//	error:   An error if the input slice is empty (errors.New("empty slice for mean calculation")).
func Mean(data []float64) (float64, error) {
	if len(data) == 0 {
		return 0.0, errors.New("empty slice for mean calculation")
	}
	sum := 0.0
	for _, val := range data {
		sum += val
	}
	return sum / float64(len(data)), nil
}

// Variance calculates the population variance of a slice of float64 values.
// It uses the formula: Sum((x - mean)^2) / N.
//
// Parameters:
//
//	data: The slice of float64 for which to calculate the variance.
//
// Returns:
//
//	float64: The calculated population variance.
//	error:   An error if the input slice is empty or contains only one element
//	         (errors.New("not enough points for variance calculation")).
func Variance(data []float64) (float64, error) {
	if len(data) < 1 { // Variance needs at least one point for mean, but usually more. NumPy uses N for N=1.
		return 0.0, errors.New("not enough points for variance calculation")
	}

	mean, err := Mean(data)
	if err != nil {
		return 0.0, err // Should not happen if len(data) > 0
	}

	sumSquaredDiff := 0.0
	for _, val := range data {
		diff := val - mean
		sumSquaredDiff += diff * diff
	}
	return sumSquaredDiff / float64(len(data)), nil
}
