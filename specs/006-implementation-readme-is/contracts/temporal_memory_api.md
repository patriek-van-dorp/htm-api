# API Contracts: Temporal Memory Endpoints

## POST /api/v1/temporal-memory/process

Process spatial pooler output through temporal memory to generate temporal sequences and predictions.

### Request Schema

```json
{
  "type": "object",
  "required": ["input", "instance_id"],
  "properties": {
    "input": {
      "type": "object",
      "required": ["spatial_sdr", "metadata"],
      "properties": {
        "spatial_sdr": {
          "type": "object",
          "required": ["width", "active_bits"],
          "properties": {
            "width": {
              "type": "integer",
              "minimum": 1,
              "maximum": 10000,
              "description": "Total number of bits in spatial SDR"
            },
            "active_bits": {
              "type": "array",
              "items": {
                "type": "integer",
                "minimum": 0
              },
              "minItems": 1,
              "maxItems": 500,
              "description": "Indices of active bits (2-5% sparsity)"
            },
            "sparsity": {
              "type": "number",
              "minimum": 0.015,
              "maximum": 0.05,
              "description": "Actual sparsity level (validation)"
            }
          }
        },
        "metadata": {
          "type": "object",
          "properties": {
            "sensor_id": {"type": "string"},
            "timestamp": {"type": "string", "format": "date-time"},
            "sequence_id": {"type": "string"},
            "step_number": {"type": "integer", "minimum": 0}
          }
        }
      }
    },
    "instance_id": {
      "type": "string",
      "pattern": "^[a-zA-Z0-9-_]{1,64}$",
      "description": "Temporal memory instance identifier"
    },
    "config_override": {
      "type": "object",
      "properties": {
        "learning_enabled": {"type": "boolean"},
        "prediction_steps": {"type": "integer", "minimum": 1, "maximum": 10}
      }
    }
  }
}
```

### Response Schema

```json
{
  "type": "object",
  "required": ["temporal_result", "predictions", "status", "processing_time"],
  "properties": {
    "temporal_result": {
      "type": "object",
      "required": ["active_cells", "predictive_cells", "winner_cells"],
      "properties": {
        "active_cells": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "column": {"type": "integer"},
              "cell": {"type": "integer"}
            }
          }
        },
        "predictive_cells": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "column": {"type": "integer"},
              "cell": {"type": "integer"}
            }
          }
        },
        "winner_cells": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "column": {"type": "integer"},
              "cell": {"type": "integer"}
            }
          }
        }
      }
    },
    "predictions": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "predicted_sdr": {
            "type": "object",
            "properties": {
              "width": {"type": "integer"},
              "active_bits": {"type": "array", "items": {"type": "integer"}}
            }
          },
          "confidence": {"type": "number", "minimum": 0, "maximum": 1},
          "steps_ahead": {"type": "integer", "minimum": 1}
        }
      }
    },
    "status": {
      "type": "string",
      "enum": ["success", "warning", "error"]
    },
    "processing_time": {
      "type": "number",
      "minimum": 0,
      "description": "Processing time in milliseconds"
    },
    "htm_properties": {
      "type": "object",
      "properties": {
        "cell_activation_rate": {"type": "number"},
        "prediction_accuracy": {"type": "number"},
        "sequence_stability": {"type": "number"}
      }
    },
    "errors": {
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
  "error": "invalid_input",
  "message": "Spatial SDR sparsity outside HTM compliance range",
  "details": {
    "actual_sparsity": 0.12,
    "required_range": "0.015-0.05"
  }
}
```

**404 Not Found**:
```json
{
  "error": "instance_not_found",
  "message": "Temporal memory instance not found",
  "instance_id": "invalid-instance"
}
```

**500 Internal Server Error**:
```json
{
  "error": "processing_error",
  "message": "Temporal memory processing failed",
  "details": "Cell activation threshold exceeded"
}
```

---

## GET /api/v1/temporal-memory/config

Retrieve current temporal memory configuration.

### Request Parameters

- `instance_id` (query, required): Temporal memory instance identifier

### Response Schema

```json
{
  "type": "object",
  "required": ["config", "status"],
  "properties": {
    "config": {
      "type": "object",
      "properties": {
        "cells_per_column": {"type": "integer", "minimum": 1, "maximum": 64},
        "activation_threshold": {"type": "integer", "minimum": 1},
        "learning_threshold": {"type": "integer", "minimum": 1},
        "initial_permanence": {"type": "number", "minimum": 0, "maximum": 1},
        "connected_permanence": {"type": "number", "minimum": 0, "maximum": 1},
        "min_threshold": {"type": "integer", "minimum": 1},
        "max_new_synapse_count": {"type": "integer", "minimum": 1},
        "permanence_increment": {"type": "number", "minimum": 0, "maximum": 0.5},
        "permanence_decrement": {"type": "number", "minimum": 0, "maximum": 0.5},
        "predicted_segment_decrement": {"type": "number", "minimum": 0, "maximum": 0.5}
      }
    },
    "status": {"type": "string", "enum": ["success"]},
    "instance_info": {
      "type": "object",
      "properties": {
        "id": {"type": "string"},
        "created_at": {"type": "string", "format": "date-time"},
        "total_cells": {"type": "integer"},
        "total_segments": {"type": "integer"},
        "total_synapses": {"type": "integer"}
      }
    }
  }
}
```

---

## PUT /api/v1/temporal-memory/config

Update temporal memory configuration.

### Request Schema

```json
{
  "type": "object",
  "required": ["instance_id", "config"],
  "properties": {
    "instance_id": {"type": "string"},
    "config": {
      "type": "object",
      "properties": {
        "cells_per_column": {"type": "integer", "minimum": 1, "maximum": 64},
        "activation_threshold": {"type": "integer", "minimum": 1, "maximum": 50},
        "learning_threshold": {"type": "integer", "minimum": 1, "maximum": 50},
        "initial_permanence": {"type": "number", "minimum": 0.1, "maximum": 0.9},
        "connected_permanence": {"type": "number", "minimum": 0.1, "maximum": 0.9},
        "permanence_increment": {"type": "number", "minimum": 0.01, "maximum": 0.5},
        "permanence_decrement": {"type": "number", "minimum": 0.01, "maximum": 0.5}
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
      "$ref": "#/components/schemas/TemporalMemoryConfig"
    },
    "validation_errors": {
      "type": "array",
      "items": {"type": "string"}
    }
  }
}
```

---

## GET /api/v1/temporal-memory/status

Get temporal memory instance status and health metrics.

### Request Parameters

- `instance_id` (query, required): Temporal memory instance identifier

### Response Schema

```json
{
  "type": "object",
  "required": ["status", "health", "metrics"],
  "properties": {
    "status": {"type": "string", "enum": ["healthy", "degraded", "error"]},
    "health": {
      "type": "object",
      "properties": {
        "cell_activation_rate": {"type": "number", "minimum": 0, "maximum": 1},
        "prediction_accuracy": {"type": "number", "minimum": 0, "maximum": 1},
        "sequence_stability": {"type": "number", "minimum": 0, "maximum": 1},
        "memory_usage_mb": {"type": "number", "minimum": 0},
        "processing_latency_ms": {"type": "number", "minimum": 0}
      }
    },
    "metrics": {
      "type": "object",
      "properties": {
        "total_processed": {"type": "integer", "minimum": 0},
        "sequences_learned": {"type": "integer", "minimum": 0},
        "predictions_generated": {"type": "integer", "minimum": 0},
        "avg_confidence": {"type": "number", "minimum": 0, "maximum": 1},
        "error_rate": {"type": "number", "minimum": 0, "maximum": 1}
      }
    },
    "instance_info": {
      "type": "object",
      "properties": {
        "id": {"type": "string"},
        "uptime_seconds": {"type": "integer"},
        "last_activity": {"type": "string", "format": "date-time"}
      }
    }
  }
}
```