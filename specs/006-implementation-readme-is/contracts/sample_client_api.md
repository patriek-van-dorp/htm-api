# API Contracts: Sample Client Application

## Sample Client API Overview

The sample client application provides a standalone Go application that demonstrates the complete HTM pipeline by hosting multiple sensor types and integrating with the HTM API. This document defines the client's internal API structure and external HTM API integration patterns.

## Client Application Endpoints

### POST /client/v1/sensors/start

Start sensor data generation and processing.

### Request Schema

```json
{
  "type": "object",
  "required": ["sensor_types", "test_scenario"],
  "properties": {
    "sensor_types": {
      "type": "array",
      "minItems": 1,
      "maxItems": 4,
      "items": {
        "type": "string",
        "enum": ["temperature", "text", "image", "audio"]
      }
    },
    "test_scenario": {
      "type": "object",
      "required": ["name", "duration"],
      "properties": {
        "name": {"type": "string"},
        "description": {"type": "string"},
        "duration": {"type": "string", "pattern": "^[0-9]+[smh]$"},
        "data_generation_rate": {"type": "string", "default": "1s"},
        "expected_outcomes": {"type": "array", "items": {"type": "string"}}
      }
    },
    "htm_api_config": {
      "type": "object",
      "properties": {
        "base_url": {"type": "string", "format": "uri", "default": "http://localhost:8080"},
        "pipeline_id": {"type": "string"},
        "timeout": {"type": "string", "default": "30s"},
        "retry_attempts": {"type": "integer", "minimum": 1, "maximum": 5, "default": 3}
      }
    },
    "performance_monitoring": {
      "type": "object",
      "properties": {
        "collect_metrics": {"type": "boolean", "default": true},
        "metrics_interval": {"type": "string", "default": "5s"},
        "log_level": {"type": "string", "enum": ["debug", "info", "warn", "error"], "default": "info"}
      }
    }
  }
}
```

### Response Schema

```json
{
  "type": "object",
  "required": ["status", "session_id", "active_sensors"],
  "properties": {
    "status": {"type": "string", "enum": ["started", "error"]},
    "session_id": {"type": "string"},
    "active_sensors": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "sensor_id": {"type": "string"},
          "type": {"type": "string"},
          "status": {"type": "string", "enum": ["active", "error"]},
          "data_rate": {"type": "string"},
          "next_reading_at": {"type": "string", "format": "date-time"}
        }
      }
    },
    "htm_connection": {
      "type": "object",
      "properties": {
        "api_url": {"type": "string"},
        "pipeline_id": {"type": "string"},
        "connection_status": {"type": "string", "enum": ["connected", "error"]},
        "last_response_time": {"type": "number"}
      }
    },
    "test_scenario_info": {
      "type": "object",
      "properties": {
        "name": {"type": "string"},
        "estimated_end_time": {"type": "string", "format": "date-time"},
        "total_expected_readings": {"type": "integer"}
      }
    }
  }
}
```

---

### GET /client/v1/session/{session_id}/status

Get current session status and metrics.

### Response Schema

```json
{
  "type": "object",
  "required": ["session_info", "sensor_status", "htm_metrics"],
  "properties": {
    "session_info": {
      "type": "object",
      "properties": {
        "session_id": {"type": "string"},
        "status": {"type": "string", "enum": ["running", "completed", "error", "paused"]},
        "start_time": {"type": "string", "format": "date-time"},
        "elapsed_time": {"type": "string"},
        "estimated_remaining": {"type": "string"}
      }
    },
    "sensor_status": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "sensor_id": {"type": "string"},
          "type": {"type": "string"},
          "status": {"type": "string"},
          "readings_generated": {"type": "integer"},
          "readings_processed": {"type": "integer"},
          "last_reading": {"type": "object"},
          "error_count": {"type": "integer"},
          "success_rate": {"type": "number", "minimum": 0, "maximum": 1}
        }
      }
    },
    "htm_metrics": {
      "type": "object",
      "properties": {
        "total_requests": {"type": "integer"},
        "successful_requests": {"type": "integer"},
        "failed_requests": {"type": "integer"},
        "avg_response_time": {"type": "number"},
        "pipeline_success_rate": {"type": "number"},
        "motor_commands_received": {"type": "integer"},
        "motor_commands_executed": {"type": "integer"}
      }
    },
    "performance_summary": {
      "type": "object",
      "properties": {
        "end_to_end_latency": {"type": "number"},
        "htm_compliance_score": {"type": "number"},
        "biological_constraints_met": {"type": "boolean"},
        "sensor_data_quality": {"type": "number"},
        "prediction_accuracy": {"type": "number"}
      }
    },
    "recent_motor_commands": {
      "type": "array",
      "maxItems": 10,
      "items": {
        "type": "object",
        "properties": {
          "command_id": {"type": "string"},
          "type": {"type": "string"},
          "parameters": {"type": "object"},
          "confidence": {"type": "number"},
          "execution_status": {"type": "string"},
          "timestamp": {"type": "string", "format": "date-time"}
        }
      }
    }
  }
}
```

---

### POST /client/v1/session/{session_id}/stop

Stop the current testing session.

### Request Schema

```json
{
  "type": "object",
  "properties": {
    "graceful_shutdown": {"type": "boolean", "default": true},
    "save_results": {"type": "boolean", "default": true},
    "generate_report": {"type": "boolean", "default": true}
  }
}
```

### Response Schema

```json
{
  "type": "object",
  "required": ["status", "session_summary"],
  "properties": {
    "status": {"type": "string", "enum": ["stopped", "error"]},
    "session_summary": {
      "type": "object",
      "properties": {
        "total_duration": {"type": "string"},
        "total_readings": {"type": "integer"},
        "successful_pipelines": {"type": "integer"},
        "failed_pipelines": {"type": "integer"},
        "motor_commands_generated": {"type": "integer"},
        "motor_commands_executed": {"type": "integer"},
        "overall_success_rate": {"type": "number"}
      }
    },
    "test_results": {
      "type": "object",
      "properties": {
        "scenario_name": {"type": "string"},
        "objectives_met": {"type": "array", "items": {"type": "string"}},
        "objectives_failed": {"type": "array", "items": {"type": "string"}},
        "htm_compliance": {"type": "boolean"},
        "performance_within_targets": {"type": "boolean"}
      }
    },
    "generated_files": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "type": {"type": "string", "enum": ["report", "metrics", "logs", "data"]},
          "path": {"type": "string"},
          "size": {"type": "integer"}
        }
      }
    }
  }
}
```

---

## HTM API Integration Patterns

### Sensor Data to Pipeline Processing

The sample client integrates with the HTM API using the following patterns:

1. **Sensor Reading Generation**: Each sensor type generates data according to its configuration
2. **HTM Pipeline Processing**: Sensor data is sent to `/api/v1/pipeline/process`
3. **Motor Command Execution**: Commands returned from the pipeline are executed
4. **Feedback Submission**: Execution results are sent to `/api/v1/motor-output/feedback`

### Data Flow Schema

```json
{
  "sensor_reading": {
    "type": "object",
    "properties": {
      "sensor_id": {"type": "string"},
      "type": {"type": "string"},
      "data": {"type": "object"},
      "timestamp": {"type": "string", "format": "date-time"}
    }
  },
  "htm_request": {
    "type": "object",
    "properties": {
      "sensor_data": {"$ref": "#/sensor_reading"},
      "pipeline_id": {"type": "string"},
      "processing_options": {"type": "object"}
    }
  },
  "htm_response": {
    "type": "object",
    "properties": {
      "pipeline_result": {"type": "object"},
      "motor_commands": {"type": "array"},
      "processing_summary": {"type": "object"}
    }
  },
  "command_execution": {
    "type": "object",
    "properties": {
      "command_id": {"type": "string"},
      "execution_result": {"type": "object"},
      "sensor_feedback": {"type": "object"}
    }
  }
}
```

## Sensor Type Specifications

### Temperature Sensor

Generates temperature readings with configurable range and variation patterns.

**Configuration**:
```json
{
  "type": "temperature",
  "config": {
    "base_temperature": {"type": "number", "default": 20.0},
    "variation_range": {"type": "number", "default": 5.0},
    "unit": {"type": "string", "enum": ["celsius", "fahrenheit"], "default": "celsius"},
    "noise_level": {"type": "number", "minimum": 0, "maximum": 1, "default": 0.1},
    "trend": {"type": "string", "enum": ["stable", "increasing", "decreasing", "cyclical"], "default": "stable"}
  }
}
```

### Text Sensor

Generates text data from predefined patterns or templates.

**Configuration**:
```json
{
  "type": "text",
  "config": {
    "source": {"type": "string", "enum": ["patterns", "templates", "random"], "default": "patterns"},
    "max_length": {"type": "integer", "minimum": 10, "maximum": 1000, "default": 100},
    "vocabulary_size": {"type": "integer", "default": 1000},
    "language": {"type": "string", "default": "en"},
    "complexity": {"type": "string", "enum": ["simple", "medium", "complex"], "default": "medium"}
  }
}
```

### Image Sensor

Generates synthetic images with controllable features.

**Configuration**:
```json
{
  "type": "image",
  "config": {
    "width": {"type": "integer", "minimum": 32, "maximum": 512, "default": 128},
    "height": {"type": "integer", "minimum": 32, "maximum": 512, "default": 128},
    "format": {"type": "string", "enum": ["jpeg", "png"], "default": "jpeg"},
    "pattern_type": {"type": "string", "enum": ["geometric", "natural", "noise"], "default": "geometric"},
    "color_mode": {"type": "string", "enum": ["grayscale", "rgb"], "default": "grayscale"}
  }
}
```

### Audio Sensor

Generates audio signals with configurable properties.

**Configuration**:
```json
{
  "type": "audio",
  "config": {
    "sample_rate": {"type": "integer", "enum": [8000, 16000, 44100], "default": 16000},
    "duration_ms": {"type": "integer", "minimum": 100, "maximum": 5000, "default": 1000},
    "signal_type": {"type": "string", "enum": ["sine", "noise", "chirp", "pulse"], "default": "sine"},
    "frequency_range": {"type": "object", "properties": {"min": {"type": "number"}, "max": {"type": "number"}}},
    "amplitude": {"type": "number", "minimum": 0, "maximum": 1, "default": 0.5}
  }
}
```

## Error Handling

### Client-Side Error Responses

```json
{
  "error_types": {
    "sensor_error": {
      "properties": {
        "sensor_id": {"type": "string"},
        "error_type": {"type": "string"},
        "message": {"type": "string"},
        "recoverable": {"type": "boolean"}
      }
    },
    "htm_api_error": {
      "properties": {
        "endpoint": {"type": "string"},
        "http_status": {"type": "integer"},
        "error_message": {"type": "string"},
        "retry_after": {"type": "string"}
      }
    },
    "processing_error": {
      "properties": {
        "stage": {"type": "string"},
        "error_details": {"type": "string"},
        "data_context": {"type": "object"}
      }
    }
  }
}
```

## Performance Monitoring

### Metrics Collection Schema

```json
{
  "performance_metrics": {
    "type": "object",
    "properties": {
      "sensor_metrics": {
        "type": "object",
        "properties": {
          "data_generation_rate": {"type": "number"},
          "data_quality_score": {"type": "number"},
          "error_rate": {"type": "number"}
        }
      },
      "api_metrics": {
        "type": "object",
        "properties": {
          "request_latency": {"type": "number"},
          "success_rate": {"type": "number"},
          "throughput": {"type": "number"}
        }
      },
      "pipeline_metrics": {
        "type": "object",
        "properties": {
          "end_to_end_latency": {"type": "number"},
          "htm_compliance_score": {"type": "number"},
          "prediction_accuracy": {"type": "number"}
        }
      },
      "motor_output_metrics": {
        "type": "object",
        "properties": {
          "commands_per_second": {"type": "number"},
          "execution_success_rate": {"type": "number"},
          "feedback_accuracy": {"type": "number"}
        }
      }
    }
  }
}
```