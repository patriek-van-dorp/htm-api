# API Contracts: Motor Output Endpoints

## POST /api/v1/motor-output/process

Convert temporal memory predictions into actionable motor commands.

### Request Schema

```json
{
  "type": "object",
  "required": ["predictions", "instance_id"],
  "properties": {
    "predictions": {
      "type": "array",
      "minItems": 1,
      "maxItems": 10,
      "items": {
        "type": "object",
        "required": ["predicted_sdr", "confidence", "steps_ahead"],
        "properties": {
          "predicted_sdr": {
            "type": "object",
            "required": ["width", "active_bits"],
            "properties": {
              "width": {"type": "integer", "minimum": 1, "maximum": 10000},
              "active_bits": {
                "type": "array",
                "items": {"type": "integer", "minimum": 0},
                "minItems": 1,
                "maxItems": 500
              }
            }
          },
          "confidence": {
            "type": "number",
            "minimum": 0,
            "maximum": 1,
            "description": "Prediction confidence from temporal memory"
          },
          "steps_ahead": {
            "type": "integer",
            "minimum": 1,
            "maximum": 10,
            "description": "Prediction horizon in time steps"
          },
          "context": {
            "type": "object",
            "properties": {
              "sensor_id": {"type": "string"},
              "sequence_id": {"type": "string"},
              "temporal_context": {"type": "object"}
            }
          }
        }
      }
    },
    "instance_id": {
      "type": "string",
      "pattern": "^[a-zA-Z0-9-_]{1,64}$",
      "description": "Motor output instance identifier"
    },
    "action_constraints": {
      "type": "object",
      "properties": {
        "allowed_types": {
          "type": "array",
          "items": {"type": "string", "enum": ["movement", "audio", "visual", "control"]}
        },
        "max_commands": {"type": "integer", "minimum": 1, "maximum": 10},
        "priority_threshold": {"type": "integer", "minimum": 1, "maximum": 10}
      }
    }
  }
}
```

### Response Schema

```json
{
  "type": "object",
  "required": ["motor_commands", "status", "processing_time"],
  "properties": {
    "motor_commands": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["id", "type", "parameters", "confidence", "priority"],
        "properties": {
          "id": {"type": "string"},
          "type": {"type": "string", "enum": ["movement", "audio", "visual", "control"]},
          "parameters": {
            "type": "object",
            "description": "Type-specific command parameters",
            "oneOf": [
              {
                "title": "Movement Command",
                "properties": {
                  "direction": {"type": "string", "enum": ["forward", "backward", "left", "right", "up", "down"]},
                  "magnitude": {"type": "number", "minimum": 0, "maximum": 1},
                  "duration_ms": {"type": "integer", "minimum": 100, "maximum": 5000}
                }
              },
              {
                "title": "Audio Command", 
                "properties": {
                  "frequency": {"type": "number", "minimum": 20, "maximum": 20000},
                  "amplitude": {"type": "number", "minimum": 0, "maximum": 1},
                  "duration_ms": {"type": "integer", "minimum": 100, "maximum": 5000},
                  "pattern": {"type": "string"}
                }
              },
              {
                "title": "Visual Command",
                "properties": {
                  "color": {"type": "string", "pattern": "^#[0-9A-Fa-f]{6}$"},
                  "brightness": {"type": "number", "minimum": 0, "maximum": 1},
                  "pattern": {"type": "string"},
                  "duration_ms": {"type": "integer", "minimum": 100, "maximum": 5000}
                }
              },
              {
                "title": "Control Command",
                "properties": {
                  "signal": {"type": "string"},
                  "value": {"type": "number"},
                  "target": {"type": "string"}
                }
              }
            ]
          },
          "confidence": {"type": "number", "minimum": 0, "maximum": 1},
          "priority": {"type": "integer", "minimum": 1, "maximum": 10},
          "created_at": {"type": "string", "format": "date-time"},
          "deadline": {"type": "string", "format": "date-time"},
          "source_prediction": {
            "type": "object",
            "properties": {
              "prediction_id": {"type": "string"},
              "steps_ahead": {"type": "integer"}
            }
          }
        }
      }
    },
    "status": {"type": "string", "enum": ["success", "partial_success", "error"]},
    "processing_time": {"type": "number", "minimum": 0},
    "command_summary": {
      "type": "object",
      "properties": {
        "total_commands": {"type": "integer"},
        "by_type": {
          "type": "object",
          "properties": {
            "movement": {"type": "integer"},
            "audio": {"type": "integer"},
            "visual": {"type": "integer"},
            "control": {"type": "integer"}
          }
        },
        "avg_confidence": {"type": "number"},
        "max_priority": {"type": "integer"}
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
  "error": "invalid_predictions",
  "message": "Prediction confidence below motor output threshold",
  "details": {
    "min_confidence": 0.75,
    "received_confidence": 0.45
  }
}
```

**429 Too Many Requests**:
```json
{
  "error": "command_limit_exceeded",
  "message": "Maximum concurrent commands exceeded",
  "details": {
    "max_concurrent": 5,
    "current_active": 5
  }
}
```

---

## POST /api/v1/motor-output/feedback

Provide feedback on motor command execution for learning and adaptation.

### Request Schema

```json
{
  "type": "object",
  "required": ["command_id", "execution_result"],
  "properties": {
    "command_id": {
      "type": "string",
      "description": "ID of the executed motor command"
    },
    "execution_result": {
      "type": "object",
      "required": ["success", "execution_time"],
      "properties": {
        "success": {"type": "boolean"},
        "execution_time": {"type": "number", "minimum": 0},
        "error_details": {"type": "string"},
        "actuator_response": {
          "type": "object",
          "description": "Actuator-specific response data"
        },
        "performance_metrics": {
          "type": "object",
          "properties": {
            "accuracy": {"type": "number", "minimum": 0, "maximum": 1},
            "efficiency": {"type": "number", "minimum": 0, "maximum": 1},
            "user_satisfaction": {"type": "number", "minimum": 0, "maximum": 1}
          }
        }
      }
    },
    "sensor_feedback": {
      "type": "object",
      "description": "Optional sensor readings after command execution",
      "properties": {
        "sensor_id": {"type": "string"},
        "reading": {"type": "object"},
        "timestamp": {"type": "string", "format": "date-time"}
      }
    },
    "learning_context": {
      "type": "object",
      "properties": {
        "environment_state": {"type": "object"},
        "expected_outcome": {"type": "object"},
        "actual_outcome": {"type": "object"}
      }
    }
  }
}
```

### Response Schema

```json
{
  "type": "object",
  "required": ["status", "learning_update"],
  "properties": {
    "status": {"type": "string", "enum": ["feedback_received", "learning_updated", "error"]},
    "learning_update": {
      "type": "object",
      "properties": {
        "action_mapping_updated": {"type": "boolean"},
        "confidence_adjustment": {"type": "number"},
        "pattern_strength_change": {"type": "number"},
        "new_associations_formed": {"type": "integer"}
      }
    },
    "command_history": {
      "type": "object",
      "properties": {
        "command_id": {"type": "string"},
        "success_rate": {"type": "number"},
        "avg_execution_time": {"type": "number"},
        "total_executions": {"type": "integer"}
      }
    },
    "prediction_accuracy": {
      "type": "object",
      "properties": {
        "temporal_accuracy": {"type": "number"},
        "action_relevance": {"type": "number"},
        "outcome_prediction": {"type": "number"}
      }
    }
  }
}
```

---

## GET /api/v1/motor-output/config

Retrieve motor output configuration.

### Request Parameters

- `instance_id` (query, required): Motor output instance identifier

### Response Schema

```json
{
  "type": "object",
  "required": ["config", "status"],
  "properties": {
    "config": {
      "type": "object",
      "properties": {
        "action_types": {
          "type": "array",
          "items": {"type": "string"}
        },
        "confidence_threshold": {"type": "number", "minimum": 0.5, "maximum": 0.95, "default": 0.75},
        "max_concurrent_commands": {"type": "integer", "minimum": 1},
        "command_timeout": {"type": "string"},
        "feedback_timeout": {"type": "string"},
        "learning_enabled": {"type": "boolean"},
        "default_actions": {"type": "object"}
      }
    },
    "status": {"type": "string", "enum": ["success"]},
    "action_mappings": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "pattern": {"type": "string"},
          "action_type": {"type": "string"},
          "confidence": {"type": "number"},
          "success_rate": {"type": "number"}
        }
      }
    }
  }
}
```

---

## GET /api/v1/motor-output/status

Get motor output status and active commands.

### Request Parameters

- `instance_id` (query, required): Motor output instance identifier

### Response Schema

```json
{
  "type": "object",
  "required": ["status", "active_commands", "metrics"],
  "properties": {
    "status": {"type": "string", "enum": ["ready", "busy", "error"]},
    "active_commands": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "command_id": {"type": "string"},
          "type": {"type": "string"},
          "status": {"type": "string", "enum": ["pending", "executing", "completed", "failed"]},
          "created_at": {"type": "string", "format": "date-time"},
          "execution_start": {"type": "string", "format": "date-time"},
          "deadline": {"type": "string", "format": "date-time"}
        }
      }
    },
    "metrics": {
      "type": "object",
      "properties": {
        "total_commands_generated": {"type": "integer"},
        "commands_executed": {"type": "integer"},
        "success_rate": {"type": "number", "minimum": 0, "maximum": 1},
        "avg_execution_time": {"type": "number"},
        "avg_confidence": {"type": "number", "minimum": 0, "maximum": 1},
        "learning_accuracy": {"type": "number", "minimum": 0, "maximum": 1}
      }
    },
    "performance": {
      "type": "object",
      "properties": {
        "command_generation_latency": {"type": "number"},
        "feedback_processing_time": {"type": "number"},
        "memory_usage_mb": {"type": "number"}
      }
    }
  }
}
```