# Slice 81 Evidence

Date: 2026-06-04

Scope: `internal/packet` row validation, contradiction attribution, and
row/ref/gap lookup shards `packet_102` through `packet_121` and `packet_138`
through `packet_144`.

Prompt class: staged implementation review after local verification.

## Implementation Summary

- Removed numbered row validation, contradiction attribution, and row/ref/gap
  lookup files from `internal/packet`.
- Added responsibility-named locality files for row validation, row indexing,
  row field checks, row evidence refs, pass evidence refs, contradiction
  checks, contradiction target selection, required-row lookup, extension-row
  lookup, deterministic row sorting, exact row-ref matching, and gap lookup.
- Added focused regression tests for row diagnostics, contradiction target
  selection, required-row precedence over extension rows for shared refs,
  row-ref attribution ordering, exact row-ref matching, and residual gap reason
  behavior.
- `cmd/sdp-trace/FAMILY_INDEX.md` is not applicable because Slice 81 is scoped
  to the non-command `internal/packet` package.

## Source Shape Evidence

Source-shape command:

```sh
find internal/packet -maxdepth 1 -type f | sed 's#^#/#' | rg '/packet_10[2-9]_[^/]+\.go$|/packet_1[01][0-9]_[^/]+\.go$|/packet_12[0-1]_[^/]+\.go$|/packet_13[8-9]_[^/]+\.go$|/packet_14[0-4]_[^/]+\.go$' || true
```

Result: verified pass, no output.

Replacement locality files:

- `internal/packet/validation_rows.go`
- `internal/packet/validation_row_index.go`
- `internal/packet/validation_row_presence.go`
- `internal/packet/validation_row_id.go`
- `internal/packet/validation_row_fields.go`
- `internal/packet/validation_row_evidence.go`
- `internal/packet/validation_contradictions.go`
- `internal/packet/validation_contradiction_requirements.go`
- `internal/packet/validation_row_lookup.go`
- `internal/packet/validation_extension_row_lookup.go`
- `internal/packet/validation_gap_lookup.go`

Result: verified pass by direct file existence after implementation.

## Focused Verification

- `gofmt -w internal/packet/validation_*.go internal/packet/packet_test.go`
  - verified pass.
- `go test ./internal/packet`
  - verified pass.
- `go test -list 'Test(Validate|RowID|Gap)' ./internal/packet`
  - verified pass.
  - Included new/expanded Slice 81 coverage:
    `TestValidateRowDiagnostics`,
    `TestValidateContradictionTargetSelection`,
    `TestRowIDForRefUsesRequiredRowOrder`,
    `TestRowIDForRefUsesExtensionFallbackOrderAndExactRefs` covering
    required-row precedence over extension rows on a shared ref plus extension
    fallback ordering and exact ref matching, and
    `TestGapForRowRequiresReasonForResidualCoverage`.

## Repository Verification

- `go test ./...`
  - verified pass.
- `go vet ./...`
  - verified pass.
- `golangci-lint run`
  - verified pass.
- `go run ./tools/doccheck && go run ./tools/hygienecheck && jq empty schema/*.json && git diff --check`
  - verified pass.

## Quality Gates

- `go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools`
  - verified pass.
- `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`
  - verified pass.
- `go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`
  - verified pass.
- `go test -count=1 ./... -coverprofile=coverage.out`
  - verified pass.
- `go tool cover -func=coverage.out > coverage-func.txt`
  - verified pass.
- `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`
  - verified pass.
- `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`
  - verified pass.

## Boundary Evidence

Staged boundary command:

```sh
git diff --cached --name-only -- internal/packet/packet_12{2..9}_*.go internal/packet/packet_13{0..7}_*.go internal/packet/packet_14{5..9}_*.go internal/packet/packet_1{50..92}_*.go internal/packet/packet_19{3..9}_*.go internal/packet/packet_200_*.go || true
```

Result: verified pass, no output. Excluded finding/gap/decision-owner
validation, shared validator error accumulation, shared artifact usability
helpers, GitHub bundle construction helpers, and rendering files are not moved
or edited in Slice 81.

## Implementation Review

Round 1:

- Hegel the 2nd
  - Harness: Codex subagent
  - Agent id: `019e93ba-d399-7283-9c5d-3b7a36014149`
  - Model/provider: not_assessed
  - Timeout/retries/fallback: 600000 ms / 0 / none
  - Result: `LGTM`
- Cicero the 2nd
  - Harness: Codex subagent
  - Agent id: `019e93ba-f954-7953-9f20-f7431abac098`
  - Model/provider: not_assessed
  - Timeout/retries/fallback: 600000 ms / 0 / none
  - Result: findings
  - Finding: major, missing focused regression coverage for required-row
    precedence over extension rows when both cite the same ref.
  - Fix: expanded `TestRowIDForRefUsesExtensionFallbackOrderAndExactRefs` so a
    required row and extension rows cite `custom:ref`, verified the required row
    wins, then removed the required row to verify deterministic extension
    fallback and exact matching.
- Godel the 2nd
  - Harness: Codex subagent
  - Agent id: `019e93bb-1c67-77f1-92c2-dc83b317a16e`
  - Model/provider: not_assessed
  - Timeout/retries/fallback: 600000 ms / 0 / none
  - Result: findings
  - Finding: major, the first replacement structure recreated the tiny-file
    antipattern with many 7-to-8-line wrapper files.
  - Fix: consolidated closely coupled helpers into cohesive locality files:
    row field checks, row evidence checks, contradiction targeting,
    contradiction requirements, required-row lookup, extension-row lookup, and
    gap lookup. MI and CRAP gates pass without baseline changes.

Round 2:

- Dewey the 2nd
  - Harness: Codex subagent
  - Agent id: `019e93c1-881a-7e13-86f4-b44092348cf3`
  - Model/provider: not_assessed
  - Timeout/retries/fallback: 600000 ms / 0 / none
  - Result: `LGTM`
- Feynman the 2nd
  - Harness: Codex subagent
  - Agent id: `019e93c1-ab06-71b0-83de-be2bed04e592`
  - Model/provider: not_assessed
  - Timeout/retries/fallback: 600000 ms / 0 / none
  - Result: `LGTM`
- Bernoulli the 2nd
  - Harness: Codex subagent
  - Agent id: `019e93c1-d151-7dd2-bf94-919ac85764fb`
  - Model/provider: not_assessed
  - Timeout/retries/fallback: 600000 ms / 0 / none
  - Result: `LGTM`

Implementation review verdict: pass.
