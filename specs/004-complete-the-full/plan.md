
# Implementation Plan: Complete Spatial Pooler Engine Integration

**Branch**: `004-complete-the-full` | **Date**: October 2, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-complete-the-full/spec.md`

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
Complete the full integration by implementing the actual spatial pooler engine and wiring everything together in main.go to enable TDD tests to run against the actual implementation. This involves creating a production-ready HTM spatial pooler engine that adheres to biological principles while providing high-performance computation for real-time processing.

## Technical Context
**Language/Version**: Go 1.23+ (existing project standard)  
**Primary Dependencies**: gonum (matrix calculations), Gin framework (HTTP), testify (testing), go-playground/validator (validation)  
**Storage**: In-memory processing, no persistent storage required for spatial pooling operations  
**Testing**: testify with comprehensive contract, integration, and unit test suites  
**Target Platform**: Linux server, containerized deployment  
**Project Type**: single - HTM API service with clean architecture  
**Performance Goals**: <100ms response time for real-time interactive use, handle up to 100 concurrent requests  
**Constraints**: <10MB input datasets, deterministic behavior for reproducible results, biologically-inspired architecture  
**Scale/Scope**: Production-ready spatial pooler engine, complete HTM processing pipeline, comprehensive test coverage

**HTM Theory Context**: Must adhere to core HTM principles including sparse distributed representations, spatial pooling with inhibition, synaptic learning, and biological plausibility. Implementation should incorporate insights from Numenta's Thousand Brains Project (tbp.monty) evolution while adapting to Go's performance characteristics.

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. Research-Driven Development ✅
- Feature builds on existing HTM spatial pooler implementation with research foundation
- Must research Thousand Brains Project (tbp.monty) evolution and HTM theory advances
- Implementation approach validated against neuroscience literature

### II. Biologically-Inspired Architecture ✅  
- Maintains HTM core principles: sparse distributed representations, spatial pooling, hierarchical structure
- Existing architecture already mirrors biological systems with clear spatial/temporal separation
- Spatial pooler implementation follows biological inhibition and learning mechanisms

### III. Test-Driven Scientific Development ✅
- Comprehensive TDD test suite already exists (contract, integration, unit tests)
- Feature specifically aims to make existing tests pass against actual implementation
- Scientific validation through behavioral testing against known HTM properties

### IV. Scalable Cloud Architecture for AI Research ✅
- Go implementation provides high-performance computing suitable for Azure deployment
- Architecture supports distributed processing and research collaboration
- Performance targets align with cloud-scale AI workloads

### V. Performance & Scientific Rigor by Design ✅
- Performance optimization (<100ms) and biological accuracy both prioritized
- Implementation includes computational complexity considerations (gonum matrices)
- Benchmarks validate both performance and biological plausibility

### Azure & Microsoft Technology Standards ⚠️
- **JUSTIFIED DEVIATION**: Using Go instead of .NET ecosystem for high-performance HTM computations
- **Rationale**: HTM algorithms require intensive matrix operations where Go+gonum provides demonstrated performance benefits for this specific domain
- **Compliance**: Still targeting Azure deployment, monitoring, and AI service integration

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
├── api/
│   └── main.go              # Application entry point - wire components

internal/
├── api/
│   └── router.go           # HTTP routing configuration
├── cortical/
│   ├── sdr/               # Sparse Distributed Representation operations
│   └── spatial/           # Spatial pooler engine implementation
├── domain/
│   └── htm/               # HTM domain models and interfaces
├── handlers/              # HTTP request handlers
├── infrastructure/        # Configuration and validation
├── ports/                 # Interface definitions
├── sensors/               # Input data encoding
└── services/              # Business logic services

pkg/
├── client/                # API client libraries
├── config/                # Configuration management
└── sensors/               # Public sensor interfaces

tests/
├── contract/              # API contract tests
├── integration/           # End-to-end integration tests
└── unit/                  # Unit tests for components
```

**Structure Decision**: Single project with clean architecture pattern. The existing structure already supports the spatial pooler integration with clear separation between domain logic, infrastructure, and presentation layers. Main integration point is `cmd/api/main.go` for wiring components together.

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
- Load `.specify/templates/tasks-template.md` as base for task structure
- Generate integration tasks from Phase 1 design docs (contracts, data model, quickstart)
- Each API endpoint → integration test task to validate actual implementation [P]
- Each configuration component → service wiring task in main.go [P]
- Each HTM validation requirement → property validation test task
- Production deployment tasks for operational readiness

**Integration-Specific Task Categories**:
1. **Main.go Integration Tasks**: Wire spatial pooler service, configure dependencies, initialize components
2. **Service Implementation Tasks**: Replace mock implementations with actual spatial pooler engine
3. **Test Integration Tasks**: Enable existing TDD test suite against actual implementation
4. **Performance Validation Tasks**: Validate <100ms response time and 100 concurrent request requirements
5. **HTM Property Tasks**: Ensure sparsity, overlap, and learning validation pass with actual implementation
6. **Operational Readiness Tasks**: Health checks, metrics collection, error handling

**Ordering Strategy**:
- **Phase 1**: Foundation tasks (service integration, main.go wiring) 
- **Phase 2**: Test enablement tasks (replace mocks, validate contracts)
- **Phase 3**: Performance and HTM validation tasks
- **Phase 4**: Operational and production readiness tasks
- Mark [P] for parallel execution where services are independent
- Sequential execution for dependent components (service before tests)

**Special Considerations**:
- Integration tasks must ensure actual implementation, not mocks
- Performance tasks must validate against real processing times
- HTM validation tasks must check biological and mathematical properties
- Test tasks must execute existing comprehensive test suite
- All tasks must maintain backward compatibility with existing API

**Estimated Output**: 50+ numbered, ordered tasks focusing on integration completion and test validation

**Task Dependencies**:
- Spatial pooler engine wiring must complete before handler integration
- Service implementation must complete before test enablement
- Configuration tasks must complete before performance validation
- Integration tasks must complete before operational readiness

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
- [x] Complexity deviations documented (Go language usage justified)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
