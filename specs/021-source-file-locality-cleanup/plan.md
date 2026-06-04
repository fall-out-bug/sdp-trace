# Plan: Source File Locality Cleanup

Status: in_progress

## Workstreams

### WS-021-A: Slice Selection

Owned files:

- `cmd/sdp-trace` command-surface registry helpers for Slice 1

Deliverable:

- Group the selected command-surface helper shards into
  `cmd/sdp-trace/command_surface_registry_helpers.go` and
  `cmd/sdp-trace/command_surface_metadata_helpers.go`.

### WS-021-B: Behavior-Named Grouping

Owned files:

- selected command-surface registry helper files

Deliverable:

- Move related functions from numbered shards into cohesive behavior-named
  files while preserving tests and public behavior.

### WS-021-C: Verification And Docs

Owned files:

- selected package docs if ownership or dependency direction changes
- `docs/package-ownership-map.md` when needed

Deliverable:

- Verify behavior and update ownership docs only when the cleanup changes
  package-level boundaries.

Slice 1 is expected not to change package boundaries or dependency direction,
so `docs/package-ownership-map.md` requires no content change unless review or
verification finds otherwise.

Slice 2 groups command-surface usage-drift helpers into collection and diff
files. A single combined drift file was rejected because it would fall below the
absolute file-MI threshold and force a mixed code/baseline PR.

Slice 3 groups only command-surface list helpers into
`cmd/sdp-trace/command_surface_list_helpers.go`. A broader argument-helper file
including `isHelp` and `isBoolLiteral` was rejected because local pre-change MI
analysis measured it below the absolute file-MI threshold.

Slice 4 groups command-surface schema type definitions and runner/JSON
functions into `command_surface_schema.go` and `command_surface_runner.go`.
Broader metadata/registry grouping was rejected because local pre-change MI
analysis measured candidate files below the absolute file-MI threshold.

Slice 5 groups the command-surface core command metadata shards into
`cmd/sdp-trace/command_surface_core_commands.go`. The combined file has no
functions and local pre-change MI analysis measured file MI `100.0`.

Slice 6 groups the command-surface observe command metadata shards into
`cmd/sdp-trace/command_surface_observe_commands.go`. The combined file has no
functions and local pre-change MI analysis measured file MI `100.0`.

Slice 7 groups the command-surface assess command metadata shards into
`cmd/sdp-trace/command_surface_assess_commands.go`. The combined file has no
functions and local pre-change MI analysis measured file MI `100.0`.

Slice 8 groups the command-surface packet command metadata shards into
`cmd/sdp-trace/command_surface_packet_commands.go`. The combined file has no
functions and local pre-change MI analysis measured file MI `100.0`.

Slice 9 groups the command-surface other command metadata shards into
`cmd/sdp-trace/command_surface_other_commands.go`. The combined file has no
functions and local pre-change MI analysis measured file MI `100.0`.

Slice 10 groups command-surface catalog metadata into
`cmd/sdp-trace/command_surface_catalog.go`. Broader registry/constants/catalog
grouping was rejected because pre-change MI analysis measured the candidate
file below the absolute file-MI threshold.

Slice 11 removes the remaining numbered files by splitting them into
responsibility-named files: CLI argument helpers, command-surface metadata
accessors, and command-surface registry assembly. A single combined file remains
rejected because it would mix unrelated behavior and risks MI regression.

Slice 12 removes the remaining numbered doctor local report/check shards in
`cmd/sdp-trace`. It keeps the behavior-preserving cleanup split by local doctor
responsibility instead of making one large file: report assembly, contract
checks, writable path checks, expected evidence checks, event vocabulary, CI
prerequisites, preview metadata, usage text, and small local result helpers. The
slice intentionally excludes other families (`assess`, `core`, `gate`,
`harnessobs`, `packet`, and `prreview`) so each family can keep a focused review
and verification trail.

Slice 13 removes the numbered `assess` command shards in `cmd/sdp-trace`. It
keeps the cleanup behavior-preserving by splitting along user-visible assess
responsibilities: command dispatch and flags, profile assessment runs, input
loading, artifact writers, preview dispatch, preview report shapes, and
profile-specific preview metadata. The slice intentionally excludes other
remaining families (`core`, `gate`, `harnessobs`, `packet`, and `prreview`) so
each family keeps a focused review and verification trail.

Slice 14 starts the numbered `core` cleanup with only the CLI kernel shards in
`cmd/sdp-trace`: exit constants, process runtime variables, handler types and registry,
top-level command dispatch, subcommand dispatch, required-flag helpers, JSON
payload helpers, and string exit-code mapping. It intentionally excludes the
remaining core assessment explain, preview action, export, and fixture helper
shards so each area can keep a focused review and verification trail.

Slice 15 continues the numbered `core` cleanup with the assessment explain,
assessment preview setup, and assessment exit-code mapping shards in
`cmd/sdp-trace`. It moves those declarations into responsibility-named
`assess_explain_*`, `assess_preview_*`, and `assess_exit_code*` files while
preserving command behavior, artifact interpretation, preview remediation text,
and exit-code semantics. It intentionally excludes release-proof, witness,
export/posture, and fixture expectation shards so those smaller responsibilities
can be reviewed independently.

Slice 16 removes the remaining numbered `core` shards in `cmd/sdp-trace` by
moving release-proof write helpers, witness target parsing, export command
helpers, telemetry export helpers, cross-repo posture export/explain helpers,
and fixture expectation metadata into responsibility-named files. The slice
closes the `core` numbered-file family while preserving CLI behavior, output
contracts, package boundaries, dependency direction, and MI/CRAP gates.

Slice 17 starts `internal/harnessobs` cleanup with numbered type, option,
context, scanner, and lookup-map shards (`harnessobs_011` through
`harnessobs_033`). It moves declarations into responsibility-named files for
options/context, session model, validation lookup sets, event reference checks,
existing path specs, shell field scanning, and isolation rule installers. It
intentionally excludes observation execution, parsing, path safety, and event
validation behavior shards so those behavior-heavy areas can keep focused
review trails.

Slice 18 continues `internal/harnessobs` cleanup with observe and session setup
entrypoint shards (`harnessobs_034` through `harnessobs_045`). It moves observe
entrypoint, option validation, path resolution, source loading, event writing,
run construction, observation context/time, and session setup entrypoint logic
into responsibility-named files. The slice preserves behavior and intentionally
keeps lower-level session setup execution, OpenCode normalization, raw signal
extraction, token safety, and validation behavior in later slices.

Slice 19 continues `internal/harnessobs` cleanup with OpenCode normalization
and session command model event shards (`harnessobs_046` through
`harnessobs_064`). It moves raw OpenCode line normalization, family detection
rules, event construction, observed-at/actor extraction, session command model
facts, and normalized event construction into responsibility-named files. It
intentionally excludes raw signal traversal (`harnessobs_065` onward), token
safety, validation loading, and event scanning so those behavior-heavy areas can
keep separate verification and review trails.

Slice 20 continues `internal/harnessobs` cleanup with raw signal traversal,
recursive key lookup, and timestamp extraction shards (`harnessobs_065` through
`harnessobs_091`). It moves raw signal dispatch, map/slice/string/scalar signal
collection, exact and prefix signal matching, key presence traversal, generic
key lookup, string and numeric key extraction, and timestamp parsing into
responsibility-named files. It intentionally excludes token safety, run loading,
validation loading, and event scanning for later slices.

Slice 21 continues `internal/harnessobs` cleanup with token safety, loaded run
reading, validation summary/loading, and profile validation shards
(`harnessobs_092` through `harnessobs_110`). It moves mutation-tool detection,
safe token rendering, run JSON/event loading, validation summary rendering,
validation artifact loading, profile metadata/family validation, and degradation
rule validation into responsibility-named files. It intentionally excludes event
source scanning and parsed-event validation (`harnessobs_111` onward) so the
line-scanning behavior can keep a separate focused review trail.

Slice 22 continues `internal/harnessobs` cleanup with event source scanning,
safe event decoding, parsed event validation, event identity/ref/content
validation, and unavailable-field validation shards (`harnessobs_111` through
`harnessobs_132`). It moves JSONL file scanning, line limits, source hashing,
safe raw-event rejection, typed event decoding, digest validation, observed-at
validation, reference checks, content-state checks, and unavailable-field checks
into responsibility-named files. It intentionally excludes evaluation/dimension
composition (`harnessobs_133` onward) and path-safety helpers so those can keep
separate review trails.

Slice 23 continues `internal/harnessobs` cleanup with evaluation and dimension
composition shards (`harnessobs_133` through `harnessobs_142`). It moves
validation assembly and event-family dimension composition into
responsibility-named files while preserving state ranking, schema-version
mismatch handling, event-family counts, dimension ordering, and validation
digest generation. It intentionally excludes path-safety helpers
(`harnessobs_143` onward) so filesystem trust boundaries keep a separate
review trail.

Slice 24 continues `internal/harnessobs` cleanup with existing path safety
helpers (`harnessobs_143` through `harnessobs_149`). It moves existing
file/directory validation, traversal rejection, symlink resolution,
working-directory containment, and expected path-type checks into
responsibility-named files. It intentionally excludes output file and output
directory path helpers (`harnessobs_150` onward) so creation-path behavior keeps
a separate review trail.

Slice 25 continues `internal/harnessobs` cleanup with output file and parent
path helpers (`harnessobs_150` through `harnessobs_161`). It moves output file
path validation, output basename validation, parent path normalization,
missing-parent resolution, symlink-aware parent containment, and working
directory relative path conversion into responsibility-named files. It
intentionally excludes output directory creation/emptiness helpers
(`harnessobs_162` onward) so directory materialization behavior keeps a
separate review trail.

Slice 26 continues `internal/harnessobs` cleanup with output directory safety
and emptiness helpers (`harnessobs_162` through `harnessobs_173`). It moves
output directory traversal rejection, existing-output symlink resolution,
missing-output parent containment, working-directory escape checks, and
empty-or-missing directory validation into responsibility-named files. It
intentionally excludes JSON writing, event refs, and digest helpers
(`harnessobs_174` onward) so artifact serialization keeps a separate review
trail.

Slice 27 continues `internal/harnessobs` cleanup with artifact serialization,
event reference, and digest helpers (`harnessobs_174` through
`harnessobs_180`). It moves JSON writing, event reference rendering and
validation, raw-line digesting, validation digesting, and command digesting into
responsibility-named files. It intentionally excludes command model extraction
and shell field parsing (`harnessobs_181` onward) so command parsing behavior
keeps a separate review trail.

Slice 28 continues `internal/harnessobs` cleanup with command model extraction
and controlled shell field parsing (`harnessobs_181` through
`harnessobs_197`). It moves command model extraction from argv and `sh`/`bash
-c`, model flag parsing, and shell field scanner stages into
responsibility-named files. It intentionally excludes command model safety,
source commit, session setup, raw event safety, and validation helpers
(`harnessobs_198` onward) so safety checks and session behavior keep separate
review trails.

Slice 29 continues `internal/harnessobs` cleanup with command model safety and
source-bound digest primitives (`harnessobs_198` through `harnessobs_203`). It
moves command model safety normalization/rejection, file digest calculation,
and current source commit discovery into responsibility-named files. It
may reuse the existing package-local SHA-256 helper and keep source commit hash
validation package-local if focused regression evidence confirms unchanged
behavior.
It intentionally excludes unsafe raw-event traversal (`harnessobs_204` onward),
session setup, collection, and validation helpers so raw event safety and
session behavior keep separate review trails.

Slice 30 continues `internal/harnessobs` cleanup with unsafe value traversal
entrypoints and the shared value dispatcher (`harnessobs_204` through
`harnessobs_208`). It moves generic unsafe traversal entrypoints, raw-event
unsafe traversal entrypoints, mode-specific wrappers, and the type dispatcher
into a responsibility-named file while preserving path rendering, reason-code
delegation, raw-event mode, and map/slice/string delegation behavior. It
intentionally excludes session setup (`harnessobs_209` onward) and the
map/slice/string-specific unsafe rule shards (`harnessobs_223` onward) so
workflow execution and rule semantics keep separate review trails.

Slice 31 continues `internal/harnessobs` cleanup with session setup path,
run construction, command metadata, time fallback, and session JSON writing
shards (`harnessobs_209` through `harnessobs_216`). It moves setup path
validation, setup run preparation, command digest/model assignment, session
time fallback, and setup `session.json` writing into responsibility-named files
while preserving required option errors, path safety, output creation, isolation
rule installation, command digest/model states, and JSON output behavior. It
intentionally excludes session collection (`harnessobs_217` onward) and
raw-event unsafe rule semantics (`harnessobs_223` onward) so collection and
rule-specific behavior keep separate review trails.

Slice 32 is a corrective cleanup for Slice 27 microfile drift. It consolidates
artifact JSON writing, event reference rendering/safety, and digest helpers back
into cohesive responsibility groups after review found that the final Slice 27
implementation reintroduced one-helper microfiles despite the reviewed plan.
The correction preserves behavior and MI/CRAP gates while treating MI pressure
as a reason for meaningful comments or tighter grouping, not as a blanket
justification for shard churn.

Slice 33 continues `internal/harnessobs` cleanup with session collection
entrypoint and input-loading shards (`harnessobs_217` through
`harnessobs_222`). It moves collection dispatch, collect-option validation
handoff, profile/session input loading, profile mismatch rejection, harness
profile loading handoff, context construction, and collection time fallback into
cohesive session collection files. It intentionally excludes event source
resolution, source normalization, observed-run writing, process execution, and
raw-event unsafe rule semantics so source collection and rule-specific behavior
keep separate review trails.

Slice 34 continues `internal/harnessobs` cleanup with raw unsafe traversal rule
semantics (`harnessobs_223` through `harnessobs_244`). It moves map/slice/string
unsafe traversal, raw-event skip rules, path/token/url checks, digest-field and
raw path-like exemptions into cohesive safety-rule files. It intentionally
excludes validation enum helpers (`harnessobs_245` onward), session collect
option validation, event source resolution, and raw-event normalization flow so
schema validation and collection orchestration keep separate review trails.

Slice 35 continues `internal/harnessobs` cleanup with validation enum, safe
reference, non-authority boundary text, and validation decoding helpers
(`harnessobs_245` through `harnessobs_255`). It moves validation lookup
accessors, generic/source/task/operation/actor reference safety helpers,
digest-mismatch event rendering, the evidence-only non-authority boundary, and
validation JSON decoding into cohesive validation helper files. It intentionally
excludes session collect option validation (`harnessobs_256` onward), profile
loading, event source resolution, runtime collection, and validation command
execution so option validation and orchestration keep separate review trails.

Slice 36 continues `internal/harnessobs` cleanup with session collection
validation, harness profile/source resolution, observed-run materialization, and
runtime session command execution shards (`harnessobs_256` through
`harnessobs_280`). It moves collect option requirements, harness profile
loading, raw-event source normalization, source-unavailable fallbacks, observed
run/event writing, collected-session finalization, `RunSession` setup/command
execution, process metadata, and finished-session collection into cohesive
session collection/runtime files. The first four-file grouping was rejected
after local MI checks because it produced dense low-MI files; the accepted split
uses responsibility-named files for collect options, source resolution, raw
normalization, source-unavailable fallback, observed collection/finalization,
observed output serialization, runtime setup, runtime finish, process execution,
and process metadata. Maintainability review rejected a follow-up split that
left observed source collection and finalization as tiny non-numbered helper
files; a one-file observed collection merge then failed the local MI gate. The
accepted split keeps source collection/finalization in
`session_collect_observed.go` and observed event/run serialization in
`session_observed_output.go`, avoiding both one-helper drift and an MI
regression. It intentionally excludes validation command execution
(`harnessobs_281` onward), validation evaluation, session profile validation,
isolation rule installation, loaded session run validation, and raw
normalization internals beyond dispatch from session collection so those
behavior-heavy areas keep separate review trails.

Slice 37 continues `internal/harnessobs` cleanup with validation command
entrypoint, validation input requirements, safe validation path resolution,
run-loading evaluation fallback, source-unavailable validation construction, and
optional validation artifact writing shards (`harnessobs_281` through
`harnessobs_290`). It moves the `Validate` orchestration, required option
checks, shared non-blank helper, profile/run/out path resolution,
`LoadRun`-backed evaluation, cannot-verify fallback, and optional output write
into cohesive validation command files. Local MI rejected a single validation
input file, so required-option/non-blank checks and safe path resolution are
kept in separate responsibility-named files. It intentionally excludes profile
loading and session profile validation (`harnessobs_291` onward), raw-event
config validation, isolation rule validation/installation, loaded session run
validation, and path-safety primitives so those behavior-heavy areas keep
separate review trails.

Slice 38 continues `internal/harnessobs` cleanup with harness/session profile
loading, session profile identity/path validation, and raw-event config pair
validation shards (`harnessobs_291` through `harnessobs_299`). It moves
`LoadProfile`, `LoadSessionProfile`, session profile orchestration, identity
checks, required session path checks, raw event format/source pairing, and raw
event format support checks into cohesive profile loading and profile
validation files. It intentionally excludes stream capture normalization
(`harnessobs_300` onward), session setup action validation, isolation rule
validation/installation, loaded session run validation, session run
construction, source commit discovery, and raw-event normalization execution so
those behavior-heavy areas keep separate review trails.

Slice 39 continues `internal/harnessobs` cleanup with session stream capture,
session setup action validation, and isolation rule validation shards
(`harnessobs_300` through `harnessobs_309`). It moves stream capture defaulting
and unsupported-mode errors, setup action count/id/kind validation, isolation
rule list/id/pattern/target/kind validation, and unsafe isolation pattern checks
into cohesive session profile rule validation files. It intentionally excludes
isolation rule target resolution and installation (`harnessobs_310` onward),
line/JSON rule materialization, loaded session run validation, session run
construction, source commit discovery, and raw-event normalization execution so
filesystem mutation and run-loading behavior keep separate review trails.

Slice 40 continues `internal/harnessobs` cleanup with isolation rule target
resolution, line/JSON rule materialization, verification readback, and isolation
result digest construction shards (`harnessobs_310` through `harnessobs_334`).
It moves profile-relative isolation file safety, parent/filename validation,
installer dispatch, line rule append/read/write helpers, JSON read-deny object
loading and mutation helpers, readback verification, cannot-verify fallback,
and isolation result digest assignment into cohesive isolation installation
files. The implementation keeps JSON mutation separate from optional object
loading and keeps readback result construction separate from presence checks to
meet the repository MI gate without returning to numbered one-helper shards.
JSON readback continues to call the existing package-local JSON reader helpers,
which remain numbered for a later slice. It intentionally excludes
loaded session run validation (`harnessobs_335` through `harnessobs_336`),
shared JSON read/decode helpers (`harnessobs_337` through `harnessobs_339`),
session run construction (`harnessobs_340` onward), source commit discovery,
event source reading, profile-relative source/output file safety, and raw-event
normalization execution so run-loading, shared JSON IO, and event
materialization behavior keep separate review trails.

Slice 41 continues `internal/harnessobs` cleanup with loaded session run
loading/validation and shared existing-JSON reader helpers (`harnessobs_335`
through `harnessobs_339`). It moves `LoadSessionRun`,
`validateLoadedSessionRun`, permissive existing JSON loading, strict existing
JSON loading, and strict decoder trailing-data rejection into cohesive session
loading and JSON loading files. It intentionally excludes session run
construction (`harnessobs_340` through `harnessobs_342`), source commit
discovery (`harnessobs_343`), event source reading (`harnessobs_344`),
profile-relative source/output file safety (`harnessobs_345` through
`harnessobs_347`), and raw-event normalization execution (`harnessobs_348`
onward) so construction, source discovery, path safety, and raw normalization
keep separate review trails.

Slice 42 continues `internal/harnessobs` cleanup with session run construction
shards (`harnessobs_340` through `harnessobs_342`). It moves `newSessionRun`,
`newSessionRunRecord`, and setup action ID collection/sorting into one cohesive
session run construction file. It intentionally excludes source commit
discovery (`harnessobs_343`), event source reading (`harnessobs_344`),
profile-relative source/output file safety (`harnessobs_345` through
`harnessobs_347`), and raw-event normalization execution (`harnessobs_348`
onward) so construction defaults, source provenance, path safety, and raw
normalization keep separate review trails.

Slice 43 continues `internal/harnessobs` cleanup with source commit state
mapping (`harnessobs_343`). It moves `currentSourceCommitState` into the
existing `source_commit.go` file next to `sourceCommit`, preserving the
fail-closed mapping from an empty/invalid commit to `cannot_verify` and a valid
commit to `pass`. It intentionally excludes event source reading
(`harnessobs_344`), profile-relative source/output file safety
(`harnessobs_345` through `harnessobs_347`), and raw-event normalization
execution (`harnessobs_348` onward) so source provenance state, event reading,
path safety, and normalization keep separate review trails.

Slice 44 continues `internal/harnessobs` cleanup with event source reading
handoff (`harnessobs_344`). It moves `readEventsFromPath` into the existing
`event_scan_input.go` file next to `readEvents`, preserving profile loading,
event scan delegation, source digest return, and error propagation. It
intentionally excludes profile-relative source/output path safety
(`harnessobs_345` through `harnessobs_347`) and raw-event normalization
execution (`harnessobs_348` onward) so event reading handoff, path safety, and
normalization keep separate review trails.

Slice 45 continues `internal/harnessobs` cleanup with profile-relative
source/output path safety (`harnessobs_345` through `harnessobs_347`). It moves
`safeProfileRelativeFile`, `safeProfileRelativeOutFile`, and
`unsafeProfileRelativePath` into a cohesive `session_profile_paths.go` file
used by session collection, raw normalization, setup isolation, and isolation
rule validation. It preserves absolute-path, URL-like, traversal, base
directory, existing-file, and output-file policy behavior. It intentionally
excludes raw-event normalization execution (`harnessobs_348` onward) so source
path safety and event normalization keep separate review trails.

Slice 46 continues `internal/harnessobs` cleanup with raw-event normalization
execution (`harnessobs_348` through `harnessobs_360`). It moves raw
normalization orchestration, input validation, OpenCode JSONL scanning, raw
line decoding, unsafe raw-event rejection, normalized source digesting, and
normalized event writing into cohesive raw-normalization files split by
responsibility to satisfy MI without baseline changes. It moves shared blank
JSONL line handling into the neutral event-line parsing file because the helper
is used by both normal event parsing and raw normalization. It preserves
supported-format gating, raw/output same-file rejection,
zero-time fallback, scanner limits, malformed JSONL errors, unsafe-input
errors, source digest calculation, output parent creation, and JSONL write
format. It intentionally excludes generic unsafe raw-value discovery and
OpenCode event construction helpers already housed in non-numbered files.

Slice 47 starts `cmd/sdp-trace` cleanup with checkpoint CLI command handling
(`gate_271` through `gate_289`). It moves checkpoint subcommand routing,
create flag parsing, checkpoint creation/write handoff, verify flag parsing,
checkpoint/policy input loading, and verify exit-code mapping into cohesive
checkpoint CLI files split by command responsibility. It preserves create and
verify flag names/defaults, usage errors, stderr/stdout behavior, JSON artifact
read/write behavior, optional policy semantics, checkpoint verification result
rendering, exit-code mapping, package boundary, dependency direction, and MI
baselines. It intentionally excludes protected-gate checkpoint policy/witness
logic (`gate_302` onward) and shared JSON/text file helpers (`gate_360`
onward) so checkpoint CLI behavior, protected-gate trust rules, and generic IO
helpers keep separate review trails.

Slice 48 continues `cmd/sdp-trace` cleanup with protected-gate core execution
and input loading (`gate_302` through `gate_311`). It moves protected gate
run/resolve/evaluate orchestration, protected checkpoint replay handoff, result
writing, required checkpoint/policy/witness input loading, protected row
loading, and witness expectation loading into `protected_gate_core.go`,
`protected_gate_inputs.go`, and `protected_gate_loaders.go`. It preserves
fail-closed setup errors for missing/malformed external trust inputs,
contract/row/run-dir/witness expectation error codes,
`PolicyProvided: true`, UTC evaluation time, result JSON writing/rendering,
gate exit-code behavior, package boundary, dependency direction, and MI
baselines. It intentionally excludes protected checkpoint trust matching
(`gate_334` through `gate_344`), demo witness construction (`gate_345` onward),
gate explain (`gate_312` through `gate_324`), preview (`gate_325` onward),
override request handling (`gate_352` onward), and shared JSON/text helpers
(`gate_360` onward) so setup/evaluation orchestration, trust matching, explain
rendering, preview, override, and generic IO keep separate review trails.

Slice 49 continues `cmd/sdp-trace` cleanup with gate explain CLI and rendering
(`gate_312` through `gate_324`). It moves read-only explain argument parsing,
gate-result artifact loading/schema validation, summary/protected detail
rendering, collection rendering, reasons, and next actions into
`gate_explain_cli.go`, `gate_explain_renderer.go`, and
`gate_explain_collections.go`, and shared reason/next-action rendering into
the neutral `explain_common_collections.go` because those helpers are also used
by assessment explain renderers. It preserves usage errors, missing/malformed
gate-result artifact `cannot_verify`, unsupported schema `cannot_verify`,
read-only behavior, legacy protected-field absence output, protected
checkpoint/condition detail lines, required-run/witness/override/
missing-evidence/reason/next-action rendering, secret non-disclosure by not
printing raw run commands, package boundary, dependency direction, and MI
baselines. It intentionally excludes gate preview (`gate_325` onward),
protected run-dir/trust matching (`gate_333` onward), override request handling
(`gate_352` onward), and shared JSON/text helpers (`gate_360` onward) so
explanation, preview, protected trust matching, override, and generic IO keep
separate review trails.

Slice 50 continues `cmd/sdp-trace` cleanup with gate preview and protected
target selection (`gate_325` through `gate_333`). It moves standard preview
report types, preview argument parsing, standard preview execution/report
building, protected preview execution/report construction, and the shared
protected run-dir selector into `gate_preview_cli.go`,
`gate_preview_args.go`, `gate_preview_standard.go`, `gate_preview_reports.go`,
`gate_preview_protected.go`, and neutral `protected_gate_run_dir.go`. It
preserves preview read-only behavior, target
arity usage errors, contract load failure behavior, standard preview report
fields (`required_runs`, `required_evidence`, `witness_inspectable`,
`witness_mismatches`, and `claim`), witness mismatch reporting without issuing
any gate verdict fields, protected preview absent/unreadable/malformed input
statuses and `cannot_verify` setup exit, secret non-disclosure by not printing
raw run commands, protected single-run selection semantics, package boundary,
dependency direction, and MI baselines. It intentionally excludes protected
checkpoint trust matching (`gate_334` through `gate_344`), demo witness
construction (`gate_345` onward), protected preview status/action helpers
(`gate_349` through `gate_351`, dependency only), override request handling
(`gate_352` onward), shared JSON/text helpers (`gate_360` onward), and preview
mode/required-ID helper shards (`gate_365` onward) so preview reporting,
protected trust matching, witness construction, protected input-status
rendering, override, generic IO, and shared preview helpers keep separate
review trails.

Slice 51 continues `cmd/sdp-trace` cleanup with protected checkpoint trust
matching (`gate_334` through `gate_344`). It moves protected checkpoint upgrade
selection, signer policy matching, witness protected-trust/source matching,
artifact digest matching, and optional source-field comparison into
`protected_checkpoint_trust.go`, `protected_checkpoint_signer.go`,
`protected_witness_match.go`, and `protected_witness_artifacts.go`. It
preserves the rule that explicit checkpoint failures cannot be upgraded,
protected trust requires CI-isolated signer authority, signer policy must bind
signer id, authority, and public key, witness status/source must match before
artifact matching, empty expected source fields remain wildcards, artifact
counts must match exactly, and no protected gate/evidence schema behavior or MI
baseline is changed. It intentionally excludes demo witness expectation and
artifact construction (`gate_345` through `gate_348`), protected preview
status/action helpers (`gate_349` through `gate_351`), override request
handling (`gate_352` onward), shared JSON/text helpers (`gate_360` onward), and
preview mode/required-ID helper shards (`gate_365` onward) so protected trust
matching, demo witness construction, preview status rendering, override, and
generic IO keep separate review trails.

Slice 52 continues `cmd/sdp-trace` cleanup with demo witness expectation and
artifact digest construction (`gate_345` through `gate_348`). It moves
protected-gate witness expectation loading, run artifact list construction,
retained `run.json` digest calculation, and per-run artifact metadata into
`protected_witness_expectation.go` and `protected_witness_digest.go`. It
preserves deriving expectations from observed run directories rather than the
supplied witness summary, using the first discovered run ID as the expected
run ID, retaining one `<run-dir-base>/run.json` artifact path per discovered
run, calculating SHA-256 from retained file bytes, propagating discovery/open
errors as `cannot_verify` through the existing loader, package boundary,
dependency direction, and MI baselines. It intentionally excludes protected
preview status/action helpers (`gate_349` through `gate_351`), override request
handling (`gate_352` onward), shared JSON/text helpers (`gate_360` onward), and
preview mode/required-ID helper shards (`gate_365` onward) so witness
expectation construction, preview status rendering, override, and generic IO
keep separate review trails.

Slice 53 continues `cmd/sdp-trace` cleanup with protected preview input status
and remediation actions (`gate_349` through `gate_351`). It moves protected
input status classification, input error status mapping, and stable preview
next-action generation into `protected_preview_inputs.go`. It preserves blank
input as `absent`, missing/permission-denied input as `present_unreadable`,
malformed JSON as `present_malformed`, readable JSON as `present_readable`,
stable action ordering for `checkpoint`, `checkpoint_policy`, and `witness`,
protected preview `cannot_verify` setup exit for unreadable/malformed inputs,
package boundary, dependency direction, and MI baselines. It intentionally
excludes override request handling (`gate_352` onward), shared JSON/text
helpers (`gate_360` onward), and preview mode/required-ID helper shards
(`gate_365` onward) so preview input readiness, override, and generic IO keep
separate review trails.

Slice 54 continues `cmd/sdp-trace` cleanup with override request CLI handling
(`gate_352` through `gate_359`). It first attempted to move override request
dispatch, trace-event append, payload construction, request subcommand
detection, flag parsing, required-flag validation, and required flag diagnostics
into `override_request.go`; file-level MI failed, so the implementation split
the same responsibility into `override_request.go`,
`override_request_payload.go`, and `override_request_flags.go` without changing
behavior or MI baselines. It preserves that only `override request` is accepted,
parser errors and missing required flags exit with `exitUsage` before appending
any trace event, positional text is rejected, persisted payload keys remain
stable (`override_id`, `producer`, `origin`, `requested_by`, `reason`,
`source_ref`, `scope`, `created_at`, optional `external_reference`),
`created_at` remains UTC RFC3339Nano, stdout prints only `override_event:
<event-id>`, append failure prints stderr and returns exit code 1, override
events remain advisory/non-upgrading, package boundary, dependency direction,
and MI baselines. It intentionally excludes shared JSON/text helpers
(`gate_360` onward) and preview mode/required-ID helper shards (`gate_365`
onward) so override write semantics, generic IO, and shared preview helpers
keep separate review trails.

Slice 55 continues `cmd/sdp-trace` cleanup with shared artifact IO helpers
(`gate_360` through `gate_364`). It moves shared JSON reads, pretty JSON writes,
atomic text writes, temp-file close handling, chmod-before-rename, and rename
publication first into `artifact_io.go`; file-level MI failed, so the
implementation split the same responsibility into `artifact_json_io.go` and
`artifact_text_io.go` without changing behavior or MI baselines. It preserves
read error and unmarshal error propagation, JSON parent-directory creation,
two-space JSON indentation with a trailing newline, JSON `os.WriteFile`
requested `0o644` mode subject to process umask, text parent-directory creation,
sibling temp-file creation, temp cleanup after successful rename or failure,
close-on-write-error behavior, `0o644` text file mode before rename, atomic
rename publication, package boundary, dependency direction, and MI baselines.
It intentionally excludes preview mode and required-ID helper shards
(`gate_365` onward) so generic artifact IO and preview-specific argument
helpers keep separate review trails.

Slice 56 continues `cmd/sdp-trace` cleanup with gate preview contract helpers
(`gate_365` through `gate_367`). It moves preview gate mode selection,
required run ID extraction, and required evidence ID extraction into
`gate_preview_contract.go`. It preserves observation mode as the default,
advisory CI mode when at least one required run asks for advisory CI,
protected-future dominance over advisory CI regardless of order, required run
ID order with empty IDs omitted, required evidence ID order with empty IDs
omitted, package boundary, dependency direction, and MI baselines. It
intentionally excludes packet/PR review shards (`packet_031` onward) so gate
preview contract display helpers and packet workflows keep separate review
trails.

Slice 57 continues `cmd/sdp-trace` cleanup with packet command surface helpers
(`packet_031` through `packet_032`). It moves the packet subcommand handler
registry and packet required-flag definitions into `packet_command_surface.go`.
It preserves the exact packet subcommands (`build-pr`, `build-github`,
`validate`, `check-demo`, and `render`), each handler binding, required flag
names and diagnostic messages, package boundary, dependency direction, and MI
baselines. It intentionally excludes packet command execution and artifact
building shards (`packet_040` onward) so command-surface metadata and packet
workflow implementation keep separate review trails.

Slice 58 continues `cmd/sdp-trace` cleanup with the `packet build-pr` command
flow and PR packet artifact publication helpers (`packet_040` through
`packet_050`). It moves packet command dispatch, `build-pr` option parsing,
PR packet result construction, markdown rendering, output-directory creation,
artifact file list construction, and sequential artifact writes into
`packet_command_dispatch.go`, `packet_build_pr_run.go`,
`packet_build_pr_options.go`, `packet_build_pr_result.go`,
`packet_build_pr_artifact_render.go`, `packet_build_pr_artifact_write.go`, and
`packet_build_pr_artifact_files.go`; the first single-file attempt at
`packet_build_pr_command.go` failed file-level MI at 56.6, and the second
three-file split still failed at 67.7/64.1 for command/artifact files, so the
final implementation keeps root dispatch, build-pr run orchestration, option
parsing, result construction, artifact rendering, artifact writing, and
artifact file-list construction separate. It preserves the packet
missing-subcommand diagnostic, `packet build-pr` flag defaults and required flag
validation, JSON `cannot_verify` output on input reconstruction and render
failure, validation plus live-gate cannot_verify result state and error
aggregation, output paths (`bundle.json`, `change-evidence-packet.md`, and
`build-pr-result.json`), output-directory creation with mode subject to process
umask, artifact write
labels and ordering, first-write-failure short-circuiting, package boundary,
dependency direction, and MI baselines. It
intentionally excludes PR input reconstruction, route/error classification,
GitHub Actions hydration, GitHub API access, fixture loading, and packet exits
(`packet_051` onward) so build-pr orchestration and evidence-source hydration
keep separate review trails.

Slice 59 continues `cmd/sdp-trace` cleanup with `packet build-pr` live gate
error aggregation (`packet_051` through `packet_053`). It moves packet row
indexing, route-readiness error construction, and verification-readiness error
construction into `packet_build_pr_gate_errors.go`. It preserves row lookup by
packet row ID, route rows passing when `PC-AGENT-ROUTE` is `pass` or
`partial`, verification passing only when `PC-VERIFICATION` is `pass`, route
error ordering before verification error ordering, diagnostic strings including
row reasons, package boundary, dependency direction, and MI baselines. It
intentionally excludes PR input reconstruction/source loading (`packet_054`
through `packet_059`), event conversion (`packet_060` onward), and GitHub
Actions hydration/API helpers so live gate readiness and evidence-source
loading keep separate review trails.

Slice 60 continues `cmd/sdp-trace` cleanup with `packet build-pr` PR input
source loading and optional evidence enrichment (`packet_054` through
`packet_059`). It moves PR input reconstruction from parsed options, source
validation, event path selection, fixture/event loading, optional checks and
artifact JSON reads, GitHub Actions hydration dispatch, route-manifest loading,
and route application into `packet_build_pr_input_source.go` and
`packet_build_pr_input_enrichment.go`; the first single-file attempt at
`packet_build_pr_input_loading.go` failed file-level MI at 65.4, so the final
implementation keeps source/event loading separate from optional evidence,
hydration dispatch, and route enrichment. It preserves the
allowed source set (`github-actions`, `github-fixture`), unsupported-source and
missing-event diagnostics, GitHub Actions fallback to `GITHUB_EVENT_PATH` only
when `--github-event` is empty, fixture mode using explicit event paths, local
optional evidence error prefixes (`read checks json`, `read artifacts json`),
route manifest error prefix (`read route manifest`), route application after
hydration, fixture-mode hermeticity, package boundary, dependency direction,
and MI baselines. It intentionally excludes event-to-input conversion
(`packet_060` through `packet_062`), GitHub Actions hydration implementation
(`packet_063` onward), route application internals (`packet_066`), and shared
optional JSON reading (`packet_095`) so input loading, event conversion,
hydration, route mutation, and shared IO keep separate review trails.

Slice 61 continues `cmd/sdp-trace` cleanup with `packet build-pr`
event-to-input mapping (`packet_060` through `packet_062`). It moves GitHub PR
evidence input construction from a loaded PR event, PR field projection, and
commit-range projection into `packet_build_pr_event_mapping.go`. It preserves
schema version `github-pr-evidence-input.v0`, prompt-boundary requirement
defaulting to true, GitHub Actions workflow run ID from `GITHUB_RUN_ID`,
fixture workflow run ID from the event payload, PR number/URL/title/body
ref/author/base ref/head ref/head SHA mapping, commit base/head SHA mapping,
changed-files ref from the event diff URL, package boundary, dependency
direction, and MI baselines. It intentionally excludes input source loading
(`packet_054` through `packet_059`), GitHub Actions hydration (`packet_063`
onward), and packet fixture type/loading (`packet_093` onward) so source
selection, event projection, hydration, and fixture IO keep separate review
trails.

Slice 62 continues `cmd/sdp-trace` cleanup with `packet build-pr` GitHub
Actions hydration dispatch (`packet_063` through `packet_064`). It moves
source-gated hydration and artifact backfill dispatch into
`packet_build_pr_actions_hydration.go`. It preserves fixture-mode no-network
behavior, `github-actions` source gating for live hydration, explicit artifact
JSON precedence over live discovery, live artifact discovery error propagation,
package boundary, dependency direction, and MI baselines. It intentionally
excludes route manifest loading/application (`packet_065` through
`packet_066`) and live artifact discovery/context/API helpers (`packet_067`
onward) so route mutation and live API access keep separate review trails.

Slice 63 continues `cmd/sdp-trace` cleanup with `packet build-pr` route manifest
helpers (`packet_065` through `packet_066`). It moves optional route manifest
reading and route, prompt-boundary, integration-action, and review field
application into `packet_build_pr_route.go`. It preserves optional JSON read
behavior, route manifest field overwrite semantics, PR identity and CI evidence
preservation, package boundary, dependency direction, and MI baselines. It
intentionally excludes GitHub Actions artifact
discovery/context/API helpers (`packet_067` onward) and shared optional JSON IO
(`packet_095`) so route mutation and shared IO keep separate review trails.

Slice 64 continues `cmd/sdp-trace` cleanup with the GitHub Actions artifact
discovery facade and response type shards (`packet_067` through `packet_068`).
It moves live artifact discovery orchestration and GitHub artifact payload/type
definitions into `packet_build_pr_actions_artifacts.go`. It preserves validated
context construction before network fetch, fetch error propagation, retained
artifact filtering, fail-closed empty-retained-set diagnostics, package
boundary, dependency direction, and MI baselines. The GitHub build and
prompt-boundary classification areas may split into multiple named locality
files to keep MI above the absolute threshold without baseline changes. It
intentionally excludes
artifact context validation and URL/token policy (`packet_071` through
`packet_085`), HTTP request/fetch/decode/retention helpers (`packet_086`
through `packet_092`), fixture IO (`packet_093` onward), and shared optional
JSON IO (`packet_095`) so live API policy and response processing keep separate
review trails.

Slice 65 continues `cmd/sdp-trace` cleanup with GitHub Actions artifact context
construction and source selection shards (`packet_071` through `packet_075`).
It moves artifact context construction, context validation, and missing identity
detection into `packet_build_pr_actions_context.go`, with token selection and
API URL selection in `packet_build_pr_actions_source.go` after a single-file
attempt failed MI. It preserves API URL flag-over-env-over-default precedence,
URL validation before returned context construction, trailing-slash trimming,
`GITHUB_TOKEN` precedence over `GH_TOKEN`, missing repo/run and token
diagnostics, package boundary, dependency direction, and MI baselines. It
intentionally excludes URL parsing, trust-target policy,
HTTPS/loopback/host validation internals (`packet_076` through `packet_085`),
HTTP request/fetch/decode/retention helpers (`packet_086` through
`packet_092`), fixture IO (`packet_093` onward), and shared optional JSON IO
(`packet_095`) so security policy and response processing keep separate review
trails.

Slice 66 continues `cmd/sdp-trace` cleanup with GitHub Actions API URL
validation and trust-target policy shards (`packet_076` through `packet_085`).
It moves API URL validation flow into `packet_build_pr_actions_url_policy.go`,
API URL parsing, HTTPS enforcement, local HTTP detection, and loopback host
detection into `packet_build_pr_actions_url_parse.go`, and public/Enterprise
host binding plus configured server host extraction into
`packet_build_pr_actions_url_host.go` after a single-file attempt failed MI. It
preserves syntax diagnostics, credential-leak prevention, credential rejection
winning over HTTPS errors for mixed-invalid URLs, loopback-only HTTP test
allowance, public GitHub `github.com` to `api.github.com` mapping,
Enterprise exact-host binding, malformed server URL fallback behavior, package
boundary, dependency direction, and MI baselines. It intentionally excludes
context/source selection (`packet_build_pr_actions_context.go` and
`packet_build_pr_actions_source.go`), HTTP request/fetch/decode/retention
helpers (`packet_086` through `packet_092`), fixture IO (`packet_093` onward),
and shared optional JSON IO (`packet_095`) so URL policy and response
processing keep separate review trails.

Slice 67 continues `cmd/sdp-trace` cleanup with GitHub Actions artifact HTTP
request/fetch/decode and retained artifact shaping shards (`packet_086`
through `packet_092`). A single combined response-processing file was rejected
because pre-change MI analysis measured file MI `66.6`, below the absolute
threshold. The initial HTTP/retention split also failed when the HTTP file
measured file MI `67.5`, so the slice instead moves live artifact HTTP fetch,
request construction, HTTPS-only token attachment, status/decode helpers, and
retained artifact shaping into `packet_build_pr_actions_fetch.go`,
`packet_build_pr_actions_request.go`,
`packet_build_pr_actions_authorization.go`,
`packet_build_pr_actions_decode.go`, and
`packet_build_pr_actions_retention.go`. It preserves credential fail-closed
behavior for malformed or non-HTTPS API URLs, loopback-test no-token behavior,
GitHub media type headers, non-2xx fail-closed diagnostics, JSON decode
diagnostics, expired artifact filtering, artifact URL precedence over
synthesized resolver URLs, missing artifact ID empty-resolver behavior, package
boundary, dependency direction, and MI baselines. It intentionally excludes
fixture event loading (`packet_093` through `packet_094`), shared optional JSON
IO (`packet_095`), and CLI exit helpers (`packet_096`) so live response
processing and fixture/shared IO keep separate review trails.

Slice 68 continues `cmd/sdp-trace` cleanup with the remaining numbered packet
fixture and validation exit helper shards (`packet_093` through `packet_096`).
A single combined fixture IO file was rejected because pre-change MI analysis
measured file MI `66.3`. The slice instead moves PR fixture event shape and
validation into `packet_build_pr_fixture_event.go`, shared optional JSON
loading into `packet_build_pr_optional_json.go`, and packet validate/check-demo
exit mapping into `packet_validation_exits.go`. It preserves required PR
fixture identity validation, optional empty-path no-op behavior, JSON
read/unmarshal error propagation, `pass` to zero exit mapping, packet
validation failure to `cannot_verify`, demo gate failure to `fail`, package
boundary, dependency direction, and MI baselines. It intentionally excludes
later numbered packet command shards so the next packet responsibilities can
keep separate review trails.

Slice 69 starts the numbered `pr_review` cleanup in `cmd/sdp-trace` with
top-level `pr-review` dispatch and packet subcommand setup shards
(`pr_review_030`, `pr_review_037` through `pr_review_039`,
`pr_review_098` through `pr_review_104`, and `pr_review_138`). A combined top-level command/packet
flag file was rejected because pre-change MI analysis measured file MI `56.6`,
and a combined packet command/options file was rejected at file MI `68.1`.
The initial packet command file also failed after adding required-input checks,
measuring file MI `68.6`. The slice instead moves top-level dispatch into
`pr_review_command.go`, packet flag metadata into `pr_review_packet_flags.go`,
packet argument parsing and required-input checks into
`pr_review_packet_args.go`, packet execution into `pr_review_packet_run.go`,
and packet option construction into `pr_review_packet_options.go`. It preserves
subcommand routing, usage and missing-subcommand diagnostics, required packet
flags and defaults, repeated context/verification flag reconstruction, optional
metadata mapping, positional-argument rejection, missing packet-anchor usage
errors, packet build failures mapping to `cannot_verify`, packet stdout
rendering, provenance anchor mapping, package boundary, dependency direction,
and MI baselines. It intentionally excludes
`pr-review run`, synthesize, validate, summarize, check, shared file helpers,
and runner helpers so later review workflows can keep separate review trails.

Slice 70 continues numbered `pr_review` cleanup in `cmd/sdp-trace` with
`pr-review run` execution shards (`pr_review_105` through `pr_review_108`). A
single combined `pr_review_run_command.go` was rejected because pre-change MI
analysis measured file MI `66.1`. The slice instead moves run execution and
reviewer runner handoff into `pr_review_run_command.go`, run argument parsing
into `pr_review_run_args.go`, and preview-vs-run output rendering into
`pr_review_run_output.go`. It preserves packet/profile loading, work-dir
directory validation, repeated allowed-runner reconstruction from raw args,
preview mode, not-assessed reason propagation, parse errors and positional
argument rejection as usage errors, packet/profile read failures mapping to
`cannot_verify`, preview output not implying evidence production, package
boundary, dependency direction, and MI baselines. It intentionally excludes
`pr-review synthesize`, validate, summarize, check, shared packet/profile
readers, shared repeated-flag helpers, runner sets, and file helpers so each
review workflow responsibility keeps a separate review trail.

Slice 71 continues numbered `pr_review` cleanup in `cmd/sdp-trace` with
`pr-review synthesize` ledger collation shards (`pr_review_109` through
`pr_review_114`). A two-file boundary was selected after pre-change MI analysis
measured `pr_review_synthesize_command.go` at file MI `73.5` and
`pr_review_synthesis_inputs.go` at file MI `72.8`. The slice moves synthesize
execution and argument parsing into `pr_review_synthesize_command.go`, and
packet/run/existing-ledger input loading into `pr_review_synthesis_inputs.go`.
It preserves mandatory output path validation, artifact-read failures mapping
to `cannot_verify`, optional existing-ledger empty-path behavior, prior ledger
as input rather than authority, durable ledger writes before stdout mirroring,
package boundary, dependency direction, and MI baselines. It intentionally
excludes `pr-review validate`, summarize, check, shared JSON writers, shared
file helpers, packet/profile readers, and runner helpers so later review
workflows can keep separate review trails.

Slice 72 continues numbered `pr_review` cleanup in `cmd/sdp-trace` with
`pr-review validate` artifact-validation shards (`pr_review_115` through
`pr_review_120`). A two-file boundary was selected after pre-change MI analysis
measured `pr_review_validate_command.go` at file MI `72.0` and
`pr_review_validation_inputs.go` at file MI `73.2`. The slice moves validate
execution, argument parsing, output path validation, durable validation writes,
stdout mirroring, and validation-exit mapping into
`pr_review_validate_command.go`, and packet/profile/run-set/ledger input
loading into `pr_review_validation_inputs.go`. It preserves mandatory output
path validation, positional-argument rejection, packet/profile/run-set/ledger
read failures mapping to `cannot_verify`, validation verdict persistence before
stdout mirroring, package-owned validation verdicts, `reviewValidationExitCode`
mapping through `exitCannotVerify`, package boundary, dependency direction, and
MI baselines. It intentionally excludes `pr-review summarize`, check, shared
JSON writers, shared file helpers, packet/profile shared readers, repeated flag
helpers, runner helpers, and generic validation exit helpers so later review
workflow responsibilities keep separate review trails.

Slice 73 continues numbered `pr_review` cleanup in `cmd/sdp-trace` with
`pr-review summarize` human-readable summary shards (`pr_review_121` through
`pr_review_125`). A two-file boundary was selected after pre-change MI analysis
measured `pr_review_summarize_command.go` at file MI `73.5` and
`pr_review_summary_io.go` at file MI `76.6`. The slice moves summarize
execution and argument parsing into `pr_review_summarize_command.go`, and
validation/ledger input loading plus optional summary-file writing into
`pr_review_summary_io.go`. It preserves summary text as UX-only output rather
than proof or approval, validation/ledger read failures mapping to
`cannot_verify`, positional-argument rejection, optional output path behavior,
write-once refusal for existing summary files, stdout mirroring even when a
durable summary file is requested, package boundary, dependency direction, and
MI baselines. It intentionally excludes `pr-review check`, shared JSON writers,
generic validation exit helpers, shared file helpers, packet/profile readers,
repeated flag helpers, and runner helpers so later workflow responsibilities
keep separate review trails.

Slice 74 continues numbered `pr_review` cleanup in `cmd/sdp-trace` with
`pr-review check` one-shot review workflow shards (`pr_review_126` through
`pr_review_136`). A five-file boundary was selected after an initial
three-file consolidation failed the file MI gate (`pr_review_check_command.go`
at `69.7` and `pr_review_check_publication.go` at `66.7`). The slice moves
command orchestration into `pr_review_check_command.go`, flag parsing and
required inputs into `pr_review_check_args.go`, packet/profile preparation and
runner execution into `pr_review_check_workflow.go`, preview/summary
publication and exit mapping into `pr_review_check_publication.go`, and durable
artifact writes into `pr_review_check_artifacts.go`. It preserves flag-only
parsing, required `--out` and packet anchors, packet/profile/readiness failures
mapping to `cannot_verify`, work-dir directory validation, repeated
allowed-runner reconstruction from raw args, preview output as non-persisted
planning data, run-set persistence before ledger and validation publication,
summary text only after durable artifacts, validation verdict exit mapping
through `exitCannotVerify`, package boundary, dependency direction, and MI
baselines. It intentionally excludes shared JSON pretty printing, shared file
helpers, packet/profile shared readers, repeated flag helpers, runner sets,
packet-dir helpers, and exit-code helpers so shared utilities can keep separate
review trails.

Slice 75 completes the remaining numbered shared helper cleanup in
`cmd/sdp-trace` with `pr_review_137`, `pr_review_139`, `pr_review_142`, and
`pr_review_144` through `pr_review_149`. `writeIndentedPayload` is treated as a
generic CLI JSON output helper, not as a `pr-review`-only helper, because it is
also used by protected gate output. A single catch-all helper file is rejected
because these helpers serve different trust boundaries: terminal JSON rendering,
write-once and work-dir safety, validation-exit mapping, packet/profile loading,
repeated flag reconstruction, runner allow-list normalization, and
packet-directory derivation. The slice renames them into cohesive shared
locality files while preserving call sites and behavior: JSON payloads remain
stdout copies only for both `pr-review` and protected gate callers, output files
remain write-once, work-dir validation keeps `work-dir:` diagnostics,
validation `cannot_verify` and `coverage_unresolved` still map to
`exitCannotVerify`, packet/profile loading keeps packet errors first and avoids
partial input mixing, repeated raw flags preserve order and parsed fallback
semantics, runner allow-lists ignore empty comma entries without creating
wildcards, packet-dir derivation accepts either a directory or `packet.json`
path, package boundary, dependency direction, and MI baselines. It intentionally
excludes `internal/packet`, release, witness, wrap, query, and numbered gate
family cleanup so later families keep separate review trails.

Slice 76 starts the numbered `internal/packet` cleanup with the packet contract
data model and catalog shards (`packet_001` through `packet_032`). The slice is
limited to constants, enum/catalog maps, core packet/bundle structs, GitHub PR
input structs, prompt boundary and prompt-boundary classification structs,
integration-action structs, build result struct, and the validation result
type. A single `packet.go` catch-all is rejected because it
would mix catalog policy, core packet JSON shape, bundle manifest shape, and
GitHub source input shape into one review surface. The slice instead plans
cohesive locality files for contract constants/catalogs, core packet types,
bundle manifest types, GitHub source input types, and validation result types.
It preserves schema-version values, required row and decision ordering, known
state/catalog membership, JSON field names and `omitempty` behavior, source
input shape for GitHub PR evidence, package boundary, dependency direction, and
MI baselines. It intentionally excludes packet validation, GitHub bundle
building, prompt-boundary classification behavior, rendering, digesting,
loading, and later numbered packet files so behavioral packet responsibilities
keep separate review trails.

Slice 77 continues `internal/packet` cleanup with packet entrypoint and prompt
boundary behavior shards (`packet_033` through `packet_046`). The slice is
limited to validator context structs, JSON loading entrypoints, GitHub bundle
assembly entrypoint helpers, packet shell construction, prompt-boundary theater
finding append, bundle manifest construction, prompt-boundary classification,
metadata completeness checks, text contamination checks, and the recorder-duty
phrase catalog. It preserves loader error behavior, generated packet/bundle IDs,
generated-at UTC formatting, default packet profile fields, integration-action
extension behavior, prompt-boundary verdicts/reasons/route proof effects,
contamination theater finding shape, bundle manifest digesting, package
boundary, dependency direction, and MI baselines. It intentionally excludes
rendering (`packet_047` through `packet_058`), packet digesting (`packet_059`),
validation execution (`packet_060` onward), and downstream GitHub row/entry
helpers (`packet_148` onward) so those behavior-heavy responsibilities keep
separate review trails.

Slice 78 continues `internal/packet` cleanup with packet rendering and digest
helper shards (`packet_047` through `packet_059`). The slice is limited to
clean theater rendering, theater finding rendering, row lookup, decision table
rendering, evidence table rendering, residual-gap rendering, non-proof text,
required-row ordering lookup, resolver fallback lookup, Markdown cell escaping,
and packet digest generation. It preserves rendered section headers, table
column order, markdown escaping for pipes and newlines, `none` rendering for
blank cells, clean-theater fallback row shape, resolver fallback behavior,
required row ordering semantics, non-approval fallback text, digest prefix and
determinism, package boundary, dependency direction, and MI baselines. It
intentionally excludes top-level markdown orchestration (`packet_193` onward),
validation execution (`packet_060` onward), and downstream GitHub row/entry
helpers (`packet_148` onward) so those behavior-heavy responsibilities keep
separate review trails.

Slice 79 continues `internal/packet` cleanup with packet validation entrypoint
and demo-first gate shards (`packet_060` through `packet_084`). The slice is
limited to the public validation entrypoint, demo-first validation entrypoint,
demo-first checker orchestration, row/manifest indexing, tool-generated
requirement, pass-or-partial row count requirement, retained row evidence
requirement, retained structured OpenCode/GSD/MiniMax route evidence
requirement, assessed verification-or-review requirement, cannot-verify closure
cap, demo-usable manifest entry helpers, route component matching, and
synthetic digest rejection. It preserves general `Validate` delegation through
`bundleValidator`, demo gate error accumulation from base validation, required
demo row and evidence semantics, pass/partial/assessed state semantics,
closure-path counting, route evidence source/kind/component/digest
requirements, package boundary, dependency direction, and MI baselines. It
intentionally excludes general bundle metadata, manifest, row, finding, gap,
decision-owner, and shared validation helper shards (`packet_085` onward) so
the larger `bundleValidator` responsibility can keep a separate review trail.

Slice 80 continues `internal/packet` cleanup with general packet validator
orchestration, metadata validation, and manifest/resolver indexing shards
(`packet_085` through `packet_101`). The slice is limited to `bundleValidator`
orchestration, schema metadata validation, packet/bundle identity checks,
required string checks, packet digest validation, packet state and authoring
method catalog checks, projection metadata checks, manifest entry indexing,
resolver entry indexing, and manifest retained-form/redaction-status enum
validation. It preserves validation phase order, error accumulation, packet
digest mismatch behavior, canonical/non-canonical projection semantics,
manifest entry empty-ref rejection, manifest resolver fallback and override
semantics, empty resolver-ref ignoring, manifest enum diagnostics, package
boundary, dependency direction, and MI baselines. It intentionally excludes row
validation and contradiction attribution (`packet_102` through `packet_121`),
finding/gap/decision owner validation (`packet_122` onward), and shared
evidence-ref artifact usability helpers (`packet_145` onward) so those
behavior-heavy responsibilities keep separate review trails.

Slice 81 continues `internal/packet` cleanup with row validation,
contradiction attribution, and row/ref/gap lookup shards (`packet_102` through
`packet_121` and `packet_138` through `packet_144`). The slice is limited to
row indexing, required row presence checks, row ID validation, row required
field validation, row reason/state/summary/owner validation, row evidence-ref
validation handoff, pass-row evidence usability handoff, contradiction target
selection, contradiction state/gap checks, required-row precedence for evidence
ref attribution, extension-row fallback ordering, exact row-ref matching, and
residual gap lookup by reason. It preserves row validation order, duplicate and
unknown row diagnostics, missing required row diagnostics, pass row evidence
requirements, resolver-backed evidence-ref diagnostics, pass-row expired and
unverifiable artifact diagnostics, contradiction row precedence and fallback
semantics, contradiction partial-state and residual-gap requirements, gap
reason semantics, package boundary, dependency direction, and MI baselines. It
intentionally excludes finding/gap/decision-owner validation (`packet_122`
through `packet_136`), shared validator error accumulation (`packet_137`),
shared artifact usability helper implementations
(`packet_145` through `packet_147`), GitHub bundle construction helpers
(`packet_148` through `packet_192`), and rendering (`packet_193` onward) so those
responsibilities keep separate review trails.

Slice 82 continues `internal/packet` cleanup with finding/gap/decision-owner
validation and shared validator error accumulation shards (`packet_122` through
`packet_137`). The slice is limited to theater finding validation, theater row
state validation, residual gap validation, residual coverage enforcement,
decision-owner indexing and required-decision checks, named decision-owner
field validation, and `bundleValidator.add` error accumulation. It preserves
validation phase order through the existing `bundleValidator.validate`
orchestration, theater reason-code diagnostics, trigger evidence-ref validation
handoff, residual gap unknown-row and missing-reason diagnostics, residual
coverage exemptions for `PC-RESIDUAL-GAPS` and pass rows, allowed theater row
states when findings are present, required decision owner ordering, decision
owner trimming semantics, duplicate decision-owner last-valid-owner-wins
overwrite behavior, decision owner state/reason diagnostics, accumulated error
formatting, package boundary, dependency direction, and MI baselines. It
intentionally excludes shared artifact usability helpers (`packet_145` through
`packet_147`), GitHub bundle construction helpers (`packet_148` through
`packet_192`), rendering (`packet_193` onward), and `internal/prreview` numbered
files so those responsibilities keep separate review trails.

Slice 83 continues `internal/packet` cleanup with shared artifact usability
helper shards (`packet_145` through `packet_147`). The slice is limited to
entry expiry checks, pass-evidence unverifiable checks, and artifact access
unverifiable classification. It folds those helpers into a cohesive artifact
evidence usability locality instead of creating standalone one-function
replacement files, because these helpers are shared by pass-row evidence
validation and the demo-first evidence gate. It preserves blank and whitespace-only
`expires_at` non-expiry behavior, malformed timestamp expiry behavior, RFC3339
expiry comparison semantics, `redaction_status=cannot_verify` and
`retained_form=not_retained` unverifiable behavior, artifact access
classification for `expired`, `inaccessible`, `malformed`, `not_assessed`, and
`cannot_verify`, default artifact access pass-through behavior, demo-first row
evidence and retained route evidence usability semantics, package boundary,
dependency direction, and MI baselines. It intentionally excludes GitHub bundle construction helpers
(`packet_148` through `packet_192`), rendering (`packet_193` onward), and
`internal/prreview` numbered files so those responsibilities keep separate
review trails.

Slice 84 continues `internal/packet` cleanup with GitHub source-change and row
construction shards (`packet_148` through `packet_173`). The slice is limited
to projecting GitHub PR input into source-change metadata, packet rows,
prompt-boundary route rows, verification rows, retained-artifact row evidence,
review rows, and the shared GitHub row constructor. It must preserve row IDs,
states, summaries, evidence refs, reasons, owner assignment, prompt-boundary
classification behavior, workflow-run wording, artifact-ref de-duplication,
retained artifact filtering, review pass/partial/not-assessed behavior, package
boundary, dependency direction, and MI baselines. It should fold the numbered
helpers into cohesive GitHub row/source locality files instead of creating
standalone one-function replacements. It intentionally excludes GitHub bundle
manifest entry helpers (`packet_174` through `packet_192`), rendering
(`packet_193` onward), and `internal/prreview` numbered files so those
responsibilities keep separate review trails.

## Verification

```text
gofmt -w <changed-go-files>
go test ./...
go vet ./...
golangci-lint run # when available
go run ./tools/doccheck
go run ./tools/hygienecheck
jq empty schema/*.json
git diff --check
```

CRAP and MI gates are required before any PR claim. If a consolidated file
creates a new MI-baseline entry or stale ratchet behavior, split the slice more
cohesively or move baseline changes into a separate reviewed PR.
