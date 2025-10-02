# Tasks: Complete Spatial Pooler Engine Integration

**Input**: Design documents from `/specs/004-complete-the-full/`
**Prerequisites**: plan.md (✓), research.md (✓), data-model.md (✓), contracts/ (✓), quickstart.md (✓)

## Execution Flow (main)
```
1. Load plan.md from feature directory ✓
   → Tech stack: Go 1.23+, gonum, Gin, testify
   → Structure: Single project with clean architecture
2. Load design documents: ✓
   → data-model.md: IntegrationContext, ApplicationConfig, ServerConfig, PerformanceConfig, EndpointHandler
   → contracts/: Integration endpoints (health, status, process, config, validation)
   → research.md: TBP insights, production integration, performance optimization
3. Generate tasks by category:
   → Setup: Go module verification, dependency wiring
   → Tests: Contract tests for integration endpoints, TDD test enablement
   → Core: Main.go integration, service wiring, actual implementation
   → Integration: Health checks, metrics, performance validation
   → Polish: Performance tests, HTM validation, documentation
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Main.go integration = sequential (shared file)
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Focus: Complete integration enabling existing test suite
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Repository root**: `cmd/`, `internal/`, `tests/`
- **Main integration**: `cmd/api/main.go`
- **Services**: `internal/services/`
- **Tests**: `tests/contract/`, `tests/integration/`, `tests/unit/`

## Phase 3.1: Setup & Verification
- [x] T001 Verify Go 1.23+ and dependencies (gonum, gin, testify) in go.mod
- [x] T002 [P] Verify existing spatial pooler foundation in `internal/cortical/spatial/`
- [x] T003 [P] Create integration configuration structures in `internal/infrastructure/config/integration_config.go`

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [x] T004 [P] Contract test GET /api/v1/health with spatial pooler status in `tests/contract/health_integration_test.go`
- [x] T005 [P] Contract test GET /api/v1/spatial-pooler/status in `tests/contract/spatial_pooler_status_test.go`
- [x] T006 [P] Contract test POST /api/v1/spatial-pooler/process with actual engine in `tests/contract/spatial_pooler_process_integration_test.go`
- [x] T007 [P] Contract test PUT /api/v1/spatial-pooler/config in `tests/contract/spatial_pooler_config_integration_test.go`
- [x] T008 [P] Contract test GET /api/v1/spatial-pooler/validation/htm-properties in `tests/contract/htm_validation_test.go`
- [x] T009 [P] Integration test complete pipeline from HTTP to actual spatial pooler in `tests/integration/complete_pipeline_test.go`
- [x] T010 [P] Integration test performance requirements (<100ms) in `tests/integration/performance_integration_test.go`
- [x] T011 [P] Integration test concurrent requests (100 max) in `tests/integration/concurrency_integration_test.go`
- [x] T012 [P] Integration test HTM properties validation with actual implementation in `tests/integration/htm_properties_integration_test.go`

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Main Integration
- [x] T013 Create application configuration management in `cmd/api/main.go` (dependency injection setup)
- [x] T014 Wire actual spatial pooler service implementation in `cmd/api/main.go` (replace mocks)
- [x] T015 Configure HTTP router with actual handlers in `cmd/api/main.go` (production routing)
- [x] T016 Initialize metrics collection and health checks in `cmd/api/main.go` (operational readiness)

### Service Layer Integration
- [ ] T017 [P] Create IntegrationContext manager in `internal/services/integration_context.go`
- [ ] T018 [P] Implement ApplicationConfig with validation in `internal/infrastructure/config/application_config.go`
- [ ] T019 [P] Implement ServerConfig for production deployment in `internal/infrastructure/config/server_config.go`
- [ ] T020 [P] Implement PerformanceConfig for optimization in `internal/infrastructure/config/performance_config.go`

### Handler Layer Integration
- [x] T021 Update health handler to include spatial pooler status in `internal/handlers/health_metrics_handler.go`
- [x] T022 Create spatial pooler status handler in `internal/handlers/spatial_pooler_status_handler.go`
- [x] T023 Update spatial pooler processing handler with actual implementation in `internal/handlers/spatial_pooler_handler.go`
- [x] T024 [P] Create spatial pooler configuration handler in `internal/handlers/spatial_pooler_config_handler.go`
- [x] T025 [P] Create HTM validation handler in `internal/handlers/htm_validation_handler.go`

### Service Implementation
- [x] T026 Update spatial pooling service to use actual spatial pooler engine in `internal/services/spatial_pooling_service.go`
- [x] T027 [P] Implement metrics collection service in `internal/services/metrics_service.go`
- [x] T028 [P] Implement health check service in `internal/services/health_service.go`

## Phase 3.4: Integration & Routing
- [x] T029 Update API router with integration endpoints in `internal/api/router.go`
- [x] T030 Configure middleware for performance monitoring in `internal/api/router.go`
- [x] T031 Add request validation middleware in `internal/api/router.go`
- [x] T032 Configure error handling middleware in `internal/api/router.go`

## Phase 3.5: Production Optimization
- [x] T033 [P] Implement matrix pooling for performance in `internal/services/matrix_pool.go`
- [x] T034 [P] Add garbage collection optimization in `internal/infrastructure/config/gc_optimization.go`
- [x] T035 [P] Implement concurrent request limiting in `internal/infrastructure/middleware/concurrency_limiter.go`

## Phase 3.6: HTM Validation & Testing Integration
- [x] T036 Enable existing TDD test suite against actual implementation by removing mocks (specifically replace mocks in `internal/services/spatial_pooling_service.go`, `internal/handlers/spatial_pooler_handler.go`, and wire actual implementations in `cmd/api/main.go`)
- [x] T037 [P] Validate HTM sparsity properties (2-5%) with actual spatial pooler engine in `tests/unit/htm_sparsity_test.go`
- [x] T038 [P] Validate HTM overlap properties with actual implementation in `tests/unit/htm_overlap_test.go`
- [x] T039 [P] Validate HTM learning properties with actual spatial pooler engine in `tests/unit/htm_learning_test.go`
- [x] T040 [P] Validate deterministic behavior with identical inputs and configuration in `tests/unit/deterministic_behavior_test.go`

## Phase 3.7: Performance & Operational Validation
- [ ] T041 [P] Performance test <100ms response time requirement in `tests/performance/response_time_test.go`
- [ ] T042 [P] Performance test 100 concurrent requests limit in `tests/performance/concurrency_test.go`
- [ ] T043 [P] Performance test 10MB input dataset limit in `tests/performance/dataset_size_test.go`
- [ ] T044 [P] Memory usage validation and optimization testing in `tests/performance/memory_test.go`
- [ ] T045 [P] Memory pressure testing with 10MB datasets under concurrent load in `tests/performance/memory_pressure_test.go`
- [ ] T046 [P] Validate sparsity threshold handling (<0.5% sparse, >10% dense) in `tests/unit/sparsity_threshold_test.go`

## Phase 3.8: Polish & Documentation
- [ ] T047 [P] Update API documentation with integration endpoints
- [ ] T048 [P] Create production deployment guide in `docs/deployment.md`
- [ ] T049 [P] Validate error handling and logging throughout pipeline
- [ ] T050 Run complete integration validation using quickstart guide
- [ ] T051 [P] Performance benchmarking and optimization validation

## Dependencies
- **Setup** (T001-T003) before everything
- **Tests** (T004-T012) before implementation (T013-T048)
- **Main Integration** (T013-T016) sequential (shared main.go file)
- **Service Layer** (T017-T020) can be parallel
- **Handler Layer** (T021-T025) depends on T017-T020
- **Service Implementation** (T026-T028) depends on handlers
- **Router Integration** (T029-T032) depends on handlers and services
- **Production Optimization** (T033-T035) can be parallel after core implementation
- **HTM Validation** (T036-T040) depends on actual implementation (T026)
- **Performance Validation** (T041-T046) depends on complete integration
- **Polish** (T047-T051) depends on complete implementation

## Parallel Execution Examples

### Phase 3.2 - Contract Tests (All Parallel)
```bash
# Launch all contract tests together:
go test -run TestHealthIntegration tests/contract/health_integration_test.go
go test -run TestSpatialPoolerStatus tests/contract/spatial_pooler_status_test.go  
go test -run TestSpatialPoolerProcess tests/contract/spatial_pooler_process_integration_test.go
go test -run TestSpatialPoolerConfig tests/contract/spatial_pooler_config_integration_test.go
go test -run TestHTMValidation tests/contract/htm_validation_test.go
```

### Phase 3.3 - Service Layer (Parallel after T016)
```bash
# Launch service implementations together:
Task: "Create IntegrationContext manager in internal/services/integration_context.go"
Task: "Implement ApplicationConfig in internal/infrastructure/config/application_config.go"  
Task: "Implement ServerConfig in internal/infrastructure/config/server_config.go"
Task: "Implement PerformanceConfig in internal/infrastructure/config/performance_config.go"
```

### Phase 3.7 - Performance Tests (All Parallel)
```bash
# Launch performance validation together:
go test -run TestResponseTime tests/performance/response_time_test.go
go test -run TestConcurrency tests/performance/concurrency_test.go
go test -run TestDatasetSize tests/performance/dataset_size_test.go
go test -run TestMemoryUsage tests/performance/memory_test.go
go test -run TestMemoryPressure tests/performance/memory_pressure_test.go
go test -run TestSparsityThreshold tests/unit/sparsity_threshold_test.go
```

## Critical Success Factors
1. **Actual Implementation**: All tasks must use actual spatial pooler engine, not mocks
2. **TDD Compliance**: All tests (T004-T012) must fail before implementation begins
3. **HTM Validation**: Spatial pooler engine must maintain biological properties (sparsity, overlap)
4. **Performance Requirements**: <100ms response time, 100 concurrent requests, 10MB datasets
5. **Test Suite Integration**: Existing TDD test suite must pass against actual implementation
6. **Production Readiness**: Complete health checks, metrics, error handling, and monitoring

## Integration Validation Checklist
- [ ] All existing tests pass with actual spatial pooler engine implementation
- [ ] Health endpoint reports spatial pooler engine status correctly
- [ ] Performance requirements met under load testing
- [ ] HTM properties (sparsity 2-5%, overlap patterns) validated
- [ ] Error handling works with actual implementation
- [ ] Metrics collection captures real performance data
- [ ] Configuration updates work without restart
- [ ] Concurrent requests handled safely without data corruption
- [ ] Deterministic behavior verified with identical inputs
- [ ] Sparsity thresholds properly handled (<0.5% sparse, >10% dense)
- [ ] Memory pressure tested under 10MB dataset loads

## Notes
- **[P] tasks** = different files, no dependencies, can run in parallel
- **Sequential tasks** = same file or dependencies, must run in order
- **Main.go integration** = sequential due to shared file
- **Test-first approach** = all contract and integration tests must fail before implementation
- **HTM compliance** = all implementations must maintain biological and mathematical properties
- **Performance focus** = all tasks must consider <100ms requirement and optimization needs
- **Error handling** = standardized error codes (INVALID_INPUT_DATA, INVALID_CONFIGURATION, MEMORY_LIMIT_EXCEEDED, SPARSITY_OUT_OF_RANGE)
- **Sparsity validation** = handle patterns outside normal range (sparse <0.5%, dense >10%)