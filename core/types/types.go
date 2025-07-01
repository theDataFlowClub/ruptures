// Package types provides fundamental data type definitions used across the ruptures library.
// These types enhance code clarity and ensure type safety for signals, vectors, and breakpoint lists.
// They act as aliases for built-in Go types, making the intent of various functions and interfaces
// more explicit, similar to how domain-specific types might be used in other languages.
package types

// Matrix represents a 2D dataset, typically a signal with multiple features.
// It is structured as [][]float64, where the first dimension corresponds to
// the number of samples (n_samples) and the second to the number of features (n_features).
type Matrix = [][]float64

// Vector represents a 1D dataset, typically a univariate signal.
// It is structured as []float64, where the length corresponds to the
// number of samples (n_samples).
type Vector = []float64

// Signal is a semantic alias for Matrix. It can be used interchangeably with Matrix
// when the context specifically refers to the input data as a "signal",
// emphasizing its role in signal processing.
type Signal = [][]float64

// Breakpoints represents a sorted list of integers indicating the indices
// where change points (breakpoints) are detected within a signal.
// By convention, the last element in a Breakpoints slice often corresponds
// to the total number of samples (n_samples), marking the end of the last segment.
type Breakpoints = []int
