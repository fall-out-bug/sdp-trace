# Slice 66 Plan Review

Status: pass

Date: 2026-06-04T16:20:20+03:00

## Scope

Slice 66 is bounded to numbered `cmd/sdp-trace` GitHub Actions API URL
validation and trust-target policy shards:

- `packet_076_validategithubapiurl.go`
- `packet_077_validateparsedgithubapiurl.go`
- `packet_078_validategithubapiurltrusttarget.go`
- `packet_079_requirehttpsgithubapi.go`
- `packet_080_parsegithubapiurl.go`
- `packet_081_localhttpgithubapi.go`
- `packet_082_githubapihostallowed.go`
- `packet_083_publicgithubserverhost.go`
- `packet_084_githubserverhost.go`
- `packet_085_loopbackhost.go`

Planned target: `cmd/sdp-trace/packet_build_pr_actions_url_policy.go`.

Excluded from this slice:

- context/source selection (`packet_build_pr_actions_context.go` and
  `packet_build_pr_actions_source.go`)
- HTTP request, fetch, decode, retained-artifact, and resolver helpers
  (`packet_086` through `packet_092`)
- packet fixture type/loading (`packet_093` onward)
- shared optional JSON IO (`packet_095`)

## Behavior Preservation Claims

- syntax diagnostics stay unchanged for malformed API URLs
- embedded credentials are still rejected before trust-target validation
- credential rejection still wins over HTTPS errors for mixed-invalid URLs
- non-local API targets still require HTTPS after credential rejection passes
- HTTP remains allowed only for loopback/local test targets
- public GitHub server host still maps only to `api.github.com`
- Enterprise server hosts still bind exactly to the configured server host
- malformed configured server URLs still fall back to an empty server host
- package boundary, dependency direction, and MI baselines stay unchanged

## Planned Focused Evidence

- `TestGitHubAPIURLPolicy`
- `TestGitHubAPIURLPolicyHelpersKeepTrustBoundaries`
- `TestParseGitHubAPIURLKeepsSyntaxDiagnostics`
- `TestGitHubServerHostKeepsFallbackBehavior`

Focused commands:

```text
go test ./cmd/sdp-trace -list 'Test(GitHubAPIURLPolicy|GitHubAPIURLPolicyHelpersKeepTrustBoundaries|ParseGitHubAPIURLKeepsSyntaxDiagnostics|GitHubServerHostKeepsFallbackBehavior)$' | awk '/^Test/ {print}' > /tmp/slice-66-tests.txt && diff -u <(printf 'TestGitHubAPIURLPolicy\nTestGitHubAPIURLPolicyHelpersKeepTrustBoundaries\nTestParseGitHubAPIURLKeepsSyntaxDiagnostics\nTestGitHubServerHostKeepsFallbackBehavior\n') /tmp/slice-66-tests.txt
go test ./cmd/sdp-trace -run 'Test(GitHubAPIURLPolicy|GitHubAPIURLPolicyHelpersKeepTrustBoundaries|ParseGitHubAPIURLKeepsSyntaxDiagnostics|GitHubServerHostKeepsFallbackBehavior)$'
```

The exact-list command is post-implementation focused evidence; the three new
focused tests are expected to be added during Slice 66 implementation before the
command is run.

## Review Lanes

- scope/correctness: LGTM after one credential/HTTPS ordering correction
- trust/evidence: LGTM after clarifying post-implementation test evidence
- maintainability/DX: LGTM after one credential/HTTPS ordering correction

## Reviewer Metadata

- harness: multi_agent_v1
- model/provider: not_assessed; the harness does not expose exact
  provider-qualified model IDs in this session
- prompt class: Slice 66 SpecKit plan/task review
- timeout: 600000ms per wait
- retries: 1 for all lanes after clarifying credential-before-HTTPS behavior
  and post-implementation exact test list evidence
- fallback: none
