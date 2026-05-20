# Slice 4 Review Context: WS-017-C/D Policy and Supply-Chain Prototypes

## Files Changed
- `examples/oss-policy/adapter.rego`
- `examples/oss-policy/test-fixture.json`
- `examples/oss-policy/README.md`
- `examples/oss-supply-chain/local-dsse.json`
- `examples/oss-supply-chain/README.md`
- `docs/oss-policy-prototype.md`
- `docs/oss-supply-chain-prototype.md`

## Review Axes
1. **Quality:** Correctness of prototypes, no overclaim, honest state markers.
2. **UX:** Documentation clarity, command copy-pasteability.
3. **DX:** File organization, example structure.
4. **Security:** No secrets in fixtures, no unsafe commands, no misleading trust claims.

## Verification Commands
```text
go test -count=1 ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

## Reviewer Instructions
- Examine prototypes for hidden overclaim or trust boundary violations.
- Check that all verifier states are canonical (pass, fail, cannot_verify, not_assessed).
- Verify no secrets or private key material in fixtures.
- Report file:line evidence for every actionable finding.
- Categorize each finding by axis (quality/ux/dx/security).
- Output exactly `LGTM` only if zero findings across all axes.
