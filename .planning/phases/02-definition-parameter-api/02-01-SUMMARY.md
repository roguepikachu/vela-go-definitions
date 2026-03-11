---
phase: 02-definition-parameter-api
plan: 01
subsystem: docs
tags: [defkit, kubevela, html, api-reference, definition-builders]

requires:
  - phase: 01-page-shell
    provides: docs/index.html shell with CSS, sidebar navigation, and placeholder sections

provides:
  - Definition Builders section in docs/index.html with 4 constructor cards and 3 special method cards
  - NewComponent card with chain method table and Go+CUE side-by-side example
  - NewTrait card with AppliesTo/ConflictsWith/PodDisruptive chain methods and example
  - NewPolicy card with minimal chain set and example
  - NewWorkflowStep card with Scope/Category chain methods and example
  - AutodetectWorkload special method card with when-to-use guidance
  - RawCUE special method card with bypass vs fluent builder comparison
  - WithImports special method card with CUE import syntax example

affects: [02-02-PLAN.md, 03-template-resource-value-velactx, 04-full-examples]

tech-stack:
  added: []
  patterns:
    - "method-card pattern: method-header (sig + meta) + method-body (applies pills + chain table + code-pair)"
    - "code-pair CSS grid: two code-block divs with .code-label.go / .code-label.cue headers"
    - "color-coded type badges: pill-component (#3b82f6), pill-trait (#fb923c), pill-policy (#f472b6), pill-workflow (#22d3ee)"

key-files:
  created: []
  modified:
    - docs/index.html

key-decisions:
  - "One method-card per constructor (NewComponent, NewTrait, NewPolicy, NewWorkflowStep) with inline chain method table rather than separate cards per chain method"
  - "Special methods (AutodetectWorkload, RawCUE, WithImports) grouped under a 'Special Methods' h3 heading after the constructor cards"
  - "CUE code blocks show generated schema to clarify the Go-to-CUE mapping, not raw CUE written by developers"

patterns-established:
  - "method-card: standard card container for all API entries — matches style used in parameter-builders and template-methods sections"
  - "code-pair: Go snippet on left, CUE equivalent on right — established as the canonical documentation format"

requirements-completed: [DEF-01, DEF-02, DEF-03, DEF-04, DEF-05, DEF-06, DEF-07, DEF-08]

duration: 15min
completed: 2026-03-11
---

# Phase 2 Plan 01: Definition Builders Summary

**Definition Builders section injected into docs/index.html — 4 constructor cards (NewComponent, NewTrait, NewPolicy, NewWorkflowStep) with per-type chain method tables, plus AutodetectWorkload, RawCUE, and WithImports special method cards**

## Performance

- **Duration:** 15 min
- **Started:** 2026-03-11T05:50:00Z
- **Completed:** 2026-03-11T06:05:00Z
- **Tasks:** 1 (+ 1 checkpoint)
- **Files modified:** 1

## Accomplishments
- Replaced the definition-builders placeholder with 7 fully-documented method-cards
- Each constructor card includes a chain method table (color-coded by definition type) and a side-by-side Go+CUE code pair
- Special methods section documents when to use AutodetectWorkload vs Workload(), RawCUE bypass behavior, and WithImports CUE import path syntax

## Task Commits

1. **Task 1: Write Definition Builders HTML** - `34deb64` (feat)

## Files Created/Modified
- `docs/index.html` — definition-builders section content injected (placeholder removed)

## Decisions Made
- Chain methods documented in a compact table per constructor (not separate cards) to keep the section scannable
- AutodetectWorkload, RawCUE, WithImports treated as "special methods" under their own h3 to distinguish them from the four main constructors

## Deviations from Plan
None - plan executed exactly as written.

## Issues Encountered
None

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- docs/index.html definition-builders section is complete and human-verified
- Parameter Builders section (02-02) can be started immediately
- CSS classes (.method-card, .code-pair, .pill-*) established and ready for reuse

---
*Phase: 02-definition-parameter-api*
*Completed: 2026-03-11*
