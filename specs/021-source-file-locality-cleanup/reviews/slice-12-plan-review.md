# Slice 12 Plan Review: Doctor Local Report Cleanup

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target files:
  - `cmd/sdp-trace/doctor_local.go`
  - `cmd/sdp-trace/doctor_report.go`
  - `cmd/sdp-trace/doctor_report_builder.go`
  - `cmd/sdp-trace/doctor_report_checks.go`
  - `cmd/sdp-trace/doctor_contract.go`
  - `cmd/sdp-trace/doctor_event_types.go`
  - `cmd/sdp-trace/doctor_writable_path.go`
  - `cmd/sdp-trace/doctor_writable_probe.go`
  - `cmd/sdp-trace/doctor_writable_results.go`
  - `cmd/sdp-trace/doctor_expected_evidence.go`
  - `cmd/sdp-trace/doctor_expected_gaps.go`
  - `cmd/sdp-trace/doctor_ci.go`
  - `cmd/sdp-trace/doctor_ci_env.go`
  - `cmd/sdp-trace/doctor_preview.go`
  - `cmd/sdp-trace/doctor_preview_offline.go`
  - `cmd/sdp-trace/doctor_usage.go`
  - `cmd/sdp-trace/doctor_usage_primary.go`
  - `cmd/sdp-trace/doctor_usage_trust.go`
  - `cmd/sdp-trace/doctor_usage_packet.go`
- Source shards:
  - `cmd/sdp-trace/doctor_420_runlocaldoctor.go`
  - `cmd/sdp-trace/doctor_479_doctorreport.go`
  - `cmd/sdp-trace/doctor_480_doctorcheck.go`
  - `cmd/sdp-trace/doctor_481_doctoroptions.go`
  - `cmd/sdp-trace/doctor_482_previewboundary.go`
  - `cmd/sdp-trace/doctor_483_previewofflineimplication.go`
  - `cmd/sdp-trace/doctor_484_const.go`
  - `cmd/sdp-trace/doctor_485_builddoctorreport.go`
  - `cmd/sdp-trace/doctor_486_updatedoctorexitforlocalchecks.go`
  - `cmd/sdp-trace/doctor_487_doctorenvironmentchecks.go`
  - `cmd/sdp-trace/doctor_488_doctorcontrolpointchecks.go`
  - `cmd/sdp-trace/doctor_489_doctordefaultcontractcheck.go`
  - `cmd/sdp-trace/doctor_490_doctorcontractcheck.go`
  - `cmd/sdp-trace/doctor_491_loadeddoctorcontractresult.go`
  - `cmd/sdp-trace/doctor_492_defaultdoctorcontractresult.go`
  - `cmd/sdp-trace/doctor_493_unreadabledoctorcontractresult.go`
  - `cmd/sdp-trace/doctor_494_updatedoctorexitforcheck.go`
  - `cmd/sdp-trace/doctor_495_writablepathcheck.go`
  - `cmd/sdp-trace/doctor_496_probewritablepath.go`
  - `cmd/sdp-trace/doctor_497_emptywritablepathcheck.go`
  - `cmd/sdp-trace/doctor_498_writablepathpasscheck.go`
  - `cmd/sdp-trace/doctor_499_writableprobetarget.go`
  - `cmd/sdp-trace/doctor_500_writableprobeparent.go`
  - `cmd/sdp-trace/doctor_501_expectedevidencereferencecheck.go`
  - `cmd/sdp-trace/doctor_502_expectedevidencenorequiredeventscheck.go`
  - `cmd/sdp-trace/doctor_503_expectedevidenceunsupportedreferencecheck.go`
  - `cmd/sdp-trace/doctor_504_expectedevidencereferencegaps.go`
  - `cmd/sdp-trace/doctor_505_expectedevidencegaps.go`
  - `cmd/sdp-trace/doctor_506_knowneventtype.go`
  - `cmd/sdp-trace/doctor_507_ciwitnessprerequisitecheck.go`
  - `cmd/sdp-trace/doctor_508_missingciwitnessfields.go`
  - `cmd/sdp-trace/doctor_509_requiredciwitnessenvfields.go`
  - `cmd/sdp-trace/doctor_510_missingenvfields.go`
  - `cmd/sdp-trace/doctor_511_saferetentionmodes.go`
  - `cmd/sdp-trace/doctor_512_previewboundaryrows.go`
  - `cmd/sdp-trace/doctor_513_previewboundaries.go`
  - `cmd/sdp-trace/doctor_514_previewofflineimplications.go`
  - `cmd/sdp-trace/doctor_515_usagetext.go`
  - `cmd/sdp-trace/doctor_516_printusage.go`

## Decision Gate

- Simpler/Faster: move the remaining doctor local report/check declarations into
  responsibility-named files without changing function bodies, command behavior,
  JSON fields, package boundaries, or dependencies.
- Blocking Edge Cases: one combined doctor file would couple unrelated
  responsibilities and risks maintainability-index regression. Splitting by
  doctor responsibility keeps review small enough to compare mechanically.
- Existing Open Source: not applicable; this is a local Go file locality cleanup
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate for new doctor files
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- zero remaining `cmd/sdp-trace/doctor_[0-9]*_*.go` files
- three independent staged-diff reviewer lanes after implementation

## Review Result

- Plan/task scope: pass
- Risk if wrong: medium; behavior-preserving file moves can still drift if a
  declaration is omitted or help text chunk ordering changes.
- Reversibility: high; all changes are local file moves/splits in one package.
- Confidence: high after focused compile, MI, full repository gates, and review
  correction for over-decomposed usage text.
