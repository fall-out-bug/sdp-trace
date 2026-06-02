# Slice 54 Plan Review

Status: pass

## Scope

Slice 54 is bounded to override request CLI handling shards:

- `cmd/sdp-trace/gate_352_runoverride.go`
- `cmd/sdp-trace/gate_353_appendoverriderequestevent.go`
- `cmd/sdp-trace/gate_354_overriderequestpayload.go`
- `cmd/sdp-trace/gate_355_parseoverriderequestargs.go`
- `cmd/sdp-trace/gate_356_isoverriderequest.go`
- `cmd/sdp-trace/gate_357_parseoverriderequestflags.go`
- `cmd/sdp-trace/gate_358_requireoverriderequestflags.go`
- `cmd/sdp-trace/gate_359_overriderequestrequiredflags.go`

Initial planned cohesive file:

- `cmd/sdp-trace/override_request.go`

Final implementation split after file-level MI failure:

- `cmd/sdp-trace/override_request.go`
- `cmd/sdp-trace/override_request_payload.go`
- `cmd/sdp-trace/override_request_flags.go`

Explicit exclusions:

- shared JSON/text helpers (`gate_360` onward)
- preview mode and required-ID helper shards (`gate_365` onward)
- gate rendering and demo override extraction outside the CLI request writer

## Behavior To Preserve

- Only `override request` is accepted under the override namespace.
- Parser errors return `exitUsage` and print diagnostics before any trace event
  is appended.
- Missing required flags return `exitUsage` before any trace event is appended.
- Positional text is rejected with `override request accepts only flags`.
- Required flag diagnostics remain stable:
  `--out`, `--id`, `--by`, `--reason`, `--source-ref`, and `--scope`.
- Successful requests append `trace.EventPolicyOverrideRequested` through
  `trace.AppendRunEvent`.
- Payload keys remain stable: `override_id`, `producer`, `origin`,
  `requested_by`, `reason`, `source_ref`, `scope`, `created_at`, and optional
  `external_reference`.
- `producer` remains `sdp-trace-cli`; `origin` remains `native_cli`;
  `created_at` remains UTC RFC3339Nano.
- Successful stdout remains `override_event: <event-id>`.
- Append failures print stderr and return exit code 1.
- Override events remain advisory/non-upgrading and do not approve a gate or
  policy state by themselves.
- No package boundary, dependency direction, or MI baseline change is planned.

## Planned Regression Evidence

- Existing `TestOverrideRequestAppendsFlightRecorderEvent`.
- Existing `TestGateOutputIncludesOverrideWithoutPassingMissingEvidence`.
- Add `TestOverrideRequestPersistsExternalReferencePayload` for optional
  `--external-reference`, stable payload fields, and UTC RFC3339Nano
  `created_at` parsing.
- Add `TestOverrideRequestRejectsInvalidRequestArgsBeforeAppend` for
  non-`request` override subcommand rejection, unknown-flag parser errors,
  positional text rejection, and missing required flag diagnostics without
  appending events.
- Add `TestOverrideRequestAppendFailureDoesNotPrintEvent` for append failure
  stderr/exit behavior and no `override_event` stdout claim.

Focused test existence and execution will use:

```text
go test ./cmd/sdp-trace -list '^(TestOverrideRequestAppendsFlightRecorderEvent|TestOverrideRequestPersistsExternalReferencePayload|TestOverrideRequestRejectsInvalidRequestArgsBeforeAppend|TestOverrideRequestAppendFailureDoesNotPrintEvent|TestGateOutputIncludesOverrideWithoutPassingMissingEvidence)$'
go test ./cmd/sdp-trace -run 'Test(OverrideRequestAppendsFlightRecorderEvent|OverrideRequestPersistsExternalReferencePayload|OverrideRequestRejectsInvalidRequestArgsBeforeAppend|OverrideRequestAppendFailureDoesNotPrintEvent|GateOutputIncludesOverrideWithoutPassingMissingEvidence)$'
```

## Plan Review Findings

- trust/evidence lane initial finding: major; unknown-flag parser-error case was
  not explicitly planned. Resolution: added
  `TestOverrideRequestRejectsInvalidRequestArgsBeforeAppend` coverage for
  unknown flags before append.
- trust/evidence lane initial finding: major; exact test existence was required
  but new test names and commands were not planned. Resolution: named all
  focused tests and commands above.
- maintainability/DX lane initial finding: major; exact test names were missing.
  Resolution: named all focused tests and commands above.
- maintainability/DX lane initial finding: minor; Slice 54 task block was
  inserted before Slice 48. Resolution: moved Slice 54 after Slice 53.

## Review Lanes

- scope/correctness reviewer: multi_agent_v1,
  `019e8858-c8b3-7063-a2eb-192f2a8bbc77`; date:
  `2026-06-02T16:40:16+03:00`; prompt class:
  `plan scope/correctness review`; timeout: `600000ms`; retries: `0`;
  fallback: `none`; result `LGTM`
- trust/evidence reviewer: multi_agent_v1,
  `019e8858-ccec-7211-9d43-eaf682f92e18`; date:
  `2026-06-02T16:40:16+03:00`; prompt class:
  `plan trust/evidence review plus targeted re-review`; timeout: `600000ms`;
  retries: `1`; fallback: `none`; result `LGTM`
- maintainability/DX reviewer: multi_agent_v1,
  `019e8858-d21f-7ad0-a8be-8d1e3e48dbcc`; date:
  `2026-06-02T16:40:16+03:00`; prompt class:
  `plan maintainability/DX review plus targeted re-review`; timeout:
  `600000ms`; retries: `1`; fallback: `none`; result `LGTM`
- model/provider: `not_assessed` for all lanes because the harness does not
  expose provider-qualified model IDs
- requested external/provider-qualified lanes: `not_assessed` because no
  callable provider-qualified model surface is exposed in this session
