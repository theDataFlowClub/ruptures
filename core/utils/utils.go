// Package utils provides miscellaneous helper functions used across the ruptures library.
// These functions perform common tasks such as data reshaping, argument validation,
// and utility operations that support the core change point detection algorithms.
package utils

import "math"

// Pair represents a generic ordered pair of two elements of potentially different types.
// It serves as a common building block in Go, especially when translating concepts
// from languages like Python where "tuples" are native.
//
// Using a struct for pairs, rather than a slice or array, offers:
//   - Enhanced readability: Access elements by name (e.g., `p.First`, `p.Second`)
//     instead of less descriptive indices (e.g., `p[0]`, `p[1]`).
//   - Type safety and clarity: Explicitly defines the types of both elements
//     at compile time, preventing common runtime errors.
//   - Extensibility: Easily supports heterogeneous pairs (e.g., `Pair[int, float64]`).
//
// Example:
//
//	myPair := Pair[int, string]{First: 10, Second: "hello"}
//	value1 := myPair.First // value1 is an int
//	value2 := myPair.Second // value2 is a string
//
// .
type Pair[T1, T2 any] struct {
	First  T1
	Second T2
}

// Pairwise returns an iterator-like slice of consecutive, non-overlapping pairs from the input slice.
// For a given slice S = [s0, s1, s2, s3, ...], Pairwise generates a new slice of pairs:
// [(s0,s1), (s1,s2), (s2,s3), ...].
//
// This function is particularly useful for operations that require inspecting adjacent
// elements, such as calculating differences or identifying local patterns.
//
// If the input slice has fewer than two elements, an empty (nil) slice of pairs is returned.
//
// Parameters:
//
//	slice: The input slice of integers from which pairs will be generated.
//
// Returns:
//
//	A slice of Pair[int, int] containing consecutive pairs from the input.
//	Returns nil if the input slice has less than 2 elements.
//
// .
func Pairwise(slice []int) []Pair[int, int] {
	if len(slice) < 2 {
		return nil
	}
	result := make([]Pair[int, int], 0, len(slice)-1)
	for i := 0; i < len(slice)-1; i++ {
		result = append(result, Pair[int, int]{slice[i], slice[i+1]})
	}
	return result
}

// Unzip separates a slice of Pair[int, int] into two distinct slices.
// It effectively reverses the operation of a "zip" function, where the first
// elements of all pairs are collected into one slice, and the second elements
// into another.
//
// Parameters:
//
//	pairs: The input slice of Pair[int, int] to be unzipped.
//
// Returns:
//
//	Two slices of integers:
//	  - The first slice contains all 'First' elements from the input pairs, in order.
//	  - The second slice contains all 'Second' elements from the input pairs, in order.
//
// .
func Unzip(pairs []Pair[int, int]) ([]int, []int) {
	a := make([]int, len(pairs))
	b := make([]int, len(pairs))
	for i, p := range pairs {
		a[i] = p.First
		b[i] = p.Second
	}
	return a, b
}

// SanityCheck validates if a proposed partitioning of a signal is mathematically possible
// given a set of segmentation parameters. This function ensures that there are enough
// data points to accommodate the specified number of breakpoints, minimum segment size,
// and jump constraint.
//
// It performs two main checks:
//  1. Ensures that the number of requested breakpoints (`nBkps`) does not exceed the
//     maximum possible number of "admissible" breakpoints given the `jump` constraint
//     (`n_samples / jump`).
//  2. Verifies that the total required points for all segments (including the minimum
//     size for each, and accounting for the `jump` constraint) do not exceed the
//     total number of samples (`nSamples`). This calculation ensures that even with
//     the smallest possible segments, the configuration fits within the signal length.
//
// This preliminary check helps to prevent invalid or unfeasible segmentation
// configurations from being passed to the more complex change point detection algorithms,
// thereby improving robustness and efficiency.
//
// Parameters:
//
//	nSamples (int): The total number of points in the signal.
//	nBkps (int): The desired number of breakpoints (which implies `nBkps + 1` segments).
//	jump (int): The step size between admissible start indices of segments.
//				A start index of segment must be a multiple of `jump`.
//	minSize (int): The minimum allowable size for any segment.
//
// Returns:
//
//	bool: True if a valid breakpoint configuration is possible with the given parameters,
//	      False otherwise.
//
// .
func SanityCheck(nSamples, nBkps, jump, minSize int) bool {
	// nAdmissibleBkps calculates the maximum number of breakpoints that can be placed
	// given the total number of samples and the jump constraint.
	nAdmissibleBkps := nSamples / jump

	// Check 1: Ensure the requested number of breakpoints doesn't exceed
	// the maximum possible admissible breakpoints.
	if nBkps > nAdmissibleBkps {
		return false
	}

	// Check 2: Calculate the minimum total points required for the specified
	// number of breakpoints and minimum segment sizes, considering the jump.
	// math.Ceil is used to ensure that even partial `minSize / jump` segments
	// are rounded up to the next full jump step.
	// The formula `nBkps * ceil(min_size / jump) * jump + min_size` ensures that
	// all `nBkps` segments can at least accommodate `minSize` points,
	// aligned with `jump` boundaries, plus the very last segment.
	requiredPoints := nBkps*int(math.Ceil(float64(minSize)/float64(jump)))*jump + minSize
	if requiredPoints > nSamples {
		return false
	}

	return true
}
