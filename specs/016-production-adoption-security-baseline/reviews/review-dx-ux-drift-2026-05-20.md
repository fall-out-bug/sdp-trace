# Повторное ревью: DX, UX, Drift

Дата: 2026-05-20
База: `origin/main` (`18cd4a0`)
Ветка: `feat/016-production-adoption-security-baseline` (`aa1e7e3`)

## Drift

### D1 — critical → verified fixed after rebase
** readiness.md отличается от command surface registry в main.**
- Проверено: после rebase `export-telemetry` присутствует (state `partial`), `override` — `not_assessed`.
- Итоговая строка корректна: "6 complete, 11 partial, 1 not_assessed (18 total)".
- Файлы: `docs/production-adoption-readiness.md`.

### D2 — minor → verified fixed after rebase
**Устаревший комментарий в readiness.md.**
- Проверено: stale-комментарий отсутствует; advisory scans указаны как `# Advisory scans` без ссылки на незавершённый WS.
- Файл: `docs/production-adoption-readiness.md`.

### D3 — info (не требует исправления)
**spec.md probe evidence содержит исторические данные о 10 gitleaks findings.**
- Это корректно как probe evidence от момента intake (до allowlist). Текущее состояние (0 findings с `.gitleaks.toml`) задокументировано в `security-baseline.md`.

## UX

### U1 — minor (advisory)
**Таблица Readiness Matrix имеет 6 колонок.**
- На узких экранах (Telegram, мобильный браузер) будет перенос строк.
- Это стандартное ограничение markdown; в source-формате читаемо.
- Действие: advisory, без изменений.

### U2 — positive
**Command Family Readiness таблица имеет 3 колонки — читаема.**
- Trust Note даёт контекст для каждой команды.

### U3 — positive
**Все относительные ссылки в docs/README.md валидны.**
- `production-adoption-readiness.md`, `security-baseline.md`, `../.github/SECURITY.md` — файлы существуют.

## DX

### DX1 — positive
**Verification commands copy-pasteable в `security-baseline.md` и `.github/SECURITY.md`.**

### DX2 — positive
**`.gitleaks.toml` находится в корне репозитория.**
- `gitleaks detect --source .` автоматически подхватывает конфиг.

### DX3 — positive
**Review-артефакты находятся в `specs/016-production-adoption-security-baseline/reviews/`.**
- Соответствует `hygienecheck` и `docs/README.md`.

## Codex Review Findings (2026-05-20) — All Fixed

| Finding | Severity | Fix | Commit |
|---------|----------|-----|--------|
| Security contact `security+report@fall-out-bug.dev` undeliverable (NXDOMAIN) | major | Removed email; fallback to GitHub private vulnerability reporting; marked dedicated email `not_assessed` | `aca55ca` |
| Command readiness undercounts (27 commands vs 18 families) | major | Added Utility Commands section (9 commands, all `complete`); fixed `export-telemetry` → `export` | `aca55ca` |
| Security scanner ledger inconsistent (12 vs 0 findings) | major | Split default-config (12 findings) vs configured-scan (0 findings); replaced stale line-number table with category-based triage | `aca55ca` |
| Spec workflow state drift (roadmap still `draft`) | major | Updated `docs/roadmap.md` spec 016 status → `in_progress` | `aca55ca` |
| Duplicate entries in `docs/README.md` | minor | Removed duplicate links | `aca55ca` |

## Резюме

- D1 — verified fixed (export family name, utility commands added).
- D2 — verified fixed (stale comment absent).
- All Codex review findings fixed in `aca55ca`.
- U1 — advisory.
- Остальное — positive или info.
