# Block 23 Quality Report

Date: 2026-05-08

This report is current local evidence for Block 23. It is not external CI proof.

## Command Results

| command | result | evidence |
| --- | --- | --- |
| `go test ./...` | pass | all packages passed; `internal/contract`, `internal/export`, and `internal/policy` have no test files |
| `go vet ./...` | pass | no output |
| `staticcheck ./...` | pass | no output |
| `golangci-lint run ./...` | pass | no output |
| `gofmt -l $(rg --files -g '*.go')` | pass | no output |
| `jq empty schema/*.json` | pass | no output |
| `git diff --check HEAD` | pass | no output |
| `go run ./cmd/sdp-trace --help` | pass | command list includes the Block 23 documented MVP surface |
| `bd ready` | pass | no ready work found; all remaining issues have blocking dependencies |

## Coverage

Coverage was generated with:

```bash
go test -coverprofile=/tmp/sdp-trace-block23-cover.out ./...
go tool cover -func=/tmp/sdp-trace-block23-cover.out
```

Relevant changed-file rows:

| file/function | coverage |
| --- | ---: |
| `internal/releaseproof/releaseproof.go:Evaluate` | 88.2% |
| `internal/releaseproof/releaseproof.go:sourceCommitState` | 80.0% |
| `internal/releaseproof/releaseproof.go:artifactCountsForState` | 100.0% |
| `internal/releaseproof/releaseproof.go:artifactState` | 80.0% |
| `internal/releaseproof/releaseproof.go:combineState` | 100.0% |
| `internal/releaseproof/releaseproof.go:applyDirtyState` | 100.0% |
| `internal/releaseproof/releaseproof.go:artifactCounts` | 100.0% |
| `internal/releaseproof/releaseproof.go:sourceCommitExists` | 100.0% |
| `internal/releaseproof/releaseproof.go:artifactBytes` | 100.0% |
| `internal/releaseproof/releaseproof.go:worktreeDirty` | 83.3% |
| `internal/releaseproof/releaseproof.go:Write` | 71.4% |
| `internal/releaseproof/releaseproof.go:Read` | 71.4% |
| `internal/releaseproof/releaseproof.go:RepoRoot` | 77.8% |
| `internal/releaseproof` package | 87.2% |
| repository total | 66.4% |

The changed production file clears the Block 23 70% changed-file threshold.
Repository-wide 70% coverage is not claimed.

## Complexity

Changed production file check:

```bash
gocyclo -over 15 internal/releaseproof/releaseproof.go
```

Result: pass, no output.

Full repository trust-adjacent scan:

```bash
gocyclo -over 15 $(rg --files -g '*.go' | rg -v '_test\.go$')
```

Result: fail for legacy hotspots outside the changed releaseproof file:

| complexity | function | file |
| ---: | --- | --- |
| 28 | `trace.writeCanonicalJSON` | `internal/trace/event.go` |
| 25 | `witness.BuildCustomerPKI` | `internal/witness/profiles.go` |
| 24 | `forensic.rawReferenceCondition` | `internal/forensic/forensic.go` |
| 23 | `demo.witnessBindingState` | `internal/demo/demo.go` |
| 23 | `main.run` | `cmd/sdp-trace/main.go` |
| 22 | `posture.metricMatches` | `internal/posture/posture.go` |
| 19 | `adaptercapture.runBindingCondition` | `internal/adaptercapture/adaptercapture.go` |
| 19 | `posture.Build` | `internal/posture/posture.go` |
| 19 | `main.(*flagSet).parse` | `cmd/sdp-trace/main.go` |
| 18 | `adaptercapture.overclaimCondition` | `internal/adaptercapture/adaptercapture.go` |
| 18 | `recorder.Run` | `internal/recorder/recorder.go` |
| 18 | `main.runGateExplain` | `cmd/sdp-trace/main.go` |
| 17 | `managed.witnessCondition` | `internal/managed/managed.go` |
| 17 | `main.witnessMatchesProtectedInput` | `cmd/sdp-trace/main.go` |
| 16 | `witness.validateCIEnvelope` | `internal/witness/profiles.go` |

These are legacy exceptions for Block 23. They are not changed production
functions and are not claimed CRAP-clean.

## CRAP Rows

Formula: `CRAP = complexity^2 * (1 - coverage)^3 + complexity`.

| function | complexity | coverage | CRAP | result |
| --- | ---: | ---: | ---: | --- |
| `Evaluate` | 4 | 0.882 | 4.03 | pass |
| `sourceCommitState` | 3 | 0.800 | 3.07 | pass |
| `artifactCountsForState` | 2 | 1.000 | 2.00 | pass |
| `artifactState` | 3 | 0.800 | 3.07 | pass |
| `combineState` | 3 | 1.000 | 3.00 | pass |
| `applyDirtyState` | 3 | 1.000 | 3.00 | pass |
| `Write` | 3 | 0.714 | 3.21 | pass |
| `Read` | 3 | 0.714 | 3.21 | pass |
| `RepoRoot` | 3 | 0.778 | 3.10 | pass |

Block 23 changed releaseproof functions are below CRAP 5. Repository-wide CRAP
is not assessed.

## Dead Code And Parked Packages

Command:

```bash
go run golang.org/x/tools/cmd/deadcode@latest ./...
```

Result: deadcode reported unreachable functions in several packages, including
the Block 23 named staged packages:

| package | deadcode result | disposition |
| --- | --- | --- |
| `internal/contract` | all exported functions reported unreachable | staged non-MVP package; not current closure evidence |
| `internal/export` | all exported functions reported unreachable | staged non-MVP package; not current closure evidence |
| `internal/policy` | policy loader/validator functions reported unreachable | staged non-MVP package; not current closure evidence |

`rg -n 'internal/(contract|export|policy)' --glob '*.go' .` returned no import
call sites. These packages remain explicit quality exceptions, not green proof.

## Exception Rows

| id | scope | metric | observed | threshold | reason | owner | follow-up |
| --- | --- | --- | --- | --- | --- | --- | --- |
| B23-Q-01 | repository-wide trust-adjacent complexity | `gocyclo` | 15 legacy functions over 15 | no changed file over 15 | Block 23 changed releaseproof only; legacy hotspots need separate design/testing slices | `role:sdp-trace-maintainer` | future quality-hardening block |
| B23-Q-02 | `internal/contract`, `internal/export`, `internal/policy` | `coverage/deadcode` | 0% coverage and unreachable by `deadcode` | current closure evidence must be reachable or explicitly parked | staged non-MVP packages are not part of current MVP proof | `role:sdp-trace-maintainer` | decide remove vs wire vs document in a later block |
| B23-Q-03 | repository-wide CRAP | `crap` | not measured for every production function | `<5` | Block 23 measured changed releaseproof functions only | `role:sdp-trace-maintainer` | future repo-wide quality gate |
