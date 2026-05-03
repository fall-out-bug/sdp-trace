#!/usr/bin/env node
import fs from 'node:fs';
import process from 'node:process';
import { spawnSync } from 'node:child_process';

const ledgerPath = process.argv[2] || 'specs/001-sdp-trace-time-series-evidence-substrate/blocks/07-minimum-trust-kernel-review-ledger.json';

function fail(message) {
  console.error(message);
  process.exitCode = 1;
}

function run(command, args) {
  const result = spawnSync(command, args, {
    encoding: 'utf8',
    stdio: ['ignore', 'pipe', 'pipe']
  });
  if (result.status !== 0) {
    const output = `${result.stderr || result.stdout}`.split(/\r?\n/).find(Boolean) || `exit ${result.status}`;
    fail(`${command} ${args.join(' ')} failed: ${output}`);
    return null;
  }
  return result.stdout;
}

const schemaResult = run('node', ['scripts/validate-json-schema.mjs', 'schema/review-ledger.schema.json', ledgerPath]);
if (schemaResult === null) process.exit(process.exitCode || 1);

const ledger = JSON.parse(fs.readFileSync(ledgerPath, 'utf8'));

const openCriticalOrMajor = ledger.findings.filter((finding) =>
  ['critical', 'major'].includes(finding.severity) && finding.closure !== 'closed' && finding.id !== 'B07-S7-F002'
);
if (openCriticalOrMajor.length > 0) {
  fail(`critical/major findings remain unclosed: ${openCriticalOrMajor.map((finding) => finding.id).join(', ')}`);
}

const externalBlocker = ledger.findings.find((finding) => finding.id === 'B07-S7-F002');
if (!externalBlocker || externalBlocker.closure !== 'blocking') {
  fail('external production trust blocker B07-S7-F002 must remain recorded as blocking');
}

const profiles = Object.fromEntries(ledger.profile_results.map((profile) => [profile.profile, profile.result]));
if (profiles.repo_baseline_structural !== 'pass') fail('repo_baseline_structural must be pass for Block 07 closure');
if (profiles.source_bound_local_release !== 'pass') fail('source_bound_local_release must be pass for Block 07 closure');
if (profiles.external_production_trust !== 'fail') fail('external_production_trust must be fail until external evidence exists');

if (ledger.activation_decision?.block_08_allowed !== true) {
  fail('Block 08 activation must be explicitly allowed by the closure ledger');
}

if (!process.exitCode) {
  console.log('review ledger valid');
}
