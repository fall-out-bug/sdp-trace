# Block 19 Adapter Capture Fixtures

These fixtures are generated from `internal/adaptercapture` cases and show
adapter capture-depth assessment results.

They are generic by design: no fixture encodes a specific harness, model
provider, Git host, CI service, or build system as a product concept.

Covered states:

- valid adapter capture with same-chain binding;
- valid adapter capture with adapter-bundle binding;
- missing required tool telemetry as `missing_telemetry`;
- unsupported observer capability as `unsupported`;
- gateway absence as `not_integrated`;
- agent-reported test evidence rejected as executed proof;
- harness-observed test evidence kept correlation-only without execution proof;
- file mutation evidence correlated to source and run identity;
- task supersession evidence with explicit actor attribution;
- provider-neutral source/change/review refs;
- late adapter event as `cannot_verify`;
- unsafe provider reference rejected before it can leak into output.
- conflicting adapter events as `cannot_verify`;
- capture-depth overclaim rejected as `fail`;
- unsafe event-level provider reference rejected as `fail`;
- duplicate empty correlation refs do not create a false conflict;
- redacted captured events can pass only when the retention cap is visible.

Some pass fixtures are intentionally byte-identical because they exercise
different condition IDs over the same valid generic adapter input. The fixture
name and test binding, not payload uniqueness, identify the acceptance row.
