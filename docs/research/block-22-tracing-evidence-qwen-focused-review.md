## Tracing/Evidence Re-Review: Block 22

Checked all prior critical and major findings from the tracing/evidence plane (S22-TE-01 through S22-TE-06) plus the additional dimensions called out in scope.

| Prior ID | Issue | Status |
| --- | --- | --- |
| S22-TE-01 (critical) | Closed reason-code registry missing | **Fixed** — `Closed Reason Codes` table present with default verifier states for every identity, signer, freshness, artifact, source, run, policy, topology, safety, PKI-specific, and unsupported case. |
| S22-TE-02 (critical) | CLI exit codes undefined | **Fixed** — `CLI Boundary` defines exit `0` (pass), `1` (fail), `2` (usage error), `3` (cannot_verify). |
| S22-TE-03 (critical) | Trust-scope determination rules missing | **Fixed** — `Trust-Scope Determination` matrix maps required facts + independence state to each closed scope (`external_witnessed`, `ci_witnessed`, `local_observed`, `cannot_verify`, `not_assessed`, `fail`). |
| S22-TE-04 (major) | Source/run/policy binding in result but not contract | **Fixed** — All three added to `Witness Profile Contract` table with purpose and non-pass states. |
| S22-TE-05 (major) | Fixture matrix lacked expected outputs | **Fixed** — `Minimum expected fixture rows` table lists profile, verifier state, established scope, and required reason code per fixture. |
| S22-TE-06 (major) | Cross-surface consumption underspecified | **Fixed** — `Cross-Surface Consumption` table maps every normalized field to gate/protected, managed harness, and export consumption semantics. |

Additional checked dimensions:

- **Customer PKI input contract** — `CLI Boundary` names four explicit public-input flags, rejects implicit scanning and private key inputs. ✓
- **Air-gapped identity** — `air-gapped-v1` is documentation + fixture only, not a `--kind`; network calls forbidden; fixture validation uses repository command. ✓
- **Source/run/policy binding symmetry** — Contract table, normalized result, trust-scope matrix, and fixture matrix all reference the same three bindings consistently. ✓
- **Fixture expected outputs** — All 22 minimum rows present. Two rows use "or" for profile or reason code (`ci-run-id-mismatch` and `air-gapped-offline-not-assessed`); this is documentation convenience, not a defect — implementation will materialize concrete profile instances. Minor.
- **Freshness semantics** — `Freshness Evaluation` section defines shared semantics; customer PKI freshness requires explicit signed evidence bound to payload digest, run id, policy digest, signer identity, issued time, and nonce/sequence. ✓
- **Independence / trust-scope separation** — `independence_state` is a closed topology enum that constrains but never equals trust scope; matrix derives scope from facts, not profile name. ✓

No remaining critical or major tracing/evidence findings.

**APPROVE**
