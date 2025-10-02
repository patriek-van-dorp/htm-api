# Integration API Contracts: Complete Spatial Pooler Engine Integration

**Date**: October 2, 2025  
**Feature**: Complete spatial pooler engine integration  
**Purpose**: API contracts for production-ready spatial pooler integration

## Integration Endpoints

### GET /api/v1/health
**Purpose**: Complete health check including spatial pooler engine status

**Response Body (Success - 200)**:
```json
{
  "status": "healthy",
  "timestamp": "2025-10-02T10:30:00Z",
  "version": "1.0.0",
  "components": {
    "spatial_pooler": {
      "status": "healthy",
      "engine_type": "production",
      "implementation": "gonum-optimized",
      "performance": {
        "avg_processing_time_ms": 45.2,
        "requests_processed": 15847,
        "error_rate": 0.001
      }
    },
    "http_server": {
      "status": "healthy",
      "concurrent_requests": 23,
      "max_concurrent": 100
    },
    "memory": {
      "status": "healthy",
      "heap_size_mb": 128.5,
      "gc_pressure": "low"
    }
  },
  "system_metrics": {
    "uptime_seconds": 3600,
    "total_requests": 15847,
    "current_load": "medium"
  }
}
```

**Response Body (Degraded - 503)**:
```json
{
  "status": "degraded",
  "timestamp": "2025-10-02T10:30:00Z",
  "components": {
    "spatial_pooler": {
      "status": "degraded",
      "issues": [
        "processing_time_above_threshold",
        "memory_pressure_detected"
      ],
      "performance": {
        "avg_processing_time_ms": 150.3,
        "error_rate": 0.05
      }
    }
  },
  "recommendations": [
    "reduce_concurrent_load",
    "restart_spatial_pooler_service"
  ]
}
```

### GET /api/v1/spatial-pooler/status
**Purpose**: Detailed spatial pooler engine status and configuration

**Response Body (Success - 200)**:
```json
{
  "engine_status": {
    "implementation": "production",
    "engine_type": "SpatialPooler",
    "version": "1.0.0",
    "initialization_time": "2025-10-02T09:00:00Z",
    "current_state": "ready"
  },
  "configuration": {
    "input_width": 2048,
    "column_count": 2048,
    "sparsity_ratio": 0.02,
    "learning_enabled": true,
    "mode": "deterministic",
    "performance_mode": "optimized"
  },
  "runtime_metrics": {
    "total_requests_processed": 15847,
    "average_processing_time_ms": 45.2,
    "peak_processing_time_ms": 89.7,
    "memory_usage_mb": 64.2,
    "current_learning_iteration": 15847
  },
  "htm_properties": {
    "current_sparsity_level": 0.0195,
    "average_overlap_similarity": 0.65,
    "learning_convergence_rate": 0.92,
    "column_utilization": 0.87
  }
}
```

### POST /api/v1/spatial-pooler/process
**Purpose**: Process data through production spatial pooler engine (extends existing endpoint)

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
      "encoder_version": "1.0.0",
      "request_source": "production_client"
    }
  },
  "config": {
    "mode": "deterministic",
    "learning_enabled": true,
    "performance_mode": "optimized"
  },
  "request_id": "req-uuid-string",
  "client_id": "client-123",
  "integration_context": {
    "test_mode": false,
    "validation_level": "production"
  }
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
    "processing_time_ms": 45.2,
    "active_columns": [12, 34, 78, 123, 234, 345, 567, 789, 890, 1001],
    "avg_overlap": 0.65,
    "sparsity_level": 0.0244140625,
    "learning_applied": true,
    "column_boost_applied": false
  },
  "metadata": {
    "request_id": "req-uuid-string",
    "processing_timestamp": "2025-10-02T10:30:00.123Z",
    "engine_version": "1.0.0",
    "htm_validation": {
      "sparsity_valid": true,
      "overlap_valid": true,
      "learning_valid": true
    }
  },
  "performance": {
    "actual_processing_time_ms": 45.2,
    "target_processing_time_ms": 100.0,
    "memory_used_mb": 2.3,
    "cpu_time_ms": 42.1
  }
}
```

**Response Body (Performance Warning - 200)**:
```json
{
  "result": {
    "sdr": {
      "width": 2048,
      "active_bits": [12, 34, 78, 123, 234, 345, 567, 789, 890, 1001],
      "sparsity": 0.0244140625
    },
    "input_id": "uuid-string",
    "processing_time_ms": 87.5,
    "active_columns": [12, 34, 78, 123, 234, 345, 567, 789, 890, 1001],
    "avg_overlap": 0.65,
    "sparsity_level": 0.0244140625,
    "learning_applied": true,
    "column_boost_applied": false
  },
  "metadata": {
    "request_id": "req-uuid-string",
    "processing_timestamp": "2025-10-02T10:30:00.123Z",
    "engine_version": "1.0.0",
    "htm_validation": {
      "sparsity_valid": true,
      "overlap_valid": true,
      "learning_valid": true
    }
  },
  "performance": {
    "actual_processing_time_ms": 87.5,
    "target_processing_time_ms": 100.0,
    "memory_used_mb": 2.3,
    "cpu_time_ms": 84.2
  },
  "warnings": [
    {
      "type": "performance",
      "message": "Processing time approaching threshold",
      "recommendation": "Monitor system load and consider optimization"
    }
  ]
}
```

**Response Body (Error - 400)**:
```json
{
  "error": {
    "code": "SPATIAL_POOLER_ERROR",
    "message": "Spatial pooler processing failed",
    "details": {
      "validation_errors": [
        "Input width exceeds maximum allowed (10MB limit)",
        "Invalid sparsity level in encoder output"
      ],
      "processing_stage": "input_validation"
    }
  },
  "metadata": {
    "request_id": "req-uuid-string",
    "timestamp": "2025-10-02T10:30:00.123Z",
    "engine_version": "1.0.0"
  }
}
```

**Response Body (Error - 500)**:
```json
{
  "error": {
    "code": "SPATIAL_POOLER_ENGINE_ERROR",
    "message": "Internal spatial pooler engine error",
    "details": {
      "error_type": "matrix_operation_failure",
      "stage": "overlap_computation",
      "recovery_attempted": true
    }
  },
  "metadata": {
    "request_id": "req-uuid-string",
    "timestamp": "2025-10-02T10:30:00.123Z",
    "engine_version": "1.0.0",
    "correlation_id": "error-123-456"
  }
}
```

### PUT /api/v1/spatial-pooler/config
**Purpose**: Update spatial pooler configuration in production environment

**Request Body**:
```json
{
  "config": {
    "learning_enabled": true,
    "sparsity_ratio": 0.025,
    "boost_strength": 1.5,
    "learning_rate": 0.1,
    "performance_mode": "optimized"
  },
  "validation": {
    "validate_only": false,
    "backup_current": true
  },
  "metadata": {
    "updated_by": "admin-user",
    "reason": "Performance optimization",
    "environment": "production"
  }
}
```

**Response Body (Success - 200)**:
```json
{
  "result": {
    "configuration_updated": true,
    "previous_config_backed_up": true,
    "restart_required": false,
    "validation_passed": true
  },
  "updated_config": {
    "learning_enabled": true,
    "sparsity_ratio": 0.025,
    "boost_strength": 1.5,
    "learning_rate": 0.1,
    "performance_mode": "optimized",
    "last_updated": "2025-10-02T10:30:00Z"
  },
  "metadata": {
    "update_timestamp": "2025-10-02T10:30:00Z",
    "config_version": "2.1.0",
    "backup_id": "backup-789"
  }
}
```

## Integration-Specific Error Codes

### Spatial Pooler Engine Errors
- `SPATIAL_POOLER_NOT_INITIALIZED` - Engine not properly initialized
- `SPATIAL_POOLER_ENGINE_ERROR` - Internal engine processing error
- `SPATIAL_POOLER_MEMORY_ERROR` - Memory allocation or management error
- `SPATIAL_POOLER_PERFORMANCE_ERROR` - Processing time exceeded threshold
- `SPATIAL_POOLER_CONCURRENCY_ERROR` - Concurrent access violation

### Input Validation Errors
- `INVALID_INPUT_DATA` - Input data validation failed (invalid dimensions, corrupted data)
- `INVALID_CONFIGURATION` - Configuration parameter validation failed
- `SPARSITY_OUT_OF_RANGE` - Input sparsity outside acceptable range (<0.5% or >10%)
- `INPUT_SIZE_LIMIT_EXCEEDED` - Input data exceeds 10MB limit

### System Resource Errors
- `MEMORY_LIMIT_EXCEEDED` - System memory limit exceeded during processing
- `PERFORMANCE_THRESHOLD_EXCEEDED` - Processing time above 100ms limit
- `CONCURRENCY_LIMIT_EXCEEDED` - Too many concurrent requests (>100)

### Integration Errors
- `INTEGRATION_DEPENDENCY_ERROR` - Missing or invalid service dependency
- `INTEGRATION_CONFIGURATION_ERROR` - Invalid integration configuration
- `INTEGRATION_TEST_MODE_ERROR` - Test mode configuration error
- `INTEGRATION_VALIDATION_ERROR` - HTM property validation failure

## Contract Validation Rules

### Request Validation
- Input encoder output must have valid width and active_bits array
- Sparsity ratio must be between 0.01 and 0.05 (1-5%), with special handling for out-of-range patterns (<0.5% or >10%)
- Request size must not exceed 10MB
- Client ID and request ID must be valid UUIDs

### Response Validation
- Processing time must be reported accurately in milliseconds
- HTM properties (sparsity, overlap) must be within valid ranges
- Error responses must include correlation IDs for debugging
- Performance metrics must be accurate and complete
- Error responses must use standardized error codes for proper client handling

### Integration Validation
- All responses must include engine version and implementation type
- Production responses must include HTM validation status
- Performance warnings must be triggered at 80% of time threshold
- Error responses must provide actionable debugging information

## Test Integration Contracts

### Contract Test Requirements
- All endpoints must be testable against actual implementation
- Mock implementations must be completely replaced with production code
- Test responses must validate HTM mathematical properties
- Performance tests must validate actual timing requirements

### Validation Endpoints for Testing
```json
GET /api/v1/spatial-pooler/validation/htm-properties
{
  "current_state": {
    "sparsity_level": 0.0195,
    "overlap_patterns": "valid",
    "learning_convergence": 0.92,
    "column_utilization": 0.87
  },
  "validation_results": {
    "sparsity_in_range": true,
    "overlap_semantic_preservation": true,
    "learning_adaptation": true,
    "performance_requirements": true
  }
}
```