# Tasks: Spatial Pooler (HTM Theory) Component

**Input**: Design documents from `/specs/003-spatial-pooler-htm/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Extract: Go 1.23+, gonum, Gin framework, testify
   → Structure: Single project extension with cortical package
2. Load design documents:
   → data-model.md: Extract entities → SpatialPooler, SpatialPoolerConfig, PoolingInput/Result
   → contracts/: API endpoints → POST /process, GET/PUT /config, GET /metrics
   → quickstart.md: Extract test scenarios → deterministic behavior, semantic similarity
3. Generate tasks by category:
   → Setup: project structure, dependencies, linting
   → Tests: contract tests, integration tests (TDD approach)
   → Core: domain models, cortical algorithms, services
   → Integration: handlers, API routes, pipeline
   → Polish: unit tests, performance validation, documentation
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. HTM-specific ordering:
   → Domain types before algorithms
   → Cortical algorithms before services
   → Services before handlers
   → All core before integration
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
Single project extension - paths relative to repository root:
- **Source**: `internal/`, `cmd/`
- **Tests**: `tests/contract/`, `tests/integration/`, `tests/unit/`
- **New packages**: `internal/cortical/`, `internal/domain/htm/`

## Phase 3.1: Setup
- [x] T001 Create cortical package structure: `internal/cortical/spatial/`, `internal/cortical/sdr/`
- [x] T002 Initialize Go module dependencies: gonum for matrix operations, go-playground/validator
- [x] T003 [P] Configure linting and formatting tools for new cortical package

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [x] T004 [P] Contract test POST /api/v1/spatial-pooler/process in `tests/contract/spatial_pooler_process_test.go`
- [x] T005 [P] Contract test GET /api/v1/spatial-pooler/config in `tests/contract/spatial_pooler_config_get_test.go`
- [x] T006 [P] Contract test PUT /api/v1/spatial-pooler/config in `tests/contract/spatial_pooler_config_put_test.go`
- [x] T007 [P] Contract test GET /api/v1/spatial-pooler/metrics in `tests/contract/spatial_pooler_metrics_test.go`
- [x] T008 [P] Integration test deterministic behavior in `tests/integration/spatial_pooler_deterministic_test.go`
- [x] T009 [P] Integration test semantic similarity preservation in `tests/integration/spatial_pooler_semantic_test.go`
- [x] T010 [P] Integration test encoder-to-spatial-pooler pipeline in `tests/integration/spatial_pooler_pipeline_test.go`
- [x] T011 [P] Performance test <10ms processing time in `tests/integration/spatial_pooler_performance_test.go`

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Domain Models
- [x] T012: Create SpatialPoolerMode enum in internal/domain/htm/spatial_pooler.go
- [x] T013 [P] SpatialPoolerConfig struct with validation in `internal/domain/htm/spatial_pooler.go` (consolidated)
- [x] T014 [P] PoolingInput struct with validation in `internal/domain/htm/spatial_pooler.go`
- [x] T015 [P] PoolingResult struct with validation in `internal/domain/htm/spatial_pooler.go`
- [x] T016 [P] PoolingError types and handling in `internal/domain/htm/spatial_pooler.go`
- [x] T017 [P] SpatialPoolerMetrics struct in `internal/domain/htm/spatial_pooler.go`

### Cortical Algorithm Implementation
- [x] T018 [P] SDR operations in `internal/cortical/sdr/operations.go` (migrate from sensors)
- [x] T019 [P] SDR representation in `internal/cortical/sdr/representation.go` (migrate from sensors)
- [x] T020 [P] SDR similarity calculations in `internal/cortical/sdr/similarity.go` (migrate from sensors)
- [x] T021 SpatialPooler core algorithm in `internal/cortical/spatial/pooler.go`
- [x] T022 Spatial pooler parameters management in `internal/cortical/spatial/parameters.go`
- [x] T023 Competitive inhibition algorithms in `internal/cortical/spatial/inhibition.go`
- [x] T024 Learning algorithms and adaptation in `internal/cortical/spatial/learning.go`

### Service Layer
- [x] T025 Spatial pooling port interface in `internal/ports/spatial_pooling.go`
- [x] T026 Spatial pooling service implementation in `internal/services/spatial_pooling_service.go`

### HTTP Layer
- [x] T027 Spatial pooler HTTP handler in `internal/handlers/spatial_pooler_handler.go`
- [x] T028 Update API router in `internal/api/router.go` (add spatial pooler routes)

## Phase 3.4: Integration
- [ ] T029 Wire spatial pooling service with cortical algorithms
- [ ] T030 Integrate spatial pooler handler with existing validation infrastructure
- [ ] T031 Add spatial pooler endpoints to main API router
- [ ] T032 Test end-to-end pipeline: sensor encoding → spatial pooling → SDR response

## Phase 3.5: Polish
- [ ] T033 [P] Unit tests for SpatialPooler algorithm in `tests/unit/cortical/spatial/pooler_test.go`
- [ ] T034 [P] Unit tests for inhibition algorithms in `tests/unit/cortical/spatial/inhibition_test.go`
- [ ] T035 [P] Unit tests for learning algorithms in `tests/unit/cortical/spatial/learning_test.go`
- [ ] T036 [P] Unit tests for spatial pooling service in `tests/unit/services/spatial_pooling_service_test.go`
- [ ] T037 [P] Benchmark tests for <10ms processing requirement
- [ ] T038 [P] HTM biological property validation tests (sparsity 2-5%, semantic similarity)
- [ ] T039 [P] Error handling and edge case tests
- [ ] T040 [P] Update API documentation with spatial pooler endpoints
- [ ] T041 Remove code duplication and optimize performance
- [ ] T042 Run quickstart.md manual testing scenarios
- [ ] T043 [P] Load testing for 1000-5000 req/sec throughput validation (FR-013)
- [ ] T044 [P] Default SDR pattern generation for invalid inputs (FR-012)
- [ ] T045 [P] Input size validation and error handling for oversized inputs (FR-006, FR-014)

## Dependencies
- Tests (T004-T011) before implementation (T012-T028)
- Domain models (T012-T017) before algorithms (T018-T024)
- SDR migration (T018-T020) before spatial pooler core (T021)
- T021 blocks T022, T023, T024 (pooler.go is foundation)
- Algorithms (T018-T024) before service (T025-T026)
- Service (T025-T026) before handler (T027)
- Handler (T027) before router integration (T028)
- Core implementation before integration (T029-T032)
- Implementation before polish (T033-T045)

## Parallel Example
```
# Launch T004-T011 together (contract and integration tests):
Task: "Contract test POST /api/v1/spatial-pooler/process in tests/contract/spatial_pooler_process_test.go"
Task: "Contract test GET /api/v1/spatial-pooler/config in tests/contract/spatial_pooler_config_get_test.go"
Task: "Integration test deterministic behavior in tests/integration/spatial_pooler_deterministic_test.go"
Task: "Performance test <10ms processing time in tests/integration/spatial_pooler_performance_test.go"

# Launch T012-T017 together (domain models in different files):
Task: "SpatialPoolerMode enum in internal/domain/htm/spatial_pooler.go"
Task: "SpatialPoolerConfig struct with validation in internal/domain/htm/pooling_config.go"

# Launch T018-T020 together (SDR migration to cortical package):
Task: "SDR operations in internal/cortical/sdr/operations.go"
Task: "SDR representation in internal/cortical/sdr/representation.go"
Task: "SDR similarity calculations in internal/cortical/sdr/similarity.go"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- HTM-specific validation: 2-5% sparsity, <10ms processing
- Maintain backward compatibility with existing sensor package
- Architectural separation: sensors (encoding) vs cortical (SDR + HTM algorithms)
- Commit after each task
- SDR package migration from sensors to cortical maintains clean HTM architecture
- New tasks T043-T045 address missing coverage gaps: throughput testing, fallback patterns, input validation

## HTM-Specific Validation Requirements
- Sparsity levels must be 2-5% for all SDR outputs
- Processing time must be <10ms per spatial pooling operation
- Semantic similarity preservation: similar inputs must have 30-70% SDR overlap, different inputs <20% overlap
- Deterministic behavior when learning is disabled
- Proper competitive inhibition (global vs local)
- Learning adaptation without breaking existing patterns
- Matrix operations optimized with gonum for performance
- Memory usage scales linearly with column dimensions
- Throughput validation: 1000-5000 requests/second under normal conditions
- Default SDR pattern generation for invalid inputs (all zeros, all ones, corrupted data)
- Input validation: reject inputs exceeding expected dimensions with clear error messages