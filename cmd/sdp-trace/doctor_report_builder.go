package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func buildDoctorReport(options doctorOptions) (doctorReport, int) {
	defaultContract := trace.DefaultContract
	result := "offline_dev"
	exitCode := 0
	// Doctor reports local readiness; it never upgrades offline evidence to a
	// CI or external witness claim.
	contract, contractCheck := doctorContractCheck(options.ContractPath, defaultContract)
	result, exitCode = updateDoctorExitForCheck(result, exitCode, contractCheck)
	ciCheck := ciWitnessPrerequisiteCheck(options.Env)
	// Writable path checks prove only local filesystem readiness for future
	// artifacts; they do not create recorder evidence.
	outputDirCheck := writablePathCheck("output_directory", options.OutputDir, "run artifact output directory is writable")
	reportDirCheck := writablePathCheck("report_directory", options.ReportDir, "report artifact directory is writable")
	expectedEvidenceCheck := expectedEvidenceReferenceCheck(contract)
	// CI prerequisites are reported, but local doctor does not require them for
	// offline development readiness.
	// The contract check runs before evidence-reference checks so a bad
	// override cannot be hidden by default-contract coverage.
	result, exitCode = updateDoctorExitForLocalChecks(result, exitCode, outputDirCheck, reportDirCheck, expectedEvidenceCheck)
	// Contract loading and writable probes are live checks; CI identity remains
	// separately marked cannot_verify when absent.
	// Doctor assembles named facts only; the result field reflects verifier
	// state without becoming an opaque score.
	report := doctorReport{
		Command:     "doctor",
		Result:      result,
		Environment: doctorEnvironmentChecks(),
		// Environment and control-point sections stay separate so local process
		// facts cannot be mistaken for gate evidence.
		// Control points mix local passes and cannot-verify prerequisites so the
		// JSON report does not collapse them into one health score.
		ControlPoints:      doctorControlPointChecks(defaultContract, ciCheck, outputDirCheck, reportDirCheck, contractCheck, expectedEvidenceCheck),
		SafeRetentionModes: safeRetentionModes(),
	}
	// The report intentionally contains no aggregate score; callers inspect
	// named control points.
	return report, exitCode
}

func updateDoctorExitForLocalChecks(result string, exitCode int, checks ...doctorCheck) (string, int) {
	for _, check := range checks {
		// Only control points that the local process can inspect affect the
		// offline doctor exit code.
		result, exitCode = updateDoctorExitForCheck(result, exitCode, check)
	}
	return result, exitCode
}
