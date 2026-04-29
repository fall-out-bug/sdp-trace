# sdp-trace: кратко для CTO

AI coding легко начать и трудно контролировать.

Риск не в том, что агент пишет код. Риск в том, что потом никто не может объяснить scope, происхождение, evidence и quality decision по изменению.

`sdp-trace` — переносимый trust layer для AI-assisted delivery. Он позволяет оставить текущий harness, но сделать изменения трассируемыми, доказуемыми и проверяемыми через gates.

## Что контролирует

- Scope: какое изменение было задумано.
- Provenance: кто или что выполняло работу.
- Evidence: тесты, CI, review, команды, diff и файлы.
- Gate verdict: `pass`, `warn`, `fail` или `not_assessed`.
- Decision record: почему изменение приняли, заблокировали или пропустили через override.

## Чего не обещает

- Не заменяет code review.
- Не гарантирует compliance.
- Не доказывает отсутствие багов.
- Не требует заменять Claude Code, Codex, OpenCode, Cursor или внутренний harness.

## Первый шаг внедрения

Возьмите один repo и одно реальное AI-assisted изменение. Соберите evidence bundle и gate verdict. Если verdict помогает reviewer принять решение быстрее и увереннее, расширяйте сценарий через `sdp-gate`.
