#!/usr/bin/env node
import fs from 'node:fs';
import path from 'node:path';
import crypto from 'node:crypto';
import process from 'node:process';
import Ajv2020 from 'ajv/dist/2020.js';

const ZERO_HASH = '0'.repeat(64);
const PROFILE_ALIASES = {
  local: 'local',
  flight_recorder_local: 'local',
  witnessed: 'witnessed',
  flight_recorder_witnessed: 'witnessed',
  forensic: 'forensic',
  flight_recorder_forensic: 'forensic'
};
const RUN_PROFILE_REQUIREMENT = {
  local_development_recorder: 'local',
  witnessed_run_recorder: 'witnessed',
  externally_witnessed_run: 'forensic'
};
const PROFILE_RANK = {
  local: 0,
  witnessed: 1,
  forensic: 2
};

const usage = () => {
  console.error('Usage: verify-flight-recorder.mjs --profile <local|witnessed|forensic|flight_recorder_local|flight_recorder_witnessed|flight_recorder_forensic> <fixture-dir>');
  console.error('       verify-flight-recorder.mjs --event-file <event-json>');
  process.exit(2);
};

const args = process.argv.slice(2);
let profile = null;
let eventFile = null;
let target = null;
for (let i = 0; i < args.length; i += 1) {
  if (args[i] === '--profile') {
    profile = args[i + 1];
    i += 1;
  } else if (args[i] === '--event-file') {
    eventFile = args[i + 1];
    i += 1;
  } else if (!target) {
    target = args[i];
  } else {
    usage();
  }
}

if ((!eventFile && (!profile || !target)) || (eventFile && (profile || target))) {
  usage();
}

const readJson = (file) => JSON.parse(fs.readFileSync(file, 'utf8'));
const sha256 = (value) => crypto.createHash('sha256').update(value).digest('hex');
const canonical = (value) => {
  if (Array.isArray(value)) {
    return `[${value.map(canonical).join(',')}]`;
  }
  if (value && typeof value === 'object') {
    return `{${Object.keys(value)
      .sort()
      .map((key) => `${JSON.stringify(key)}:${canonical(value[key])}`)
      .join(',')}}`;
  }
  return JSON.stringify(value);
};

const ajv = new Ajv2020({
  allErrors: true,
  strict: false,
  validateSchema: true
});

for (const file of fs.readdirSync('schema').filter((name) => name.endsWith('.json')).sort()) {
  const schema = readJson(path.join('schema', file));
  if (schema.$id) {
    ajv.addSchema(schema);
  }
}

const eventSchema = ajv.getSchema('https://schemas.sdp.dev/trace/flight-recorder-event.schema.json');
const runSchema = ajv.getSchema('https://schemas.sdp.dev/trace/flight-recorder-run.schema.json');
const witnessSchema = ajv.getSchema('https://schemas.sdp.dev/trace/flight-recorder-witness.schema.json');

const finding = (code, message, file, details = undefined) => ({
  code,
  severity: 'fail',
  message,
  ...(file ? { file } : {}),
  ...(details ? { details } : {})
});

const state = (stateValue, reason = undefined) => ({
  state: stateValue,
  ...(reason ? { reason } : {})
});

const validateWithSchema = (validate, data, file, findings) => {
  if (validate(data)) {
    return true;
  }
  findings.push(finding('schema_invalid', 'JSON document does not match its flight-recorder schema.', file, validate.errors));
  return false;
};

const recomputeEvent = (event) => {
  const payloadDigest = sha256(canonical(event.event_payload));
  const hashSubject = { ...event };
  delete hashSubject.event_hash;
  const eventHash = sha256(canonical(hashSubject));
  return { payloadDigest, eventHash };
};

const verifyEventObject = (event, file) => {
  const findings = [];
  validateWithSchema(eventSchema, event, file, findings);

  if (event.event_payload) {
    const { payloadDigest, eventHash } = recomputeEvent(event);
    if (event.event_payload_digest !== payloadDigest) {
      findings.push(finding('payload_digest_mismatch', 'event_payload_digest does not match canonical event_payload.', file, {
        expected: payloadDigest,
        actual: event.event_payload_digest
      }));
    }
    if (event.event_hash !== eventHash) {
      findings.push(finding('event_hash_mismatch', 'event_hash does not match canonical event without event_hash.', file, {
        expected: eventHash,
        actual: event.event_hash
      }));
    }
  }

  return findings;
};

if (eventFile) {
  const event = readJson(eventFile);
  const findings = verifyEventObject(event, eventFile);
  const passed = findings.length === 0;
  console.log(JSON.stringify({
    target: eventFile,
    verifier_states: {
      event_schema_valid: state(passed ? 'pass' : 'fail', passed ? undefined : 'Event file failed standalone verification.')
    },
    findings
  }, null, 2));
  process.exit(passed ? 0 : 1);
}

if (!PROFILE_ALIASES[profile]) {
  usage();
}
profile = PROFILE_ALIASES[profile];

const runFile = path.join(target, 'run.json');
const findings = [];
const run = readJson(runFile);
validateWithSchema(runSchema, run, runFile, findings);
const requiredProfile = RUN_PROFILE_REQUIREMENT[run.profile];
let profileCompatible = true;
if (requiredProfile && PROFILE_RANK[profile] < PROFILE_RANK[requiredProfile]) {
  profileCompatible = false;
  findings.push(finding('profile_downgrade', 'Selected verifier profile is weaker than run.json profile requires.', runFile, {
    selected_profile: profile,
    run_profile: run.profile,
    required_profile: requiredProfile
  }));
}

const eventRefs = Array.isArray(run.event_refs) ? run.event_refs : [];
const events = [];
for (const ref of eventRefs) {
  const file = path.join(target, ref.uri);
  if (!fs.existsSync(file)) {
    findings.push(finding('event_ref_missing', 'run.json references an event file that does not exist.', file, ref));
    continue;
  }
  const event = readJson(file);
  events.push({ ref, event, file });
  findings.push(...verifyEventObject(event, file));
}

let previousHash = ZERO_HASH;
let chainStructurallyValid = true;
for (let index = 0; index < events.length; index += 1) {
  const { ref, event, file } = events[index];

  if (event.sequence !== index || ref.sequence !== index) {
    chainStructurallyValid = false;
    findings.push(finding('event_sequence_mismatch', 'Event sequence must be contiguous and match run event_refs order.', file, {
      expected: index,
      event_sequence: event.sequence,
      ref_sequence: ref.sequence
    }));
  }
  if (ref.event_hash !== event.event_hash) {
    chainStructurallyValid = false;
    findings.push(finding('event_ref_hash_mismatch', 'run.json event_ref hash does not match the referenced event file.', file, {
      expected: event.event_hash,
      actual: ref.event_hash
    }));
  }
  if (event.prev_event_hash !== previousHash) {
    chainStructurallyValid = false;
    findings.push(finding('prev_event_hash_mismatch', 'Event prev_event_hash does not match the previous event hash in run order.', file, {
      expected: previousHash,
      actual: event.prev_event_hash
    }));
  }
  if (index > 0) {
    const previousTime = events[index - 1].event.event_time;
    if (previousTime && event.event_time < previousTime) {
      chainStructurallyValid = false;
      findings.push(finding('event_time_order_mismatch', 'Event time moves backward relative to the previous event in chain order.', file, {
        previous_event_time: previousTime,
        event_time: event.event_time
      }));
    }
  }
  previousHash = event.event_hash;
}

const lastEvent = events.at(-1)?.event;
if (run.event_chain_head && lastEvent && run.event_chain_head !== lastEvent.event_hash) {
  chainStructurallyValid = false;
  findings.push(finding('run_chain_head_mismatch', 'run.json event_chain_head does not match the final event hash.', runFile, {
    expected: lastEvent.event_hash,
    actual: run.event_chain_head
  }));
}

const runClosedEvent = events.find(({ event }) => event.event_type === 'run_closed');
if (runClosedEvent) {
  const runClosedIndex = events.indexOf(runClosedEvent);
  if (runClosedIndex !== events.length - 1) {
    chainStructurallyValid = false;
    findings.push(finding('run_closed_not_terminal', 'run_closed must be the terminal event in the recorded chain.', runClosedEvent.file, {
      run_closed_sequence: runClosedEvent.event.sequence,
      final_sequence: events.at(-1)?.event.sequence
    }));
  }
  const previousBeforeClosure = runClosedIndex > 0 ? events[runClosedIndex - 1].event.event_hash : ZERO_HASH;
  if (runClosedEvent.event.event_payload.chain_head !== previousBeforeClosure) {
    findings.push(finding('run_closed_chain_head_mismatch', 'run_closed payload must record the chain head before closure.', runClosedEvent.file, {
      expected: previousBeforeClosure,
      actual: runClosedEvent.event.event_payload.chain_head
    }));
  }
}

if (findings.some((item) => ['schema_invalid', 'payload_digest_mismatch', 'event_hash_mismatch'].includes(item.code))) {
  chainStructurallyValid = false;
}

const hasSourceBaseline = events.some(({ event }) => event.event_type === 'source_baseline_recorded');
const hasTaskLocked = events.some(({ event }) => event.event_type === 'task_locked' || event.event_type === 'expectation_locked');
const hasRunClosed = Boolean(runClosedEvent);
let sourceBaselineState = state(hasSourceBaseline ? 'pass' : 'fail', hasSourceBaseline ? undefined : 'No source_baseline_recorded event found.');
let taskLockedState = state(hasTaskLocked ? 'pass' : 'fail', hasTaskLocked ? undefined : 'No task_locked or expectation_locked event found.');
const sourceBaselineEvents = events.filter(({ event }) => event.event_type === 'source_baseline_recorded');
if (sourceBaselineEvents.length > 0) {
  const matchesRunSummary = sourceBaselineEvents.some(({ event }) =>
    event.event_payload.source_baseline_hash === run.source_summary?.source_baseline_hash &&
    event.event_payload.source_kind === run.source_summary?.source_kind
  );
  if (!matchesRunSummary) {
    sourceBaselineState = state('fail', 'Run source_summary is not bound to a source_baseline_recorded event.');
    findings.push(finding('source_summary_mismatch', 'run.json source_summary must match a source_baseline_recorded event.', runFile, {
      source_summary: run.source_summary,
      source_events: sourceBaselineEvents.map(({ event }) => ({
        sequence: event.sequence,
        source_baseline_hash: event.event_payload.source_baseline_hash,
        source_kind: event.event_payload.source_kind
      }))
    }));
  }
}
const taskLockEvents = events.filter(({ event }) => event.event_type === 'task_locked' || event.event_type === 'expectation_locked');
if (taskLockEvents.length > 0) {
  const matchesRunSummary = taskLockEvents.some(({ event }) =>
    event.event_payload.task_hash === run.task_summary?.task_hash &&
    event.event_payload.task_ref === run.task_summary?.task_ref
  );
  if (!matchesRunSummary) {
    taskLockedState = state('fail', 'Run task_summary is not bound to a task_locked or expectation_locked event.');
    findings.push(finding('task_summary_mismatch', 'run.json task_summary must match a task_locked or expectation_locked event.', runFile, {
      task_summary: run.task_summary,
      task_events: taskLockEvents.map(({ event }) => ({
        sequence: event.sequence,
        task_hash: event.event_payload.task_hash,
        task_ref: event.event_payload.task_ref
      }))
    }));
  }
}
let lateAttachState = state('not_assessed', 'No run_started event was available for late-attach assessment.');
const runStartedEvent = events.find(({ event }) => event.event_type === 'run_started');
if (runStartedEvent?.event.event_payload?.attachment_mode === 'full_run') {
  lateAttachState = state('pass');
} else if (runStartedEvent?.event.event_payload?.attachment_mode === 'mid_flight') {
  if (runStartedEvent.event.event_payload.late_attach_boundary?.state === 'not_assessed') {
    lateAttachState = state('pass');
  } else {
    lateAttachState = state('fail', 'Mid-flight recorder attachment must expose a not_assessed late_attach_boundary.');
    findings.push(finding('late_attach_boundary_missing', 'run_started mid_flight event lacks an explicit not_assessed boundary.', runStartedEvent.file));
  }
}

let redactionState = state('pass');
for (const { event, file } of events) {
  if (['cannot_verify', 'not_assessed'].includes(event.redaction_state?.state)) {
    redactionState = state('fail', 'At least one event has unresolved redaction state.');
    findings.push(finding('redaction_unresolved', 'Event redaction_state is unresolved and cannot support a closed evidence claim.', file, event.redaction_state));
  }
  const retention = event.event_payload?.output_retention;
  if (profile === 'forensic' && event.event_payload?.forensic_importance === 'critical' && retention?.mode === 'digest_only') {
    redactionState = state('fail', 'Forensic profile cannot accept critical evidence retained only as a digest.');
    findings.push(finding('forensic_digest_only_critical', 'Critical forensic evidence is digest-only and cannot support reconstruction.', file, retention));
  }
}
let witnessState = state('not_assessed', 'Local profile does not claim witnessed accountability.');

if (profile === 'witnessed' || profile === 'forensic') {
  witnessState = state('pass');
  if (!run.witness_ref?.uri) {
    witnessState = state('fail', 'Witnessed profile requires run.witness_ref.uri.');
    findings.push(finding('witness_ref_missing', 'Witnessed profile requires a witness_ref with a readable witness artifact.', runFile));
  } else {
    const witnessFile = path.join(target, run.witness_ref.uri);
    if (!fs.existsSync(witnessFile)) {
      witnessState = state('fail', 'Witness artifact referenced by run.json does not exist.');
      findings.push(finding('witness_ref_missing', 'run.json references witness material that does not exist.', witnessFile, run.witness_ref));
    } else {
      const witness = readJson(witnessFile);
      if (!validateWithSchema(witnessSchema, witness, witnessFile, findings)) {
        witnessState = state('fail', 'Witness artifact does not match flight-recorder witness schema.');
      }
      const expectedWitnessScope = profile === 'forensic' ? 'external_witness_extension' : 'local_file_witness';
      if (witness.witness_scope !== expectedWitnessScope) {
        witnessState = state('fail', 'Witness scope does not match verifier profile expectation.');
        findings.push(finding('witness_scope_mismatch', 'Witness scope must match the selected verifier profile.', witnessFile, {
          expected: expectedWitnessScope,
          actual: witness.witness_scope,
          profile
        }));
      }
      const witnessChecks = [
        ['run_id', run.run_id, witness.run_id, 'witness_run_id_mismatch', 'Witness run_id does not match run.json.'],
        ['source_baseline_hash', run.source_summary?.source_baseline_hash, witness.source_baseline_hash, 'witness_source_baseline_mismatch', 'Witness source baseline hash does not match run.json.'],
        ['task_hash', run.task_summary?.task_hash, witness.task_hash, 'witness_task_hash_mismatch', 'Witness task hash does not match run.json.'],
        ['recorder_version', events[0]?.event.recorder_identity?.recorder_version, witness.recorder_version, 'witness_recorder_version_mismatch', 'Witness recorder version does not match the recorded chain.'],
        ['chain_head', run.event_chain_head, witness.chain_head, 'witness_chain_head_mismatch', 'Witness chain_head does not match run.json event_chain_head.']
      ];
      for (const [field, expected, actual, code, message] of witnessChecks) {
        if (expected !== actual) {
          witnessState = state('fail', 'Witness artifact does not bind the run manifest and event chain.');
          findings.push(finding(code, message, witnessFile, { field, expected, actual }));
        }
      }
    }
  }
}

const verifierStates = {
  event_chain_structurally_valid: state(chainStructurallyValid ? 'pass' : 'fail', chainStructurallyValid ? undefined : 'Event chain failed structural verification.'),
  source_baseline_recorded: sourceBaselineState,
  task_locked: taskLockedState,
  run_closed: state(hasRunClosed ? 'pass' : 'fail', hasRunClosed ? undefined : 'No run_closed event found.'),
  event_chain_witnessed: witnessState,
  late_attach_boundary_explicit: lateAttachState,
  model_identity_recorded: state('not_assessed', 'Model identity verification is handled by a later slice.'),
  command_events_bound: state('not_assessed', 'Command binding verification is handled by a later slice.'),
  file_mutations_bound: state('not_assessed', 'File mutation binding verification is handled by a later slice.'),
  redaction_resolved: redactionState
};

const failedCoreState = [
  verifierStates.event_chain_structurally_valid,
  verifierStates.source_baseline_recorded,
  verifierStates.task_locked,
  verifierStates.run_closed,
  verifierStates.event_chain_witnessed,
  verifierStates.late_attach_boundary_explicit,
  verifierStates.redaction_resolved
].some((item) => item.state === 'fail') || !profileCompatible || findings.some((item) => item.code === 'run_closed_not_terminal');

console.log(JSON.stringify({
  profile,
  target,
  run_id: run.run_id,
  verifier_states: verifierStates,
  findings
}, null, 2));

process.exit(failedCoreState ? 1 : 0);
