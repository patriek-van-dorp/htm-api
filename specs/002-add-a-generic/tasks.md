# Tasks: Generic HTM Sensor Package

**Input**: Design documents from `/specs/002-add-a-generic/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/, quickstart.md

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Tech stack: Go 1.21+, gonum, testify
   → Structure: Single library package with internal/pkg separation
2. Load design documents:
   → data-model.md: 8 core entities identified
   → contracts/interfaces.md: 4 core interfaces + supporting types
   → quickstart.md: Integration scenarios and usage examples
   → research.md: Performance and validation requirements
3. Generate tasks by category:
   → Setup: Go module, dependencies, tooling
   → Tests: Contract tests, performance tests, integration tests
   → Core: SDR, interfaces, encoders, registry, validation
   → Integration: Multi-sensor scenarios, performance validation
   → Polish: Documentation, benchmarks, examples
4. Apply task rules:
   → Different packages/files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
   → Performance requirements: sub-millisecond encoding
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph with performance focus
7. Create parallel execution examples
8. SUCCESS: 42 tasks ready for execution
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **Performance Critical**: Sub-millisecond encoding requirement
- **Single-threaded**: No concurrency support needed
- **Silent Failure**: Empty SDR on encoding failures

## Path Conventions (Single Go Library)
- **Core**: `internal/sensors/` for implementation
- **Public API**: `pkg/sensors/` for exports
- **Tests**: `tests/contract/`, `tests/integration/`, `tests/unit/`

## Phase 3.1: Setup & Foundation
- [X] T001 Create Go module structure per implementation plan (go.mod, directory layout)
- [X] T002 Initialize Go dependencies (gonum, testify, performance profiling tools)
- [X] T003 [P] Configure Go tooling (linting with golangci-lint, formatting with gofmt)
- [X] T004 [P] Setup performance benchmarking framework for sub-millisecond validation

## Phase 3.2: Contract Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Interface Contract Tests
- [X] T005 [P] SDR interface contract test in tests/contract/sdr_interface_test.go
- [X] T006 [P] SensorInterface contract test in tests/contract/sensor_interface_test.go  
- [X] T007 [P] SensorConfig contract test in tests/contract/sensor_config_test.go
- [X] T008 [P] SensorRegistry contract test in tests/contract/sensor_registry_test.go

### Performance Contract Tests (Sub-millisecond requirement)
- [X] T009 [P] Numeric encoder performance contract test in tests/contract/numeric_performance_test.go
- [X] T010 [P] Categorical encoder performance contract test in tests/contract/categorical_performance_test.go
- [X] T011 [P] Text encoder performance contract test in tests/contract/text_performance_test.go
- [X] T012 [P] Spatial encoder performance contract test in tests/contract/spatial_performance_test.go

### Integration Test Scenarios (From quickstart.md)
- [X] T013 [P] Registry setup and sensor registration test in tests/integration/registry_setup_test.go
- [X] T014 [P] Numeric data encoding pipeline test in tests/integration/numeric_pipeline_test.go
- [X] T015 [P] Categorical data encoding pipeline test in tests/integration/categorical_pipeline_test.go
- [X] T016 [P] Text data encoding pipeline test in tests/integration/text_pipeline_test.go
- [X] T017 [P] Multi-sensor sequential processing test in tests/integration/multi_sensor_test.go
- [X] T017A [P] Batch processing integration test - validate processing arrays of inputs with consistent SDR output, memory efficiency for large batches, and performance scaling in tests/integration/batch_processing_test.go
- [X] T017B [P] Sensor composition integration test - validate combining multiple sensor types (numeric+text, spatial+categorical) with weighted SDR fusion, conflict resolution, and semantic coherence in tests/integration/sensor_composition_test.go
- [X] T017C [P] Multi-SDR collection encoding test - validate spatial subdivision of large inputs (images) into coherent SDR collections, topology preservation, and hierarchical encoding patterns in tests/integration/multi_sdr_test.go
- [X] T018 [P] Silent failure mode validation test in tests/integration/silent_failure_test.go
- [X] T019 [P] Input size limit (1MB) validation test in tests/integration/size_limit_test.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### SDR Core Implementation
- [X] T020 [P] SDR representation struct in internal/sensors/sdr/representation.go
- [X] T021 [P] SDR sparsity management in internal/sensors/sdr/sparsity.go
- [X] T022 [P] SDR similarity calculations in internal/sensors/sdr/similarity.go

### Core Interfaces Implementation  
- [X] T023 [P] Core sensor interfaces in internal/sensors/interfaces.go
- [X] T024 [P] Sensor configuration implementation in internal/sensors/config.go
- [X] T025 [P] SDR validation logic in internal/sensors/validation.go

### Sensor Registry Implementation
- [X] T026 Sensor registry with factory pattern in internal/sensors/registry.go

### High-Performance Encoder Implementation
- [ ] T027 [P] Numeric encoder with sub-ms optimization in internal/sensors/encoders/numeric.go
- [ ] T028 [P] Categorical encoder with hash collision handling in internal/sensors/encoders/categorical.go
- [ ] T029 [P] Text encoder with 1MB document support in internal/sensors/encoders/text.go
- [ ] T030 [P] Spatial encoder with topology preservation in internal/sensors/encoders/spatial.go
- [ ] T030A [P] Multi-SDR spatial encoder for image subdivision in internal/sensors/encoders/spatial_multi.go

### Public API Layer
- [ ] T031 Public API exports and convenience functions in pkg/sensors/public_api.go
- [ ] T031A [P] Batch processing API - implement efficient batch encoding with configurable batch sizes, memory pooling, progress tracking, and error aggregation for multiple inputs in pkg/sensors/batch_api.go
- [ ] T031B [P] Sensor composition utilities - implement weighted sensor fusion, conflict resolution strategies, and semantic coherence validation for multi-SDR processing in pkg/sensors/composition.go

## Phase 3.4: Integration & Validation

### Single-threaded Integration
- [ ] T032 Sequential multi-sensor pipeline implementation in internal/sensors/pipeline.go
- [ ] T032A Batch processing implementation - implement memory-efficient batch operations with configurable chunk sizes, progress tracking, error collection, and memory usage optimization for encoder collections in internal/sensors/batch.go
- [ ] T032B Sensor composition framework - implement weighted combination strategies, semantic conflict detection, coherence validation, and fusion algorithms for multi-SDR input handling in internal/sensors/composition.go
- [ ] T033A Silent failure error handling in SDR operations (internal/sensors/sdr/)
- [ ] T033B Silent failure error handling in encoder implementations (internal/sensors/encoders/)
- [ ] T033C Silent failure error handling in registry operations (internal/sensors/registry.go)
- [ ] T033D Silent failure error handling in validation logic (internal/sensors/validation.go)
- [ ] T034A Input size validation enforcement in numeric encoder (internal/sensors/encoders/numeric.go)
- [ ] T034B Input size validation enforcement in categorical encoder (internal/sensors/encoders/categorical.go)
- [ ] T034C Input size validation enforcement in text encoder (internal/sensors/encoders/text.go)
- [ ] T034D Input size validation enforcement in spatial encoder (internal/sensors/encoders/spatial.go)

### Performance Validation (Critical Requirements)
- [ ] T035 Sub-millisecond encoding benchmark suite in tests/integration/performance_test.go
- [ ] T035A Batch processing performance benchmarks in tests/integration/batch_performance_test.go
- [ ] T035B Multi-SDR collection performance validation in tests/integration/multi_sdr_performance_test.go
- [ ] T036 Memory usage profiling for 1MB inputs in tests/integration/memory_test.go
- [ ] T037 Sparsity validation across all encoder types in tests/integration/sparsity_validation_test.go

## Phase 3.5: Polish & Documentation

### Unit Test Coverage
- [ ] T038 [P] Numeric encoder edge cases in tests/unit/numeric_encoder_test.go
- [ ] T039 [P] Categorical encoder edge cases in tests/unit/categorical_encoder_test.go  
- [ ] T040 [P] Text encoder edge cases in tests/unit/text_encoder_test.go
- [ ] T041 [P] SDR operations edge cases in tests/unit/sdr_operations_test.go
- [ ] T041A [P] Batch processing edge cases in tests/unit/batch_test.go
- [ ] T041B [P] Sensor composition edge cases in tests/unit/composition_test.go
- [ ] T041C [P] Multi-SDR collection edge cases in tests/unit/multi_sdr_test.go

### Final Integration
- [ ] T042 Complete quickstart.md validation and example testing

## Dependencies

### Phase Dependencies
- Setup (T001-T004) before Tests (T005-T019)
- Tests (T005-T019) before Implementation (T020-T031)
- Core Implementation (T020-T031) before Integration (T032-T037)
- Integration (T032-T037) before Polish (T038-T042)

### Critical Path
- T020 (SDR representation) blocks T021, T022, T023
- T023 (interfaces) blocks T024, T025, T026
- T026 (registry) blocks T031 (public API)
- All encoders (T027-T030) block T032 (pipeline)
- T032 (pipeline) blocks T035-T037 (performance validation)
- T031A (batch API) enables T032A (batch implementation)
- T031B (composition utilities) enables T032B (composition framework)
- T030A (multi-SDR spatial) enables T017C (multi-SDR testing)

### Performance Critical Path
- T004 (benchmark framework) enables T009-T012 (performance tests)
- T027-T030 (encoder implementations) must pass T035 (sub-ms benchmark)
- T036 (memory profiling) validates 1MB input handling

## Parallel Execution Examples

### Phase 3.2: Contract Tests (All Parallel)
```bash
# Launch contract tests together:
go test -v tests/contract/sdr_interface_test.go
go test -v tests/contract/sensor_interface_test.go  
go test -v tests/contract/sensor_config_test.go
go test -v tests/contract/sensor_registry_test.go
```

### Phase 3.3: Core Implementation (SDR Components)
```bash
# Launch SDR implementation in parallel:
# Task: "SDR representation struct in internal/sensors/sdr/representation.go"
# Task: "SDR sparsity management in internal/sensors/sdr/sparsity.go"  
# Task: "SDR similarity calculations in internal/sensors/sdr/similarity.go"
```

### Phase 3.3: Encoder Implementation (All Independent)
```bash
# Launch encoder implementations in parallel:
# Task: "Numeric encoder with sub-ms optimization in internal/sensors/encoders/numeric.go"
# Task: "Categorical encoder with hash collision handling in internal/sensors/encoders/categorical.go"
# Task: "Text encoder with 1MB document support in internal/sensors/encoders/text.go"
# Task: "Spatial encoder with topology preservation in internal/sensors/encoders/spatial.go"
```

## Performance Requirements Integration

### Sub-millisecond Encoding Targets
- All encoder implementations (T027-T030) must pass performance contracts (T009-T012)
- Benchmark validation (T035) enforces <1ms encoding time
- Memory profiling (T036) ensures efficient 1MB input handling

### Single-threaded Operation
- No concurrent access patterns in implementation tasks
- Registry (T026) and pipeline (T032) designed for sequential operation
- Test scenarios (T017) validate sequential multi-sensor processing

### Silent Failure Mode
- Error handling integration (T033) implements empty SDR fallback
- Silent failure validation (T018) ensures proper behavior
- All encoder edge cases (T038-T041) test failure scenarios

## Notes
- **[P] tasks**: Different files/packages, no shared dependencies
- **Performance Critical**: Sub-millisecond requirement drives optimization
- **HTM Compliance**: All implementations must maintain 2-5% sparsity
- **Single-threaded**: Simplified implementation without concurrency complexity  
- **Silent Failure**: Empty SDR return on encoding failures
- **Input Limits**: 1MB maximum per encoding operation
- **TDD Required**: All tests must fail before implementation begins
- **Biological Fidelity**: Maintain HTM theoretical properties throughout