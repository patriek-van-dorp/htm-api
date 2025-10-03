# Data Model: Complete HTM Pipeline with Sensor-to-Motor Integration

**Feature**: Complete HTM Pipeline Implementation  
**Date**: October 2, 2025  
**Status**: Phase 1 - Design Artifacts

## Core Entities

### 1. Temporal Memory Components

#### TemporalMemory
Represents the HTM temporal memory processing engine that learns sequences from spatial pooler output.

**Fields**:
- `ID`: string - Unique temporal memory instance identifier
- `Config`: TemporalMemoryConfig - Configuration parameters
- `Cells`: [][]Cell - 2D array of cells (columns × cells_per_column)
- `ActiveCells`: []CellIndex - Currently active cells
- `PredictiveCells`: []CellIndex - Cells in predictive state
- `WinnerCells`: []CellIndex - Winner cells from current processing
- `Connections`: ConnectionMatrix - Distal dendrite connections between cells
- `Metrics`: TemporalMemoryMetrics - Performance and HTM compliance metrics
- `State`: ProcessingState - Current processing state
- `CreatedAt`: time.Time - Instance creation timestamp
- `UpdatedAt`: time.Time - Last processing timestamp

**Relationships**:
- Receives input from SpatialPooler (existing)
- Outputs predictions to MotorOutput
- Maintains synaptic connections between processing cycles

**Validation Rules**:
- ActiveCells must maintain HTM sparsity constraints (2-5% of total cells)
- PredictiveCells must be subset of non-active cells
- ConnectionMatrix must maintain biological synapse limits
- Config values must validate against HTM theoretical constraints

#### TemporalMemoryConfig
Configuration parameters for temporal memory processing.

**Fields**:
- `CellsPerColumn`: int - Number of cells per spatial pooler column (default: 32)
- `ActivationThreshold`: int - Minimum synapses for cell activation (default: 13)
- `LearningThreshold`: int - Minimum synapses for learning (default: 10)
- `InitialPermanence`: float64 - Initial synapse permanence (default: 0.21)
- `ConnectedPermanence`: float64 - Permanence threshold for connected synapses (default: 0.50)
- `MinThreshold`: int - Minimum synapses for cell prediction (default: 9)
- `MaxNewSynapseCount`: int - Maximum new synapses per learning cycle (default: 20)
- `PermanenceIncrement`: float64 - Learning increment value (default: 0.10)
- `PermanenceDecrement`: float64 - Learning decrement value (default: 0.10)
- `PredictedSegmentDecrement`: float64 - Decrement for incorrectly predicted segments (default: 0.02)
- `MaxSegmentsPerCell`: int - Maximum distal segments per cell (default: 255)
- `MaxSynapsesPerSegment`: int - Maximum synapses per segment (default: 255)

**Validation Rules**:
- CellsPerColumn: 1-64 (biological constraint)
- Permanence values: 0.0-1.0 range
- Threshold values: positive integers
- Learning rates: 0.0-0.5 range for stability

#### Cell
Individual processing unit within temporal memory columns.

**Fields**:
- `Index`: CellIndex - Unique cell identifier (column, cell_index)
- `Active`: bool - Current activation state
- `Predictive`: bool - Current predictive state
- `Winner`: bool - Winner cell status for current cycle
- `DistalSegments`: []Segment - Distal dendrite segments
- `LastActiveCycle`: int - Cycle number of last activation
- `ActivationCount`: int64 - Total activation count for metrics

**State Transitions**:
- Inactive → Active: When sufficient distal synapse support
- Active → Predictive: When temporal patterns indicate future activation
- Predictive → Active: When prediction is confirmed
- Any → Winner: When selected for learning during active state

#### Segment
Represents distal dendrite segments that detect temporal context.

**Fields**:
- `ID`: string - Unique segment identifier
- `Synapses`: []Synapse - Synaptic connections to other cells
- `CreationCycle`: int - Cycle when segment was created
- `LastActiveCycle`: int - Last cycle with sufficient active synapses
- `ActivationCount`: int64 - Total activation count

#### Synapse
Individual synaptic connection between cells.

**Fields**:
- `PresynapticCell`: CellIndex - Source cell for connection
- `Permanence`: float64 - Connection strength (0.0-1.0)
- `CreationCycle`: int - When synapse was formed
- `LastActiveCycle`: int - Last time presynaptic cell was active

### 2. Motor Output Components

#### MotorOutput
Represents the motor output processing engine that converts temporal memory predictions into actionable commands.

**Fields**:
- `ID`: string - Unique motor output instance identifier
- `Config`: MotorOutputConfig - Configuration parameters
- `ActionMappings`: []ActionMapping - Prediction-to-action mappings
- `ActiveCommands`: []Command - Currently executing commands
- `CommandHistory`: []CommandExecution - Recent command history
- `Metrics`: MotorOutputMetrics - Performance metrics
- `State`: ProcessingState - Current processing state
- `CreatedAt`: time.Time - Instance creation timestamp
- `UpdatedAt`: time.Time - Last processing timestamp

**Relationships**:
- Receives predictions from TemporalMemory
- Generates commands for external actuators
- Receives feedback for learning and adaptation

#### MotorOutputConfig
Configuration parameters for motor output processing.

**Fields**:
- `ActionTypes`: []string - Supported action categories (movement, audio, visual, control)
- `ConfidenceThreshold`: float64 - Minimum prediction confidence for action (default: 0.75, range: 0.5-0.95, HTM-tuned for biological timing)
- `MaxConcurrentCommands`: int - Maximum simultaneous commands (default: 5)
- `CommandTimeout`: time.Duration - Maximum command execution time (default: 5s)
- `FeedbackTimeout`: time.Duration - Maximum feedback wait time (default: 10s)
- `LearningEnabled`: bool - Whether to learn from feedback (default: true)
- `DefaultActions`: map[string]interface{} - Default actions for low-confidence predictions

#### Command
Represents an actionable command generated from temporal memory predictions.

**Fields**:
- `ID`: string - Unique command identifier
- `Type`: ActionType - Action category (movement, audio, visual, control)
- `Parameters`: map[string]interface{} - Type-specific parameters
- `Confidence`: float64 - Prediction confidence that generated command
- `Priority`: int - Execution priority (1-10, higher = more important)
- `CreatedAt`: time.Time - Command generation timestamp
- `Deadline`: time.Time - Latest acceptable execution time
- `Status`: CommandStatus - Current execution status
- `Source`: PredictionSource - Temporal memory prediction that generated command

#### CommandExecution
Records the execution and feedback for motor output commands.

**Fields**:
- `CommandID`: string - Reference to executed command
- `StartTime`: time.Time - Execution start timestamp
- `EndTime`: time.Time - Execution completion timestamp
- `Success`: bool - Whether command executed successfully
- `Feedback`: map[string]interface{} - Execution feedback data
- `ErrorDetails`: string - Error information if execution failed
- `Duration`: time.Duration - Total execution time
- `ActuatorResponse`: interface{} - Response from target actuator

### 3. Pipeline Orchestration

#### HTMPipeline
Represents the complete sensor-to-motor processing pipeline.

**Fields**:
- `ID`: string - Unique pipeline instance identifier
- `Config`: PipelineConfig - Pipeline configuration
- `SpatialPooler`: SpatialPoolerReference - Reference to spatial pooler instance
- `TemporalMemory`: TemporalMemory - Temporal memory instance
- `MotorOutput`: MotorOutput - Motor output instance
- `ActiveSessions`: []ProcessingSession - Currently active processing sessions
- `Metrics`: PipelineMetrics - End-to-end performance metrics
- `State`: ProcessingState - Overall pipeline state
- `CreatedAt`: time.Time - Pipeline creation timestamp
- `UpdatedAt`: time.Time - Last activity timestamp

#### ProcessingSession
Represents a single sensor-to-motor processing session.

**Fields**:
- `ID`: string - Unique session identifier
- `SensorID`: string - Source sensor identifier
- `InputData`: interface{} - Raw sensor input
- `SpatialPoolerOutput`: SDR - Spatial pooler result
- `TemporalMemoryOutput`: Predictions - Temporal memory result
- `MotorCommands`: []Command - Generated motor commands
- `StartTime`: time.Time - Session start timestamp
- `EndTime`: time.Time - Session completion timestamp
- `Duration`: time.Duration - Total processing time
- `Success`: bool - Whether session completed successfully
- `Errors`: []string - Any processing errors

### 4. Sample Client Components

#### SampleClient
Represents the sample client application for testing HTM pipeline.

**Fields**:
- `ID`: string - Unique client instance identifier
- `Config`: ClientConfig - Client configuration
- `Sensors`: []SensorInstance - Available sensor instances
- `APIClient`: HTTPClient - HTM API communication client
- `TestScenarios`: []TestScenario - Available test scenarios
- `ActiveSessions`: []ClientSession - Currently running sessions
- `Metrics`: ClientMetrics - Performance and success metrics
- `State`: ClientState - Current operational state

#### SensorInstance
Represents a sensor instance within the sample client.

**Fields**:
- `ID`: string - Unique sensor identifier
- `Type`: SensorType - Sensor category (temperature, text, image, audio)
- `Config`: SensorConfig - Sensor-specific configuration
- `Generator`: DataGenerator - Data generation logic
- `State`: SensorState - Current sensor state
- `LastReading`: SensorReading - Most recent sensor data
- `ReadingHistory`: []SensorReading - Recent reading history
- `Metrics`: SensorMetrics - Sensor performance metrics

#### TestScenario
Represents automated test scenarios for HTM pipeline validation.

**Fields**:
- `ID`: string - Unique scenario identifier
- `Name`: string - Human-readable scenario name
- `Description`: string - Scenario purpose and expectations
- `SensorTypes`: []SensorType - Required sensor types
- `Duration`: time.Duration - Expected scenario duration
- `Steps`: []TestStep - Ordered test steps
- `ExpectedOutcomes`: []Outcome - Expected results
- `Validation`: ValidationCriteria - Success/failure criteria

## Entity Relationships

```
SensorData → Encoder → SpatialPooler → TemporalMemory → MotorOutput → Commands
     ↑                                       ↓
SampleClient ←---- Feedback ←---- CommandExecution
     ↓
TestScenarios → ValidationResults
```

## Data Flow Patterns

### 1. Forward Processing Flow
1. **Sensor Input**: Raw data from sample client sensors
2. **Encoding**: Convert to SDR format (existing encoder functionality)  
3. **Spatial Pooling**: Generate normalized HTM SDR (existing implementation)
4. **Temporal Memory**: Process temporal sequences, generate predictions
5. **Motor Output**: Convert predictions to actionable commands
6. **Command Execution**: Execute commands via sample client actuators

### 2. Feedback Learning Flow
1. **Command Execution**: Motor commands executed by sample client
2. **Feedback Collection**: Success/failure results collected
3. **Learning Update**: Temporal memory and motor output learn from results
4. **Adaptation**: Future predictions and commands improve based on feedback

### 3. Validation Flow
1. **Test Scenario Execution**: Automated scenarios run via sample client
2. **Pipeline Processing**: Complete sensor-to-motor processing
3. **Result Validation**: Outcomes compared to expected results
4. **Metrics Collection**: Performance and accuracy metrics recorded
5. **HTM Compliance Check**: Biological constraints validated

## State Management

### Processing States
- `Inactive`: Component not processing
- `Active`: Currently processing data
- `Learning`: Actively updating internal state
- `Error`: Processing failed, requires intervention
- `Maintenance`: Temporary unavailable for updates

### Session States
- `Created`: Session initialized, ready to process
- `Processing`: Data flowing through pipeline stages
- `Completed`: Session finished successfully
- `Failed`: Session encountered unrecoverable error
- `Timeout`: Session exceeded maximum processing time

## Performance Considerations

### Memory Management
- **Temporal Memory**: Estimated 2-3x spatial pooler memory usage
- **Motor Output**: Minimal memory overhead, primarily computation-bound
- **Pipeline Sessions**: Limited concurrent sessions to prevent memory exhaustion
- **History Retention**: Configurable retention periods for metrics and history

### Concurrency Patterns
- **Sensor Processing**: Multiple sensors can process simultaneously
- **Pipeline Stages**: Sequential processing within each session
- **Command Execution**: Parallel motor command execution
- **Feedback Processing**: Asynchronous feedback integration

## Validation Rules Summary

### HTM Compliance
- Temporal memory sparsity: 2-5% active cells
- Sequence learning: Proper temporal continuity
- Prediction accuracy: Measurable improvement over time
- Biological constraints: Cell counts, connection limits, learning rates

### Performance Requirements
- End-to-end latency: <100ms for complete pipeline
- Concurrent sensors: Support up to 25 simultaneous sensors
- Memory usage: <500MB for complete pipeline
- Throughput: Handle real-time sensor data rates

### API Consistency
- REST endpoints follow existing spatial pooler patterns
- Request/response schemas maintain backward compatibility
- Error handling consistent with existing API behavior
- Authentication and authorization inherit from existing implementation

---

*Data model completed: October 2, 2025*  
*Next: API contracts and quickstart guide*