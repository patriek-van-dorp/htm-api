package sdr

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

// SDR represents a Sparse Distributed Representation as defined by HTM theory
// This is the enhanced version for cortical operations
type SDR struct {
	Width      int     `json:"width"`
	ActiveBits []int   `json:"active_bits"`
	Sparsity   float64 `json:"sparsity"`
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
		Width:      width,
		ActiveBits: sortedActiveBits,
		Sparsity:   sparsity,
	}, nil
}

// NewEmptySDR creates an empty SDR
func NewEmptySDR(width int) (*SDR, error) {
	if width <= 0 {
		return nil, errors.New("SDR width must be positive")
	}

	return &SDR{
		Width:      width,
		ActiveBits: []int{},
		Sparsity:   0.0,
	}, nil
}

// NewSDRFromPattern creates SDR from binary pattern (for testing/debugging)
func NewSDRFromPattern(pattern []bool) (*SDR, error) {
	width := len(pattern)
	if width == 0 {
		return nil, errors.New("pattern cannot be empty")
	}

	var activeBits []int
	for i, bit := range pattern {
		if bit {
			activeBits = append(activeBits, i)
		}
	}

	return NewSDR(width, activeBits)
}

// IsActive returns true if the bit at given index is active
func (s *SDR) IsActive(index int) bool {
	if index < 0 || index >= s.Width {
		return false
	}

	// Binary search since ActiveBits is sorted
	return binarySearch(s.ActiveBits, index)
}

// Overlap calculates the number of shared active bits with another SDR
func (s *SDR) Overlap(other *SDR) int {
	if s.Width != other.Width {
		return 0 // Different widths have no meaningful overlap
	}

	// Use two-pointer technique on sorted arrays
	i, j := 0, 0
	overlap := 0

	for i < len(s.ActiveBits) && j < len(other.ActiveBits) {
		if s.ActiveBits[i] == other.ActiveBits[j] {
			overlap++
			i++
			j++
		} else if s.ActiveBits[i] < other.ActiveBits[j] {
			i++
		} else {
			j++
		}
	}

	return overlap
}

// OverlapRatio returns normalized overlap (0.0-1.0) with another SDR
func (s *SDR) OverlapRatio(other *SDR) float64 {
	if s.Width != other.Width {
		return 0.0 // Different widths have no meaningful similarity
	}

	if len(s.ActiveBits) == 0 || len(other.ActiveBits) == 0 {
		return 0.0 // Empty SDRs have no similarity
	}

	overlap := s.Overlap(other)

	// Use the minimum active bits count for normalization
	// This is consistent with HTM theory for similarity calculation
	minActiveBits := len(s.ActiveBits)
	if len(other.ActiveBits) < minActiveBits {
		minActiveBits = len(other.ActiveBits)
	}

	return float64(overlap) / float64(minActiveBits)
}

// JaccardSimilarity calculates Jaccard index (intersection/union)
func (s *SDR) JaccardSimilarity(other *SDR) float64 {
	if s.Width != other.Width {
		return 0.0
	}

	if len(s.ActiveBits) == 0 && len(other.ActiveBits) == 0 {
		return 1.0 // Both empty, considered identical
	}

	intersection := float64(s.Overlap(other))
	union := float64(len(s.ActiveBits)+len(other.ActiveBits)) - intersection

	if union == 0 {
		return 0.0
	}

	return intersection / union
}

// CosineSimilarity calculates cosine similarity between SDRs
func (s *SDR) CosineSimilarity(other *SDR) float64 {
	if s.Width != other.Width {
		return 0.0
	}

	if len(s.ActiveBits) == 0 || len(other.ActiveBits) == 0 {
		return 0.0
	}

	// For binary vectors, cosine similarity = overlap / sqrt(|A| * |B|)
	overlap := float64(s.Overlap(other))
	magnitude := math.Sqrt(float64(len(s.ActiveBits)) * float64(len(other.ActiveBits)))

	return overlap / magnitude
}

// Union creates a new SDR with bits active in either SDR
func (s *SDR) Union(other *SDR) (*SDR, error) {
	if s.Width != other.Width {
		return nil, errors.New("SDRs must have same width for union")
	}

	// Merge sorted arrays
	result := make([]int, 0, len(s.ActiveBits)+len(other.ActiveBits))
	i, j := 0, 0

	for i < len(s.ActiveBits) && j < len(other.ActiveBits) {
		if s.ActiveBits[i] == other.ActiveBits[j] {
			result = append(result, s.ActiveBits[i])
			i++
			j++
		} else if s.ActiveBits[i] < other.ActiveBits[j] {
			result = append(result, s.ActiveBits[i])
			i++
		} else {
			result = append(result, other.ActiveBits[j])
			j++
		}
	}

	// Add remaining elements
	for i < len(s.ActiveBits) {
		result = append(result, s.ActiveBits[i])
		i++
	}
	for j < len(other.ActiveBits) {
		result = append(result, other.ActiveBits[j])
		j++
	}

	return NewSDR(s.Width, result)
}

// Intersection creates a new SDR with bits active in both SDRs
func (s *SDR) Intersection(other *SDR) (*SDR, error) {
	if s.Width != other.Width {
		return nil, errors.New("SDRs must have same width for intersection")
	}

	// Find intersection of sorted arrays
	result := make([]int, 0, min(len(s.ActiveBits), len(other.ActiveBits)))
	i, j := 0, 0

	for i < len(s.ActiveBits) && j < len(other.ActiveBits) {
		if s.ActiveBits[i] == other.ActiveBits[j] {
			result = append(result, s.ActiveBits[i])
			i++
			j++
		} else if s.ActiveBits[i] < other.ActiveBits[j] {
			i++
		} else {
			j++
		}
	}

	return NewSDR(s.Width, result)
}

// ToBinaryArray converts SDR to binary array representation
func (s *SDR) ToBinaryArray() []bool {
	result := make([]bool, s.Width)
	for _, bit := range s.ActiveBits {
		result[bit] = true
	}
	return result
}

// String returns a string representation of the SDR for debugging
func (s *SDR) String() string {
	return fmt.Sprintf("SDR(width=%d, active=%d, sparsity=%.3f, bits=%v)",
		s.Width, len(s.ActiveBits), s.Sparsity, s.ActiveBits)
}

// ValidateHTMCompliance checks if the SDR meets HTM theory requirements
func (s *SDR) ValidateHTMCompliance() error {
	if s.Sparsity < 0.01 {
		return fmt.Errorf("sparsity %.3f below HTM minimum of 1%%", s.Sparsity)
	}
	if s.Sparsity > 0.10 {
		return fmt.Errorf("sparsity %.3f above HTM maximum of 10%%", s.Sparsity)
	}
	return nil
}

// ValidateSpatialPoolerCompliance checks if SDR meets spatial pooler requirements (2-5%)
func (s *SDR) ValidateSpatialPoolerCompliance() error {
	if s.Sparsity < 0.02 {
		return fmt.Errorf("sparsity %.3f below spatial pooler minimum of 2%%", s.Sparsity)
	}
	if s.Sparsity > 0.05 {
		return fmt.Errorf("sparsity %.3f above spatial pooler maximum of 5%%", s.Sparsity)
	}
	return nil
}

// Clone creates a deep copy of the SDR
func (s *SDR) Clone() *SDR {
	activeBits := make([]int, len(s.ActiveBits))
	copy(activeBits, s.ActiveBits)

	return &SDR{
		Width:      s.Width,
		ActiveBits: activeBits,
		Sparsity:   s.Sparsity,
	}
}

// IsEmpty returns true if SDR has no active bits
func (s *SDR) IsEmpty() bool {
	return len(s.ActiveBits) == 0
}

// IsSimilarTo returns true if overlap ratio meets the threshold
func (s *SDR) IsSimilarTo(other *SDR, threshold float64) bool {
	return s.OverlapRatio(other) >= threshold
}

// IsDistinctFrom returns true if overlap ratio is below the threshold
func (s *SDR) IsDistinctFrom(other *SDR, threshold float64) bool {
	return s.OverlapRatio(other) <= threshold
}

// SpatialPoolerOperations contains specialized operations for spatial pooling
type SpatialPoolerOperations struct{}

// NewSpatialPoolerOperations creates a new spatial pooler operations helper
func NewSpatialPoolerOperations() *SpatialPoolerOperations {
	return &SpatialPoolerOperations{}
}

// NormalizeSparsity ensures SDR meets exact sparsity target for spatial pooler output
func (spo *SpatialPoolerOperations) NormalizeSparsity(sdr *SDR, targetSparsity float64) (*SDR, error) {
	if targetSparsity < 0.02 || targetSparsity > 0.05 {
		return nil, fmt.Errorf("target sparsity %.4f outside spatial pooler range [0.02, 0.05]", targetSparsity)
	}

	targetActiveBits := int(float64(sdr.Width) * targetSparsity)
	if targetActiveBits == 0 {
		targetActiveBits = 1 // Ensure at least 1 active bit
	}

	currentActiveBits := len(sdr.ActiveBits)

	if currentActiveBits == targetActiveBits {
		return sdr.Clone(), nil
	}

	if currentActiveBits > targetActiveBits {
		// Randomly select subset (deterministic based on input pattern)
		selectedBits := make([]int, targetActiveBits)
		step := float64(currentActiveBits) / float64(targetActiveBits)
		for i := 0; i < targetActiveBits; i++ {
			selectedBits[i] = sdr.ActiveBits[int(float64(i)*step)]
		}
		return NewSDR(sdr.Width, selectedBits)
	} else {
		// Need to add more bits - this should be handled by spatial pooler algorithm
		// For normalization, we return error
		return nil, fmt.Errorf("cannot normalize: SDR has %d bits but needs %d bits (expansion not supported in normalization)",
			currentActiveBits, targetActiveBits)
	}
}

// CalculateSemanticContinuity measures how well spatial pooler preserves semantic relationships
func (spo *SpatialPoolerOperations) CalculateSemanticContinuity(
	inputSDRs []*SDR,
	outputSDRs []*SDR,
	similarThreshold float64,
	differentThreshold float64) (float64, error) {

	if len(inputSDRs) != len(outputSDRs) {
		return 0.0, errors.New("input and output SDR arrays must have same length")
	}

	if len(inputSDRs) < 2 {
		return 1.0, nil // Perfect continuity for single SDR
	}

	totalComparisons := 0
	continuityViolations := 0

	for i := 0; i < len(inputSDRs); i++ {
		for j := i + 1; j < len(inputSDRs); j++ {
			inputSimilarity := inputSDRs[i].OverlapRatio(inputSDRs[j])
			outputSimilarity := outputSDRs[i].OverlapRatio(outputSDRs[j])

			totalComparisons++

			// Check semantic continuity violations
			if inputSimilarity >= similarThreshold && outputSimilarity <= differentThreshold {
				// Similar inputs became different outputs - violation
				continuityViolations++
			} else if inputSimilarity <= differentThreshold && outputSimilarity >= similarThreshold {
				// Different inputs became similar outputs - violation
				continuityViolations++
			}
		}
	}

	if totalComparisons == 0 {
		return 1.0, nil
	}

	continuityScore := 1.0 - (float64(continuityViolations) / float64(totalComparisons))
	return continuityScore, nil
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
