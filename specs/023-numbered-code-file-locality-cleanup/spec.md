# Spec 023: Numbered Code File Locality Cleanup

Status: complete

## Objective

Replace remaining numbered Go source shards such as `family_123_name.go` with
cohesive behavior-named files across active code packages while preserving
behavior, tests, quality gates, and package boundaries.

## Background

Spec 021 removed the `main_[0-9]*.go` command-surface shards but repository
inventory still shows numbered Go files in active product code:

- `cmd/sdp-trace`: 430 files
- `internal/harnessobs`: 350 files
- `internal/packet`: 200 files
- `internal/prreview`: 192 files

The total remaining active numbered Go file count at intake is 1172.

## Requirements

- FR-023-001: Split cleanup by package, command family, or cohesive behavior;
  no repo-wide rename sweep.
- FR-023-002: Preserve current public command behavior and output contracts.
- FR-023-003: Group functions into cohesive files named after behavior,
  command family, or domain concept.
- FR-023-004: Keep complexity, MI, and CRAP gates as verification signals, not
  as a reason for metric-gaming file moves.
- FR-023-005: Preserve package boundaries and dependency direction from
  `docs/package-ownership-map.md`.

## Non-Goals

- No command behavior change.
- No package split unless another reviewed spec owns it.
- No broad formatting churn outside touched packages.
- No production trust, release approval, or external attestation claim.

## Acceptance Criteria

- Each cleanup slice lists the exact package or command family it touches.
- Tests and quality gates pass after each slice.
- Each slice has three independent staged-diff reviewer lanes before commit.
- Remaining active numbered Go file count decreases monotonically.

## Completion Evidence

Status updated to `complete` on 2026-06-05 after current-state inventory found
no remaining numbered Go source files in active product code paths:

```sh
test -z "$(rg --files -g '*.go' | rg '(^|/)[0-9]+|_type\.go$|_[0-9]+\.go$|[0-9].*\.go$' || true)"
```

Exit status: 0.

The completion claim does not include example event JSON, SpecKit artifact
names, checked-in review artifact numbering, PR merge approval, release
approval, or external attestation.

## Active Slice 1

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` release-proof command shards only.

Files selected for grouping:

- `cmd/sdp-trace/release_151_run.go`
- `cmd/sdp-trace/release_154_parseargs.go`
- `cmd/sdp-trace/release_155_flagsandexits.go`

Target files:

- `cmd/sdp-trace/release_proof_run.go`
- `cmd/sdp-trace/release_proof_args.go`
- `cmd/sdp-trace/release_proof_policy.go`

Rejected grouping:

- A single `release_proof_command.go` file was rejected because local
  pre-change MI analysis measured file MI `65.3`, below the absolute threshold.
- Combining args and exits was rejected because local pre-change MI analysis
  measured file MI `69.2`, below the absolute threshold.

Intended behavior boundary: this slice should only move release-proof command
runner, argument parsing, required flags, and exit-code declarations into
behavior-named files. No CLI behavior, JSON field, schema contract, or command
metadata value should change.

## Active Slice 2

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` observe command adapter and exit policy shards only.

Files selected for grouping:

- `cmd/sdp-trace/observe_012_runcommand.go`
- `cmd/sdp-trace/observe_013_runharness.go`
- `cmd/sdp-trace/observe_028_harnessstateexits.go`

Target files:

- `cmd/sdp-trace/observe_command_adapters.go`
- `cmd/sdp-trace/observe_exit_policy.go`

Rejected grouping:

- A single observe command file was rejected because local pre-change MI
  analysis measured file MI `64.1`, below the absolute threshold.

Intended behavior boundary: this slice should only move observe command adapter
functions and harness state exit-code policy into behavior-named files. No CLI
behavior, JSON field, schema contract, or command metadata value should change.

## Active Slice 3

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` envelope summarize command shards only.

Files selected for grouping:

- `cmd/sdp-trace/envelope_173_run.go`
- `cmd/sdp-trace/envelope_174_requiredflags.go`
- `cmd/sdp-trace/envelope_174_writeoptionaljsonfile.go`
- `cmd/sdp-trace/envelope_175_parsesummarizeargs.go`

Target files:

- `cmd/sdp-trace/envelope_summary_run.go`
- `cmd/sdp-trace/envelope_summary_args.go`

Rejected grouping:

- A single `envelope_command.go` file was rejected because local MI analysis
  measured file MI `69.0`, below the absolute threshold.

Intended behavior boundary: this slice should only move the envelope summarize
runner, argument parsing, required flags, and optional JSON output helper into
behavior-named command files. No CLI behavior, JSON field, schema contract, or
command metadata value should change.

## Active Slice 4

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` export command dispatcher shards only.

Files selected for grouping:

- `cmd/sdp-trace/export_452_run.go`
- `cmd/sdp-trace/export_453_telemetryrequested.go`
- `cmd/sdp-trace/export_454_crossrepopostureexplainrequested.go`
- `cmd/sdp-trace/export_455_crossrepoposturerequested.go`

Target file:

- `cmd/sdp-trace/export_command.go`

Intended behavior boundary: this slice should only move the export dispatcher
and export subcommand predicates into a behavior-named command file. No CLI
behavior, JSON field, schema contract, or command metadata value should change.

## Active Slice 5

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` fixture validation command shards only.

Files selected for grouping:

- `cmd/sdp-trace/fixture_472_run.go`
- `cmd/sdp-trace/fixture_473_rootarg.go`
- `cmd/sdp-trace/fixture_474_validatefixtureruns.go`
- `cmd/sdp-trace/fixture_475_validatefixturerun.go`
- `cmd/sdp-trace/fixture_476_expectationfailed.go`
- `cmd/sdp-trace/fixture_477_expectedresultfailed.go`
- `cmd/sdp-trace/fixture_478_unexpectedresultfailed.go`

Target files:

- `cmd/sdp-trace/fixture_validation_run.go`
- `cmd/sdp-trace/fixture_validation_args.go`
- `cmd/sdp-trace/fixture_expectation_policy.go`

Rejected grouping:

- A single `fixture_validation_command.go` file was rejected because local MI
  analysis measured file MI `65.4`, below the absolute threshold.
- Keeping fixture root argument parsing in `fixture_validation_run.go` was
  rejected because local MI analysis measured file MI `68.6`, below the
  absolute threshold.

Intended behavior boundary: this slice should only move fixture validation
runner, root selection, per-run validation, and fixture expectation policy into
behavior-named command files. No CLI behavior, verifier artifact behavior, JSON
field, schema contract, or command metadata value should change.

## Active Slice 6

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` interaction command shards only.

Files selected for grouping:

- `cmd/sdp-trace/interaction_158_run.go`
- `cmd/sdp-trace/interaction_159_runrelay.go`
- `cmd/sdp-trace/interaction_160_relayoptions.go`
- `cmd/sdp-trace/interaction_161_parserelayargs.go`
- `cmd/sdp-trace/interaction_162_newrelayflagset.go`
- `cmd/sdp-trace/interaction_163_relaystringflags.go`
- `cmd/sdp-trace/interaction_164_requirerest.go`
- `cmd/sdp-trace/interaction_165_requireonlyflagscode.go`
- `cmd/sdp-trace/interaction_166_requiredflags.go`
- `cmd/sdp-trace/interaction_167_runimporttranscript.go`
- `cmd/sdp-trace/interaction_168_writeimportedtranscript.go`
- `cmd/sdp-trace/interaction_169_importtranscriptfromoptions.go`
- `cmd/sdp-trace/interaction_170_parseimporttranscriptargs.go`
- `cmd/sdp-trace/interaction_171_runsummarize.go`
- `cmd/sdp-trace/interaction_172_parsesummarizeargs.go`

Target files:

- `cmd/sdp-trace/interaction_command.go`
- `cmd/sdp-trace/interaction_relay.go`
- `cmd/sdp-trace/interaction_relay_args.go`
- `cmd/sdp-trace/interaction_transcript_import.go`
- `cmd/sdp-trace/interaction_transcript_import_args.go`
- `cmd/sdp-trace/interaction_summary.go`
- `cmd/sdp-trace/cli_flag_requirements.go`

Rejected grouping:

- Keeping transcript import runner, writer, options mapping, and args parsing
  in one file was rejected because local MI analysis measured file MI `69.5`,
  below the absolute threshold.

Intended behavior boundary: this slice should only move interaction router,
relay, transcript import, summarize, and directly related CLI flag helpers into
behavior-named command files. No CLI behavior, command forwarding behavior,
JSON field, trace import, summary, schema contract, or command metadata value
should change.

## Active Slice 7

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` wrap, run, preview, and dry-run command shards only.

Files selected for grouping:

- `cmd/sdp-trace/wrap_399_run.go`
- `cmd/sdp-trace/wrap_400_runlegacyrecorder.go`
- `cmd/sdp-trace/wrap_401_parseargs.go`
- `cmd/sdp-trace/wrap_402_command.go`
- `cmd/sdp-trace/wrap_403_writerunresult.go`
- `cmd/sdp-trace/wrap_404_runwrappedcommand.go`
- `cmd/sdp-trace/wrap_405_runtaskrecorder.go`
- `cmd/sdp-trace/wrap_406_parsewrappedcommandargs.go`
- `cmd/sdp-trace/wrap_407_requirewrappedcommandargs.go`
- `cmd/sdp-trace/wrap_408_missingrequiredcontract.go`
- `cmd/sdp-trace/wrap_409_rundryrun.go`
- `cmd/sdp-trace/wrap_410_runpreview.go`
- `cmd/sdp-trace/wrap_411_runpreviewcommand.go`
- `cmd/sdp-trace/wrap_412_writepreviewcommandpayload.go`
- `cmd/sdp-trace/wrap_413_previewcommandpayload.go`
- `cmd/sdp-trace/wrap_414_parsepreviewcommandargs.go`
- `cmd/sdp-trace/wrap_415_loadpreviewcontract.go`

Target files:

- `cmd/sdp-trace/wrap_legacy.go`
- `cmd/sdp-trace/wrap_recorder.go`
- `cmd/sdp-trace/wrap_run.go`
- `cmd/sdp-trace/wrap_run_args.go`
- `cmd/sdp-trace/wrap_preview.go`
- `cmd/sdp-trace/wrap_preview_args.go`
- `cmd/sdp-trace/wrap_preview_payload.go`

Rejected grouping:

- Keeping preview runner, payload writing, and payload construction in one file
  was rejected because local MI analysis measured file MI `67.9`, below the
  absolute threshold.

Intended behavior boundary: this slice should only move legacy wrap, modern
run, preview, dry-run, recorder execution, and directly related argument and
payload helpers into behavior-named command files. No CLI behavior, recorder
artifact behavior, preview JSON field, schema contract, or command metadata
value should change.

## Active Slice 8

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` query, verify, explain, and query-pack command shards
only.

Files selected for grouping:

- `cmd/sdp-trace/query_430_runverify.go`
- `cmd/sdp-trace/query_431_parseverifyargs.go`
- `cmd/sdp-trace/query_432_existingdirectory.go`
- `cmd/sdp-trace/query_433_verifierresultexitcodes.go`
- `cmd/sdp-trace/query_434_verifierresultexitcode.go`
- `cmd/sdp-trace/query_435_runexplain.go`
- `cmd/sdp-trace/query_436_runquery.go`
- `cmd/sdp-trace/query_437_runnamedquery.go`
- `cmd/sdp-trace/query_438_capturedepthquery.go`
- `cmd/sdp-trace/query_439_missingevidencequery.go`
- `cmd/sdp-trace/query_440_runquerypack.go`
- `cmd/sdp-trace/query_441_runquerypackbuild.go`
- `cmd/sdp-trace/query_442_writequerypackartifact.go`
- `cmd/sdp-trace/query_443_runquerypackexplain.go`
- `cmd/sdp-trace/query_444_options.go`
- `cmd/sdp-trace/query_446_parsequerypackargs.go`
- `cmd/sdp-trace/query_447_parsequerypackexplainargs.go`
- `cmd/sdp-trace/query_448_validatequerypackoptions.go`
- `cmd/sdp-trace/query_449_requirequerypackrequiredinputs.go`
- `cmd/sdp-trace/query_450_readquerypackresult.go`
- `cmd/sdp-trace/query_451_validatequerypackexplainresult.go`

Target files:

- `cmd/sdp-trace/query_verify.go`
- `cmd/sdp-trace/query_verify_args.go`
- `cmd/sdp-trace/query_verify_exit.go`
- `cmd/sdp-trace/query_explain.go`
- `cmd/sdp-trace/query_run.go`
- `cmd/sdp-trace/query_dispatch.go`
- `cmd/sdp-trace/query_pack.go`
- `cmd/sdp-trace/query_pack_build.go`
- `cmd/sdp-trace/query_pack_explain.go`
- `cmd/sdp-trace/query_pack_args.go`
- `cmd/sdp-trace/query_pack_validation.go`

Rejected grouping:

- Keeping verify runner, args, directory check, and exit policy in one file was
  rejected because local MI analysis measured file MI `66.7`, below the
  absolute threshold.
- Keeping query runner, dispatch, and diagnostic query execution in one file
  was rejected because local MI analysis measured file MI `68.5`, below the
  absolute threshold.
- Keeping query-pack routing, build, artifact write, and explain in one file
  was rejected because local MI analysis measured file MI `66.7`, below the
  absolute threshold.

Intended behavior boundary: this slice should only move query, verify, explain,
query-pack build, query-pack explain, parsing, validation, and exit policy into
behavior-named command files. No CLI behavior, verifier artifact behavior,
query JSON payload, query-pack artifact, schema contract, or command metadata
value should change.

## Active Slice 9

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` witness command shards only.

Files selected for grouping:

- `cmd/sdp-trace/witness_374_run.go`
- `cmd/sdp-trace/witness_375_options.go`
- `cmd/sdp-trace/witness_376_parseoptions.go`
- `cmd/sdp-trace/witness_377_parseflagset.go`
- `cmd/sdp-trace/witness_378_optionsfromflags.go`
- `cmd/sdp-trace/witness_379_requiredfields.go`
- `cmd/sdp-trace/witness_380_requiredfieldsfromflags.go`
- `cmd/sdp-trace/witness_381_kindoutfromflags.go`
- `cmd/sdp-trace/witness_382_optionsfromrequiredfields.go`
- `cmd/sdp-trace/witness_384_kindfromflags.go`
- `cmd/sdp-trace/witness_385_outfromflags.go`
- `cmd/sdp-trace/witness_386_validatekindflags.go`
- `cmd/sdp-trace/witness_387_missingkindflags.go`
- `cmd/sdp-trace/witness_388_buildrecord.go`
- `cmd/sdp-trace/witness_389_builders.go`
- `cmd/sdp-trace/witness_391_buildgithubactions.go`
- `cmd/sdp-trace/witness_392_buildenvelope.go`
- `cmd/sdp-trace/witness_393_buildcustomerpki.go`
- `cmd/sdp-trace/witness_394_writerecordoutput.go`
- `cmd/sdp-trace/witness_395_missingcustomerpkiflags.go`
- `cmd/sdp-trace/witness_396_appendmissingstringflags.go`
- `cmd/sdp-trace/witness_397_missingcustomerpkipubliccredential.go`
- `cmd/sdp-trace/witness_398_allowedwitnesskind.go`

Target files:

- `cmd/sdp-trace/witness_command.go`
- `cmd/sdp-trace/witness_options.go`
- `cmd/sdp-trace/witness_options_parse.go`
- `cmd/sdp-trace/witness_options_build.go`
- `cmd/sdp-trace/witness_flag_set.go`
- `cmd/sdp-trace/witness_required_fields.go`
- `cmd/sdp-trace/witness_kind_validation.go`
- `cmd/sdp-trace/witness_record_builders.go`
- `cmd/sdp-trace/witness_output.go`
- `cmd/sdp-trace/witness_customer_pki.go`

Rejected grouping:

- Keeping all option parsing and option construction in one file was rejected
  because local MI analysis measured file MI `68.4`, below the absolute
  threshold.
- Keeping required-field and kind validation in one file was rejected because
  local MI analysis measured file MI `65.2`, below the absolute threshold.

Intended behavior boundary: this slice should only move witness command
parsing, required-field validation, witness builders, customer-PKI missing flag
helpers, and output rendering into behavior-named command files. No CLI
behavior, witness JSON field, schema contract, or command metadata value should
change.

## Active Slice 10

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` doctor repo-observer and install command shards only.
The local doctor report/check shards remain outside this slice.

Files selected for grouping:

- `cmd/sdp-trace/doctor_416_rundoctor.go`
- `cmd/sdp-trace/doctor_417_parsedoctorargs.go`
- `cmd/sdp-trace/doctor_418_runrepoobserverdoctor.go`
- `cmd/sdp-trace/doctor_419_writerepoobserverdoctor.go`
- `cmd/sdp-trace/doctor_421_runinstall.go`
- `cmd/sdp-trace/doctor_422_handlerepoobserverinstallerror.go`
- `cmd/sdp-trace/doctor_423_repoobserverinstallexitcode.go`
- `cmd/sdp-trace/doctor_424_parseinstallrepoobserverargs.go`
- `cmd/sdp-trace/doctor_425_requireinstallrepoobserverflags.go`
- `cmd/sdp-trace/doctor_426_hasinstallrepoobserversubcommand.go`
- `cmd/sdp-trace/doctor_427_installrepoobserverflagset.go`
- `cmd/sdp-trace/doctor_428_repoobserveroptionsfromflags.go`
- `cmd/sdp-trace/doctor_429_repoobserverexitcode.go`

Target files:

- `cmd/sdp-trace/doctor_command.go`
- `cmd/sdp-trace/doctor_repo_observer.go`
- `cmd/sdp-trace/doctor_install.go`
- `cmd/sdp-trace/doctor_install_args.go`
- `cmd/sdp-trace/doctor_install_options.go`

Rejected grouping:

- Keeping install argument parsing, flag registration, and option conversion in
  one file was rejected because local MI analysis measured file MI `68.1`,
  below the absolute threshold.

Intended behavior boundary: this slice should only move doctor routing,
doctor argument parsing, repo-observer doctor output, repo-observer install
execution, install argument parsing, and install option construction into
behavior-named command files. No CLI behavior, repo-observer JSON field,
schema contract, or command metadata value should change.
