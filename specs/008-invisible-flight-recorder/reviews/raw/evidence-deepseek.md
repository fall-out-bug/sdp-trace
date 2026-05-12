### Evidence/Trust Review Findings

#### Findings:

1. **Weak Proof of Non-Interference (High Severity)**:
   - The spec claims "non-interfering observer contract" without defining the scope or boundaries of this contract. Without explicit details on how `sdp-trace` ensures non-interference (e.g., no side effects on the developer prompt or runtime behavior), this claim is vague and risks overpromising.
   - **Recommendation**: Define clear boundaries for the observer contract and explicitly state what guarantees are provided.

2. **Manual Evidence Contamination Risk (High Severity)**:
   - While the spec emphasizes that manually edited evidence should not satisfy `PC-AGENT-ROUTE`, it does not detail how `sdp-trace` prevents such contamination or enforces append-only evidence creation.
   - **Recommendation**: Clarify the mechanisms (e.g., file permissions, immutable metadata) that ensure recorder-owned evidence is append-only and cannot be tampered with.

3. **Authority Confusion in Evidence Lifecycle (Medium Severity)**:
   - Multiple actors (recorder, CI, operator, integrators) are involved in evidence lifecycle, but the spec does not explicitly define how each actor's authority is verified or how conflicts between actor-generated evidence are resolved.
   - **Recommendation**: Define authority verification mechanisms for each actor and specify how conflicts or contradictory evidence are handled.

4. **Lack of Detail on Prompt Boundary Classification (Medium Severity)**:
   - The spec introduces "prompt boundary classification" to detect contamination but does not specify the criteria or methodology for classification. This could lead to inconsistent or incorrect classification of prompt contamination.
   - **Recommendation**: Detail the classification criteria and provide examples of both clean and contaminated prompts.

5. **Ambiguity in GitHub Source Discovery (Low Severity)**:
   - The spec requires GitHub source discovery but does not clarify whether the builder should fall back to manual input (`github-input.json`) if GitHub context is incomplete or unavailable.
   - **Recommendation**: Clarify the fallback behavior and ensure it aligns with the principle of non-interference.

#### Summary:
The spec demonstrates a strong focus on separating developer responsibility from recorder duties, but it leaves critical gaps in defining mechanisms to ensure non-interference, prevent evidence contamination, and verify actor authority. These gaps introduce risks of overclaiming and weak proof, which could undermine the trustworthiness of the recorder. Detailed definitions and explicit safeguards are needed to address these issues.
