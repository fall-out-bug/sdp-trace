You are an adversarial reviewer for the sdp-trace repository. Review the following spec for trust-sensitive implementation. Be concise, critical, and evidence-backed. Do not give generic praise.

---

## Repository Rules
- Machine proof wins over prose.
- No deferred trust closure.
- Source-bound proof requires clean immutable source commit.
- No TODO/FIXME markers in new Go code.
- No Node.js/TS/JS in active product path.
- Bash only as thin command launcher.
- Keep root router under 100 lines.
- Module with >10 skills is too large.
- Test-first behavior for behavior changes.
- Go is the target product code.

## Current State
- `cmd/sdp-trace/` contains 500+ small Go files (~17.7K lines total).
- Files are numbered (`main_001_const.go` … `main_546_commandsurfacejson.go`) plus named files.
- Quality gates pass: strict MI 70.1, cyclomatic/cognitive <=10, CRAP strict <5, coverage-backed.
- The codebase passed strict metrics through broad same-package file splitting (PR #43).
- Review accepted result but recorded navigation and contributor-DX debt.

## Spec

### Feature Specification: Command Package Organization
**Feature Branch**: `010-command-package-organization`
**Created**: 2026-05-13
**Status**: Draft for PI review
**Input**: PR #43 closed numeric quality gates partly through broad same-package file splitting. Review accepted the result but recorded navigation and contributor-DX debt.

#### Product Boundary
This slice improves human maintainability and navigation of the CLI implementation without changing user-facing behavior.

Allowed:
- same-behavior Go refactors;
- command-family grouping;
- package or file organization changes;
- small tests that guard behavior while moving code.

Not allowed:
- changing CLI output or exit semantics unless separately specified;
- lowering quality gates;
- using file movement to hide unverified behavior;
- adding non-Go tooling.

#### Problem Statement
The codebase now passes strict metrics, but the CLI package contains hundreds of small files. That is acceptable as a metric remediation endpoint, but it is not ideal for new contributors or review agents.

The next improvement should preserve the gates while making command-family ownership visible:
- reviewers should find packet, PR review, assess, witness, query, and recorder command code quickly;
- future edits should not require scanning hundreds of similarly named files;
- package boundaries should not introduce cycles or harness-specific dependencies.

#### Core Claim
This slice may claim:
> The CLI implementation is organized by command family or another reviewed structure while preserving command behavior and quality gates.

This slice must not claim:
- a new CLI architecture;
- better metrics unless replayed;
- semantic behavior changes;
- full Clean Architecture completion beyond the reviewed organization scope.

#### Required User Stories

**US-001 - Contributor Navigation (P0)**
A contributor can identify the implementation area for a command family without scanning hundreds of files.

**US-002 - Behavior Preservation (P0)**
Existing CLI commands, help text, exit codes, and docs checks remain unchanged unless a change is explicitly reviewed.

**US-003 - Quality Gate Preservation (P0)**
CRAP, cognitive complexity, cyclomatic complexity, and strict MI remain green after any regrouping.

**US-004 - Dependency Direction (P0)**
If subpackages are introduced, dependency direction stays clear and does not make core trace/evidence packages depend on CLI packages.

#### Functional Requirements
- **FR-001**: Choose and document one organization strategy before moving code: command-family subpackages, family-prefixed files, or generated index.
- **FR-002**: Move or index code in small slices so regressions are attributable.
- **FR-003**: Keep command handlers discoverable from a single registry.
- **FR-004**: Preserve existing command names, flags, outputs, and exit code semantics.
- **FR-005**: Record any remaining high-file-count package as an advisory follow-up, not as hidden debt.

#### Acceptance Criteria
- No behavior-changing diff is introduced without a specific spec delta.
- Full local verification passes.
- PI review covers code/correctness, Clean Architecture, DX, and requirements-vs-implementation.
- Any advisory follow-ups are recorded separately from blockers.

## Plan

### Phase 0 - Organization Design Review
Before moving code, write a short design note choosing one strategy:
- command-family subpackages;
- command-family file prefixes with generated index;
- hybrid grouping for only the largest families.
Run PI review on the design before implementation.

### Phase 1 - One Vertical Slice
Move or index one command family with strong tests, preferably a smaller family first. Verify behavior, help/docs, quality gates.

### Phase 2 - Repeat By Family
Proceed family by family. Do not batch unrelated command families in one commit.

### Phase 3 - Final Audit
Run full verification, PI review planes, update contributor-facing docs only if navigation contract changed.

## Tasks
- T001 Run PI review on whether this organization work is worth doing now.
- T002 Choose one organization strategy and record the rationale.
- T003 Review the chosen strategy for dependency-cycle and behavior-preservation risk.
- T004 Stop for explicit approval before moving code.
- T010 Pick one small command family.
- T011 Move or index that family only.
- T012 Prove help, docs, exit codes, and tests remain unchanged.
- T020 Repeat by command family with scoped commits.
- T021 Keep command handlers discoverable from one registry.
- T022 Record any remaining high-file-count area as advisory debt.
- T030 Run full local verification.
- T031 Run PI code/correctness, Clean Architecture, DX, and requirements review.
- T032 Update contributor navigation docs only if the navigation contract changed.

## Review Instructions
Review this spec for whether it can reduce CLI navigation debt without undoing the quality-gate work. Focus on:
1. whether the proposed organization is small enough for safe slices;
2. whether it risks dependency cycles or import churn;
3. whether behavior preservation is testable;
4. whether the spec prevents metric-gaming from being mistaken for Clean Code.

Output format:
- Summary first, ordered by severity.
- File/line evidence where applicable.
- Disposition per finding: accepted, accepted_fixed, rejected_false_positive, deferred_not_assessed, cannot_verify, advisory.
- Note any `not_assessed` areas.
