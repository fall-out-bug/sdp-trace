# sdp-trace для CTO за одну минуту

AI-assisted delivery ускоряет работу, но ломает управляемость: через неделю трудно понять, что именно было обещано, кто или что это сделал, какие доказательства есть, где данные отсутствуют и кто отвечает за следующий шаг.

`sdp-trace` нужен не для того, чтобы сказать "все хорошо" или "мы деградируем". Он нужен, чтобы зафиксировать проверяемую цепочку:

```text
идея -> спецификация -> задача -> изменение -> evidence -> provenance -> accountability -> движение метрик -> проверенный контракт
```

Для CTO это означает: процесс можно смотреть во времени. Не через ощущения и не через opaque score, а через prior/current/delta, coverage evidence и явные `not_assessed` gaps.

Для CEO это означает: "ответственный" не ИИ. У каждого значимого артефакта есть human-held DRI, approver, risk owner и escalation.

Для CIO это означает: контракт нельзя тихо упростить. Schemas, docs,
validation commands, fixtures и release-proof artifacts можно проверять через
явные contract manifests и authority scopes.

`sdp-trace` сознательно не принимает policy decisions: pass/fail, readiness, degradation, thresholds и overrides принадлежат `sdp-gate` или другому внешнему policy consumer.

Текущая поверхность продукта годится для контролируемых пилотов: local trace
capture, report packages, явные missing telemetry states, assessment profiles,
CI/customer witness artifacts при наличии evidence и source-bound local release
proof.

Это не broad production trust, не universal harness compatibility и не
автоматическое обнаружение каждого локального agent run вне wrapper или adapter.

Начинайте с `docs/README.md`, `docs/cto-adoption-guide.ru.md` и
`docs/customer-questions.ru.md`. Development specs и research notes полезны для
audit history, но это не onboarding path для покупателя.
