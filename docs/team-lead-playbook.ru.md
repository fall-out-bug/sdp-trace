# Playbook для тимлида

Используйте `sdp-trace`, когда у команды уже есть AI coding workflow и нужен общий quality contract.

## Ежедневный сценарий

1. Зафиксировать scope.
2. Записать provenance: человек, агент, модель, tools и команды.
3. Приложить evidence: тесты, CI, review comments, файлы и diff.
4. Записать accountability: human-held DRI, approver, risk owner и escalation.
5. Собрать assessment input с evidence, observations, movement data и `not_assessed` gaps.
6. Записать gate verdict только как external verdict input от `sdp-gate` или другого policy consumer.

## Team defaults

Договоритесь:

- какое evidence требуется для разных типов изменений
- какая внешняя policy блокирует merge
- кто может approve или override в policy layer
- какие harness поддерживаются
- что в вашей команде означает `not_assessed`

## Как читать external verdict

External verdict может использовать `pass`, `warn`, `fail` или `not_assessed`, но это не native decision `sdp-trace`.

`not_assessed` — не pass. Missing evidence должно оставаться видимым в assessment input.
