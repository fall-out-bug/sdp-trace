# Playbook для тимлида

Используйте `sdp-trace`, когда у команды уже есть AI coding workflow и нужен общий quality contract.

Этот playbook описывает текущую поверхность Block 12: process wrapping,
local report/gate artifacts и GitHub Actions OIDC witness. Он еще не обещает
полные harness internals, automatic file mutation tracing, fail-closed managed
harness enforcement или external signed timelines.

## Ежедневный сценарий

1. Зафиксировать scope.
2. Записать provenance: человек, агент, модель, tools и команды.
3. Приложить evidence: тесты, CI, review comments, файлы и diff.
4. Записать accountability: human-held DRI, approver, risk owner и escalation.
5. Собрать assessment input с evidence, observations, movement data и `not_assessed` gaps.
6. Записать gate verdict только как external verdict input от `sdp-gate` или другого policy consumer.

Текущая capture boundary:

- `wrap` наблюдает lifecycle wrapped process и command-level events.
- Он не видит automatically internal tool calls внутри harness, если harness не
  отправляет adapter events.
- Он не доказывает, что никто не запускал агента вне wrapper.
- Missing expected evidence должно оставаться `missing_telemetry` или
  `cannot_verify`, а не превращаться в pass.

## Team defaults

Договоритесь:

- какое evidence требуется для разных типов изменений
- какая внешняя policy блокирует merge
- кто может approve или override в policy layer
- какие harness поддерживаются
- что в вашей команде означает `not_assessed`

## Настройка репозитория

Для каждого репозитория добавьте:

- expected-evidence contract, которым владеет команда;
- `.sdp-trace-runs/` для wrapped local/CI runs;
- `.sdp-trace-report/` для CI artifacts;
- CI steps для `report`, `gate`, `witness` и `gate --witness`.

Последовательность внедрения:

```text
sdp-trace wrap --name <workflow-name> --contract <contract> --output-dir .sdp-trace-runs/<run-id> -- <existing command...>
sdp-trace report --out .sdp-trace-report --contract <contract> .sdp-trace-runs
sdp-trace gate --out .sdp-trace-report/gate-result.json --contract <contract> .sdp-trace-runs
sdp-trace witness --kind github-actions --out .sdp-trace-report/ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs
sdp-trace gate --out .sdp-trace-report/gate-result.json --contract <contract> --witness .sdp-trace-report/ci-witness.json .sdp-trace-runs
```

Для GitHub Actions включите OIDC:

```text
permissions:
  id-token: write
  contents: read
```

Без OIDC `ci_witness_gate` остается `cannot_verify`.

Не коммитьте `.sdp-trace-report/ci-witness.json` с машины разработчика как
trusted evidence. Генерируйте его в CI и храните как CI artifact.

## Privacy и redaction

По умолчанию Block 12 report artifacts сохраняют command metadata и
stdout/stderr digests, а не raw stdout/stderr bodies. OIDC request tokens
используются только для запроса GitHub OIDC token и не должны записываться в
`.sdp-trace-runs/` или `.sdp-trace-report/`.

Перед rollout договоритесь:

- можно ли вообще сохранять prompts, source snippets или tool payloads;
- какие outputs должны оставаться digest-only;
- кто может разрешить raw capture для узкого incident window;
- как redaction decisions фиксируются для расследования.

Если команде нужен raw prompt, source или model-response capture, это отдельный
adapter/redaction profile. Не считайте это включенным этим playbook.
Любой будущий raw-capture profile должен делать redaction до persistent write.
Block 12 не дает raw-capture mode и поэтому не обещает безопасность post-hoc
redaction.

## Emergency changes

Не прячьте emergency bypass. Если production urgency требует shipping с missing
telemetry, запишите это как policy override во внешнем policy layer и оставьте
`cannot_verify` или `missing_telemetry` видимым в report. Bypass допустим
только если организация позже увидит, кто его approved, почему и какое evidence
отсутствовало.

В Block 12 еще нет native `policy_override_requested` trace event. Пока он не
реализован, override record должен жить во внешней policy/change management
системе и ссылаться на report artifacts.

## Offline и failure modes

- `wrap`, `report` и local contract gate могут работать без network access.
- `witness --kind github-actions` требует GitHub Actions OIDC и поэтому не может
  пройти offline.
- Если `witness` завершился exit `3`, смотрите `ci-witness.json`: `reason` и
  `missing_identity_fields`.
- Если `gate` fail, смотрите `gate-result.json`, `missing-telemetry.json` и
  per-run `verifier/` outputs.
- Используйте `sdp-trace dry-run --contract <contract> -- <command...>`, чтобы
  preview command и contract без записи run artifacts.

Gate debugging checklist:

1. Проверьте `gate-result.json`: `local_gate`, `ci_witness_gate`,
   `audit_grade_gate`.
2. Сравните `required_evidence` и `observed_evidence`.
3. Проверьте `missing-telemetry.json` на absent contract evidence.
4. Проверьте `ci-witness.json`: `reason`, `missing_identity_fields`, OIDC state.
5. Проверьте `verifier/` output каждого run перед изменением contract.

## Как читать external verdict

External verdict может использовать `pass`, `warn`, `fail` или `not_assessed`, но это не native decision `sdp-trace`.

`not_assessed` — не pass. Missing evidence должно оставаться видимым в assessment input.

## Что проверять

Проверяйте эти артефакты по repo и commit:

- `.sdp-trace-report/summary.json`
- `.sdp-trace-report/evidence-table.json`
- `.sdp-trace-report/missing-telemetry.json`
- `.sdp-trace-report/gate-result.json`
- `.sdp-trace-report/ci-witness.json`

Главные вопросы:

- Expected evidence contract соответствует работе?
- Все required evidence ids observed?
- Trace только `local_observed` или уже `ci_witnessed`?
- Missing telemetry не спрятано?
- `audit_grade_gate` все еще `cannot_verify`, потому что external witness еще не
  реализован?

## Retention

Для расследований храните `.sdp-trace-report/` и соответствующий
`.sdp-trace-runs/` не меньше, чем CI logs и review records. Если у организации
есть incident/audit retention requirements, сохраняйте report directory в
immutable artifact store. Block 12 сам retention не реализует.

## Текущие продуктовые gaps

- fail-closed managed harness enforcement;
- signed timeline / DSSE / external transparency witness;
- cross-repository dashboard и degradation analytics;
- richer query surface для monthly или incident-wide investigation;
- redaction audit trail сверх digest-only defaults;
- support для non-GitHub CI witness profiles.
- native `policy_override_requested` trace event;
- measured wrapper overhead и latency budgets.
