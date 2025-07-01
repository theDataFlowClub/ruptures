// Package exceptions provides custom error types used throughout the ruptures library.
// These errors help to signal specific conditions that prevent operations from completing successfully,
// such as insufficient data for calculations or invalid segmentation parameters.
package exceptions

import "errors"

// ErrNotEnoughPoints is an error returned when there are insufficient data points
// to perform a required calculation, typically within a cost function or an algorithm
// that expects a minimum segment size.
var ErrNotEnoughPoints = errors.New("not enough points to calculate cost or perform operation")

// ErrBadSegmentationParameters is an error returned when a segmentation is not possible
// given the parameters, typically caught by preliminary sanity checks.
var ErrBadSegmentationParameters = errors.New("segmentation not possible given the parameters")

// ErrSegmentOutOfBounds is an error returned when the provided segment indices (start, end)
// are out of the valid range of the signal, or when start >= end (invalid segment definition).
var ErrSegmentOutOfBounds = errors.New("segment indices out of bounds or invalid") // NUEVO ERROR EXPORTADO

var ErrInvalidSignal = errors.New("rupture: invalid signal (nil or empty)")
