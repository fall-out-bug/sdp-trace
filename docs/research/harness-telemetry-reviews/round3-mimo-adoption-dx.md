```text
VERDICT: CHANGES_REQUIRED

Reviewer: staff engineer / platform owner, responsible for wrapper CI,
secrets management, gateway ops, and repo template maintenance.

============================================================
SECTION 1: DEVELOPER RESISTANCE POINTS
============================================================

R1. Shell wrapper is a non-starter for adoption on developer laptops.
    Developers will disable or alias-out any wrapper that intercepts
    their shell. The brief says "attach recorder to existing harness"
    but a shell wrapper is not an attachment to a harness—it is a
    modification of the developer's environment. On macOS, this means
    modifying PATH, creating shell functions, or installing a shim
    binary. Every one of these breaks when:
    - Developer uses nvm, pyenv, rbenv, or asdf (which also shim PATH)
    - Developer uses a subshell, tmux pane, or SSH session
    - Developer runs inside a container or devcontainer
    - Developer uses `exec` or replaces the shell

    I have maintained PATH-based shims before. They break quarterly
    and generate support tickets. The brief does not acknowledge this
    operational cost. Expected developer reaction: alias
    sdp-trace-shell-wrapper="bash" within 48 hours of deployment.

R2. Tool wrapper assumes the harness exposes a tool-call API that can
    be wrapped. For pi, tool calls go through the pi binary directly.
    For OpenCode, tools are invoked inside the agent's runtime. I
    cannot wrap a tool call without either:
    (a) modifying the harness source to emit an event, or
    (b) replacing the tool binary with a shim.

    Option (a) requires a PR to each harness. I maintain zero of those
    repos. Option (b) means every tool (bash, file_read, file_write,
    grep, etc.) needs a wrapper binary on PATH. Combined with R1, this
    doubles the PATH pollution. I will not deploy this.

R3. File watcher (fsevents/inotify) on a repo with node_modules,
    .git, build output, and IDE state will produce thousands of events
    per minute during a normal agent session. Filtering to relevant
    events is a non-trivial engineering problem. The brief assumes
    "file watcher captures mutations" as if it were a clean signal.
    It is not. I would need:
    - ignore patterns per repo (not specified)
    - debounce logic (not specified)
    - attribution logic (which command caused which mutation)
    - deduplication for rapid writes (save-on-type in IDE)

    The brief gives me a table cell. I need a 3-page spec.

R4. LLM gateway proxy means I am now a man-in-the-middle for all
    model calls. My security team will ask:
    - What happens when the proxy is down? (model calls fail)
    - Where are the API keys stored? (proxy needs all provider keys)
    - Who has access to the logs? (prompt/response digests)
    - Is this SOC2 compliant? (no answer in brief)
    - Does this add latency? (yes, every call now has a proxy hop)

    The brief says "commit to digest-only with stated algorithm for
    Slice 2." I cannot sell a gateway to my org without the privacy
    model defined in Slice 1. The CISO will block the deployment.

R5. CI signing requires a CI-side secret (signing key). Who provisions
    it? Who rotates it? The brief says "CI identity signs chain head."
    In GitHub Actions, this is a repository or org secret. In GitLab
    CI, it is a CI/CD variable. In Jenkins, it is a credentials plugin
    entry. Each CI system has different secret injection semantics.
    The brief assumes "CI signs" is a single operation. It is three
    different operations per CI platform, each with different trust
    properties.

    GitHub Actions OIDC tokens can prove "this workflow ran in this
    repo at this time." That is the closest to a signing identity I
    can get without managing a key. The brief should adopt OIDC or
    equivalent, not raw key signing.


============================================================
SECTION 2: OPERATIONAL BURDEN
============================================================

O1. Four interception layers (harness adapter + tool wrapper + shell
    wrapper + file watcher) is four things to install, configure,
    monitor, debug, and keep compatible across:
    - 3+ harnesses (pi, OpenCode, Kilo)
    - 3+ OSes (macOS, Linux, Windows/WSL)
    - 3+ CI systems (GitHub Actions, GitLab CI, Jenkins)
    - 4+ shell environments (bash, zsh, fish, PowerShell)

    That is 108+ combinations for the "minimum viable" stack. I
    maintain platform tooling for ~200 developers. I cannot support
    this matrix with fewer than 2 FTEs dedicated to the wrapper
    layer alone. The brief does not acknowledge this staffing cost.

O2. The event volume is unbounded. A 30-minute agent session with
    pi generates:
    - 50-200 tool calls (file reads, searches, edits)
    - 10-50 shell commands
    - 1000+ file system events (watcher)
    - 1-5 model calls (gateway)

    Multiply by hash computation per event. Multiply by signature
    at checkpoints. The brief says nothing about storage format,
    compression, retention, or rotation. My CI artifact storage
    costs will spike. My repo size will grow. The brief needs a
    storage budget.

O3. Debugging a failed verification is currently impossible. The
    verifier emits 8 states with no event-level explanation. When a
    run fails verification, the developer asks "why." My answer is
    currently "the verifier says so." I need:
    - which event broke the chain
    - which expected event is missing
    - which layer disagrees with which layer
    - what the developer should do to fix it

    The brief defines what the verifier outputs. It does not define
    what the verifier explains.

O4. Key management is a blocker. The brief lists 4 signing levels.
    Each requires a key. I need:
    - Local recorder: key generated per machine, stored in what?
      macOS Keychain? ~/.sdp-trace/? Agent process memory?
    - CI: secret managed by CI platform.
    - External witness: API key or TLS client cert for the witness
      service (which does not exist).
    - Agent-signed: the agent generates its own key? Where? How does
      the verifier distinguish agent key from recorder key?

    I cannot deploy any signing level without a key lifecycle spec.
    The brief should mark all signing claims as `not_assessed` until
    this exists.


============================================================
SECTION 3: FIRST-WEEK INTEGRATION (what I would actually build)
============================================================

Day 1-2: Shell wrapper prototype (bash only, macOS only)
  - Shim on PATH that logs: timestamp, cwd, argv digest, exit code
  - Store as JSONL in ~/.sdp-trace/sessions/<session-id>.jsonl
  - Session-id derived from shell PID or env var
  - Result: ~60% of build/test commands captured for bash users

Day 3-4: Git hook post-commit listener
  - On commit: capture commit hash, branch, diff stat, changed files
  - Store alongside session events
  - Result: file mutation evidence tied to commit, not to tool call

Day 5: Manual verification script
  - Read session JSONL
  - Check hash chain integrity
  - Report: event count, hash breaks, missing timestamps
  - Result: basic structural verification, no trust claims

What I cannot build in week 1:
  - Tool wrapper (requires harness API research per harness)
  - File watcher (requires filtering spec I do not have)
  - Gateway proxy (requires privacy model I do not have)
  - CI signing (requires signing spec I do not have)
  - Remote witness (service does not exist)
  - Harness adapter (requires harness PR or plugin API)

Week 1 deliverable: shell wrapper + git hook + hash chain + verification
script. This is the "observe-only" stage. It produces local telemetry
with zero trust claims. The brief says "observe-only: no blocking." I
agree. It also produces no value beyond a structured log. The CTO query
surface cannot be answered from shell-wrapper-only telemetry. The brief
does not acknowledge this value gap.


============================================================
SECTION 4: REQUIRED CONTRACTS (what I need before I build anything)
============================================================

C1. EVENT SCHEMA (JSON Schema, versioned, mandatory fields)
    Minimum fields per event:
    - event_id (UUID, mandatory)
    - timestamp (ISO 8601 with timezone, mandatory)
    - source (enum: agent | recorder | gateway | ci | witness, mandatory)
    - event_type (enum, mandatory)
    - session_id (UUID, mandatory)
    - harness_id (string, mandatory)
    - parent_session_id (UUID, optional, for subagents)
    - request_id (UUID, optional, for gateway correlation)
    - payload (object, schema varies by event_type)
    - prev_hash (hex string, mandatory for chain linking)

    Without this schema, I cannot store, validate, or query events.
    The brief gives me prose. I need a file.

C2. VERIFIER STATE MACHINE
    Three-axis model (from GLM round 2):
    - Verdict: pass | fail | cannot_verify | not_assessed
    - Scope: local_only | ci_witnessed | externally_witnessed
    - Completeness: complete | partial | unknown

    With composition rules for multi-layer conflicts. Without this,
    my verification script emits ambiguous states that developers
    cannot act on.

C3. COMPLETENESS CONTRACT
    Per run type, define:
    - Required events (e.g., "at least one session_start, at least
      one session_end, at least one command event for builds")
    - Optional events (e.g., "gateway events are optional for
      local-only runs")
    - Absent-event rules (e.g., "if no gateway events and run type
      is 'ci_witnessed', completeness is partial")

    Without this, the verifier cannot distinguish "short run" from
    "truncated run." The brief names this gap. It remains a gap.

C4. SIGNING SPEC
    Per signing level:
    - Key format (Ed25519? RSA? P-256?)
    - Key generation (who, when, where stored)
    - Key rotation (how often, how triggered)
    - Key revocation (how detected by verifier)
    - Signature format (JWS? raw detached? COSE?)
    - What is signed (event hash? chain head? checkpoint?)

    Without this, the signing model is vapor. The brief lists 4
    levels with a 2-column table. I need 4 pages.

C5. HARNESS ADAPTER API
    Per supported harness:
    - What events does the harness natively emit?
    - What events can a plugin capture?
    - What events require external interception?
    - What is the plugin installation mechanism?
    - What is the plugin bypass detection mechanism?

    Without this, "harness adapter" is a placeholder. I cannot
    write a pi plugin without knowing pi's plugin API (if it has
    one). I cannot write an OpenCode adapter without knowing
    OpenCode's extension points. The brief should name the first
    harness and the first adapter target. It names neither.


============================================================
SECTION 5: UNSAFE TELEMETRY (what the design makes dangerous)
============================================================

U1. PROMPT HASH INVERSION
    The brief proposes prompt hashes at the LLM gateway. Short
    prompts (< 20 chars) are trivially invertible via rainbow table.
    Even salted hashes of common prompts ("fix this bug", "refactor
    this function") are low-entropy. If the gateway log leaks, an
    adversary reconstructs prompts. My CISO will classify this as a
    data exfiltration risk and block gateway deployment.

    Required: minimum prompt length for hashing, or HMAC with a
    per-deployment key that is rotated and not stored in the log.

U2. FILE PATH EXPOSURE
    Tool wrapper and file watcher capture file paths. In a monorepo,
    file paths reveal project structure, team names, feature branches,
    and internal codenames. If telemetry is stored in a shared location
    or transmitted to a witness, file paths are a data classification
    risk. The brief does not mention path redaction or classification.

U3. COMMAND ARGV EXPOSURE
    Shell wrapper captures argv. Commands like `curl -H "Authorization:
    Bearer <token>" ...` or `psql -p 5432 -U admin ...` contain
    secrets in argv. The brief does not define argv scrubbing. My
    developers will run commands with inline secrets. The wrapper will
    log them. This is a secret leak vector.

U4. GATEWAY LOG RETENTION
    If the LLM gateway stores prompt/response digests, how long are
    they retained? Who can access them? Is there an audit trail for
    access? The brief does not define retention policy. My compliance
    team requires a retention spec before any log storage is approved.

U5. LOCAL TELEMETRY STORE PERMISSIONS
    Session JSONL files on developer machines contain all captured
    events. If a developer's machine is compromised, these files
    contain command history, file paths, and potentially scrubbed
    prompts. The brief does not define:
    - file permissions (should be 0600, owner-only)
    - encryption at rest (optional but recommended)
    - automatic cleanup (retention period)

    I have seen developer laptops exfiltrated. Telemetry files are
    a high-value target.


============================================================
SECTION 6: FAILURE MODES (what breaks in production)
============================================================

F1. SHELL WRAPPER BREAKS SCRIPTS
    Any script that parses `$BASH_VERSION` or `$0` will break if
    the wrapper changes the shell invocation context. Scripts that
    use `exec` will replace the wrapper. Scripts that source other
    scripts will see the wrapper's PATH. This is not theoretical;
    I have seen it with nvm, pyenv, and direnv.

F2. TOOL WRAPPER BREAKS HARNESS
    If I shim a tool binary, the harness may detect the shim (wrong
    binary path, wrong version, unexpected stderr output). pi and
    OpenCode may reject the tool invocation. The agent session fails
    and the developer blames the telemetry system.

F3. FILE WATCHER DRAINS BATTERY / CPU
    fsevents on macOS is efficient for small repos. On a monorepo
    with 100k files, the watcher will consume measurable CPU. During
    a `git checkout` or `npm install`, the event flood will overwhelm
    any event processing pipeline. I need backpressure logic and
    event batching. The brief does not mention either.

F4. GATEWAY PROXY ADDS LATENCY
    Every model call now has an extra network hop. For streaming
    responses (SSE/WebSocket), the proxy must be transparent. If
    the proxy buffers or inspects streaming chunks, latency increases
    and developer experience degrades. The brief does not address
    streaming proxy semantics.

F5. CI SIGNING FAILS SILENTLY
    If the CI signing step fails (key not found, signing binary
    crashes, network timeout to witness), the build passes but the
    run is unsigned. The verifier later sees an unsigned run and
    emits `not_assessed`. The developer does not know why their run
    is untrusted. The brief does not define CI signing failure
    semantics (should the build fail? should it warn? should it
    proceed unsigned?).

F6. VERIFIER RUNS TOO LATE
    If the verifier runs after the session closes (post-hoc), it
    cannot detect post-hoc chain fabrication. If it runs during the
    session (inline), it adds latency to every event. The brief does
    not specify when the verifier runs. This is a fundamental
    architecture decision that affects every failure mode above.


============================================================
SECTION 7: DEMO SCOPE CUTS (what I would demo in 2 weeks)
============================================================

Demo scope (achievable in 2 weeks with 1 engineer):

1. Shell wrapper for bash on macOS
   - Logs: timestamp, cwd, argv digest (SHA-256, first 16 chars),
     exit code
   - JSONL format in ~/.sdp-trace/sessions/
   - Session ID from $SDP_TRACE_SESSION_ID env var

2. Git post-commit hook
   - Logs: commit hash, branch, diff stat, file list
   - Correlated to session ID via env var

3. Hash chain
   - Each event linked to previous via SHA-256(prev_event_json)
   - Genesis event has prev_hash = "0"*64
   - Chain stored as JSONL with embedded hash

4. Minimal verifier
   - Reads session JSONL
   - Validates hash chain integrity
   - Reports: event count, chain integrity (pass/fail), event
     source breakdown
   - Emits: { verdict, scope, completeness } per the 3-axis model

5. One worked example
   - Agent runs `npm test` via shell wrapper
   - Git hook captures commit
   - Verifier validates chain
   - Output shown as structured report (not raw JSON)

NOT in demo:
  - Tool wrapper (no harness adapter spec)
  - File watcher (no filtering spec)
  - LLM gateway (no privacy model)
  - CI signing (no signing spec)
  - Remote witness (service does not exist)
  - Harness adapter (no plugin API)
  - Multi-session / subagent correlation (no schema)
  - CTO dashboard (no query surface spec)
  - Key management (no lifecycle spec)
  - Anti-forgery beyond hash integrity (no completeness contract)

This demo shows: local event capture, hash chain integrity, and
basic verification. It does not show: trust, provenance, witness,
or CTO value. That is honest. The brief should acknowledge that
the first demo is a structured logger, not a trust layer.


============================================================
SECTION 8: MINIMUM VIABLE CHANGES
============================================================

M1. Commit to one integration point first.
    Recommendation: shell wrapper + git hook. This is cross-harness,
    cross-platform (with effort), and does not require harness
    cooperation. The brief should name this as the Slice 1 target
    and defer tool wrapper, file watcher, and gateway to Slice 2+.

M2. Write the event schema.
    JSON Schema, version 0.1, with the fields from C1 above. Validate
    it with the existing baseline verifier. This is the contract that
    everything else builds on.

M3. Write the verifier state machine.
    Three-axis model (verdict × scope × completeness). Define
    composition rules for multi-layer conflicts. Define what each
    state means for a developer and for a CTO. This is the output
    contract.

M4. Define argv scrubbing rules.
    Before the shell wrapper logs argv, it must redact:
    - env vars in argv (pattern: all-caps=value)
    - URLs with credentials (pattern: https://user:pass@...)
    - Base64 blobs (pattern: long base64 strings)
    - File paths to known secret locations (~/.ssh/, ~/.aws/, etc.)
    Without this, the wrapper is a secret leak vector.

M5. Acknowledge the demo is a structured logger, not a trust layer.
    The brief's product language says "flight recorder, telemetry
    contract, verifier, witness model, query surface." The first
    demo delivers the first two at local-only scope with no trust
    claims. The brief should be explicit about this scope to avoid
    the overclaim the AGENTS.md trust rules are designed to prevent.

M6. Defer gateway to a follow-up brief.
    Gateway integration requires privacy model, streaming proxy
    semantics, key management, and CISO approval. None of these
    exist. Including gateway in the same brief as shell wrapper
    conflates a 2-week build with a 3-month build. Split the brief.

M7. Defer witness to a follow-up brief.
    Remote witness requires a protocol, a service, a key management
    model, and a deployment story. None of these exist. The brief
    names it as "strongest anti-forgery boundary" without a protocol.
    This is a promise that cannot be kept in the current design
    cycle. Remove it from the demo scope and mark it `not_assessed`
    in all signing claims.


============================================================
SECTION 9: QUESTIONS BEFORE I INVEST TIME
============================================================

Q1. Does pi have a plugin API that emits tool-call events?
    If not, the tool wrapper layer is dead for pi. I need a yes/no
    before I write a pi adapter.

Q2. Does OpenCode have an extension point for telemetry hooks?
    Same question. If neither harness has an API, "harness adapter"
    is a concept, not a deliverable.

Q3. Who is the first target user? A single developer? A team? An
    org? The adoption story changes completely:
    - Single developer: shell wrapper + local verifier, no sharing
    - Team: shared CI signing, team-level witness
    - Org: gateway, policy engine, harness registry

    The brief does not specify. I cannot design for all three.

Q4. What is the retention requirement for session telemetry?
    7 days? 90 days? Forever? This affects storage, cost, and
    compliance. The brief says nothing.

Q5. What happens when the developer runs the agent outside the
    wrapper? The brief says "managed mode: verifier fails or blocks"
    but "observe-only: no blocking." If I deploy observe-only, the
    developer can run without wrapper and the verifier emits
    `not_assessed`. This is the correct behavior but it means the
    CTO cannot trust any run's completeness in observe-only mode.
    Is that acceptable?

Q6. What CI platforms are in scope for Slice 1?
    GitHub Actions? GitLab CI? Jenkins? Each requires different
    signing integration. I need to know which one to build first.

Q7. Is the verifier allowed to fail the CI build?
    If yes, developers will hate it. If no, the verifier is advisory
    only and provides no enforcement. The brief needs to pick a side.
```

</file>

I'll save this review to the reviews directory.

```bash
cat > /Users/fall_out_bug/projects/vibe_coding/sdp-trace/docs/research/harness-telemetry-reviews/round3-pi-staff-engineer-adoption-dx.md << 'REVIEW_EOF'
# Round 3 Adoption/DX Review: Staff Engineer / Platform Owner

**Source text**: `docs/research/harness-telemetry-trust-brief.md`
**Reviewer role**: staff engineer and platform owner responsible for wrappers, CI, secrets, gateway, and repo templates
**Date**: 2026-05-05

---

```text
VERDICT: CHANGES_REQUIRED

Reviewer: staff engineer / platform owner, responsible for wrapper CI,
secrets management, gateway ops, and repo template maintenance.

============================================================
SECTION 1: DEVELOPER RESISTANCE POINTS
============================================================

R1. Shell wrapper is a non-starter for adoption on developer laptops.
    Developers will disable or alias-out any wrapper that intercepts
    their shell. The brief says "attach recorder to existing harness"
    but a shell wrapper is not an attachment to a harness—it is a
    modification of the developer's environment. On macOS, this means
    modifying PATH, creating shell functions, or installing a shim
    binary. Every one of these breaks when:
    - Developer uses nvm, pyenv, rbenv, or asdf (which also shim PATH)
    - Developer uses a subshell, tmux pane, or SSH session
    - Developer runs inside a container or devcontainer
    - Developer uses `exec` or replaces the shell

    I have maintained PATH-based shims before. They break quarterly
    and generate support tickets. The brief does not acknowledge this
    operational cost. Expected developer reaction: alias
    sdp-trace-shell-wrapper="bash" within 48 hours of deployment.

R2. Tool wrapper assumes the harness exposes a tool-call API that can
    be wrapped. For pi, tool calls go through the pi binary directly.
    For OpenCode, tools are invoked inside the agent's runtime. I
    cannot wrap a tool call without either:
    (a) modifying the harness source to emit an event, or
    (b) replacing the tool binary with a shim.

    Option (a) requires a PR to each harness. I maintain zero of those
    repos. Option (b) means every tool (bash, file_read, file_write,
    grep, etc.) needs a wrapper binary on PATH. Combined with R1, this
    doubles the PATH pollution. I will not deploy this.

R3. File watcher (fsevents/inotify) on a repo with node_modules,
    .git, build output, and IDE state will produce thousands of events
    per minute during a normal agent session. Filtering to relevant
    events is a non-trivial engineering problem. The brief assumes
    "file watcher captures mutations" as if it were a clean signal.
    It is not. I would need:
    - ignore patterns per repo (not specified)
    - debounce logic (not specified)
    - attribution logic (which command caused which mutation)
    - deduplication for rapid writes (save-on-type in IDE)

    The brief gives me a table cell. I need a 3-page spec.

R4. LLM gateway proxy means I am now a man-in-the-middle for all
    model calls. My security team will ask:
    - What happens when the proxy is down? (model calls fail)
    - Where are the API keys stored? (proxy needs all provider keys)
    - Who has access to the logs? (prompt/response digests)
    - Is this SOC2 compliant? (no answer in brief)
    - Does this add latency? (yes, every call now has a proxy hop)

    The brief says "commit to digest-only with stated algorithm for
    Slice 2." I cannot sell a gateway to my org without the privacy
    model defined in Slice 1. The CISO will block the deployment.

R5. CI signing requires a CI-side secret (signing key). Who provisions
    it? Who rotates it? The brief says "CI identity signs chain head."
    In GitHub Actions, this is a repository or org secret. In GitLab
    CI, it is a CI/CD variable. In Jenkins, it is a credentials plugin
    entry. Each CI system has different secret injection semantics.
    The brief assumes "CI signs" is a single operation. It is three
    different operations per CI platform, each with different trust
    properties.

    GitHub Actions OIDC tokens can prove "this workflow ran in this
    repo at this time." That is the closest to a signing identity I
    can get without managing a key. The brief should adopt OIDC or
    equivalent, not raw key signing.


============================================================
SECTION 2: OPERATIONAL BURDEN
============================================================

O1. Four interception layers (harness adapter + tool wrapper + shell
    wrapper + file watcher) is four things to install, configure,
    monitor, debug, and keep compatible across:
    - 3+ harnesses (pi, OpenCode, Kilo)
    - 3+ OSes (macOS, Linux, Windows/WSL)
    - 3+ CI systems (GitHub Actions, GitLab CI, Jenkins)
    - 4+ shell environments (bash, zsh, fish, PowerShell)

    That is 108+ combinations for the "minimum viable" stack. I
    maintain platform tooling for ~200 developers. I cannot support
    this matrix with fewer than 2 FTEs dedicated to the wrapper
    layer alone. The brief does not acknowledge this staffing cost.

O2. The event volume is unbounded. A 30-minute agent session with
    pi generates:
    - 50-200 tool calls (file reads, searches, edits)
    - 10-50 shell commands
    - 1000+ file system events (watcher)
    - 1-5 model calls (gateway)

    Multiply by hash computation per event. Multiply by signature
    at checkpoints. The brief says nothing about storage format,
    compression, retention, or rotation. My CI artifact storage
    costs will spike. My repo size will grow. The brief needs a
    storage budget.

O3. Debugging a failed verification is currently impossible. The
    verifier emits 8 states with no event-level explanation. When a
    run fails verification, the developer asks "why." My answer is
    currently "the verifier says so." I need:
    - which event broke the chain
    - which expected event is missing
    - which layer disagrees with which layer
    - what the developer should do to fix it

    The brief defines what the verifier outputs. It does not define
    what the verifier explains.

O4. Key management is a blocker. The brief lists 4 signing levels.
    Each requires a key. I need:
    - Local recorder: key generated per machine, stored in what?
      macOS Keychain? ~/.sdp-trace/? Agent process memory?
    - CI: secret managed by CI platform.
    - External witness: API key or TLS client cert for the witness
      service (which does not exist).
    - Agent-signed: the agent generates its own key? Where? How does
      the verifier distinguish agent key from recorder key?

    I cannot deploy any signing level without a key lifecycle spec.
    The brief should mark all signing claims as `not_assessed` until
    this exists.


============================================================
SECTION 3: FIRST-WEEK INTEGRATION (what I would actually build)
============================================================

Day 1-2: Shell wrapper prototype (bash only, macOS only)
  - Shim on PATH that logs: timestamp, cwd, argv digest, exit code
  - Store as JSONL in ~/.sdp-trace/sessions/<session-id>.jsonl
  - Session-id derived from shell PID or env var
  - Result: ~60% of build/test commands captured for bash users

Day 3-4: Git hook post-commit listener
  - On commit: capture commit hash, branch, diff stat, changed files
  - Store alongside session events
  - Result: file mutation evidence tied to commit, not to tool call

Day 5: Manual verification script
  - Read session JSONL
  - Check hash chain integrity
  - Report: event count, hash breaks, missing timestamps
  - Result: basic structural verification, no trust claims

What I cannot build in week 1:
  - Tool wrapper (requires harness API research per harness)
  - File watcher (requires filtering spec I do not have)
  - Gateway proxy (requires privacy model I do not have)
  - CI signing (requires signing spec I do not have)
  - Remote witness (service does not exist)
  - Harness adapter (requires harness PR or plugin API)

Week 1 deliverable: shell wrapper + git hook + hash chain + verification
script. This is the "observe-only" stage. It produces local telemetry
with zero trust claims. The brief says "observe-only: no blocking." I
agree. It also produces no value beyond a structured log. The CTO query
surface cannot be answered from shell-wrapper-only telemetry. The brief
does not acknowledge this value gap.


============================================================
SECTION 4: REQUIRED CONTRACTS (what I need before I build anything)
============================================================

C1. EVENT SCHEMA (JSON Schema, versioned, mandatory fields)
    Minimum fields per event:
    - event_id (UUID, mandatory)
    - timestamp (ISO 8601 with timezone, mandatory)
    - source (enum: agent | recorder | gateway | ci | witness, mandatory)
    - event_type (enum, mandatory)
    - session_id (UUID, mandatory)
    - harness_id (string, mandatory)
    - parent_session_id (UUID, optional, for subagents)
    - request_id (UUID, optional, for gateway correlation)
    - payload (object, schema varies by event_type)
    - prev_hash (hex string, mandatory for chain linking)

    Without this schema, I cannot store, validate, or query events.
    The brief gives me prose. I need a file.

C2. VERIFIER STATE MACHINE
    Three-axis model (from GLM round 2):
    - Verdict: pass | fail | cannot_verify | not_assessed
    - Scope: local_only | ci_witnessed | externally_witnessed
    - Completeness: complete | partial | unknown

    With composition rules for multi-layer conflicts. Without this,
    my verification script emits ambiguous states that developers
    cannot act on.

C3. COMPLETENESS CONTRACT
    Per run type, define:
    - Required events (e.g., "at least one session_start, at least
      one session_end, at least one command event for builds")
    - Optional events (e.g., "gateway events are optional for
      local-only runs")
    - Absent-event rules (e.g., "if no gateway events and run type
      is 'ci_witnessed', completeness is partial")

    Without this, the verifier cannot distinguish "short run" from
    "truncated run." The brief names this gap. It remains a gap.

C4. SIGNING SPEC
    Per signing level:
    - Key format (Ed25519? RSA? P-256?)
    - Key generation (who, when, where stored)
    - Key rotation (how often, how triggered)
    - Key revocation (how detected by verifier)
    - Signature format (JWS? raw detached? COSE?)
    - What is signed (event hash? chain head? checkpoint?)

    Without this, the signing model is vapor. The brief lists 4
    levels with a 2-column table. I need 4 pages.

C5. HARNESS ADAPTER API
    Per supported harness:
    - What events does the harness natively emit?
    - What events can a plugin capture?
    - What events require external interception?
    - What is the plugin installation mechanism?
    - What is the plugin bypass detection mechanism?

    Without this, "harness adapter" is a placeholder. I cannot
    write a pi plugin without knowing pi's plugin API (if it has
    one). I cannot write an OpenCode adapter without knowing
    OpenCode's extension points. The brief should name the first
    harness and the first adapter target. It names neither.


============================================================
SECTION 5: UNSAFE TELEMETRY (what the design makes dangerous)
============================================================

U1. PROMPT HASH INVERSION
    The brief proposes prompt hashes at the LLM gateway. Short
    prompts (< 20 chars) are trivially invertible via rainbow table.
    Even salted hashes of common prompts ("fix this bug", "refactor
    this function") are low-entropy. If the gateway log leaks, an
    adversary reconstructs prompts. My CISO will classify this as a
    data exfiltration risk and block gateway deployment.

    Required: minimum prompt length for hashing, or HMAC with a
    per-deployment key that is rotated and not stored in the log.

U2. FILE PATH EXPOSURE
    Tool wrapper and file watcher capture file paths. In a monorepo,
    file paths reveal project structure, team names, feature branches,
    and internal codenames. If telemetry is stored in a shared location
    or transmitted to a witness, file paths are a data classification
    risk. The brief does not mention path redaction or classification.

U3. COMMAND ARGV EXPOSURE
    Shell wrapper captures argv. Commands like `curl -H "Authorization:
    Bearer <token>" ...` or `psql -p 5432 -U admin ...` contain
    secrets in argv. The brief does not define argv scrubbing. My
    developers will run commands with inline secrets. The wrapper will
    log them. This is a secret leak vector.

U4. GATEWAY LOG RETENTION
    If the LLM gateway stores prompt/response digests, how long are
    they retained? Who can access them? Is there an audit trail for
    access? The brief does not define retention policy. My compliance
    team requires a retention spec before any log storage is approved.

U5. LOCAL TELEMETRY STORE PERMISSIONS
    Session JSONL files on developer machines contain all captured
    events. If a developer's machine is compromised, these files
    contain command history, file paths, and potentially scrubbed
    prompts. The brief does not define:
    - file permissions (should be 0600, owner-only)
    - encryption at rest (optional but recommended)
    - automatic cleanup (retention period)

    I have seen developer laptops exfiltrated. Telemetry files are
    a high-value target.


============================================================
SECTION 6: FAILURE MODES (what breaks in production)
============================================================

F1. SHELL WRAPPER BREAKS SCRIPTS
    Any script that parses `$BASH_VERSION` or `$0` will break if
    the wrapper changes the shell invocation context. Scripts that
    use `exec` will replace the wrapper. Scripts that source other
    scripts will see the wrapper's PATH. This is not theoretical;
    I have seen it with nvm, pyenv, and direnv.

F2. TOOL WRAPPER BREAKS HARNESS
    If I shim a tool binary, the harness may detect the shim (wrong
    binary path, wrong version, unexpected stderr output). pi and
    OpenCode may reject the tool invocation. The agent session fails
    and the developer blames the telemetry system.

F3. FILE WATCHER DRAINS BATTERY / CPU
    fsevents on macOS is efficient for small repos. On a monorepo
    with 100k files, the watcher will consume measurable CPU. During
    a `git checkout` or `npm install`, the event flood will overwhelm
    any event processing pipeline. I need backpressure logic and
    event batching. The brief does not mention either.

F4. GATEWAY PROXY ADDS LATENCY
    Every model call now has an extra network hop. For streaming
    responses (SSE/WebSocket), the proxy must be transparent. If
    the proxy buffers or inspects streaming chunks, latency increases
    and developer experience degrades. The brief does not address
    streaming proxy semantics.

F5. CI SIGNING FAILS SILENTLY
    If the CI signing step fails (key not found, signing binary
    crashes, network timeout to witness), the build passes but the
    run is unsigned. The verifier later sees an unsigned run and
    emits `not_assessed`. The developer does not know why their run
    is untrusted. The brief does not define CI signing failure
    semantics (should the build fail? should it warn? should it
    proceed unsigned?).

F6. VERIFIER RUNS TOO LATE
    If the verifier runs after the session closes (post-hoc), it
    cannot detect post-hoc chain fabrication. If it runs during the
    session (inline), it adds latency to every event. The brief does
    not specify when the verifier runs. This is a fundamental
    architecture decision that affects every failure mode above.


============================================================
SECTION 7: DEMO SCOPE CUTS (what I would demo in 2 weeks)
============================================================

Demo scope (achievable in 2 weeks with 1 engineer):

1. Shell wrapper for bash on macOS
   - Logs: timestamp, cwd, argv digest (SHA-256, first 16 chars),
     exit code
   - JSONL format in ~/.sdp-trace/sessions/
   - Session ID from $SDP_TRACE_SESSION_ID env var

2. Git post-commit hook
   - Logs: commit hash, branch, diff stat, file list
   - Correlated to session ID via env var

3. Hash chain
   - Each event linked to previous via SHA-256(prev_event_json)
   - Genesis event has prev_hash = "0"*64
   - Chain stored as JSONL with embedded hash

4. Minimal verifier
   - Reads session JSONL
   - Validates hash chain integrity
   - Reports: event count, chain integrity (pass/fail), event
     source breakdown
   - Emits: { verdict, scope, completeness } per the 3-axis model

5. One worked example
   - Agent runs `npm test` via shell wrapper
   - Git hook captures commit
   - Verifier validates chain
   - Output shown as structured report (not raw JSON)

NOT in demo:
  - Tool wrapper (no harness adapter spec)
  - File watcher (no filtering spec)
  - LLM gateway (no privacy model)
  - CI signing (no signing spec)
  - Remote witness (service does not exist)
  - Harness adapter (no plugin API)
  - Multi-session / subagent correlation (no schema)
  - CTO dashboard (no query surface spec)
  - Key management (no lifecycle spec)
  - Anti-forgery beyond hash integrity (no completeness contract)

This demo shows: local event capture, hash chain integrity, and
basic verification. It does not show: trust, provenance, witness,
or CTO value. That is honest. The brief should acknowledge that
the first demo is a structured logger, not a trust layer.


============================================================
SECTION 8: MINIMUM VIABLE CHANGES
============================================================

M1. Commit to one integration point first.
    Recommendation: shell wrapper + git hook. This is cross-harness,
    cross-platform (with effort), and does not require harness
    cooperation. The brief should name this as the Slice 1 target
    and defer tool wrapper, file watcher, and gateway to Slice 2+.

M2. Write the event schema.
    JSON Schema, version 0.1, with the fields from C1 above. Validate
    it with the existing baseline verifier. This is the contract that
    everything else builds on.

M3. Write the verifier state machine.
    Three-axis model (verdict × scope × completeness). Define
    composition rules for multi-layer conflicts. Define what each
    state means for a developer and for a CTO. This is the output
    contract.

M4. Define argv scrubbing rules.
    Before the shell wrapper logs argv, it must redact:
    - env vars in argv (pattern: all-caps=value)
    - URLs with credentials (pattern: https://user:pass@...)
    - Base64 blobs (pattern: long base64 strings)
    - File paths to known secret locations (~/.ssh/, ~/.aws/, etc.)
    Without this, the wrapper is a secret leak vector.

M5. Acknowledge the demo is a structured logger, not a trust layer.
    The brief's product language says "flight recorder, telemetry
    contract, verifier, witness model, query surface." The first
    demo delivers the first two at local-only scope with no trust
    claims. The brief should be explicit about this scope to avoid
    the overclaim the AGENTS.md trust rules are designed to prevent.

M6. Defer gateway to a follow-up brief.
    Gateway integration requires privacy model, streaming proxy
    semantics, key management, and CISO approval. None of these
    exist. Including gateway in the same brief as shell wrapper
    conflates a 2-week build with a 3-month build. Split the brief.

M7. Defer witness to a follow-up brief.
    Remote witness requires a protocol, a service, a key management
    model, and a deployment story. None of these exist. The brief
    names it as "strongest anti-forgery boundary" without a protocol.
    This is a promise that cannot be kept in the current design
    cycle. Remove it from the demo scope and mark it `not_assessed`
    in all signing claims.


============================================================
SECTION 9: QUESTIONS BEFORE I INVEST TIME
============================================================

Q1. Does pi have a plugin API that emits tool-call events?
    If not, the tool wrapper layer is dead for pi. I need a yes/no
    before I write a pi adapter.

Q2. Does OpenCode have an extension point for telemetry hooks?
    Same question. If neither harness has an API, "harness adapter"
    is a concept, not a deliverable.

Q3. Who is the first target user? A single developer? A team? An
    org? The adoption story changes completely:
    - Single developer: shell wrapper + local verifier, no sharing
    - Team: shared CI signing, team-level witness
    - Org: gateway, policy engine, harness registry

    The brief does not specify. I cannot design for all three.

Q4. What is the retention requirement for session telemetry?
    7 days? 90 days? Forever? This affects storage, cost, and
    compliance. The brief says nothing.

Q5. What happens when the developer runs the agent outside the
    wrapper? The brief says "managed mode: verifier fails or blocks"
    but "observe-only: no blocking." If I deploy observe-only, the
    developer can run without wrapper and the verifier emits
    `not_assessed`. This is the correct behavior but it means the
    CTO cannot trust any run's completeness in observe-only mode.
    Is that acceptable?

Q6. What CI platforms are in scope for Slice 1?
    GitHub Actions? GitLab CI? Jenkins? Each requires different
    signing integration. I need to know which one to build first.

Q7. Is the verifier allowed to fail the CI build?
    If yes, developers will hate it. If no, the verifier is advisory
    only and provides no enforcement. The brief needs to pick a side.
```
REVIEW_EOF

echo "Saved to docs/research/harness-telemetry-reviews/round3-pi-staff-engineer-adoption-dx.md"
```
