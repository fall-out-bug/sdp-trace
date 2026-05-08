# Forensic Retention Fixtures

These fixtures are generated from `internal/forensic.Evaluate` and checked by
`TestBlock18CommittedFixturesHaveForensicAssessmentShape`.

They cover the assessment-result shape for explicit
`assess --profile forensic-retention` use. They are verifier facts only:
they do not decide legal, incident, readiness, merge, release, or risk
acceptance outcomes.

| Case | Fixture |
| --- | --- |
| Valid forensic retention assessment | `valid-forensic-retention.assessment-result.json` |
| Digest-only critical evidence fails with cap | `digest-only-critical-fail.assessment-result.json` |
| External access present but unverifiable cannot verify | `external-access-unverifiable-cannot-verify.assessment-result.json` |
| Weak raw-reference digest fails | `weak-digest-fail.assessment-result.json` |
| Self-claimed redaction authority cannot verify | `authority-self-claim-cannot-verify.assessment-result.json` |
| Withheld evidence maps to `not_assessed` and fails forensic retention | `withhold-not-assessed-fail.assessment-result.json` |
| Missing redaction policy cannot verify | `missing-policy-cannot-verify.assessment-result.json` |
