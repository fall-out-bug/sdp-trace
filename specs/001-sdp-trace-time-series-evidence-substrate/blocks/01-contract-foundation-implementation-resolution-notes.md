# Block 01 Implementation Critic Resolution Notes

Critic: MiniMax via `pi`, no tools, no repository context beyond attached files.

## Accepted and Fixed

- `examples/contract-foundation/not-assessed-assessment-input.json`: added the missing `sample-log-availability-previous` sample so the comparison references resolve structurally.
- `examples/contract-foundation/positive-assessment-input.json` and `examples/contract-foundation/not-assessed-assessment-input.json`: removed circular `evidence_bundle_refs` that pointed an assessment input back to itself.
- `schema/contract-release-verification.schema.json` and `examples/contract-foundation/contract-release-verification.example.json`: added and used `transparency_log_status: "compensating_control_recorded"` for the local private-equivalent verification record.
- `docs/agent-onboarding.md` and `docs/agent-onboarding.md`: added direct pointers to SpecKit, schema, fixture, and validation evidence locations.

## Rejected With Evidence

- Missing schema/script/DSSE/public-key findings are rejected. The files exist in the checkout and local verification commands read them successfully: `schema/consumer-schema-version-declaration.schema.json`, `scripts/verify-release-signature.sh`, `scripts/generate-local-dev-dsse.sh`, `examples/contract-foundation/contract-release.dsse.json`, and `examples/contract-foundation/local-dev-signing-public.pem`.
- The `trusted_contract_release: true` plus `identity_policy_status: "mismatch"` concern is rejected as an implementation defect. `schema/contract-release-verification.schema.json` constrains trusted releases to matched identity policy, and `scripts/validate-contracts.sh` requires `examples/contract-foundation/negative-unauthorized-signer.json` to fail validation.
- The AI-only accountability concern is rejected as an implementation defect. `schema/accountability.schema.json` resolves to accountable actor types `human_user`, `human_role`, and `human_group`; `scripts/validate-contracts.sh` requires `examples/contract-foundation/negative-ai-sole-accountable-owner.json` to fail validation.

## Residual Risk

- The positive manifest uses `source_commit: "working-tree-block-01"` because this verification was produced before a commit exists. This is acceptable for the local development private-equivalent proof, but a production trusted release must use an immutable source commit or signed source reference.
- `docs/process-metric-catalog.md` remains open as T011 and is not claimed complete in Block 01.
