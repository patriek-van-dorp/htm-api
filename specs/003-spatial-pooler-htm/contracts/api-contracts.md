# API Contracts: Spatial Pooler (HTM Theory) Component

**Date**: October 1, 2025  
**Feature**: Spatial pooler HTM API endpoints  
**Integration**: Extends existing HTM API with spatial pooler functionality

## REST API Endpoints

### POST /api/v1/spatial-pooler/process
**Purpose**: Process sensor encoder output through spatial pooler to create proper HTM-compliant SDRs

**Request Body**:
```json
{
  "input": {
    "encoder_output": {
      "width": 2048,
      "active_bits": [10, 25, 67, 89, 134, 256, 445, 678, 789, 901],
      "sparsity": 0.0048828125
    },
    "input_id": "uuid-string",
    "learning_enabled": true,
    "metadata": {
      "sensor_type": "categorical",
      "encoder_version": "1.0.0"
    }
  },
  "config": {
    "mode": "deterministic",
    "learning_enabled": true
  },
  "request_id": "req-uuid-string",
  "client_id": "client-123"
}
```

**Response Body (Success - 200)**:
```json
{
  "result": {
    "sdr": {
      "width": 2048,
      "active_bits": [12, 34, 78, 123, 234, 345, 567, 789, 890, 1001],
      "sparsity": 0.0244140625
    },
    "input_id": "uuid-string",
    "processing_time_ms": 8.5,
    "active_columns": [12, 34, 78, 123, 234, 345, 567, 789, 890, 1001],
    "avg_overlap": 0.75,
    "sparsity_level": 0.0244140625,
    "learning_occurred": true,
    "boosting_applied": false
  },
  "request_id": "req-uuid-string",
  "processing_time_ms": 9.2,
  "status": "success"
}
```
```

**Response Body (Error - 400)**:
```json
{
  "error": {
    "type": "invalid_input",
    "message": "Encoder output exceeds expected input dimensions",
    "input_id": "uuid-string",
    "details": {
      "expected_width": 2048,
      "actual_width": 4096,
      "config_field": "input_dimensions"
    }
  },
  "request_id": "req-uuid-string",
  "status": "error"
}
```

**Response Body (Error - 500)**:
```json
{
  "error": {
    "type": "processing_error", 
    "message": "Spatial pooler processing failed",
    "input_id": "uuid-string",
    "details": {
      "processing_stage": "overlap_calculation",
      "internal_error": "matrix dimension mismatch"
    }
  },
  "request_id": "req-uuid-string", 
  "status": "error"
}
```

### GET /api/v1/spatial-pooler/config
**Purpose**: Retrieve current spatial pooler configuration

**Response Body (Success - 200)**:
```json
{
  "config": {
    "input_dimensions": [2048],
    "column_dimensions": [2048],
    "potential_radius": 0.1,
    "potential_pct": 0.5,
    "global_inhibition": true,
    "num_active_columns": 50,
    "stimulus_threshold": 1,
    "syn_perm_inc": 0.05,
    "syn_perm_dec": 0.008,
    "syn_perm_connected": 0.1,
    "syn_perm_trim_threshold": 0.01,
    "duty_cycle_period": 1000,
    "boost_strength": 2.0,
    "learning_enabled": true,
    "mode": "deterministic"
  },
  "status": "success"
}
```

### PUT /api/v1/spatial-pooler/config
**Purpose**: Update spatial pooler configuration

**Request Body**:
```json
{
  "config": {
    "global_inhibition": true,
    "num_active_columns": 100,
    "learning_enabled": false,
    "boost_strength": 1.5,
    "mode": "randomized"
  }
}
```

**Response Body (Success - 200)**:
```json
{
  "config": {
    "input_dimensions": [2048],
    "column_dimensions": [2048], 
    "potential_radius": 0.1,
    "potential_pct": 0.5,
    "global_inhibition": true,
    "num_active_columns": 100,
    "stimulus_threshold": 1,
    "syn_perm_inc": 0.05,
    "syn_perm_dec": 0.008,
    "syn_perm_connected": 0.1,
    "syn_perm_trim_threshold": 0.01,
    "duty_cycle_period": 1000,
    "boost_strength": 1.5,
    "learning_enabled": false,
    "mode": "randomized"
  },
  "message": "Configuration updated successfully",
  "status": "success"
}
```

### GET /api/v1/spatial-pooler/metrics
**Purpose**: Retrieve spatial pooler performance and behavioral metrics

**Response Body (Success - 200)**:
```json
{
  "metrics": {
    "total_processed": 15423,
    "average_processing_time_ms": 7.8,
    "learning_iterations": 12456,
    "average_sparsity": 0.0245,
    "boosting_events": 234,
    "error_counts": {
      "invalid_input": 12,
      "configuration_error": 0,
      "processing_error": 3,
      "performance_error": 1
    },
    "column_usage_stats": {
      "min_usage": 0.0123,
      "max_usage": 0.0456,
      "std_deviation": 0.0067
    },
    "overlap_score_stats": {
      "min_overlap": 0.0,
      "max_overlap": 1.0,
      "average_overlap": 0.687
    }
  },
  "status": "success"
}
```

### POST /api/v1/spatial-pooler/reset
**Purpose**: Reset spatial pooler state and learning

**Request Body**:
```json
{
  "reset_learning": true,
  "reset_duty_cycles": true,
  "reset_boost_factors": true,
  "preserve_config": true
}
```

**Response Body (Success - 200)**:
```json
{
  "message": "Spatial pooler reset successfully",
  "reset_components": [
    "learning_state",
    "duty_cycles", 
    "boost_factors"
  ],
  "preserved_components": [
    "configuration",
    "permanence_matrix"
  ],
  "status": "success"
}
```

## Request/Response Schemas

### Common Types
```json
{
  "SDR": {
    "type": "object",
    "properties": {
      "width": {"type": "integer", "minimum": 1},
      "active_bits": {
        "type": "array",
        "items": {"type": "integer", "minimum": 0}
      },
      "sparsity": {"type": "number", "minimum": 0, "maximum": 1}
    },
    "required": ["width", "active_bits", "sparsity"]
  },
  
  "SpatialPoolerMode": {
    "type": "string",
    "enum": ["deterministic", "randomized"]
  },
  
  "PoolingErrorType": {
    "type": "string", 
    "enum": ["invalid_input", "configuration_error", "processing_error", "performance_error", "learning_error"]
  }
}
```

### Request Validation Rules
- `encoder_output.width` must match configured input dimensions
- `encoder_output.active_bits` must be sorted array of unique integers
- `encoder_output.sparsity` must match calculated sparsity from active_bits
- `input_id` must be non-empty string (UUID format recommended)
- `request_id` must be non-empty string and unique per request
- `learning_enabled` must be boolean (optional, defaults to config value)
- `config.mode` must be valid SpatialPoolerMode enum value

### Response Validation Rules
- `normalized_sdr.sparsity` must be between 0.02 and 0.05 (2-5%)
- `processing_time_ms` should be under 10ms for performance compliance
- `active_columns` must be sorted array with length matching sparsity target
- `sparsity_level` must match `normalized_sdr.sparsity`
- Error responses must include appropriate HTTP status codes (400, 500)
- All timestamps must be in ISO 8601 format where applicable

## Integration with Existing API

### Compatibility with Existing Endpoints
- Spatial pooler endpoints extend current HTM API structure
- Existing `/api/v1/process` endpoint remains unchanged
- New endpoints follow existing authentication/authorization patterns
- Request/response format matches existing HTM API conventions

### Common Headers
```
Content-Type: application/json
Accept: application/json
Authorization: Bearer {token}  // If authentication enabled
X-Request-ID: {uuid}          // Request tracking
X-Client-Version: {version}   // Client version for compatibility
```

### Error Response Format (Consistent with Existing API)
```json
{
  "error": {
    "type": "error_category",
    "message": "Human readable message", 
    "details": {
      "field": "specific_error_context"
    }
  },
  "request_id": "request-uuid",
  "timestamp": "2025-10-01T10:30:00Z",
  "status": "error"
}
```

## Performance Contracts

### Processing Time Guarantees
- Spatial pooler processing: <10ms per request
- Configuration updates: <100ms  
- Metrics retrieval: <50ms
- Reset operations: <200ms

### Throughput Requirements
- Support 1,000-5,000 requests/second under normal load
- Graceful degradation under high load with appropriate error responses
- Memory usage scales linearly with configured column dimensions

### Reliability Guarantees  
- Deterministic mode produces identical outputs for identical inputs
- Learning state persists across configuration updates (unless explicitly reset)
- Error recovery maintains spatial pooler state consistency
- No memory leaks during extended operation