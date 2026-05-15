---
description: DeepSeek reasoning review
argument-hint: "<artifact-description> <contract>"
---
Adversarial review via DeepSeek. Find what is wrong with this artifact. Assume the author is overconfident.

Look for: unstated assumptions, missing evidence, schema drift, source/provenance gaps, incorrect gate semantics, stale checked-in proof, external trust overclaim, security/forgery risk, and project convention violations.

Do not validate. Do not summarize. Return only actionable issues with file/line or artifact references, or state that you cannot find any after checking the contract.

ARTIFACT:
$1

CONTRACT:
$2
