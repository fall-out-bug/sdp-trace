#!/usr/bin/env node
import fs from 'node:fs';
import path from 'node:path';
import process from 'node:process';
import { spawnSync } from 'node:child_process';

const ROOT = process.cwd();
const REQUIRED_KEYS = ['claim', 'subject', 'state', 'profile', 'evidence'];
const ALLOWED = {
  claim: new Set(['task_closed', 'command_verified', 'profile_passed', 'trust_not_assessed']),
  state: new Set(['pass', 'fail', 'not_assessed', 'stale', 'cannot_verify']),
  profile: new Set(['repo_baseline', 'source_bound_local_release', 'external_production_trust', 'observed_slice'])
};
const ALLOWED_EVIDENCE = new Set(['command_set:block04-t070', 'state:claim_tags_consistent']);

const COMMAND_SETS = {
  'block04-t070': [
    ['npm', ['run', 'validate']],
    ['git', ['diff', '--check']],
    ['scripts/verify-self-attestation.sh', ['--all']],
    ['scripts/finalize-source-bound-release.sh', [
      '--manifest',
      'examples/contract-foundation/contract-manifest.example.json',
      '--source-ref',
      'HEAD'
    ]]
  ]
};

function fail(message) {
  console.error(message);
  process.exitCode = 1;
}

function usage() {
  console.error('Usage: scripts/validate-claim-consistency.mjs [--files <path>...]');
  process.exit(2);
}

function walk(dir, files = []) {
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    if (['.git', '.beads', 'node_modules', 'benchmarks'].includes(entry.name)) continue;
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      walk(fullPath, files);
    } else if (entry.isFile() && fullPath.endsWith('.md')) {
      files.push(fullPath);
    }
  }
  return files;
}

function parseArgs(argv) {
  if (argv.length === 0) {
    return walk(path.join(ROOT, 'docs')).concat(walk(path.join(ROOT, 'specs')));
  }
  if (argv[0] !== '--files' || argv.length === 1) usage();
  return argv.slice(1);
}

function stripQuoted(value) {
  const trimmed = value.trim();
  if (trimmed.startsWith('"') && trimmed.endsWith('"')) {
    return trimmed.slice(1, -1).replace(/\\"/g, '"');
  }
  return trimmed;
}

function splitFields(body) {
  const fields = [];
  let current = '';
  let quoted = false;
  for (let i = 0; i < body.length; i += 1) {
    const char = body[i];
    if (char === '"' && body[i - 1] !== '\\') quoted = !quoted;
    if (char === ';' && !quoted) {
      fields.push(current.trim());
      current = '';
    } else {
      current += char;
    }
  }
  if (current.trim()) fields.push(current.trim());
  return fields;
}

function parseClaimTag(body, location) {
  const parsed = {};
  for (const field of splitFields(body)) {
    const eq = field.indexOf('=');
    if (eq === -1) {
      fail(`${location}: malformed claim field: ${field}`);
      continue;
    }
    const key = field.slice(0, eq).trim();
    const value = stripQuoted(field.slice(eq + 1));
    if (!REQUIRED_KEYS.includes(key)) {
      fail(`${location}: unsupported claim key: ${key}`);
      continue;
    }
    if (Object.hasOwn(parsed, key)) {
      fail(`${location}: duplicate claim key: ${key}`);
      continue;
    }
    parsed[key] = value;
  }
  for (const key of REQUIRED_KEYS) {
    if (!Object.hasOwn(parsed, key)) fail(`${location}: missing claim key: ${key}`);
  }
  return parsed;
}

function validateClaimValues(claim, location) {
  for (const [key, allowed] of Object.entries(ALLOWED)) {
    if (!allowed.has(claim[key])) {
      fail(`${location}: unsupported ${key} value: ${claim[key]}`);
    }
  }
  const evidence = claim.evidence;
  if (!ALLOWED_EVIDENCE.has(evidence)) {
    fail(`${location}: unsupported evidence value: ${evidence}`);
  }
}

function runCommandSet(name, location) {
  const commands = COMMAND_SETS[name];
  if (!commands) {
    fail(`${location}: unsupported command_set: ${name}`);
    return;
  }
  for (const [command, args] of commands) {
    const result = spawnSync(command, args, {
      cwd: ROOT,
      encoding: 'utf8',
      stdio: ['ignore', 'pipe', 'pipe']
    });
    if (result.status !== 0) {
      const firstLine = `${result.stderr || result.stdout}`.split(/\r?\n/).find(Boolean) || `exit ${result.status}`;
      fail(`${location}: command_set:${name} did not pass: ${command} ${args.join(' ')}: ${firstLine}`);
      return;
    }
  }
}

function validateClaimEvidence(claim, location) {
  if (claim.state === 'pass' && claim.evidence.startsWith('command_set:')) {
    runCommandSet(claim.evidence.slice('command_set:'.length), location);
  }
  if (claim.state === 'pass' && claim.evidence === 'state:claim_tags_consistent') {
    fail(`${location}: pass claims require executable evidence in Slice 1`);
  }
}

function scanFile(file) {
  const text = fs.readFileSync(file, 'utf8');
  const lines = text.split(/\r?\n/);
  let inFence = false;
  lines.forEach((line, index) => {
    if (/^\s*```/.test(line)) {
      inFence = !inFence;
      return;
    }
    if (inFence) return;
    const matches = line.matchAll(/<!--\s*sdp-trace-claim:\s*(.*?)\s*-->/g);
    for (const match of matches) {
      const location = `${file}:${index + 1}`;
      const claim = parseClaimTag(match[1], location);
      validateClaimValues(claim, location);
      validateClaimEvidence(claim, location);
    }
  });
}

for (const file of parseArgs(process.argv.slice(2))) {
  scanFile(file);
}

if (!process.exitCode) {
  console.log('claim tags valid');
}
