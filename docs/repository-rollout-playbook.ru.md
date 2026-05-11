# Playbook rollout для репозитория

Используйте `sdp-trace`, когда у команды уже есть AI coding workflow и нужен
общий evidence contract без замены текущего harness.

Давайте coding agents ссылку на [Agent Onboarding](agent-onboarding.md) как
первый вход в репозиторий. Этот playbook нужен для подключения конкретного
репозитория к evidence path.

Этот playbook описывает текущую pilot-поверхность: wrapping, task-linked runs,
previews, local verification, reports, queries, query packs, assessment
profiles, advisory/protected gate facts, CI/customer witness profiles,
source-bound release proof и fixture validation. Он не обещает full harness
internals, автоматическое обнаружение каждого bypass, external production trust
или policy decisions.

## Ежедневный сценарий

1. Зафиксировать spec, plan, task и expected evidence.
2. Записать provenance: human, agent, model, tools, commands и source context.
3. Запустить workflow через `wrap`, `run` или adapter, если workflow можно наблюдать.
4. Приложить evidence: tests, CI, review comments, files, diffs и retained artifacts.
5. Записать accountability: human-held DRI, approver, risk owner и escalation path.
6. Собрать report, query, assessment, gate facts и witness artifacts, где они доступны.
7. Оставить `not_assessed` и `cannot_verify` states видимыми.

Текущая capture boundary:

- `wrap` наблюдает lifecycle wrapped process и command-level events.
- Adapter profiles могут оценивать более богатые harness events, если harness
  их emits.
- `sdp-trace` не доказывает, что никто не запускал агента вне wrapper.
- Missing expected evidence должно оставаться `missing_telemetry`,
  `not_assessed` или `cannot_verify`, а не превращаться в pass.

## Team defaults

Договоритесь:

- какое evidence требуется для разных типов изменений;
- какие assessment profiles применяются к каким workflow;
- какая внешняя policy блокирует merge или release;
- кто может approve или override в policy layer;
- какие harness и CI systems поддерживаются;
- что в customer handoff означает `not_assessed`.

## Настройка репозитория

Для каждого репозитория добавьте:

- expected evidence contract, которым владеет команда;
- `.sdp-trace-runs/` для wrapped local/CI runs;
- `.sdp-trace-report/` для report artifacts;
- optional adapter registry, managed policy, redaction policy и witness policy files;
- CI steps для `report`, `gate` и выбранного `witness` kind.

Минимальная implementation sequence:

```text
sdp-trace wrap --name <workflow-name> --output-dir .sdp-trace-runs/<run-id> -- <existing command...>
sdp-trace report --out .sdp-trace-report .sdp-trace-runs
sdp-trace gate --out .sdp-trace-report/gate-result.json .sdp-trace-runs
sdp-trace witness --kind github-actions --out .sdp-trace-report/ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs
```

Полезные local checks:

```text
sdp-trace doctor
sdp-trace verify .sdp-trace-runs/<run-id>
sdp-trace explain .sdp-trace-runs/<run-id>
sdp-trace query --query missing-evidence .sdp-trace-runs/<run-id>
sdp-trace query-pack --pack forensics-basic-v1 --run .sdp-trace-runs/<run-id> --out query-pack.json
```

## Какие профили использовать

| Need | Command |
| --- | --- |
| Adapter coverage и overclaim review | `sdp-trace assess --profile adapter-capture --out assessment.json --run .sdp-trace-runs/<run-id>` |
| Managed harness profile | `sdp-trace assess --profile managed-harness --out assessment.json --contract contract.json --run .sdp-trace-runs/<run-id> --adapter-registry registry.json --managed-policy policy.json --managed-witness witness.json` |
| Forensic retention profile | `sdp-trace assess --profile forensic-retention --out assessment.json --run .sdp-trace-runs/<run-id> --redaction-policy redaction.json` |
| Assessment explanation | `sdp-trace assess explain --assessment-result assessment.json` |
| Source-bound release proof | `sdp-trace release-proof --manifest contract-manifest.json --out release-proof.json` |

`managed-harness` выдает verifier facts и exit behavior. Он не решает
merge/release readiness. Missing managed witness evidence обычно остается
`cannot_verify`.

## Witness profiles

Поддерживаемые значения `witness --kind`:

- `github-actions`
- `gitlab-ci`
- `buildkite`
- `customer-pki`

Для GitHub Actions включите OIDC:

```text
permissions:
  id-token: write
  contents: read
```

Без required identity или binding evidence witness output должен оставаться
`cannot_verify`. Не коммитьте witness file с машины разработчика как trusted
evidence. Генерируйте его в CI или в customer-approved PKI process и храните как
protected artifact.

Air-gapped evidence не отдельный command kind. Считайте его customer
policy/private-equivalent guidance: explicit authority policy, payload digest,
freshness или timestamp evidence и retained audit references. Если какой-то
обязательной части нет, записывайте `not_assessed` или `cannot_verify`.

## Gate debugging

`gate` output - это verifier-derived evidence, а не native policy decision.

Debugging checklist:

1. Проверьте `gate-result.json`: selected mode, required runs, required evidence и reason rows.
2. Проверьте `.sdp-trace-report/missing-telemetry.json` на absent contract evidence.
3. Проверьте witness output: source, run, freshness и identity binding state.
4. Проверьте `assess explain` output для profile-specific conditions.
5. Проверьте verifier output каждого run перед изменением contract.

## Privacy и redaction

По умолчанию используйте safe-to-commit artifacts: metadata, digests, sanitized
excerpts и external references. Не коммитьте raw customer source, private
prompts, credentials, provider tokens, authenticated URLs, OIDC request tokens
или raw logs. Если raw capture нужен для incident, сначала определите
retention/redaction profile и human owner.

## Emergency changes

Не прячьте emergency bypass. Если production urgency требует shipping с missing
telemetry, запишите override во внешней policy/change management system и
оставьте `cannot_verify` или `missing_telemetry` видимым в report. `sdp-trace`
записывает evidence; он не approve risk.

## Что проверять

Проверяйте эти artifacts по repo и commit:

- `.sdp-trace-report/summary.json`
- `.sdp-trace-report/evidence-table.json`
- `.sdp-trace-report/missing-telemetry.json`
- `.sdp-trace-report/gate-result.json`
- selected witness artifact, например `.sdp-trace-report/ci-witness.json`
- assessment result files для `adapter-capture`, `managed-harness` или `forensic-retention`
- query-pack output для incident или forensic review
- release-proof output для source-bound local release claims

Главные вопросы:

- Expected evidence contract соответствует работе?
- Все required evidence ids observed?
- Trace только local или уже bound to CI/customer evidence?
- Missing telemetry не спрятано?
- Какие states `not_assessed` или `cannot_verify`, и кто owns follow-up?
- Нет ли абзаца, который превращает local gate, witness или release proof в external production trust?

## Retention

Для расследований храните `.sdp-trace-report/`, selected witness artifacts и
соответствующий `.sdp-trace-runs/` не меньше, чем CI logs и review records. Если
есть incident/audit retention requirements, храните report directory в immutable
artifact store.

## Текущие продуктовые gaps

- native dashboards и policy decisions;
- guaranteed detection of every unwrapped agent run;
- external production trust без selected passing external profile;
- universal air-gapped witness command;
- raw prompt/source/model-response capture без отдельного redaction profile;
- measured wrapper overhead и latency budgets.
