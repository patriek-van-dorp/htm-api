package sdr

import (
	"fmt"
	"math"
)

// RepresentationAnalyzer provides analysis tools for SDR representations
type RepresentationAnalyzer struct{}

// NewRepresentationAnalyzer creates a new representation analyzer
func NewRepresentationAnalyzer() *RepresentationAnalyzer {
	return &RepresentationAnalyzer{}
}

// AnalyzeDistribution analyzes the distribution of active bits in an SDR
func (ra *RepresentationAnalyzer) AnalyzeDistribution(sdr *SDR) *DistributionAnalysis {
	if sdr.IsEmpty() {
		return &DistributionAnalysis{
			UniformityScore: 0.0,
			ClusteringIndex: 0.0,
			MaxGap:          sdr.Width,
			MinGap:          sdr.Width,
			AverageGap:      float64(sdr.Width),
		}
	}

	gaps := ra.calculateGaps(sdr)

	return &DistributionAnalysis{
		UniformityScore: ra.calculateUniformityScore(gaps, sdr.Width),
		ClusteringIndex: ra.calculateClusteringIndex(gaps),
		MaxGap:          ra.maxGap(gaps),
		MinGap:          ra.minGap(gaps),
		AverageGap:      ra.averageGap(gaps),
	}
}

// DistributionAnalysis contains metrics about SDR bit distribution
type DistributionAnalysis struct {
	UniformityScore float64 `json:"uniformity_score"` // 0.0 (clustered) to 1.0 (uniform)
	ClusteringIndex float64 `json:"clustering_index"` // Measure of bit clustering
	MaxGap          int     `json:"max_gap"`          // Largest gap between active bits
	MinGap          int     `json:"min_gap"`          // Smallest gap between active bits
	AverageGap      float64 `json:"average_gap"`      // Average gap size
}

// ValidateRepresentation performs comprehensive validation of SDR representation
func (ra *RepresentationAnalyzer) ValidateRepresentation(sdr *SDR) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Issues:  []string{},
	}

	// Check basic constraints
	if sdr.Width <= 0 {
		result.IsValid = false
		result.Issues = append(result.Issues, "SDR width must be positive")
	}

	// Check active bits are in valid range and sorted
	for i, bit := range sdr.ActiveBits {
		if bit < 0 || bit >= sdr.Width {
			result.IsValid = false
			result.Issues = append(result.Issues, fmt.Sprintf("active bit %d out of range [0, %d)", bit, sdr.Width))
		}
		if i > 0 && bit <= sdr.ActiveBits[i-1] {
			result.IsValid = false
			result.Issues = append(result.Issues, "active bits must be sorted and unique")
			break
		}
	}

	// Check sparsity calculation
	expectedSparsity := float64(len(sdr.ActiveBits)) / float64(sdr.Width)
	if math.Abs(sdr.Sparsity-expectedSparsity) > 0.000001 {
		result.IsValid = false
		result.Issues = append(result.Issues, fmt.Sprintf("sparsity mismatch: expected %.6f, got %.6f", expectedSparsity, sdr.Sparsity))
	}

	// HTM compliance checks
	if err := sdr.ValidateHTMCompliance(); err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("HTM compliance: %v", err))
	}

	// Spatial pooler compliance checks
	if err := sdr.ValidateSpatialPoolerCompliance(); err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("Spatial pooler compliance: %v", err))
	}

	return result
}

// ValidationResult contains the result of SDR validation
type ValidationResult struct {
	IsValid bool     `json:"is_valid"`
	Issues  []string `json:"issues"`
}

// CalculateCapacity estimates the representational capacity of SDR parameters
func (ra *RepresentationAnalyzer) CalculateCapacity(width, activeBits int) *CapacityAnalysis {
	if width <= 0 || activeBits <= 0 || activeBits > width {
		return &CapacityAnalysis{
			IsValid: false,
			Error:   "invalid parameters",
		}
	}

	// Calculate combinatorial capacity: C(width, activeBits)
	capacity := ra.combinations(width, activeBits)
	sparsity := float64(activeBits) / float64(width)

	return &CapacityAnalysis{
		IsValid:               true,
		CombinationalCapacity: capacity,
		Sparsity:              sparsity,
		Width:                 width,
		ActiveBits:            activeBits,
		RecommendedForHTM:     sparsity >= 0.01 && sparsity <= 0.10,
		RecommendedForSpatial: sparsity >= 0.02 && sparsity <= 0.05,
	}
}

// CapacityAnalysis contains representational capacity metrics
type CapacityAnalysis struct {
	IsValid               bool    `json:"is_valid"`
	Error                 string  `json:"error,omitempty"`
	CombinationalCapacity float64 `json:"combinational_capacity"`
	Sparsity              float64 `json:"sparsity"`
	Width                 int     `json:"width"`
	ActiveBits            int     `json:"active_bits"`
	RecommendedForHTM     bool    `json:"recommended_for_htm"`
	RecommendedForSpatial bool    `json:"recommended_for_spatial"`
}

// CompareRepresentations compares two SDRs across multiple dimensions
func (ra *RepresentationAnalyzer) CompareRepresentations(sdr1, sdr2 *SDR) *ComparisonResult {
	if sdr1.Width != sdr2.Width {
		return &ComparisonResult{
			IsComparable: false,
			Error:        "SDRs must have same width for comparison",
		}
	}

	overlap := sdr1.Overlap(sdr2)
	overlapRatio := sdr1.OverlapRatio(sdr2)
	jaccardSim := sdr1.JaccardSimilarity(sdr2)
	cosineSim := sdr1.CosineSimilarity(sdr2)

	// Calculate distribution similarity
	dist1 := ra.AnalyzeDistribution(sdr1)
	dist2 := ra.AnalyzeDistribution(sdr2)
	distributionSimilarity := 1.0 - math.Abs(dist1.UniformityScore-dist2.UniformityScore)

	return &ComparisonResult{
		IsComparable:           true,
		Overlap:                overlap,
		OverlapRatio:           overlapRatio,
		JaccardSimilarity:      jaccardSim,
		CosineSimilarity:       cosineSim,
		DistributionSimilarity: distributionSimilarity,
		SparsityDifference:     math.Abs(sdr1.Sparsity - sdr2.Sparsity),
		SizeComparison: map[string]int{
			"sdr1_active": len(sdr1.ActiveBits),
			"sdr2_active": len(sdr2.ActiveBits),
			"difference":  len(sdr1.ActiveBits) - len(sdr2.ActiveBits),
		},
	}
}

// ComparisonResult contains comprehensive comparison metrics
type ComparisonResult struct {
	IsComparable           bool           `json:"is_comparable"`
	Error                  string         `json:"error,omitempty"`
	Overlap                int            `json:"overlap"`
	OverlapRatio           float64        `json:"overlap_ratio"`
	JaccardSimilarity      float64        `json:"jaccard_similarity"`
	CosineSimilarity       float64        `json:"cosine_similarity"`
	DistributionSimilarity float64        `json:"distribution_similarity"`
	SparsityDifference     float64        `json:"sparsity_difference"`
	SizeComparison         map[string]int `json:"size_comparison"`
}

// Helper methods

func (ra *RepresentationAnalyzer) calculateGaps(sdr *SDR) []int {
	if len(sdr.ActiveBits) <= 1 {
		return []int{}
	}

	gaps := make([]int, len(sdr.ActiveBits)-1)
	for i := 1; i < len(sdr.ActiveBits); i++ {
		gaps[i-1] = sdr.ActiveBits[i] - sdr.ActiveBits[i-1] - 1
	}
	return gaps
}

func (ra *RepresentationAnalyzer) calculateUniformityScore(gaps []int, width int) float64 {
	if len(gaps) == 0 {
		return 1.0
	}

	// Calculate ideal gap size for uniform distribution
	idealGap := float64(width) / float64(len(gaps)+1)

	// Calculate variance from ideal
	variance := 0.0
	for _, gap := range gaps {
		diff := float64(gap) - idealGap
		variance += diff * diff
	}
	variance /= float64(len(gaps))

	// Convert variance to uniformity score (0 = clustered, 1 = uniform)
	maxVariance := idealGap * idealGap // Maximum possible variance
	if maxVariance == 0 {
		return 1.0
	}

	return math.Max(0, 1.0-(variance/maxVariance))
}

func (ra *RepresentationAnalyzer) calculateClusteringIndex(gaps []int) float64 {
	if len(gaps) == 0 {
		return 0.0
	}

	// Count small gaps (indicating clustering)
	smallGaps := 0
	for _, gap := range gaps {
		if gap < 2 { // Adjacent or nearly adjacent bits
			smallGaps++
		}
	}

	return float64(smallGaps) / float64(len(gaps))
}

func (ra *RepresentationAnalyzer) maxGap(gaps []int) int {
	if len(gaps) == 0 {
		return 0
	}

	max := gaps[0]
	for _, gap := range gaps[1:] {
		if gap > max {
			max = gap
		}
	}
	return max
}

func (ra *RepresentationAnalyzer) minGap(gaps []int) int {
	if len(gaps) == 0 {
		return 0
	}

	min := gaps[0]
	for _, gap := range gaps[1:] {
		if gap < min {
			min = gap
		}
	}
	return min
}

func (ra *RepresentationAnalyzer) averageGap(gaps []int) float64 {
	if len(gaps) == 0 {
		return 0.0
	}

	sum := 0
	for _, gap := range gaps {
		sum += gap
	}
	return float64(sum) / float64(len(gaps))
}

func (ra *RepresentationAnalyzer) combinations(n, k int) float64 {
	if k > n || k < 0 {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}

	// Use logarithms to avoid overflow for large numbers
	result := 0.0
	for i := 0; i < k; i++ {
		result += math.Log(float64(n-i)) - math.Log(float64(i+1))
	}
	return math.Exp(result)
}
