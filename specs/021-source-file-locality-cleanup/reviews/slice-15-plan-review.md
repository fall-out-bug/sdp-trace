# Slice 15 Plan Review: Core Assessment Explain And Preview Cleanup

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Source shards:
  - `cmd/sdp-trace/core_234_previewinputcannotverify.go`
  - `cmd/sdp-trace/core_235_runassessexplain.go`
  - `cmd/sdp-trace/core_236_parseassessexplainargs.go`
  - `cmd/sdp-trace/core_237_explainassessmentresult.go`
  - `cmd/sdp-trace/core_238_dispatchassessmentexplanation.go`
  - `cmd/sdp-trace/core_239_assessmentexplainhandler.go`
  - `cmd/sdp-trace/core_240_assessmentexplainhandlers.go`
  - `cmd/sdp-trace/core_241_explaintypedassessment_t_any.go`
  - `cmd/sdp-trace/core_242_explainadaptercaptureassessment.go`
  - `cmd/sdp-trace/core_243_explainadaptercapturecondition.go`
  - `cmd/sdp-trace/core_244_explainmanagedassessment.go`
  - `cmd/sdp-trace/core_245_explainforensicassessment.go`
  - `cmd/sdp-trace/core_246_explainforensiccondition.go`
  - `cmd/sdp-trace/core_247_explainciartifactobservation.go`
  - `cmd/sdp-trace/core_248_explainciartifactfamilies.go`
  - `cmd/sdp-trace/core_249_explainauthorityevaluation.go`
  - `cmd/sdp-trace/core_250_explainauthorityactionevaluations.go`
  - `cmd/sdp-trace/core_251_explainauthoritybindingevaluations.go`
  - `cmd/sdp-trace/core_252_managedinputstatus.go`
  - `cmd/sdp-trace/core_253_jsonreadablestatus.go`
  - `cmd/sdp-trace/core_254_managedpreviewactions.go`
  - `cmd/sdp-trace/core_255_forensicpreviewactions.go`
  - `cmd/sdp-trace/core_256_previewactionsforinputs.go`
  - `cmd/sdp-trace/core_257_adaptercapturepreviewactions.go`
  - `cmd/sdp-trace/core_258_ciartifactpreviewactions.go`
  - `cmd/sdp-trace/core_259_authoritypreviewactions.go`
  - `cmd/sdp-trace/core_260_previewactionforinputstate.go`
  - `cmd/sdp-trace/core_261_adaptercaptureexitcode.go`
  - `cmd/sdp-trace/core_262_adaptercaptureexitcodes.go`
  - `cmd/sdp-trace/core_263_managedexitcode.go`
  - `cmd/sdp-trace/core_264_managedexitcodes.go`
  - `cmd/sdp-trace/core_265_forensicexitcode.go`
  - `cmd/sdp-trace/core_266_forensicexitcodes.go`
  - `cmd/sdp-trace/core_267_ciartifactexitcode.go`
  - `cmd/sdp-trace/core_268_ciartifactexitcodes.go`
  - `cmd/sdp-trace/core_269_authorityexitcodes.go`
  - `cmd/sdp-trace/core_270_authorityexitcode.go`
- Target files:
  - `cmd/sdp-trace/assess_explain_command.go`
  - `cmd/sdp-trace/assess_explain_loader.go`
  - `cmd/sdp-trace/assess_explain_registry.go`
  - `cmd/sdp-trace/assess_explain_adapter.go`
  - `cmd/sdp-trace/assess_explain_managed.go`
  - `cmd/sdp-trace/assess_explain_forensic.go`
  - `cmd/sdp-trace/assess_explain_ci_artifact.go`
  - `cmd/sdp-trace/assess_explain_authority.go`
  - `cmd/sdp-trace/assess_preview_input_status.go`
  - `cmd/sdp-trace/assess_preview_actions.go`
  - `cmd/sdp-trace/assess_preview_action_helpers.go`
  - `cmd/sdp-trace/assess_exit_codes.go`
  - `cmd/sdp-trace/assess_exit_codes_artifacts.go`
  - `cmd/sdp-trace/assess_exit_code_lookup.go`
  - `cmd/sdp-trace/assess_exit_code_lookup_artifacts.go`

## Decision Gate

- Simpler/Faster: move declarations into responsibility files in the same
  package without changing function bodies, command names, explanation output,
  preview remediation text, exit-code values, or dependencies.
- Blocking Edge Cases: one combined assessment explain/preview/exit-code file
  mixes user-facing explanation, setup preview, and exit-code mapping
  responsibilities. A later reviewer-proposed single exit-code file was also
  measured and rejected because it failed file MI with
  `maintainability index 54.9 under threshold 70.0`.
- Existing Open Source: not applicable; this is local Go file locality cleanup
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate for changed assessment files
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- no remaining `cmd/sdp-trace/core_23*_*.go`,
  `cmd/sdp-trace/core_24*_*.go`, `cmd/sdp-trace/core_25*_*.go`, or
  `cmd/sdp-trace/core_26*_*.go` files from this slice scope
- three independent staged-diff reviewer lanes after implementation
