# Block 17 Managed Harness Fixtures

These fixtures are committed examples for the managed harness assessment output
shape. They are not external audit proof and do not make `sdp-trace` a policy
decision engine.

| Scenario | Fixture | Evidence |
|---|---|---|
| Valid managed profile assessment | `valid-managed-profile.assessment-result.json` | Committed JSON fixture plus Go domain/CLI tests. |
| Unmanaged run fails | `unmanaged-run-fail.assessment-result.json` | Fixture generated from `internal/managed.Evaluate`; checked by `TestBlock17CommittedFixturesHaveManagedAssessmentShape`. |
| Late enrollment fails | `late-enrollment-fail.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Post-hoc policy fails | `post-hoc-policy-fail.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Post-hoc registry fails | `post-hoc-registry-fail.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Unauthorized adapter fails | `unauthorized-adapter-fail.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Adapter disconnect fails | `adapter-disconnect-fail.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Missing capability cannot verify | `missing-capability-cannot-verify.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Missing harness event cannot verify | `missing-harness-event-cannot-verify.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Missing tool event cannot verify | `missing-tool-event-cannot-verify.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Missing file mutation event cannot verify | `missing-file-mutation-event-cannot-verify.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Missing test telemetry cannot verify | `missing-test-telemetry-cannot-verify.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Agent-reported test evidence fails | `agent-reported-test-evidence-fail.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Policy-authorized suppression passes | `policy-authorized-suppression-pass.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Suppression without policy fails | `suppression-without-policy-fail.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Witness missing cannot verify | `witness-missing-cannot-verify.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Stale witness cannot verify | `stale-witness-cannot-verify.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Witness mismatch fails | `witness-mismatch-fail.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Override present does not upgrade | `override-present-non-upgrading-pass.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
| Override upgrade attempt fails | `override-upgrade-fail.assessment-result.json` | Fixture generated from evaluator and checked by fixture matrix test. |
