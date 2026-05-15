# Agent Entrypoint

Use these commands and profile IDs to select what `sdp-trace` can prove without
introducing harness assumptions. The active entrypoint is the Go CLI reported by
`sdp-trace --help`; update this document only when the help surface
changes in the same slice.

For the sidecar-first CI pilot evidence package, start with
`examples/pilot-runs/opencode-minimax-kotlin-bazel/README.md`. It shows the
exact observed slice and its current verifier limitations without turning the
run into broad harness or model support.

## Profile Selection

Each assertion is anchored to one of these profile IDs:

- `repo_baseline_structural`
- `source_bound_local_release`
- `external_production_trust`

Do not infer profile from role. Choose the profile directly from the claim:

- `repo_baseline_structural`: structural command, fixture, and local trace integrity.
- `source_bound_local_release`: local manifest, source commit, artifact digest, and DSSE/source-bound checks.
- `external_production_trust`: external identity, protected source, transparency or customer audit evidence, approval, and production release verification.

Dirty-checkout baseline output is only valid under the
`local_dirty_structural_only` authority scope. It is not a profile ID and must
not be used to close `source_bound_local_release` or
`external_production_trust`.

## Result, Trust Scope, And Authority Scope

Keep these vocabularies separate:

- Result state: the verifier outcome for a selected command or profile, such as
  `observed`, `pass`, `fail`, `not_assessed`, or `cannot_verify`.
- Trust scope: the evidence boundary recorded by a run, witness, or assessment,
  such as `local_observed` or `ci_witnessed`.
- Authority scope: the reporting boundary for a committed package, such as
  `demo_pilot_only`.

Known trust scopes used by the current pilot docs:

- `local_observed`: local run/report evidence was captured and checked, but it
  is not CI-witnessed or external production trust.
- `ci_witnessed`: available CI identity and artifact binding evidence supported
  the witness profile for the exact CI topology under assessment.
- `external_witnessed`: external witness evidence was supplied and accepted by
  the selected profile.

Known authority scopes used by current docs:

- `demo_pilot_only`: sanitized demo-repository evidence. It can support pilot
  mechanics and state interpretation only; it does not establish customer
  production trust, owner independence, non-GitHub portability, or release
  binary acquisition.
- `local_dirty_structural_only`: dirty-checkout structural output. It can
  support local shape/debug inspection only; it cannot support source-bound or
  external trust closure.

`not_assessed` can return exit `0` only when the selected command completed and
the unassessed state is explicitly outside the selected profile or run scope.
Pipeline authors must inspect emitted JSON/state fields instead of treating exit
`0` as proof that every possible trust state passed. If a required check for the
selected profile cannot run or lacks required evidence, the command must use
`cannot_verify` and exit `3`.

## Command Contract

### Machine-Readable Command Surface

`sdp-trace command-surface` emits machine-readable JSON describing the current
command surface. This is the authoritative source of truth for
`tools/doccheck` drift checks.

**Stability**: experimental. The `schema_version` prefix is stable but the full
surface shape may change between releases without a semver bump. Agents
consuming this surface should be resilient to new fields.

Current output: `go run ./cmd/sdp-trace command-surface`

## Local Quality Gates

Use these checks for Go-first product-path changes before claiming local
CRAP, cyclomatic-complexity, or cognitive-complexity pass:

- `go test -count=1 ./... -coverprofile=/tmp/sdp-trace-cover.out`
- `go tool cover -func=/tmp/sdp-trace-cover.out > /tmp/sdp-trace-cover-func.txt`
- `go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools`
- `go run ./tools/qualitycheck -fail-only -function-mi-under 70 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`
- `go run ./tools/qualitycheck -fail-only -mi-under 70 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`
- `go run ./tools/qualitycheck -gocyclo cmd internal tools > /tmp/sdp-trace-gocyclo.txt`
- `go run ./tools/crapcheck -cover-func /tmp/sdp-trace-cover-func.txt -gocyclo /tmp/sdp-trace-gocyclo.txt -threshold 5 -strict-less`
- `go run ./tools/qualitycheck -mi-under 70 cmd internal tools` as the advisory absolute executable-file MI check; this is expected to remain `assessed_gap` until the historical baseline is retired

The CRAP, cyclomatic, and cognitive gates cover Go files under `cmd`,
`internal`, and `tools`; test files are excluded. The enforced complexity
threshold is `<= 10`, and the CRAP threshold is strict: a function with CRAP
`5.00` or greater is not a pass. CI enforces
CRAP/cyclomatic/cognitive gates for `cmd`, `internal`, and `tools`, the
function-MI ratchet for `cmd` and `internal`, and the file-MI ratchet for
`cmd`, `internal`, and `tools`. Maintainability Index remains `assessed_gap`
for absolute `> 70`: the baselines prevent regression and block new
below-threshold functions/files, but historical below-threshold functions/files
are not claimed as passing.

Treat Maintainability Index as a local refactor heuristic, not a standalone
quality proof. File-level MI is sensitive to file size and can improve after a
same-package split without behavior changing; keep CRAP, complexity, coverage,
spec drift, and review evidence in the same claim.

Spec drift, work-without-spec, and docs completeness are review gates, not
machine-proof gates today. For behavior or trust-claim changes, compare the
implementation against the active spec, plan, and task, then update this command
contract and affected docs in the same slice. Missing spec or stale docs keep
the gate `cannot_verify` or `not_assessed`; do not convert them into pass by
local prose.

Current command surface:

- `go test -count=1 ./...`
- `sdp-trace --help`
- `sdp-trace command-surface`
- `sdp-trace version`
- `sdp-trace wrap --name <name> [--contract <file>] [--output-dir <dir>] -- <command...>`
- `sdp-trace run --task <task-ref> [--contract <file> | --use-default-contract] -- <command...>`
- `sdp-trace dry-run [--contract <file> | --use-default-contract] -- <command...>`
- `sdp-trace preview [--contract <file> | --use-default-contract] -- <command...>`
- `sdp-trace doctor [--contract <file>]`
- `sdp-trace doctor --profile github-actions-git-hooks-v1 [--out <file>]`
- `sdp-trace install repo-observer --profile github-actions-git-hooks-v1 [--repository-id <safe-id>] [--write] [--force] [--out <file>]`
- `sdp-trace interaction relay --task-id <safe-id> --event-type <type> --out <file> -- <forward-command...>`
- `sdp-trace interaction import-transcript --source preclassified-transcript-import --task-id <safe-id> --events-jsonl <file> --out <file>`
- `sdp-trace interaction summarize --trace <file> [--out <file>]`
- `sdp-trace observe setup --profile <session-profile.json> --out <run-dir> [--command <harness-command-preview>]`
- `sdp-trace observe collect --profile <session-profile.json> --run <run-dir>`
- `sdp-trace observe session --profile <session-profile.json> --out <run-dir> -- <harness-command...>`
- `sdp-trace harness observe --profile <harness-profile.json> --source <harness-events.jsonl> --out <run-dir>`
- `sdp-trace harness validate --profile <harness-profile.json> --run <run-dir> --out <validation.json>`
- `sdp-trace harness summarize --validation <validation.json>`
- `sdp-trace envelope summarize --envelope <file> [--out <file>]`
- `sdp-trace verify <run-dir>`
- `sdp-trace explain <run-dir>`
- `sdp-trace query --query <missing-evidence|capture-depth> <run-dir>`
- `sdp-trace query-pack --pack forensics-basic-v1 --run <run-dir> --out <file>`
- `sdp-trace query-pack explain --result <file>`
- `sdp-trace export cross-repo-posture --profile cross-repo-evidence-posture-v1 --selection <file> --out <file>`
- `sdp-trace export cross-repo-posture explain --result <file>`
- `sdp-trace export telemetry --profile prometheus-text-v1 --cross-repo-posture <file> --out <file|->`
- `sdp-trace assess --profile adapter-capture --out <file> --run <run-dir>`
- `sdp-trace assess --profile managed-harness --out <file> --contract <file> --run <run-dir> --adapter-registry <file> --managed-policy <file> --managed-witness <file>`
- `sdp-trace assess --profile forensic-retention --out <file> --run <run-dir> --redaction-policy <file>`
- `sdp-trace assess --profile ci-artifact-observation --out <file> --artifact-manifest <file>`
- `sdp-trace assess --profile authority-envelope --out <file> --authority-package <file>`
- `sdp-trace assess preview --profile <adapter-capture|managed-harness|forensic-retention|ci-artifact-observation|authority-envelope> [profile inputs]`
- `sdp-trace assess explain --assessment-result <file>`
- `sdp-trace report --out <dir> <runs-root-or-run-dir>`
- `sdp-trace gate --out <file> <runs-root-or-run-dir>`
- `sdp-trace witness --kind <github-actions|gitlab-ci|buildkite|customer-pki> --out <file> [--report-dir <dir>] [--witness-envelope <file>] [--customer-pki-authority-policy <file>] [--customer-pki-public-cert <file> | --customer-pki-public-key <file>] [--customer-pki-payload-digest <sha256>] [--customer-pki-freshness-evidence <file>] <runs-root-or-run-dir>`
- `sdp-trace release-proof --manifest <file> --out <file>`
- `sdp-trace pr-review packet --out <dir> --repo-id <safe-id> --change-ref <pr|mr|change-id> --base <sha> --head <sha> --diff <file> [--ci-state <state>] [--created-by <actor>]`
- `sdp-trace pr-review run --packet <dir> --profile <file> --out <dir> [--preview] [--work-dir <dir>] [--allow-external-runner <runner>]...`
- `sdp-trace pr-review synthesize --packet <dir> --runs <dir> --out <file>`
- `sdp-trace pr-review validate --packet <dir> --profile <file> --runs <dir> --ledger <file> --out <file>`
- `sdp-trace pr-review summarize --validation <file> --ledger <file> [--out <file>]`
- `sdp-trace pr-review check --out <dir> --repo-id <safe-id> --change-ref <pr|mr|change-id> --base <sha> --head <sha> --diff <file> --profile <file> [--work-dir <dir>] [--allow-external-runner <runner>]...`
- `sdp-trace packet build-pr --source <github-actions|github-fixture> --out <dir> [--github-event <file>] [--checks-json <file>] [--artifacts-json <file>] [--route-manifest <file>] [--github-api-url <url>]`
- `sdp-trace packet build-github --github-input <file> --out <file>`
- `sdp-trace packet validate --bundle <file>`
- `sdp-trace packet check-demo --bundle <file>`
- `sdp-trace packet render --bundle <file> --out <file>`
- `sdp-trace validate-fixtures [root-dir]`

Do not add aliases, hidden switches, or workflow-specific wrappers as product
entrypoints unless this document and `--help` are updated in the same change.

## English Command Reference

| Command/profile | Purpose | Minimum invocation | Output and trust boundary |
| --- | --- | --- | --- |
| `wrap` | Observe one existing command as a trace run. | `sdp-trace wrap --name smoke -- /bin/echo ok` | Writes run artifacts; local observation only unless later bound by report/witness/profile checks. |
| `run` | Run a task-referenced command with an optional contract. | `sdp-trace run --task T1 --use-default-contract -- /bin/echo ok` | Writes task-linked run artifacts; missing contract evidence remains visible. |
| `dry-run` | Show what would run without writing run artifacts. | `sdp-trace dry-run -- /bin/echo ok` | Preview only; cannot support proof closure. |
| `preview` | Preview command/contract implications before execution. | `sdp-trace preview -- /bin/echo ok` | Read-only preview; any unavailable profile remains `not_assessed`. |
| `doctor` | Inspect local environment and contract prerequisites. | `sdp-trace doctor` | Emits structural readiness; offline or missing prerequisites can produce `cannot_verify`. |
| `doctor --profile github-actions-git-hooks-v1` | Inspect repository observer installation and proof state without relying on agent prompts. | `sdp-trace doctor --profile github-actions-git-hooks-v1 --out repo-observer-status.json` | Emits machine JSON plus a human table. Local hooks/config are `local_structural`; CI artifact proof remains `not_assessed` until uploaded artifacts are observed. |
| `install repo-observer` | Install portable repo observer files for local git hooks and GitHub Actions artifact upload. | `sdp-trace install repo-observer --profile github-actions-git-hooks-v1 --write --out repo-observer-status.json` | Dry-run by default. With `--write`, writes only the documented allowlist and refuses existing hooks-path conflicts unless `--force` is used after review. |
| `interaction relay` | Record an interaction event before forwarding message content to a command. | `printf 'correct plan boundary\n' \| sdp-trace interaction relay --task-id T1 --event-type corrective_feedback --out trace.json -- /bin/cat` | Writes a local interaction trace before delivery. If recording fails, forwarding does not happen. This records observation, not agent compliance. |
| `interaction import-transcript` | Import pre-classified interaction events from JSONL. | `sdp-trace interaction import-transcript --source preclassified-transcript-import --task-id T1 --events-jsonl events.jsonl --out trace.json` | Post-hoc import only. Source completeness can be `complete`, `partial`, `not_assessed`, or `cannot_verify`; free-form chat parsing is out of scope. |
| `interaction summarize` | Summarize interaction trace events and friction counts. | `sdp-trace interaction summarize --trace trace.json --out summary.json` | Read-only summary; friction counts are facts, not model, employee, or spec scores. |
| `observe setup` | Prepare a first-run harness observation session before delivery starts. | `sdp-trace observe setup --profile session-profile.json --out session-run --command 'opencode run ...'` | Writes setup metadata only. Setup is bounded by the session profile; no harness command runs here. Optional command preview is stored as a digest, not raw command output. |
| `observe collect` | Collect declared harness observation output after the normal harness command ran. | `sdp-trace observe collect --profile session-profile.json --run session-run` | Normalizes profile-declared `harness-event-v1` output, or supported raw OpenCode JSONL declared by `raw_event_format`, into an observed harness run without asking the operator to author events. Missing source output is `cannot_verify`; missing event families remain `not_assessed` at validation. |
| `observe session` | Convenience wrapper for a controlled first-run harness command. | `sdp-trace observe session --profile session-profile.json --out session-run -- <harness-command>` | Runs the command without stdin injection and does not retain stdout/stderr bodies by default. It records command/process provenance and then collects declared event output or supported raw OpenCode JSONL. |
| `harness observe` | Import a local harness lifecycle export without running or modifying the harness. | `sdp-trace harness observe --profile profile.json --source events.jsonl --out harness-run` | Reads explicit files only. Unsafe raw prompts, model responses, tokens, authenticated URLs, private paths, digest mismatches, and malformed JSONL fail before a run is written. |
| `harness validate` | Validate observed harness events against a selected profile. | `sdp-trace harness validate --profile profile.json --run harness-run --out validation.json` | Emits evidence facts with `pass`, `fail`, `not_assessed`, or `cannot_verify`; missing required event families are not passes. |
| `harness summarize` | Render a safe human summary of harness validation. | `sdp-trace harness summarize --validation validation.json` | Summary output is non-authoritative: it does not claim harness compliance, feature delivery, PR approval, merge approval, release readiness, or production trust. |
| `envelope summarize` | Summarize a delivery trace envelope across task, run, promise, LLM, mutation, stage, and friction refs. | `sdp-trace envelope summarize --envelope envelope.json --out summary.json` | Read-only over refs, including recorder run refs. It reports linked and `not_assessed` areas without readiness or quality verdicts. |
| `verify` | Verify one recorded run directory. | `sdp-trace verify .sdp-trace-runs/run-1` | Supports local structural assertions only; use JSON/state output for exact `observed`, `fail`, `not_assessed`, or `cannot_verify` interpretation. |
| `explain` | Render human-readable explanation for one run. | `sdp-trace explain .sdp-trace-runs/run-1` | Explanation is derived from run artifacts; it does not upgrade trust scope. |
| `query` | Query missing evidence or capture depth for a run. | `sdp-trace query --query missing-evidence .sdp-trace-runs/run-1` | Highlights gaps; missing rows are not passes. |
| `query-pack` | Build a forensic query package. | `sdp-trace query-pack --pack forensics-basic-v1 --run .sdp-trace-runs/run-1 --out query-pack.json` | Produces investigation package according to retained evidence; digest-only or redacted data limits reconstruction. |
| `query-pack explain` | Explain a forensic query-pack result. | `sdp-trace query-pack explain --result query-pack.json` | Explanation only; no new evidence is created. |
| `export cross-repo-posture` | Export cross-repository evidence posture for downstream analysis. | `sdp-trace export cross-repo-posture --profile cross-repo-evidence-posture-v1 --selection selection.json --out posture.json` | Aggregates selected inputs; degradation decisions remain outside `sdp-trace`. |
| `export cross-repo-posture explain` | Explain a cross-repo posture export. | `sdp-trace export cross-repo-posture explain --result posture.json` | Human-readable explanation; no policy decision. |
| `export telemetry` | Export posture facts as standard telemetry. | `sdp-trace export telemetry --profile prometheus-text-v1 --cross-repo-posture posture.json --out <file\|->` | Prometheus text output only; dashboards, alerts, reports, thresholds, and policy decisions remain downstream. |
| `assess --profile adapter-capture` | Assess adapter-capture evidence and overclaim risk. | `sdp-trace assess --profile adapter-capture --out assessment.json --run .sdp-trace-runs/run-1` | Can fail or return `cannot_verify` when adapter events are absent, agent-reported, or insufficient. |
| `assess --profile managed-harness` | Assess managed harness evidence against policy, registry, and witness inputs. | `sdp-trace assess --profile managed-harness --out assessment.json --contract contract.json --run .sdp-trace-runs/run-1 --adapter-registry registry.json --managed-policy policy.json --managed-witness witness.json` | Verifier facts plus exit behavior; external CI or policy owns block/allow decisions. Missing or stale witness usually produces `cannot_verify`. |
| `assess --profile forensic-retention` | Assess whether retained evidence supports forensic reconstruction. | `sdp-trace assess --profile forensic-retention --out assessment.json --run .sdp-trace-runs/run-1 --redaction-policy redaction.json` | Digest-only, unresolved redaction, or missing retention may fail or remain `cannot_verify`. |
| `assess --profile ci-artifact-observation` | Assess whether selected artifact families are CI-uploaded or lower-authority facts. | `sdp-trace assess --profile ci-artifact-observation --out observation.json --artifact-manifest artifact-manifest.json` | Fact package only; checked-in, local, prose, and agent-reported claims cannot satisfy selected `ci_uploaded` requirements. |
| `assess --profile authority-envelope` | Assess observed actions against a caller-selected authority envelope. | `sdp-trace assess --profile authority-envelope --authority-package authority-package.json --out authority-evaluation.json` | Emits authority facts only. `outside_authority`, `not_assessed`, and `cannot_verify` are not native merge, demo, discipline, or readiness decisions. |
| `assess preview` | Preview required inputs for an assessment profile. | `sdp-trace assess preview --profile managed-harness --run .sdp-trace-runs/run-1` | Read-only; does not evaluate authority or write proof. |
| `assess explain` | Explain an assessment result. | `sdp-trace assess explain --assessment-result assessment.json` | Explanation only; unsupported schema versions can produce `cannot_verify`. |
| `report` | Build `.sdp-trace-report/` from one run or run root. | `sdp-trace report --out .sdp-trace-report .sdp-trace-runs` | Packages observed data and gaps; report presence is not proof of completeness. |
| `gate` | Produce advisory/protected gate facts for a run root. | `sdp-trace gate --out .sdp-trace-report/gate-result.json .sdp-trace-runs` | Caveat: `gate` emits verifier facts and states; it is not a native merge, release, risk, or production-trust decision. Missing evidence stays `fail` or `cannot_verify`. |
| `witness --kind github-actions` | Bind report/run evidence to GitHub Actions identity when OIDC evidence is available. | `sdp-trace witness --kind github-actions --out ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs` | Caveat: CI witness is not external production trust by itself. Missing OIDC or binding data yields `cannot_verify`. |
| `witness --kind gitlab-ci` | Record GitLab CI witness profile evidence. | `sdp-trace witness --kind gitlab-ci --out gitlab-witness.json --witness-envelope envelope.json .sdp-trace-runs` | Caveat: profile output depends on supplied envelope and policy evidence; unsupported or incomplete fields stay `cannot_verify`. |
| `witness --kind buildkite` | Record Buildkite witness profile evidence. | `sdp-trace witness --kind buildkite --out buildkite-witness.json --witness-envelope envelope.json .sdp-trace-runs` | Caveat: CI-bound evidence is not a transparency log or release approval. |
| `witness --kind customer-pki` | Record customer PKI/private-equivalent witness evidence. | `sdp-trace witness --kind customer-pki --out customer-pki-witness.json --customer-pki-authority-policy policy.json --customer-pki-public-cert cert.pem --customer-pki-payload-digest <sha256> --customer-pki-freshness-evidence freshness.json .sdp-trace-runs` | Caveat: customer PKI requires accepted customer policy, key/cert material, payload digest, and freshness evidence; missing pieces are not production trust. |
| `release-proof` | Verify a source-bound local release manifest and proof artifact. | `sdp-trace release-proof --manifest examples/contract-foundation/contract-manifest.example.json --out release-proof.json` | Caveat: `source_bound_local_release` is narrower than external production trust; dirty/stale source, manifest mismatch, or absent source commit fails or returns `cannot_verify`. |
| `pr-review` | Build, run, synthesize, validate, and summarize automated PR review evidence over a frozen packet. | `sdp-trace pr-review check --out review --repo-id demo_repo --change-ref pr-123 --base <sha> --head <sha> --diff change.diff --profile examples/pr-review/trust-sensitive-default.profile.json` | Emits review-record facts only. `coverage_satisfied` is not merge approval; merge, release, risk acceptance, and human approval remain external. |
| `packet build-pr` | Build live PR packet artifacts from GitHub Actions context/API or offline GitHub fixtures without checked-in packet edits. | `sdp-trace packet build-pr --source github-actions --route-manifest route.json --out packet-artifacts` | Authoritative live-demo packet path for this slice. Writes `bundle.json`, `change-evidence-packet.md`, and `build-pr-result.json`; `PC-VERIFICATION` must bind to workflow/artifact evidence or remains `cannot_verify`. If `--checks-json` and `--artifacts-json` are omitted in GitHub Actions, the builder uses `GITHUB_TOKEN`/`GH_TOKEN` and the current run context to discover retained artifacts. |
| `packet build-github` | Build a packet bundle from a curated GitHub input fixture. | `sdp-trace packet build-github --github-input examples/change-evidence-packet/github-input.json --out packet-bundle.json` | Backfill/fixture authority only. It is useful for examples and regression tests, but not sufficient as live demo flight-recorder proof. |
| `packet` | Validate, demo-check, and render Change Evidence Packet v0 bundles. | `sdp-trace packet render --bundle examples/change-evidence-packet/happy-path.bundle.json --out change-evidence-packet.md` | Canonical packet rendering only; row states and residual gaps are not merge, release, compliance, production trust, or semantic quality approval. `packet check-demo` is limited to the 007 first-packet minimum bar and is not a general approval gate. |
| `validate-fixtures` | Validate checked trace-run fixture directories. | `sdp-trace validate-fixtures examples/agentic-sdlc` | Structural fixture validation only; does not prove customer production readiness. Point it at a fixture root that contains run directories. |

## Russian Command Reference

State: `deferred_scope` for full bilingual parity. This table preserves the
current high-use Russian quick reference, but the English command reference
above is the canonical command contract until the Russian table is expanded to
cover every live `--help` command family. Do not claim bilingual command parity
from this section.

| Команда/профиль | Назначение | Минимальный запуск | Граница вывода и доверия |
| --- | --- | --- | --- |
| `wrap` | Наблюдает одну существующую команду как trace run. | `sdp-trace wrap --name smoke -- /bin/echo ok` | Пишет run artifacts; это local observation, пока отчет, witness или профиль не добавят другую проверку. |
| `run` | Запускает команду, связанную с task ref. | `sdp-trace run --task T1 --use-default-contract -- /bin/echo ok` | Пишет task-linked run artifacts; missing contract evidence остается видимым. |
| `dry-run` | Показывает запуск без записи run artifacts. | `sdp-trace dry-run -- /bin/echo ok` | Только preview; proof closure не поддерживает. |
| `preview` | Показывает command/contract implications до выполнения. | `sdp-trace preview -- /bin/echo ok` | Read-only preview; недоступный профиль остается `not_assessed`. |
| `doctor` | Проверяет локальную среду и prerequisites. | `sdp-trace doctor` | Structural readiness; offline или missing prerequisites могут дать `cannot_verify`. |
| `doctor --profile github-actions-git-hooks-v1` | Проверяет установку repo observer и proof state без опоры на промпты агентов. | `sdp-trace doctor --profile github-actions-git-hooks-v1 --out repo-observer-status.json` | Пишет machine JSON и human table. Local hooks/config остаются `local_structural`; CI artifact proof остается `not_assessed`, пока uploaded artifacts не наблюдены. |
| `install repo-observer` | Устанавливает portable repo observer files для git hooks и GitHub Actions artifact upload. | `sdp-trace install repo-observer --profile github-actions-git-hooks-v1 --write --out repo-observer-status.json` | По умолчанию dry-run. С `--write` пишет только documented allowlist и отказывается от hooks-path conflicts без reviewed `--force`. |
| `interaction relay` | Записывает interaction event до передачи message content в команду. | `printf 'correct plan boundary\n' \| sdp-trace interaction relay --task-id T1 --event-type corrective_feedback --out trace.json -- /bin/cat` | Пишет local interaction trace до delivery. Если запись не удалась, forwarding не происходит. Это observation, не compliance. |
| `interaction import-transcript` | Импортирует pre-classified interaction events из JSONL. | `sdp-trace interaction import-transcript --source preclassified-transcript-import --task-id T1 --events-jsonl events.jsonl --out trace.json` | Только post-hoc import. Completeness источника остается явной: `complete`, `partial`, `not_assessed` или `cannot_verify`. |
| `interaction summarize` | Суммирует interaction trace events и friction counts. | `sdp-trace interaction summarize --trace trace.json --out summary.json` | Read-only summary; friction counts являются фактами, не score модели, сотрудника или спеки. |
| `envelope summarize` | Суммирует delivery trace envelope по task, run, promise, LLM, mutation, stage и friction refs. | `sdp-trace envelope summarize --envelope envelope.json --out summary.json` | Read-only по refs, включая recorder run refs. Показывает linked и `not_assessed` зоны без readiness/quality verdict. |
| `verify` | Проверяет один recorded run directory. | `sdp-trace verify .sdp-trace-runs/run-1` | Поддерживает local structural assertions; для точных states используйте JSON/state output. |
| `explain` | Объясняет один run. | `sdp-trace explain .sdp-trace-runs/run-1` | Объяснение не повышает trust scope. |
| `query` | Ищет missing evidence или capture depth. | `sdp-trace query --query missing-evidence .sdp-trace-runs/run-1` | Показывает gaps; missing rows не являются pass. |
| `query-pack` | Собирает forensic query package. | `sdp-trace query-pack --pack forensics-basic-v1 --run .sdp-trace-runs/run-1 --out query-pack.json` | Расследование ограничено тем, какое evidence было retained, redacted или digest-only. |
| `query-pack explain` | Объясняет query-pack result. | `sdp-trace query-pack explain --result query-pack.json` | Только explanation; новое evidence не создается. |
| `export cross-repo-posture` | Экспортирует cross-repo evidence posture. | `sdp-trace export cross-repo-posture --profile cross-repo-evidence-posture-v1 --selection selection.json --out posture.json` | Агрегирует выбранные inputs; degradation decision остается вне `sdp-trace`. |
| `export cross-repo-posture explain` | Объясняет cross-repo export. | `sdp-trace export cross-repo-posture explain --result posture.json` | Human-readable explanation, не policy decision. |
| `export telemetry` | Экспортирует posture facts в стандартную telemetry surface. | `sdp-trace export telemetry --profile prometheus-text-v1 --cross-repo-posture posture.json --out <file\|->` | Только Prometheus text; dashboards, alerts, reports, thresholds и policy decisions остаются downstream. |
| `assess --profile adapter-capture` | Проверяет adapter-capture evidence и overclaim risk. | `sdp-trace assess --profile adapter-capture --out assessment.json --run .sdp-trace-runs/run-1` | Может дать `fail` или `cannot_verify`, если adapter events absent, agent-reported или insufficient. |
| `assess --profile managed-harness` | Проверяет managed harness evidence по policy, registry и witness. | `sdp-trace assess --profile managed-harness --out assessment.json --contract contract.json --run .sdp-trace-runs/run-1 --adapter-registry registry.json --managed-policy policy.json --managed-witness witness.json` | Это verifier facts и exit behavior; block/allow решает external CI или policy. |
| `assess --profile forensic-retention` | Проверяет, хватает ли retained evidence для forensic reconstruction. | `sdp-trace assess --profile forensic-retention --out assessment.json --run .sdp-trace-runs/run-1 --redaction-policy redaction.json` | Digest-only, unresolved redaction или missing retention могут дать `fail` или `cannot_verify`. |
| `assess --profile ci-artifact-observation` | Проверяет, какие artifact families реально наблюдаются как uploaded CI artifacts, а какие остаются lower-authority facts. | `sdp-trace assess --profile ci-artifact-observation --out observation.json --artifact-manifest artifact-manifest.json` | Только fact package; downstream gates принимают policy decisions отдельно. |
| `assess --profile authority-envelope` | Проверяет observed actions по caller-selected authority envelope. | `sdp-trace assess --profile authority-envelope --authority-package authority-package.json --out authority-evaluation.json` | Пишет только authority facts; block/contamination/readiness decisions остаются downstream. |
| `assess preview` | Показывает required inputs для assessment profile. | `sdp-trace assess preview --profile managed-harness --run .sdp-trace-runs/run-1` | Read-only; authority не проверяет и proof не пишет. |
| `assess explain` | Объясняет assessment result. | `sdp-trace assess explain --assessment-result assessment.json` | Только explanation; unsupported schema может дать `cannot_verify`. |
| `report` | Собирает `.sdp-trace-report/` из run или run root. | `sdp-trace report --out .sdp-trace-report .sdp-trace-runs` | Упаковывает observed data и gaps; сам отчет не доказывает полноту. |
| `gate` | Пишет advisory/protected gate facts. | `sdp-trace gate --out .sdp-trace-report/gate-result.json .sdp-trace-runs` | Caveat: `gate` не является native merge, release, risk или production-trust decision. |
| `witness --kind github-actions` | Связывает evidence с GitHub Actions identity при наличии OIDC. | `sdp-trace witness --kind github-actions --out ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs` | Caveat: CI witness сам по себе не external production trust; missing OIDC дает `cannot_verify`. |
| `witness --kind gitlab-ci` | Записывает GitLab CI witness profile evidence. | `sdp-trace witness --kind gitlab-ci --out gitlab-witness.json --witness-envelope envelope.json .sdp-trace-runs` | Incomplete envelope или policy evidence остаются `cannot_verify`. |
| `witness --kind buildkite` | Записывает Buildkite witness profile evidence. | `sdp-trace witness --kind buildkite --out buildkite-witness.json --witness-envelope envelope.json .sdp-trace-runs` | CI-bound evidence не является transparency log или release approval. |
| `witness --kind customer-pki` | Записывает customer PKI/private-equivalent witness evidence. | `sdp-trace witness --kind customer-pki --out customer-pki-witness.json --customer-pki-authority-policy policy.json --customer-pki-public-cert cert.pem --customer-pki-payload-digest <sha256> --customer-pki-freshness-evidence freshness.json .sdp-trace-runs` | Требует accepted customer policy, key/cert material, payload digest и freshness evidence. |
| `release-proof` | Проверяет source-bound local release manifest. | `sdp-trace release-proof --manifest examples/contract-foundation/contract-manifest.example.json --out release-proof.json` | Caveat: `source_bound_local_release` уже, чем external production trust; dirty/stale source или manifest mismatch не проходят. |
| `pr-review` | Собирает, запускает, синтезирует, валидирует и суммирует automated PR review evidence. | `sdp-trace pr-review check --out review --repo-id demo_repo --change-ref pr-123 --base <sha> --head <sha> --diff change.diff --profile examples/pr-review/trust-sensitive-default.profile.json` | Review evidence only; не merge approval, не human approval и не release decision. |
| `validate-fixtures` | Валидирует trace-run fixture directories. | `sdp-trace validate-fixtures examples/agentic-sdlc` | Только structural fixture validation; не доказывает customer production readiness. Указывайте fixture root, где есть run directories. |

## State And Exit Code Contract

### Result States

These are the verifier result states. Every command that reports verifier
outcome uses one of these. They map to exit codes.

| Result state | Exit code | Meaning |
| --- | --- | --- |
| `observed` | `0` | Verifier evidence met required checks for the selected local profile. |
| `pass` | `0` | Selected profile concluded successfully where the command contract uses pass/fail states. |
| `fail` | `1` | Verifier evidence conflicted or was insufficient for required checks. |
| `not_assessed` | `0` | State was outside the run scope; it does not imply success or evidence. May return `0` only when the command completed and the unassessed state is explicitly outside the selected profile or run scope. |
| `cannot_verify` | `3` | Verifier could not execute a required check or lacked required evidence. |

Exit code `2` is reserved for usage error / invalid command invocation and is
not a verifier result state.

### Telemetry Labels

These labels describe evidence availability, not verifier outcomes. They do
not have exit-code mappings.

- `missing_telemetry`: a telemetry stream or metric was expected but not found.
  Used by query-pack and managed-harness capture-depth reporting. It is not a
  verifier result state; the corresponding verifier result may be `not_assessed`
  or `cannot_verify` depending on whether the telemetry was required.

### Completeness Markers

These describe source or input completeness, not verifier outcomes.

- `complete`, `partial`: source completeness for `interaction import-transcript`.
  They describe the imported data set, not a verifier pass/fail.

### PR-Review Sub-States

These are command-specific coverage states reported by `pr-review`. They are
not verifier result states and do not have exit-code mappings.

- `coverage_satisfied`: the review packet reached the coverage threshold for the profile.
- `coverage_partial`: some planes were reviewed but coverage did not reach the threshold.
- `coverage_unresolved`: coverage could not be determined.

### External Verdict Sub-States

These appear in policy-consumer output and concept docs. They are not verifier
result states.

- `warn`: evidence exists but risk remains. Defined in `docs/concepts.md` as an
  External Verdict value. The corresponding verifier result is typically
  `observed` or `pass` with an advisory note.

### Integration And Adapter Labels

These describe integration or adapter status. They are not verifier result states.

- `not_integrated`: an expected adapter or integration is absent.
- `unsupported`: a format, schema version, or configuration is not supported.

### Authority Scope Labels

These describe the reporting boundary, not a verifier result.

- `outside_authority`: an observed action is outside the caller-selected
  authority envelope. Emitted by `assess --profile authority-envelope`.
- `local_dirty_structural_only`: dirty-checkout structural output. Valid only
  for local shape/debug inspection; cannot support source-bound or external
  trust closure.

### Claim Tag States

These appear in claim tags (`docs/claim-authoring.md`). They describe the
freshness of a historical claim, not a verifier result.

- `stale`: a historical closure record contradicts the current verifier output.
  Used in `sdp-trace-claim` tags to mark claims that need re-verification.

### Quality Gate Statuses

These appear in local quality gate discussions. They describe a standing
assessment condition, not a verifier result.

- `assessed_gap`: a metric or threshold is below target but tracked as a known
  gap. Used in Maintainability Index discussions to indicate historical
  below-threshold functions/files that are tracked but not claimed as passing.

A checked-in proof JSON is an audit artifact, not authority. Authority is replayed
only from live Go verifier output and the canonical command/state contract above.

## Air-Gapped Fixture Guidance

Air-gapped evidence is a fixture and customer-policy pattern, not a native
`witness --kind air-gapped` command. Use `customer-pki` or an accepted private
equivalent with explicit authority policy, payload digest, freshness evidence,
and retained audit references. If those are absent, record `not_assessed` or
`cannot_verify`; do not claim external production trust.

## Overclaim And Forbidden Claims

See [`docs/overclaim-checklist.md`](overclaim-checklist.md) for the canonical
overclaim checklist.
