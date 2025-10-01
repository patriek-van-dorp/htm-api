package sdr

import (
	"errors"
	"fmt"
	"sort"
)

// SDR represents a Sparsely Distributed Representation as defined by HTM theory
type SDR struct {
	width      int     // Total number of bits in the representation
	activeBits []int   // Indices of active (1) bits, maintained in sorted order
	sparsity   float64 // Cached sparsity calculation
}

// NewSDR creates a new SDR with specified width and active bit indices
func NewSDR(width int, activeBits []int) (*SDR, error) {
	if width <= 0 {
		return nil, errors.New("SDR width must be positive")
	}

	// Validate and sort active bits
	if err := validateActiveBits(activeBits, width); err != nil {
		return nil, err
	}

	// Create a copy and sort to ensure immutability and order
	sortedActiveBits := make([]int, len(activeBits))
	copy(sortedActiveBits, activeBits)
	sort.Ints(sortedActiveBits)

	// Remove duplicates
	sortedActiveBits = removeDuplicates(sortedActiveBits)

	sparsity := float64(len(sortedActiveBits)) / float64(width)

	return &SDR{
		width:      width,
		activeBits: sortedActiveBits,
		sparsity:   sparsity,
	}, nil
}

// NewEmptySDR creates an empty SDR (for silent failure mode)
func NewEmptySDR(width int) (*SDR, error) {
	if width <= 0 {
		return nil, errors.New("SDR width must be positive")
	}

	return &SDR{
		width:      width,
		activeBits: []int{},
		sparsity:   0.0,
	}, nil
}

// Width returns the total number of bits in the representation
func (s *SDR) Width() int {
	return s.width
}

// ActiveBits returns indices of active (1) bits in sorted order
func (s *SDR) ActiveBits() []int {
	// Return a copy to maintain immutability
	result := make([]int, len(s.activeBits))
	copy(result, s.activeBits)
	return result
}

// Sparsity returns the percentage of active bits (0.0-1.0)
func (s *SDR) Sparsity() float64 {
	return s.sparsity
}

// IsActive returns true if the bit at given index is active
func (s *SDR) IsActive(index int) bool {
	if index < 0 || index >= s.width {
		return false
	}

	// Binary search since activeBits is sorted
	return binarySearch(s.activeBits, index)
}

// Overlap calculates the number of shared active bits with another SDR
func (s *SDR) Overlap(other *SDR) int {
	if s.width != other.width {
		return 0 // Different widths have no meaningful overlap
	}

	// Use two-pointer technique on sorted arrays
	i, j := 0, 0
	overlap := 0

	for i < len(s.activeBits) && j < len(other.activeBits) {
		if s.activeBits[i] == other.activeBits[j] {
			overlap++
			i++
			j++
		} else if s.activeBits[i] < other.activeBits[j] {
			i++
		} else {
			j++
		}
	}

	return overlap
}

// Similarity returns normalized overlap (0.0-1.0) with another SDR
func (s *SDR) Similarity(other *SDR) float64 {
	if s.width != other.width {
		return 0.0 // Different widths have no meaningful similarity
	}

	if len(s.activeBits) == 0 || len(other.activeBits) == 0 {
		return 0.0 // Empty SDRs have no similarity
	}

	overlap := s.Overlap(other)

	// Use the minimum active bits count for normalization
	// This is consistent with HTM theory for similarity calculation
	minActiveBits := len(s.activeBits)
	if len(other.activeBits) < minActiveBits {
		minActiveBits = len(other.activeBits)
	}

	return float64(overlap) / float64(minActiveBits)
}

// String returns a string representation of the SDR for debugging
func (s *SDR) String() string {
	return fmt.Sprintf("SDR(width=%d, active=%d, sparsity=%.3f, bits=%v)",
		s.width, len(s.activeBits), s.sparsity, s.activeBits)
}

// ValidateHTMCompliance checks if the SDR meets HTM theory requirements
func (s *SDR) ValidateHTMCompliance() error {
	if s.sparsity < 0.01 {
		return fmt.Errorf("sparsity %.3f below HTM minimum of 1%%", s.sparsity)
	}
	if s.sparsity > 0.10 {
		return fmt.Errorf("sparsity %.3f above HTM maximum of 10%%", s.sparsity)
	}
	return nil
}

// Helper functions

func validateActiveBits(activeBits []int, width int) error {
	for _, bit := range activeBits {
		if bit < 0 || bit >= width {
			return fmt.Errorf("active bit index %d out of range [0, %d)", bit, width)
		}
	}
	return nil
}

func removeDuplicates(sorted []int) []int {
	if len(sorted) <= 1 {
		return sorted
	}

	j := 0
	for i := 1; i < len(sorted); i++ {
		if sorted[i] != sorted[j] {
			j++
			sorted[j] = sorted[i]
		}
	}
	return sorted[:j+1]
}

func binarySearch(sorted []int, target int) bool {
	left, right := 0, len(sorted)-1

	for left <= right {
		mid := left + (right-left)/2
		if sorted[mid] == target {
			return true
		} else if sorted[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return false
}
