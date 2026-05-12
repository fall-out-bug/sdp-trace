# Evidence Bundle Manifest v0

`evidence-bundle-manifest.v0` binds packet row refs to retained evidence.

Each entry has:

- `ref`: evidence ref used by packet rows or theater findings.
- `source_class`: source class such as `git`, `ci`, `harness`, `review`,
  `witness`, `change_host`, `external_assertion`, or `manual`.
- `digest`: digest of retained evidence or retained metadata when available.
- `retained_form`: `raw`, `redacted`, `digest_only`, `external_ref`, or
  `not_retained`.
- `redaction_status`: `not_needed`, `redacted`, `digest_only`, `withheld`, or
  `cannot_verify`.
- `resolver`: how a reviewer can resolve the ref.
- `expires_at` and `artifact_access`: optional artifact retention state.
- `packet_digest`: digest binding for the canonical packet artifact.
- `actor`: optional actor label such as `recorder`, `ci_packet_builder`,
  `operator`, `integration`, or `developer_route`.
- `write_authority`: optional authority label such as `recorder_owned`,
  `ci_generated`, `operator_authored`, or `integration_authored`.
- `generated_by`, `source_commit_state`, and `source_ref`: optional source
  binding metadata for generated packet artifacts.

If a packet row cites a ref absent from the manifest, the validator rejects the
bundle. If the ref exists but has no resolver, the validator rejects the bundle.
If a `pass` row cites an expired or unverifiable artifact, the validator rejects
the bundle.
The packet `bundle_ref` must match the manifest `bundle_id`.

PR comments and PR bodies are projections. They can point at the canonical
packet artifact, but they are not canonical packet evidence by themselves.
Operator-authored and integration-authored entries may explain packet history,
but they do not satisfy developer-route proof by themselves.
