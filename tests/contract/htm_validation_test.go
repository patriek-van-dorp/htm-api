package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/htm-project/neural-api/internal/api"
	"github.com/htm-project/neural-api/internal/handlers"
	"github.com/htm-project/neural-api/internal/infrastructure/config"
	"github.com/htm-project/neural-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HTMValidationResponse represents the HTM properties validation response
type HTMValidationResponse struct {
	Valid                bool                 `json:"valid"`
	Timestamp            string               `json:"timestamp"`
	BiologicalCompliance BiologicalCompliance `json:"biological_compliance"`
	SparsityAnalysis     SparsityAnalysis     `json:"sparsity_analysis"`
	OverlapAnalysis      OverlapAnalysis      `json:"overlap_analysis"`
	LearningAnalysis     LearningAnalysis     `json:"learning_analysis"`
	InhibitionAnalysis   InhibitionAnalysis   `json:"inhibition_analysis"`
	NetworkTopology      NetworkTopology      `json:"network_topology"`
	Recommendations      []string             `json:"recommendations,omitempty"`
}

// BiologicalCompliance represents overall biological compliance
type BiologicalCompliance struct {
	OverallScore         float64              `json:"overall_score"`
	ComplianceLevel      string               `json:"compliance_level"`
	CriticalIssues       []string             `json:"critical_issues,omitempty"`
	BiologicalPrinciples BiologicalPrinciples `json:"biological_principles"`
}

// BiologicalPrinciples represents adherence to biological principles
type BiologicalPrinciples struct {
	SparseCoding        bool `json:"sparse_coding"`
	LocalInhibition     bool `json:"local_inhibition"`
	HebbianLearning     bool `json:"hebbian_learning"`
	CompetitiveLearning bool `json:"competitive_learning"`
	AdaptiveBoosting    bool `json:"adaptive_boosting"`
	DendriteComputation bool `json:"dendrite_computation"`
}

// SparsityAnalysis represents sparsity validation analysis
type SparsityAnalysis struct {
	CurrentSparsity float64                 `json:"current_sparsity"`
	TargetSparsity  float64                 `json:"target_sparsity"`
	BiologicalRange BiologicalSparsityRange `json:"biological_range"`
	SparsityValid   bool                    `json:"sparsity_valid"`
	SparsityHealth  string                  `json:"sparsity_health"`
	SparsityHistory []SparsityDataPoint     `json:"sparsity_history"`
	Issues          []string                `json:"issues,omitempty"`
}

// BiologicalSparsityRange represents the biological sparsity range
type BiologicalSparsityRange struct {
	MinSparsity     float64 `json:"min_sparsity"`
	MaxSparsity     float64 `json:"max_sparsity"`
	OptimalSparsity float64 `json:"optimal_sparsity"`
	Source          string  `json:"source"`
}

// SparsityDataPoint represents a historical sparsity measurement
type SparsityDataPoint struct {
	Timestamp string  `json:"timestamp"`
	Sparsity  float64 `json:"sparsity"`
}

// OverlapAnalysis represents overlap pattern validation analysis
type OverlapAnalysis struct {
	OverlapDistribution []int                    `json:"overlap_distribution"`
	MeanOverlap         float64                  `json:"mean_overlap"`
	OverlapVariance     float64                  `json:"overlap_variance"`
	MinimumOverlap      int                      `json:"minimum_overlap"`
	MaximumOverlap      int                      `json:"maximum_overlap"`
	OverlapValid        bool                     `json:"overlap_valid"`
	PatternSeparation   PatternSeparationMetrics `json:"pattern_separation"`
	Issues              []string                 `json:"issues,omitempty"`
}

// PatternSeparationMetrics represents pattern separation analysis
type PatternSeparationMetrics struct {
	SeparationScore    float64 `json:"separation_score"`
	UniquenessScore    float64 `json:"uniqueness_score"`
	OverlapSimilarity  float64 `json:"overlap_similarity"`
	DistributionHealth string  `json:"distribution_health"`
}

// LearningAnalysis represents learning mechanism validation
type LearningAnalysis struct {
	LearningActive    bool                  `json:"learning_active"`
	LearningRate      float64               `json:"learning_rate"`
	PermanenceUpdates PermanenceMetrics     `json:"permanence_updates"`
	SynapticChanges   SynapticChangeMetrics `json:"synaptic_changes"`
	LearningHealth    string                `json:"learning_health"`
	AdaptationMetrics AdaptationMetrics     `json:"adaptation_metrics"`
	Issues            []string              `json:"issues,omitempty"`
}

// PermanenceMetrics represents permanence update analysis
type PermanenceMetrics struct {
	IncrementRate      float64 `json:"increment_rate"`
	DecrementRate      float64 `json:"decrement_rate"`
	ConnectedThreshold float64 `json:"connected_threshold"`
	PermanenceHealth   string  `json:"permanence_health"`
	UpdateFrequency    float64 `json:"update_frequency"`
}

// SynapticChangeMetrics represents synaptic plasticity analysis
type SynapticChangeMetrics struct {
	NewConnections    int     `json:"new_connections"`
	LostConnections   int     `json:"lost_connections"`
	StableConnections int     `json:"stable_connections"`
	PlasticityRate    float64 `json:"plasticity_rate"`
	StabilityScore    float64 `json:"stability_score"`
}

// AdaptationMetrics represents adaptation mechanism analysis
type AdaptationMetrics struct {
	DutyCycleBalance      float64 `json:"duty_cycle_balance"`
	BoostingEffectiveness float64 `json:"boosting_effectiveness"`
	CompetitionHealth     string  `json:"competition_health"`
	ResourceUtilization   float64 `json:"resource_utilization"`
}

// InhibitionAnalysis represents inhibition mechanism validation
type InhibitionAnalysis struct {
	InhibitionActive   bool                    `json:"inhibition_active"`
	InhibitionType     string                  `json:"inhibition_type"`
	InhibitionStrength float64                 `json:"inhibition_strength"`
	WinnerSelection    WinnerSelectionMetrics  `json:"winner_selection"`
	LocalCompetition   LocalCompetitionMetrics `json:"local_competition"`
	InhibitionHealth   string                  `json:"inhibition_health"`
	Issues             []string                `json:"issues,omitempty"`
}

// WinnerSelectionMetrics represents winner selection analysis
type WinnerSelectionMetrics struct {
	SelectionAccuracy  float64 `json:"selection_accuracy"`
	ConsistencyScore   float64 `json:"consistency_score"`
	CompetitionBalance float64 `json:"competition_balance"`
	WinnerDistribution []int   `json:"winner_distribution"`
}

// LocalCompetitionMetrics represents local competition analysis
type LocalCompetitionMetrics struct {
	NeighborhoodSize    int     `json:"neighborhood_size"`
	CompetitionRadius   float64 `json:"competition_radius"`
	LocalBalance        float64 `json:"local_balance"`
	SpatialDistribution string  `json:"spatial_distribution"`
}

// NetworkTopology represents network structure validation
type NetworkTopology struct {
	InputDimensions     []int               `json:"input_dimensions"`
	ColumnDimensions    []int               `json:"column_dimensions"`
	ConnectivityPattern ConnectivityPattern `json:"connectivity_pattern"`
	TopologyHealth      string              `json:"topology_health"`
	ScalingProperties   ScalingProperties   `json:"scaling_properties"`
	Issues              []string            `json:"issues,omitempty"`
}

// ConnectivityPattern represents connection pattern analysis
type ConnectivityPattern struct {
	PotentialConnections int     `json:"potential_connections"`
	ActiveConnections    int     `json:"active_connections"`
	ConnectivityRatio    float64 `json:"connectivity_ratio"`
	PatternType          string  `json:"pattern_type"`
	DistributionQuality  string  `json:"distribution_quality"`
}

// ScalingProperties represents scalability analysis
type ScalingProperties struct {
	ComputationalComplexity  string  `json:"computational_complexity"`
	MemoryScaling            string  `json:"memory_scaling"`
	ParallelizationPotential float64 `json:"parallelization_potential"`
	ScalabilityScore         float64 `json:"scalability_score"`
}

// TestHTMValidationHealthy tests HTM properties validation when system is healthy
func TestHTMValidationHealthy(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create HTM validation handler
	htmHandler := handlers.NewHTMValidationHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupHTMValidationRoutes(router, htmHandler)

	// Execute: Make HTM validation request
	req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/validation/htm-properties", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected successful HTM validation")

	// Verify: Response body structure
	var response HTMValidationResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Verify: Overall validation
	assert.True(t, response.Valid, "Expected HTM properties to be valid")
	assert.NotEmpty(t, response.Timestamp, "Expected timestamp to be present")

	// Verify: Biological compliance
	compliance := response.BiologicalCompliance
	assert.GreaterOrEqual(t, compliance.OverallScore, 0.8, "Expected high biological compliance score")
	assert.Contains(t, []string{"high", "excellent"}, compliance.ComplianceLevel, "Expected good compliance level")
	assert.Empty(t, compliance.CriticalIssues, "Expected no critical biological issues")

	// Verify: Biological principles adherence
	principles := compliance.BiologicalPrinciples
	assert.True(t, principles.SparseCoding, "Expected sparse coding principle to be followed")
	assert.True(t, principles.LocalInhibition, "Expected local inhibition principle to be followed")
	assert.True(t, principles.HebbianLearning, "Expected Hebbian learning principle to be followed")
	assert.True(t, principles.CompetitiveLearning, "Expected competitive learning principle to be followed")
	assert.True(t, principles.AdaptiveBoosting, "Expected adaptive boosting principle to be followed")

	// Verify: Sparsity analysis
	sparsity := response.SparsityAnalysis
	assert.True(t, sparsity.SparsityValid, "Expected sparsity to be valid")
	assert.GreaterOrEqual(t, sparsity.CurrentSparsity, 0.015, "Expected sparsity >= 1.5%")
	assert.LessOrEqual(t, sparsity.CurrentSparsity, 0.05, "Expected sparsity <= 5%")
	assert.Equal(t, integrationConfig.Application.SpatialPooler.SparsityRatio, sparsity.TargetSparsity, "Expected target sparsity to match config")
	assert.Contains(t, []string{"excellent", "good", "healthy"}, sparsity.SparsityHealth, "Expected good sparsity health")

	// Verify: Biological sparsity range
	biologyRange := sparsity.BiologicalRange
	assert.Equal(t, 0.015, biologyRange.MinSparsity, "Expected biological minimum sparsity of 1.5%")
	assert.Equal(t, 0.05, biologyRange.MaxSparsity, "Expected biological maximum sparsity of 5%")
	assert.GreaterOrEqual(t, biologyRange.OptimalSparsity, biologyRange.MinSparsity, "Expected optimal >= min")
	assert.LessOrEqual(t, biologyRange.OptimalSparsity, biologyRange.MaxSparsity, "Expected optimal <= max")
	assert.NotEmpty(t, biologyRange.Source, "Expected biological range source")

	// Verify: Overlap analysis
	overlap := response.OverlapAnalysis
	assert.True(t, overlap.OverlapValid, "Expected overlap to be valid")
	assert.NotEmpty(t, overlap.OverlapDistribution, "Expected overlap distribution data")
	assert.Greater(t, overlap.MeanOverlap, 0.0, "Expected positive mean overlap")
	assert.GreaterOrEqual(t, overlap.OverlapVariance, 0.0, "Expected non-negative overlap variance")
	assert.GreaterOrEqual(t, overlap.MinimumOverlap, 0, "Expected non-negative minimum overlap")
	assert.GreaterOrEqual(t, overlap.MaximumOverlap, overlap.MinimumOverlap, "Expected max >= min overlap")

	// Verify: Pattern separation metrics
	separation := overlap.PatternSeparation
	assert.GreaterOrEqual(t, separation.SeparationScore, 0.7, "Expected good pattern separation")
	assert.GreaterOrEqual(t, separation.UniquenessScore, 0.7, "Expected good pattern uniqueness")
	assert.Contains(t, []string{"excellent", "good", "healthy"}, separation.DistributionHealth, "Expected good distribution health")

	// Verify: Learning analysis
	learning := response.LearningAnalysis
	assert.True(t, learning.LearningActive, "Expected learning to be active")
	assert.Greater(t, learning.LearningRate, 0.0, "Expected positive learning rate")
	assert.Contains(t, []string{"excellent", "good", "healthy"}, learning.LearningHealth, "Expected good learning health")

	// Verify: Permanence metrics
	permanence := learning.PermanenceUpdates
	assert.Greater(t, permanence.IncrementRate, 0.0, "Expected positive increment rate")
	assert.Greater(t, permanence.DecrementRate, 0.0, "Expected positive decrement rate")
	assert.Greater(t, permanence.ConnectedThreshold, 0.0, "Expected positive connected threshold")
	assert.LessOrEqual(t, permanence.ConnectedThreshold, 1.0, "Expected connected threshold <= 1.0")
	assert.Contains(t, []string{"excellent", "good", "healthy"}, permanence.PermanenceHealth, "Expected good permanence health")

	// Verify: Synaptic changes
	synaptic := learning.SynapticChanges
	assert.GreaterOrEqual(t, synaptic.NewConnections, 0, "Expected non-negative new connections")
	assert.GreaterOrEqual(t, synaptic.LostConnections, 0, "Expected non-negative lost connections")
	assert.GreaterOrEqual(t, synaptic.StableConnections, 0, "Expected non-negative stable connections")
	assert.GreaterOrEqual(t, synaptic.PlasticityRate, 0.0, "Expected non-negative plasticity rate")
	assert.GreaterOrEqual(t, synaptic.StabilityScore, 0.0, "Expected non-negative stability score")
	assert.LessOrEqual(t, synaptic.StabilityScore, 1.0, "Expected stability score <= 1.0")

	// Verify: Adaptation metrics
	adaptation := learning.AdaptationMetrics
	assert.GreaterOrEqual(t, adaptation.DutyCycleBalance, 0.0, "Expected non-negative duty cycle balance")
	assert.GreaterOrEqual(t, adaptation.BoostingEffectiveness, 0.0, "Expected non-negative boosting effectiveness")
	assert.Contains(t, []string{"excellent", "good", "healthy"}, adaptation.CompetitionHealth, "Expected good competition health")
	assert.GreaterOrEqual(t, adaptation.ResourceUtilization, 0.0, "Expected non-negative resource utilization")
	assert.LessOrEqual(t, adaptation.ResourceUtilization, 1.0, "Expected resource utilization <= 1.0")

	// Verify: Inhibition analysis
	inhibition := response.InhibitionAnalysis
	assert.True(t, inhibition.InhibitionActive, "Expected inhibition to be active")
	assert.Contains(t, []string{"global", "local", "hybrid"}, inhibition.InhibitionType, "Expected valid inhibition type")
	assert.Greater(t, inhibition.InhibitionStrength, 0.0, "Expected positive inhibition strength")
	assert.Contains(t, []string{"excellent", "good", "healthy"}, inhibition.InhibitionHealth, "Expected good inhibition health")

	// Verify: Winner selection metrics
	winner := inhibition.WinnerSelection
	assert.GreaterOrEqual(t, winner.SelectionAccuracy, 0.8, "Expected high selection accuracy")
	assert.GreaterOrEqual(t, winner.ConsistencyScore, 0.7, "Expected good consistency")
	assert.GreaterOrEqual(t, winner.CompetitionBalance, 0.7, "Expected good competition balance")
	assert.NotEmpty(t, winner.WinnerDistribution, "Expected winner distribution data")

	// Verify: Network topology
	topology := response.NetworkTopology
	assert.NotEmpty(t, topology.InputDimensions, "Expected input dimensions")
	assert.NotEmpty(t, topology.ColumnDimensions, "Expected column dimensions")
	assert.Contains(t, []string{"excellent", "good", "healthy"}, topology.TopologyHealth, "Expected good topology health")

	// Verify: Connectivity pattern
	connectivity := topology.ConnectivityPattern
	assert.Greater(t, connectivity.PotentialConnections, 0, "Expected positive potential connections")
	assert.Greater(t, connectivity.ActiveConnections, 0, "Expected positive active connections")
	assert.LessOrEqual(t, connectivity.ActiveConnections, connectivity.PotentialConnections, "Expected active <= potential")
	assert.Greater(t, connectivity.ConnectivityRatio, 0.0, "Expected positive connectivity ratio")
	assert.LessOrEqual(t, connectivity.ConnectivityRatio, 1.0, "Expected connectivity ratio <= 1.0")

	// Verify: Scaling properties
	scaling := topology.ScalingProperties
	assert.NotEmpty(t, scaling.ComputationalComplexity, "Expected computational complexity info")
	assert.NotEmpty(t, scaling.MemoryScaling, "Expected memory scaling info")
	assert.GreaterOrEqual(t, scaling.ParallelizationPotential, 0.0, "Expected non-negative parallelization potential")
	assert.GreaterOrEqual(t, scaling.ScalabilityScore, 0.0, "Expected non-negative scalability score")
	assert.LessOrEqual(t, scaling.ScalabilityScore, 1.0, "Expected scalability score <= 1.0")
}

// TestHTMValidationSparsityIssues tests HTM validation when sparsity is outside biological range
func TestHTMValidationSparsityIssues(t *testing.T) {
	// Setup: Create integration configuration with problematic sparsity
	integrationConfig := config.NewDefaultIntegrationConfig()
	integrationConfig.Application.SpatialPooler.SparsityRatio = 0.001 // Too sparse (<0.5%)

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	htmHandler := handlers.NewHTMValidationHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupHTMValidationRoutes(router, htmHandler)

	// Execute: Make validation request
	req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/validation/htm-properties", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status (should still be 200, but validation should show issues)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected response even with HTM issues")

	// Verify: Response body
	var response HTMValidationResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Verify: Validation should show issues with sparsity
	sparsity := response.SparsityAnalysis
	assert.False(t, sparsity.SparsityValid, "Expected sparsity validation to fail for too sparse configuration")
	assert.NotEmpty(t, sparsity.Issues, "Expected sparsity issues to be reported")

	// Verify: Overall biological compliance should be lower
	assert.LessOrEqual(t, response.BiologicalCompliance.OverallScore, 0.7, "Expected lower compliance score with sparsity issues")

	// Verify: Recommendations should be provided
	assert.NotEmpty(t, response.Recommendations, "Expected recommendations for fixing sparsity issues")
}

// TestHTMValidationBiologicalPrinciples tests specific biological principle validation
func TestHTMValidationBiologicalPrinciples(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	htmHandler := handlers.NewHTMValidationHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupHTMValidationRoutes(router, htmHandler)

	// Execute: Make validation request
	req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/validation/htm-properties", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Verify: Biological principles validation
	var response HTMValidationResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	principles := response.BiologicalCompliance.BiologicalPrinciples

	// Verify: Core HTM principles are implemented
	assert.True(t, principles.SparseCoding, "Expected sparse coding to be implemented")
	assert.True(t, principles.CompetitiveLearning, "Expected competitive learning to be implemented")

	// Verify: Learning mechanisms follow biological principles
	assert.True(t, principles.HebbianLearning, "Expected Hebbian learning to be implemented")
	assert.True(t, principles.AdaptiveBoosting, "Expected adaptive boosting to be implemented")

	// Verify: Spatial processing follows biological principles
	if integrationConfig.Application.SpatialPooler.GlobalInhibition {
		// Global inhibition should still respect local computation principles
		assert.True(t, principles.LocalInhibition, "Expected local inhibition principles even with global inhibition")
	}
}
