# Block 22 Focused Product-Boundary Re-Review

## Issues Checked Against Fixes

| ID | Finding | Fix claimed | Re-review verdict |
| --- | --- | --- | --- |
| S22-PB-01 | Customer PKI freshness without live PKI could accept self-claimed freshness | Added signed freshness evidence contract: payload digest, run id, policy digest, signer identity, issued time, valid-until, nonce/sequence | **FIXED** — "self-claimed timestamp in unsigned witness JSON is not authority" explicitly ruled out |
| S22-PB-02 | Air-gapped profile left reviewers unable to answer what command or validation path can be run | `air-gapped-v1` is documentation+fixture id only; reviewer validates via repository fixture validation command; network calls forbidden | **FIXED** — runnable validation path is now fixture validation; explicit network-call prohibition added |
| S22-PB-03 | Independence thresholds were abstract; profile implementers could invent trust ceilings | Closed `independence_state` enum added; Trust-Scope Determination matrix added; GitLab section names topology caps | **FIXED** — closed enum and matrix bind implementer choices |
| S22-PB-04 | Buildkite trust ceiling did not state whether `external_witnessed` was achievable | `buildkite-v1` caps at `ci_witnessed`; external requires a later reviewed profile | **FIXED** — explicit ceiling stated |
| S22-PB-05 | CLI error handling for unknown kind, missing inputs, exit states undefined | CLI Boundary defines closed kinds, single-kind, exit 2/1/3/0 | **FIXED** — exit codes and usage errors defined |
| S22-PB-06 | Customer PKI public input flags and safe path handling unnamed | Named `--customer-pki-*` flags; rejects implicit scanning, private keys, provider tokens, customer directories | **FIXED** — flags named and boundary stated |

## Remaining Findings

**None.** All six prior product-boundary critical and major findings have been addressed in the spec text with specific sections, enumerated states, and explicit non-goals that bind implementers.

## Recommendation

**APPROVE.** The reviewed direction is ready for implementation. Implementation must not start until explicit user approval of this reviewed direction is confirmed.
