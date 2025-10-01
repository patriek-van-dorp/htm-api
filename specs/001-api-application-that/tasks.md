# Tasks: HTM Neural Network API Core

**Input**: Design documents from `/specs/001-api-application-that/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, services, CLI commands
   → Integration: DB, middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All endpoints implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: Go project structure at repository root
- Paths follow hexagonal architecture with cmd/, internal/, pkg/, tests/

## Phase 3.1: Setup
- [x] T001 Create Go project structure with cmd/api/, internal/{api,domain,infrastructure,ports}/, pkg/client/, tests/{contract,integration,unit}/
- [x] T002 Initialize Go module with dependencies: gin-gonic/gin, gonum.org/v1/gonum, go-playground/validator/v10, google/uuid
- [x] T003 [P] Configure Go project files: go.mod, go.sum, .gitignore, README.md, Dockerfile
- [x] T004 [P] Setup basic configuration in internal/infrastructure/config/config.go
- [x] T005 [P] Create validation utilities in internal/infrastructure/validation/validator.go

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [x] T006 [P] Contract test POST /api/v1/process in tests/contract/process_post_test.go
- [x] T007 [P] Contract test GET /health in tests/contract/health_get_test.go
- [x] T008 [P] Contract test GET /metrics in tests/contract/metrics_get_test.go
- [x] T009 [P] Integration test HTM input processing workflow in tests/integration/htm_processing_test.go
- [x] T010 [P] Integration test error handling and validation in tests/integration/error_handling_test.go
- [x] T011 [P] Integration test concurrent requests in tests/integration/concurrent_requests_test.go
- [x] T012 [P] Integration test response time requirements in tests/integration/performance_test.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T013 [P] HTMInput model struct in internal/domain/htm/input.go
- [ ] T014 [P] InputMetadata model struct in internal/domain/htm/metadata.go
- [ ] T015 [P] ProcessingResult model struct in internal/domain/htm/result.go
- [ ] T016 [P] ProcessingStatus enum in internal/domain/htm/status.go
- [ ] T017 [P] APIRequest model struct in internal/domain/htm/request.go
- [ ] T018 [P] APIResponse model struct in internal/domain/htm/response.go
- [ ] T019 [P] APIError model struct in internal/domain/htm/error.go
- [ ] T020 [P] RequestPriority enum in internal/domain/htm/priority.go
- [ ] T021 [P] Matrix processing service interface in internal/ports/processing.go
- [ ] T022 [P] HTTP handler interfaces in internal/ports/http.go
- [ ] T023 Matrix processing service implementation in internal/domain/processing/htm_processor.go
- [ ] T024 Input validation service in internal/domain/processing/validator.go
- [ ] T025 Process HTTP handler in internal/api/handlers/process.go
- [ ] T026 Health HTTP handler in internal/api/handlers/health.go
- [ ] T027 Metrics HTTP handler in internal/api/handlers/metrics.go
- [ ] T028 HTTP router setup in internal/api/router.go
- [ ] T029 Main application entry point in cmd/api/main.go

## Phase 3.4: Integration
- [ ] T030 [P] Request logging middleware in internal/api/middleware/logging.go
- [ ] T031 [P] Error handling middleware in internal/api/middleware/error.go
- [ ] T032 [P] Metrics collection middleware in internal/api/middleware/metrics.go
- [ ] T033 [P] CORS middleware in internal/api/middleware/cors.go
- [ ] T034 Connect HTTP handlers to processing services in internal/api/handlers/process.go
- [ ] T035 Integrate validation with HTTP handlers
- [ ] T036 Implement retry logic for transient failures
- [ ] T037 Add request timeout handling

## Phase 3.5: Polish
- [ ] T038 [P] Unit tests for HTM models in tests/unit/htm_models_test.go
- [ ] T039 [P] Unit tests for processing service in tests/unit/processing_service_test.go
- [ ] T040 [P] Unit tests for validation logic in tests/unit/validation_test.go
- [ ] T041 [P] Unit tests for HTTP handlers in tests/unit/handlers_test.go
- [ ] T042 [P] Unit tests for middleware in tests/unit/middleware_test.go
- [ ] T043 Performance benchmarks for matrix operations in tests/unit/processing_benchmark_test.go
- [ ] T044 [P] Go client library in pkg/client/htm_client.go
- [ ] T045 [P] OpenAPI documentation generation
- [ ] T046 Code formatting and linting with gofmt and golangci-lint
- [ ] T047 Execute quickstart.md validation scenarios

## Dependencies
**Setup Phase**: T001 → T002 → T003-T005 (parallel)
**Test Phase**: T006-T012 (all parallel, must complete before implementation)
**Models Phase**: T013-T020 (parallel, foundation for services)
**Interfaces Phase**: T021-T022 (parallel, required for implementations)
**Implementation Phase**: 
  - T023-T024 depend on T013-T022
  - T025-T027 depend on T021-T022
  - T028-T029 depend on T025-T027
**Integration Phase**: 
  - T030-T033 (parallel middleware)
  - T034-T037 depend on T025-T029
**Polish Phase**: T038-T047 (mostly parallel, T047 depends on all implementation)

## Parallel Example
```
# Phase 3.2 - Launch all contract tests together:
Task: "Contract test POST /api/v1/process in tests/contract/process_post_test.go"
Task: "Contract test GET /health in tests/contract/health_get_test.go"
Task: "Contract test GET /metrics in tests/contract/metrics_get_test.go"
Task: "Integration test HTM input processing workflow in tests/integration/htm_processing_test.go"
Task: "Integration test error handling and validation in tests/integration/error_handling_test.go"
Task: "Integration test concurrent requests in tests/integration/concurrent_requests_test.go"
Task: "Integration test response time requirements in tests/integration/performance_test.go"

# Phase 3.3 - Launch all model definitions together:
Task: "HTMInput model struct in internal/domain/htm/input.go"
Task: "InputMetadata model struct in internal/domain/htm/metadata.go"
Task: "ProcessingResult model struct in internal/domain/htm/result.go"
Task: "ProcessingStatus enum in internal/domain/htm/status.go"
Task: "APIRequest model struct in internal/domain/htm/request.go"
Task: "APIResponse model struct in internal/domain/htm/response.go"
Task: "APIError model struct in internal/domain/htm/error.go"
Task: "RequestPriority enum in internal/domain/htm/priority.go"

# Phase 3.4 - Launch all middleware together:
Task: "Request logging middleware in internal/api/middleware/logging.go"
Task: "Error handling middleware in internal/api/middleware/error.go"
Task: "Metrics collection middleware in internal/api/middleware/metrics.go"
Task: "CORS middleware in internal/api/middleware/cors.go"
```

## Notes
- [P] tasks = different files, no dependencies between them
- Verify all tests fail before implementing corresponding features
- Follow TDD: write failing tests, then implement to make them pass
- Commit after each task completion for clear progress tracking
- All tasks specify exact file paths for unambiguous implementation
- Performance requirement: <100ms acknowledgment for /api/v1/process endpoint
- Matrix operations use gonum library for computational efficiency
- Hexagonal architecture maintains clean separation of concerns

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts (openapi.yaml)**:
   - POST /api/v1/process → T006, T025
   - GET /health → T007, T026
   - GET /metrics → T008, T027

2. **From Data Model (data-model.md)**:
   - HTMInput entity → T013
   - InputMetadata entity → T014
   - ProcessingResult entity → T015
   - ProcessingStatus enum → T016
   - APIRequest entity → T017
   - APIResponse entity → T018
   - APIError entity → T019
   - RequestPriority enum → T020

3. **From User Stories (quickstart.md)**:
   - Basic HTM processing workflow → T009
   - Error handling scenarios → T010
   - Concurrent request handling → T011
   - Performance validation → T012
   - Quickstart validation → T047

4. **Ordering**:
   - Setup (T001-T005) → Tests (T006-T012) → Models (T013-T020) → Interfaces (T021-T022) → Implementation (T023-T029) → Integration (T030-T037) → Polish (T038-T047)

## Validation Checklist
*GATE: Checked by main() before returning*

- [x] All contracts have corresponding tests (T006-T008 cover all endpoints)
- [x] All entities have model tasks (T013-T020 cover all data models)
- [x] All tests come before implementation (T006-T012 before T013+)
- [x] Parallel tasks truly independent (different files, no shared state)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Performance requirements addressed (T012, T043)
- [x] Go-specific patterns followed (hexagonal architecture, standard project layout)
- [x] TDD principle enforced (failing tests before implementation)