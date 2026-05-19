---
description: Kimi wide-context code and spec review
argument-hint: "<artifact-description> <contract>"
---
Adversarial review via Kimi. Find what is wrong with this artifact. Assume the author is overconfident.

Look for: missing edge cases, incomplete negative tests, stale evidence, source/provenance gaps, incorrect gate semantics, over-broad refactors, unchecked external trust, schema drift, and project convention violations.

Do not validate. Do not summarize. Return only actionable issues with file/line or artifact references, or state that you cannot find any after checking the contract.

ARTIFACT:
$1

CONTRACT:
$2
