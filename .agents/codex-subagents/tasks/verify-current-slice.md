Run the full sdp-trace verification suite on the current branch and report every gate as pass/fail/cannot_verify/not_assessed.

Verification steps:
1. go build ./...
2. go test -count=1 ./... -coverprofile=/tmp/sdp-trace-cover.out
3. go tool cover -func=/tmp/sdp-trace-cover.out > /tmp/sdp-trace-cover-func.txt
4. go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools
5. go run ./tools/qualitycheck -fail-only -function-mi-under 70 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal
6. go run ./tools/qualitycheck -fail-only -mi-under 70 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools
7. go run ./tools/qualitycheck -gocyclo cmd internal tools > /tmp/sdp-trace-gocyclo.txt
8. go run ./tools/crapcheck -cover-func /tmp/sdp-trace-cover-func.txt -gocyclo /tmp/sdp-trace-gocyclo.txt -threshold 5 -strict-less
9. go vet ./...
10. go run ./tools/doccheck
11. git diff --check
12. diff -q .sdp-trace-cmd-surface-before.json <(go run ./cmd/sdp-trace command-surface)
13. diff -q .sdp-trace-help-before.txt <(go run ./cmd/sdp-trace --help)

Report each step with its state.