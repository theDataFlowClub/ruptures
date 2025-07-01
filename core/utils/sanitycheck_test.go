package utils_test

import (
	"testing"

	"github.com/theDataFlowClub/ruptures/core/utils"
)

func TestSanityCheck(t *testing.T) {
	testCases := []struct {
		name     string
		nSamples int
		nBkps    int
		jump     int
		minSize  int
		expected bool
	}{
		{
			name:     "Valid_SimpleCase",
			nSamples: 100,
			nBkps:    1,
			jump:     1,
			minSize:  10,
			expected: true, // Segment 0-10, 10-100 (90 points) is possible
		},
		{
			name:     "Valid_MultipleBkps",
			nSamples: 100,
			nBkps:    3,
			jump:     1,
			minSize:  10,
			expected: true, // 3 bkps, 4 segments. Smallest possible: 3*10 + 10 = 40 points
		},
		{
			name:     "Invalid_TooManyBkps",
			nSamples: 50,
			nBkps:    5,
			jump:     1,
			minSize:  10,
			expected: false, // 5 bkps (6 segments) * 10 minSize = 60 points > 50 samples
		},
		{
			name:     "Invalid_TooManyBkpsWithJump",
			nSamples: 100,
			nBkps:    10, // Max admissible bkps for jump 10 is 100/10 = 10
			jump:     10,
			minSize:  5,
			expected: false, // 10 bkps * ceil(5/10)*10 + 5 = 10*1*10 + 5 = 105 > 100
		},
		{
			name:     "Valid_WithJumpConstraint",
			nSamples: 100,
			nBkps:    3,
			jump:     10,
			minSize:  10,
			expected: true, // 3 bkps, 4 segments. Smallest: 3*ceil(10/10)*10 + 10 = 3*10 + 10 = 40. OK.
		},
		{
			name:     "Invalid_MinSizeTooLarge",
			nSamples: 20,
			nBkps:    1,
			jump:     1,
			minSize:  15,
			expected: false, // 1 bkp, 2 segments. Smallest: 1*15 + 15 = 30 > 20
		},
		{
			name:     "EdgeCase_ZeroBkps",
			nSamples: 10,
			nBkps:    0,
			jump:     1,
			minSize:  1,
			expected: true, // 0 bkps, 1 segment (10 points). OK.
		},
		{
			name:     "EdgeCase_nSamplesLessThanMinSize",
			nSamples: 5,
			nBkps:    0,
			jump:     1,
			minSize:  10,
			expected: false, // 0 bkps, but minSize (10) > nSamples (5)
		},
		{
			name:     "EdgeCase_MinimumPossible",
			nSamples: 2,
			nBkps:    0,
			jump:     1,
			minSize:  2,
			expected: true, // 1 segment of size 2.
		},
		{
			name:     "EdgeCase_JumpGreaterThanMinSize",
			nSamples: 100,
			nBkps:    1,
			jump:     20,
			minSize:  10,
			expected: true, // 1 bkp, 2 segments. minPoints = 1*ceil(10/20)*20 + 10 = 1*1*20 + 10 = 30. OK.
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.SanityCheck(tc.nSamples, tc.nBkps, tc.jump, tc.minSize)
			if result != tc.expected {
				t.Errorf("SanityCheck(nSamples:%d, nBkps:%d, jump:%d, minSize:%d) = %v; want %v",
					tc.nSamples, tc.nBkps, tc.jump, tc.minSize, result, tc.expected)
			}
		})
	}
}
