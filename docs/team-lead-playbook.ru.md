# Playbook для тимлида

Используйте `sdp-trace`, когда у команды уже есть AI coding workflow и нужен общий quality contract.

## Ежедневный сценарий

1. Зафиксировать scope.
2. Записать provenance: человек, агент, модель, tools и команды.
3. Приложить evidence: тесты, CI, review comments, файлы и diff.
4. Запустить или записать gate checks.
5. Опубликовать verdict: `pass`, `warn`, `fail` или `not_assessed`.
6. Записать decision и причину override, если он был.

## Team defaults

Договоритесь:

- какое evidence требуется для разных типов изменений
- что блокирует merge
- кто может делать override
- какие harness поддерживаются
- что в вашей команде означает `not_assessed`

## Как читать verdict

- `pass`: evidence удовлетворяет gate.
- `warn`: evidence неполное или есть риск, но команда может продолжить.
- `fail`: gate criteria не выполнены.
- `not_assessed`: gate не смог принять доказуемое решение.

`not_assessed` — не pass.
