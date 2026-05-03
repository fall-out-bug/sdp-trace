# Block 04 Socratic Dialogue: Release Finalization and External Trust

Date: 2026-05-01
Block: `04-release-finalization-external-trust`
Beads mirror: `sdp-trace-cdn.21`

## Consensus Candidate

Block 04 should include local release finalization plus external trust design, but it must not claim production trust without real Sigstore/Rekor or customer PKI evidence.

## Consensus Result

Consensus is accepted for implementation. The selected scope is option 2: source-bound local finalization plus external trust schema/fixtures. Implementation proceeds without any `trusted_contract_release: true` claim until all required proof states are assessed and successful.

## Q1: Is local DSSE enough to call the release trusted?

**Critic**: Local DSSE proves a private key signed a manifest, but it does not prove the release came from the protected source, approved workflow, or external release authority.

**Answer**: No. Local DSSE can support `locally_attested` only after source-content verification succeeds. It cannot support `production_release_verified` or `trusted_contract_release` by itself.

**Resolution**: Use "source-bound local release" for local proof and reserve "externally trusted production release" for accepted external trust anchors.

## Q2: Can Block 04 close the current `source_commit_artifacts_verified: not_assessed` state?

**Critic**: The current source commit does not contain final artifacts, so any attempt to mark source-content proof as assessed would repeat the Block 03 bug.

**Answer**: Yes, but only after a final committed source reference exists and manifest artifacts are regenerated against that source reference.

**Resolution**: The finalization command must refuse dirty-tree source-content claims and fail closed without emitting a source proof artifact. Dirty-tree self-attestation results may still record `source_commit_artifacts_verified: not_assessed`, but finalization itself does not create that artifact.

## Q3: Should Block 04 implement real Sigstore/Rekor verification immediately?

**Critic**: A design-only Sigstore section can become decorative if it does not produce executable proof. But implementing public Sigstore now may force GitHub and network assumptions into a portable repo.

**Answer**: Not immediately. The repo can define the evidence contract, negative fixtures, and `not_assessed` behavior without making Sigstore a hard runtime dependency.

**Resolution**: Block 04 implementation should add external trust evidence fields and fixtures first. Real public verification remains optional until the repo has accepted external evidence.

## Q4: Should customer PKI be treated as weaker than Sigstore/Rekor?

**Critic**: If private customers cannot use public Rekor, the product cannot make public transparency mandatory for all production releases.

**Answer**: Customer PKI can be an accepted external trust anchor if the policy names the identity, source reference, audit/timestamp evidence, and approval evidence.

**Resolution**: Do not call customer PKI "local." It is an external trust profile with a different evidence shape.

## Q5: Does this unblock pilot claims?

**Critic**: A source-bound local release improves repository integrity, but it does not prove harness/model compatibility or customer pilot readiness.

**Answer**: No. It unblocks a cleaner foundation for pilot artifacts; it does not replace run-cards, compatibility matrices, or customer evidence packages.

**Resolution**: Keep `sdp-trace-cdn.6`, `.7`, `.9`, and `.10` open. Compatibility claims still require committed pilot evidence or `not_assessed`.

## Q6: Should `trusted_contract_release` be true for customer PKI if Rekor is unavailable?

**Critic**: If the schema requires public Rekor for true, private customers are blocked. If it does not require any transparency substitute, the field becomes too weak.

**Answer**: `trusted_contract_release` can be true for customer PKI only when the selected trusted identity policy explicitly accepts the private profile and the result records assessed audit/timestamp evidence.

**Resolution**: The proof state must include `external_trust_profile`, `transparency_evidence_ref` or accepted substitute, and `approval_status`.

## Q7: What is the main implementation risk?

**Critic**: A release-finalization script can create circular digest churn if it includes artifacts that embed the manifest digest or signer fingerprint.

**Answer**: The manifest artifact list must continue to exclude generated release verification artifacts that embed the manifest digest, DSSE envelope, or local public key fingerprint.

**Resolution**: Document and test the non-circular artifact boundary. The manifest covers source contracts and examples; generated release proof artifacts verify the manifest but are not themselves manifest subjects.

## Q8: What consensus is required before implementation?

**Critic**: "Local plus external trust" is broad. Without a sequence, implementation will sprawl.

**Answer**: The implementation sequence must be:

1. Source-bound local finalization.
2. External trust evidence schema/fixtures.
3. Real external verifier only if accepted evidence is available.

**Resolution**: Proceed with implementation only if the user accepts this staged delivery and the claim boundary.
