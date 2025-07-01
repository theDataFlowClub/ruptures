package utils_test // Note: _test suffix for package name when testing external functions

import (
	"reflect" // Used for deep comparison of slices and structs
	"testing"

	"github.com/theDataFlowClub/ruptures/core/utils" // Import the package being tested
)

func TestPairwise(t *testing.T) {
	// Test case 1: Empty slice
	t.Run("EmptySlice", func(t *testing.T) {
		input := []int{}
		expected := []utils.Pair[int, int](nil) // nil slice for expected empty result
		result := utils.Pairwise(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Pairwise(%v) = %v; want %v", input, result, expected)
		}
	})

	// Test case 2: Single element slice
	t.Run("SingleElementSlice", func(t *testing.T) {
		input := []int{1}
		expected := []utils.Pair[int, int](nil)
		result := utils.Pairwise(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Pairwise(%v) = %v; want %v", input, result, expected)
		}
	})

	// Test case 3: Standard slice
	t.Run("StandardSlice", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		expected := []utils.Pair[int, int]{
			{First: 1, Second: 2},
			{First: 2, Second: 3},
			{First: 3, Second: 4},
		}
		result := utils.Pairwise(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Pairwise(%v) = %v; want %v", input, result, expected)
		}
	})

	// Test case 4: Slice with duplicate elements
	t.Run("DuplicateElements", func(t *testing.T) {
		input := []int{5, 5, 6, 6}
		expected := []utils.Pair[int, int]{
			{First: 5, Second: 5},
			{First: 5, Second: 6},
			{First: 6, Second: 6},
		}
		result := utils.Pairwise(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Pairwise(%v) = %v; want %v", input, result, expected)
		}
	})
}

func TestUnzip(t *testing.T) {
	// Test case 1: Empty slice of pairs
	t.Run("EmptyPairs", func(t *testing.T) {
		input := []utils.Pair[int, int]{}
		expectedA := []int{}
		expectedB := []int{}
		resultA, resultB := utils.Unzip(input)
		if !reflect.DeepEqual(resultA, expectedA) || !reflect.DeepEqual(resultB, expectedB) {
			t.Errorf("Unzip(%v) = (%v, %v); want (%v, %v)", input, resultA, resultB, expectedA, expectedB)
		}
	})

	// Test case 2: Standard slice of pairs
	t.Run("StandardPairs", func(t *testing.T) {
		input := []utils.Pair[int, int]{
			{First: 1, Second: 10},
			{First: 2, Second: 20},
			{First: 3, Second: 30},
		}
		expectedA := []int{1, 2, 3}
		expectedB := []int{10, 20, 30}
		resultA, resultB := utils.Unzip(input)
		if !reflect.DeepEqual(resultA, expectedA) || !reflect.DeepEqual(resultB, expectedB) {
			t.Errorf("Unzip(%v) = (%v, %v); want (%v, %v)", input, resultA, resultB, expectedA, expectedB)
		}
	})

	// Test case 3: Pairs with duplicate values
	t.Run("DuplicateValues", func(t *testing.T) {
		input := []utils.Pair[int, int]{
			{First: 7, Second: 7},
			{First: 8, Second: 9},
		}
		expectedA := []int{7, 8}
		expectedB := []int{7, 9}
		resultA, resultB := utils.Unzip(input)
		if !reflect.DeepEqual(resultA, expectedA) || !reflect.DeepEqual(resultB, expectedB) {
			t.Errorf("Unzip(%v) = (%v, %v); want (%v, %v)", input, resultA, resultB, expectedA, expectedB)
		}
	})
}
