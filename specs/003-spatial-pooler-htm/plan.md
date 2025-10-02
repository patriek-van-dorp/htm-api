
# Implementation Plan: Spatial Pooler (HTM Theory) Component

**Branch**: `003-spatial-pooler-htm` | **Date**: October 1, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-spatial-pooler-htm/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from file system structure or context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Implement a spatial pooler component as the first processing layer in the HTM cortical column that transforms encoder outputs (dense/sparse bit arrays) into proper sparse distributed representations (SDRs) with 2-5% sparsity. The spatial pooler will be implemented in the API layer as a cortical algorithm that processes raw encoder outputs from existing sensor encoders, producing true SDRs with semantic continuity for future temporal memory components. This maintains clear separation: sensors handle raw input encoding to bit arrays, while the cortical column (starting with spatial pooler) produces SDRs and handles HTM processing algorithms.

## Technical Context
**Language/Version**: Go 1.23+ (matching existing HTM API project)  
**Primary Dependencies**: gonum (matrix calculations), Gin framework (HTTP), testify (testing), go-playground/validator (validation)  
**Storage**: In-memory processing, no persistent storage required for spatial pooling operations  
**Testing**: go test with testify for unit/integration tests, benchmark tests for performance validation  
**Target Platform**: Linux/Windows server deployment via Azure  
**Project Type**: Single project - extending existing HTM API  
**Performance Goals**: <10ms spatial pooling processing time, 1,000-5,000 requests/second throughput  
**Constraints**: 2-5% SDR sparsity levels, <100ms total API response time, semantic continuity preservation  
**Scale/Scope**: Integration with existing sensor package, preparation for future temporal memory component

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Research-Driven Development**: ✅ PASS - Feature requires research into HTM spatial pooler algorithms and neuroscience principles
**Biologically-Inspired Architecture**: ✅ PASS - Implementation follows HTM's spatial pooling and sparse distributed representations principles  
**Test-Driven Scientific Development**: ✅ PASS - TDD approach planned with scientific validation against HTM properties and performance benchmarks
**Scalable Cloud Architecture**: ✅ PASS - Extends existing Azure-ready HTM API with Go performance optimization
**Performance & Scientific Rigor**: ✅ PASS - Performance requirements (<10ms) and biological plausibility (sparsity levels) defined
**Technology Stack Compliance**: ✅ PASS - Uses Go 1.23+ (acceptable for high-performance libraries per constitution) with existing HTM API stack

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
cmd/
├── api/                     # Existing API entry point
│   └── main.go             # No changes required

internal/
├── api/                     # Existing API router
│   └── router.go           # Add spatial pooler endpoints
├── domain/
│   └── htm/                # Existing HTM domain types
│       ├── spatial_pooler.go    # New: Spatial pooler domain types
│       └── pooling_config.go    # New: Pooling configuration types
├── handlers/               # Existing HTTP handlers
│   └── spatial_pooler_handler.go  # New: Spatial pooler HTTP endpoint
├── infrastructure/         # Existing infrastructure
│   ├── config/
│   └── validation/
├── ports/                  # Existing port interfaces
│   └── spatial_pooling.go  # New: Spatial pooling port interface
├── cortical/               # New: Cortical column components (HTM algorithms)
│   ├── spatial/            # New: Spatial pooling algorithms
│   │   ├── pooler.go       # New: Core spatial pooler implementation
│   │   ├── parameters.go   # New: Pooling parameters and configuration
│   │   ├── learning.go     # New: Learning algorithms for adaptation
│   │   └── inhibition.go   # New: Competitive inhibition algorithms
│   └── sdr/                # New: SDR representation and operations (migrated from sensors)
│       ├── representation.go  # SDR data structure and basic operations
│       ├── similarity.go      # SDR similarity and overlap calculations
│       └── operations.go      # SDR manipulation for cortical algorithms
├── sensors/                # Existing sensor package - NO CHANGES
│   ├── encoders/           # Existing encoders (categorical, numeric, etc.)
│   │                       # Output: consistent byte arrays (dense/sparse bit patterns)
│   ├── interfaces.go       # Existing sensor interfaces
│   ├── config.go          # Existing sensor configuration
│   └── validation.go      # Existing sensor validation
│                           # NOTE: SDR package to be migrated to cortical package
└── services/               # Existing services
    └── spatial_pooling_service.go  # New: Spatial pooling business logic

tests/
├── contract/               # Existing contract tests
│   └── spatial_pooler_test.go    # New: Spatial pooler API contract tests
├── integration/            # Existing integration tests
│   └── spatial_pooling_test.go   # New: End-to-end spatial pooling tests
└── unit/                   # Existing unit tests
    └── cortical/           # New: Cortical algorithm unit tests
        └── spatial/        # New: Spatial pooler unit tests
            ├── pooler_test.go
            ├── parameters_test.go
            ├── learning_test.go
            └── inhibition_test.go
```

**Structure Decision**: Single project extension with clear architectural separation:
- **Sensors**: Handle only raw input encoding to consistent byte arrays (existing functionality, SDR package to be migrated out)
- **Cortical**: New package for HTM cortical column algorithms - contains SDR representation and spatial pooler that produces true SDRs
- **API**: Orchestrates the pipeline: sensor encoding → spatial pooling → SDR output → future temporal memory
- **Services**: Business logic to coordinate between sensors and cortical components
- **Migration**: Move SDR functionality from sensors to cortical package for proper HTM architecture

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/powershell/update-agent-context.ps1 -AgentType copilot`
     **IMPORTANT**: Execute it exactly as specified above. Do not add or remove any arguments.
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base template for task structure
- Generate tasks from Phase 1 design documents (data-model.md, contracts/, quickstart.md)
- Each API endpoint → contract test task [P] + handler implementation task
- Each domain entity → model creation task [P] + validation task [P]
- Each user story from quickstart.md → integration test task + implementation task
- Core spatial pooler algorithm → algorithm implementation tasks in new cortical package

**Ordering Strategy**:
- TDD order: Tests before implementation for all components
- Dependency order: 
  1. Domain models (SpatialPooler, Config, Input/Output types in `internal/domain/htm`)
  2. Cortical algorithm implementation (new `internal/cortical/spatial` package)
  3. Service layer integration (spatial_pooling_service.go coordinates sensors + cortical)
  4. HTTP handlers (spatial_pooler_handler.go)
  5. API router integration (router.go updates)
  6. Integration tests (sensor → spatial pooler → response pipeline)
- Mark [P] for parallel execution (independent files/packages)
- **Architectural Validation**: Ensure sensors remain unchanged, cortical algorithms isolated

**HTM-Specific Task Categories**:
- **Cortical Algorithm Tasks**: Core spatial pooler in new cortical package
- **SDR Migration Tasks**: Move SDR functionality from sensors to cortical package
- **Architectural Separation Tasks**: Ensure clean boundaries between sensors (encoding) and cortical (SDR + HTM processing)
- **Integration Tasks**: Update API to use cortical SDRs instead of sensor SDRs
- **Biological Validation Tasks**: Property tests for sparsity, semantic similarity, learning
- **Performance Tasks**: Latency benchmarks, memory usage validation, throughput testing
- **Scientific Validation Tasks**: HTM property verification, adaptation behavior testing

**Estimated Output**: 40-45 numbered, ordered tasks in tasks.md covering:
- 6-8 SDR migration tasks (move from sensors to cortical package)
- 8-10 domain model tasks (types, validation, configuration)
- 12-15 cortical algorithm tasks (new package: spatial pooler, learning, inhibition)
- 6-8 service/handler tasks (business logic coordination, HTTP endpoints)
- 4-5 integration tasks (API router, sensor-cortical pipeline)
- 8-10 testing tasks (unit, integration, performance, biological validation)
- **Architecture Validation**: Tasks to ensure proper HTM separation (sensors=encoding, cortical=SDRs+algorithms)

**Task Template Structure**:
```
## Task [N]: [Component] - [Action]
**Priority**: [High/Medium/Low]
**Type**: [Implementation/Test/Documentation/Research]
**Parallel**: [Yes/No - can be done independently]
**Dependencies**: [List of prerequisite task numbers]
**Validation**: [How to verify completion]
**HTM Requirements**: [Specific HTM theory compliance requirements]
```

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
