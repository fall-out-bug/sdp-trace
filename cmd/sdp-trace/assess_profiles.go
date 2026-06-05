package main

import (
	"fmt"
	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
	"github.com/fall_out_bug/sdp-trace/internal/managed"
	"io"
)

func runAdapterCaptureAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireAdapterCaptureAssessInputs(opts, stderr) {
		// Missing durable input/output flags are usage failures before
		// adaptercapture can evaluate run evidence.
		return exitUsage
	}
	// Adapter-capture assessment is run-bound only; missing run evidence is a
	// usage error before any verdict artifact exists.
	input, err := loadAdapterCaptureInput(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := adaptercapture.Evaluate(input)
	return writeAssessmentArtifact(opts.stringValue("out"), result, stdout, stderr, adapterCaptureExitCode)
}

func runManagedAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireManagedAssessInputs(opts, stderr) {
		// Managed assessment has no implicit defaults for registry, policy, or
		// witness authority.
		return exitUsage
	}
	// Managed-harness assessment joins contract, policy, registry, run, and
	// witness evidence before deriving a trust state.
	input, err := loadManagedInput(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := managed.Evaluate(input)
	return writeAssessmentArtifact(opts.stringValue("out"), result, stdout, stderr, managedExitCode)
}

func runForensicAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireForensicAssessInputs(opts, stderr) {
		// A forensic verdict without a redaction policy would overclaim retention
		// coverage.
		return exitUsage
	}
	// Forensic retention assessment requires both policy and run evidence so
	// missing redaction rules cannot be treated as passing defaults.
	input, err := loadForensicInput(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := forensic.Evaluate(input)
	return writeAssessmentArtifact(opts.stringValue("out"), result, stdout, stderr, forensicExitCode)
}
