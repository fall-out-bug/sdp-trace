#!/usr/bin/env node
import fs from 'node:fs';
import path from 'node:path';
import process from 'node:process';

let selfTraceDir = 'examples/self-trace';
let e2ePackageDir = 'examples/pilot-runs/opencode-minimax-kotlin-bazel';

function usage() {
  console.error('Usage: scripts/validate-cross-references.mjs [--self-trace-dir path] [--e2e-package-dir path]');
  process.exit(2);
}

for (let i = 2; i < process.argv.length; i += 1) {
  const arg = process.argv[i];
  if (arg === '--self-trace-dir') {
    if (i + 1 >= process.argv.length) usage();
    selfTraceDir = process.argv[i + 1];
    i += 1;
  } else if (arg === '--e2e-package-dir') {
    if (i + 1 >= process.argv.length) usage();
    e2ePackageDir = process.argv[i + 1];
    i += 1;
  } else {
    usage();
  }
}

const failures = [];

function readJson(file) {
  return JSON.parse(fs.readFileSync(file, 'utf8'));
}

function requireRef(refs, validIds, location, kind) {
  for (const ref of refs || []) {
    if (!validIds.has(ref)) {
      failures.push(`${location}: missing ${kind} ref: ${ref}`);
    }
  }
}

function validateObservationRefs(observations, evidenceIds, provenanceIds, label) {
  for (const observation of observations) {
    if ((observation.assessment_status || observation.assessment_state) === 'not_assessed') continue;
    requireRef(observation.evidence_refs, evidenceIds, `${label}:${observation.id}`, 'evidence');
    requireRef(observation.provenance_refs, provenanceIds, `${label}:${observation.id}`, 'provenance');
  }
}

function validateMetricRefs(streams, evidenceIds, provenanceIds, label) {
  for (const stream of streams) {
    for (const sample of stream.samples || []) {
      if (sample.assessment_state === 'not_assessed') continue;
      requireRef(sample.evidence_refs, evidenceIds, `${label}:${stream.id}/${sample.id}`, 'evidence');
      requireRef(sample.provenance_refs, provenanceIds, `${label}:${stream.id}/${sample.id}`, 'provenance');
    }
  }
}

function validateTraceEdges(trace, label) {
  const nodeIds = new Set((trace.nodes || []).map((node) => node.id));
  for (const edge of trace.edges || []) {
    if (!nodeIds.has(edge.from)) failures.push(`${label}: edge from missing node: ${edge.from}`);
    if (!nodeIds.has(edge.to)) failures.push(`${label}: edge to missing node: ${edge.to}`);
  }
}

function validateSelfTrace(dir) {
  const evidence = readJson(path.join(dir, 'evidence-events.json'));
  const provenance = readJson(path.join(dir, 'provenance-records.json'));
  const observations = readJson(path.join(dir, 'observations.json'));
  const metrics = readJson(path.join(dir, 'metric-stream.json'));
  const trace = readJson(path.join(dir, 'trace-snapshot.json'));

  const evidenceIds = new Set(evidence.map((event) => event.id));
  const provenanceIds = new Set(provenance.map((record) => record.id));

  validateObservationRefs(observations, evidenceIds, provenanceIds, 'self-trace observation');
  validateMetricRefs(metrics, evidenceIds, provenanceIds, 'self-trace metric');
  validateTraceEdges(trace, 'self-trace trace-snapshot');
}

function validateE2ePackage(dir) {
  const evidence = readJson(path.join(dir, 'evidence/evidence-events.json'));
  const provenance = readJson(path.join(dir, 'evidence/provenance-records.json'));
  const observations = readJson(path.join(dir, 'evidence/observations.json'));
  const metrics = readJson(path.join(dir, 'evidence/metric-stream.json'));
  const trace = readJson(path.join(dir, 'evidence/trace-snapshot.json'));
  const proofStates = readJson(path.join(dir, 'evidence/proof-states.json'));

  const evidenceIds = new Set(evidence.map((event) => event.id));
  const successEvidenceIds = new Set(evidence.filter((event) => event.status === 'success').map((event) => event.id));
  const provenanceIds = new Set(provenance.map((record) => record.id));

  validateObservationRefs(observations, evidenceIds, provenanceIds, 'e2e observation');
  validateMetricRefs(metrics, evidenceIds, provenanceIds, 'e2e metric');
  validateTraceEdges(trace, 'e2e trace-snapshot');

  for (const state of proofStates.states || []) {
    if (state.state !== 'observed') continue;
    for (const ref of state.evidence_refs || []) {
      if (!successEvidenceIds.has(ref)) {
        failures.push(`e2e proof-state:${state.name}: missing successful evidence ref: ${ref}`);
      }
    }
  }
}

validateSelfTrace(selfTraceDir);
validateE2ePackage(e2ePackageDir);

if (failures.length > 0) {
  for (const failure of failures) console.error(failure);
  process.exit(1);
}

console.log('cross references valid');
