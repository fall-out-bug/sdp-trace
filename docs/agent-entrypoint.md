# Agent Entrypoint

Use these commands and profile IDs to select what `sdp-trace` can prove without
introducing harness assumptions. The active entrypoint is the Go CLI reported by
`go run ./cmd/sdp-trace --help`; update this document only when the help surface
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

Current command surface:

- `go test ./...`
- `go run ./cmd/sdp-trace --help`
- `go run ./cmd/sdp-trace wrap --name <name> [--contract <file>] [--output-dir <dir>] -- <command...>`
- `go run ./cmd/sdp-trace run --task <task-ref> [--contract <file> | --use-default-contract] -- <command...>`
- `go run ./cmd/sdp-trace dry-run [--contract <file> | --use-default-contract] -- <command...>`
- `go run ./cmd/sdp-trace preview [--contract <file> | --use-default-contract] -- <command...>`
- `go run ./cmd/sdp-trace doctor [--contract <file>]`
- `go run ./cmd/sdp-trace doctor --profile github-actions-git-hooks-v1 [--out <file>]`
- `go run ./cmd/sdp-trace install repo-observer --profile github-actions-git-hooks-v1 [--repository-id <safe-id>] [--write] [--force] [--out <file>]`
- `go run ./cmd/sdp-trace interaction relay --task-id <safe-id> --event-type <type> --out <file> -- <forward-command...>`
- `go run ./cmd/sdp-trace interaction import-transcript --source preclassified-transcript-import --task-id <safe-id> --events-jsonl <file> --out <file>`
- `go run ./cmd/sdp-trace interaction summarize --trace <file> [--out <file>]`
- `go run ./cmd/sdp-trace harness observe --profile <harness-profile.json> --source <harness-events.jsonl> --out <run-dir>`
- `go run ./cmd/sdp-trace harness validate --profile <harness-profile.json> --run <run-dir> --out <validation.json>`
- `go run ./cmd/sdp-trace harness summarize --validation <validation.json>`
- `go run ./cmd/sdp-trace envelope summarize --envelope <file> [--out <file>]`
- `go run ./cmd/sdp-trace verify <run-dir>`
- `go run ./cmd/sdp-trace explain <run-dir>`
- `go run ./cmd/sdp-trace query --query <missing-evidence|capture-depth> <run-dir>`
- `go run ./cmd/sdp-trace query-pack --pack forensics-basic-v1 --run <run-dir> --out <file>`
- `go run ./cmd/sdp-trace query-pack explain --result <file>`
- `go run ./cmd/sdp-trace export cross-repo-posture --profile cross-repo-evidence-posture-v1 --selection <file> --out <file>`
- `go run ./cmd/sdp-trace export cross-repo-posture explain --result <file>`
- `go run ./cmd/sdp-trace export telemetry --profile prometheus-text-v1 --cross-repo-posture <file> --out <file|->`
- `go run ./cmd/sdp-trace assess --profile adapter-capture --out <file> --run <run-dir>`
- `go run ./cmd/sdp-trace assess --profile managed-harness --out <file> --contract <file> --run <run-dir> --adapter-registry <file> --managed-policy <file> --managed-witness <file>`
- `go run ./cmd/sdp-trace assess --profile forensic-retention --out <file> --run <run-dir> --redaction-policy <file>`
- `go run ./cmd/sdp-trace assess preview --profile <adapter-capture|managed-harness|forensic-retention> [profile inputs]`
- `go run ./cmd/sdp-trace assess explain --assessment-result <file>`
- `go run ./cmd/sdp-trace report --out <dir> <runs-root-or-run-dir>`
- `go run ./cmd/sdp-trace gate --out <file> <runs-root-or-run-dir>`
- `go run ./cmd/sdp-trace witness --kind <github-actions|gitlab-ci|buildkite|customer-pki> --out <file> [--report-dir <dir>] [--witness-envelope <file>] [--customer-pki-authority-policy <file>] [--customer-pki-public-cert <file> | --customer-pki-public-key <file>] [--customer-pki-payload-digest <sha256>] [--customer-pki-freshness-evidence <file>] <runs-root-or-run-dir>`
- `go run ./cmd/sdp-trace release-proof --manifest <file> --out <file>`
- `go run ./cmd/sdp-trace validate-fixtures [root-dir]`

Do not add aliases, hidden switches, or workflow-specific wrappers as product
entrypoints unless this document and `--help` are updated in the same change.

## English Command Reference

| Command/profile | Purpose | Minimum invocation | Output and trust boundary |
| --- | --- | --- | --- |
| `wrap` | Observe one existing command as a trace run. | `go run ./cmd/sdp-trace wrap --name smoke -- /bin/echo ok` | Writes run artifacts; local observation only unless later bound by report/witness/profile checks. |
| `run` | Run a task-referenced command with an optional contract. | `go run ./cmd/sdp-trace run --task T1 -- /bin/echo ok` | Writes task-linked run artifacts; missing contract evidence remains visible. |
| `dry-run` | Show what would run without writing run artifacts. | `go run ./cmd/sdp-trace dry-run -- /bin/echo ok` | Preview only; cannot support proof closure. |
| `preview` | Preview command/contract implications before execution. | `go run ./cmd/sdp-trace preview -- /bin/echo ok` | Read-only preview; any unavailable profile remains `not_assessed`. |
| `doctor` | Inspect local environment and contract prerequisites. | `go run ./cmd/sdp-trace doctor` | Emits structural readiness; offline or missing prerequisites can produce `cannot_verify`. |
| `doctor --profile github-actions-git-hooks-v1` | Inspect repository observer installation and proof state without relying on agent prompts. | `go run ./cmd/sdp-trace doctor --profile github-actions-git-hooks-v1 --out repo-observer-status.json` | Emits machine JSON plus a human table. Local hooks/config are `local_structural`; CI artifact proof remains `not_assessed` until uploaded artifacts are observed. |
| `install repo-observer` | Install portable repo observer files for local git hooks and GitHub Actions artifact upload. | `go run ./cmd/sdp-trace install repo-observer --profile github-actions-git-hooks-v1 --write --out repo-observer-status.json` | Dry-run by default. With `--write`, writes only the documented allowlist and refuses existing hooks-path conflicts unless `--force` is used after review. |
| `interaction relay` | Record an interaction event before forwarding message content to a command. | `printf 'correct plan boundary\n' \| go run ./cmd/sdp-trace interaction relay --task-id T1 --event-type corrective_feedback --out trace.json -- /bin/cat` | Writes a local interaction trace before delivery. If recording fails, forwarding does not happen. This records observation, not agent compliance. |
| `interaction import-transcript` | Import pre-classified interaction events from JSONL. | `go run ./cmd/sdp-trace interaction import-transcript --source preclassified-transcript-import --task-id T1 --events-jsonl events.jsonl --out trace.json` | Post-hoc import only. Source completeness can be `complete`, `partial`, `not_assessed`, or `cannot_verify`; free-form chat parsing is out of scope. |
| `interaction summarize` | Summarize interaction trace events and friction counts. | `go run ./cmd/sdp-trace interaction summarize --trace trace.json --out summary.json` | Read-only summary; friction counts are facts, not model, employee, or spec scores. |
| `harness observe` | Import a local harness lifecycle export without running or modifying the harness. | `go run ./cmd/sdp-trace harness observe --profile profile.json --source events.jsonl --out harness-run` | Reads explicit files only. Unsafe raw prompts, model responses, tokens, authenticated URLs, private paths, digest mismatches, and malformed JSONL fail before a run is written. |
| `harness validate` | Validate observed harness events against a selected profile. | `go run ./cmd/sdp-trace harness validate --profile profile.json --run harness-run --out validation.json` | Emits evidence facts with `pass`, `fail`, `not_assessed`, or `cannot_verify`; missing required event families are not passes. |
| `harness summarize` | Render a safe human summary of harness validation. | `go run ./cmd/sdp-trace harness summarize --validation validation.json` | Summary output is non-authoritative: it does not claim harness compliance, feature delivery, PR approval, merge approval, release readiness, or production trust. |
| `envelope summarize` | Summarize a delivery trace envelope across task, run, promise, LLM, mutation, stage, and friction refs. | `go run ./cmd/sdp-trace envelope summarize --envelope envelope.json --out summary.json` | Read-only over refs, including recorder run refs. It reports linked and `not_assessed` areas without readiness or quality verdicts. |
| `verify` | Verify one recorded run directory. | `go run ./cmd/sdp-trace verify .sdp-trace-runs/run-1` | Supports local structural assertions only; use JSON/state output for exact `observed`, `fail`, `not_assessed`, or `cannot_verify` interpretation. |
| `explain` | Render human-readable explanation for one run. | `go run ./cmd/sdp-trace explain .sdp-trace-runs/run-1` | Explanation is derived from run artifacts; it does not upgrade trust scope. |
| `query` | Query missing evidence or capture depth for a run. | `go run ./cmd/sdp-trace query --query missing-evidence .sdp-trace-runs/run-1` | Highlights gaps; missing rows are not passes. |
| `query-pack` | Build a forensic query package. | `go run ./cmd/sdp-trace query-pack --pack forensics-basic-v1 --run .sdp-trace-runs/run-1 --out query-pack.json` | Produces investigation package according to retained evidence; digest-only or redacted data limits reconstruction. |
| `query-pack explain` | Explain a forensic query-pack result. | `go run ./cmd/sdp-trace query-pack explain --result query-pack.json` | Explanation only; no new evidence is created. |
| `export cross-repo-posture` | Export cross-repository evidence posture for downstream analysis. | `go run ./cmd/sdp-trace export cross-repo-posture --profile cross-repo-evidence-posture-v1 --selection selection.json --out posture.json` | Aggregates selected inputs; degradation decisions remain outside `sdp-trace`. |
| `export cross-repo-posture explain` | Explain a cross-repo posture export. | `go run ./cmd/sdp-trace export cross-repo-posture explain --result posture.json` | Human-readable explanation; no policy decision. |
| `export telemetry` | Export posture facts as standard telemetry. | `go run ./cmd/sdp-trace export telemetry --profile prometheus-text-v1 --cross-repo-posture posture.json --out <file\|->` | Prometheus text output only; dashboards, alerts, reports, thresholds, and policy decisions remain downstream. |
| `assess --profile adapter-capture` | Assess adapter-capture evidence and overclaim risk. | `go run ./cmd/sdp-trace assess --profile adapter-capture --out assessment.json --run .sdp-trace-runs/run-1` | Can fail or return `cannot_verify` when adapter events are absent, agent-reported, or insufficient. |
| `assess --profile managed-harness` | Assess managed harness evidence against policy, registry, and witness inputs. | `go run ./cmd/sdp-trace assess --profile managed-harness --out assessment.json --contract contract.json --run .sdp-trace-runs/run-1 --adapter-registry registry.json --managed-policy policy.json --managed-witness witness.json` | Verifier facts plus exit behavior; external CI or policy owns block/allow decisions. Missing or stale witness usually produces `cannot_verify`. |
| `assess --profile forensic-retention` | Assess whether retained evidence supports forensic reconstruction. | `go run ./cmd/sdp-trace assess --profile forensic-retention --out assessment.json --run .sdp-trace-runs/run-1 --redaction-policy redaction.json` | Digest-only, unresolved redaction, or missing retention may fail or remain `cannot_verify`. |
| `assess --profile ci-artifact-observation` | Assess whether selected artifact families are CI-uploaded or lower-authority facts. | `go run ./cmd/sdp-trace assess --profile ci-artifact-observation --out observation.json --artifact-manifest artifact-manifest.json` | Fact package only; checked-in, local, prose, and agent-reported claims cannot satisfy selected `ci_uploaded` requirements. |
| `assess --profile authority-envelope` | Assess observed actions against a caller-selected authority envelope. | `go run ./cmd/sdp-trace assess --profile authority-envelope --authority-package authority-package.json --out authority-evaluation.json` | Emits authority facts only. `outside_authority`, `not_assessed`, and `cannot_verify` are not native merge, demo, discipline, or readiness decisions. |
| `assess preview` | Preview required inputs for an assessment profile. | `go run ./cmd/sdp-trace assess preview --profile managed-harness --run .sdp-trace-runs/run-1` | Read-only; does not evaluate authority or write proof. |
| `assess explain` | Explain an assessment result. | `go run ./cmd/sdp-trace assess explain --assessment-result assessment.json` | Explanation only; unsupported schema versions can produce `cannot_verify`. |
| `report` | Build `.sdp-trace-report/` from one run or run root. | `go run ./cmd/sdp-trace report --out .sdp-trace-report .sdp-trace-runs` | Packages observed data and gaps; report presence is not proof of completeness. |
| `gate` | Produce advisory/protected gate facts for a run root. | `go run ./cmd/sdp-trace gate --out .sdp-trace-report/gate-result.json .sdp-trace-runs` | Caveat: `gate` emits verifier facts and states; it is not a native merge, release, risk, or production-trust decision. Missing evidence stays `fail` or `cannot_verify`. |
| `witness --kind github-actions` | Bind report/run evidence to GitHub Actions identity when OIDC evidence is available. | `go run ./cmd/sdp-trace witness --kind github-actions --out ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs` | Caveat: CI witness is not external production trust by itself. Missing OIDC or binding data yields `cannot_verify`. |
| `witness --kind gitlab-ci` | Record GitLab CI witness profile evidence. | `go run ./cmd/sdp-trace witness --kind gitlab-ci --out gitlab-witness.json --witness-envelope envelope.json .sdp-trace-runs` | Caveat: profile output depends on supplied envelope and policy evidence; unsupported or incomplete fields stay `cannot_verify`. |
| `witness --kind buildkite` | Record Buildkite witness profile evidence. | `go run ./cmd/sdp-trace witness --kind buildkite --out buildkite-witness.json --witness-envelope envelope.json .sdp-trace-runs` | Caveat: CI-bound evidence is not a transparency log or release approval. |
| `witness --kind customer-pki` | Record customer PKI/private-equivalent witness evidence. | `go run ./cmd/sdp-trace witness --kind customer-pki --out customer-pki-witness.json --customer-pki-authority-policy policy.json --customer-pki-public-cert cert.pem --customer-pki-payload-digest <sha256> --customer-pki-freshness-evidence freshness.json .sdp-trace-runs` | Caveat: customer PKI requires accepted customer policy, key/cert material, payload digest, and freshness evidence; missing pieces are not production trust. |
| `release-proof` | Verify a source-bound local release manifest and proof artifact. | `go run ./cmd/sdp-trace release-proof --manifest examples/contract-foundation/contract-manifest.example.json --out release-proof.json` | Caveat: `source_bound_local_release` is narrower than external production trust; dirty/stale source, manifest mismatch, or absent source commit fails or returns `cannot_verify`. |
| `validate-fixtures` | Validate checked fixture directories. | `go run ./cmd/sdp-trace validate-fixtures examples` | Structural fixture validation only; does not prove customer production readiness. |

## Russian Command Reference

| Команда/профиль | Назначение | Минимальный запуск | Граница вывода и доверия |
| --- | --- | --- | --- |
| `wrap` | Наблюдает одну существующую команду как trace run. | `go run ./cmd/sdp-trace wrap --name smoke -- /bin/echo ok` | Пишет run artifacts; это local observation, пока отчет, witness или профиль не добавят другую проверку. |
| `run` | Запускает команду, связанную с task ref. | `go run ./cmd/sdp-trace run --task T1 -- /bin/echo ok` | Пишет task-linked run artifacts; missing contract evidence остается видимым. |
| `dry-run` | Показывает запуск без записи run artifacts. | `go run ./cmd/sdp-trace dry-run -- /bin/echo ok` | Только preview; proof closure не поддерживает. |
| `preview` | Показывает command/contract implications до выполнения. | `go run ./cmd/sdp-trace preview -- /bin/echo ok` | Read-only preview; недоступный профиль остается `not_assessed`. |
| `doctor` | Проверяет локальную среду и prerequisites. | `go run ./cmd/sdp-trace doctor` | Structural readiness; offline или missing prerequisites могут дать `cannot_verify`. |
| `doctor --profile github-actions-git-hooks-v1` | Проверяет установку repo observer и proof state без опоры на промпты агентов. | `go run ./cmd/sdp-trace doctor --profile github-actions-git-hooks-v1 --out repo-observer-status.json` | Пишет machine JSON и human table. Local hooks/config остаются `local_structural`; CI artifact proof остается `not_assessed`, пока uploaded artifacts не наблюдены. |
| `install repo-observer` | Устанавливает portable repo observer files для git hooks и GitHub Actions artifact upload. | `go run ./cmd/sdp-trace install repo-observer --profile github-actions-git-hooks-v1 --write --out repo-observer-status.json` | По умолчанию dry-run. С `--write` пишет только documented allowlist и отказывается от hooks-path conflicts без reviewed `--force`. |
| `interaction relay` | Записывает interaction event до передачи message content в команду. | `printf 'correct plan boundary\n' \| go run ./cmd/sdp-trace interaction relay --task-id T1 --event-type corrective_feedback --out trace.json -- /bin/cat` | Пишет local interaction trace до delivery. Если запись не удалась, forwarding не происходит. Это observation, не compliance. |
| `interaction import-transcript` | Импортирует pre-classified interaction events из JSONL. | `go run ./cmd/sdp-trace interaction import-transcript --source preclassified-transcript-import --task-id T1 --events-jsonl events.jsonl --out trace.json` | Только post-hoc import. Completeness источника остается явной: `complete`, `partial`, `not_assessed` или `cannot_verify`. |
| `interaction summarize` | Суммирует interaction trace events и friction counts. | `go run ./cmd/sdp-trace interaction summarize --trace trace.json --out summary.json` | Read-only summary; friction counts являются фактами, не score модели, сотрудника или спеки. |
| `envelope summarize` | Суммирует delivery trace envelope по task, run, promise, LLM, mutation, stage и friction refs. | `go run ./cmd/sdp-trace envelope summarize --envelope envelope.json --out summary.json` | Read-only по refs, включая recorder run refs. Показывает linked и `not_assessed` зоны без readiness/quality verdict. |
| `verify` | Проверяет один recorded run directory. | `go run ./cmd/sdp-trace verify .sdp-trace-runs/run-1` | Поддерживает local structural assertions; для точных states используйте JSON/state output. |
| `explain` | Объясняет один run. | `go run ./cmd/sdp-trace explain .sdp-trace-runs/run-1` | Объяснение не повышает trust scope. |
| `query` | Ищет missing evidence или capture depth. | `go run ./cmd/sdp-trace query --query missing-evidence .sdp-trace-runs/run-1` | Показывает gaps; missing rows не являются pass. |
| `query-pack` | Собирает forensic query package. | `go run ./cmd/sdp-trace query-pack --pack forensics-basic-v1 --run .sdp-trace-runs/run-1 --out query-pack.json` | Расследование ограничено тем, какое evidence было retained, redacted или digest-only. |
| `query-pack explain` | Объясняет query-pack result. | `go run ./cmd/sdp-trace query-pack explain --result query-pack.json` | Только explanation; новое evidence не создается. |
| `export cross-repo-posture` | Экспортирует cross-repo evidence posture. | `go run ./cmd/sdp-trace export cross-repo-posture --profile cross-repo-evidence-posture-v1 --selection selection.json --out posture.json` | Агрегирует выбранные inputs; degradation decision остается вне `sdp-trace`. |
| `export cross-repo-posture explain` | Объясняет cross-repo export. | `go run ./cmd/sdp-trace export cross-repo-posture explain --result posture.json` | Human-readable explanation, не policy decision. |
| `export telemetry` | Экспортирует posture facts в стандартную telemetry surface. | `go run ./cmd/sdp-trace export telemetry --profile prometheus-text-v1 --cross-repo-posture posture.json --out <file\|->` | Только Prometheus text; dashboards, alerts, reports, thresholds и policy decisions остаются downstream. |
| `assess --profile adapter-capture` | Проверяет adapter-capture evidence и overclaim risk. | `go run ./cmd/sdp-trace assess --profile adapter-capture --out assessment.json --run .sdp-trace-runs/run-1` | Может дать `fail` или `cannot_verify`, если adapter events absent, agent-reported или insufficient. |
| `assess --profile managed-harness` | Проверяет managed harness evidence по policy, registry и witness. | `go run ./cmd/sdp-trace assess --profile managed-harness --out assessment.json --contract contract.json --run .sdp-trace-runs/run-1 --adapter-registry registry.json --managed-policy policy.json --managed-witness witness.json` | Это verifier facts и exit behavior; block/allow решает external CI или policy. |
| `assess --profile forensic-retention` | Проверяет, хватает ли retained evidence для forensic reconstruction. | `go run ./cmd/sdp-trace assess --profile forensic-retention --out assessment.json --run .sdp-trace-runs/run-1 --redaction-policy redaction.json` | Digest-only, unresolved redaction или missing retention могут дать `fail` или `cannot_verify`. |
| `assess --profile ci-artifact-observation` | Проверяет, какие artifact families реально наблюдаются как uploaded CI artifacts, а какие остаются lower-authority facts. | `go run ./cmd/sdp-trace assess --profile ci-artifact-observation --out observation.json --artifact-manifest artifact-manifest.json` | Только fact package; downstream gates принимают policy decisions отдельно. |
| `assess --profile authority-envelope` | Проверяет observed actions по caller-selected authority envelope. | `go run ./cmd/sdp-trace assess --profile authority-envelope --authority-package authority-package.json --out authority-evaluation.json` | Пишет только authority facts; block/contamination/readiness decisions остаются downstream. |
| `assess preview` | Показывает required inputs для assessment profile. | `go run ./cmd/sdp-trace assess preview --profile managed-harness --run .sdp-trace-runs/run-1` | Read-only; authority не проверяет и proof не пишет. |
| `assess explain` | Объясняет assessment result. | `go run ./cmd/sdp-trace assess explain --assessment-result assessment.json` | Только explanation; unsupported schema может дать `cannot_verify`. |
| `report` | Собирает `.sdp-trace-report/` из run или run root. | `go run ./cmd/sdp-trace report --out .sdp-trace-report .sdp-trace-runs` | Упаковывает observed data и gaps; сам отчет не доказывает полноту. |
| `gate` | Пишет advisory/protected gate facts. | `go run ./cmd/sdp-trace gate --out .sdp-trace-report/gate-result.json .sdp-trace-runs` | Caveat: `gate` не является native merge, release, risk или production-trust decision. |
| `witness --kind github-actions` | Связывает evidence с GitHub Actions identity при наличии OIDC. | `go run ./cmd/sdp-trace witness --kind github-actions --out ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs` | Caveat: CI witness сам по себе не external production trust; missing OIDC дает `cannot_verify`. |
| `witness --kind gitlab-ci` | Записывает GitLab CI witness profile evidence. | `go run ./cmd/sdp-trace witness --kind gitlab-ci --out gitlab-witness.json --witness-envelope envelope.json .sdp-trace-runs` | Incomplete envelope или policy evidence остаются `cannot_verify`. |
| `witness --kind buildkite` | Записывает Buildkite witness profile evidence. | `go run ./cmd/sdp-trace witness --kind buildkite --out buildkite-witness.json --witness-envelope envelope.json .sdp-trace-runs` | CI-bound evidence не является transparency log или release approval. |
| `witness --kind customer-pki` | Записывает customer PKI/private-equivalent witness evidence. | `go run ./cmd/sdp-trace witness --kind customer-pki --out customer-pki-witness.json --customer-pki-authority-policy policy.json --customer-pki-public-cert cert.pem --customer-pki-payload-digest <sha256> --customer-pki-freshness-evidence freshness.json .sdp-trace-runs` | Требует accepted customer policy, key/cert material, payload digest и freshness evidence. |
| `release-proof` | Проверяет source-bound local release manifest. | `go run ./cmd/sdp-trace release-proof --manifest examples/contract-foundation/contract-manifest.example.json --out release-proof.json` | Caveat: `source_bound_local_release` уже, чем external production trust; dirty/stale source или manifest mismatch не проходят. |
| `validate-fixtures` | Валидирует fixture directories. | `go run ./cmd/sdp-trace validate-fixtures examples` | Только structural fixture validation; не доказывает customer production readiness. |

## State And Exit Code Contract

- `0`: `observed`, `pass`, or explicitly scoped `not_assessed`
- `1`: `fail`
- `2`: usage error / invalid command invocation
- `3`: `cannot_verify`

- `observed`: verifier evidence met required checks for the selected local profile.
- `pass`: selected profile concluded successfully where the command contract uses pass/fail states.
- `fail`: verifier evidence conflicted or was insufficient for required checks.
- `not_assessed`: state was outside the run scope; it does not imply success or evidence.
- `cannot_verify`: verifier could not execute a required check or lacked required evidence.

A checked-in proof JSON is an audit artifact, not authority. Authority is replayed
only from live Go verifier output and the canonical command/state contract above.

## Air-Gapped Fixture Guidance

Air-gapped evidence is a fixture and customer-policy pattern, not a native
`witness --kind air-gapped` command. Use `customer-pki` or an accepted private
equivalent with explicit authority policy, payload digest, freshness evidence,
and retained audit references. If those are absent, record `not_assessed` or
`cannot_verify`; do not claim external production trust.

## Forbidden Claims

Do not emit these in this repo surface:

- `external_production_trust=true` without a live `external_production_trust` profile pass.
- `trusted_contract_release=true` without live external trust closure.
- `production_release_verified=true` outside a concluded `external_production_trust` run.
- Claims that treat `repo_baseline_structural` or `source_bound_local_release` outputs as production trust.
- Dirty-checkout structural output as source-bound or external-trust evidence.
