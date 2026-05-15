# Claim Authoring

Authoritative repository claims must use `sdp-trace-claim` comments. Prose, checkboxes, reports, and review discussion are commentary unless they carry a claim tag.

## Grammar

```text
<!-- sdp-trace-claim: claim=<kind>; subject=<id>; state=<state>; profile=<profile>; evidence=<evidence-ref> -->
```

Allowed Slice 1 values:

- `claim`: `task_closed`, `command_verified`, `profile_passed`, `trust_not_assessed`
- `state`: `pass`, `fail`, `not_assessed`, `stale`, `cannot_verify`
- `profile`: `repo_baseline_structural`, `source_bound_local_release`, `external_production_trust`, `observed_slice`
- `evidence`: `command_set:block04-t070` or `state:claim_tags_consistent`

`state=pass` with `command_set:block04-t070` replays the current closure commands. If any command fails or cannot verify, the claim fails. Use `state=stale` for historical closure records that contradict the current verifier.

Slice 1 does not accept arbitrary `state:*`, `proof:*`, or `none` evidence. Those require cross-reference verification in a later slice.

Untagged prose that says a claim is false, stale, verified, complete, or trusted is not machine-authoritative. It may still be misleading to a human reviewer, but it cannot close proof.
