
# Implementation Plan: Generic HTM Sensor Package

**Branch**: `002-add-a-generic` | **Date**: October 1, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-add-a-generic/spec.md`

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
Generic sensor package that converts any input data type into Sparsely Distributed Representations (SDRs) for HTM algorithms. Provides extensible architecture with built-in encoders for common data types (numeric, categorical, text, spatial) and interface for custom sensor implementations. Based on Hierarchical Temporal Memory theory and neurobiological principles from cortical columns.

**Clarifications Applied**: Sub-millisecond performance targets, serializable Go type inputs up to 1MB, single-threaded operation, silent failure error handling.

## Technical Context
**Language/Version**: Go 1.21+ (matching existing HTM API project)  
**Primary Dependencies**: gonum (matrix calculations), Go standard library for interfaces  
**Storage**: N/A (stateless encoding operations)  
**Testing**: testify (consistent with existing project), property-based testing for SDR validation  
**Target Platform**: Cross-platform Go library (Linux, Windows, macOS)
**Project Type**: Single library package (sensor processing component)  
**Performance Goals**: Sub-millisecond encoding (<1ms per operation) for high-frequency systems  
**Constraints**: HTM theory compliance (2-5% sparsity), neurobiological fidelity, single-threaded operation, input data <1MB, silent failure mode  
**Scale/Scope**: Support for extensible sensor types, configurable SDR dimensions up to 10,000 bits, batch processing capability, sensor composition for multi-SDR inputs, multi-SDR collection support for spatial subdivision

**User Requirements**: Go implementation with maintainable, self-documenting code practices and comprehensive documentation for developer learning.

**Clarifications Applied**: 
- Performance: Sub-millisecond latency requirement for high-frequency trading/control systems
- Input constraints: Any serializable Go type up to 1MB per operation
- Concurrency: Single-threaded operation only (no concurrent access support)
- Error handling: Silent failure mode returning empty/default SDR on encoding failures
- Data volume: Medium data support (documents, small images <1MB)

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**I. Research-Driven Development**: ✅ PASS
- Feature requires research into HTM theory, neuroscience principles, and existing implementations
- Research documented with clear hypotheses and theoretical foundations
- All decisions validated against current HTM literature and theory

**II. Biologically-Inspired Architecture**: ✅ PASS  
- Implementation adheres to HTM core principles: sparse distributed representations, biological encoding
- Architecture mirrors cortical column processing with clear separation of concerns
- Maintains biological fidelity in SDR operations and sparsity requirements

**III. Test-Driven Scientific Development**: ✅ PASS (with clarifications)
- Mandatory TDD approach with hypothesis → test → implementation → validation
- Comprehensive unit tests, integration tests, and behavioral validation against HTM properties
- Performance benchmarks for sub-millisecond latency and biological plausibility
- **Clarification applied**: Silent failure mode aligns with scientific validation approach

**IV. Scalable Cloud Architecture for AI Research**: ✅ PASS (with clarifications)
- Go library design supports Azure deployment and distributed processing
- Architecture enables integration with Azure AI/ML services for large-scale experiments
- **Clarification applied**: Single-threaded operation simplifies cloud deployment patterns

**V. Performance & Scientific Rigor by Design**: ✅ PASS (with clarifications)
- Performance optimization focus with computational complexity analysis
- Scientific accuracy validation against HTM theory and neuroscience literature
- **Clarification applied**: Sub-millisecond performance targets quantified
- **Clarification applied**: Input size constraints (1MB) support performance goals

**Result**: PASS - All constitutional principles satisfied with clarifications integrated

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
internal/
├── sensors/
│   ├── interfaces.go          # Core sensor contracts and SDR types
│   ├── registry.go           # Sensor registration and discovery
│   ├── validation.go         # SDR quality validation
│   ├── encoders/
│   │   ├── numeric.go        # Numeric data encoder (sub-ms performance)
│   │   ├── categorical.go    # Categorical data encoder
│   │   ├── text.go          # Text tokenization and encoding (<1MB support)
│   │   └── spatial.go       # Spatial/image data encoder
│   └── sdr/
│       ├── representation.go # SDR data structure and operations
│       ├── sparsity.go      # Sparsity management and validation
│       └── similarity.go    # Semantic similarity calculations

pkg/
├── sensors/
│   ├── public_api.go         # Public API exports for external use
│   ├── batch_api.go          # Batch processing API for multiple inputs
│   └── composition.go        # Sensor composition utilities for multi-SDR processing

tests/
├── contract/
│   ├── sensor_interface_test.go    # Single-threaded operation validation
│   ├── sdr_validation_test.go      # Silent failure mode testing
│   └── encoder_contract_test.go    # Sub-ms performance validation
├── integration/
│   ├── multi_sensor_test.go        # Sequential sensor processing
│   ├── batch_processing_test.go    # 1MB input volume testing
│   ├── sensor_composition_test.go  # Multi-SDR sensor composition testing
│   ├── multi_sdr_test.go          # Multi-SDR collection testing (spatial subdivision)
│   └── performance_test.go         # Sub-millisecond latency validation
└── unit/
    ├── numeric_encoder_test.go     # Precision and performance tests
    ├── categorical_encoder_test.go
    ├── text_encoder_test.go        # 1MB document processing tests
    ├── spatial_encoder_test.go
    ├── batch_test.go               # Batch processing edge cases
    ├── composition_test.go         # Sensor composition edge cases
    ├── multi_sdr_test.go          # Multi-SDR collection edge cases
    └── sdr_operations_test.go      # Silent failure and sparsity tests
```

**Structure Decision**: Single Go library structure selected with clarifications applied:
- Performance-focused encoder implementations for sub-millisecond operation
- Single-threaded design eliminates concurrency complexity  
- Silent failure mode integrated into validation and testing layers
- Input size validation (<1MB) built into encoder interfaces
- Test structure validates all clarified requirements

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

## Phase 2: Task Planning Approach (Updated with Clarifications)
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy** (Incorporating Clarifications):
Based on the Phase 1 design artifacts and clarified requirements, the /tasks command will generate approximately 42 implementation tasks:

1. **Performance-Critical Contract Tests**: Each interface generates performance validation tests
   - SDR interface compliance with sub-millisecond operations [P]
   - SensorInterface performance benchmarks (<1ms encoding) [P] 
   - Input size validation tests (1MB constraint) [P]
   - Silent failure mode validation tests [P]

2. **Core Entity Implementation** (Single-threaded focus):
   - SDR representation with performance optimization [P]
   - SensorConfig with input size validation [P]
   - SensorRegistry (sequential operation only) 
   - Error types with silent failure support [P]

3. **High-Performance Encoder Implementation**:
   - NumericEncoder with sub-ms optimization [P]
   - CategoricalEncoder with size constraints [P]
   - TextEncoder with 1MB document support [P]
   - SpatialEncoder with memory optimization [P]
   - Multi-SDR SpatialEncoder for image subdivision [P]

4. **Enhanced Integration Tasks**:
   - Sequential multi-sensor pipeline (no concurrency)
   - Batch processing API and implementation
   - Sensor composition framework for multi-SDR processing
   - Multi-SDR collection handling and validation
   - Silent failure integration testing
   - 1MB input volume stress testing
   - Sub-millisecond performance validation

**Ordering Strategy** (Updated):
1. **Foundation Layer**: Performance-critical interfaces and SDR types
2. **Validation Layer**: Silent failure and size constraint validation
3. **Encoder Layer**: High-performance individual sensor implementations (parallel)
4. **Integration Layer**: Sequential processing and performance validation
5. **Documentation Layer**: Performance characteristics and constraint documentation

**Performance Validation Tasks**:
- Sub-millisecond latency benchmarks for each encoder
- Memory usage profiling for 1MB inputs
- Silent failure behavior validation
- Single-threaded operation verification

**Estimated Output**: 
- 8 performance-critical contract tasks (T005-T012)
- 4 setup and foundation tasks (T001-T004)
- 7 integration test scenarios (T013-T019)
- 12 core implementation tasks (T020-T031B)
- 7 integration and validation tasks (T032-T037)
- 4 unit test and polish tasks (T038-T042)
**Total**: ~42 ordered, performance-focused tasks in tasks.md

**IMPORTANT**: The above planning incorporates all clarifications and will be executed by the /tasks command

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


## Progress Tracking (Updated with Clarifications)
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command) → research.md created
- [x] Phase 1: Design complete (/plan command) → data-model.md, contracts/, quickstart.md, .github/copilot-instructions.md updated
- [x] Phase 2: Task planning approach updated with clarifications (/plan command)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS (all HTM principles satisfied)
- [x] Post-Design Constitution Check: PASS (design maintains constitutional compliance with clarifications)
- [x] All NEEDS CLARIFICATION resolved (5 critical clarifications applied)
- [x] Complexity deviations documented (none required - follows standard patterns)

**Clarifications Integration Status**:
- [x] Performance requirements integrated: Sub-millisecond encoding targets
- [x] Input constraints clarified: Serializable Go types up to 1MB
- [x] Concurrency model specified: Single-threaded operation only
- [x] Error handling strategy defined: Silent failure with empty SDR fallback
- [x] Data volume limits established: Medium data support (<1MB)

**Artifacts Generated/Updated**:
- [x] research.md - HTM theory research and Go implementation decisions
- [x] data-model.md - Core entities and relationships for sensor package
- [x] contracts/interfaces.md - Go interface definitions and contracts
- [x] quickstart.md - Comprehensive usage guide with examples
- [x] .github/copilot-instructions.md - Updated with new sensor package context
- [x] plan.md - Updated with all clarifications integrated

**Clarifications Applied Successfully**: All 5 critical ambiguities resolved and integrated into design artifacts.

**Ready for /tasks command** - All planning artifacts complete, clarifications integrated, and performance requirements specified.

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
