# Block 24 Demo Repository CI Artifact Index

Authority scope: `demo_pilot_only`

Raw artifacts remain in the private demo repository GitHub Actions artifact
store. This file records sanitized references and selected digests only.

## Run Reference

| field | value |
| --- | --- |
| Demo repository | `fall-out-bug/sdp-trace-demo-ci-pilot` |
| Demo commit | `e370d1c00df8a7e7859adc284480563a269e64ca` |
| Workflow | `sdp-trace demo pilot` |
| Workflow run | `25548285336` |
| Run URL | `https://github.com/fall-out-bug/sdp-trace-demo-ci-pilot/actions/runs/25548285336` |
| `sdp-trace` source ref | `codex/block-24-demo-ci-trace-pilot` |
| Artifact retention | 14 days |

## GitHub Actions Artifacts

| artifact | id | expiration | state |
| --- | --- | --- | --- |
| `sdp-trace-demo-clean-report` | `6876152707` | 2026-05-22T09:34:56Z | `available` at collection time |
| `sdp-trace-demo-no-oidc-report` | `6876147448` | 2026-05-22T09:34:38Z | `available` at collection time |

If either artifact expires before review, the corresponding evidence state
becomes `cannot_verify` until the workflow is rerun.

## Clean Run Summary

| run | result | trust scope | exit | stdout sha256 | stderr sha256 |
| --- | --- | --- | --- | --- | --- |
| `clean-feature-flag` | `observed` | `local_observed` | 0 | `1eb1944e354ef3306bda0ee3d6d2fb1073091f34d82684c77f649f13771b7df2` | `c33931e982426586358825549e0f05a03e9d3553f29dd4d39f0214890f8a1ec6` |
| `clean-entitlement-matrix` | `observed` | `local_observed` | 0 | `b2775d4cd02392656ac0bf2cc6b689f722ef13926045177ffb1bcdbdf2194e64` | `46d08915c993f62225d01d893d69944137a99cb3ec06a49c5ec46bcb167ab3ea` |
| `clean-audit-scope` | `observed` | `local_observed` | 0 | `d83410094344b9afb3b8ddd6b15343d747f6b57e609f1efe71ae8ac6a68a8143` | `fb4cf942a8a671e34e898cbbc147ae802de72cac2bc2900ce9ce6cec98619494` |

`summary.json` reported `run_count=3`, `observed_count=3`,
`cannot_verify_count=0`, `not_assessed_count=0`, and
`trust_scope=local_observed` for the run summary. The separate witness profile
raised the exact CI run to `ci_witnessed`; it does not raise external
production trust.

## Witness Extracts

| path | status | established trust scope | reason | exit |
| --- | --- | --- | --- | --- |
| `ci-witness.json` | `pass` | `ci_witnessed` | `ci_identity_present` | 0 |
| `ci-witness-no-oidc.json` | `cannot_verify` | `cannot_verify` | `missing_ci_oidc` | 3 |

Allowed extracted witness fields:

- `kind=github-actions`
- `source.repository=fall-out-bug/sdp-trace-demo-ci-pilot`
- `source.commit_sha=e370d1c00df8a7e7859adc284480563a269e64ca`
- `ci.run_id=25548285336`
- `ci.run_attempt=1`
- `oidc.subject=repo:fall-out-bug/sdp-trace-demo-ci-pilot:ref:refs/heads/main`
- `profile_states.identity_state=pass` for the clean witness
- `profile_states.identity_state=cannot_verify` for the no-OIDC witness
- `missing_identity_fields=[ACTIONS_ID_TOKEN_REQUEST_TOKEN,ACTIONS_ID_TOKEN_REQUEST_URL]` for the no-OIDC witness

Forbidden raw fields were not copied: raw OIDC JWTs, CI request tokens, provider
credentials, authenticated artifact URLs, and raw logs.

## Report Artifact Digests

Selected clean-report digests from `ci-witness.json`:

| artifact | sha256 |
| --- | --- |
| `clean-feature-flag/run.json` | `2ae1db9059b636834ff2433007342e6a3f3114c9e8488533531711eb48d99e02` |
| `clean-entitlement-matrix/run.json` | `d1d0d39aad8674b76261134fb767a9aee3fc97e6451b1500aa60157a4aa20ee6` |
| `clean-audit-scope/run.json` | `e87338d2c7f1d4fe8dd760356e2e1d62e24db42621e043cf8f1b1326563b7031` |
| `summary.json` | `c25dabe501be03467caba0a31f30e68c7b35d852ec85095bbac06262ce254738` |
| `gate-result.json` | `397fe813deec6d4f141996d0a20ba934b9ae65f47b9e3b1eef85168e4c2bac3c` |
| `verify.txt` | `0cb8d646bcc8e7f24ed322fedaea9585f382b889fefddeaf572b93ca322620c9` |
| `explain.txt` | `d4814e71d5e87fcf034353f76a387a49492f4bfeaf398c5805c7ad8c310b8534` |

Selected no-OIDC digests:

| artifact | sha256 |
| --- | --- |
| `no-oidc-feature-flag/run.json` | `074c29e8be3f566806b5e984ce9a20265be6008f2e837152f65cb0f23380bd11` |
| `no-oidc-env-check.txt` | `c5161741271f8eee4fb85cfe461ac58c1b93549827ec2af4db9ef85e20147ae2` |

## Safety Scan

| field | value |
| --- | --- |
| Pattern file | `docs/research/block-24-redaction-denylist.patterns` |
| Pattern sha256 | `c5ba21129cbc0c969a2d02b46a15bed1cf8c3d48d51643b4eeb6f899150cbbb7` |
| Command shape | `rg --pcre2 -n -f docs/research/block-24-redaction-denylist.patterns <artifact-root>` |
| Clean artifact result | `pass` |
| No-OIDC artifact result | `pass` |

## Dishonest Trace Cases

| case | recorded state | reason codes | customer interpretation |
| --- | --- | --- | --- |
| `dishonest-source-run-mismatch` | `cannot_verify` | `source_run_binding_mismatch`, `ci_witness_not_upgraded` | The command may have run, but copied source/run binding is inconsistent. Treat as incomplete evidence, not trusted CI-witnessed trace. |
| `dishonest-stale-digest-index` | `fail` | `artifact_digest_mismatch`, `stale_index` | The artifact index no longer matches referenced bytes. Treat the artifact as tampered or stale until rerun and reindexed. |

