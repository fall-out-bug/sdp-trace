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

## Verification

```text
gofmt -w <changed-go-files>
go test ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

CRAP and MI gates are required before any PR claim. If a consolidated file
creates a new MI-baseline entry or stale ratchet behavior, split the slice more
cohesively or move baseline changes into a separate reviewed PR.
