package sdr

import (
	"fmt"
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

// CalculateAllSimilarities computes all available similarity metrics
func (sc *SimilarityCalculator) CalculateAllSimilarities(sdr1, sdr2 *SDR) *SimilarityMetrics {
	if sdr1 == nil || sdr2 == nil {
		return &SimilarityMetrics{
			IsValid: false,
			Error:   "one or both SDRs are nil",
		}
	}

	if sdr1.Width != sdr2.Width {
		return &SimilarityMetrics{
			IsValid: false,
			Error:   "SDRs must have same width",
		}
	}

	return &SimilarityMetrics{
		IsValid:           true,
		OverlapSimilarity: sc.OverlapSimilarity(sdr1, sdr2),
		JaccardSimilarity: sc.JaccardSimilarity(sdr1, sdr2),
		CosineSimilarity:  sc.CosineSimilarity(sdr1, sdr2),
		DiceSimilarity:    sc.DiceSimilarity(sdr1, sdr2),
		HammingDistance:   sc.HammingDistance(sdr1, sdr2),
		EuclideanDistance: sc.EuclideanDistance(sdr1, sdr2),
		OverlapCount:      sdr1.Overlap(sdr2),
	}
}

// SimilarityMetrics contains all computed similarity measures
type SimilarityMetrics struct {
	IsValid           bool    `json:"is_valid"`
	Error             string  `json:"error,omitempty"`
	OverlapSimilarity float64 `json:"overlap_similarity"`
	JaccardSimilarity float64 `json:"jaccard_similarity"`
	CosineSimilarity  float64 `json:"cosine_similarity"`
	DiceSimilarity    float64 `json:"dice_similarity"`
	HammingDistance   int     `json:"hamming_distance"`
	EuclideanDistance float64 `json:"euclidean_distance"`
	OverlapCount      int     `json:"overlap_count"`
}

// OverlapSimilarity calculates overlap-based similarity (0.0-1.0)
// This is the standard HTM similarity metric
func (sc *SimilarityCalculator) OverlapSimilarity(sdr1, sdr2 *SDR) float64 {
	if sdr1 == nil || sdr2 == nil {
		return 0.0
	}
	return sdr1.OverlapRatio(sdr2)
}

// JaccardSimilarity calculates Jaccard index (intersection/union)
func (sc *SimilarityCalculator) JaccardSimilarity(sdr1, sdr2 *SDR) float64 {
	if sdr1 == nil || sdr2 == nil {
		return 0.0
	}
	return sdr1.JaccardSimilarity(sdr2)
}

// CosineSimilarity calculates cosine similarity between SDRs
func (sc *SimilarityCalculator) CosineSimilarity(sdr1, sdr2 *SDR) float64 {
	if sdr1 == nil || sdr2 == nil {
		return 0.0
	}
	return sdr1.CosineSimilarity(sdr2)
}

// DiceSimilarity calculates Dice coefficient (Sørensen-Dice index)
func (sc *SimilarityCalculator) DiceSimilarity(sdr1, sdr2 *SDR) float64 {
	if sdr1 == nil || sdr2 == nil || sdr1.Width != sdr2.Width {
		return 0.0
	}

	if len(sdr1.ActiveBits) == 0 && len(sdr2.ActiveBits) == 0 {
		return 1.0 // Both empty
	}

	intersection := float64(sdr1.Overlap(sdr2))
	totalActiveBits := float64(len(sdr1.ActiveBits) + len(sdr2.ActiveBits))

	if totalActiveBits == 0 {
		return 0.0
	}

	return (2.0 * intersection) / totalActiveBits
}

// HammingDistance calculates Hamming distance between SDRs
func (sc *SimilarityCalculator) HammingDistance(sdr1, sdr2 *SDR) int {
	if sdr1 == nil || sdr2 == nil || sdr1.Width != sdr2.Width {
		return -1 // Invalid distance
	}

	// For sparse representations, Hamming distance = bits only in sdr1 + bits only in sdr2
	overlap := sdr1.Overlap(sdr2)
	return len(sdr1.ActiveBits) + len(sdr2.ActiveBits) - 2*overlap
}

// EuclideanDistance calculates Euclidean distance between SDRs
func (sc *SimilarityCalculator) EuclideanDistance(sdr1, sdr2 *SDR) float64 {
	if sdr1 == nil || sdr2 == nil || sdr1.Width != sdr2.Width {
		return -1.0 // Invalid distance
	}

	// For binary vectors, Euclidean distance = sqrt(Hamming distance)
	hammingDist := sc.HammingDistance(sdr1, sdr2)
	if hammingDist < 0 {
		return -1.0
	}

	return math.Sqrt(float64(hammingDist))
}

// BatchSimilarity calculates similarity matrix for a batch of SDRs
func (sc *SimilarityCalculator) BatchSimilarity(sdrs []*SDR, metric SimilarityMetric) ([][]float64, error) {
	if len(sdrs) == 0 {
		return [][]float64{}, nil
	}

	// Validate all SDRs have same width
	width := sdrs[0].Width
	for i, sdr := range sdrs {
		if sdr.Width != width {
			return nil, fmt.Errorf("SDR %d has width %d, expected %d", i, sdr.Width, width)
		}
	}

	n := len(sdrs)
	matrix := make([][]float64, n)
	for i := range matrix {
		matrix[i] = make([]float64, n)
	}

	// Calculate upper triangle (symmetric matrix)
	for i := 0; i < n; i++ {
		matrix[i][i] = 1.0 // Self-similarity is 1.0
		for j := i + 1; j < n; j++ {
			similarity := sc.calculateSimilarity(sdrs[i], sdrs[j], metric)
			matrix[i][j] = similarity
			matrix[j][i] = similarity // Symmetric
		}
	}

	return matrix, nil
}

// SimilarityMetric represents different similarity calculation methods
type SimilarityMetric int

const (
	OverlapMetric SimilarityMetric = iota
	JaccardMetric
	CosineMetric
	DiceMetric
)

// SemanticSimilarityAnalyzer analyzes semantic relationships between SDRs
type SemanticSimilarityAnalyzer struct {
	similarThreshold   float64
	differentThreshold float64
}

// NewSemanticSimilarityAnalyzer creates a new semantic similarity analyzer
func NewSemanticSimilarityAnalyzer(similarThreshold, differentThreshold float64) *SemanticSimilarityAnalyzer {
	return &SemanticSimilarityAnalyzer{
		similarThreshold:   similarThreshold,
		differentThreshold: differentThreshold,
	}
}

// AnalyzeSemanticPreservation analyzes how well spatial pooler preserves semantic relationships
func (ssa *SemanticSimilarityAnalyzer) AnalyzeSemanticPreservation(inputSDRs, outputSDRs []*SDR) (*SemanticAnalysis, error) {
	if len(inputSDRs) != len(outputSDRs) {
		return nil, fmt.Errorf("input and output SDR arrays must have same length")
	}

	calc := NewSimilarityCalculator()
	inputMatrix, err := calc.BatchSimilarity(inputSDRs, OverlapMetric)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate input similarity matrix: %v", err)
	}

	outputMatrix, err := calc.BatchSimilarity(outputSDRs, OverlapMetric)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate output similarity matrix: %v", err)
	}

	return ssa.analyzeMatrices(inputMatrix, outputMatrix), nil
}

// SemanticAnalysis contains results of semantic preservation analysis
type SemanticAnalysis struct {
	TotalComparisons        int     `json:"total_comparisons"`
	SemanticViolations      int     `json:"semantic_violations"`
	PreservationScore       float64 `json:"preservation_score"`
	SimilarPairsPreserved   int     `json:"similar_pairs_preserved"`
	SimilarPairsTotal       int     `json:"similar_pairs_total"`
	DifferentPairsPreserved int     `json:"different_pairs_preserved"`
	DifferentPairsTotal     int     `json:"different_pairs_total"`
	CorrelationCoefficient  float64 `json:"correlation_coefficient"`
}

// TemporalSimilarityTracker tracks similarity changes over time
type TemporalSimilarityTracker struct {
	history    []SimilaritySnapshot
	windowSize int
}

// SimilaritySnapshot represents similarity state at a point in time
type SimilaritySnapshot struct {
	Timestamp    int64       `json:"timestamp"`
	Similarities [][]float64 `json:"similarities"`
	SDRCount     int         `json:"sdr_count"`
}

// NewTemporalSimilarityTracker creates a new temporal similarity tracker
func NewTemporalSimilarityTracker(windowSize int) *TemporalSimilarityTracker {
	return &TemporalSimilarityTracker{
		history:    make([]SimilaritySnapshot, 0),
		windowSize: windowSize,
	}
}

// AddSnapshot adds a new similarity snapshot
func (tst *TemporalSimilarityTracker) AddSnapshot(timestamp int64, sdrs []*SDR) error {
	calc := NewSimilarityCalculator()
	similarities, err := calc.BatchSimilarity(sdrs, OverlapMetric)
	if err != nil {
		return err
	}

	snapshot := SimilaritySnapshot{
		Timestamp:    timestamp,
		Similarities: similarities,
		SDRCount:     len(sdrs),
	}

	tst.history = append(tst.history, snapshot)

	// Maintain window size
	if len(tst.history) > tst.windowSize {
		tst.history = tst.history[1:]
	}

	return nil
}

// GetStabilityMetrics calculates temporal stability metrics
func (tst *TemporalSimilarityTracker) GetStabilityMetrics() *StabilityMetrics {
	if len(tst.history) < 2 {
		return &StabilityMetrics{
			HasSufficientData: false,
		}
	}

	variance := tst.calculateSimilarityVariance()
	trend := tst.calculateSimilarityTrend()

	return &StabilityMetrics{
		HasSufficientData:  true,
		SimilarityVariance: variance,
		SimilarityTrend:    trend,
		SnapshotCount:      len(tst.history),
		TimeSpan:           tst.history[len(tst.history)-1].Timestamp - tst.history[0].Timestamp,
	}
}

// StabilityMetrics contains temporal stability analysis
type StabilityMetrics struct {
	HasSufficientData  bool    `json:"has_sufficient_data"`
	SimilarityVariance float64 `json:"similarity_variance"`
	SimilarityTrend    float64 `json:"similarity_trend"`
	SnapshotCount      int     `json:"snapshot_count"`
	TimeSpan           int64   `json:"time_span"`
}

// Helper methods

func (sc *SimilarityCalculator) calculateSimilarity(sdr1, sdr2 *SDR, metric SimilarityMetric) float64 {
	switch metric {
	case OverlapMetric:
		return sc.OverlapSimilarity(sdr1, sdr2)
	case JaccardMetric:
		return sc.JaccardSimilarity(sdr1, sdr2)
	case CosineMetric:
		return sc.CosineSimilarity(sdr1, sdr2)
	case DiceMetric:
		return sc.DiceSimilarity(sdr1, sdr2)
	default:
		return sc.OverlapSimilarity(sdr1, sdr2)
	}
}

func (ssa *SemanticSimilarityAnalyzer) analyzeMatrices(inputMatrix, outputMatrix [][]float64) *SemanticAnalysis {
	n := len(inputMatrix)
	totalComparisons := 0
	semanticViolations := 0
	similarPairsPreserved := 0
	similarPairsTotal := 0
	differentPairsPreserved := 0
	differentPairsTotal := 0

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			inputSim := inputMatrix[i][j]
			outputSim := outputMatrix[i][j]
			totalComparisons++

			if inputSim >= ssa.similarThreshold {
				similarPairsTotal++
				if outputSim >= ssa.similarThreshold {
					similarPairsPreserved++
				} else {
					semanticViolations++
				}
			} else if inputSim <= ssa.differentThreshold {
				differentPairsTotal++
				if outputSim <= ssa.differentThreshold {
					differentPairsPreserved++
				} else {
					semanticViolations++
				}
			}
		}
	}

	preservationScore := 1.0
	if totalComparisons > 0 {
		preservationScore = 1.0 - (float64(semanticViolations) / float64(totalComparisons))
	}

	correlation := ssa.calculateCorrelation(inputMatrix, outputMatrix)

	return &SemanticAnalysis{
		TotalComparisons:        totalComparisons,
		SemanticViolations:      semanticViolations,
		PreservationScore:       preservationScore,
		SimilarPairsPreserved:   similarPairsPreserved,
		SimilarPairsTotal:       similarPairsTotal,
		DifferentPairsPreserved: differentPairsPreserved,
		DifferentPairsTotal:     differentPairsTotal,
		CorrelationCoefficient:  correlation,
	}
}

func (ssa *SemanticSimilarityAnalyzer) calculateCorrelation(matrix1, matrix2 [][]float64) float64 {
	n := len(matrix1)
	if n == 0 {
		return 0.0
	}

	var sumX, sumY, sumXY, sumX2, sumY2 float64
	count := 0

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			x := matrix1[i][j]
			y := matrix2[i][j]

			sumX += x
			sumY += y
			sumXY += x * y
			sumX2 += x * x
			sumY2 += y * y
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	// Pearson correlation coefficient
	numerator := float64(count)*sumXY - sumX*sumY
	denominator := math.Sqrt((float64(count)*sumX2 - sumX*sumX) * (float64(count)*sumY2 - sumY*sumY))

	if denominator == 0 {
		return 0.0
	}

	return numerator / denominator
}

func (tst *TemporalSimilarityTracker) calculateSimilarityVariance() float64 {
	if len(tst.history) < 2 {
		return 0.0
	}

	// Calculate variance of average similarities over time
	averages := make([]float64, len(tst.history))
	for i, snapshot := range tst.history {
		total := 0.0
		count := 0
		n := len(snapshot.Similarities)

		for row := 0; row < n; row++ {
			for col := row + 1; col < n; col++ {
				total += snapshot.Similarities[row][col]
				count++
			}
		}

		if count > 0 {
			averages[i] = total / float64(count)
		}
	}

	// Calculate variance
	mean := 0.0
	for _, avg := range averages {
		mean += avg
	}
	mean /= float64(len(averages))

	variance := 0.0
	for _, avg := range averages {
		diff := avg - mean
		variance += diff * diff
	}
	variance /= float64(len(averages))

	return variance
}

func (tst *TemporalSimilarityTracker) calculateSimilarityTrend() float64 {
	if len(tst.history) < 2 {
		return 0.0
	}

	// Simple linear trend: (last - first) / time_span
	first := tst.getAverageSimilarity(tst.history[0])
	last := tst.getAverageSimilarity(tst.history[len(tst.history)-1])
	timeSpan := float64(tst.history[len(tst.history)-1].Timestamp - tst.history[0].Timestamp)

	if timeSpan == 0 {
		return 0.0
	}

	return (last - first) / timeSpan
}

func (tst *TemporalSimilarityTracker) getAverageSimilarity(snapshot SimilaritySnapshot) float64 {
	total := 0.0
	count := 0
	n := len(snapshot.Similarities)

	for row := 0; row < n; row++ {
		for col := row + 1; col < n; col++ {
			total += snapshot.Similarities[row][col]
			count++
		}
	}

	if count == 0 {
		return 0.0
	}
	return total / float64(count)
}
