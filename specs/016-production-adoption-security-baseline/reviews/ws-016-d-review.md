# WS-016-D Review: Repository Security Policy

Date: 2026-05-20
Files: `.github/SECURITY.md`, `docs/security-baseline.md`, `docs/README.md`

## Quality

| Finding | Severity | Evidence | Disposition |
|---------|----------|----------|-------------|
| Q1: SECURITY.md correctly states only `main` is supported (pre-1.0). | positive | Matches project state. | accepted. |
| Q2: In-scope list covers all product areas: cmd, internal, tools, schema, docs. | positive | Complete. | accepted. |
| Q3: Out-of-scope explicitly names local agent run directories. | positive | Matches FR-016-003. | accepted. |
| Q4: Reporting email uses project-specific address. | info | `security+report@fall-out-bug.dev` | accepted. |

## UX

| Finding | Severity | Evidence | Disposition |
|---------|----------|----------|-------------|
| U1: Response timeline is clear (ack 5 days, assessment 10 days). | positive | Reporter knows what to expect. | accepted. |
| U2: Trust boundary paragraph clarifies that product does not approve changes. | positive | Prevents overclaim in security context too. | accepted. |

## DX

| Finding | Severity | Evidence | Disposition |
|---------|----------|----------|-------------|
| D1: Verification commands block is copy-pasteable. | positive | Includes `gitleaks detect --source . --config .gitleaks.toml`. | accepted. |
| D2: docs/README.md link added in natural alphabetical position. | positive | Between Security Baseline and CI Check Policy. | accepted. |

## Security

| Finding | Severity | Evidence | Disposition |
|---------|----------|----------|-------------|
| S1: Document does NOT claim security guarantee. "best-effort advisories only." | positive | Honest scope. | accepted. |
| S2: Out-of-scope correctly excludes local ignored clutter. | positive | Matches security baseline policy. | accepted. |
| S3: Trust boundary section redirects security reports toward evidence/provenance issues rather than generic code bugs. | positive | Focused scope. | accepted. |

## Synthesis

- All findings positive.
- No fixes required.
- No blockers.
