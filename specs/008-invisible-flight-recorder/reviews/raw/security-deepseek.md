### Security/Non-Interference Review Findings

#### **1. Filesystem Risks**  
   - **Severity**: Medium  
   - **Description**: The recorder writes append-only run directories and manifests. While this is secure by design, there is a risk that incorrect permissions or misconfiguration could allow unauthorized modification of evidence files.  
   - **Mitigation**: Ensure that recorder-owned evidence directories are created with strict permissions (`chmod 700`) and that manifests are signed or use integrity checks.

#### **2. CI Secret Exposure**  
   - **Severity**: High  
   - **Description**: The `packet build-pr` command relies on GitHub Actions context, which may involve accessing GitHub API tokens or secrets. If misconfigured, sensitive credentials could be exposed in logs or artifacts.  
   - **Mitigation**: Ensure GitHub secrets are masked in logs and not exposed in packet artifacts. Use HTTPS for API calls and validate inputs to prevent accidental leakage.

#### **3. Prompt Injection Risks**  
   - **Severity**: Medium  
   - **Description**: The recorder detects prompt contamination by analyzing developer prompts. Malicious or malformed prompts could exploit this mechanism to bypass contamination checks or inject false evidence.  
   - **Mitigation**: Validate and sanitize prompt inputs strictly. Use regex or whitelisting to ensure only expected content is processed.

#### **4. Harness Interference Risks**  
   - **Severity**: Low  
   - **Description**: The recorder must run passively without interfering with the developer harness. However, there is a risk of race conditions or resource contention if the recorder is not isolated properly.  
   - **Mitigation**: Ensure the recorder operates in a separate process or thread with minimal resource usage. Use locking mechanisms for shared resources.

#### **5. Unsafe Command Risks**  
   - **Severity**: Medium  
   - **Description**: The CLI commands (`sdp-trace recorder run`, `sdp-trace packet build-pr`) could be misused or abused if inputs are not validated properly. For example, passing untrusted arguments to `--developer-command` could lead to command injection.  
   - **Mitigation**: Validate and sanitize all CLI inputs. Use a whitelist of allowed commands or arguments where possible.

#### **6. Evidence Tampering Risks**  
   - **Severity**: High  
   - **Description**: While manually authored `.evidence` or `.sdp-trace` files are discouraged, they could still be used to tamper with evidence if not properly marked as operator/integration artifacts.  
   - **Mitigation**: Enforce strict validation of evidence files and ensure they are explicitly marked with authority metadata. Use cryptographic checksums or signatures to verify integrity.

#### **7. GitHub API Abuse Risks**  
   - **Severity**: Medium  
   - **Description**: Accessing GitHub API endpoints for PR metadata and artifact discovery could lead to rate limiting or API abuse if not handled properly.  
   - **Mitigation**: Implement rate limiting and retry logic for API calls. Use caching where possible to reduce redundant requests.

#### **8. CLI Help Misuse Risks**  
   - **Severity**: Low  
   - **Description**: The CLI help text (`docs/agent-entrypoint.md`, `docs/change-evidence-packet.md`) could inadvertently expose sensitive details or mislead users about command authority.  
   - **Mitigation**: Regularly review and update CLI help text to ensure clarity and avoid exposing sensitive information.

### Summary  
The highest risks involve CI secret exposure and evidence tampering, which could compromise security and trust. Prompt injection and unsafe command risks also require careful mitigation. Filesystem and harness interference risks are lower but should still be addressed to ensure robustness. Regular reviews and strict validation are key to maintaining security and non-interference.
