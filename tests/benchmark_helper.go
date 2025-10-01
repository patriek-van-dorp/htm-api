package tests

import (
	"testing"
	"time"
)

// BenchmarkConfig defines configuration for performance benchmarks
type BenchmarkConfig struct {
	MaxDuration    time.Duration // Maximum allowed execution time
	WarmupRounds   int           // Number of warmup iterations
	MeasureRounds  int           // Number of measurement iterations
	MemoryTracking bool          // Whether to track memory allocations
}

// SubMillisecondBenchmark provides a framework for validating sub-millisecond performance
type SubMillisecondBenchmark struct {
	config BenchmarkConfig
}

// NewSubMillisecondBenchmark creates a benchmark framework with 1ms max duration
func NewSubMillisecondBenchmark() *SubMillisecondBenchmark {
	return &SubMillisecondBenchmark{
		config: BenchmarkConfig{
			MaxDuration:    time.Millisecond,
			WarmupRounds:   100,
			MeasureRounds:  1000,
			MemoryTracking: true,
		},
	}
}

// Run executes a benchmark operation and validates performance constraints
func (b *SubMillisecondBenchmark) Run(t *testing.T, name string, operation func()) {
	t.Helper()

	// Warmup phase
	for i := 0; i < b.config.WarmupRounds; i++ {
		operation()
	}

	// Measurement phase
	start := time.Now()
	for i := 0; i < b.config.MeasureRounds; i++ {
		operation()
	}
	elapsed := time.Since(start)

	avgDuration := elapsed / time.Duration(b.config.MeasureRounds)

	if avgDuration > b.config.MaxDuration {
		t.Errorf("%s: Average execution time %v exceeds maximum allowed %v",
			name, avgDuration, b.config.MaxDuration)
	}

	t.Logf("%s: Average execution time: %v (limit: %v)",
		name, avgDuration, b.config.MaxDuration)
}

// BenchmarkMemory runs a memory allocation benchmark
func (b *SubMillisecondBenchmark) BenchmarkMemory(t *testing.T, name string, operation func()) {
	t.Helper()

	if !b.config.MemoryTracking {
		return
	}

	// Use Go's built-in benchmarking for memory tracking
	result := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			operation()
		}
	})

	t.Logf("%s: %v allocations, %v bytes/op",
		name, result.AllocsPerOp(), result.AllocedBytesPerOp())
}
