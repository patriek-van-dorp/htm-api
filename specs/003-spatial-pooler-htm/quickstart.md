# Quickstart: Spatial Pooler (HTM Theory) Component

**Date**: October 1, 2025  
**Feature**: Spatial pooler HTM implementation  
**Audience**: Developers integrating spatial pooler functionality

## Overview

This quickstart guide demonstrates how to integrate and use the spatial pooler component in the HTM API. The spatial pooler is implemented as the first layer of the cortical column, processing outputs from existing sensor encoders to produce proper sparse distributed representations (SDRs) with 2-5% sparsity levels while preserving semantic relationships.

**Architecture Flow**: Raw Input → Sensor Encoding → **Spatial Pooler** → True SDR → Future Temporal Memory

**Note**: While proper HTM theory suggests encoders should output raw bit arrays and spatial poolers should create SDRs, the current implementation works with existing sensor SDR outputs and normalizes them into proper HTM-compliant SDRs.

## Prerequisites

- HTM API is running and accessible
- Existing sensor encoders producing SDR output (current implementation - will migrate to bit arrays in future)
- Basic understanding of sparse distributed representations (SDRs)
- HTTP client for API requests (curl, Postman, or similar)
- **Note**: Spatial pooler normalizes existing sensor SDR outputs into proper HTM-compliant SDRs

## Quick Start Steps

### 1. Verify API Availability

First, confirm the spatial pooler endpoints are available:

```bash
curl -X GET http://localhost:8080/api/v1/spatial-pooler/config \
  -H "Content-Type: application/json"
```

Expected response:
```json
{
  "config": {
    "input_dimensions": [2048],
    "column_dimensions": [2048],
    "global_inhibition": true,
    "num_active_columns": 50,
    "learning_enabled": true,
    "mode": "deterministic"
  },
  "status": "success"
}
```

### 2. Process Your First Encoder Output

Send a typical encoder output through the spatial pooler:

```bash
curl -X POST http://localhost:8080/api/v1/spatial-pooler/process \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "encoder_output": {
        "width": 2048,
        "active_bits": [10, 25, 67, 89, 134, 256, 445, 678, 789, 901],
        "sparsity": 0.0048828125
      },
      "input_id": "quickstart-test-001",
      "learning_enabled": true,
      "metadata": {
        "sensor_type": "categorical",
        "test_data": true
      }
    },
    "config": {
      "mode": "deterministic",
      "learning_enabled": true
    },
    "request_id": "quickstart-req-001",
    "client_id": "quickstart-client"
  }'
```

Expected response:
```json
{
  "result": {
    "normalized_sdr": {
      "width": 2048,
      "active_bits": [12, 34, 78, 123, 234, 345, 567, 789, 890, 1001],
      "sparsity": 0.0244140625
    },
    "input_id": "quickstart-test-001",
    "processing_time_ms": 8.5,
    "sparsity_level": 0.0244140625,
    "learning_occurred": true,
    "boosting_applied": false
  },
  "request_id": "quickstart-req-001",
  "processing_time_ms": 9.2,
  "status": "success"
}
```

### 3. Verify Deterministic Behavior

Send the same input again to verify deterministic processing:

```bash
curl -X POST http://localhost:8080/api/v1/spatial-pooler/process \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "encoder_output": {
        "width": 2048,
        "active_bits": [10, 25, 67, 89, 134, 256, 445, 678, 789, 901],
        "sparsity": 0.0048828125
      },
      "input_id": "quickstart-test-002",
      "learning_enabled": false
    },
    "config": {
      "mode": "deterministic",
      "learning_enabled": false
    },
    "request_id": "quickstart-req-002"
  }'
```

The `normalized_sdr.active_bits` should be identical to the previous response when learning is disabled.

### 4. Test Semantic Similarity

Process a similar input to verify semantic continuity:

```bash
curl -X POST http://localhost:8080/api/v1/spatial-pooler/process \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "encoder_output": {
        "width": 2048,
        "active_bits": [10, 25, 67, 89, 134, 256, 445, 678, 799, 911],
        "sparsity": 0.0048828125
      },
      "input_id": "quickstart-test-003",
      "learning_enabled": false
    },
    "config": {
      "mode": "deterministic",
      "learning_enabled": false
    },
    "request_id": "quickstart-req-003"
  }'
```

The output should have significant overlap with the previous result, demonstrating semantic similarity preservation.

### 5. Check Performance Metrics

Monitor spatial pooler performance:

```bash
curl -X GET http://localhost:8080/api/v1/spatial-pooler/metrics \
  -H "Content-Type: application/json"
```

Expected response:
```json
{
  "metrics": {
    "total_processed": 3,
    "average_processing_time_ms": 8.2,
    "learning_iterations": 1,
    "average_sparsity": 0.0244,
    "boosting_events": 0,
    "error_counts": {
      "invalid_input": 0,
      "processing_error": 0
    }
  },
  "status": "success"
}
```

## Common Use Cases

### Use Case 1: Categorical Data Processing

For categorical sensor data (e.g., categories, discrete values):

```bash
# Example: Processing a categorical encoder output
curl -X POST http://localhost:8080/api/v1/spatial-pooler/process \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "encoder_output": {
        "width": 2048,
        "active_bits": [100, 101, 102, 103, 104],
        "sparsity": 0.002441
      },
      "input_id": "cat-item-red",
      "metadata": {
        "category": "color",
        "value": "red"
      }
    },
    "request_id": "cat-req-001"
  }'
```

### Use Case 2: Numeric Data Processing

For numeric sensor data (e.g., scalar encoders):

```bash
# Example: Processing a numeric encoder output
curl -X POST http://localhost:8080/api/v1/spatial-pooler/process \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "encoder_output": {
        "width": 2048,
        "active_bits": [500, 501, 502, 503, 504, 505, 506, 507, 508, 509],
        "sparsity": 0.004883
      },
      "input_id": "num-temp-72.5",
      "metadata": {
        "sensor": "temperature",
        "value": 72.5,
        "units": "fahrenheit"
      }
    },
    "request_id": "num-req-001"
  }'
```

### Use Case 3: Configuration Customization

Adjust spatial pooler parameters for specific requirements:

```bash
# Update configuration for higher sparsity
curl -X PUT http://localhost:8080/api/v1/spatial-pooler/config \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "num_active_columns": 100,
      "learning_enabled": true,
      "boost_strength": 2.0,
      "mode": "randomized"
    }
  }'
```

## Integration Examples

### Go Client Integration

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type SpatialPoolerRequest struct {
    Input struct {
        EncoderOutput struct {
            Width      int     `json:"width"`
            ActiveBits []int   `json:"active_bits"`
            Sparsity   float64 `json:"sparsity"`
        } `json:"encoder_output"`
        InputID  string `json:"input_id"`
    } `json:"input"`
    RequestID string `json:"request_id"`
}

func processSpatialPooler(encoderOutput []int, width int, inputID string) error {
    request := SpatialPoolerRequest{}
    request.Input.EncoderOutput.Width = width
    request.Input.EncoderOutput.ActiveBits = encoderOutput
    request.Input.EncoderOutput.Sparsity = float64(len(encoderOutput)) / float64(width)
    request.Input.InputID = inputID
    request.RequestID = "go-client-" + inputID

    jsonData, err := json.Marshal(request)
    if err != nil {
        return err
    }

    resp, err := http.Post(
        "http://localhost:8080/api/v1/spatial-pooler/process",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Handle response...
    fmt.Printf("Status: %s\n", resp.Status)
    return nil
}
```

### Python Client Integration

```python
import requests
import json

def process_spatial_pooler(encoder_output, width, input_id):
    """Process encoder output through spatial pooler"""
    
    sparsity = len(encoder_output) / width
    
    request_data = {
        "input": {
            "encoder_output": {
                "width": width,
                "active_bits": encoder_output,
                "sparsity": sparsity
            },
            "input_id": input_id
        },
        "request_id": f"python-client-{input_id}"
    }
    
    response = requests.post(
        "http://localhost:8080/api/v1/spatial-pooler/process",
        headers={"Content-Type": "application/json"},
        json=request_data
    )
    
    if response.status_code == 200:
        result = response.json()
        return result["result"]["normalized_sdr"]
    else:
        raise Exception(f"API error: {response.status_code}")

# Example usage
encoder_bits = [10, 25, 67, 89, 134]
normalized_sdr = process_spatial_pooler(encoder_bits, 2048, "test-001")
print(f"Normalized SDR: {normalized_sdr}")
```

## Testing and Validation

### Functional Test Script

```bash
#!/bin/bash
# Spatial pooler functional test script

echo "Testing spatial pooler functionality..."

# Test 1: Basic processing
echo "Test 1: Basic processing"
response=$(curl -s -X POST http://localhost:8080/api/v1/spatial-pooler/process \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "encoder_output": {
        "width": 2048,
        "active_bits": [1, 2, 3, 4, 5],
        "sparsity": 0.00244
      },
      "input_id": "test-basic"
    },
    "request_id": "test-basic-req"
  }')

status=$(echo $response | jq -r '.status')
if [ "$status" = "success" ]; then
    echo "✓ Basic processing test passed"
else
    echo "✗ Basic processing test failed"
    echo $response
fi

# Test 2: Sparsity validation
sparsity=$(echo $response | jq -r '.result.sparsity_level')
if (( $(echo "$sparsity >= 0.02" | bc -l) )) && (( $(echo "$sparsity <= 0.05" | bc -l) )); then
    echo "✓ Sparsity level test passed ($sparsity)"
else
    echo "✗ Sparsity level test failed ($sparsity)"
fi

# Test 3: Performance validation
processing_time=$(echo $response | jq -r '.result.processing_time_ms')
if (( $(echo "$processing_time <= 10" | bc -l) )); then
    echo "✓ Performance test passed (${processing_time}ms)"
else
    echo "✗ Performance test failed (${processing_time}ms)"
fi

echo "Functional tests completed."
```

## Troubleshooting

### Common Issues

1. **High Processing Time**
   - Check system load and available memory
   - Consider reducing column dimensions or input size
   - Monitor metrics for performance trends

2. **Sparsity Out of Range**
   - Verify `num_active_columns` configuration
   - Check input encoder sparsity levels
   - Adjust `local_area_density` for local inhibition

3. **Memory Issues**
   - Monitor column dimensions × input dimensions memory usage
   - Consider batch processing for large inputs
   - Reset spatial pooler state periodically if needed

4. **Learning Not Occurring**
   - Verify `learning_enabled` is true in both config and request
   - Check that inputs are sufficiently different
   - Monitor duty cycle periods for adaptation time

### Error Response Handling

```bash
# Example error response
{
  "error": {
    "type": "invalid_input",
    "message": "Encoder output exceeds expected input dimensions",
    "input_id": "oversized-input",
    "details": {
      "expected_width": 2048,
      "actual_width": 4096
    }
  },
  "request_id": "error-req-001",
  "status": "error"
}
```

## Next Steps

1. **Integration with Temporal Memory**: Once implemented, chain spatial pooler output to temporal memory components
2. **Custom Configurations**: Experiment with different parameter settings for your specific use case
3. **Performance Monitoring**: Set up monitoring for processing times and sparsity levels
4. **Batch Processing**: Implement batch processing for high-throughput scenarios
5. **A/B Testing**: Compare results with and without spatial pooler normalization

## Support

- **API Documentation**: See `/contracts/openapi.yaml` for complete API specification
- **Performance Benchmarks**: Monitor `/api/v1/spatial-pooler/metrics` for operational metrics  
- **Configuration Reference**: Check `/api/v1/spatial-pooler/config` for current parameter values
- **Error Debugging**: Enable detailed logging and check error response details

This quickstart provides a foundation for integrating spatial pooler functionality into your HTM-based applications. The spatial pooler serves as the critical first step in creating consistent, semantically meaningful representations for temporal memory processing.