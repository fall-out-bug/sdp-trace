# Руководство CTO по внедрению sdp-trace

`sdp-trace` — это контрольный слой поверх существующей AI-assisted delivery. Он
не заменяет ваш harness, prompts, agents, CI, review process или repo
templates.

В текущей реализации Block 12 "контрольный слой" означает capture evidence,
явную missing telemetry и CI-witnessed reporting. Это еще не automatic merge
blocking, не org-wide degradation dashboard, не внешнее нотариальное
доказательство и не гарантированное обнаружение каждого agent run вне wrapper.

## Что получает CTO

По каждому репозиторию и коммиту организация видит:

- какой agent/human workflow был наблюдаем;
- какой evidence contract ожидался;
- какие команды и test/build шаги были наблюдаемы;
- какие артефакты сохранены как digest;
- какое evidence отсутствует;
- trace только local, CI-witnessed или externally witnessed;
- local contract gate прошел или нет, есть ли CI witness, и остается ли
  audit-grade evidence в состоянии `cannot_verify`.

Opaque score нет. Missing telemetry остается видимым.

## Модель внедрения

Путь внедрения sidecar-first:

```text
существующий harness / agent / prompt flow
        |
        v
sdp-trace wrap
        |
        v
local run artifacts
        |
        v
report + local gate
        |
        v
CI witness
        |
        v
evidence package для CTO/team по repo и commit
```

Минимальный integration contract:

- положить expected-evidence contract в каждый repo;
- завернуть существующую harness-команду;
- генерировать report/gate artifacts в CI;
- добавить CI witness record;
- явно показывать local-only и missing telemetry states.

Если агент или разработчик не запускается через `sdp-trace wrap` или adapter,
`sdp-trace` не видит эту локальную работу напрямую. Первый обнаруживаемый
сигнал появляется на expected evidence boundary: CI/report/gate покажет, что
required run artifacts или contract evidence отсутствуют. Это полезный сигнал
control posture, но не полный agent activity log.

Минимальная CI-последовательность:

```text
sdp-trace report --out .sdp-trace-report --contract <contract> .sdp-trace-runs
sdp-trace gate --out .sdp-trace-report/gate-result.json --contract <contract> .sdp-trace-runs
sdp-trace witness --kind github-actions --out .sdp-trace-report/ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs
sdp-trace gate --out .sdp-trace-report/gate-result.json --contract <contract> --witness .sdp-trace-report/ci-witness.json .sdp-trace-runs
```

## Уровни доверия

- `local_observed`: полезно для расследования, но не gate-grade.
- `ci_witnessed`: CI связал evidence package с repo, commit, workflow, job и
  run id.
- `external_witnessed`: будущий профиль независимого timestamp/log witness.

CI witness нужен покупателю, потому что локальная история перестает выглядеть
как gate-grade trace. Но сам по себе CI witness не доказывает честность агента
или качество релиза.

Для GitHub Actions workflow должен разрешить OIDC (`id-token: write`). Без OIDC
`sdp-trace witness` записывает `cannot_verify`, а не делает вид, что обычных
environment variables достаточно.

Block 12 CI witness — это не external trust. Это CI-generated JSON artifact,
который связывает report/run digests с GitHub Actions OIDC claims, если он
создан в protected workflow. Это не public transparency log, не DSSE envelope и
не court-ready signed timeline. Нельзя считать witness-файл, закоммиченный
агентом или разработчиком, authority; генерируйте его внутри CI и храните как
protected CI artifact.

Policy interpretation:

| State | Что поддерживает | Чего не поддерживает |
|---|---|---|
| `local_observed` | Local reconstruction и feedback разработчику | Merge/release trust сам по себе |
| `ci_witnessed` | CI-bound evidence package для repo/commit | Agent honesty, достаточность тестов, external audit proof |
| `external_witnessed` | Будущий external timestamp/log profile | Не реализовано в Block 12 |

## Что пока не видно

Block 12 пока не дает:

- автоматического обнаружения, что агент использовался вне wrapper;
- internal tool-call telemetry для harness без adapter;
- raw prompt/model response capture;
- file mutation или VCS event capture сверх текущих recorder events;
- signed timeline или append-only transparency log;
- automatic degradation analytics по всем репозиториям;
- dashboard/query surface кроме generated artifacts и существующих local query
  commands.

Это явные gaps, а не скрытый pass.

## Что передать инженерному лиду

Передайте лиду:

- шаблон expected-evidence contract;
- wrapper command;
- CI witness command;
- policy для report directory;
- правило: `cannot_verify` не является pass.

CTO должен смотреть `.sdp-trace-report/` по repo и commit, а не raw JSON на
машинах разработчиков.

Для rollout по нескольким репозиториям требуйте, чтобы команды публиковали
CI-generated report directory как retained CI artifact. Начните с отслеживания
доли `local_observed`, `ci_witnessed`, `cannot_verify` и missing evidence по
репозиториям. Это первый честный degradation signal; более широкий
`sdp-report` analytics layer — следующий продуктовый слой.
