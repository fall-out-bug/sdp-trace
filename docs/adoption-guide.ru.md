# Руководство по внедрению sdp-trace

`sdp-trace` добавляет рядом с существующим AI-workflow слой записей: что было
сделано, какие доказательства есть, чего не хватает и кто принимает следующий
человеческий decision. Он не заменяет агента, CI, code review или release
governance.

Core pilot path: trace capture, local verification, explanation, report и
missing-evidence query. Assessment profiles, advisory/protected gate facts,
CI/customer witness profiles, forensic query packs, cross-repository posture
export, PR review packets и source-bound release proof являются extension
surfaces. См.
[`command-stability-matrix.md`](command-stability-matrix.md) для
current stability classification каждого command family. Это не automatic merge
blocking, не production release approval, не external audit proof и не гарантия
обнаружения каждого unwrapped agent run.

## Что дает sdp-trace

По каждому репозиторию и коммиту организация может проверить:

- какой agent или human workflow был observed;
- какой task, command, model, harness и source context записаны;
- какой evidence contract или assessment profile ожидался;
- какие artifacts retained, redacted или digest-only;
- какое evidence missing, `not_assessed` или `cannot_verify`;
- run только local, CI-witnessed, customer-PKI witnessed или пока только
  documentation/fixture guidance;
- прошел ли source-bound local release proof без claims про external production
  trust.

Opaque score нет. Missing telemetry остается видимой.

## Как это дополняет CI logs, git diff и review comments

CI logs показывают command output. Git diff показывает file changes. Review
comments показывают human discussion. `sdp-trace` добавляет portable evidence
contract поверх этих поверхностей:

- provenance связывает, кто или что произвел evidence;
- trace runs сохраняют command и task context;
- assessment profiles объясняют, почему evidence достаточно, отсутствует,
  устарело или unverifiable;
- witness records связывают выбранное evidence с CI или customer authority,
  когда профилю хватает данных;
- release proof проверяет manifest subjects against source commit вместо
  доверия к prose.

Это все еще evidence, а не policy decision. CI, release management, customer
governance или другой внешний policy consumer решает, что блокировать.

## Модель внедрения

Путь внедрения sidecar-first:

```text
existing harness / agent / prompt flow
        |
        v
sdp-trace wrap / adapter events
        |
        v
.sdp-trace-runs/
        |
        v
verify, explain, report, missing-evidence query
        |
        v
optional extension: assess, gate, witness, packet, release-proof, export
        |
        v
evidence package per repo and commit
```

Минимальная core command sequence:

```text
sdp-trace wrap --name <workflow-name> --output-dir .sdp-trace-runs/<run-id> -- <existing command...>
sdp-trace verify .sdp-trace-runs/<run-id>
sdp-trace explain .sdp-trace-runs/<run-id>
sdp-trace report --out .sdp-trace-report .sdp-trace-runs
sdp-trace query --query missing-evidence .sdp-trace-runs/<run-id>
```

Optional extension sequence, если downstream policy consumer требует эти факты:

```text
sdp-trace gate --out .sdp-trace-report/gate-result.json .sdp-trace-runs
sdp-trace witness --kind github-actions --out .sdp-trace-report/ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs
```

Если агент или разработчик не запускается через `sdp-trace wrap` или adapter,
`sdp-trace` не видит эту локальную работу напрямую. Обнаруживаемый сигнал
появляется на expected evidence boundary: required run artifacts, adapter
events, witness bindings или profile inputs отсутствуют. Результат verifier —
`not_assessed` или `cannot_verify`. Метки telemetry, такие как
`missing_telemetry`, описывают характер gap, но не являются result states.
См. `docs/agent-entrypoint.md` для canonical state contract.

## Текущие профили и границы

| Surface | Что поддерживает сейчас | Caveat |
| --- | --- | --- |
| Core run/report/query | Records, verifies, explains, reports и queries missing evidence для local run. | Local evidence only; это не CI witness и не production trust. |
| `adapter-capture` | Проверяет adapter event coverage и overclaim risk. | Missing adapter events не доказывают, что агент не использовался; они доказывают, что профилю не хватает evidence. |
| `managed-harness` | Проверяет managed policy, adapter registry, run и witness evidence. | Выдает verifier facts; block/allow решает external CI или policy. |
| `forensic-retention` | Проверяет, хватает ли retained/redacted evidence для reconstruction. | Digest-only или unresolved redaction могут блокировать forensic claims. |
| `gate` | Выдает advisory/protected gate facts и reasons. | Это не native merge, release, readiness, degradation, override или risk decision. |
| `witness` | Создает witness artifacts для `github-actions`, `gitlab-ci`, `buildkite` и `customer-pki`. | CI/customer witness не external production trust, пока внешний trust profile не прошел. |
| `release-proof` | Проверяет source-bound local release manifests against source commit. | `source_bound_local_release` не равен `external_production_trust`. |
| Air-gapped guidance | Использует customer policy/private-equivalent evidence patterns. | Нет `witness --kind air-gapped`; unsupported evidence остается `not_assessed` или `cannot_verify`. |

## Что смотреть в репозитории

- `.sdp-trace-report/summary.json`: run и report summary.
- `.sdp-trace-report/evidence-table.json`: observed evidence rows.
- `.sdp-trace-report/missing-telemetry.json`: required evidence not observed.
- `sdp-trace query --query missing-evidence <run-dir>`: core missing-evidence table.
- `.sdp-trace-report/gate-result.json`: optional advisory/protected gate facts и reasons.
- `.sdp-trace-report/ci-witness.json` или другой witness artifact: CI/customer binding state.
- `.sdp-trace-runs/<run-id>/`: raw run package с учетом retention/redaction policy.
- `query-pack` output: incident или forensic reconstruction package.
- `release-proof` output: source-bound local release state.
- Workflow docs: spec, plan, tasks, evidence, decisions и deferred gaps из
  SpecKit, gsd, Superpowers, Oh My OpenAgent, ticket tracker или кастомного
  planning flow команды.

## Как читать missing states

- `not_assessed`: state был вне scope текущего run. Нужно решить, приемлем ли
  этот scope, или потребовать follow-up profile.
- `cannot_verify`: verifier пытался проверить state, но не хватило evidence,
  environment или consistency. Для trust claims это fail-closed.
- `fail`: evidence противоречит выбранному profile.
- `observed` или `pass`: выбранная local/profile проверка завершилась, но
  только внутри указанного trust scope.

## Privacy и non-capture

Текущая pilot-поверхность не должна требовать committed raw customer source, private prompts,
credentials, provider tokens или raw logs. Предпочитайте digests, sanitized
excerpts, encrypted external refs и explicit redaction notes. Если команде нужен
raw capture для incident, это отдельный retention/redaction profile с human
owner.

Trace metadata включает provenance fields вроде invocation working directory.
Запускайте `wrap` из project directory, который трассируется, и редактируйте
или не передавайте trace packages наружу, если local paths чувствительны.

## Rollout inputs

Перед rollout зафиксируйте в репозитории:

- expected evidence contract или profile inputs;
- wrapper или adapter integration command;
- report и retention policy;
- witness profile, если он нужен;
- правило: `not_assessed` и `cannot_verify` не являются pass;
- retained report, witness и assessment artifacts для конкретного проверяемого
  run.

Для rollout по нескольким репозиториям требуйте, чтобы команды публиковали
report directory и witness artifacts как retained CI/customer-policy artifacts.
Отслеживайте долю `observed`, `pass`, `fail`, `not_assessed` и `cannot_verify`
по репозиториям. Dashboards и policy decisions должны жить вне `sdp-trace`.
