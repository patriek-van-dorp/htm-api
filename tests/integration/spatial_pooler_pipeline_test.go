package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSpatialPoolerPipeline tests the complete encoder-to-spatial-pooler pipeline integration
func TestSpatialPoolerPipeline(t *testing.T) {
	t.Skip("Skipping until spatial pooler and sensor integration is complete")

	// This test validates the complete pipeline: Raw Input → Sensor Encoding → Spatial Pooling → SDR Response

	t.Run("end_to_end_categorical_pipeline", func(t *testing.T) {
		// Test complete pipeline with categorical sensor encoder

		// Setup: Categorical encoder + Spatial pooler
		// categoricalEncoder := setupCategoricalEncoder(categories=["red", "green", "blue"], width=2048)
		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		// Test input: categorical value
		inputValue := "red"

		// Step 1: Encode using categorical encoder
		// encoderOutput := categoricalEncoder.Encode(inputValue)
		// assert.Equal(t, 2048, encoderOutput.Width, "Encoder output should have correct width")
		// assert.Greater(t, len(encoderOutput.ActiveBits), 0, "Encoder should produce active bits")

		// Step 2: Process through spatial pooler
		// poolingInput := createPoolingInput(encoderOutput, "categorical-test")
		// poolingResult := spatialPooler.Process(poolingInput)

		// Step 3: Validate pipeline output
		// assert.Equal(t, "categorical-test", poolingResult.InputID, "Input ID should be preserved")
		// assert.Equal(t, 2048, poolingResult.NormalizedSDR.Width, "SDR width should match encoder width")
		// assert.GreaterOrEqual(t, poolingResult.SparsityLevel, 0.02, "Sparsity should be >= 2%")
		// assert.LessOrEqual(t, poolingResult.SparsityLevel, 0.05, "Sparsity should be <= 5%")
		// assert.LessOrEqual(t, poolingResult.ProcessingTime.Milliseconds(), int64(10), "Processing should be <= 10ms")

		// Validate that spatial pooler normalized the encoder output
		// encoderSparsity := float64(len(encoderOutput.ActiveBits)) / float64(encoderOutput.Width)
		// if encoderSparsity < 0.02 || encoderSparsity > 0.05 {
		//     // If encoder output was outside HTM range, spatial pooler should have normalized it
		//     assert.NotEqual(t, encoderOutput.ActiveBits, poolingResult.NormalizedSDR.ActiveBits,
		//         "Spatial pooler should normalize encoder output when outside HTM sparsity range")
		// }
	})

	t.Run("end_to_end_numeric_pipeline", func(t *testing.T) {
		// Test complete pipeline with numeric sensor encoder

		// Setup: Numeric encoder + Spatial pooler
		// numericEncoder := setupNumericEncoder(min=0, max=100, width=2048, resolution=1.0)
		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		testValues := []float64{25.5, 50.0, 75.5}

		for i, value := range testValues {
			// Step 1: Encode using numeric encoder
			// encoderOutput := numericEncoder.Encode(value)

			// Step 2: Process through spatial pooler
			// poolingInput := createPoolingInput(encoderOutput, fmt.Sprintf("numeric-test-%d", i))
			// poolingResult := spatialPooler.Process(poolingInput)

			// Step 3: Validate pipeline output
			// assert.Equal(t, fmt.Sprintf("numeric-test-%d", i), poolingResult.InputID)
			// assert.GreaterOrEqual(t, poolingResult.SparsityLevel, 0.02, "Value %.1f: Sparsity should be >= 2%", value)
			// assert.LessOrEqual(t, poolingResult.SparsityLevel, 0.05, "Value %.1f: Sparsity should be <= 5%", value)
			// assert.LessOrEqual(t, poolingResult.ProcessingTime.Milliseconds(), int64(10), "Value %.1f: Processing should be <= 10ms", value)
		}
	})

	t.Run("pipeline_with_different_encoders", func(t *testing.T) {
		// Test that spatial pooler can handle outputs from different encoder types

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		// Setup different encoders
		// categoricalEncoder := setupCategoricalEncoder(categories=["A", "B", "C"], width=2048)
		// numericEncoder := setupNumericEncoder(min=0, max=100, width=2048, resolution=1.0)
		// dateTimeEncoder := setupDateTimeEncoder(width=2048)

		// Test inputs
		testCases := []struct {
			name      string
			inputType string
			value     interface{}
		}{
			{"categorical_input", "categorical", "A"},
			{"numeric_input", "numeric", 42.0},
			{"datetime_input", "datetime", time.Now()},
		}

		var results []TestResult
		for _, tc := range testCases {
			var encoderOutput EncoderOutput

			switch tc.inputType {
			case "categorical":
				// encoderOutput = categoricalEncoder.Encode(tc.value.(string))
			case "numeric":
				// encoderOutput = numericEncoder.Encode(tc.value.(float64))
			case "datetime":
				// encoderOutput = dateTimeEncoder.Encode(tc.value.(time.Time))
			}

			// poolingInput := createPoolingInput(encoderOutput, tc.name)
			// result := spatialPooler.Process(poolingInput)
			// results = append(results, result)

			// Validate that each encoder type produces valid HTM-compliant SDRs
			// assert.GreaterOrEqual(t, result.SparsityLevel, 0.02, "%s: Sparsity should be >= 2%", tc.name)
			// assert.LessOrEqual(t, result.SparsityLevel, 0.05, "%s: Sparsity should be <= 5%", tc.name)
		}

		// All results should be valid HTM SDRs regardless of encoder type
		for _, result := range results {
			// assert.Equal(t, 2048, result.NormalizedSDR.Width, "All outputs should have consistent width")
			// assert.Greater(t, len(result.NormalizedSDR.ActiveBits), 0, "All outputs should have active bits")
		}
	})

	t.Run("pipeline_error_handling", func(t *testing.T) {
		// Test pipeline error handling for invalid encoder outputs

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		errorCases := []struct {
			name          string
			encoderOutput EncoderOutput
			expectedError string
		}{
			{
				name: "oversized_input",
				encoderOutput: EncoderOutput{
					Width:      4096, // Exceeds expected 2048
					ActiveBits: []int{1, 2, 3},
					Sparsity:   0.001,
				},
				expectedError: "invalid_input",
			},
			{
				name: "empty_active_bits",
				encoderOutput: EncoderOutput{
					Width:      2048,
					ActiveBits: []int{}, // Empty active bits
					Sparsity:   0.0,
				},
				expectedError: "invalid_input",
			},
			{
				name: "invalid_active_bits",
				encoderOutput: EncoderOutput{
					Width:      2048,
					ActiveBits: []int{-1, 2048, 3000}, // Out of range bits
					Sparsity:   0.001,
				},
				expectedError: "invalid_input",
			},
		}

		for _, tc := range errorCases {
			t.Run(tc.name, func(t *testing.T) {
				// poolingInput := createPoolingInput(tc.encoderOutput, tc.name)

				// Should return error instead of result
				// result, err := spatialPooler.ProcessWithError(poolingInput)
				// assert.Error(t, err, "%s should produce an error", tc.name)
				// assert.Nil(t, result, "%s should not produce a result on error", tc.name)

				// Error should be of expected type
				// if poolingErr, ok := err.(*PoolingError); ok {
				//     assert.Equal(t, tc.expectedError, poolingErr.ErrorType, "%s should produce expected error type", tc.name)
				// }
			})
		}
	})

	t.Run("pipeline_performance_requirements", func(t *testing.T) {
		// Test that the complete pipeline meets performance requirements

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)
		// categoricalEncoder := setupCategoricalEncoder(categories=["test"], width=2048)

		// Measure performance over multiple iterations
		iterationCount := 100
		var totalTime time.Duration

		for i := 0; i < iterationCount; i++ {
			start := time.Now()

			// Complete pipeline: encode + spatial pool
			// encoderOutput := categoricalEncoder.Encode("test")
			// poolingInput := createPoolingInput(encoderOutput, fmt.Sprintf("perf-test-%d", i))
			// result := spatialPooler.Process(poolingInput)

			elapsed := time.Since(start)
			totalTime += elapsed

			// Individual processing should be under 10ms
			// assert.LessOrEqual(t, result.ProcessingTime.Milliseconds(), int64(10),
			//     "Iteration %d: Spatial pooling should be <= 10ms", i)
		}

		averageTime := totalTime / time.Duration(iterationCount)

		// Average complete pipeline should be well under 100ms total API response time budget
		assert.LessOrEqual(t, averageTime.Milliseconds(), int64(50),
			"Average complete pipeline should be <= 50ms (leaving room for other API processing)")
	})

	t.Run("pipeline_semantic_preservation", func(t *testing.T) {
		// Test that semantic relationships are preserved through the complete pipeline

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)
		// categoricalEncoder := setupCategoricalEncoder(categories=["red", "green", "blue", "crimson", "navy"], width=2048)

		// Test similar and different categorical values
		similarInputs := []string{"red", "crimson"} // Similar colors
		differentInputs := []string{"red", "blue"}  // Different colors

		// Process similar inputs through complete pipeline
		var similarResults []TestResult
		for _, input := range similarInputs {
			// encoderOutput := categoricalEncoder.Encode(input)
			// poolingInput := createPoolingInput(encoderOutput, input)
			// result := spatialPooler.Process(poolingInput)
			// similarResults = append(similarResults, result)
		}

		// Process different inputs through complete pipeline
		var differentResults []TestResult
		for _, input := range differentInputs {
			// encoderOutput := categoricalEncoder.Encode(input)
			// poolingInput := createPoolingInput(encoderOutput, input)
			// result := spatialPooler.Process(poolingInput)
			// differentResults = append(differentResults, result)
		}

		// Calculate overlaps
		// similarOverlap := calculateSDROverlapPercentage(similarResults[0].NormalizedSDR, similarResults[1].NormalizedSDR)
		// differentOverlap := calculateSDROverlapPercentage(differentResults[0].NormalizedSDR, differentResults[1].NormalizedSDR)

		// Semantic relationships should be preserved through the pipeline
		// assert.Greater(t, similarOverlap, differentOverlap, "Similar inputs should have higher overlap than different inputs")
		// assert.GreaterOrEqual(t, similarOverlap, 0.30, "Similar inputs should have >= 30% overlap")
		// assert.Less(t, differentOverlap, 0.20, "Different inputs should have < 20% overlap")
	})
}

// TestSpatialPoolerIntegrationWithExistingSensors tests integration with existing sensor package
func TestSpatialPoolerIntegrationWithExistingSensors(t *testing.T) {
	t.Skip("Skipping until sensor integration is complete")

	t.Run("integration_with_existing_sensor_api", func(t *testing.T) {
		// Test integration with existing sensor API endpoints

		// This test assumes existing sensor endpoints that produce SDR outputs
		// The spatial pooler should be able to consume these existing outputs

		// Setup HTTP client for sensor API
		// sensorClient := setupSensorAPIClient()
		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		// Call existing sensor endpoint
		// sensorResponse := sensorClient.ProcessCategorical("test-value")
		// assert.Equal(t, 200, sensorResponse.StatusCode, "Sensor API should respond successfully")

		// Extract SDR from sensor response
		// sensorSDR := extractSDRFromSensorResponse(sensorResponse)

		// Convert sensor SDR to spatial pooler input format
		// poolingInput := createPoolingInputFromSensorSDR(sensorSDR, "sensor-integration-test")

		// Process through spatial pooler
		// result := spatialPooler.Process(poolingInput)

		// Validate integration
		// assert.Equal(t, "sensor-integration-test", result.InputID)
		// assert.GreaterOrEqual(t, result.SparsityLevel, 0.02, "Integrated output should have >= 2% sparsity")
		// assert.LessOrEqual(t, result.SparsityLevel, 0.05, "Integrated output should have <= 5% sparsity")
	})

	t.Run("backward_compatibility_with_sensor_outputs", func(t *testing.T) {
		// Test that spatial pooler maintains backward compatibility with existing sensor outputs

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		// Test with various sensor output formats that might exist
		existingSensorOutputs := []SensorOutput{
			// createExistingSensorOutput("categorical", "red"),
			// createExistingSensorOutput("numeric", 42.0),
			// createExistingSensorOutput("datetime", time.Now()),
		}

		for _, sensorOutput := range existingSensorOutputs {
			// Convert to spatial pooler format
			// poolingInput := convertSensorOutputToPoolingInput(sensorOutput)

			// Should process without errors
			// result := spatialPooler.Process(poolingInput)

			// Should produce valid HTM SDR
			// assert.GreaterOrEqual(t, result.SparsityLevel, 0.02, "Should normalize sensor output to HTM sparsity")
			// assert.LessOrEqual(t, result.SparsityLevel, 0.05, "Should normalize sensor output to HTM sparsity")
		}
	})
}

// Helper types and functions for pipeline testing

type SensorOutput struct {
	Type  string
	Value interface{}
	SDR   SDR
}

func createPoolingInput(encoderOutput EncoderOutput, inputID string) interface{} {
	// Convert encoder output to pooling input format
	// This would be implemented based on actual types
	return struct {
		EncoderOutput EncoderOutput
		InputID       string
	}{
		EncoderOutput: encoderOutput,
		InputID:       inputID,
	}
}

func createPoolingInputFromSensorSDR(sdr SDR, inputID string) interface{} {
	// Convert existing sensor SDR to pooling input format
	return createPoolingInput(EncoderOutput{
		Width:      sdr.Width,
		ActiveBits: sdr.ActiveBits,
		Sparsity:   sdr.Sparsity,
	}, inputID)
}

func extractSDRFromSensorResponse(response interface{}) SDR {
	// Extract SDR from sensor API response
	// This would be implemented based on actual sensor API format
	return SDR{
		Width:      2048,
		ActiveBits: []int{1, 2, 3, 4, 5},
		Sparsity:   0.00244,
	}
}

func convertSensorOutputToPoolingInput(sensorOutput SensorOutput) interface{} {
	// Convert existing sensor output format to spatial pooler input
	return createPoolingInputFromSensorSDR(sensorOutput.SDR, "converted-input")
}
