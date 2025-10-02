package integration

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSpatialPoolerPerformance tests that spatial pooler meets performance requirements
func TestSpatialPoolerPerformance(t *testing.T) {
	t.Skip("Skipping until spatial pooler implementation is complete")

	// This test validates FR-010: System MUST provide processing performance under 10ms
	// to maintain overall API response time targets under 100ms

	t.Run("single_processing_under_10ms", func(t *testing.T) {
		// Test that individual spatial pooling operations complete under 10ms

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		testCases := []struct {
			name        string
			inputSize   int
			activeCount int
		}{
			{"small_input", 1024, 5},
			{"medium_input", 2048, 10},
			{"large_input", 4096, 20},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := createTestEncoderOutput(generateSequentialBits(tc.activeCount), tc.inputSize)

				// Measure processing time
				start := time.Now()
				// result := spatialPooler.Process(input)
				elapsed := time.Since(start)

				// Processing should be under 10ms
				assert.LessOrEqual(t, elapsed.Milliseconds(), int64(10),
					"%s: Processing time should be <= 10ms, got %dms", tc.name, elapsed.Milliseconds())

				// Reported processing time should also be under 10ms
				// assert.LessOrEqual(t, result.ProcessingTime.Milliseconds(), int64(10),
				//     "%s: Reported processing time should be <= 10ms", tc.name)
			})
		}
	})

	t.Run("batch_processing_performance", func(t *testing.T) {
		// Test performance with multiple inputs processed in sequence

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		batchSize := 100
		inputs := make([]TestInput, batchSize)
		for i := 0; i < batchSize; i++ {
			inputs[i] = createTestEncoderOutput(generateSequentialBits(5), 2048)
		}

		start := time.Now()
		var totalProcessingTime time.Duration

		for i, input := range inputs {
			iterStart := time.Now()
			// result := spatialPooler.Process(input)
			iterElapsed := time.Since(iterStart)

			// Each individual operation should still be under 10ms
			assert.LessOrEqual(t, iterElapsed.Milliseconds(), int64(10),
				"Batch item %d: Processing should be <= 10ms even in batch", i)

			// totalProcessingTime += result.ProcessingTime
		}

		totalElapsed := time.Since(start)
		avgTime := totalElapsed / time.Duration(batchSize)

		// Average time should be well under 10ms
		assert.LessOrEqual(t, avgTime.Milliseconds(), int64(8),
			"Average processing time should be <= 8ms in batch processing")

		// Total time should be reasonable for the batch size
		expectedMaxTotal := time.Duration(batchSize) * 10 * time.Millisecond
		assert.LessOrEqual(t, totalElapsed, expectedMaxTotal,
			"Total batch processing should not exceed %v", expectedMaxTotal)
	})

	t.Run("performance_with_learning_enabled", func(t *testing.T) {
		// Test that performance is maintained even with learning enabled

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=true)

		testInput := createTestEncoderOutput([]int{10, 20, 30, 40, 50}, 2048)

		// Measure performance over multiple iterations (learning may slow down initially)
		iterationCount := 50
		var processingTimes []time.Duration

		for i := 0; i < iterationCount; i++ {
			start := time.Now()
			// result := spatialPooler.Process(testInput)
			elapsed := time.Since(start)

			processingTimes = append(processingTimes, elapsed)

			// Even with learning, should still be under 10ms
			assert.LessOrEqual(t, elapsed.Milliseconds(), int64(10),
				"Iteration %d with learning: Processing should be <= 10ms", i)
		}

		// Calculate statistics
		var total time.Duration
		for _, duration := range processingTimes {
			total += duration
		}
		average := total / time.Duration(len(processingTimes))

		// Average should be well under the limit
		assert.LessOrEqual(t, average.Milliseconds(), int64(7),
			"Average processing time with learning should be <= 7ms")
	})

	t.Run("performance_with_different_configurations", func(t *testing.T) {
		// Test performance with various spatial pooler configurations

		configs := []struct {
			name             string
			columnCount      int
			globalInhibition bool
			learningEnabled  bool
		}{
			{"small_columns_global", 1024, true, false},
			{"medium_columns_global", 2048, true, false},
			{"large_columns_global", 4096, true, false},
			{"medium_columns_local", 2048, false, false},
			{"medium_learning", 2048, true, true},
		}

		for _, config := range configs {
			t.Run(config.name, func(t *testing.T) {
				// spatialPooler := setupSpatialPoolerWithConfig(config)

				testInput := createTestEncoderOutput([]int{1, 2, 3, 4, 5}, 2048)

				// Measure performance
				start := time.Now()
				// result := spatialPooler.Process(testInput)
				elapsed := time.Since(start)

				// All configurations should meet performance requirements
				assert.LessOrEqual(t, elapsed.Milliseconds(), int64(10),
					"Config %s: Processing should be <= 10ms", config.name)
			})
		}
	})

	t.Run("throughput_performance", func(t *testing.T) {
		// Test that spatial pooler can handle required throughput (1000-5000 req/sec)
		// This validates FR-013: System MUST handle throughput of 1,000-5,000 requests per second

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		testInput := createTestEncoderOutput([]int{1, 2, 3, 4, 5}, 2048)

		// Test for 1 second duration
		duration := 1 * time.Second
		start := time.Now()
		requestCount := 0

		for time.Since(start) < duration {
			iterStart := time.Now()
			// spatialPooler.Process(testInput)
			iterElapsed := time.Since(iterStart)

			// Each request should still meet individual performance requirements
			assert.LessOrEqual(t, iterElapsed.Milliseconds(), int64(10),
				"Request %d: Individual processing should be <= 10ms", requestCount)

			requestCount++
		}

		actualDuration := time.Since(start)
		requestsPerSecond := float64(requestCount) / actualDuration.Seconds()

		// Should achieve at least 1000 requests per second
		assert.GreaterOrEqual(t, requestsPerSecond, 1000.0,
			"Should achieve >= 1000 requests/second, got %.1f", requestsPerSecond)

		// Should ideally achieve 5000 requests per second for good performance
		t.Logf("Achieved throughput: %.1f requests/second", requestsPerSecond)
		if requestsPerSecond >= 5000.0 {
			t.Logf("✓ Excellent: Achieved target of 5000+ requests/second")
		} else if requestsPerSecond >= 2500.0 {
			t.Logf("✓ Good: Achieved mid-range performance")
		} else {
			t.Logf("⚠ Minimum: Achieved minimum requirement of 1000+ requests/second")
		}
	})

	t.Run("concurrent_processing_performance", func(t *testing.T) {
		// Test performance under concurrent load

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		concurrentWorkers := 10
		requestsPerWorker := 100

		var wg sync.WaitGroup
		var mu sync.Mutex
		var totalRequests int
		var maxProcessingTime time.Duration

		start := time.Now()

		for worker := 0; worker < concurrentWorkers; worker++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				for req := 0; req < requestsPerWorker; req++ {
					testInput := createTestEncoderOutput([]int{workerID, req, req + 1, req + 2, req + 3}, 2048)

					reqStart := time.Now()
					// spatialPooler.Process(testInput)
					reqElapsed := time.Since(reqStart)

					mu.Lock()
					totalRequests++
					if reqElapsed > maxProcessingTime {
						maxProcessingTime = reqElapsed
					}
					mu.Unlock()

					// Each request should still meet performance requirements under concurrent load
					assert.LessOrEqual(t, reqElapsed.Milliseconds(), int64(15), // Slightly relaxed under concurrent load
						"Worker %d, Request %d: Concurrent processing should be <= 15ms", workerID, req)
				}
			}(worker)
		}

		wg.Wait()
		totalElapsed := time.Since(start)

		// Calculate concurrent throughput
		concurrentThroughput := float64(totalRequests) / totalElapsed.Seconds()

		// Should maintain good performance under concurrent load
		assert.GreaterOrEqual(t, concurrentThroughput, 800.0, // Slightly reduced expectation under concurrency
			"Concurrent throughput should be >= 800 req/sec, got %.1f", concurrentThroughput)

		assert.LessOrEqual(t, maxProcessingTime.Milliseconds(), int64(15),
			"Maximum processing time under concurrent load should be <= 15ms")

		t.Logf("Concurrent performance: %.1f req/sec with %d workers", concurrentThroughput, concurrentWorkers)
		t.Logf("Maximum individual processing time: %dms", maxProcessingTime.Milliseconds())
	})

	t.Run("memory_usage_performance", func(t *testing.T) {
		// Test that memory usage remains reasonable and doesn't impact performance

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		// Process many different inputs to test memory behavior
		inputCount := 1000
		var processingTimes []time.Duration

		for i := 0; i < inputCount; i++ {
			// Create varied inputs to test memory allocation patterns
			activeBits := generateVariedActiveBits(i, 5, 2048)
			testInput := createTestEncoderOutput(activeBits, 2048)

			start := time.Now()
			// spatialPooler.Process(testInput)
			elapsed := time.Since(start)

			processingTimes = append(processingTimes, elapsed)

			// Performance should not degrade over time due to memory issues
			if i > 0 && i%100 == 0 {
				// Check recent average vs initial average
				recentAvg := calculateAverage(processingTimes[i-100:])
				initialAvg := calculateAverage(processingTimes[0:100])

				// Recent performance should not be significantly worse than initial
				assert.LessOrEqual(t, recentAvg.Milliseconds(), initialAvg.Milliseconds()+2,
					"Performance should not degrade significantly over %d iterations", i)
			}
		}

		// Overall average should still meet requirements
		overallAvg := calculateAverage(processingTimes)
		assert.LessOrEqual(t, overallAvg.Milliseconds(), int64(8),
			"Overall average processing time should be <= 8ms over %d iterations", inputCount)
	})
}

// TestSpatialPoolerPerformanceRegression tests for performance regressions
func TestSpatialPoolerPerformanceRegression(t *testing.T) {
	t.Skip("Skipping until spatial pooler implementation is complete")

	t.Run("baseline_performance_benchmark", func(t *testing.T) {
		// Establish baseline performance metrics for regression testing

		// spatialPooler := setupSpatialPooler(deterministic=true, learning=false)

		// Standard test input
		standardInput := createTestEncoderOutput([]int{10, 20, 30, 40, 50}, 2048)

		// Warm up
		for i := 0; i < 10; i++ {
			// spatialPooler.Process(standardInput)
		}

		// Measure baseline
		iterationCount := 1000
		start := time.Now()

		for i := 0; i < iterationCount; i++ {
			// spatialPooler.Process(standardInput)
		}

		elapsed := time.Since(start)
		avgTime := elapsed / time.Duration(iterationCount)

		// Log baseline for future regression testing
		t.Logf("Baseline performance: %.3fms average over %d iterations",
			float64(avgTime.Microseconds())/1000.0, iterationCount)

		// Ensure baseline meets requirements
		assert.LessOrEqual(t, avgTime.Milliseconds(), int64(5),
			"Baseline performance should be <= 5ms for regression testing")
	})
}

// Helper functions for performance testing

func generateSequentialBits(count int) []int {
	bits := make([]int, count)
	for i := 0; i < count; i++ {
		bits[i] = i * 10
	}
	return bits
}

func generateVariedActiveBits(seed, count, width int) []int {
	bits := make([]int, count)
	for i := 0; i < count; i++ {
		// Generate pseudo-random but deterministic positions
		bits[i] = (seed*13 + i*17) % width
	}
	return bits
}

func calculateAverage(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

// setupSpatialPoolerWithConfig is a helper to create spatial pooler with specific configuration
func setupSpatialPoolerWithConfig(config struct {
	name             string
	columnCount      int
	globalInhibition bool
	learningEnabled  bool
}) interface{} {
	// This would return a configured spatial pooler instance
	// Implementation depends on actual spatial pooler API
	return nil
}
