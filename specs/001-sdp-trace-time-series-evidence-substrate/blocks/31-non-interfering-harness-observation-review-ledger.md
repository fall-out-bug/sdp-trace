# Block 31 Review Ledger: Non-Interfering Harness Observation

Status: draft revised after Socratic review. Implementation remains blocked
until explicit approval of the reviewed direction.

Reviewed artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/31-non-interfering-harness-observation.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/31-non-interfering-harness-observation-implementation-plan.md`
- `docs/reviews/demo-jvm-gsd-observation-ledger.md`

## Intake Findings

| ID | Severity | Plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| P0-001 | critical | product-boundary / trace-evidence | No non-interfering OpenCode/GSD observation path exists for the required demo proof. Existing surfaces either inspect repository setup, wrap commands, relay prompts, or preview missing adapter-capture input. | accepted_open | `docs/reviews/demo-jvm-gsd-observation-ledger.md` lines 19-117; Block 31 draft created to define a portable harness observation intake before implementation. |

## Socratic Review Findings

| ID | Severity | Plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| S31-PUX-01 | critical | product-boundary / UX-DX | The draft said no OpenCode/GSD dependency but made the first OpenCode/GSD profile look canonical through fixtures and acceptance criteria. | accepted_fixed | Spec now states OpenCode/GSD is a dogfood exemplar, no command/schema/package/test may require it, and fixtures must include a harness-generic complete export. |
| S31-PUX-02 | major | product-boundary / UX-DX | Empty or zero-event harness exports had undefined validation behavior. | accepted_fixed | Validation mapping now distinguishes unreadable/no supplied source as `cannot_verify` and supplied zero-event runs as required-family `not_assessed`. |
| S31-PUX-03 | major | UX-DX / safety | `observe` did not define whether unsafe input fails, redacts, warns, or skips. | accepted_fixed | First surface is strict: unsafe input exits non-zero and writes no observed run; errors use safe ids and reason codes only. |
| S31-TE-01 | major | trace/evidence | Top-level validation composition only said what must not happen; mixed states were undefined. | accepted_fixed | Spec now composes required dimensions by strictest state and keeps optional dimensions dimension-level unless unsafe input fails. |
| S31-TE-02 | major | trace/evidence | `source_digest` had no canonical byte input definition. | accepted_fixed | Spec now defines SHA-256 over canonical event JSON with `source_digest` itself blanked to avoid self-referential hashing. |
| S31-TE-03 | major | trace/evidence | `redaction_state` mixed data-treatment and assessment states. | accepted_fixed | Event contract now uses `content_state`; assessment states live in validation output. |
| S31-SAFE-01 | critical | safety/privacy | Token-like values and authenticated URL detection were undefined. | accepted_fixed | Safety requirements now define first-pass token/auth URL heuristics, digest-field allowlist, and non-echoing errors. |
| S31-SAFE-02 | critical | safety/privacy | Digest canonicalization was unspecified. | accepted_fixed | Same fix as S31-TE-02. |
| S31-SAFE-03 | major | safety/privacy | `--source` and `--profile` path handling did not cover symlink escapes. | accepted_fixed | Safety requirements now require symlink evaluation and rejection of escapes; fixtures include symlink escape. |
| S31-SAFE-04 | major | safety/privacy | Summary filtering for `not_assessed` and `cannot_verify` states was undefined. | accepted_fixed | Summary output section now requires per-dimension state/reason/event counts and explicit non-authority statement without raw unsafe values. |
| S31-IMPL-01 | critical | implementation feasibility | JSON Schema validation was only syntax-checked with `jq empty`; no Go validation strategy existed. | accepted_fixed | Implementation plan now requires a Go validator choice, defaulting to `github.com/santhosh-tekuri/jsonschema/v6` behind a wrapper or equivalent focused Go contract validation. |
| S31-IMPL-02 | critical | implementation feasibility | `degradation_rules` had no concrete grammar. | accepted_fixed | Spec now defines a closed JSON degradation rule grammar with fixed states and reason codes. |
| S31-IMPL-03 | major | implementation feasibility | JSONL ingestion had no streaming or size-bound strategy. | accepted_fixed | Implementation plan now requires line-by-line bounded processing with default 1 MiB line and 100,000 event caps. |
| S31-IMPL-04 | major | implementation feasibility | `gap` event family overlapped with unavailable fields and synthesized gaps. | accepted_fixed | Event contract now distinguishes harness-emitted `gap` events from validator-derived dimensions. |
| S31-IMPL-05 | major | implementation feasibility | New package/module scope was not explicit. | accepted_fixed | Implementation plan keeps `internal/harnessobs` inside the existing Go module and prohibits harness-specific SDK imports through the existing no-dependency boundary. |
| S31-MINOR-01 | minor | product / trace / implementation | Several minor terms were underspecified: safe refs, profile/event unavailable field relationship, profile versioning, observed_at format, overwrite policy, and summarize format. | accepted_fixed | Spec and plan now define safe refs, unsupported vs unavailable fields, schema version mismatch behavior, RFC3339 timestamps, non-empty out-dir refusal, and safe human summary output. |
| S31-TE-04 | minor | trace/evidence | Reviewer proposed a separate `no_run` state for absent runs. | rejected_policy_conflict | Repository policy currently constrains live verifier states to `pass`, `fail`, `cannot_verify`, and `not_assessed`; spec instead maps absent/unreadable source to `cannot_verify` with reason and zero-event source to dimension `not_assessed`. |

## Closure State

- Spec direction: `revised_after_socratic_review`
- Implementation: `implemented_pending_pr_review`
- Demo P0 closure: `not_assessed`
- Reason: Block 31 Socratic review returned `REVISE`; accepted critical and
  major findings were fixed in the draft. The initial generic harness
  observation implementation exists, but there is no verifier-backed OpenCode/GSD
  observation evidence yet.

## Implementation Review Findings

| ID | Severity | Plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| I31-CODE-01 | major | code/correctness | CLI `harness observe` and `harness validate` relied on lower-level errors for missing documented flags. | accepted_fixed | `cmd/sdp-trace/main.go` now checks required flags before calling `internal/harnessobs`; `cmd/sdp-trace/harness_cli_test.go` covers missing flags. |
| I31-CODE-02 | major | code/correctness / safety | CLI discarded `json.MarshalIndent` errors when printing run and validation output. | accepted_fixed | `cmd/sdp-trace/main.go` now handles marshal errors and returns `exitCannotVerify`. |
| I31-TE-01 | major | trace/evidence | `not_assessed` and `cannot_verify` were both handled by a default exit-code branch, hiding intent. | accepted_fixed | `runHarnessValidate` now has explicit `StateNotAssessed` and `StateCannotVerify` cases; JSON remains authoritative for exact state. |
| I31-SAFE-01 | critical | safety/privacy | Re-review found file-write/read traversal risk through event refs and file-name ids. | accepted_fixed | `internal/harnessobs` now uses a stricter file id grammar for `event_id`, validates event refs before reading, and tests unsafe ids/refs. |
| I31-SAFE-02 | major | safety/privacy | `--out` did not reject symlink escapes before writing observed runs. | accepted_fixed | `safeOutDir` now evaluates existing symlinks and rejects paths escaping the working directory. |
| I31-SAFE-03 | major | correctness / evidence persistence | `harness validate --out` did not persist diagnostic JSON when the run was missing. | accepted_fixed | `Validate` now writes `cannot_verify` output for missing/unreadable runs when `--out` is supplied; focused test added. |
| I31-TE-02 | major | trace/evidence | First review plane could not assess implementation because the initial review packet omitted untracked files. | accepted_narrower_fixed | Re-review used a full intent-to-add diff including `internal/harnessobs`, schemas, examples, docs, and specs; trace/evidence re-review returned `APPROVE`. |
| I31-SAFE-04 | minor | safety/privacy | Focused re-review noted a possible `--out a/b` escape if `a` was an existing symlink outside the working directory. | accepted_fixed | `safeOutDir` now evaluates existing parent components for symlink escapes before `MkdirAll`; focused test added. |
| PR31-CODE-01 | critical | PR code/correctness | PR-level review flagged `validationDigest` implementation as mutating its input and confusing self-digest semantics. | accepted_fixed | `validationDigest` now hashes a local copy with `ValidationDigest` blanked. |
| PR31-CODE-02 | major | PR code/correctness | `LoadRun` did not validate `run.json` schema version. | accepted_fixed | `LoadRun` now rejects unsupported `RunSchemaVersion`; schema-version tests added. |
| PR31-REQ-01 | critical | PR requirements/safety | `harness validate --profile`, `--run`, `--out`, and `harness summarize --validation` did not all use the same path containment rules as observe. | accepted_fixed | `Validate`, `LoadProfile`, `LoadValidation`, `safeExistingDir`, and `safeOutFile` now reject traversal, absolute paths, symlink escapes, and unsafe output filenames; focused tests added. |
| PR31-TE-01 | minor | PR trace/evidence | `event_refs` schema allowed dots/colons while Go `safeFileIDPattern` did not. | accepted_fixed | `schema/harness-observation-run.schema.json` now matches Go event-ref grammar. |
| PR31-TE-02 | minor | PR trace/evidence | Unknown JSONL fields were safety-scanned but then dropped, while schema says `additionalProperties: false`. | unresolved_minor | Current Go intake remains safety-first and stores only known fields; schema validation remains stricter. No critical/major reviewer blocked merge on this. |
