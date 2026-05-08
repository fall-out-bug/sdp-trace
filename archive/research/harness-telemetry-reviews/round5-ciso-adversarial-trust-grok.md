```
VERDICT: CHANGES_REQUIRED

Critical blockers:
- No explicit requirement for signing key separation from agent workspace (signing and verification; observer authority). Agent/harness observer must not access keys for local_observed or higher scopes.
- MVP lacks real external witness (CI/gate anchoring; demo credibility). Simulated CI labeled "not gate-grade" is fine, but must mandate OIDC path with non-local signing.

Major gaps:
- Replay prevention weak: nonce per event/checkpoint insufficient without timestamp authority or Rekor bundle (trace model; signing and verification).
- Verifier lacks check for run deletion/restart: no binding to workspace state or VCS ref at start (provenance model; CI/gate anchoring).
- Quiet telemetry drop possible if agent skips events without explicit missing_telemetry fact signed by authorized observer (trace model; expected evidence contracts).

False assumptions:
- Local hash chains + checkpoints suffice pre-CI; attacker recomputes full chain post-facto (evidence model; observer authority).
- Agent self-report + local_observed upgrades to harness_observed without distinct signer (product layering; signing and verification).

Standards we should reuse instead of building:
- SLSA Provenance for run start/end binding to VCS/workspace snapshot (provenance model).
- Sigstore Fulcio/DSSE with Rekor transparency for all checkpoints, not optional (signing and verification).

Minimum viable product correction:
- Add keyless OIDC signing fixture in CI path only; local verifier simulates but flags as local_only (signing and verification; demo credibility).
- Mandate explicit "run_start" event with VCS ref digest, signed by harness or external (trace model; CI/gate anchoring).

Questions before implementation:
- How does verifier enforce signer cannot access payload during signing for local_observed? (observer authority)
- What binds trace to non-resettable workspace state? (provenance model; privacy and retention)

Attack scenarios not covered:
- Developer deletes workspace, regenerates identical trace with new timestamps/nonces post-CI (trace model; CI/gate anchoring).
- Compromised harness signs fabricated local_observed events as ci_witnessed (observer authority; signing and verification).
- Agent drops model/gateway events, signs agent_reported chain as complete (expected evidence contracts; product layering).

Demo changes required:
- Demo 5: Show attacker replaying old trace to new PR; verifier fails on VCS ref mismatch or timestamp (demo credibility).
- All demos: Explicit verifier output with "key_accessible_to_agent: true/false" and "replay_risk: high/medium/low" (CTO usefulness; verifier states).
```
