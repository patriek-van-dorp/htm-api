package sdr

import (
	"errors"
	"fmt"
	"math"
)

// SparsityManager handles SDR sparsity calculations and validation
type SparsityManager struct {
	targetSparsity float64
	minSparsity    float64
	maxSparsity    float64
}

// NewSparsityManager creates a new sparsity manager with HTM-compliant defaults
func NewSparsityManager(targetSparsity float64) (*SparsityManager, error) {
	if targetSparsity < 0.01 || targetSparsity > 0.10 {
		return nil, fmt.Errorf("target sparsity %.3f outside HTM range [0.01, 0.10]", targetSparsity)
	}

	return &SparsityManager{
		targetSparsity: targetSparsity,
		minSparsity:    0.01, // HTM minimum: 1%
		maxSparsity:    0.10, // HTM maximum: 10%
	}, nil
}

// CalculateActiveBitsCount returns the number of active bits for given width and target sparsity
func (sm *SparsityManager) CalculateActiveBitsCount(width int) int {
	if width <= 0 {
		return 0
	}

	activeBits := int(math.Round(float64(width) * sm.targetSparsity))

	// Ensure at least 1 active bit for non-zero sparsity
	if activeBits == 0 && sm.targetSparsity > 0 {
		activeBits = 1
	}

	// Ensure we don't exceed width
	if activeBits > width {
		activeBits = width
	}

	return activeBits
}

// ValidateSparsity checks if the given sparsity meets HTM compliance
func (sm *SparsityManager) ValidateSparsity(actualSparsity float64) error {
	if actualSparsity < sm.minSparsity {
		return fmt.Errorf("sparsity %.3f below HTM minimum %.3f", actualSparsity, sm.minSparsity)
	}
	if actualSparsity > sm.maxSparsity {
		return fmt.Errorf("sparsity %.3f above HTM maximum %.3f", actualSparsity, sm.maxSparsity)
	}
	return nil
}

// ValidateSDRSparsity checks if an SDR meets HTM sparsity requirements
func (sm *SparsityManager) ValidateSDRSparsity(sdr *SDR) error {
	if sdr == nil {
		return errors.New("SDR cannot be nil")
	}
	return sm.ValidateSparsity(sdr.Sparsity())
}

// AdjustActiveBitsForWidth adjusts the number of active bits to maintain target sparsity for a given width
func (sm *SparsityManager) AdjustActiveBitsForWidth(currentActiveBits, newWidth int) int {
	if newWidth <= 0 {
		return 0
	}

	// Calculate what the active bits should be for the new width
	targetActiveBits := sm.CalculateActiveBitsCount(newWidth)
	return targetActiveBits
}

// GetTargetSparsity returns the configured target sparsity
func (sm *SparsityManager) GetTargetSparsity() float64 {
	return sm.targetSparsity
}

// GetSparsityRange returns the valid sparsity range for HTM compliance
func (sm *SparsityManager) GetSparsityRange() (min, max float64) {
	return sm.minSparsity, sm.maxSparsity
}

// IsSparsityInRange checks if sparsity is within HTM-compliant range
func (sm *SparsityManager) IsSparsityInRange(sparsity float64) bool {
	return sparsity >= sm.minSparsity && sparsity <= sm.maxSparsity
}

// CalculateSparsityForBits calculates sparsity given active bits count and width
func CalculateSparsityForBits(activeBitsCount, width int) float64 {
	if width <= 0 {
		return 0.0
	}
	return float64(activeBitsCount) / float64(width)
}

// SparsityStatistics provides statistical information about sparsity distribution
type SparsityStatistics struct {
	Mean     float64
	Min      float64
	Max      float64
	StdDev   float64
	Count    int
	InRange  int
	OutRange int
}

// SparsityAnalyzer analyzes sparsity patterns across multiple SDRs
type SparsityAnalyzer struct {
	sparsities []float64
	manager    *SparsityManager
}

// NewSparsityAnalyzer creates a new sparsity analyzer
func NewSparsityAnalyzer(manager *SparsityManager) *SparsityAnalyzer {
	return &SparsityAnalyzer{
		sparsities: make([]float64, 0),
		manager:    manager,
	}
}

// AddSDR adds an SDR's sparsity to the analysis
func (sa *SparsityAnalyzer) AddSDR(sdr *SDR) {
	if sdr != nil {
		sa.sparsities = append(sa.sparsities, sdr.Sparsity())
	}
}

// GetStatistics calculates and returns sparsity statistics
func (sa *SparsityAnalyzer) GetStatistics() SparsityStatistics {
	if len(sa.sparsities) == 0 {
		return SparsityStatistics{}
	}

	stats := SparsityStatistics{Count: len(sa.sparsities)}

	// Calculate min, max, and sum
	sum := 0.0
	stats.Min = sa.sparsities[0]
	stats.Max = sa.sparsities[0]

	for _, sparsity := range sa.sparsities {
		sum += sparsity
		if sparsity < stats.Min {
			stats.Min = sparsity
		}
		if sparsity > stats.Max {
			stats.Max = sparsity
		}

		// Count in-range vs out-of-range
		if sa.manager.IsSparsityInRange(sparsity) {
			stats.InRange++
		} else {
			stats.OutRange++
		}
	}

	// Calculate mean
	stats.Mean = sum / float64(len(sa.sparsities))

	// Calculate standard deviation
	sumSquaredDiff := 0.0
	for _, sparsity := range sa.sparsities {
		diff := sparsity - stats.Mean
		sumSquaredDiff += diff * diff
	}
	stats.StdDev = math.Sqrt(sumSquaredDiff / float64(len(sa.sparsities)))

	return stats
}

// Reset clears all collected sparsity data
func (sa *SparsityAnalyzer) Reset() {
	sa.sparsities = sa.sparsities[:0]
}
