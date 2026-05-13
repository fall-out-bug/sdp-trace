# Maintainability Index Gap Report

Date: 2026-05-13

This report records the current absolute MI `> 70` gap. It is improvement
planning evidence only; it is not a pass claim.

## Commands

- `go run ./tools/qualitycheck -mi-under 70 cmd internal tools`
- `go run ./tools/qualitycheck -function-mi-under 70 cmd internal tools`

Both commands exited `1` in the fresh 2026-05-13 replay.

## Current Failure Shape

| Scope | Failing rows | Notes |
| --- | ---: | --- |
| File MI | 20 failure rows plus the raw `exit status 1` line | The failures include historical large production files and remaining command/tool surfaces; the raw command output also includes one `exit status 1` line. |
| Function MI | 1275 failure rows plus the raw `exit status 1` line | The absolute function-level MI check fails again in the current tree; ratchet baselines pass, but absolute MI is not closed. |

Function-level failures by top-level area:

| Area | Failing functions |
| --- | ---: |
| `cmd` | not recounted in this replay |
| `internal` | not recounted in this replay |
| `tools` | not recounted in this replay |

Largest function-level failure clusters:

| File | Failing functions |
| --- | ---: |
| `internal/adaptercapture/adaptercapture.go` | present in latest absolute output |
| `internal/authority/authority.go` | present in latest absolute output |
| `internal/checkpoint/checkpoint.go` | present in latest absolute output |

The `internal/harnessobs/harnessobs.go` trust-boundary comment pass documented
observation setup/collection, raw-event normalization, replay validation,
filesystem containment, digest calculation, command-model extraction, and unsafe
retention checks without changing behavior. Focused `internal/harnessobs` tests
pass, and the latest absolute function-MI replay reduced total rows from 858 to
545 while this file dropped from 214 to 89. File-MI is now at 94 rows, so this
is documentation and MI-debt reduction evidence only.

The `cmd/sdp-trace/main.go` packet/PR CLI trust-boundary comment pass documented
top-level dispatch, flag-only command contracts, structured `cannot_verify`
output, packet artifact publication ordering, PR event authority, fixture versus
live GitHub Actions hydration, GitHub API trust-target checks, credential
handling, and retained artifact refs without changing behavior. Focused
`cmd/sdp-trace` tests pass, and the latest absolute function-MI replay reduced
total rows to 545 while `cmd/sdp-trace/main.go` dropped from 165 to 132.
File-MI is now at 94 rows, so this is CLI trust-documentation and MI-debt
reduction evidence only.

The `cmd/sdp-trace/main.go` PR-review CLI trust-boundary comment pass documented
packet/run/synthesis/validation/check artifact ordering, runner allow-list
reconstruction, write-once review summary behavior, preview non-evidence scope,
validation exit semantics, work-dir requirements, and repeated flag ordering
without changing behavior. Focused `cmd/sdp-trace` tests pass, and the latest
absolute function-MI replay reduced total rows to 545 while
`cmd/sdp-trace/main.go` dropped from 132 to 92 after the packet CLI split and
export/doctor comment pass.
File-MI is now at 94 rows, so this is CLI PR-review and export/doctor
trust-documentation plus MI-debt reduction evidence only.

The packet `build-pr` GitHub input path pass documented CLI-to-evidence
handoff, required flag validation, fixture/live source boundaries, event-file
authority, request construction before credential egress, and GitHub artifact
response decoding. The dense PR event mapping was split into PR identity and
commit-range helpers without changing behavior. Focused `cmd/sdp-trace` tests
pass, and the latest absolute function-MI replay reduced total rows from 344 to
318 while `cmd/sdp-trace/main.go` dropped from 81 to 79 and the targeted
`runPacketBuildPR`, `parsePacketBuildPROptions`, `loadPRInputSourceEvent`,
`githubPRInputFromEvent`, and `fetchGitHubActionsArtifacts` rows cleared the
absolute function-MI failure output. File-MI remains at 91 rows, so this is
packet trust-documentation and MI-debt reduction evidence only.

The follow-up ratchet fix corrected the posture export header helper so the
current dirty checkout compiles after the prior header split, then documented
override request parsing rather than weakening the baseline. Focused
`internal/posture` and `cmd/sdp-trace` tests pass, and the latest absolute
that function-MI replay reduced total rows from 318 to 283 while
`cmd/sdp-trace/main.go` dropped from 79 to 75. File-MI remained at 91 rows, so
that was compile restoration and MI-ratchet enforcement evidence only.

The parallel CLI/posture worker plus telemetry comment pass extracted shared
CLI assess-preview output handling, moved static preview vocabularies to
package scope, reused override required-field validation, completed posture
helper extraction to clear all `internal/posture/posture.go` function-MI rows,
and documented telemetry aggregation/rendering boundaries. Focused tests and
strict CRAP for `cmd/sdp-trace`, `internal/posture`, and `internal/telemetry`
pass. That replay reported 265 function rows and 91 file rows:
`cmd/sdp-trace/main.go` was 67 rows, `internal/posture/posture.go` was 0
rows, and `internal/telemetry/prometheus.go` was 4 rows. This remains
MI-debt reduction evidence only, not absolute MI closure.

The PR-review synthesis/validation and interaction CLI pass documented review
artifact collation, validation persistence, summary non-proof scope, interaction
router contracts, relay command boundaries, transcript import paths, and summary
derivation. The relay option/default construction was split into focused
helpers without changing behavior. In parallel, a bounded worker extracted
helpers and added trust comments in `internal/prreview/prreview.go`; focused
`cmd/sdp-trace` and `internal/prreview` tests pass, and
`internal/prreview/prreview.go` now has zero absolute function-MI rows. The
latest absolute function-MI replay reduced total rows from 265 to 257 while
`cmd/sdp-trace/main.go` dropped from 67 to 60. File-MI remains at 91 rows, so
this is review/interaction trust-documentation and MI-debt reduction evidence
only.

The strict CRAP repair pass added focused tests for harness setup action,
isolation/degradation rule, normalized event write, shell parsing,
source-commit, validation decoding, and demo checkpoint-state mapping helper
branches. Full local verification passed after regenerating MI ratchets from the
then-current checkout. That absolute function-MI replay was 220 rows and 91 file
rows; `cmd/sdp-trace/main.go` remains at 60 function rows. This is coverage and
ratchet-stabilization evidence only, not absolute MI closure.

The follow-up parallel MI wave cleared `internal/harnessobs/harnessobs.go`,
`internal/demo/demo.go`, `internal/witness/profiles.go`,
`internal/checkpoint/checkpoint.go`, and `internal/policy/authority.go` from the
absolute function-MI output. Local slices also reduced
`internal/ciartifact/ciartifact.go` from 8 to 3 rows,
`internal/recorder/recorder.go` from 8 to 3 rows,
`internal/adaptercapture/adaptercapture.go` from 7 to 4 rows,
`internal/contract/contract.go` from 5 to 1 row, and
`cmd/sdp-trace/main.go` from 60 to 55 rows. Focused tests and package strict
CRAP checks pass for the touched packages. The current absolute replay is 190
function rows and 91 file rows, so this remains MI-debt reduction evidence only.

The assess/envelope CLI pass extracted shared assess flag registration and
preview report builders, documented assess/envelope dispatch boundaries, and
kept preview output explicitly non-verdict. In parallel, a bounded worker
reduced `internal/witness/profiles.go` to zero function-MI rows by extracting
Customer PKI validation preparation and authority issue ordering helpers while
preserving precedence. Focused `cmd/sdp-trace` and `internal/witness` tests
pass. The latest absolute function-MI replay is 170 rows and 91 file rows, with
`cmd/sdp-trace/main.go` at 45 rows. This is helper-extraction and
trust-documentation evidence only, not absolute MI closure.

The CLI explanation/checkpoint pass reused shared reason/action renderers,
extracted CI artifact family explanation rendering, moved checkpoint flag
defaults into helper tables, and documented checkpoint verification boundaries.
In parallel, a bounded worker cleared all function-MI rows in
`internal/verifier/explain.go` through explanation helper cleanup. Focused
`cmd/sdp-trace` and `internal/verifier` tests pass. The latest absolute
function-MI replay is 150 rows and 91 file rows, with `cmd/sdp-trace/main.go`
at 39 rows. This is helper-extraction and trust-documentation evidence only,
not absolute MI closure.

The CLI/verifier/query/authority/witness continuation documented remaining CLI
verifier, query, query-pack, doctor, CI-witness, OIDC, and witness artifact
trust boundaries. Bounded workers cleared `internal/verifier/verify.go`,
`internal/verifier/explain.go`, `internal/query/query.go`,
`internal/query/querypack.go`, and `internal/authority/authority.go` from the
absolute function-MI output; the local integrator reduced
`cmd/sdp-trace/main.go` from 45 to 42 rows and `internal/witness/witness.go`
from 4 to 3 rows. Focused package tests and focused strict CRAP replay pass for
the touched packages. The latest absolute replay is 153 function rows and 91
file rows, so this remains MI-debt reduction evidence only.

The CLI gate/checkpoint continuation reused shared flag registration, added
protected gate input and preview report helpers, and documented near-threshold
gate explanation and witness-artifact trust boundaries. A telemetry worker
cleared `internal/telemetry/prometheus.go` function-MI rows by extracting typed
aggregate helpers, and an export worker cleared `internal/export/export.go`
function-MI rows by splitting bundle construction, event selection, write, and
read paths. The local integrator also repaired the existing
`internal/repoobserver` config payload helper so the checkout compiled again.
Focused `cmd/sdp-trace`, `internal/telemetry`, `internal/export`, and
`internal/repoobserver` tests pass. A follow-up parser split restored strict
CRAP after `parseGateArgs` crossed the threshold. The latest absolute replay is
66 function rows and 90 file rows, so this remains MI-debt reduction evidence
only.

The CLI/quality-tool/small-internal continuation documented near-threshold CLI
parser and command-boundary rows, cleared assigned `tools/qualitycheck` parser
and scoring rows, cleared assigned `tools/crapcheck` and
`tools/mibaselinepolicy` rows, and directly verified unreported internal changes
in witness, repo-observer, recorder, interaction, export, and CI artifact
packages after the responsible worker failed to return a usable result.
Focused tests and strict CRAP replay pass for the touched packages. The latest
absolute replay is 0 function rows and 66 file rows, so this remains MI-debt
reduction evidence only.

The near-threshold tool file-MI continuation cleared assigned rows in
`tools/crapcheck`, `tools/mibaselinepolicy`, and `tools/qualitycheck` through
small helper splits and invariant comments. Focused tool tests pass, and the
absolute file-MI replay now reports 66 file rows. This is file-MI reduction
evidence only; absolute file MI remains open.

The current-tree ratchet repair regenerated both MI baseline files from the
actual files present on disk and fixed `tools/crapcheck` path normalization so
repo-relative paths under `cmd/sdp-trace` are not truncated during coverage /
complexity joins. Focused `tools/crapcheck` tests pass, strict CRAP replay now
passes against `/tmp/sdp-trace-cover-final-local.out`, and both MI ratchet commands
exit 0. Absolute MI remains open: function MI exits 1 with 1301 stderr rows and
file MI exits 1 with 29 stderr rows.

The `internal/trace` file-layout slice split the trace event, manifest,
contract, verifier, payload, validation, store, and source-snapshot surfaces
into same-package files without changing exported APIs. Focused
`internal/trace` tests pass, `internal/trace` absolute function MI and file MI
both exit 0, and repository absolute file MI drops to 29 stderr rows while
absolute function MI remains open at 1301 stderr rows. This is file-MI
reduction evidence only; repository-wide absolute MI remains open.

The `internal/export` bundle-layout slice split bundle construction, JSON I/O,
and run-manifest path handling into same-package files without changing the
public API. Focused `internal/export` tests pass, `internal/export` absolute
function MI and file MI both exit 0, repository absolute file MI drops to 28
stderr rows, and absolute function MI drops to 1296 stderr rows. This is
export-surface MI reduction evidence only; repository-wide absolute MI remains
open.

The `internal/contract` evidence-contract slice split digest construction,
validation, and validation field tables into same-package files without changing
the public API. Focused `internal/contract` tests pass, `internal/contract`
absolute function MI and file MI both exit 0, repository absolute file MI drops
to 27 stderr rows, and absolute function MI drops to 1288 stderr rows. This is
contract-surface MI reduction evidence only; repository-wide absolute MI remains
open. The pi worker attempt for this slice timed out and is not counted as
evidence; the ported implementation was verified locally.

The `internal/policy` and `internal/query` file-layout slice split authority
policy loading/validation/validator concerns and capture-depth query helpers
into same-package files. Focused `internal/policy` and `internal/query` tests
pass, the new split files pass the MI ratchets, and `internal/policy` plus the
former `internal/query/query.go` row drop out of the absolute file-MI output.
Repository absolute file MI drops to 24 stderr rows, and absolute function MI
drops to 1275 stderr rows. The pi query worker was cancelled because it wrote
into the main checkout and left new below-threshold rows; the local repairs are
the counted evidence.

The `cmd/sdp-trace/harness_cli.go` file-layout slice split the harness router,
observe, validate, validate argument/exit mapping, and summarize surfaces into
same-package files. Focused `cmd/sdp-trace` tests pass, all new harness files
pass absolute file/function MI, and `cmd/sdp-trace/observe_cli.go` plus
`cmd/sdp-trace/harness_cli.go` drop out of the absolute file-MI output.
Repository absolute file MI drops to 22 failure rows plus the raw
`exit status 1` line; absolute function MI remains at 1275 failure rows plus
the raw `exit status 1` line.

The `cmd/sdp-trace/packet_cli.go` and `internal/telemetry/prometheus.go`
file-layout slice split packet subcommands and Prometheus telemetry helpers into
same-package files without changing exported APIs. Focused `cmd/sdp-trace`,
`internal/telemetry`, and `internal/releaseproof` tests pass. Packet and
telemetry now pass absolute file/function MI. Repository absolute file MI drops
to 20 failure rows plus the raw `exit status 1` line; absolute function MI
remains at 1275 failure rows plus the raw `exit status 1` line. The
releaseproof worker output was discarded because it increased the file-MI
failure list.

## File-Level Failures

- `cmd/sdp-trace/main.go`
- `internal/adaptercapture/adaptercapture.go`
- `internal/authority/authority.go`
- `internal/checkpoint/checkpoint.go`
- `internal/ciartifact/ciartifact.go`
- `internal/demo/demo.go`
- `internal/forensic/forensic.go`
- `internal/harnessobs/harnessobs.go`
- `internal/interaction/interaction.go`
- `internal/managed/managed.go`
- `internal/packet/packet.go`
- `internal/posture/posture.go`
- `internal/prreview/prreview.go`
- `internal/query/querypack.go`
- `internal/recorder/recorder.go`
- `internal/releaseproof/releaseproof.go`
- `internal/repoobserver/repoobserver.go`
- `internal/verifier/verify.go`
- `internal/witness/profiles.go`
- `internal/witness/witness.go`

## Closure Strategy

1. Keep the current ratchets as the active gate. They block new below-threshold
   functions/files and regressions without pretending historical code is clean.
2. Retire the largest clusters through vertical slices, not mechanical comment
   churn. Good first candidates are CLI command-family extraction from
   `cmd/sdp-trace/main.go`, harness observation parser/package splits, and
   packet/prreview/posture package splits by command or artifact boundary.
3. Treat tool MI separately from product MI. Tool files are already under CRAP,
   complexity, coverage, and ratchet gates; absolute file MI remains low because
   the current MI formula penalizes compact parser/analyzer files heavily.
4. Update this report after each MI retirement slice with the same two commands
   above, and only remove the `assessed_gap` label when both absolute commands
   exit `0`.

## 2026-05-12 Slice Note

The observe/harness CLI helpers were moved from `cmd/sdp-trace/main.go` into
same-package `cmd/sdp-trace/observe_cli.go` and
`cmd/sdp-trace/harness_cli.go`. Local tests and MI ratchets pass after baseline
regeneration. The slice reduced `cmd/sdp-trace/main.go` from the previous 5626
lines / cyclo 735 shape to 5366 lines / cyclo 701, then split the restored
helper file from one 13.7 file-MI surface into `observe_cli.go` at 23.1 and
`harness_cli.go` at 26.3. Absolute file MI is still below 70, so this is
measurable layout progress, not MI closure.

The `tools/crapcheck/read.go` parser file was split into `coverage.go`,
`complexity.go`, and `rows.go`. Focused tests and ratchets pass, and the largest
parser file MI improved from `read.go` at 26.3 to `coverage.go` at 36.3 and
`complexity.go` at 37.6. All three files still remain below 70, so the absolute
file-MI failure count increases until deeper parser simplification or a revised
MI policy is approved.

## 2026-05-13 Slice Note

The `tools/qualitycheck/report.go` reporting file was split into `report.go`,
`report_baseline.go`, and `report_failures.go`. Focused tests and ratchets pass,
and the largest reporting file MI improved from `report.go` at 20.2 to
`report.go` at 35.0, `report_baseline.go` at 46.2, and `report_failures.go` at
32.2. All three files still remain below 70. The current absolute file-MI
failure count is 93 after later slices; this row is layout improvement evidence,
not MI closure.

The `tools/qualitycheck/complexity.go` analyzer file was split into
`complexity.go`, `complexity_clause.go`, `complexity_score.go`,
`complexity_statement.go`, `boolean.go`, and `cyclomatic.go`. Focused tests and
ratchets pass, and the largest complexity file MI improved from
`complexity.go` at 17.5 to helper files ranging from 31.3 to 53.0. These files
still remain below 70, so this remains layout progress rather than MI closure.

The `tools/qualitycheck/baseline.go` and `tools/qualitycheck/analyze.go`
helpers were split by responsibility into smaller same-package files. Focused
tests and ratchets pass. The baseline surface improved from `baseline.go` at
24.9 to files ranging from 42.6 to 50.9, and the analyzer surface improved from
`analyze.go` at 23.2 to files ranging from 38.2 to 56.3. These files still
remain below 70. The current absolute file-MI failure count is 93 after later
slices; this is layout improvement evidence, not MI closure.

The `tools/crapcheck/score.go` scoring/join surface was split into `score.go`,
`join.go`, and `scoring.go` by an isolated worker. Focused tests and ratchets
pass. The original `score.go` surface improved from 32.5 to helper files
ranging from 37.9 to 55.8. These files still remain below 70, so this is layout
improvement evidence, not MI closure.

The `tools/qualitycheck/main.go` option parsing surface was split into
`main.go`, `options.go`, and `option_flags.go`. Focused tests and ratchets pass,
and `main.go` improved from MI 27.9 to 36.0. The extracted option files remain
below 70, so the absolute file-MI failure count is now 93 and the function-MI
failure count is now 545; this is layout improvement evidence, not MI closure.

The `tools/qualitycheck/halstead.go` token and MI formula surface was split into
`halstead.go`, `maintainability.go`, and `round.go`. Focused tests and ratchets
pass, and `halstead.go` improved from MI 30.5 to 35.8. The extracted files
remain below 70, so the absolute file-MI failure count is now 93 while the
function-MI failure count remained 545; this is layout improvement evidence,
not MI closure.

The `tools/qualitycheck/report_failures.go` failure-predicate surface was split
into threshold, function-MI, and file-MI files. Focused tests and ratchets pass,
and the original `report_failures.go` surface improved from MI 32.2 to helper
files ranging from 46.4 to 47.6. These files still remain below 70, so this is
layout improvement evidence, not MI closure.

The `tools/mibaselinepolicy/policy.go` policy helper surface was split by an
isolated worker into `policy.go`, `baseline_policy.go`, and
`changed_files.go`. Focused tests and ratchets pass, and the original
`policy.go` surface improved from MI 39.3 to helper files ranging from 50.5 to
57.6. These files still remain below 70, so this is layout improvement evidence,
not MI closure.

The `tools/qualitycheck/discover.go` path helper surface was split to
`path.go`. Focused tests and ratchets pass. The extracted `path.go` has
function-level MI above 70 for its helpers, but file-level MI is still 53.3, so
it remains part of the absolute file-MI gap.

The `tools/qualitycheck/names.go` receiver parsing surface was split into
`names.go` and `receiver.go`. Focused tests and ratchets pass, and `names.go`
improved from MI 37.6 to 54.7. The extracted receiver helper remains below 70,
so the absolute file-MI failure count is now 93 while the function-MI failure
count remained 545; this is layout improvement evidence, not MI closure.

The `tools/qualitycheck/report.go` output surface was split into report
orchestration, function-report, file-report, and output helpers. Focused tests
and ratchets pass, and `report.go` improved from MI 35.0 to helper files ranging
from 47.6 to 63.0. These files still remain below 70, so this is layout
improvement evidence, not MI closure.

The `tools/crapcheck/coverage.go` coverage parser was split by an isolated
worker into `coverage.go` and `coverage_line.go`. Focused tests and ratchets
pass, and the original `coverage.go` surface improved from MI 36.3 to helper
files at MI 44.5 and 46.0. These files still remain below 70, so this is layout
improvement evidence, not MI closure.

The `tools/qualitycheck/source.go` source slicing and line-counting helpers
were split into `source.go` and `lines.go`. Focused tests and ratchets pass, and
`source.go` improved from MI 45.4 to 49.6. The extracted line helper remains
below 70, so the absolute file-MI failure count is now 93 while the function-MI
failure count is now 545; this is layout improvement evidence, not MI closure.

The `tools/qualitycheck/measure.go` function and file measurement surface was
split into `measure.go`, `measure_file.go`, and `measure_function.go`; the
previously half-applied discovery walk extraction was also completed as
`discover_walk.go`. Focused tests and ratchets pass, `measure.go` improved from
MI 38.2 to 53.4, and `discover.go` improved to 43.4. The extracted helper files
remain below 70, so the absolute file-MI failure count is now 93 and the
function-MI failure count is now 545; this is layout improvement evidence, not
MI closure.

The `tools/crapcheck/complexity.go` gocyclo parser was split by an isolated
worker into `complexity.go` and `complexity_line.go`. Focused tests and ratchets
pass, and the original `complexity.go` surface improved from MI 37.6 to helper
files at MI 46.2 and 46.6. These files still remain below 70, so this is layout
improvement evidence, not MI closure.

The `tools/qualitycheck/baseline_read.go` and `baseline_write.go` JSON plumbing
was deduplicated into `baseline_io.go`. Focused tests and ratchets pass;
`baseline_read.go` improved from MI 43.9 to 49.1 and `baseline_write.go`
improved from MI 48.7 to 60.0. The shared IO helper remains below 70, so the
absolute file-MI failure count is now 93 while the function-MI failure count
remains 545; this is duplication-reduction evidence, not MI closure.

The `tools/crapcheck/join.go` support helpers were split into `join.go`,
`coverage_match.go`, and `row_key.go`. Focused tests and ratchets pass, and
`join.go` improved from MI 37.9 to 44.6. The extracted helper files remain below
the file-level threshold, so the absolute file-MI failure count is now 93 while
the function-MI failure count remained 545; this is layout improvement evidence,
not MI closure.

The `tools/qualitycheck/boolean.go` one-use `booleanOperatorIncrement` helper was
inlined into `cyclomaticIncrement`. Focused tests and ratchets pass, and the
absolute function-MI failure count dropped from 2022 to 2021 in that historical slice; latest replay at that point was 545 while file-MI remained
at 94 rows. This is cleanup evidence, not MI closure.

The single-use `tools/qualitycheck/options.go` helpers `parsedPaths` and
`parsedOptions` were folded into `parseOptions`. Focused tests and ratchets pass,
and reduced function-MI debt before the later CRAP-directed discovery split.
File-MI is now at 94 rows. This is cleanup evidence, not MI closure.

The one-use `tools/qualitycheck/discover.go` and `discover_walk.go` helpers
`productionGoFile`, `skipDir`, `goFilesInDir`, `walkGoFiles`, `collectDir`, and
`collectFile` were folded into their callers. Focused tests and ratchets pass;
the then-current absolute function-MI failure count was 545 while file-MI is now at 94 rows. This is cleanup evidence, not MI closure.

The single-use `tools/qualitycheck/baseline_keys.go` helpers
`baselineRecordsByKey` and `fileBaselineRecordsByKey` were folded into the
baseline loaders. Focused tests and ratchets pass; the current absolute
function-MI failure count had dropped to 545 while file-MI is now at 94 rows. This
is cleanup evidence, not MI closure.

The duplicate file/function report baseline loaders were consolidated into one
generic `loadMIBaselineForReport` helper. Focused tests and ratchets pass; the
then-current absolute function-MI failure count was 545 while file-MI remained at 93
rows. This is cleanup evidence, not MI closure.

The duplicate typed MI baseline readers were consolidated into one generic
`readMIBaseline` helper while preserving the file/function schema error text.
Focused tests and ratchets pass; the current absolute function-MI failure count
was 545 while file-MI is now at 94 rows. This is cleanup evidence, not MI
closure.

The duplicate file/function MI baseline record comparison helpers were routed
through one generic `miRecordFails` helper while keeping the tested wrapper names
and error messages. Focused tests and ratchets pass; the current absolute
function-MI failure count was 545 while file-MI is now at 94 rows. This is
cleanup evidence, not MI closure.

The single-use `tools/qualitycheck` MI wrappers `functionMIRecordFails` and
`fileMIRecordFails` were folded into their threshold callers while keeping the
rounded-baseline comparison covered through the threshold-helper test. Focused
tests pass; the then-current absolute function-MI failure count was 545 while file-MI
remains at 94 rows. This is cleanup evidence, not MI closure.

The single-use `tools/qualitycheck/maintainability.go` helper `commentPercent`
was folded into `maintainabilityIndex` as an explicit guarded comment-ratio
calculation. Focused tests pass; the current absolute function-MI failure count
was 545 while file-MI is now at 94 rows. This is cleanup evidence, not MI
closure.

The branch-free `tools/qualitycheck/complexity_score.go` helper
`nestedControlScore` was folded into its three cognitive scoring call sites.
Focused tests pass and gocyclo still reports `loopStatementScore` and
`switchStatementScore` at 1 and `cognitiveIfScore` at 3. The current absolute
function-MI failure count remained 545 while file-MI is now at 94 rows. This is
simplification evidence, not MI closure.

The duplicate function/file MI baseline exit handlers in `tools/qualitycheck`
were replaced by one `writeBaselineExit` helper with direct dispatch from
`writeRequestedBaseline`. Focused tests pass, the focused package coverage is
91.0%, and the then-current absolute function-MI failure count was 545 while
file-MI is now at 94 rows. This is duplication-reduction evidence, not MI
closure.

The single-use `tools/qualitycheck` `reportExit` wrapper was folded into
`runWithOptions`. Focused tests pass; `runWithOptions` is cyclomatic 4 with
77.8% focused coverage, and the current absolute function-MI failure count is
545 while file-MI is now at 94 rows. This is cleanup evidence, not MI closure.

The single-use `tools/qualitycheck/comments.go` helper `commentGroupLines` was
folded into `commentLinesInRange`. Focused tests and ratchets pass; the current
absolute function-MI failure count was 545 while file-MI is now at 94 rows. This
is cleanup evidence, not MI closure.

The single-use `tools/crapcheck/path.go` helper `normalizeAbsoluteFile` was
folded into `normalizeFile`. Focused tests and ratchets pass; the current
absolute function-MI failure count was 545 while file-MI is now at 94 rows. This
is cleanup evidence, not MI closure.

The `tools/crapcheck` parser and input-error paths gained focused tests for
skipped coverage rows, malformed coverage percentages, unrecognized gocyclo
input, missing input files, and `loadRows` complexity read failures. Focused
package coverage reached 94.7%. This is verifier hardening evidence, not MI
closure.

Halstead token length is now tracked during scanning, which removed the separate
`totalTokenCount` helper while keeping the small skip/key helpers to preserve
strict CRAP. Focused `tools/qualitycheck` tests pass; Halstead functions remain
under the complexity ratchet at gocyclo 3, 3, 4, 2, and 2. The current absolute
function-MI failure count is 545 while file-MI is now at 94 rows. This is
function-debt reduction evidence, not MI closure.

Small witness trust-context and customer-PKI predicates were collapsed into
direct boolean expressions without changing semantics. Focused `internal/witness`
tests pass; the current absolute function-MI failure count is 545 while
file-MI is now at 94 rows. This is function-debt reduction evidence, not MI
closure.

Small `cmd/sdp-trace` read/boolean predicates were collapsed without changing
semantics. Focused `cmd/sdp-trace` tests and strict CRAP pass; the current
absolute function-MI failure count is 545 while file-MI is now at 94 rows.
This is function-debt reduction evidence, not MI closure.

Small adapter-capture overclaim predicates were collapsed without changing
semantics. Focused `internal/adaptercapture` tests and strict CRAP pass; the
current absolute function-MI failure count is 545 while file-MI remains at 93
rows. This is function-debt reduction evidence, not MI closure.

`internal/demo.isRunDirCandidate` was collapsed into a direct short-circuit
predicate without changing semantics. Focused `internal/demo` tests and strict
CRAP pass; the current absolute function-MI failure count is 545 while file-MI
remains at 94 rows. This is function-debt reduction evidence, not MI closure.

Small authority policy boolean predicates were collapsed without changing
semantics. Focused `internal/policy` tests and strict CRAP pass; the current
absolute function-MI failure count is 545 while file-MI is now at 94 rows.
This is function-debt reduction evidence, not MI closure.

Small managed boundary and witness identity predicates were collapsed without
changing semantics. Focused `internal/managed` tests and strict CRAP pass; the
current absolute function-MI failure count is 545 while file-MI remains at 93
rows. This is function-debt reduction evidence, not MI closure.

`internal/posture.checkedDigest` now performs its only digest comparison
directly, removing the single-use `checkDigestMatch` helper without changing
the returned digest or `errDigestMismatch` behavior. Focused `internal/posture`
tests and strict CRAP pass; the current absolute function-MI failure count is
545 while file-MI is now at 94 rows. This is function-debt reduction evidence,
not MI closure.

`internal/authority.resultReasons` now collects evaluation and binding reason
codes directly, removing two single-use loop helpers without changing reason
deduplication or ordering. Focused `internal/authority` tests and strict CRAP
pass; the current absolute function-MI failure count is 545 while file-MI
remains at 94 rows. This is function-debt reduction evidence, not MI closure.

`cmd/sdp-trace.demoWitnessArtifacts` now records the first run id directly in
the artifact loop, removing the single-use `firstNonEmpty` helper without
changing artifact ordering or digest handling. Focused `cmd/sdp-trace` tests and
strict CRAP pass; the current absolute function-MI failure count is 545 while
file-MI is now at 94 rows. This is function-debt reduction evidence, not MI
closure.

`internal/managed.missingEventGroupCondition` now derives the file-event reason
prefix directly, removing the single-use `groupReasonPrefix` helper while
preserving `file_mutation` reason codes. Focused `internal/managed` tests and
strict CRAP pass; the current absolute function-MI failure count is 545 while
file-MI is now at 94 rows. This is function-debt reduction evidence, not MI
closure.

`internal/harnessobs.normalizePotentialParentPath` now applies its default
current-directory value directly, removing the single-use
`defaultRelativeParentPath` helper without changing path validation or cleaning.
Focused `internal/harnessobs` tests and strict CRAP pass; the current absolute
function-MI failure count is 545 while file-MI is now at 94 rows. This is
function-debt reduction evidence, not MI closure.

`internal/query.safeToken` now returns `unknown` directly for fully sanitized
empty tokens, removing the single-use `tokenOrUnknown` helper while preserving
query-pack row token fallback behavior. Focused `internal/query` tests and
strict CRAP pass; the current absolute function-MI failure count is 545 while
file-MI is now at 94 rows. This is function-debt reduction evidence, not MI
closure.

`internal/forensic.criticalEvents` now adds policy critical-event families
directly, removing the single-use `addCriticalPolicyEvents` helper while
preserving downgrade handling through `removeDowngradedEvents`. Focused
`internal/forensic` tests and strict CRAP pass; the current absolute function-MI
failure count is 545 while file-MI is now at 94 rows. This is function-debt
reduction evidence, not MI closure.

Small retained posture refusal, digest-reason, grouping-key, and unsafe-output
predicates were folded into their only callers while restoring guard helpers
needed to keep strict CRAP below 5. Focused `internal/posture` tests and strict
CRAP pass; the current absolute function-MI failure count is 545 while file-MI
remains at 94 rows. This is function-debt reduction evidence, not MI closure.

`internal/releaseproof.releaseStateForCommitStatus` was folded into `Evaluate`
while preserving the `cannot_verify` state for missing source commits and the
source-boundary rationale. `sourceCommitExists` and `artifactBytes` also now
carry explicit source-boundary comments for the git object and source-commit
artifact reads that make local release proof source-bound. Focused
`internal/releaseproof` tests and function/file MI ratchets pass without
regenerated MI baselines. The current absolute function-MI failure count is 545 while
file-MI is now at 94 rows. This is trust-documentation and function-MI reduction
evidence, not MI closure.

`internal/releaseproof.combineState` and `applyDirtyState` now document that
artifact verification and dirty-checkout evidence can only lower source-bound
confidence. Focused `internal/releaseproof` tests pass; function/file MI ratchets
pass after regenerated MI baselines. The current absolute function-MI failure count is
545 while file-MI is now at 94 rows. This is trust-documentation and
function-MI reduction evidence, not MI closure.

Additional `internal/releaseproof` comments document manifest path normalization,
symlink containment checks, artifact-count denominator semantics, persisted proof
JSON shape, proof read error attribution, and git repository-root discovery.
Focused `internal/releaseproof` tests pass; function/file MI ratchets pass
after regenerated MI baselines. The current absolute function-MI failure count is 545
while file-MI is now at 94 rows. This is trust-documentation and function-MI
reduction evidence, not MI closure.

`internal/releaseproof.artifactCounts` now documents that digest hex case is
formatting while reported actual digests remain canonical lowercase. Focused
`internal/releaseproof` tests pass; function/file MI ratchets pass without
regenerated MI baselines. The current absolute function-MI failure count is 545 while
file-MI is now at 94 rows. This is trust-documentation and function-MI reduction
evidence, not MI closure.

Small tool policy helpers now document the active product path scope, strict
CRAP `< N` threshold semantics, skipped qualitycheck directories, Halstead token
normalization, and quiet diagnostic output behavior. Focused tool tests pass;
the current absolute function-MI failure count is 545 while file-MI
remained at 94 rows. This is verifier-policy documentation and function-MI
reduction evidence, not MI closure.

Trust-boundary comments were added to managed verdict aggregation, authority
allow/deny and evidence-resolution rules, witness OIDC audience and CI identity
handling, customer-PKI fallback selection, GitHub artifact request media type,
and human-facing explanation rows. Focused tests for `internal/managed`,
`internal/authority`, `internal/witness`, and `cmd/sdp-trace` pass; the
current absolute function-MI failure count is 545 while file-MI
remained at 94 rows. This is trust-documentation and function-MI reduction
evidence, not MI closure.

Checkpoint and CI artifact comments now document untrusted checkpoint shape,
independent binding checks, source-binding failures, artifact family access and
binding defaults, visible reason/action preservation, producer authority,
binding severity aggregation, safe authority fallback, and genesis sequence
rules. Focused `internal/checkpoint` and `internal/ciartifact` tests pass; the
then-current absolute function-MI failure count was 1849 while file-MI remained
at 94 rows. This is trust-documentation and function-MI reduction evidence, not
MI closure.

Packet and harness observation comments now document packet metadata ordering,
closed state vocabularies, pass-row evidence requirements, contradiction
handling, decision-owner accountability, replayable observation/session times,
OpenCode signal family classification, run schema loading, content-state
validation, command digesting, raw-event field skipping, operation-ref safety,
and unsafe event masking. Focused `internal/packet` and `internal/harnessobs`
tests pass; the then-current absolute function-MI failure count was 1849 while
file-MI is now at 94 rows. This is trust-documentation and function-MI
reduction evidence, not MI closure.

Adapter-capture and forensic next-action helpers now document that remediation
is derived from unresolved condition evidence and deduplicated before report
output; forensic critical-event policy additions now document that downgrades
apply only after the required event surface is expanded. Focused
`internal/adaptercapture` and `internal/forensic` tests pass; the then-current
current absolute function-MI failure count is 545 while file-MI is now at 94 rows.
This is trust-documentation and function-MI ratchet repair evidence, not MI
closure.

Posture, interaction, and query-pack comments now document selection schema and
repository requirements, digest mismatch trust mapping, safe label/path
normalization, deterministic export identifiers, partial/cannot-verify trace
state ordering, optional ID safety, content-ref binding, sorted verifier gap
rows, unavailable-state gap families, deterministic query row appends,
adapter task-family routing, and malformed required input fan-out. Focused
`internal/posture`, `internal/interaction`, and `internal/query` tests pass; the
then-current absolute function-MI failure count was 1849 while file-MI remained
at 94 rows. This is trust-documentation and function-MI reduction evidence, not
MI closure.

`internal/demo.evaluateRequiredRuns` now applies the observation-profile default directly, removing the single-use `runProfile` helper while preserving required run profile fallback behavior. Focused `internal/demo` tests and strict CRAP pass; the then-current absolute function-MI failure count was 1849 while file-MI is now at 94 rows. This is function-debt reduction evidence, not MI closure.

`internal/authority.approvalFailureState` now documents its trust boundary: missing approval evidence is outside authority, while unresolved external approval references remain `cannot_verify`. Focused `internal/authority` tests and strict CRAP pass; the current absolute function-MI failure count is 545 while file-MI is now at 94 rows. This is trust-documentation and function-debt reduction evidence, not MI closure.

Repo-observer and demo comments now document portable default profile selection,
profile validation, explicit repository identity precedence, proof aggregation,
human gap detail rendering, repository target containment, executable target
modes, URL credential redaction, repository-identity evidence scope, gitignore
marker limits, missing workflow replay semantics, new generated-file writes,
generated README evidence boundaries, checkpoint/gate state fallbacks,
required-run cannot-verify reasons, file digest error handling, and replayable
summary counters. Focused `internal/repoobserver` and `internal/demo` tests
pass; `internal/repoobserver` function-MI failures dropped to 18, and the
current absolute function-MI failure count is 545 while file-MI remains at 93
rows. This is trust-documentation and function-MI reduction evidence, not MI
closure.

Small validation, trace, review, recorder, and CLI helpers now document required
contract fields, authority-policy identifiers, optional event-chain validation,
append-only sequence/count semantics, reviewer status/default metadata,
concurrent stream digest capture, packet/review exit-code meaning, runner
allow-list parsing, known flag classification, envelope witness input, preview
install exits, and fixture-root defaults. Focused package tests pass; the
current absolute function-MI failure count is 545 while file-MI remains at 93
rows. This is trust-documentation and function-MI reduction evidence, not MI
closure.

Release-proof, witness, and harness comments now document source-commit
artifact lookup versus manifest report paths, missing and mismatched manifest
obligations, JWT-like token detection by shape, and unavailable-field metadata
safety. Focused `internal/releaseproof`, `internal/witness`, and
`internal/harnessobs` tests pass; the current absolute function-MI failure count
is 545 while file-MI is now at 94 rows. This is trust-documentation and
function-MI reduction evidence, not MI closure.

Witness profile comments now document Customer PKI failure-field fallback,
profile decision alignment, artifact path binding, demo run-set binding,
non-empty run ID handling, mandatory Customer PKI input shape, unsafe input
path/symlink rejection, public trust-anchor loading, and SHA-256-sized freshness
digests. Focused `internal/witness` tests pass; `internal/witness/profiles.go`
function-MI failures dropped to 25, and the current absolute function-MI failure
count is 545 while file-MI is now at 94 rows. This is witness trust-contract
documentation and function-MI reduction evidence, not MI closure.

Demo payload and report helper comments now document report-artifact directory
creation before partial writes, decoder-preserved JSON numeric payloads, and
typed string extraction from raw event arrays. Focused `internal/demo` tests
pass; the current absolute function-MI failure count is 545 while file-MI
remains at 94 rows. This is trust-documentation and function-MI reduction
evidence, not MI closure.

`cmd/sdp-trace/harness_cli.go` now hoists stable harness subcommand and required
flag tables out of the CLI handlers while preserving the `observe`, `validate`,
and `summarize` command contracts. Focused `cmd/sdp-trace` tests and strict CRAP
pass; the current absolute function-MI failure count is 545 while file-MI
remains at 94 rows. This is CLI cleanup and function-debt reduction evidence,
not MI closure.

`cmd/sdp-trace` now also hoists stable `pr-review`, `packet`, and `observe`
subcommand/required-flag tables out of CLI handlers, and
`requirePRReviewPacketInputs` uses the shared required-flag shape instead of an
ad hoc dynamic-value map. Focused `cmd/sdp-trace` tests and strict CRAP pass;
the current absolute function-MI failure count is 545 while file-MI is now at 94 rows. This is CLI contract cleanup and function-debt reduction evidence, not
MI closure.

Remaining release-proof, interaction, envelope, and checkpoint parser functions
now use stable required-flag tables instead of inline literals. Focused
`cmd/sdp-trace` tests and strict CRAP pass; the current absolute function-MI
failure count remains 545 while file-MI is now at 94 rows. This is
CleanCode/DX consistency evidence, not MI closure.

`cmd/sdp-trace.harnessValidationExitCode` now checks documented
`cannot_verify`-class states through a state set instead of repeated string
comparisons, lowering the function from cyclo/cognitive 4/3 to 3/2 while
preserving unknown-state warnings. Focused `cmd/sdp-trace` tests and strict CRAP
pass; the current absolute function-MI failure count remains 545 while file-MI
remains at 94 rows. This is complexity reduction evidence, not MI closure.

GitHub artifact CLI helpers now document token fallback, Enterprise API URL
precedence, URL credential rejection, loopback HTTP test scope,
public-vs-Enterprise host binding, and empty server URL handling. Focused
`cmd/sdp-trace` tests and strict CRAP pass; the current absolute function-MI
failure count is 545 while file-MI is now at 94 rows. This is security/DX
documentation and function-debt reduction evidence, not MI closure.

`internal/posture.groupingKeys` now uses a stable grouping-key map, and
digest-error reason selection uses an ordered table while preserving
`errors.Is` matching. Focused `internal/posture` tests and strict CRAP pass
after the repair; the current absolute function-MI failure count is 545 while
file-MI is now at 94 rows. This is strict-CRAP repair and function-debt
reduction evidence, not MI closure.

`internal/adaptercapture` now represents closed vocabularies for insufficient
event-family states, insufficient retention modes, sensitive event types, valid
retention modes, and digest-only-valid event types as package-level sets instead
of repeated switches. Focused `internal/adaptercapture` tests and strict CRAP
pass; the current absolute function-MI failure count is 545 while file-MI
remains at 94 rows. This is trust-vocabulary cleanup and function-debt reduction
evidence, not MI closure.

`internal/forensic.validRetentionMode` now uses a stable retention-mode set that
matches the forensic retention vocabulary. Focused `internal/forensic` tests and
strict CRAP pass; the current absolute function-MI failure count is 545 while
file-MI is now at 94 rows. This is trust-vocabulary cleanup and function-debt
reduction evidence, not MI closure.

`internal/ciartifact` now represents artifact family, producer-scope,
access-state, and binding-state vocabularies as package-level sets, with
fallback downgrades documented for unknown values. Focused `internal/ciartifact`
tests and strict CRAP pass; the current absolute function-MI failure count is
545 while file-MI is now at 94 rows. This is trust-vocabulary cleanup and
function-debt reduction evidence, not MI closure.

`internal/harnessobs` now represents digest-field, event-family, content-state,
validation-state, and rule-key schema vocabularies as package-level sets, with
terminal digest-field matching documented for nested raw-event paths. Focused
`internal/harnessobs` tests and strict CRAP pass; the current absolute
function-MI failure count is 545 while file-MI is now at 94 rows. This is
schema-vocabulary cleanup and function-debt reduction evidence, not MI closure.

Release-proof and CLI helper comments now document manifest-path preservation on
missing source-commit artifacts, digest mismatch bucketing, optional PR review
ledger and JSON-output behavior, protected input status classes, protected
profile exit precedence, customer-PKI witness flag scope, GitHub event-path and
HTTPS boundaries, query-pack explain contract checks, and safe flag defaults.
Focused `internal/releaseproof` and `cmd/sdp-trace` tests pass; the current
absolute function-MI failure count is 545 while file-MI is now at 94 rows.
This is trust/DX documentation and function-MI reduction evidence, not MI
closure.

Trace canonical replay comments now document null predecessor sentinels,
manifest/event-count binding, canonical scalar and boolean emission, event-hash
algorithm defaults, required trace identifier validation, run-directory replay
scope, event-hash mismatch evidence, public hash wrapper reuse, and argv digest
encoding. Focused `internal/trace` tests pass; the current absolute function-MI
failure count is 545 while file-MI is now at 94 rows. This is trace replay
documentation and function-MI reduction evidence, not MI closure.

Verifier replay comments now document optional final-chain-head binding,
zero-based event sequence checks, genesis previous-hash sentinels,
payload-digest trust failures, optional integrity-audit emission, and optional
explain rows for closure state and contract path context. Focused
`internal/verifier` tests pass; the current absolute function-MI failure count
is 545 while file-MI is now at 94 rows. This is verifier trust-documentation
and function-MI reduction evidence, not MI closure.

Interaction validation comments now document transcript task binding, actor ID
fallback, primary and optional ID safety, SHA-256 content digest binding, local
content-ref format, reference/LLM ref ordering, and envelope identity scope.
Focused `internal/interaction` tests pass; the current absolute function-MI
failure count is 545 while file-MI is now at 94 rows. This is interaction
evidence-contract documentation and function-MI reduction evidence, not MI
closure.

Checkpoint, CI artifact, adapter-capture, and managed-boundary helper comments
now document canonical replay format restrictions, checkpoint sequence-link
rules, checkpoint public-key fallback, source/nonce binding boundaries,
deterministic envelope error selection, CI artifact access/source/run
sanitization, adapter event contract and run-binding checks, test provenance
execution boundaries, managed adapter
authorization and capability coverage, managed witness binding, pre-run
provenance authority, artifact path binding, top-level state coercion, and
managed boundary connection/bypass/event-source precedence. Focused package tests
pass; `internal/checkpoint` and `internal/managed` function-MI failures are now
36 each, and the current absolute function-MI failure count is 545 while
file-MI is now at 94 rows. This is trust-boundary documentation and function-MI
reduction evidence, not MI closure.

Query and authority helper comments now document stable query row ordering,
retention-limited reconstructability, safe token fallback, selected authority
envelope validation, authority remediation fallback, authority required-list
policy domains, and contract non-empty gate lists. Focused `internal/query`,
`internal/authority`, `internal/policy`, and `internal/contract` tests pass; the
current absolute MI gap remains open at 519 function failure rows and 91 file
failure rows. This is trust-boundary documentation evidence, not MI closure.

Witness, telemetry, recorder, and forensic helper comments now document
CI-witnessed pass requirements, environment replay shape, telemetry profile
dispatch, recorder empty-array persistence, withholding audit facts, and
retention-mode verdict mapping. Focused `internal/witness`,
`internal/telemetry`, `internal/recorder`, and `internal/forensic` tests pass;
the current absolute MI gap remains open at 519 function failure rows and 93
file failure rows. This is trust-boundary documentation evidence, not MI
closure.

Packet validation comments now document resolver provenance anchors, separated
row trust-contract checks, decision-owner reason requirements, unverifiable pass
refs, and prompt-boundary failure precedence. Focused `internal/packet` tests
pass; the current absolute MI gap remains open at 519 function failure rows and
91 file failure rows. This is packet trust-boundary documentation evidence, not
MI closure.

Harness observation comments now document output-directory write-target checks,
raw-event source fallback limits, isolation-rule replay verification, tool-family
signal classification, structured-before-leaf raw signal extraction,
signal-value key gating, optional family vocabulary validation, missing output
directory handling, safe model argument capture, and encoded-token detection.
Focused `internal/harnessobs` tests pass; the current absolute MI gap remains
open at 519 function failure rows and 91 file failure rows. This is harness
observation trust-boundary documentation evidence, not MI closure.

Posture trust-documentation comments now document selection replay boundaries,
grouping validation, handoff shape normalization, refusal-versus-metric
aggregation, digest-before-read proof checks, optional signal semantics,
deterministic grouping/export/movement ordering, required collection shape,
closed metric and dimension vocabularies, sensitive-class absence claims,
profile validation, manifest schema drift, movement comparability, and label/path
safety. Focused `internal/posture` tests pass; the current absolute MI gap
remains open at 519 function failure rows and 91 file failure rows. This is
posture trust-documentation and function-MI reduction evidence, not MI closure.

Additional harness observation collection comments now document event-before-run
persistence, input validation before output creation, read-only source versus
write-target policy, profile-bound event parsing, portable source basenames,
explicit cannot-verify collection state, session/profile identity checks,
profile-relative harness paths, normalization re-resolution, missing-source
verdict persistence, observed-run validation profile, command end-time fallback,
source commit cannot-verify semantics, normalized event digest binding, raw-field
rejection before family expansion, and additive OpenCode family classification.
Focused `internal/harnessobs` tests pass; the current absolute MI gap remains
open at 519 function failure rows and 91 file failure rows. This is harness
observation trust-documentation and function-MI reduction evidence, not MI
closure.

Packet trust-contract comments now document draft/non-approval packet generation,
prompt-contamination theater findings, retained-text versus digest-only prompt
classification, packet digest binding, demo first-packet evidence requirements,
manifest ref usability, route-observation component proof, schema and bundle
identity binding, resolver provenance, cross-row contradiction handling, residual
gap coverage, theater state lowering, pass-row artifact access, GitHub
verification pass requirements, retained artifact deduplication,
prompt-boundary cannot-verify effects, stable entry ordering, required
prompt-boundary manifest refs, real route digests, artifact access metadata,
authority-entry limits, and synthetic digest placeholders. Focused
`internal/packet` tests pass; the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is packet
trust-documentation and function-MI reduction evidence, not MI closure.

PR review evidence-boundary comments now document packet immutability, digest
binding, unavailable-field evidence, review-record authority limits, stale
packet detection, required plane isolation, model fallback provenance, runner
allow-listing, OpenCode read-only mutation checks, raw output retention, prompt
digest-only refs, citation resolvability, command-runner setup authority,
OpenCode baseline scope, context classification limits, path-safe ID
normalization, and reviewer text sanitization.
Focused `internal/prreview` tests pass; `internal/prreview` function-MI failures
dropped from 100 to 63, and the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is PR-review
trust-documentation and function-MI reduction evidence, not MI closure.

Demo gate trust-boundary comments now document live report and gate regeneration,
contract-bound evidence classification, local gate derivation, protected-profile
lowering, witness binding and freshness, checkpoint signature/binding/signing
policy, required-run matching, override non-upgrade behavior, timeline escaping,
and stable report artifact writes. Focused `internal/demo` tests pass;
`internal/demo` function-MI failures dropped from 97 to 50, and the current
absolute MI gap remains open at 519 function failure rows and 91 file failure
rows. This is demo trust-documentation and function-MI reduction evidence, not
MI closure.

Interaction trace trust-boundary comments now document relay capture before
forwarding, transcript import source restrictions, event/source/catalog
validation, content digest and blob binding, trace and envelope read validation,
summary not-assessed propagation, JSONL streaming limits, per-source ordering,
safe ID/reference limits, and deterministic event-type help output. Focused
`internal/interaction` tests pass; `internal/interaction` function-MI failures
dropped from 68 to 25, and the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is interaction
trust-documentation and function-MI reduction evidence, not MI closure.

Repo-observer install/proof boundary comments now document read-only doctor
behavior, install preview/write rescans, optional JSON sinks, local config
safety, human table state separation, repo-root containment, hook/workflow and
artifact proof limits, monotonic state aggregation, force-mode backups and safe
summaries, gitignore managed-block handling, generated target portability,
GitHub workflow observation limits, and origin credential redaction. Focused
`internal/repoobserver` tests pass; `internal/repoobserver` function-MI failures
dropped from 64 to 18, and the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is repo-observer trust/DX
documentation and function-MI reduction evidence, not MI closure.

Witness trust-profile boundary comments now document profile publication safety,
envelope artifact comparison, missing-envelope CI downgrade semantics, unsafe
envelope rejection, Customer PKI authority/freshness/signature ordering,
revocation assessment gaps, output-safety redaction, run-ID binding, secret/JWT
detection, public-key parsing, and GitHub OIDC request/claim boundaries. Focused
`internal/witness` tests pass; `internal/witness` function-MI failures dropped
from 89 to 33, and the current absolute MI gap remains open at 519 function
failure rows and 91 file failure rows. This is witness trust-documentation and
function-MI reduction evidence, not MI closure.

Authority evaluation boundary comments now document package parse authority
limits, selected-envelope ambiguity, policy/event validation ordering,
target-rule conflicts, deterministic action evaluation, pre-decision blockers,
task scope, approval evidence, binding endpoints, gateway model attribution,
evidence-ref resolution, aggregate state ranking, diagnostic tool/model
attribution gaps, deterministic unique summaries, source coverage, and safe
reference filtering. Focused `internal/authority` tests pass;
`internal/authority` function-MI failures dropped from 49 to 5, and the current
absolute MI gap remains open at 519 function failure rows and 91 file failure
rows. This is authority trust-documentation and function-MI reduction evidence,
not MI closure.

Checkpoint replay/signature boundary comments now document local key generation
limits, signing versus signer authority, payload replay binding, run nonce
authority, envelope shape failures, sequence-link rules, set verification
ordering, digest/signature verification, policy signer binding,
authority-to-trust-scope mapping, worst-state aggregation, final result
derivation, canonicalization restrictions, and key decoding size checks. Focused
`internal/checkpoint` tests pass; `internal/checkpoint` function-MI failures
dropped from 36 to 7, and the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is checkpoint
trust-documentation and function-MI reduction evidence, not MI closure.

CI artifact observation boundary comments now document source/run identity
sanitization before proof output, closed artifact-family vocabulary,
required-family ordering, extra-family provenance limits, family evaluation
order, CI-uploaded producer authority, artifact index and output-safety
separation, top-level state derivation, binding summaries, aggregate
producer/access states, required producer defaults, safe source/run token
handling, unsafe identity marker rejection, safe class filtering, default
safety-ruleset hashing, and deterministic reason/action ordering. Focused
`internal/ciartifact` tests pass; `internal/ciartifact` function-MI failures
dropped from 34 to 8, and the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is CI artifact
trust-documentation and function-MI reduction evidence, not MI closure.

Adapter capture evidence-boundary comments now document condition-first verdict
derivation, local trust-scope lowering, adapter event contract failures,
correlation-key ambiguity, identity binding classification, run id/nonce replay
tuple, same-chain and bundle binding, task supersession attribution, missing
versus unsupported tool telemetry, file mutation correlation, gateway model
identity, test provenance overclaims, provider-ref secret scanning, redaction
metadata failures, capture-depth caps, top-level state lowering, secret markers,
and representative fixture construction. Focused `internal/adaptercapture`
tests pass; `internal/adaptercapture` function-MI failures dropped from 36 to
7, and the current absolute MI gap remains open at 519 function failure rows
and 91 file failure rows. This is adapter-capture trust-documentation and
function-MI reduction evidence, not MI closure.

Forensic retention evidence-boundary comments now document condition-row
verdict derivation, policy binding layers, policy contract replay facets,
self-claimed authority limits, run/event policy contradictions, prewrite rule
indexing, persisted secret-like value failure semantics, redaction digest and
rule-reference requirements, withholding authority, retention-mode policy
membership, critical-event reconstructability, raw-reference validation,
profile selection accountability, downgrade completeness, top-level failure
precedence, deterministic reason ordering, and representative fixture scope.
Focused `internal/forensic` tests pass; `internal/forensic` function-MI
failures dropped from 36 to 17, and the current absolute MI gap remains open at
519 function failure rows and 91 file failure rows. This is forensic
trust-documentation and function-MI reduction evidence, not MI closure.

Recorder evidence-boundary comments now document pre-execution trace setup,
explicit environment inheritance, contract digest binding, initial manifest
authority, command-failure closure semantics, fresh output-directory isolation,
run-start ordering, command start/finish evidence scope, run-closed head
binding, default-contract opt-in, digest-only stdout/stderr retention, source
snapshot limits, fixed run layout, writer-owned sequence/hash assignment,
durable event advancement, manifest head mirroring, process start/cancel/signal
classification, JSON object payload replay, portable command basenames,
fallback id limits, and pretty JSON artifact reviewability. Focused
`internal/recorder` tests pass; `internal/recorder` function-MI failures
dropped from 30 to 8, and the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is recorder
trust-documentation and function-MI reduction evidence, not MI closure.

CLI CRAP hotspot tests now exercise interaction import/summarize and relay-arg
failure paths, optional JSON output, managed assessment input failures, doctor
and install usage failures, forensic and CI artifact exit-code mappings, demo
explanation helpers, optional witness matching, protected preview actions, gate
state exit selection, fixture expectation helpers, telemetry export input
validation, unsupported witness kind handling, flag accessors, help detection,
and boolean literal parsing. Focused `cmd/sdp-trace` tests pass and the full
strict CRAP replay passes; the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is test coverage and
function-MI reduction evidence, not MI closure.

Query-pack evidence-boundary comments now document required-run artifact
authority, optional forensic and adapter artifact error handling, safe
explanation output, digest binding before JSON decode, stable empty query-group
serialization, optional input inclusion, output-safety scope, fallback timeline
rows, optional-block absence versus malformed evidence, Block 18/19 row source
refs, deterministic verifier-state ordering, general gap semantics,
unverified-claim mirroring, summary-row indexing, retention-limited critical
evidence caps, per-query row id generation, unknown state lowering, condition
family classification, rule-order classification, token sanitization, and
sensitive-class absence scope. Focused `internal/query` tests pass;
`internal/query/querypack.go` function-MI failures dropped from 33 to 9, and
the current absolute MI gap remains open at 519 function failure rows and 93
file failure rows. This is query-pack trust-documentation and function-MI
reduction evidence, not MI closure.

CLI PR-review command comments now document packet replay boundaries, runner
allow-list provenance, preview versus evidence output, synthesis and validation
artifact ordering, write-once summary behavior, combined check sequencing,
profile coverage requirements, and required packet provenance anchors. Focused
`cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go` function-MI failures
dropped from 322 to 310, and the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is CLI trust-boundary
documentation and function-MI reduction evidence, not MI closure.

CLI release-proof and interaction/envelope command comments now document
source-bound release proof evaluation, proof persistence requirements, relay
capture-before-forwarding, transcript import attribution, summary derived-view
limits, and envelope summarize command parsing. Focused `cmd/sdp-trace` tests
pass; `cmd/sdp-trace/main.go` function-MI failures dropped from 310 to 300, and
the current absolute MI gap remains open at 519 function failure rows and 93
file failure rows. This is CLI trust-boundary documentation and function-MI
reduction evidence, not MI closure.

CLI assessment command comments now document profile-specific evidence shapes,
required input boundaries, local JSON loading authority, verdict persistence,
preview non-evidence semantics, and authority-package assessment boundaries.
Focused `cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go` function-MI
failures dropped from 300 to 292, and the current absolute MI gap remains open
at 519 function failure rows and 91 file failure rows. This is CLI
trust-boundary documentation and function-MI reduction evidence, not MI closure.

CLI assessment preview comments now document that preview reports expose setup
readiness and expected vocabularies without evaluating raw payloads, fetching
artifacts, or emitting verdicts. Focused `cmd/sdp-trace` tests pass;
`cmd/sdp-trace/main.go` function-MI failures dropped from 292 to 290, and the
current absolute MI gap remains open at 519 function failure rows and 91 file
failure rows. This is CLI preview trust-language and function-MI reduction
evidence, not MI closure.

CLI assessment explanation comments now document that explanation renders
existing assessment artifacts only, uses schema-version dispatch instead of
profile-name dispatch, keeps unknown schemas as `cannot_verify`, treats typed
artifact loading as the trust boundary, preserves retention caps beside their
conditions, and keeps CI-artifact producer/access/binding/output-safety plus
authority action/provenance gaps distinct. Focused `cmd/sdp-trace` tests pass;
`cmd/sdp-trace/main.go` function-MI failures dropped from 290 to 281, and the
current absolute MI gap remains open at 519 function failure rows and 91 file
failure rows. This is CLI explanation trust-language and function-MI reduction
evidence, not MI closure.

CLI checkpoint comments now document the create-versus-verify command boundary,
subcommand flag ownership, local signing-key limits, signed artifact
persistence, verification input loading order, optional policy semantics,
required replay inputs, and checkpoint pass/cannot-verify exit-code mapping.
Focused `cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go` function-MI
failures dropped from 281 to 272, and the current absolute MI gap remains open
at 519 function failure rows and 91 file failure rows. This is CLI checkpoint
trust-language and function-MI reduction evidence, not MI closure.

CLI report/gate comments now document report artifact authority, cannot-verify
exit precedence, single-target report/gate provenance, protected-gate external
trust input requirements, protected checkpoint replay scope, witness expectation
derivation, persisted gate result authority, unsupported gate-result schema
handling, and layered gate explanation output. Focused `cmd/sdp-trace` tests
pass; `cmd/sdp-trace/main.go` function-MI failures dropped from 272 to 260, and
the current absolute MI gap remains open at 519 function failure rows and 93
file failure rows. This is CLI report/gate trust-language and function-MI
reduction evidence, not MI closure.

CLI gate preview/protected helper comments now document preview non-verdict
semantics, target-scoped preview provenance, optional witness binding preview,
protected setup cannot-verify handling, single-run protected replay, protected
trust upgrade prerequisites, signer-policy binding, witness source and artifact
matching, normalized expected artifact digests, wildcard source fields, and
run-derived witness expectations. Focused `cmd/sdp-trace` tests pass;
`cmd/sdp-trace/main.go` function-MI failures dropped from 260 to 250, and the
current absolute MI gap remains open at 519 function failure rows and 91 file
failure rows. This is CLI gate-preview trust-language and function-MI reduction
evidence, not MI closure.

CLI override/shared persistence comments now document run-artifact digest
authority for witness expectations, protected input setup gaps, deterministic
protected-preview remediation, override request non-upgrade semantics, external
reference limits, single override write action parsing, required override
provenance fields, shared JSON artifact load responsibilities, reviewable JSON
output, and atomic text artifact writes. Focused `cmd/sdp-trace` tests pass;
`cmd/sdp-trace/main.go` function-MI failures dropped from 250 to 242, and the
current absolute MI gap remains open at 519 function failure rows and 91 file
failure rows. This is CLI override/persistence trust-language and function-MI
reduction evidence, not MI closure.

CLI witness/shared helper comments now document atomic text publication,
temporary-file cleanup on write failure, gate-preview mode precedence, omission
of empty preview identifiers, protected versus non-protected exit-code
ownership, required-run participation in exit decisions, cannot-verify exit
semantics, explicit witness provenance fields, and required witness field
normalization. Focused `cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go`
function-MI failures dropped from 242 to 231, and the current absolute MI gap
remains open at 519 function failure rows and 91 file failure rows. This is
CLI witness/helper trust-language and function-MI reduction evidence, not MI
closure.

CLI witness option/output comments now document one-target witness provenance,
required output persistence, kind-specific customer-PKI validation, optional
witness field copying, closed witness kind semantics, builder ownership of
evidence interpretation, customer-PKI authority/credential/freshness inputs,
witness stdout as a rendered copy, cannot-verify/not-assessed exit behavior,
deterministic customer-PKI remediation, and allowed-kind builder contracts.
Focused `cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go` function-MI
failures dropped from 231 to 218, and the current absolute MI gap remains open
at 519 function failure rows and 91 file failure rows. This is CLI witness
trust-language and function-MI reduction evidence, not MI closure.

CLI run/wrap/preview comments now document default-contract legacy wrap
behavior, recorder-owned artifact/event writes, child-command provenance,
run-directory output authority, task-to-SpecKit binding, explicit contract
selection, preview non-write behavior, command descriptor requirements, default
contract replayability, and malformed preview-contract `cannot_verify`
semantics. Focused `cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go`
function-MI failures dropped from 218 to 213, and the current absolute MI gap
remains open at 519 function failure rows and 91 file failure rows. This is
CLI run/preview trust-language and function-MI reduction evidence, not MI
closure.

CLI doctor/install comments now document profile-scoped doctor delegation,
flag-only doctor determinism, supported repo-observer profile scope, read-only
doctor behavior, persisted doctor/install JSON artifacts, local doctor scope,
preview/write install handling, partial install status rendering,
repo-observer-only install parsing, flag-only install semantics, default
portable profile selection, shared repoobserver option conversion, and
cannot-verify install/proof exit behavior. Focused `cmd/sdp-trace` tests pass;
`cmd/sdp-trace/main.go` function-MI failures dropped from 213 to 206, and the
current absolute MI gap remains open at 519 function failure rows and 91 file
failure rows. This is CLI doctor/install trust-language and function-MI
reduction evidence, not MI closure.

CLI verify/query comments now document verifier artifact persistence even on
semantic verification errors, structured result output before diagnostics,
retained-run argument authority, missing run roots as `cannot_verify`, future
verdict fallback behavior, read-only explanation semantics, query run-source
requirements, unsupported query usage errors, query replay failures as
`cannot_verify`, query-pack build versus explain separation, and portable
query-pack artifact writing. Focused `cmd/sdp-trace` tests pass;
`cmd/sdp-trace/main.go` function-MI failures dropped from 206 to 198, and the
current absolute MI gap remains open at 519 function failure rows and 91 file
failure rows. This is CLI verify/query trust-language and function-MI reduction
evidence, not MI closure.

CLI query-pack/export comments now document explanation as rendering existing
artifacts only, explicit provenance flags for pack/run/output inputs, closed
pack and export command vocabularies, source-run and durable-output
requirements, query-pack read failures as verification failures, telemetry
rendering from an existing cross-repo posture artifact, and stdout/file output
semantics. Focused `cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go`
function-MI failures dropped from 198 to 192, and the current absolute MI gap
remains open at 519 function failure rows and 91 file failure rows. This is CLI
query-pack/export trust-language and function-MI reduction evidence, not MI
closure.

CLI cross-repo posture export comments now document selection-file authority,
artifact-path authority over stdout, flag-only provenance, validate-only
non-publication semantics, durable output requirements, supported posture
profile scope, explanation from saved artifacts only, output-safety failure
semantics, and schema/profile checks before explanation. Focused
`cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go` function-MI failures
dropped from 192 to 188, and the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is CLI cross-repo posture
trust-language and function-MI reduction evidence, not MI closure.

CLI fixture-validation comments now document fixture-root discovery scope,
continuing through all fixture runs, verifier artifact persistence during
semantic replay errors, structured-verdict comparison against fixture
expectations, invalid expectation metadata as drift, explicit expected-verdict
authority, and default failure handling for fail/cannot-verify results. Focused
`cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go` function-MI failures
dropped from 188 to 183, and the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is CLI fixture-validation
trust-language and function-MI reduction evidence, not MI closure.

CLI doctor trust-language comments now document local-readiness scope without
CI/external witness upgrades, embedded versus requested contract loading,
cannot-verify exit lowering, durable artifact path probes, non-directory and
missing-root write-target handling, contract evidence-reference support checks,
stable local event vocabulary, CI witness environment prerequisites, safe
retention modes, preview boundary classification, and offline evidence
upgrade requirements. Focused `cmd/sdp-trace` tests pass;
`cmd/sdp-trace/main.go` function-MI failures dropped from 183 to 175, and the
current absolute MI gap remains open at 519 function failure rows and 91 file
failure rows. This is CLI doctor trust-language and function-MI reduction
evidence, not MI closure.

CLI usage, fixture expectation, and flag parser comments now document global
help as the local command contract, optional expectation metadata, fixture
expectation basename keys, missing expectation defaults, malformed expectation
errors, parser-owned index advancement, `--` command payload handling,
positional argument preservation, direct and next-argument string assignment,
bare boolean semantics, and invalid boolean usage errors. Focused
`cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go` function-MI failures
dropped from 175 to 165, and the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is CLI parser
trust-language and function-MI reduction evidence, not MI closure.

Observe and harness CLI comments now document setup versus collection scope,
flag-only profile/source/output provenance, retained run evidence binding,
combined session/run output shape, command payload forwarding after `--`,
validation-artifact loading, validation-state exit mapping, and persisted
validation evidence for summaries. Focused `cmd/sdp-trace` tests pass;
`cmd/sdp-trace/observe_cli.go` function-MI failures dropped from 8 to 6,
`cmd/sdp-trace/harness_cli.go` function-MI failures dropped from 6 to 5, and
the current absolute MI gap remains open at 519 function failure rows and 93
file failure rows. This is observe/harness CLI trust-language and function-MI
reduction evidence, not MI closure.

Managed-harness evidence-boundary comments now document fixed condition-set
verdict derivation, trust-scope lowering, managed policy and registry pre-run
authority, boundary enrollment ordering and binding, adapter identity and
capability coverage, authorized capability event coverage, activation source,
required event group semantics, suppression absence versus policy-valid
suppression, test provenance authority, witness mismatch failure semantics,
witness pass/freshness and artifact requirements, independent witness binding
checks, override non-upgrade behavior, adapter authorization selection,
policy-group scope matching, suppression rule verification, artifact digest
matching, severity ordering, reason ordering, next-action deduplication, and
condition-id tie-breaking. Focused `internal/managed` tests pass;
`internal/managed/managed.go` function-MI failures dropped from 36 to 4, and
the current absolute MI gap remains open at 519 function failure rows and 93
file failure rows. This is managed-harness trust-documentation and function-MI
reduction evidence, not MI closure.

Verifier replay-boundary comments now document manifest and event loading
failure semantics, provisional observed results, chain hard-fail behavior,
default-contract digest drift, explicit contract loading, replay lowering,
missing-evidence `not_assessed` behavior, event-type coverage, integrity-audit
detail scope, event-file filtering, portable parse errors, event-level defect
precedence, sequence and previous-hash checks, payload-digest and event-hash
authority, verifier artifact separation, and human-reviewable JSON output.
Focused `internal/verifier` tests pass; `internal/verifier/verify.go`
function-MI failures dropped from 26 to 5, and the current absolute MI gap
remains open at 519 function failure rows and 91 file failure rows. This is
verifier trust-documentation and function-MI reduction evidence, not MI
closure.

Telemetry export-boundary comments now document validated posture input
requirements, empty Prometheus output shape, deterministic series ordering,
fixed HELP/TYPE authority, no-hidden-score metric gauges, movement
comparability, unsafe-label rejection before rendering, cardinality limits,
refusal and input aggregation semantics, public dimension filtering, label map
copying, closed label vocabulary, label value safety checks, duplicate series
rejection, sorted label rendering, empty-label omission, and Prometheus
label-value escaping. Focused `internal/telemetry` tests pass;
`internal/telemetry/prometheus.go` function-MI failures dropped from 25 to 11,
and the current absolute MI gap remains open at 519 function failure rows and
91 file failure rows. This is telemetry trust-documentation and function-MI
reduction evidence, not MI closure.

Interaction trust-boundary comments now document live relay evidence authority,
relay-before-forward replay semantics, transcript-import provenance limits,
observed event retention and channel-exclusivity gaps, trace completeness
folding, validation catalog separation, trace and envelope summary limits,
bounded stdin/JSONL import, exact forwarded-byte hashing, safe actor/source
catalogs, and retained content-reference semantics. Focused
`internal/interaction` tests pass; `internal/interaction/interaction.go`
function-MI failures dropped from 25 to 3, and the current absolute MI gap
remains open at 519 function failure rows and 91 file failure rows. This is
interaction trust-documentation and function-MI reduction evidence, not MI
closure.

Trace replay-boundary comments now document canonical JSON hash authority,
numeric normalization, event-hash self-exclusion, payload digest requirements,
run-layout materialization, manifest/event-chain validation separation,
append-only event sequence and chain-head updates, artifact-copy durability,
contract defaulting, missing-evidence table semantics, safe command-descriptor
retention, event defaulting, payload representation synchronization, and local
source-snapshot disclosure. Focused `internal/trace` tests pass;
`internal/trace` function-MI failures dropped from 61 to 0, and the current
absolute MI gap remains open at 519 function failure rows and 91 file failure
rows. This is trace trust-documentation and function-MI reduction evidence, not
MI closure.

Witness trust-boundary comments now document GitHub Actions upgrade
preconditions, optional report-artifact binding, OIDC request host/audience
pinning, JWT claim retention limits, run/report digest path semantics, Customer
PKI external-trust prerequisites, envelope safety stop points, Customer PKI
authority issue precedence, CI-envelope state repair order, output-safety
replacement records, run-ID binding, secret/JWT deny-list scope, and envelope
scalar safety. Focused `internal/witness` tests pass; `internal/witness`
function-MI failures dropped from 33 to 12, and the current absolute MI gap
remains open at 519 function failure rows and 91 file failure rows. This is
witness trust-documentation and function-MI reduction evidence, not MI closure.

Repo-observer DX/trust-boundary comments now document optional JSON status
sinks, local repository-root evidence, human table install/proof separation,
safe diff summaries, snapshot-only status generation, surface ordering, CI
artifact upload versus inspected bundle gaps, deterministic next actions,
pre-write hooksPath validation, safe target containment, force-mode backups,
managed `.gitignore` replacement, generated config manifest scope, and workflow
artifact-upload limits. Focused `internal/repoobserver` tests pass;
`internal/repoobserver` function-MI failures dropped from 18 to 3, and the
current absolute MI gap remains open at 519 function failure rows and 91 file
failure rows. This is repo-observer trust/DX documentation and function-MI
reduction evidence, not MI closure.

Forensic retention trust-boundary cleanup now documents condition aggregation,
policy and prewrite precedence, event-policy binding, retention-mode
vocabulary, critical evidence reconstruction, raw-reference access, key
custody, lifecycle checks, profile-selection accountability, and fixture
evidence shape. `ValidTestInput` fixture construction is split into small
same-package helpers with the same representative policy/run semantics.
Focused `internal/forensic` tests pass; `internal/forensic` function-MI
failures dropped from 17 to 0, and the current absolute MI gap remains open at
519 function failure rows and 91 file failure rows. This is forensic
trust-documentation and function-MI reduction evidence, not MI closure.

Packet CLI helper extraction moved basic `packet build-github`,
`packet validate`, `packet check-demo`, and `packet render` helpers from
`cmd/sdp-trace/main.go` into same-package `cmd/sdp-trace/packet_cli.go`.
Trust-boundary comments now document durable packet publication,
committed-bundle validation, demo first-packet checks, markdown persistence,
and flag-only packet input parsing. Focused packet CLI tests pass;
`cmd/sdp-trace/main.go` function-MI failures dropped from 132 to 92 after
the follow-on export/doctor comment pass, and `packet_cli.go` has 0
function-MI failures. The current absolute MI gap remains open at 519 function
failure rows and 91 file failure rows because `packet_cli.go` is still below
file-MI 70. This is packet CLI layout and trust-documentation evidence, not MI
closure.

CLI export/doctor comments now document telemetry profile locking, posture
artifact authority, explicit stdout/file export intent, cross-repo posture
artifact provenance, explanation non-evidence scope, local doctor write probes,
control-point exit lowering, CI witness prerequisites, and external witness
non-integration. Focused `cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go`
function-MI failures dropped from 108 to 105, and the current absolute MI gap
remains open at 519 function failure rows and 91 file failure rows. This is CLI
trust-documentation and function-MI reduction evidence, not MI closure.

Parallel posture and packet worker slices added trust-boundary comments in
`internal/posture/posture.go` and `internal/packet/packet.go`. The packet slice
also extracted prompt-boundary classification helpers without changing tested
behavior. Focused `internal/posture`, `internal/packet`, and `cmd/sdp-trace`
tests pass after integration; `internal/posture/posture.go` function-MI
failures are now 71, and `internal/packet/packet.go` function-MI failures are
now 54. The current absolute MI gap remains open at 519 function failure rows
and 91 file failure rows. This is parallel trust-documentation and packet
helper extraction evidence, not MI closure.

CLI assess/interaction comments now document interaction verb dispatch, relay
command-boundary preservation, parser usage-code handling, optional summary
publication failures, flag-only assess replayability, explicit evidence inputs
for adapter/managed/forensic/CI/authority assessments, and preview
non-evidence semantics. Focused `cmd/sdp-trace` tests pass;
`cmd/sdp-trace/main.go` function-MI failures dropped from 105 to 92, and the
current absolute MI gap remains open at 519 function failure rows and 91 file
failure rows. This is CLI trust-documentation and function-MI reduction
evidence, not MI closure.

Parallel harness and PR-review worker slices added trust-boundary comments in
`internal/harnessobs/harnessobs.go` and `internal/prreview/prreview.go`.
Focused `internal/harnessobs`, `internal/prreview`, and `cmd/sdp-trace` tests
pass after integration; `internal/harnessobs/harnessobs.go` function-MI
failures are now 75, and `internal/prreview/prreview.go` function-MI failures
are now 47. The current absolute MI gap remains open at 519 function failure
rows and 91 file failure rows. This is parallel trust-documentation evidence,
not MI closure.

Harness near-threshold trust comments now document observed run index
publication, setup-created session collection, normalized signal comparison,
run index validation before event refs, harness family identity constraints,
output path joining after safety checks, symlink target normalization, and
immutable source commit retention. Focused `internal/harnessobs` tests pass;
`internal/harnessobs/harnessobs.go` function-MI failures dropped from 83 to 75,
and the current absolute MI gap remains open at 519 function failure rows and
91 file failure rows. This is harness observation trust-documentation and
function-MI reduction evidence, not MI closure.

Posture validation/catalog trust comments now document required export
collections, closed dimension/trust/source/comparison/refusal vocabularies,
digest manifest authority, artifact byte hashing, and output-safety rejection
for URLs, credentials, and paths. Focused `internal/posture` tests pass;
`internal/posture/posture.go` function-MI failures dropped from 71 to 62, and
the current absolute MI gap remains open at 519 function failure rows and 94
file failure rows. This is posture validation trust-documentation and
function-MI reduction evidence, not MI closure.

Packet validation/render trust comments now document demo first-packet
validation accumulation, row-specific closure evidence, case-insensitive
observed component matching, manifest-only evidence refs, GitHub source-change
derivation, prompt-boundary route cannot-verify reasons, and evidence bundle
rendering. Focused `internal/packet` tests pass; `internal/packet/packet.go`
function-MI failures dropped from 50 to 43, and the current absolute MI gap
remains open at 519 function failure rows and 91 file failure rows. This is
packet validation/render trust-documentation and function-MI reduction evidence,
not MI closure.

CLI packet-build trust comments now document packet build-pr input
reconstruction before artifact writes, cannot-verify JSON publication without
complete artifacts, flag-only replayability, result-path declaration before
writes, route/CI live gate readiness, optional route enrichment,
prompt-boundary requirement, retained GitHub artifact references, API host
authorization, and non-2xx artifact discovery failures. Focused
`cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go` function-MI failures
dropped from 97 to 92, and the current absolute MI gap remains open at 545
function failure rows and 91 file failure rows. This is packet-build
trust-documentation and function-MI reduction evidence, not MI closure.

Parallel demo and qualitycheck worker slices added comments in
`internal/demo/demo.go` and `tools/qualitycheck`. Focused `internal/demo`,
`tools/qualitycheck`, and `cmd/sdp-trace` tests pass after integration;
`internal/demo/demo.go` function-MI failures are now 41, and tool-level
function-MI failures are now 78. The current absolute MI gap remains open at
519 function failure rows and 91 file failure rows. This is parallel
trust/quality-tool documentation evidence, not MI closure.

CLI query/export/doctor comments now document query mode selection, raw query
JSON emission, query-pack persistence, query-pack explain schema gating, export
subcommand closed vocabulary, telemetry atomic output, posture export build
failure semantics, posture explain artifact validation, doctor contract-check
ordering, and named control-point reporting. Focused `cmd/sdp-trace` tests
pass; `cmd/sdp-trace/main.go` function-MI failures dropped from 92 to 90, and
the current absolute MI gap remains open at 519 function failure rows and 94
file failure rows. This is CLI trust-documentation and function-MI reduction
evidence, not MI closure.

Parallel posture and packet contract comments added trust-boundary notes in
`internal/posture/posture.go` and `internal/packet/packet.go`. Focused
`internal/posture`, `internal/packet`, and `cmd/sdp-trace` tests pass after
integration; `internal/posture/posture.go` function-MI failures are now 62,
and `internal/packet/packet.go` function-MI failures are now 43. The current
absolute MI gap remains open at 519 function failure rows and 91 file failure
rows. This is parallel trust-documentation evidence, not MI closure.

CLI gate/checkpoint comments now document checkpoint verb semantics,
checkpoint verification input separation, report publication failure, gate
subcommand read-only paths, protected gate input parsing and evaluation
ordering, protected row replay failures, gate explanation artifact handling,
and preview non-evidence output. Focused `cmd/sdp-trace` tests pass;
`cmd/sdp-trace/main.go` function-MI failures dropped from 91 to 90, and the
current absolute MI gap remains open at 519 function failure rows and 91 file
failure rows. This is CLI trust-documentation and function-MI reduction
evidence, not MI closure.

Parallel harness and PR-review near-threshold comments added trust-boundary
notes in `internal/harnessobs/harnessobs.go` and
`internal/prreview/prreview.go`. Focused `internal/harnessobs`,
`internal/prreview`, and `cmd/sdp-trace` tests pass after integration;
`internal/harnessobs/harnessobs.go` function-MI failures are now 67, and
`internal/prreview/prreview.go` function-MI failures are now 47. The current
absolute MI gap remains open at 519 function failure rows and 91 file failure
rows. This is parallel trust-documentation evidence, not MI closure.

PR-review execution and validation trust comments now document packet default
semantics, required-plane filtering, digest mismatch handling, unavailable
versus failed runner states, closed CI and severity vocabularies, safe copied
input extensions, and default reviewer result timestamps/model identity.
Focused `internal/prreview` tests pass; `internal/prreview/prreview.go`
function-MI failures dropped from 47 to 38. The current absolute MI gap remains
open at 519 function failure rows and 91 file failure rows. This is PR-review
trust-documentation and function-MI reduction evidence, not MI closure.

Parallel packet/posture helper extraction plus a local PR-review trust-comment
pass reduced the current function-MI failure shape materially. The packet
worker split packet construction, prompt-boundary classification, validation,
GitHub row/entry generation, and rendering helpers until `internal/packet`
reported zero function-MI failures. The posture worker reduced
`internal/posture/posture.go` from 62 to 48 failures with digest, metric,
movement, validation, and output-safety helper extraction. The local PR-review
pass reduced `internal/prreview/prreview.go` from 38 to 21 failures with
validation, run-set, prompt, and copied-input trust-boundary comments. Focused
`internal/posture`, `internal/packet`, and `internal/prreview` tests pass; the
current absolute MI gap remains open at 344 function failure rows and 91 file
failure rows. This is helper-extraction and trust-documentation evidence, not
MI closure.

Demo gate/report trust comments and CLI usage rendering cleanup reduced current
function-MI debt without changing behavior. The demo pass documented report
artifact authority, run discovery, gate construction, protected-profile
conditions, witness binding, override extraction, payload decoding, and timeline
rendering. The CLI pass moved the global usage literal to package scope so
`printUsage` remains a simple renderer. Focused `cmd/sdp-trace` and
`internal/demo` tests pass; `internal/demo/demo.go` function-MI failures dropped
from 35 to 9 and `cmd/sdp-trace/main.go` dropped from 90 to 89. The current
absolute MI gap remains open at 344 function failure rows and 91 file failure
rows. This is trust-documentation and CLI layout evidence, not MI closure.

PR-review CLI parser trust comments documented packet/synthesize/validate/check
flag parsing, stdout-as-artifact-copy behavior, malformed profile handling,
run-set persistence ordering, packet/profile load failures, repeated flag
reconstruction, and empty external-runner allow-lists. Focused
`cmd/sdp-trace` tests pass; `cmd/sdp-trace/main.go` function-MI failures dropped
from 89 to 81. The current absolute MI gap remains open at 344 function failure
rows and 91 file failure rows. This is CLI trust-documentation evidence, not MI
closure.

Parallel harness/CLI helper extraction plus a local PR-review comment pass
reduced the current function-MI failure shape again. The harness worker reduced
`internal/harnessobs/harnessobs.go` from 67 to 12 failures with observation,
provenance, path-safety, normalization, and validation helper extraction while
preserving package CRAP `< 5`. The CLI worker reduced `cmd/sdp-trace/main.go`
from 89 to 81 failures by extracting PR-review option assembly, doctor-report
assembly, and static preview boundaries while preserving focused CLI CRAP
`< 5`. The local PR-review pass reduced `internal/prreview/prreview.go` from
21 to 10 failures with packet, citation, prompt, copied-input, and sanitization
trust-boundary comments. Focused `cmd/sdp-trace`, `internal/harnessobs`, and
`internal/prreview` tests and strict focused CRAP pass; the current absolute MI
gap remains open at 344 function failure rows and 91 file failure rows. This is
parallel helper-extraction and trust-documentation evidence, not MI closure.
