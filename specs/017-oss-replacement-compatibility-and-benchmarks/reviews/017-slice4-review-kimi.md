# Slice 4 Review: WS-017-C/D Policy and Supply-Chain Prototypes

Reviewer: Kimi (adversarial)
Date: 2026-05-20

## Review

### Correct
- **Canonical verifier states:** Probe result tables in `docs/oss-policy-prototype.md` and `docs/oss-supply-chain-prototype.md` use only canonical states (`pass`, `fail`, `cannot_verify`, `not_assessed`). No custom or misleading verdict labels in probe tables.
- **Honest state markers, no overclaim:** Both docs explicitly label prototypes as "local experiments only" and include clear Substitution Boundary and Non-Goals sections. No production migration claims.
- **No secrets in fixtures:** `test-fixture.json` contains only synthetic test data. `local-dsse.json` contains a truncated signature (`"sig": "MEUCIQDtVdxRBBT..."`) and a test-only `keyid` (`local-test-key`). Decoded payload is harmless test metadata.
- **Valid syntax:** `adapter.rego` uses correct OPA v1 syntax (`import rego.v1`, `pass if { ... }`, `fail_reason contains ... if { ... }`). JSON fixtures validate with `jq empty`.
- **Verification commands pass:** `go test -count=1 ./...`, `go vet ./...`, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, and `git diff --check` all exit 0.

### Fixed
- **Important / UX:** `examples/oss-supply-chain/README.md:47` — the `slsa-verifier` command used `--provenance-path examples/oss-supply-chain/local-dsse.json`, which fails when the reader runs the command from the `examples/oss-supply-chain/` directory (the natural working directory when following the README). Changed to `--provenance-path local-dsse.json` so the command is copy-pasteable from the README's own directory.

### Blocker
- None.

### Note
- **Advisory / UX:** The Cosign local blob verify command (`examples/oss-supply-chain/README.md:36`) may fail on newer Cosign versions with default Rekor/transparency-log settings unless the user disables the transparency log or passes `--insecure-ignore-tlog`. The doc (`docs/oss-supply-chain-prototype.md`) already notes "when transparency log is disabled," so this is documented at the spec level, but the README command itself does not include the flag. Consider adding a note inline if the target audience runs recent Cosign.
- **Advisory / DX:** The `in-toto-run` example (`examples/oss-supply-chain/README.md:26`) references `/tmp/test-key.pem` without showing key generation. Since `in-toto-run` is listed as a prerequisite, this is acceptable for a prototype, but a one-liner key-generation hint would improve copy-pasteability.
