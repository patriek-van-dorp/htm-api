package sdr

import (
	"math"
)

// SimilarityCalculator provides various similarity metrics for SDRs
type SimilarityCalculator struct {
	// Configuration for similarity calculations
}

// NewSimilarityCalculator creates a new similarity calculator
func NewSimilarityCalculator() *SimilarityCalculator {
	return &SimilarityCalculator{}
}

// OverlapSimilarity calculates overlap-based similarity (0.0-1.0)
// This is the standard HTM similarity metric
func (sc *SimilarityCalculator) OverlapSimilarity(sdr1, sdr2 *SDR) float64 {
	if sdr1 == nil || sdr2 == nil {
		return 0.0
	}
	return sdr1.Similarity(sdr2)
}

// JaccardSimilarity calculates Jaccard index (intersection/union)
func (sc *SimilarityCalculator) JaccardSimilarity(sdr1, sdr2 *SDR) float64 {
	if sdr1 == nil || sdr2 == nil || sdr1.Width() != sdr2.Width() {
		return 0.0
	}

	if len(sdr1.ActiveBits()) == 0 && len(sdr2.ActiveBits()) == 0 {
		return 1.0 // Both empty, considered identical
	}

	intersection := float64(sdr1.Overlap(sdr2))
	union := float64(len(sdr1.ActiveBits())+len(sdr2.ActiveBits())) - intersection

	if union == 0 {
		return 0.0
	}

	return intersection / union
}

// CosineSimilarity calculates cosine similarity between SDRs
func (sc *SimilarityCalculator) CosineSimilarity(sdr1, sdr2 *SDR) float64 {
	if sdr1 == nil || sdr2 == nil || sdr1.Width() != sdr2.Width() {
		return 0.0
	}

	if len(sdr1.ActiveBits()) == 0 || len(sdr2.ActiveBits()) == 0 {
		return 0.0
	}

	intersection := float64(sdr1.Overlap(sdr2))
	magnitude1 := math.Sqrt(float64(len(sdr1.ActiveBits())))
	magnitude2 := math.Sqrt(float64(len(sdr2.ActiveBits())))

	return intersection / (magnitude1 * magnitude2)
}

// HammingDistance calculates Hamming distance between SDRs (number of differing bits)
func (sc *SimilarityCalculator) HammingDistance(sdr1, sdr2 *SDR) int {
	if sdr1 == nil || sdr2 == nil || sdr1.Width() != sdr2.Width() {
		return -1 // Invalid comparison
	}

	// Hamming distance = total active bits - 2 * overlap
	// This accounts for bits that are active in one but not the other
	overlap := sdr1.Overlap(sdr2)
	totalActiveBits := len(sdr1.ActiveBits()) + len(sdr2.ActiveBits())
	hammingDistance := totalActiveBits - 2*overlap

	return hammingDistance
}

// NormalizedHammingDistance calculates normalized Hamming distance (0.0-1.0)
func (sc *SimilarityCalculator) NormalizedHammingDistance(sdr1, sdr2 *SDR) float64 {
	if sdr1 == nil || sdr2 == nil || sdr1.Width() != sdr2.Width() {
		return 1.0 // Maximum distance for invalid comparison
	}

	hammingDistance := sc.HammingDistance(sdr1, sdr2)
	if hammingDistance < 0 {
		return 1.0
	}

	maxPossibleDistance := len(sdr1.ActiveBits()) + len(sdr2.ActiveBits())
	if maxPossibleDistance == 0 {
		return 0.0 // Both SDRs are empty
	}

	return float64(hammingDistance) / float64(maxPossibleDistance)
}

// SimilarityMatrix calculates pairwise similarities for a collection of SDRs
type SimilarityMatrix struct {
	sdrs         []*SDR
	similarities [][]float64
	calculator   *SimilarityCalculator
}

// NewSimilarityMatrix creates a similarity matrix for the given SDRs
func NewSimilarityMatrix(sdrs []*SDR) *SimilarityMatrix {
	n := len(sdrs)
	similarities := make([][]float64, n)
	for i := range similarities {
		similarities[i] = make([]float64, n)
	}

	calculator := NewSimilarityCalculator()

	matrix := &SimilarityMatrix{
		sdrs:         sdrs,
		similarities: similarities,
		calculator:   calculator,
	}

	matrix.calculate()
	return matrix
}

// calculate fills the similarity matrix
func (sm *SimilarityMatrix) calculate() {
	n := len(sm.sdrs)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				sm.similarities[i][j] = 1.0 // Self-similarity is 1.0
			} else if i < j {
				// Calculate similarity
				similarity := sm.calculator.OverlapSimilarity(sm.sdrs[i], sm.sdrs[j])
				sm.similarities[i][j] = similarity
				sm.similarities[j][i] = similarity // Matrix is symmetric
			}
		}
	}
}

// GetSimilarity returns the similarity between SDRs at indices i and j
func (sm *SimilarityMatrix) GetSimilarity(i, j int) float64 {
	if i < 0 || i >= len(sm.similarities) || j < 0 || j >= len(sm.similarities) {
		return 0.0
	}
	return sm.similarities[i][j]
}

// GetAverageSimilarity returns the average similarity in the matrix
func (sm *SimilarityMatrix) GetAverageSimilarity() float64 {
	if len(sm.similarities) == 0 {
		return 0.0
	}

	sum := 0.0
	count := 0

	n := len(sm.similarities)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ { // Only upper triangle to avoid double counting
			sum += sm.similarities[i][j]
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	return sum / float64(count)
}

// FindMostSimilarPair returns indices of the most similar pair of SDRs
func (sm *SimilarityMatrix) FindMostSimilarPair() (int, int, float64) {
	maxSimilarity := -1.0
	maxI, maxJ := -1, -1

	n := len(sm.similarities)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if sm.similarities[i][j] > maxSimilarity {
				maxSimilarity = sm.similarities[i][j]
				maxI, maxJ = i, j
			}
		}
	}

	return maxI, maxJ, maxSimilarity
}

// FindLeastSimilarPair returns indices of the least similar pair of SDRs
func (sm *SimilarityMatrix) FindLeastSimilarPair() (int, int, float64) {
	minSimilarity := 2.0 // Start above max possible similarity
	minI, minJ := -1, -1

	n := len(sm.similarities)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if sm.similarities[i][j] < minSimilarity {
				minSimilarity = sm.similarities[i][j]
				minI, minJ = i, j
			}
		}
	}

	return minI, minJ, minSimilarity
}

// SimilarityThreshold applies a threshold to classify SDR pairs as similar/dissimilar
type SimilarityThreshold struct {
	threshold float64
}

// NewSimilarityThreshold creates a new threshold classifier
func NewSimilarityThreshold(threshold float64) *SimilarityThreshold {
	return &SimilarityThreshold{threshold: threshold}
}

// AreSimilar returns true if two SDRs are above the similarity threshold
func (st *SimilarityThreshold) AreSimilar(sdr1, sdr2 *SDR) bool {
	calculator := NewSimilarityCalculator()
	similarity := calculator.OverlapSimilarity(sdr1, sdr2)
	return similarity >= st.threshold
}

// ClassifySimilarity classifies similarity as "high", "medium", or "low"
func (st *SimilarityThreshold) ClassifySimilarity(sdr1, sdr2 *SDR) string {
	calculator := NewSimilarityCalculator()
	similarity := calculator.OverlapSimilarity(sdr1, sdr2)

	if similarity >= 0.8 {
		return "high"
	} else if similarity >= 0.5 {
		return "medium"
	} else {
		return "low"
	}
}
