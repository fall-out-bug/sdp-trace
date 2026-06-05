# Slice 68 Plan Review

Status: in_progress
Date: 2026-06-04

## Scope

Slice 68 is bounded to remaining numbered `cmd/sdp-trace` packet fixture and
validation exit helper shards:

- `cmd/sdp-trace/packet_093_prfixtureevent.go`
- `cmd/sdp-trace/packet_094_loadprfixtureevent.go`
- `cmd/sdp-trace/packet_095_readoptionaljson.go`
- `cmd/sdp-trace/packet_096_exits.go`

## Decision Gate

- Simpler/Faster: rename/move only; no behavior change, public API change,
  dependency change, or baseline update.
- Blocking Edge Cases: a single fixture IO file failed pre-change file MI at
  `66.3`; fixture identity validation and exit mapping affect trust verdict
  reporting and need focused regression evidence.
- Existing Open Source: no new JSON helper, CLI framework, or validation
  library is needed; Go standard library and existing package-local helpers are
  sufficient.

## Planned File Boundary

- `cmd/sdp-trace/packet_build_pr_fixture_event.go`: PR fixture event shape and
  fixture identity validation.
- `cmd/sdp-trace/packet_build_pr_optional_json.go`: shared optional JSON
  loading.
- `cmd/sdp-trace/packet_validation_exits.go`: packet validation and demo gate
  exit mapping.

## Planned Regression Evidence

- Exact focused test existence:
  - `TestLoadPRFixtureEventRequiresIdentity`
  - `TestReadOptionalJSONKeepsOptionalAndErrorBehavior`
  - `TestPacketExitMappingsKeepTrustSemantics`
- Focused behavior checks: required PR number and URL validation, optional empty
  path no-op, JSON read and unmarshal error propagation, `pass` to zero exit
  mapping, packet validation failure to `cannot_verify`, and demo gate failure
  to `fail`.
- Standard verification, conditional `golangci-lint run`, and CRAP/MI gates
  remain required before the implementation review.

## Review Rounds

### Round 1

- scope/correctness: `LGTM`
- trust/evidence: minor finding. The plan verification block omitted required
  `jq empty schema/*.json` and conditional `golangci-lint run`, while tasks
  already included them.
- maintainability/DX: major finding. `packet_096` contains packet
  validate/check-demo exit helpers, not `packet build-pr` behavior, so
  `packet_build_pr_exits.go` would encode the wrong ownership.

### Round 2

- trust/evidence: `LGTM`
- maintainability/DX: minor finding. T021-4681 still called the regression
  lane focused `packet build-pr` evidence even though it now covers packet
  validate/check-demo exit mapping.

### Round 3

- maintainability/DX: `LGTM`

Final plan-review status: `LGTM` across scope/correctness, trust/evidence,
and maintainability/DX after fixes.
