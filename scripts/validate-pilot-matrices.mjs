#!/usr/bin/env node
import fs from 'node:fs';

const allowedStates = new Set(['observed', 'not_assessed']);
const allowedReasons = new Set([
  'run_artifact_available',
  'no_run_artifact',
  'missing_export',
  'discovery_required',
  'design_fixture_only',
  'sanitization_pending',
  'unsafe_to_run'
]);

const requiredColumns = [
  'Target',
  'Scope',
  'Evidence state',
  'Reason code',
  'Artifact reference',
  'External verdict reference',
  'Gap reason',
  'Next required evidence'
];

const matrices = [
  {
    path: 'docs/harness-compatibility-matrix.md',
    requiredTargets: [
      'Superpowers-style harness pattern',
      'gsd',
      'gsd2',
      'Oh My OpenAgent'
    ]
  },
  {
    path: 'docs/model-compatibility.md',
    requiredTargets: [
      'OpenCode + MiniMax',
      'OpenCode + Kimi',
      'OpenCode + GLM'
    ]
  }
];

function fail(message) {
  console.error(message);
  process.exitCode = 1;
}

function parseRows(markdown, path) {
  const lines = markdown.split(/\r?\n/);
  const tableStart = lines.findIndex((line, index) => {
    if (!line.trim().startsWith('|')) return false;
    const next = lines[index + 1] || '';
    return next.includes('---');
  });
  if (tableStart === -1) {
    fail(`${path}: no markdown table found`);
    return [];
  }

  const headers = lines[tableStart]
    .split('|')
    .slice(1, -1)
    .map((cell) => cell.trim());

  for (const column of requiredColumns) {
    if (!headers.includes(column)) {
      fail(`${path}: missing required matrix column: ${column}`);
    }
  }

  const rows = [];
  for (let i = tableStart + 2; i < lines.length; i += 1) {
    const line = lines[i];
    if (!line.trim().startsWith('|')) break;
    const cells = line
      .split('|')
      .slice(1, -1)
      .map((cell) => cell.trim().replace(/^`|`$/g, ''));
    if (cells.length !== headers.length) {
      fail(`${path}:${i + 1}: table row has ${cells.length} cells, expected ${headers.length}`);
      continue;
    }
    rows.push(Object.fromEntries(headers.map((header, idx) => [header, cells[idx]])));
  }
  return rows;
}

function readJsonIfPossible(path) {
  if (!path.endsWith('.json')) return null;
  try {
    return JSON.parse(fs.readFileSync(path, 'utf8'));
  } catch {
    fail(`${path}: observed artifact reference is not parseable JSON`);
    return null;
  }
}

function artifactIsPlaceholder(value) {
  if (!value || typeof value !== 'object') return false;
  if (Array.isArray(value)) return value.some(artifactIsPlaceholder);
  if (value.evidence_state === 'not_assessed') return true;
  if (['design_fixture_only', 'no_run_artifact'].includes(value.reason_code)) return true;
  if (value.status === 'not_assessed' && ['design_fixture_only', 'no_run_artifact'].includes(value.reason_code)) return true;
  return Object.values(value).some(artifactIsPlaceholder);
}

for (const matrix of matrices) {
  const markdown = fs.readFileSync(matrix.path, 'utf8');
  const rows = parseRows(markdown, matrix.path);
  const targets = new Set(rows.map((row) => row.Target));

  for (const target of matrix.requiredTargets) {
    if (!targets.has(target)) {
      fail(`${matrix.path}: missing required target row: ${target}`);
    }
  }

  for (const row of rows) {
    if (!allowedStates.has(row['Evidence state'])) {
      fail(`${matrix.path}: ${row.Target}: invalid evidence state: ${row['Evidence state']}`);
    }
    if (!allowedReasons.has(row['Reason code'])) {
      fail(`${matrix.path}: ${row.Target}: invalid reason code: ${row['Reason code']}`);
    }
    if (!row['Gap reason'] || row['Gap reason'] === 'none') {
      fail(`${matrix.path}: ${row.Target}: gap reason must be explicit`);
    }
    if (!row['Next required evidence'] || row['Next required evidence'] === 'none') {
      fail(`${matrix.path}: ${row.Target}: next required evidence must be explicit`);
    }
    const artifact = row['Artifact reference'];
    if (row['Evidence state'] === 'observed') {
      if (!artifact || artifact === 'none') {
        fail(`${matrix.path}: ${row.Target}: observed rows require committed artifact reference`);
      } else if (!fs.existsSync(artifact)) {
        fail(`${matrix.path}: ${row.Target}: artifact reference does not exist: ${artifact}`);
      } else {
        const json = readJsonIfPossible(artifact);
        if (json && artifactIsPlaceholder(json)) {
          fail(`${matrix.path}: ${row.Target}: observed row cannot point to placeholder or not_assessed artifact: ${artifact}`);
        }
      }
    }
    const externalVerdict = row['External verdict reference'];
    if (externalVerdict && externalVerdict !== 'none' && !fs.existsSync(externalVerdict)) {
      fail(`${matrix.path}: ${row.Target}: external verdict reference does not exist: ${externalVerdict}`);
    }
  }
}

const claimBoundaryFiles = [
  'docs/research/run-card-template.md',
  'docs/research/opencode-model-run-card.md',
  'docs/research/harness-run-card.md',
  'docs/research/kotlin-bazel-fixture-plan.md',
  'docs/research/customer-pilot-evidence-package.md',
  'docs/harness-compatibility-matrix.md',
  'docs/model-compatibility.md',
  'docs/jvm-bazel-guide.md'
];

const boundaryToken = /\b(pass|fail|warn|blocked|readiness|ready|support|supported|supports|compatible|compatibility|TBD|external_verdict_recorded|evidence_backed)\b/i;
const boundaryAllowlist = [
  /does not/i,
  /do not/i,
  /not claim/i,
  /not evidence/i,
  /not a completed/i,
  /not native/i,
  /native .* outcome/i,
  /external/i,
  /legacy-named/i,
  /planned assessment/i,
  /auxiliary configuration/i,
  /prohibited/i,
  /banned/i,
  /policy reference/i,
  /raw prompts/i,
  /prompt.*approval/i,
  /approved for release/i,
  /outside the repository/i,
  /compatibility-matrix\.md/i,
  /model-compatibility\.md/i,
  /not_assessed/i
];

for (const path of claimBoundaryFiles) {
  const lines = fs.readFileSync(path, 'utf8').split(/\r?\n/);
  lines.forEach((line, index) => {
    if (!boundaryToken.test(line)) return;
    if (boundaryAllowlist.some((pattern) => pattern.test(line))) return;
    fail(`${path}:${index + 1}: possible native verdict/support wording without boundary: ${line.trim()}`);
  });
}

const bannedChecks = [
  {
    path: 'docs/research/run-card-template.md',
    patterns: [/pass \/ warn \/ fail/i, /Evidence\s+Verdict/i]
  },
  {
    path: 'examples/jvm-bazel/evidence-bundle.json',
    patterns: [/"status"\s*:\s*"(pass|warn|fail)"/i]
  },
  {
    path: 'docs/jvm-bazel-guide.md',
    patterns: [/first-class stack targets/i]
  }
];

for (const check of bannedChecks) {
  const text = fs.readFileSync(check.path, 'utf8');
  for (const pattern of check.patterns) {
    if (pattern.test(text)) {
      fail(`${check.path}: prohibited native verdict/support wording matched ${pattern}`);
    }
  }
}

const placeholderSelfTest = readJsonIfPossible('examples/jvm-bazel/evidence-bundle.json');
if (!artifactIsPlaceholder(placeholderSelfTest)) {
  fail('examples/jvm-bazel/evidence-bundle.json: placeholder self-test failed; observed rows could over-credit design fixtures');
}

if (!process.exitCode) {
  console.log('pilot matrices valid');
}
