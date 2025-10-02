package integration

import (
	"math"
	"sort"
	"testing"
)

// TestSpatialPoolerSemanticSimilarity tests that the spatial pooler preserves semantic relationships between inputs
func TestSpatialPoolerSemanticSimilarity(t *testing.T) {
	t.Skip("Skipping until spatial pooler implementation is complete")

	// This test validates FR-003: System MUST maintain semantic continuity where similar inputs
	// produce overlapping SDR patterns with 30-70% bit overlap, while different inputs have <20% overlap

	t.Run("similar_inputs_high_overlap", func(t *testing.T) {
		// Test that similar inputs produce SDRs with 30-70% overlap

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		// Create similar inputs (overlapping active bits)
		baseInput := createTestEncoderOutput([]int{10, 20, 30, 40, 50}, 2048)
		similarInput1 := createTestEncoderOutput([]int{10, 20, 30, 45, 55}, 2048) // 3/5 overlap with base
		similarInput2 := createTestEncoderOutput([]int{12, 22, 32, 40, 50}, 2048) // 2/5 overlap with base

		// Process inputs
		// baseResult := spatialPooler.Process(baseInput)
		// similarResult1 := spatialPooler.Process(similarInput1)
		// similarResult2 := spatialPooler.Process(similarInput2)

		// Calculate overlaps
		// overlap1 := calculateSDROverlapPercentage(baseResult.NormalizedSDR, similarResult1.NormalizedSDR)
		// overlap2 := calculateSDROverlapPercentage(baseResult.NormalizedSDR, similarResult2.NormalizedSDR)

		// Similar inputs should have 30-70% overlap in their output SDRs
		// assert.GreaterOrEqual(t, overlap1, 0.30, "Similar inputs should have >= 30% SDR overlap")
		// assert.LessOrEqual(t, overlap1, 0.70, "Similar inputs should have <= 70% SDR overlap")
		// assert.GreaterOrEqual(t, overlap2, 0.30, "Similar inputs should have >= 30% SDR overlap")
		// assert.LessOrEqual(t, overlap2, 0.70, "Similar inputs should have <= 70% SDR overlap")
	})

	t.Run("different_inputs_low_overlap", func(t *testing.T) {
		// Test that different inputs produce SDRs with <20% overlap

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		// Create very different inputs (no overlapping active bits)
		input1 := createTestEncoderOutput([]int{10, 20, 30, 40, 50}, 2048)
		input2 := createTestEncoderOutput([]int{100, 200, 300, 400, 500}, 2048)      // No overlap
		input3 := createTestEncoderOutput([]int{1000, 1100, 1200, 1300, 1400}, 2048) // No overlap

		// Process inputs
		// result1 := spatialPooler.Process(input1)
		// result2 := spatialPooler.Process(input2)
		// result3 := spatialPooler.Process(input3)

		// Calculate overlaps
		// overlap12 := calculateSDROverlapPercentage(result1.NormalizedSDR, result2.NormalizedSDR)
		// overlap13 := calculateSDROverlapPercentage(result1.NormalizedSDR, result3.NormalizedSDR)
		// overlap23 := calculateSDROverlapPercentage(result2.NormalizedSDR, result3.NormalizedSDR)

		// Different inputs should have <20% overlap in their output SDRs
		// assert.Less(t, overlap12, 0.20, "Different inputs should have < 20% SDR overlap")
		// assert.Less(t, overlap13, 0.20, "Different inputs should have < 20% SDR overlap")
		// assert.Less(t, overlap23, 0.20, "Different inputs should have < 20% SDR overlap")
	})

	t.Run("semantic_gradient_preservation", func(t *testing.T) {
		// Test that semantic gradients are preserved (more similar inputs have higher overlap)

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		baseInput := createTestEncoderOutput([]int{10, 20, 30, 40, 50}, 2048)

		// Create inputs with varying degrees of similarity
		verySimilar := createTestEncoderOutput([]int{10, 20, 30, 40, 51}, 2048)             // 4/5 overlap
		moderatelySimilar := createTestEncoderOutput([]int{10, 20, 35, 45, 55}, 2048)       // 2/5 overlap
		slightlySimilar := createTestEncoderOutput([]int{15, 25, 35, 45, 55}, 2048)         // 0/5 but close values
		veryDifferent := createTestEncoderOutput([]int{1000, 1100, 1200, 1300, 1400}, 2048) // No similarity

		// Process all inputs
		// baseResult := spatialPooler.Process(baseInput)
		// verySimilarResult := spatialPooler.Process(verySimilar)
		// moderatelySimilarResult := spatialPooler.Process(moderatelySimilar)
		// slightlySimilarResult := spatialPooler.Process(slightlySimilar)
		// veryDifferentResult := spatialPooler.Process(veryDifferent)

		// Calculate overlaps
		// overlapVerySimilar := calculateSDROverlapPercentage(baseResult.NormalizedSDR, verySimilarResult.NormalizedSDR)
		// overlapModerately := calculateSDROverlapPercentage(baseResult.NormalizedSDR, moderatelySimilarResult.NormalizedSDR)
		// overlapSlightly := calculateSDROverlapPercentage(baseResult.NormalizedSDR, slightlySimilarResult.NormalizedSDR)
		// overlapVeryDifferent := calculateSDROverlapPercentage(baseResult.NormalizedSDR, veryDifferentResult.NormalizedSDR)

		// Semantic gradient should be preserved: more similar inputs should have higher overlap
		// assert.Greater(t, overlapVerySimilar, overlapModerately, "Very similar should have higher overlap than moderately similar")
		// assert.Greater(t, overlapModerately, overlapSlightly, "Moderately similar should have higher overlap than slightly similar")
		// assert.Greater(t, overlapSlightly, overlapVeryDifferent, "Slightly similar should have higher overlap than very different")

		// Verify ranges
		// assert.GreaterOrEqual(t, overlapVerySimilar, 0.50, "Very similar inputs should have high overlap")
		// assert.LessOrEqual(t, overlapVeryDifferent, 0.20, "Very different inputs should have low overlap")
	})

	t.Run("categorical_similarity_preservation", func(t *testing.T) {
		// Test semantic similarity for categorical-like inputs

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		// Simulate categorical encoder outputs for colors
		redInputs := []TestInput{
			createTestEncoderOutput([]int{0, 1, 2, 3, 4}, 2048), // Red variant 1
			createTestEncoderOutput([]int{0, 1, 2, 3, 5}, 2048), // Red variant 2
			createTestEncoderOutput([]int{0, 1, 2, 4, 5}, 2048), // Red variant 3
		}

		blueInputs := []TestInput{
			createTestEncoderOutput([]int{100, 101, 102, 103, 104}, 2048), // Blue variant 1
			createTestEncoderOutput([]int{100, 101, 102, 103, 105}, 2048), // Blue variant 2
			createTestEncoderOutput([]int{100, 101, 102, 104, 105}, 2048), // Blue variant 3
		}

		greenInputs := []TestInput{
			createTestEncoderOutput([]int{200, 201, 202, 203, 204}, 2048), // Green variant 1
			createTestEncoderOutput([]int{200, 201, 202, 203, 205}, 2048), // Green variant 2
			createTestEncoderOutput([]int{200, 201, 202, 204, 205}, 2048), // Green variant 3
		}

		// Process all inputs
		var redResults, blueResults, greenResults []TestResult
		for _, input := range redInputs {
			// result := spatialPooler.Process(input)
			// redResults = append(redResults, result)
		}
		for _, input := range blueInputs {
			// result := spatialPooler.Process(input)
			// blueResults = append(blueResults, result)
		}
		for _, input := range greenInputs {
			// result := spatialPooler.Process(input)
			// greenResults = append(greenResults, result)
		}

		// Calculate intra-category overlaps (should be high)
		// redIntraOverlaps := calculateIntraCategoryOverlaps(redResults)
		// blueIntraOverlaps := calculateIntraCategoryOverlaps(blueResults)
		// greenIntraOverlaps := calculateIntraCategoryOverlaps(greenResults)

		// Calculate inter-category overlaps (should be low)
		// redBlueOverlaps := calculateInterCategoryOverlaps(redResults, blueResults)
		// redGreenOverlaps := calculateInterCategoryOverlaps(redResults, greenResults)
		// blueGreenOverlaps := calculateInterCategoryOverlaps(blueResults, greenResults)

		// Verify intra-category similarity (same category should have high overlap)
		for _, overlap := range redIntraOverlaps {
			// assert.GreaterOrEqual(t, overlap, 0.30, "Red variants should have high intra-category overlap")
		}
		for _, overlap := range blueIntraOverlaps {
			// assert.GreaterOrEqual(t, overlap, 0.30, "Blue variants should have high intra-category overlap")
		}
		for _, overlap := range greenIntraOverlaps {
			// assert.GreaterOrEqual(t, overlap, 0.30, "Green variants should have high intra-category overlap")
		}

		// Verify inter-category dissimilarity (different categories should have low overlap)
		for _, overlap := range redBlueOverlaps {
			// assert.Less(t, overlap, 0.20, "Red-Blue should have low inter-category overlap")
		}
		for _, overlap := range redGreenOverlaps {
			// assert.Less(t, overlap, 0.20, "Red-Green should have low inter-category overlap")
		}
		for _, overlap := range blueGreenOverlaps {
			// assert.Less(t, overlap, 0.20, "Blue-Green should have low inter-category overlap")
		}
	})
}

// TestSpatialPoolerSemanticStability tests that semantic relationships remain stable over time
func TestSpatialPoolerSemanticStability(t *testing.T) {
	t.Skip("Skipping until spatial pooler implementation is complete")

	t.Run("semantic_stability_without_learning", func(t *testing.T) {
		// Test that semantic relationships are stable when learning is disabled

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		input1 := createTestEncoderOutput([]int{10, 20, 30, 40, 50}, 2048)
		input2 := createTestEncoderOutput([]int{12, 22, 32, 42, 52}, 2048) // Similar to input1

		// Measure initial overlap
		// result1_initial := spatialPooler.Process(input1)
		// result2_initial := spatialPooler.Process(input2)
		// initialOverlap := calculateSDROverlapPercentage(result1_initial.NormalizedSDR, result2_initial.NormalizedSDR)

		// Process many other random inputs to potentially change internal state
		for i := 0; i < 1000; i++ {
			randomInput := createTestEncoderOutput(generateRandomActiveBits(5, 2048), 2048)
			// spatialPooler.Process(randomInput)
		}

		// Measure overlap again
		// result1_final := spatialPooler.Process(input1)
		// result2_final := spatialPooler.Process(input2)
		// finalOverlap := calculateSDROverlapPercentage(result1_final.NormalizedSDR, result2_final.NormalizedSDR)

		// Semantic relationship should be stable
		// assert.InDelta(t, initialOverlap, finalOverlap, 0.05, "Semantic overlap should be stable without learning")
	})

	t.Run("semantic_adaptation_with_learning", func(t *testing.T) {
		// Test that semantic relationships can adapt with learning while maintaining basic structure

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=true)

		input1 := createTestEncoderOutput([]int{10, 20, 30, 40, 50}, 2048)
		input2 := createTestEncoderOutput([]int{12, 22, 32, 42, 52}, 2048)      // Similar to input1
		input3 := createTestEncoderOutput([]int{100, 200, 300, 400, 500}, 2048) // Different from input1

		// Measure initial relationships
		// result1_initial := spatialPooler.Process(input1)
		// result2_initial := spatialPooler.Process(input2)
		// result3_initial := spatialPooler.Process(input3)

		// initialSimilarOverlap := calculateSDROverlapPercentage(result1_initial.NormalizedSDR, result2_initial.NormalizedSDR)
		// initialDifferentOverlap := calculateSDROverlapPercentage(result1_initial.NormalizedSDR, result3_initial.NormalizedSDR)

		// Process inputs repeatedly to trigger learning
		for i := 0; i < 1000; i++ {
			// spatialPooler.Process(input1)
			// spatialPooler.Process(input2)
			// spatialPooler.Process(input3)
		}

		// Measure final relationships
		// result1_final := spatialPooler.Process(input1)
		// result2_final := spatialPooler.Process(input2)
		// result3_final := spatialPooler.Process(input3)

		// finalSimilarOverlap := calculateSDROverlapPercentage(result1_final.NormalizedSDR, result2_final.NormalizedSDR)
		// finalDifferentOverlap := calculateSDROverlapPercentage(result1_final.NormalizedSDR, result3_final.NormalizedSDR)

		// Similar inputs should still have higher overlap than different inputs, even after learning
		// assert.Greater(t, finalSimilarOverlap, finalDifferentOverlap, "Learning should preserve semantic relationships")

		// Overlaps should still be in valid ranges
		// assert.GreaterOrEqual(t, finalSimilarOverlap, 0.30, "Similar inputs should maintain high overlap after learning")
		// assert.Less(t, finalDifferentOverlap, 0.20, "Different inputs should maintain low overlap after learning")
	})
}

// Helper functions for semantic similarity testing

func calculateSDROverlapPercentage(sdr1, sdr2 SDR) float64 {
	// Calculate the percentage overlap between two SDRs
	if len(sdr1.ActiveBits) == 0 || len(sdr2.ActiveBits) == 0 {
		return 0.0
	}

	// Convert to sets for easier intersection calculation
	set1 := make(map[int]bool)
	for _, bit := range sdr1.ActiveBits {
		set1[bit] = true
	}

	intersectionCount := 0
	for _, bit := range sdr2.ActiveBits {
		if set1[bit] {
			intersectionCount++
		}
	}

	// Calculate overlap as percentage of the smaller SDR
	minSize := len(sdr1.ActiveBits)
	if len(sdr2.ActiveBits) < minSize {
		minSize = len(sdr2.ActiveBits)
	}

	return float64(intersectionCount) / float64(minSize)
}

func calculateIntraCategoryOverlaps(results []TestResult) []float64 {
	var overlaps []float64
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			overlap := calculateSDROverlapPercentage(results[i].NormalizedSDR, results[j].NormalizedSDR)
			overlaps = append(overlaps, overlap)
		}
	}
	return overlaps
}

func calculateInterCategoryOverlaps(results1, results2 []TestResult) []float64 {
	var overlaps []float64
	for _, result1 := range results1 {
		for _, result2 := range results2 {
			overlap := calculateSDROverlapPercentage(result1.NormalizedSDR, result2.NormalizedSDR)
			overlaps = append(overlaps, overlap)
		}
	}
	return overlaps
}

func generateRandomActiveBits(count, width int) []int {
	// Generate random active bit positions for testing
	bits := make([]int, count)
	for i := 0; i < count; i++ {
		bits[i] = int(math.Mod(float64(i*123+456), float64(width))) // Simple deterministic "random"
	}
	sort.Ints(bits)
	return bits
}
