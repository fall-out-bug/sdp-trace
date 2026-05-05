#!/usr/bin/env node
import fs from 'node:fs';
import path from 'node:path';
import process from 'node:process';

const usage = () => {
  console.error('Usage: query-flight-recorder.mjs --query <run-summary|provenance|gaps|command-timeline|file-mutations|test-evidence|requirement-timeline|redactions|witness-state> <fixture-dir>');
  process.exit(2);
};

const args = process.argv.slice(2);
let query = null;
let target = null;
for (let i = 0; i < args.length; i += 1) {
  if (args[i] === '--query') {
    query = args[i + 1];
    i += 1;
  } else if (!target) {
    target = args[i];
  } else {
    usage();
  }
}

if (!query || !target) {
  usage();
}

const readJson = (file) => {
  try {
    return JSON.parse(fs.readFileSync(file, 'utf8'));
  } catch (error) {
    console.error(`Unable to read JSON: ${file}`);
    console.error(error.message);
    process.exit(1);
  }
};
const run = readJson(path.join(target, 'run.json'));
const events = (run.event_refs || []).flatMap((ref) => {
  const file = path.join(target, ref.uri);
  if (!fs.existsSync(file)) {
    return [{ missing: true, ref, file }];
  }
  return [{ ref, file, event: readJson(file) }];
});

const eventPayloads = (type) => events
  .filter((entry) => entry.event?.event_type === type)
  .map((entry) => ({
    ...entry.event.event_payload,
    sequence: entry.event.sequence,
    event_type: entry.event.event_type,
    event_hash: entry.event.event_hash,
    redaction_state: entry.event.redaction_state
  }));

const output = (() => {
  switch (query) {
    case 'run-summary':
      return {
        run_id: run.run_id,
        profile: run.profile,
        trust_scope: run.trust_scope,
        created_at: run.created_at,
        closed_at: run.closed_at,
        event_count: events.length,
        event_chain_head: run.event_chain_head,
        verifier_states: run.verifier_states
      };
    case 'provenance':
      return {
        run_id: run.run_id,
        source: run.source_summary,
        task: run.task_summary,
        model: run.model_summary,
        harness: run.harness_summary
      };
    case 'gaps':
      return {
        run_id: run.run_id,
        late_attach_boundaries: eventPayloads('run_started')
          .filter((payload) => payload.attachment_mode === 'mid_flight')
          .map((payload) => payload.late_attach_boundary),
        not_assessed_states: Object.entries(run.verifier_states || {})
          .filter(([, value]) => value.state === 'not_assessed' || value.state === 'cannot_verify')
          .map(([name, value]) => ({ name, ...value })),
        missing_event_refs: events.filter((entry) => entry.missing).map((entry) => entry.ref)
      };
    case 'command-timeline': {
      const commandFinishes = new Map(eventPayloads('command_finished').map((payload) => [payload.command_id, payload]));
      return {
        run_id: run.run_id,
        commands: eventPayloads('command_started').map((payload) => ({
          sequence: payload.sequence,
          event_hash: payload.event_hash,
          command_id: payload.command_id,
          argv_digest: payload.argv_digest,
          working_directory: payload.working_directory,
          started_at: payload.started_at,
          task_hash: payload.task_hash || null,
          redaction_state: payload.redaction_state,
          finished_at: commandFinishes.get(payload.command_id)?.finished_at || null,
          exit_state: commandFinishes.get(payload.command_id)?.exit_state || null,
          stdout_retention: commandFinishes.get(payload.command_id)?.stdout_retention || null,
          stderr_retention: commandFinishes.get(payload.command_id)?.stderr_retention || null
        }))
      };
    }
    case 'file-mutations':
      return {
        run_id: run.run_id,
        file_mutations: [
          ...eventPayloads('file_state_observed'),
          ...eventPayloads('file_mutation_observed')
        ].map((payload) => ({
          sequence: payload.sequence,
          event_type: payload.event_type,
          event_hash: payload.event_hash,
          path_scope: payload.path_scope,
          source_baseline_hash: payload.source_baseline_hash,
          tree_or_diff_digest: payload.tree_or_diff_digest,
          attributed_command_id: payload.attributed_command_id || null
        }))
      };
    case 'test-evidence':
      return {
        run_id: run.run_id,
        test_evidence: eventPayloads('test_output_observed').map((payload) => ({
          sequence: payload.sequence,
          event_hash: payload.event_hash,
          test_command_id: payload.test_command_id,
          output_retention: payload.output_retention,
          forensic_importance: payload.forensic_importance || null,
          redaction_state: payload.redaction_state,
          verifier_state: run.verifier_states?.redaction_resolved || null
        }))
      };
    case 'requirement-timeline':
      return {
        run_id: run.run_id,
        locked_requirements: [
          ...eventPayloads('task_locked'),
          ...eventPayloads('expectation_locked')
        ].map((payload) => ({
          sequence: payload.sequence,
          event_hash: payload.event_hash,
          task_ref: payload.task_ref,
          task_hash: payload.task_hash
        })),
        supersessions: eventPayloads('requirement_superseded').map((payload) => ({
          sequence: payload.sequence,
          superseded_event_hash: payload.superseded_event_hash,
          replacement_task_hash: payload.replacement_task_hash,
          reason: payload.reason
        }))
      };
    case 'redactions':
      return {
        run_id: run.run_id,
        evidence_retention_summary: run.evidence_retention_summary,
        redaction_events: events
          .filter((entry) => entry.event)
          .map((entry) => ({
            sequence: entry.event.sequence,
              event_type: entry.event.event_type,
            redaction_state: entry.event.redaction_state || { state: 'cannot_verify', reason: 'Event redaction_state is missing.' },
            output_retention: entry.event.event_payload.output_retention || null
          }))
          .filter((entry) => entry.redaction_state?.state !== 'not_required' || entry.output_retention)
      };
    case 'witness-state': {
      const witnessPath = run.witness_ref?.uri ? path.join(target, run.witness_ref.uri) : null;
      const witness = witnessPath && fs.existsSync(witnessPath) ? readJson(witnessPath) : null;
      return {
        run_id: run.run_id,
        witness_ref: run.witness_ref || null,
        witness_present: witnessPath ? fs.existsSync(witnessPath) : false,
        witness_scope: witness?.witness_scope || null,
        witness_agreement: witness ? {
          run_id_matches: witness.run_id === run.run_id,
          source_baseline_matches: witness.source_baseline_hash === run.source_summary?.source_baseline_hash,
          task_hash_matches: witness.task_hash === run.task_summary?.task_hash,
          chain_head_matches: witness.chain_head === run.event_chain_head
        } : null,
        witness
      };
    }
    default:
      usage();
  }
})();

console.log(JSON.stringify(output, null, 2));
