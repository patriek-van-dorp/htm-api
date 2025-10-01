# Quickstart: HTM Sensor Package

## Overview
This quickstart demonstrates how to use the Generic HTM Sensor Package to convert various data types into Sparsely Distributed Representations (SDRs) for Hierarchical Temporal Memory algorithms.

## Prerequisites
- Go 1.21+ installed
- Basic understanding of HTM theory concepts
- Familiarity with Go interfaces and error handling

## Installation
```bash
go get github.com/patriek-van-dorp/htm-api/pkg/sensors
```

## Basic Usage Example

### 1. Setup and Registration
```go
package main

import (
    "fmt"
    "log"
    
    "github.com/patriek-van-dorp/htm-api/pkg/sensors"
    "github.com/patriek-van-dorp/htm-api/internal/sensors/encoders"
)

func main() {
    // Create sensor registry
    registry := sensors.NewRegistry()
    
    // Register built-in sensor types
    registry.Register("numeric", encoders.NewNumericSensor)
    registry.Register("categorical", encoders.NewCategoricalSensor)
    registry.Register("text", encoders.NewTextSensor)
    registry.Register("spatial", encoders.NewSpatialSensor)
    
    fmt.Println("Registered sensors:", registry.List())
}
```

### 2. Numeric Data Encoding
```go
func encodeNumericData(registry *sensors.Registry) {
    // Configure numeric sensor
    config := sensors.NewConfig()
    config.SetParam("sdr_width", 2048)
    config.SetParam("target_sparsity", 0.02)  // 2% sparsity
    config.SetParam("min_value", 0.0)
    config.SetParam("max_value", 100.0)
    config.SetParam("resolution", 0.1)
    
    // Create sensor instance
    sensor, err := registry.Create("numeric", config)
    if err != nil {
        log.Fatal("Failed to create numeric sensor:", err)
    }
    
    // Encode values
    values := []float64{25.5, 50.0, 75.3, 25.5} // Note: duplicate for consistency test
    
    for i, value := range values {
        sdr, err := sensor.Encode(value)
        if err != nil {
            log.Printf("Encoding failed for value %f: %v", value, err)
            continue
        }
        
        fmt.Printf("Value %f -> SDR (width: %d, sparsity: %.3f, active bits: %d)\n",
            value, sdr.Width(), sdr.Sparsity(), len(sdr.ActiveBits()))
        
        // Validate SDR properties
        if err := sdr.Validate(); err != nil {
            log.Printf("SDR validation failed: %v", err)
        }
    }
    
    // Test consistency: same input should produce same SDR
    sdr1, _ := sensor.Encode(25.5)
    sdr2, _ := sensor.Encode(25.5)
    
    if sdr1.Overlap(sdr2) == len(sdr1.ActiveBits()) {
        fmt.Println("✓ Consistency test passed: identical inputs produce identical SDRs")
    } else {
        fmt.Println("✗ Consistency test failed")
    }
}
```

### 3. Categorical Data Encoding
```go
func encodeCategoricalData(registry *sensors.Registry) {
    // Configure categorical sensor
    config := sensors.NewConfig()
    config.SetParam("sdr_width", 2048)
    config.SetParam("target_sparsity", 0.02)
    config.SetParam("categories", []string{"red", "green", "blue", "yellow", "orange"})
    config.SetParam("bucket_size", 40) // Bits per category
    
    // Create sensor
    sensor, err := registry.Create("categorical", config)
    if err != nil {
        log.Fatal("Failed to create categorical sensor:", err)
    }
    
    // Encode categories
    categories := []string{"red", "blue", "green", "red"}
    
    var sdrs []sensors.SDR
    for _, category := range categories {
        sdr, err := sensor.Encode(category)
        if err != nil {
            log.Printf("Encoding failed for category %s: %v", category, err)
            continue
        }
        
        sdrs = append(sdrs, sdr)
        fmt.Printf("Category '%s' -> SDR (active bits: %v)\n", 
            category, sdr.ActiveBits()[:5]) // Show first 5 active bits
    }
    
    // Test semantic similarity: different categories should have low overlap
    redSDR := sdrs[0]
    blueSDR := sdrs[1]
    similarity := redSDR.Similarity(blueSDR)
    
    fmt.Printf("Similarity between 'red' and 'blue': %.3f\n", similarity)
    if similarity < 0.1 {
        fmt.Println("✓ Semantic separation test passed")
    }
}
```

### 4. Text Data Encoding
```go
func encodeTextData(registry *sensors.Registry) {
    // Configure text sensor
    config := sensors.NewConfig()
    config.SetParam("sdr_width", 2048)
    config.SetParam("target_sparsity", 0.03)
    config.SetParam("vocabulary_size", 1000)
    config.SetParam("context_window", 3)
    config.SetParam("tokenizer", "word") // or "character"
    
    // Create sensor
    sensor, err := registry.Create("text", config)
    if err != nil {
        log.Fatal("Failed to create text sensor:", err)
    }
    
    // Encode text samples
    texts := []string{
        "hello world",
        "hello universe", 
        "goodbye world",
        "hello world", // Duplicate for consistency
    }
    
    var sdrs []sensors.SDR
    for _, text := range texts {
        sdr, err := sensor.Encode(text)
        if err != nil {
            log.Printf("Encoding failed for text '%s': %v", text, err)
            continue
        }
        
        sdrs = append(sdrs, sdr)
        fmt.Printf("Text '%s' -> SDR (sparsity: %.3f)\n", text, sdr.Sparsity())
    }
    
    // Test semantic similarity
    helloWorld1 := sdrs[0]
    helloUniverse := sdrs[1]
    goodbyeWorld := sdrs[2]
    helloWorld2 := sdrs[3]
    
    fmt.Printf("Similarity 'hello world' vs 'hello universe': %.3f\n", 
        helloWorld1.Similarity(helloUniverse))
    fmt.Printf("Similarity 'hello world' vs 'goodbye world': %.3f\n", 
        helloWorld1.Similarity(goodbyeWorld))
    fmt.Printf("Consistency 'hello world' duplicate: %.3f\n", 
        helloWorld1.Similarity(helloWorld2))
}
```

### 5. Multi-Sensor Pipeline
```go
func multiSensorPipeline(registry *sensors.Registry) {
    // Create different sensors for different data types
    numericConfig := sensors.NewConfig()
    numericConfig.SetParam("sdr_width", 1024)
    numericConfig.SetParam("target_sparsity", 0.02)
    numericConfig.SetParam("min_value", 0.0)
    numericConfig.SetParam("max_value", 1.0)
    
    textConfig := sensors.NewConfig()
    textConfig.SetParam("sdr_width", 1024)
    textConfig.SetParam("target_sparsity", 0.02)
    textConfig.SetParam("vocabulary_size", 500)
    
    numericSensor, _ := registry.Create("numeric", numericConfig)
    textSensor, _ := registry.Create("text", textConfig)
    
    // Process mixed data types
    dataPoints := []struct {
        value float64
        label string
    }{
        {0.25, "low"},
        {0.75, "high"},
        {0.50, "medium"},
    }
    
    fmt.Println("\nMulti-sensor processing:")
    for _, dp := range dataPoints {
        numericSDR, _ := numericSensor.Encode(dp.value)
        textSDR, _ := textSensor.Encode(dp.label)
        
        fmt.Printf("Value: %.2f, Label: %s\n", dp.value, dp.label)
        fmt.Printf("  Numeric SDR sparsity: %.3f\n", numericSDR.Sparsity())
        fmt.Printf("  Text SDR sparsity: %.3f\n", textSDR.Sparsity())
        
        // Combined processing would concatenate or process SDRs further
        // This demonstrates how different sensors can work together
    }
}
```

### 6. Custom Sensor Implementation
```go
// Example custom sensor for timestamp encoding
type TimestampSensor struct {
    config sensors.SensorConfig
}

func NewTimestampSensor(config sensors.SensorConfig) (sensors.SensorInterface, error) {
    sensor := &TimestampSensor{config: config}
    return sensor, sensor.Validate()
}

func (t *TimestampSensor) Encode(input interface{}) (sensors.SDR, error) {
    timestamp, ok := input.(time.Time)
    if !ok {
        return nil, sensors.NewSensorError("timestamp", "encode", 
            fmt.Sprintf("%v", input), "input must be time.Time")
    }
    
    // Extract time components for encoding
    hour := float64(timestamp.Hour()) / 24.0    // 0-1 range
    minute := float64(timestamp.Minute()) / 60.0 // 0-1 range
    
    // Use composition with numeric sensors for each component
    hourSensor := createComponentSensor(t.config, "hour")
    minuteSensor := createComponentSensor(t.config, "minute")
    
    hourSDR, _ := hourSensor.Encode(hour)
    minuteSDR, _ := minuteSensor.Encode(minute)
    
    // Combine SDRs (simplified - real implementation would be more sophisticated)
    return combinedSDR(hourSDR, minuteSDR), nil
}

func (t *TimestampSensor) Configure(config sensors.SensorConfig) error {
    t.config = config
    return t.Validate()
}

func (t *TimestampSensor) Validate() error {
    // Validate configuration parameters
    if t.config.SDRWidth() < 100 {
        return sensors.NewConfigurationError("sdr_width", t.config.SDRWidth(), 
            "minimum 100", "insufficient bits for time encoding")
    }
    return nil
}

func (t *TimestampSensor) Metadata() sensors.SensorMetadata {
    return sensors.SensorMetadata{
        Type:            "timestamp",
        InputTypes:      []string{"time.Time"},
        OutputSDRWidth:  t.config.SDRWidth(),
        TargetSparsity:  t.config.TargetSparsity(),
        Description:     "Encodes timestamps into time-of-day representations",
        BiologicalBasis: "Circadian rhythm encoding in suprachiasmatic nucleus",
        Parameters: map[string]string{
            "sdr_width":       "Total bits in output SDR",
            "target_sparsity": "Desired active bit percentage",
        },
    }
}

func (t *TimestampSensor) Clone() sensors.SensorInterface {
    return &TimestampSensor{config: t.config.Clone()}
}

// Register and use custom sensor
func useCustomSensor(registry *sensors.Registry) {
    // Register custom sensor
    registry.Register("timestamp", NewTimestampSensor)
    
    // Configure and use
    config := sensors.NewConfig()
    config.SetParam("sdr_width", 1024)
    config.SetParam("target_sparsity", 0.02)
    
    sensor, _ := registry.Create("timestamp", config)
    sdr, _ := sensor.Encode(time.Now())
    
    fmt.Printf("Timestamp encoding -> SDR width: %d, sparsity: %.3f\n", 
        sdr.Width(), sdr.Sparsity())
}
```

## Testing and Validation

### Built-in Validation
```go
func validateSensorOutput(sdr sensors.SDR) {
    // Automatic validation
    if err := sdr.Validate(); err != nil {
        log.Printf("SDR validation failed: %v", err)
        return
    }
    
    // Check HTM properties
    width := sdr.Width()
    activeBits := len(sdr.ActiveBits())
    sparsity := sdr.Sparsity()
    
    fmt.Printf("SDR Properties:\n")
    fmt.Printf("  Width: %d bits\n", width)
    fmt.Printf("  Active bits: %d\n", activeBits)
    fmt.Printf("  Sparsity: %.3f (%.1f%%)\n", sparsity, sparsity*100)
    
    // Validate HTM constraints
    if sparsity < 0.01 || sparsity > 0.10 {
        fmt.Printf("⚠️  Sparsity outside recommended range (1-10%%)\n")
    }
    
    if width < 1000 || width > 10000 {
        fmt.Printf("⚠️  Width outside typical range (1000-10000 bits)\n")
    }
    
    fmt.Println("✓ SDR validation completed")
}
```

### Performance Testing
```go
func benchmarkSensorPerformance(sensor sensors.SensorInterface) {
    const iterations = 1000
    const targetLatency = time.Millisecond // 1ms target
    
    // Prepare test data
    testData := generateTestData(iterations)
    
    // Benchmark encoding performance
    start := time.Now()
    for _, data := range testData {
        _, err := sensor.Encode(data)
        if err != nil {
            log.Printf("Encoding error: %v", err)
        }
    }
    elapsed := time.Since(start)
    
    avgLatency := elapsed / iterations
    
    fmt.Printf("Performance Results:\n")
    fmt.Printf("  Iterations: %d\n", iterations)
    fmt.Printf("  Total time: %v\n", elapsed)
    fmt.Printf("  Average latency: %v\n", avgLatency)
    fmt.Printf("  Throughput: %.0f ops/sec\n", float64(iterations)/elapsed.Seconds())
    
    if avgLatency <= targetLatency {
        fmt.Println("✓ Performance target met")
    } else {
        fmt.Printf("⚠️  Performance below target (%.2fx slower)\n", 
            float64(avgLatency)/float64(targetLatency))
    }
}
```

## Complete Example

Run the complete quickstart example:

```bash
go run quickstart.go
```

Expected output:
```
Registered sensors: [numeric categorical text spatial]
Value 25.500000 -> SDR (width: 2048, sparsity: 0.020, active bits: 41)
Value 50.000000 -> SDR (width: 2048, sparsity: 0.020, active bits: 41)
Value 75.300000 -> SDR (width: 2048, sparsity: 0.020, active bits: 41)
✓ Consistency test passed: identical inputs produce identical SDRs
Category 'red' -> SDR (active bits: [42 87 123 156 203])
Similarity between 'red' and 'blue': 0.023
✓ Semantic separation test passed
Text 'hello world' -> SDR (sparsity: 0.030)
Similarity 'hello world' vs 'hello universe': 0.654
Similarity 'hello world' vs 'goodbye world': 0.234
Consistency 'hello world' duplicate: 1.000
✓ All tests passed
```

## Next Steps

1. **Read the Documentation**: Explore detailed API documentation for advanced configuration
2. **Experiment**: Try different sparsity levels and SDR widths for your use case
3. **Custom Sensors**: Implement domain-specific sensors using the SensorInterface
4. **Integration**: Connect SDR outputs to HTM spatial pooler and temporal memory algorithms
5. **Performance Tuning**: Optimize encoding parameters for your specific data characteristics

## Common Issues and Solutions

**Issue**: SDR sparsity doesn't match target
**Solution**: Adjust resolution or bucket size parameters for your encoder type

**Issue**: Poor semantic similarity preservation
**Solution**: Increase SDR width or adjust encoder-specific parameters

**Issue**: Encoding performance too slow
**Solution**: Reduce SDR width, batch process inputs, or optimize custom sensor logic

**Issue**: Inconsistent encoding results
**Solution**: Ensure all configuration parameters are deterministic and avoid floating-point precision issues