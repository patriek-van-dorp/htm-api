package spatial

import (
	"fmt"
	"math"
	"sort"

	"github.com/htm-project/neural-api/internal/domain/htm"
)

// InhibitionManager implements competitive inhibition algorithms for spatial pooling
type InhibitionManager struct {
	config            *htm.SpatialPoolerConfig
	neighborhoodCache map[int][]int // Cache for column neighborhoods
}

// NewInhibitionManager creates a new inhibition manager
func NewInhibitionManager(config *htm.SpatialPoolerConfig) *InhibitionManager {
	return &InhibitionManager{
		config:            config,
		neighborhoodCache: make(map[int][]int),
	}
}

// GlobalInhibition performs global competitive inhibition across all columns
func (im *InhibitionManager) GlobalInhibition(overlapScores []float64) ([]int, error) {
	targetActiveColumns := im.config.GetExpectedActiveColumns()

	// Create column-score pairs
	type columnScore struct {
		column int
		score  float64
	}

	candidates := make([]columnScore, 0, len(overlapScores))
	for i, score := range overlapScores {
		if score > 0 { // Only consider columns with positive overlap
			candidates = append(candidates, columnScore{column: i, score: score})
		}
	}

	// Sort by score (descending)
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].score == candidates[j].score {
			// Deterministic tie-breaking
			return candidates[i].column < candidates[j].column
		}
		return candidates[i].score > candidates[j].score
	})

	// Select top candidates
	activeCount := targetActiveColumns
	if len(candidates) < activeCount {
		activeCount = len(candidates)
	}

	activeColumns := make([]int, activeCount)
	for i := 0; i < activeCount; i++ {
		activeColumns[i] = candidates[i].column
	}

	// Sort for consistent output
	sort.Ints(activeColumns)
	return activeColumns, nil
}

// LocalInhibition performs local competitive inhibition within neighborhoods
func (im *InhibitionManager) LocalInhibition(overlapScores []float64) ([]int, error) {
	activeColumns := make([]bool, len(overlapScores))

	for col := 0; col < len(overlapScores); col++ {
		if overlapScores[col] <= 0 {
			continue // Skip columns with no overlap
		}

		// Get local neighborhood
		neighborhood := im.getNeighborhood(col)

		// Calculate desired local activity
		localDensity := im.config.LocalAreaDensity
		if localDensity == 0 {
			localDensity = float64(im.config.GetExpectedActiveColumns()) / float64(im.config.ColumnCount)
		}

		desiredLocalActivity := int(math.Ceil(float64(len(neighborhood)) * localDensity))
		if desiredLocalActivity < 1 {
			desiredLocalActivity = 1
		}

		// Check if this column should be active in its neighborhood
		if im.shouldBeActiveLocally(col, neighborhood, overlapScores, desiredLocalActivity) {
			activeColumns[col] = true
		}
	}

	// Convert boolean array to list of active column indices
	result := make([]int, 0)
	for i, isActive := range activeColumns {
		if isActive {
			result = append(result, i)
		}
	}

	return result, nil
}

// AdaptiveInhibition combines global and local inhibition strategies
func (im *InhibitionManager) AdaptiveInhibition(overlapScores []float64, globalWeight float64) ([]int, error) {
	// Perform both global and local inhibition
	globalActive, err := im.GlobalInhibition(overlapScores)
	if err != nil {
		return nil, err
	}

	localActive, err := im.LocalInhibition(overlapScores)
	if err != nil {
		return nil, err
	}

	// Combine results based on weight
	if globalWeight >= 1.0 {
		return globalActive, nil
	}
	if globalWeight <= 0.0 {
		return localActive, nil
	}

	// Weighted combination
	return im.combineInhibitionResults(globalActive, localActive, globalWeight), nil
}

// KWinnersInhibition implements k-winners-take-all inhibition
func (im *InhibitionManager) KWinnersInhibition(overlapScores []float64, k int) ([]int, error) {
	if k <= 0 {
		return []int{}, nil
	}
	if k >= len(overlapScores) {
		// All columns active
		result := make([]int, len(overlapScores))
		for i := range result {
			result[i] = i
		}
		return result, nil
	}

	// Find k-th largest score
	sortedScores := make([]float64, len(overlapScores))
	copy(sortedScores, overlapScores)
	sort.Float64s(sortedScores)

	// Get threshold (k-th largest score)
	thresholdIndex := len(sortedScores) - k
	threshold := sortedScores[thresholdIndex]

	// Select columns with scores >= threshold
	activeColumns := make([]int, 0, k)
	for i, score := range overlapScores {
		if score >= threshold && len(activeColumns) < k {
			activeColumns = append(activeColumns, i)
		}
	}

	sort.Ints(activeColumns)
	return activeColumns, nil
}

// BoostAwareInhibition performs inhibition considering boost factors
func (im *InhibitionManager) BoostAwareInhibition(overlapScores, boostFactors []float64) ([]int, error) {
	if len(overlapScores) != len(boostFactors) {
		return nil, fmt.Errorf("overlap scores and boost factors must have same length")
	}

	// Apply boost factors to overlap scores
	boostedScores := make([]float64, len(overlapScores))
	for i := range overlapScores {
		boostedScores[i] = overlapScores[i] * boostFactors[i]
	}

	// Use global inhibition on boosted scores
	return im.GlobalInhibition(boostedScores)
}

// Private helper methods

// getNeighborhood returns the neighborhood of columns for a given column
func (im *InhibitionManager) getNeighborhood(column int) []int {
	// Check cache first
	if neighbors, exists := im.neighborhoodCache[column]; exists {
		return neighbors
	}

	radius := im.config.InhibitionRadius
	if radius <= 0 {
		// No local inhibition, return just the column itself
		return []int{column}
	}

	// For 1D arrangement, neighborhood is columns within radius
	neighbors := make([]int, 0, 2*radius+1)
	start := max(0, column-radius)
	end := min(im.config.ColumnCount-1, column+radius)

	for i := start; i <= end; i++ {
		neighbors = append(neighbors, i)
	}

	// Cache the result
	im.neighborhoodCache[column] = neighbors
	return neighbors
}

// shouldBeActiveLocally determines if a column should be active within its neighborhood
func (im *InhibitionManager) shouldBeActiveLocally(column int, neighborhood []int, overlapScores []float64, desiredActivity int) bool {
	// Get scores for neighborhood
	neighborScores := make([]float64, len(neighborhood))
	for i, neighbor := range neighborhood {
		neighborScores[i] = overlapScores[neighbor]
	}

	// Sort neighborhood by score (descending)
	type neighborScore struct {
		index int
		score float64
	}

	neighbors := make([]neighborScore, len(neighborhood))
	for i, neighbor := range neighborhood {
		neighbors[i] = neighborScore{index: neighbor, score: overlapScores[neighbor]}
	}

	sort.Slice(neighbors, func(i, j int) bool {
		if neighbors[i].score == neighbors[j].score {
			// Deterministic tie-breaking
			return neighbors[i].index < neighbors[j].index
		}
		return neighbors[i].score > neighbors[j].score
	})

	// Check if this column is in the top desiredActivity
	for i := 0; i < desiredActivity && i < len(neighbors); i++ {
		if neighbors[i].index == column {
			return true
		}
	}

	return false
}

// combineInhibitionResults combines global and local inhibition results
func (im *InhibitionManager) combineInhibitionResults(globalActive, localActive []int, globalWeight float64) []int {
	// Convert to sets for easier manipulation
	globalSet := make(map[int]bool)
	for _, col := range globalActive {
		globalSet[col] = true
	}

	localSet := make(map[int]bool)
	for _, col := range localActive {
		localSet[col] = true
	}

	// Weighted combination: include columns based on probability
	combined := make(map[int]bool)

	// Add all columns that appear in both
	for col := range globalSet {
		if localSet[col] {
			combined[col] = true
		}
	}

	// Add remaining global columns with probability = globalWeight
	for col := range globalSet {
		if !combined[col] && im.shouldIncludeWithProbability(globalWeight) {
			combined[col] = true
		}
	}

	// Add remaining local columns with probability = (1 - globalWeight)
	for col := range localSet {
		if !combined[col] && im.shouldIncludeWithProbability(1.0-globalWeight) {
			combined[col] = true
		}
	}

	// Convert back to sorted slice
	result := make([]int, 0, len(combined))
	for col := range combined {
		result = append(result, col)
	}
	sort.Ints(result)

	return result
}

// shouldIncludeWithProbability returns true with given probability (deterministic for reproducibility)
func (im *InhibitionManager) shouldIncludeWithProbability(probability float64) bool {
	// For deterministic behavior, use a simple threshold
	// In a real implementation, this could use a seeded random number generator
	return probability > 0.5
}

// InhibitionMetrics contains metrics about inhibition performance
type InhibitionMetrics struct {
	TotalColumns         int     `json:"total_columns"`
	ActiveColumns        int     `json:"active_columns"`
	ActualSparsity       float64 `json:"actual_sparsity"`
	TargetSparsity       float64 `json:"target_sparsity"`
	SparsityError        float64 `json:"sparsity_error"`
	InhibitionEfficiency float64 `json:"inhibition_efficiency"`
}

// CalculateInhibitionMetrics calculates metrics for inhibition performance
func (im *InhibitionManager) CalculateInhibitionMetrics(activeColumns []int) *InhibitionMetrics {
	totalColumns := im.config.ColumnCount
	activeCount := len(activeColumns)
	actualSparsity := float64(activeCount) / float64(totalColumns)
	targetSparsity := im.config.SparsityRatio
	sparsityError := math.Abs(actualSparsity - targetSparsity)

	// Efficiency: how well we hit the target
	efficiency := 1.0 - (sparsityError / targetSparsity)
	if efficiency < 0 {
		efficiency = 0
	}

	return &InhibitionMetrics{
		TotalColumns:         totalColumns,
		ActiveColumns:        activeCount,
		ActualSparsity:       actualSparsity,
		TargetSparsity:       targetSparsity,
		SparsityError:        sparsityError,
		InhibitionEfficiency: efficiency,
	}
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
