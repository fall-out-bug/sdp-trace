# sdp-trace для CTO за одну минуту

AI-assisted delivery ускоряет работу, но ломает управляемость: через неделю трудно понять, что именно было обещано, кто или что это сделал, какие доказательства есть, где данные отсутствуют и кто отвечает за следующий шаг.

`sdp-trace` нужен не для того, чтобы сказать "все хорошо" или "мы деградируем". Он нужен, чтобы зафиксировать проверяемую цепочку:

```text
идея -> спецификация -> задача -> изменение -> evidence -> provenance -> accountability -> движение метрик -> проверенный контракт
```

Для CTO это означает: процесс можно смотреть во времени. Не через ощущения и не через opaque score, а через prior/current/delta, coverage evidence и явные `not_assessed` gaps.

Для CEO это означает: "ответственный" не ИИ. У каждого значимого артефакта есть human-held DRI, approver, risk owner и escalation.

Для CIO это означает: контракт нельзя тихо упростить. Schemas, docs, validation scripts и fixtures входят в contract manifest с digest-проверкой и release verification profile.

`sdp-trace` сознательно не принимает policy decisions: pass/fail, readiness, degradation, thresholds и overrides принадлежат `sdp-gate` или другому внешнему policy consumer.

Block 01 строит contract scaffold: evidence contracts, accountability, manifest verification, signing profile, negative fixtures и proof that missing data stays `not_assessed`.

Продукт не имеет права просить доверия у заказчика, пока не отследит сам себя. Следующее доказательство должно показать собственные spec, plan, tasks, changes, evidence, provenance, accountability, reviews, metrics и missing data этого репозитория по тем же контрактам.

Репозиторные доказательства начинаются в `specs/001-sdp-trace-time-series-evidence-substrate/`, `schema/README.md`, `examples/contract-foundation/` и `docs/research/block-01-validation-summary.md`.
