# Profile Selection Guide

This guide maps the three profile taxonomies used in `sdp-trace` and helps you
choose the right one. For the canonical state and exit-code contract, see
`docs/agent-entrypoint.md`.

## Taxonomy Overview

| Taxonomy | What it describes | Examples |
| --- | --- | --- |
| **Trust profile ID** | What level of trust the evidence can support | `repo_baseline_structural`, `source_bound_local_release`, `external_production_trust` |
| **Assessment profile** | Which kind of evidence quality check to run | `adapter-capture`, `managed-harness`, `forensic-retention`, `ci-artifact-observation`, `authority-envelope` |
| **Witness kind** | Which CI or customer system provided identity evidence | `github-actions`, `gitlab-ci`, `buildkite`, `customer-pki` |
| **Authority scope** | The reporting boundary for a committed package | `demo_pilot_only`, `local_dirty_structural_only` |

Trust profile IDs and authority scopes are **not** commands. Assessment profiles
are selected with `sdp-trace assess --profile <profile>`. Witness kinds are
selected with `sdp-trace witness --kind <kind>`.

## Trust Profile IDs

Choose the trust profile from the claim you need to make, not from your role.

| Trust profile ID | Use when | What it proves |
| --- | --- | --- |
| `repo_baseline_structural` | You need structural command, fixture, and local trace integrity. | Local shape and debug inspection. |
| `source_bound_local_release` | You need local manifest, source commit, artifact digest, and DSSE/source-bound checks. | The built artifact matches the source commit and manifest. |
| `external_production_trust` | You need external identity, protected source, transparency or customer audit evidence, approval, and production release verification. | The full external trust chain is closed. |

**Rule**: Do not use a lower trust profile to claim a higher one. Dirty-checkout
output is valid only under `local_dirty_structural_only` (an authority scope,
not a profile ID) and cannot close `source_bound_local_release` or
`external_production_trust`.

## Assessment Profiles

Choose the assessment profile from the evidence question you need answered.

| Question | Assessment profile | Typical command |
| --- | --- | --- |
| Did the adapter capture enough evidence? Is there overclaim risk? | `adapter-capture` | `sdp-trace assess --profile adapter-capture --out assessment.json --run <run-dir>` |
| Does the managed harness evidence satisfy policy, registry, and witness inputs? | `managed-harness` | `sdp-trace assess --profile managed-harness --out assessment.json --contract contract.json --run <run-dir> --adapter-registry registry.json --managed-policy policy.json --managed-witness witness.json` |
| Can retained evidence support forensic reconstruction? | `forensic-retention` | `sdp-trace assess --profile forensic-retention --out assessment.json --run <run-dir> --redaction-policy redaction.json` |
| Are selected artifact families CI-uploaded facts or lower-authority facts? | `ci-artifact-observation` | `sdp-trace assess --profile ci-artifact-observation --out observation.json --artifact-manifest artifact-manifest.json` |
| Do observed actions stay inside a caller-selected authority envelope? | `authority-envelope` | `sdp-trace assess --profile authority-envelope --authority-package authority-package.json --out authority-evaluation.json` |

Assessment profiles produce **facts only**. Block/allow, readiness, and policy
decisions remain with the downstream consumer.

## Witness Kinds

Choose the witness kind from the CI or customer system that produced the run.

| System | Witness kind | Typical command |
| --- | --- | --- |
| GitHub Actions with OIDC | `github-actions` | `sdp-trace witness --kind github-actions --out ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs` |
| GitLab CI | `gitlab-ci` | `sdp-trace witness --kind gitlab-ci --out gitlab-witness.json --witness-envelope envelope.json .sdp-trace-runs` |
| Buildkite | `buildkite` | `sdp-trace witness --kind buildkite --out buildkite-witness.json --witness-envelope envelope.json .sdp-trace-runs` |
| Customer PKI or private-equivalent | `customer-pki` | `sdp-trace witness --kind customer-pki --out customer-pki-witness.json --customer-pki-authority-policy policy.json --customer-pki-public-cert cert.pem --customer-pki-payload-digest <sha256> --customer-pki-freshness-evidence freshness.json .sdp-trace-runs` |

A CI witness file is **not** external production trust by itself. It binds
available evidence to CI identity when the required OIDC or envelope evidence
exists.

## Authority Scopes

Authority scopes describe the boundary of a report or package, not a verifier
result.

| Authority scope | Meaning |
| --- | --- |
| `demo_pilot_only` | Sanitized demo-repository evidence. Supports pilot mechanics only. |
| `local_dirty_structural_only` | Dirty-checkout structural output. Local shape/debug inspection only. |

## Decision Flow

1. **What do you need to prove?** → Choose a **trust profile ID**.
2. **What evidence do you have?** → Choose an **assessment profile** to check it.
3. **Where did the run happen?** → Choose a **witness kind** if CI identity is available.
4. **Is the checkout clean?** → If dirty, the authority scope is `local_dirty_structural_only`.

Do not mix scopes: a `repo_baseline_structural` result plus a `github-actions`
witness does not become `external_production_trust` unless every required
external-trust check passes live.
