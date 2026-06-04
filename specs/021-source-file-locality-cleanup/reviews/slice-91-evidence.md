# Slice 91 Evidence

Date: 2026-06-04

Scope: `internal/prreview` validation orchestration, digest validation, and
required-plane selection shards `prreview_049` through `prreview_060`.

## Locality

- Removed numbered validation orchestration shards `prreview_049` through
  `prreview_060`.
- Moved public validation orchestration, authority-bound result construction,
  and role lookup into `validation.go`.
- Moved packet/run/ledger digest validation into `validation_digest.go`.
- Moved required-plane set construction and sorted iteration into
  `validation_required_planes.go`.
- Moved plane validation note helpers into `validation_plane_notes.go`.
- Moved best plane-result selection into `validation_plane_selection.go`.
- Plane ranking/model identity, ledger finding validation, coverage-state,
  summary rendering, file IO, prompt generation, and reviewer execution are
  intentionally excluded from this slice.

## Source Shape

```text
$ find internal/prreview -maxdepth 1 -type f | sed 's#^#/#' | rg '/prreview_0(4[9]|5[0-9]|60)_[^/]+\.go$' || true
<no output>

$ git diff --name-only | rg '^internal/prreview/prreview_(0(6[1-9]|7[0-9]|8[0-9]|9[0-9])|1[0-9][0-9])_' || true
<no output>
```

## Verification

`pass`:

```text
gofmt -w internal/prreview/*.go
tests='TestValidateReviewStatesAndAuthorityBoundary|TestValidationAndLedgerLifecyclePreserveTrustSemantics|TestValidateUsesBestPlaneResultAcrossRetries|TestValidateCannotVerifyPerResultPacketDigestMismatch|TestValidateCoverageStatesForNoReviewersUnresolvedAndStaleDigest|TestPrreviewValidationOrchestrationPreservesDigestRequiredPlaneAndAuthorityContracts'; listed=$(go test ./internal/prreview -list "$tests" | rg "^($tests)$"); test "$(printf '%s\n' "$listed" | sed '/^$/d' | wc -l | tr -d ' ')" = 6 && go test ./internal/prreview -run "$tests" -count=1 -v
go test ./internal/prreview
go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal
go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools
```

`pass`:

```text
go test ./...
go vet ./...
golangci-lint run
go run ./tools/doccheck && go run ./tools/hygienecheck && jq empty schema/*.json && git diff --check
go test -count=1 ./... -coverprofile=coverage.out
go tool cover -func=coverage.out > coverage-func.txt
go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt
go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less
```

Temporary files `coverage.out`, `coverage-func.txt`, and `gocyclo.txt` were
removed after CRAP verification.

## Reviewer Rounds

- Round 1 Lane A behavior/correctness: Beauvoir `LGTM`.
- Round 1 Lane B locality/boundary/MI: Halley `LGTM`.
- Round 1 Lane C tests/evidence: Peirce reported two major test-evidence gaps.
  The focused Slice 91 regression did not assert the missing required-plane
  fallback `PlaneResult` shape and did not directly prove plane note
  aggregation for non-usable plane results.
- Fix: `TestPrreviewValidationOrchestrationPreservesDigestRequiredPlaneAndAuthorityContracts`
  now includes an unassessed required security plane and a parse-failed privacy
  plane, asserting sorted fallback shape, `<plane>:<status>` reasons, and
  expected next actions.
- Round 2 Lane C tests/evidence: Peirce re-review returned exactly `LGTM`
  after the focused regression asserted missing required-plane fallback shape
  and plane note aggregation.
