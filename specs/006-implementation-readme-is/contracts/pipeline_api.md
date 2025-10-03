# API Contracts: Complete Pipeline Endpoints

## POST /api/v1/pipeline/process

Process data through the complete HTM pipeline from sensor input to motor output.

### Request Schema

```json
{
  "type": "object",
  "required": ["sensor_data", "pipeline_id"],
  "properties": {
    "sensor_data": {
      "type": "object",
      "required": ["type", "data", "metadata"],
      "properties": {
        "type": {"type": "string", "enum": ["temperature", "text", "image", "audio"]},
        "data": {
          "description": "Type-specific sensor data",
          "oneOf": [
            {
              "title": "Temperature Data",
              "type": "object",
              "properties": {
                "value": {"type": "number"},
                "unit": {"type": "string", "enum": ["celsius", "fahrenheit", "kelvin"]},
                "precision": {"type": "number"}
              }
            },
            {
              "title": "Text Data",
              "type": "object", 
              "properties": {
                "text": {"type": "string", "maxLength": 1000},
                "encoding": {"type": "string", "default": "utf-8"},
                "language": {"type": "string"}
              }
            },
            {
              "title": "Image Data",
              "type": "object",
              "properties": {
                "base64": {"type": "string"},
                "format": {"type": "string", "enum": ["jpeg", "png", "bmp"]},
                "width": {"type": "integer", "minimum": 1},
                "height": {"type": "integer", "minimum": 1}
              }
            },
            {
              "title": "Audio Data",
              "type": "object",
              "properties": {
                "base64": {"type": "string"},
                "format": {"type": "string", "enum": ["wav", "mp3", "raw"]},
                "sample_rate": {"type": "integer", "minimum": 8000},
                "channels": {"type": "integer", "enum": [1, 2]}
              }
            }
          ]
        },
        "metadata": {
          "type": "object",
          "properties": {
            "sensor_id": {"type": "string"},
            "timestamp": {"type": "string", "format": "date-time"},
            "sequence_id": {"type": "string"},
            "step_number": {"type": "integer", "minimum": 0},
            "quality_score": {"type": "number", "minimum": 0, "maximum": 1}
          }
        }
      }
    },
    "pipeline_id": {
      "type": "string",
      "pattern": "^[a-zA-Z0-9-_]{1,64}$",
      "description": "Pipeline instance identifier"
    },
    "processing_options": {
      "type": "object",
      "properties": {
        "enable_learning": {"type": "boolean", "default": true},
        "prediction_steps": {"type": "integer", "minimum": 1, "maximum": 5, "default": 1},
        "motor_output_threshold": {"type": "number", "minimum": 0, "maximum": 1, "default": 0.7},
        "max_motor_commands": {"type": "integer", "minimum": 1, "maximum": 5, "default": 3}
      }
    }
  }
}
```

### Response Schema

```json
{
  "type": "object",
  "required": ["pipeline_result", "processing_summary", "status"],
  "properties": {
    "pipeline_result": {
      "type": "object",
      "required": ["encoded_sdr", "spatial_sdr", "temporal_result", "motor_commands"],
      "properties": {
        "encoded_sdr": {
          "type": "object",
          "properties": {
            "width": {"type": "integer"},
            "active_bits": {"type": "array", "items": {"type": "integer"}},
            "encoder_type": {"type": "string"}
          }
        },
        "spatial_sdr": {
          "type": "object",
          "properties": {
            "width": {"type": "integer"},
            "active_bits": {"type": "array", "items": {"type": "integer"}},
            "sparsity": {"type": "number"},
            "overlap_score": {"type": "number"}
          }
        },
        "temporal_result": {
          "type": "object",
          "properties": {
            "active_cells": {"type": "array"},
            "predictive_cells": {"type": "array"},
            "predictions": {"type": "array"},
            "sequence_stability": {"type": "number"}
          }
        },
        "motor_commands": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "id": {"type": "string"},
              "type": {"type": "string"},
              "parameters": {"type": "object"},
              "confidence": {"type": "number"},
              "priority": {"type": "integer"}
            }
          }
        }
      }
    },
    "processing_summary": {
      "type": "object",
      "required": ["total_time", "stage_times", "success"],
      "properties": {
        "total_time": {"type": "number", "description": "Total processing time in milliseconds"},
        "stage_times": {
          "type": "object",
          "properties": {
            "encoding_time": {"type": "number"},
            "spatial_pooling_time": {"type": "number"},
            "temporal_memory_time": {"type": "number"},
            "motor_output_time": {"type": "number"}
          }
        },
        "success": {"type": "boolean"},
        "warning_count": {"type": "integer"},
        "error_count": {"type": "integer"}
      }
    },
    "status": {"type": "string", "enum": ["success", "partial_success", "error"]},
    "htm_compliance": {
      "type": "object",
      "properties": {
        "encoder_sparsity": {"type": "number"},
        "spatial_sparsity": {"type": "number"},
        "temporal_sparsity": {"type": "number"},
        "biological_constraints_met": {"type": "boolean"},
        "performance_within_limits": {"type": "boolean"}
      }
    },
    "session_info": {
      "type": "object",
      "properties": {
        "session_id": {"type": "string"},
        "pipeline_id": {"type": "string"},
        "sequence_position": {"type": "integer"},
        "processing_timestamp": {"type": "string", "format": "date-time"}
      }
    },
    "errors": {
      "type": "array",
      "items": {"type": "string"}
    },
    "warnings": {
      "type": "array",
      "items": {"type": "string"}
    }
  }
}
```

### Error Responses

**400 Bad Request**:
```json
{
  "error": "invalid_sensor_data",
  "message": "Sensor data format not supported",
  "details": {
    "received_type": "video",
    "supported_types": ["temperature", "text", "image", "audio"]
  }
}
```

**413 Payload Too Large**:
```json
{
  "error": "data_size_exceeded",
  "message": "Sensor data exceeds maximum size limit",
  "details": {
    "max_size_mb": 10,
    "received_size_mb": 15.7
  }
}
```

**504 Gateway Timeout**:
```json
{
  "error": "processing_timeout",
  "message": "Pipeline processing exceeded time limit",
  "details": {
    "timeout_ms": 100,
    "processing_time_ms": 127,
    "failed_stage": "temporal_memory"
  }
}
```

---

## GET /api/v1/pipeline/status

Get pipeline instance status and health metrics.

### Request Parameters

- `pipeline_id` (query, required): Pipeline instance identifier

### Response Schema

```json
{
  "type": "object",
  "required": ["status", "health", "components"],
  "properties": {
    "status": {"type": "string", "enum": ["healthy", "degraded", "error", "maintenance"]},
    "health": {
      "type": "object",
      "properties": {
        "overall_score": {"type": "number", "minimum": 0, "maximum": 1},
        "uptime_seconds": {"type": "integer"},
        "total_sessions_processed": {"type": "integer"},
        "avg_processing_time": {"type": "number"},
        "error_rate": {"type": "number", "minimum": 0, "maximum": 1},
        "throughput_per_second": {"type": "number"}
      }
    },
    "components": {
      "type": "object",
      "properties": {
        "encoder": {
          "type": "object",
          "properties": {
            "status": {"type": "string"},
            "supported_types": {"type": "array", "items": {"type": "string"}},
            "performance": {"type": "object"}
          }
        },
        "spatial_pooler": {
          "type": "object",
          "properties": {
            "status": {"type": "string"},
            "instance_id": {"type": "string"},
            "htm_compliance": {"type": "boolean"},
            "performance": {"type": "object"}
          }
        },
        "temporal_memory": {
          "type": "object",
          "properties": {
            "status": {"type": "string"},
            "instance_id": {"type": "string"},
            "cell_count": {"type": "integer"},
            "sequence_count": {"type": "integer"},
            "performance": {"type": "object"}
          }
        },
        "motor_output": {
          "type": "object",
          "properties": {
            "status": {"type": "string"},
            "instance_id": {"type": "string"},
            "active_commands": {"type": "integer"},
            "action_types": {"type": "array", "items": {"type": "string"}},
            "performance": {"type": "object"}
          }
        }
      }
    },
    "configuration": {
      "type": "object",
      "properties": {
        "pipeline_id": {"type": "string"},
        "created_at": {"type": "string", "format": "date-time"},
        "learning_enabled": {"type": "boolean"},
        "max_concurrent_sessions": {"type": "integer"},
        "performance_targets": {"type": "object"}
      }
    },
    "recent_activity": {
      "type": "array",
      "maxItems": 10,
      "items": {
        "type": "object",
        "properties": {
          "session_id": {"type": "string"},
          "sensor_type": {"type": "string"},
          "processing_time": {"type": "number"},
          "success": {"type": "boolean"},
          "timestamp": {"type": "string", "format": "date-time"}
        }
      }
    }
  }
}
```

---

## PUT /api/v1/pipeline/config

Update pipeline configuration.

### Request Schema

```json
{
  "type": "object",
  "required": ["pipeline_id"],
  "properties": {
    "pipeline_id": {"type": "string"},
    "config": {
      "type": "object",
      "properties": {
        "learning_enabled": {"type": "boolean"},
        "max_concurrent_sessions": {"type": "integer", "minimum": 1, "maximum": 50},
        "processing_timeout_ms": {"type": "integer", "minimum": 50, "maximum": 1000},
        "performance_targets": {
          "type": "object",
          "properties": {
            "max_processing_time_ms": {"type": "number"},
            "min_htm_compliance": {"type": "number", "minimum": 0, "maximum": 1},
            "target_throughput": {"type": "number"}
          }
        },
        "component_configs": {
          "type": "object",
          "properties": {
            "spatial_pooler": {"type": "object"},
            "temporal_memory": {"type": "object"},
            "motor_output": {"type": "object"}
          }
        }
      }
    }
  }
}
```

### Response Schema

```json
{
  "type": "object",
  "required": ["status", "updated_config"],
  "properties": {
    "status": {"type": "string", "enum": ["success", "partial_success", "error"]},
    "updated_config": {
      "type": "object",
      "description": "Current configuration after update"
    },
    "validation_results": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "component": {"type": "string"},
          "parameter": {"type": "string"},
          "status": {"type": "string", "enum": ["valid", "warning", "error"]},
          "message": {"type": "string"}
        }
      }
    },
    "restart_required": {"type": "boolean"},
    "affected_components": {"type": "array", "items": {"type": "string"}}
  }
}
```

---

## POST /api/v1/pipeline/reset

Reset pipeline state and metrics.

### Request Schema

```json
{
  "type": "object",
  "required": ["pipeline_id"],
  "properties": {
    "pipeline_id": {"type": "string"},
    "reset_options": {
      "type": "object",
      "properties": {
        "reset_metrics": {"type": "boolean", "default": true},
        "reset_learning": {"type": "boolean", "default": false},
        "reset_temporal_memory": {"type": "boolean", "default": false},
        "reset_motor_mappings": {"type": "boolean", "default": false},
        "preserve_config": {"type": "boolean", "default": true}
      }
    }
  }
}
```

### Response Schema

```json
{
  "type": "object",
  "required": ["status", "reset_summary"],
  "properties": {
    "status": {"type": "string", "enum": ["success", "partial_success", "error"]},
    "reset_summary": {
      "type": "object",
      "properties": {
        "components_reset": {"type": "array", "items": {"type": "string"}},
        "metrics_cleared": {"type": "boolean"},
        "learning_state_reset": {"type": "boolean"},
        "temporal_memory_cleared": {"type": "boolean"},
        "motor_mappings_reset": {"type": "boolean"},
        "reset_timestamp": {"type": "string", "format": "date-time"}
      }
    },
    "new_state": {
      "type": "object",
      "properties": {
        "pipeline_status": {"type": "string"},
        "component_states": {"type": "object"},
        "ready_for_processing": {"type": "boolean"}
      }
    }
  }
}
```