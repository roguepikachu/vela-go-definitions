# Requirements: defkit API Documentation

**Defined:** 2026-03-11
**Core Value:** Any new developer can open one HTML page and understand how to write a defkit definition from scratch

## v1 Requirements

### Page Structure

- [ ] **STRUCT-01**: Single `docs/index.html` file with all content inline (no external files)
- [ ] **STRUCT-02**: Fixed sidebar navigation with all sections and anchors
- [ ] **STRUCT-03**: Modern design — dark sidebar, light content area, clean typography
- [ ] **STRUCT-04**: Responsive layout that works on desktop
- [ ] **STRUCT-05**: Syntax-highlighted code blocks for Go and CUE

### Definition Builders

- [ ] **DEF-01**: `defkit.NewComponent(name)` documented with description, signature, example
- [ ] **DEF-02**: `defkit.NewTrait(name)` documented with description, signature, example
- [ ] **DEF-03**: `defkit.NewPolicy(name)` documented with description, signature, example
- [ ] **DEF-04**: `defkit.NewWorkflowStep(name)` documented with description, signature, example
- [ ] **DEF-05**: Definition chain methods documented per type (Description, Workload, Params, Template, AppliesTo, etc.)
- [ ] **DEF-06**: `AutodetectWorkload()` documented — for definitions without a fixed workload type
- [ ] **DEF-07**: `RawCUE(string)` documented — inject raw CUE template body (used in k8s-objects, ref-objects)
- [ ] **DEF-08**: `WithImports(...)` documented — add CUE import paths to definition

### Parameter Builders

- [ ] **PARAM-01**: Scalar types documented: String, Int, Bool
- [ ] **PARAM-02**: Collection types documented: Array, List, StringList, StringKeyMap, Object
- [ ] **PARAM-03**: Complex types documented: Enum, Struct, OneOf
- [ ] **PARAM-04**: Parameter chain methods documented: Required, Optional, Default, Description, Values, Min, Max, WithFields, Of

### Template Methods

- [ ] **TPL-01**: `tpl.Output(resource)` documented with example
- [ ] **TPL-02**: `tpl.Outputs(name, resource)` and `tpl.OutputsIf(...)` documented
- [ ] **TPL-03**: `tpl.Patch()` builder and `tpl.SetRawPatchBlock()` documented
- [ ] **TPL-04**: `tpl.UsePatchContainer(config)` documented with PatchContainerConfig fields
- [ ] **TPL-05**: `tpl.SetRawHeaderBlock()` and `tpl.SetRawOutputsBlock()` documented
- [ ] **TPL-06**: `tpl.Helper()` builder documented with `.FromArray()`, `.FromFields()`, `.Pick()`, `.Map()`, `.Dedupe()`, `.AfterOutput()`, `.Build()`

### Resource Builders

- [ ] **RES-01**: `defkit.NewResource(apiVersion, kind)` documented
- [ ] **RES-02**: `.Set()`, `.SetIf()`, `.SetDefault()` documented with CUE equivalent
- [ ] **RES-03**: `.Directive()`, `.SetRawBlock()` documented
- [ ] **RES-04**: `.ForEach()`, `.ForEachWith()`, `.Item()`, `.ItemIf()` documented
- [ ] **RES-05**: Conditional methods `.If()/.EndIf()`, `.SpreadIf()` documented
- [ ] **RES-06**: `NewResourceWithConditionalVersion(kind)` and `.VersionIf(condition, version)` documented
- [ ] **RES-07**: `ItemBuilder` methods documented: `.Set()`, `.SetIf()`, `.IfSet()`, `.IfNotSet()`, `.If()`, `.Let()`, `.Wrap()`
- [ ] **RES-08**: `ItemFieldIsSet(fieldName)` and `FieldEquals(fieldName, value)` documented

### Value Expressions

- [ ] **VAL-01**: `defkit.Lit()`, `defkit.Reference()`, `defkit.Interpolation()` documented
- [ ] **VAL-02**: Value chain methods: `.IsSet()`, `.NotSet()`, `.Eq()`, `.Ne()`, `.Lt()`, `.Le()`, `.Gt()`, `.Ge()`, `.Field()`, `.Or()` documented
- [ ] **VAL-03**: Logical operators: `defkit.And()`, `defkit.Or()`, `defkit.Not()` documented
- [ ] **VAL-04**: String operations: `defkit.Plus()`, `defkit.Format()`, `defkit.StringsToLower()`, `defkit.StrconvFormatInt()` documented
- [ ] **VAL-05**: `defkit.PathExists(path)`, `defkit.ParamRef(param)`, `defkit.ParamNotSet(param)` documented
- [ ] **VAL-06**: `defkit.Let(name, expr)` and `defkit.LetVariable()` documented — let bindings in templates

### VelaCtx

- [ ] **CTX-01**: `defkit.VelaCtx()` methods documented: Name, AppName, Namespace, Revision, ClusterVersion().Minor()

### Health & Status DSL

- [ ] **HEALTH-01**: `defkit.Health()` builder documented with `.HealthyWhen()`, `.WithDefault()`, `.WithDisableAnnotation()`
- [ ] **HEALTH-02**: `defkit.Status()` builder documented with `.IntField()`, `.Message()`
- [ ] **HEALTH-03**: Preset health builders documented: `DeploymentHealth()`, `DaemonSetHealth()`, `JobHealth()`, `DeploymentStatus()`, `DaemonSetStatus()`

### Full Examples

- [ ] **EX-01**: Complete worked Component definition example (e.g., webservice-like)
- [ ] **EX-02**: Complete worked Trait definition example (e.g., env-like patch trait)
- [ ] **EX-03**: Complete worked Policy definition example
- [ ] **EX-04**: Complete worked WorkflowStep definition example

### Code Examples Format

- [ ] **FMT-01**: Every method has a Go defkit snippet
- [ ] **FMT-02**: Every method has a side-by-side CUE equivalent
- [ ] **FMT-03**: Full definition examples show the init() registration pattern

## v2 Requirements

### Search & Navigation

- **NAV-01**: In-page search/filter for methods
- **NAV-02**: Copy-to-clipboard button on code blocks

### Extended Content

- **EXT-01**: Migration guide from CUE to defkit
- **EXT-02**: Common patterns / cookbook section

## Out of Scope

| Feature | Reason |
|---------|--------|
| Auto-generation from AST | Hand-crafted gives better narrative flow |
| Multi-page site | Single file is simpler and shareable |
| Internal/helper methods | Docs target definition authors, not defkit maintainers |
| Build toolchain | Pure HTML to avoid infra dependencies |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| STRUCT-01 | Phase 1 | Pending |
| STRUCT-02 | Phase 1 | Pending |
| STRUCT-03 | Phase 1 | Pending |
| STRUCT-04 | Phase 1 | Pending |
| STRUCT-05 | Phase 1 | Pending |
| FMT-01 | Phase 1 | Pending |
| FMT-02 | Phase 1 | Pending |
| FMT-03 | Phase 1 | Pending |
| DEF-01 | Phase 2 | Pending |
| DEF-02 | Phase 2 | Pending |
| DEF-03 | Phase 2 | Pending |
| DEF-04 | Phase 2 | Pending |
| DEF-05 | Phase 2 | Pending |
| DEF-06 | Phase 2 | Pending |
| DEF-07 | Phase 2 | Pending |
| DEF-08 | Phase 2 | Pending |
| PARAM-01 | Phase 2 | Pending |
| PARAM-02 | Phase 2 | Pending |
| PARAM-03 | Phase 2 | Pending |
| PARAM-04 | Phase 2 | Pending |
| TPL-01 | Phase 3 | Pending |
| TPL-02 | Phase 3 | Pending |
| TPL-03 | Phase 3 | Pending |
| TPL-04 | Phase 3 | Pending |
| TPL-05 | Phase 3 | Pending |
| TPL-06 | Phase 3 | Pending |
| RES-01 | Phase 3 | Pending |
| RES-02 | Phase 3 | Pending |
| RES-03 | Phase 3 | Pending |
| RES-04 | Phase 3 | Pending |
| RES-05 | Phase 3 | Pending |
| RES-06 | Phase 3 | Pending |
| RES-07 | Phase 3 | Pending |
| RES-08 | Phase 3 | Pending |
| VAL-01 | Phase 3 | Pending |
| VAL-02 | Phase 3 | Pending |
| VAL-03 | Phase 3 | Pending |
| VAL-04 | Phase 3 | Pending |
| VAL-05 | Phase 3 | Pending |
| VAL-06 | Phase 3 | Pending |
| CTX-01 | Phase 3 | Pending |
| HEALTH-01 | Phase 3 | Pending |
| HEALTH-02 | Phase 3 | Pending |
| HEALTH-03 | Phase 3 | Pending |
| EX-01 | Phase 4 | Pending |
| EX-02 | Phase 4 | Pending |
| EX-03 | Phase 4 | Pending |
| EX-04 | Phase 4 | Pending |

**Coverage:**
- v1 requirements: 48 total
- Mapped to phases: 48
- Unmapped: 0 ✓

---
*Requirements defined: 2026-03-11*
*Last updated: 2026-03-11 — added DEF-06/07/08, TPL-05/06, RES-06/07/08, VAL-05/06, HEALTH-01/02/03 from full audit*
