# Feature Specification: Contributor Onboarding And Verification

**Feature Branch**: `013-contributor-onboarding-and-verification`
**Created**: 2026-05-14
**Status**: Draft for PI review
**Input**: DX review found no single quick-start gate for a new contributor or coding agent after reading the first docs.

## Problem Statement

The repository has multiple strong entrypoints, but a new contributor must jump between README, install docs, agent onboarding, reviewer entrypoint, and agent entrypoint to know whether their local environment is ready. Smoke commands are repeated in several places and are not all drift-checked.

## Core Claim

This slice may claim:

> A new contributor or coding agent can follow one short onboarding path and produce a local verification result before deeper work.

This slice must not claim:

- CI-backed readiness;
- production trust;
- successful setup for every OS or shell without evidence.

## Required User Stories

### US-001 - One-Page Contributor Start (P0)

A first-time contributor can open one linked page and run the minimum local commands required to know whether the repo is usable.

**Independent Test**: A cold reader can find the page from README and complete the listed commands without reading the full command reference.

### US-002 - Canonical Smoke Path (P0)

Smoke commands are defined once and referenced from README, install docs, onboarding, and reviewer docs.

**Independent Test**: A doc check fails if referenced smoke command blocks drift from the canonical source.

### US-003 - Failure Routing (P1)

When smoke setup fails, contributors know whether to inspect Go/toolchain setup, CLI usage, evidence state, or repo policy.

**Independent Test**: The quick-start includes expected pass/fail/cannot_verify behavior and the next diagnostic command.

## Functional Requirements

- **FR-001**: Add or revise a contributor start page with a short command sequence.
- **FR-002**: Link the page from README and docs map.
- **FR-003**: Avoid duplicating long command tables.
- **FR-004**: Include claim authoring caveats needed before agents write task/proof prose.
- **FR-005**: Make smoke command examples machine-checkable or explicitly non-authoritative.

## Acceptance Criteria

- New contributor path is discoverable from README in one click.
- Smoke path includes `--help`, environment check, wrap, verify, and explain/read step if available.
- Duplicated smoke command blocks are removed or checked.
- `go run ./tools/doccheck` validates the onboarding references it owns.

## PI Review Prompt

Review whether this onboarding path gives a cold contributor enough to act without forcing them through the entire docs set. Focus on duplication, command drift, and whether failure modes are actionable.
