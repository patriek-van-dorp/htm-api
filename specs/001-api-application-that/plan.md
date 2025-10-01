
# Implementation Plan: HTM Neural Network API Core

**Branch**: `001-api-application-that` | **Date**: September 30, 2025 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-api-application-that/spec.md`

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
HTM Neural Network API Core providing a fast, scalable REST API for receiving multi-dimensional array inputs representing spatial-temporal patterns and returning processed HTM neural network computation results. The API must respond within 100ms for acknowledgment, support concurrent requests from multiple sensor instances, and maintain consistent input/output formats for API chaining. Implementation uses Go with gonum for matrix calculations, following Go project structure and best practices to replace Python/numpy-based approaches.

## Technical Context
**Language/Version**: Go 1.21+  
**Primary Dependencies**: gonum (matrix calculations), Gin/Echo (HTTP framework), testify (testing)  
**Storage**: N/A (stateless API for this feature)  
**Testing**: Go standard testing package with testify  
**Target Platform**: Linux server, containerized deployment  
**Project Type**: Single project (API service)  
**Performance Goals**: <100ms API acknowledgment, concurrent request handling  
**Constraints**: Stateless design, multi-dimensional array processing, horizontal scaling support  
**Scale/Scope**: Multiple concurrent sensor instances, API chaining capability

**User-provided context**: API needs to be written in Go using libraries for matrix calculations like gonum to replace numpy (which is commonly used in Python). HTM will be implemented in a later feature. Give me the initial API structure according Go project structure and best practices and conventions.

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**I. Research-Driven Development**: ✅ PASS  
- Feature spec includes documented research and validation requirements
- HTM theory foundation will be researched and documented in research.md

**II. Biologically-Inspired Architecture**: ✅ PASS  
- API designed as foundation for future HTM implementation
- Input/output format preserves multi-dimensional spatial-temporal patterns

**III. Test-Driven Scientific Development**: ✅ PASS  
- Go testing with testify planned for comprehensive test coverage
- Contract tests and integration tests will be implemented before API logic

**IV. Scalable Cloud Architecture for AI Research**: ✅ PASS  
- Designed for horizontal scaling with multiple instances
- Stateless API suitable for containerized Azure deployment

**V. Performance & Scientific Rigor by Design**: ✅ PASS  
- 100ms response time requirement documented
- gonum library provides computational efficiency for matrix operations

**Azure & Microsoft Technology Standards**: ⚠️ PARTIAL  
- Go chosen instead of .NET ecosystem (justified by user requirements)
- Azure deployment compatible through containerization
- Future Azure AI/ML service integration planned

**Research & Development Requirements**: ✅ PASS  
- Scientific documentation planned for HTM integration
- API design supports future algorithmic implementations

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
# Go Project Structure (Single API service)
cmd/
├── api/                 # Main application entry point
│   └── main.go

internal/
├── api/                 # HTTP handlers and routing
│   ├── handlers/
│   ├── middleware/
│   └── router.go
├── domain/              # Business logic and entities
│   ├── htm/            # HTM-related domain models
│   └── processing/     # Core processing logic
├── infrastructure/      # External concerns
│   ├── config/
│   └── validation/
└── ports/               # Interface definitions
    └── http.go

pkg/                     # Public APIs (if needed for future features)
└── client/              # Go client library

tests/
├── contract/            # Contract/API tests
├── integration/         # Integration tests
└── unit/               # Unit tests

# Standard Go project files
go.mod
go.sum
Dockerfile
.gitignore
README.md
```

**Structure Decision**: Single Go project using hexagonal architecture pattern with clear separation of concerns. The `internal/` directory ensures internal packages cannot be imported by external projects, while `cmd/api/` provides the application entry point. This structure supports future HTM algorithm integration while maintaining clean boundaries between HTTP concerns, business logic, and infrastructure.

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
- Load `.specify/templates/tasks-template.md` as base template
- Generate Go-specific tasks following TDD and hexagonal architecture patterns
- Create tasks from Phase 1 design artifacts (OpenAPI contracts, data models, quickstart validation)
- Each API endpoint → contract test task [P] + handler implementation task
- Each data model entity → struct definition task [P] + validation task [P]
- Each user story from spec → integration test scenario task
- Infrastructure tasks for middleware, configuration, and error handling

**Go-Specific Task Categories**:
1. **Foundation Tasks**: Go module setup, dependency management, project structure
2. **Data Model Tasks**: Struct definitions, validation tags, JSON marshaling
3. **Contract Tasks**: OpenAPI validation, HTTP handler signatures, middleware setup
4. **Core Logic Tasks**: Matrix processing with gonum, business logic implementation
5. **Testing Tasks**: Unit tests, integration tests, contract validation tests
6. **Infrastructure Tasks**: Configuration, logging, health checks, metrics

**Ordering Strategy**:
- **TDD Principle**: All test tasks before corresponding implementation tasks
- **Dependency Order**: 
  1. Module setup and basic structure
  2. Data models and validation (foundation for everything)
  3. Contract tests (define API behavior)
  4. HTTP infrastructure (routing, middleware)
  5. Core processing logic (HTM input processing)
  6. Integration and performance validation
- **Parallel Execution Markers [P]**: Independent tasks that can run simultaneously
  - Data model struct definitions
  - Individual handler implementations
  - Unit test files for different packages
  - Contract test files for different endpoints

**Key Implementation Priorities**:
1. **High Priority**: HTTP routing, input validation, basic matrix operations
2. **Medium Priority**: Error handling, retry logic, metrics collection
3. **Low Priority**: Performance optimization, advanced monitoring

**Estimated Task Breakdown**:
- **Setup & Structure**: 5-8 tasks (module init, directories, basic config)
- **Data Models**: 8-10 tasks (structs, validation, marshaling, tests)
- **API Layer**: 12-15 tasks (handlers, middleware, routing, contract tests)
- **Processing Logic**: 8-10 tasks (matrix operations, business logic, unit tests)
- **Integration**: 6-8 tasks (end-to-end tests, quickstart validation)
- **Infrastructure**: 5-7 tasks (health checks, metrics, error handling)

**Total Estimated Tasks**: 45-58 numbered, ordered tasks in tasks.md

**Testing Strategy Integration**:
- Each contract endpoint → failing contract test task + implementation task
- Each data model → unit test task + struct implementation task
- Each user story → integration test scenario + implementation tasks
- Performance validation → benchmark test + optimization tasks

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
- [ ] Complexity deviations documented

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
