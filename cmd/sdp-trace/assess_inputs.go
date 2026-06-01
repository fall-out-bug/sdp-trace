package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
	"path/filepath"
)

func loadForensicInput(opts *flagSet) (forensic.Input, error) {
	var policy forensic.Policy
	if err := readJSONFile(opts.stringValue("redaction-policy"), &policy); err != nil {
		return forensic.Input{}, err
	}
	var runEvidence forensic.RunEvidence
	// Forensic assessment reuses the run directory contract so raw and sanitized
	// evidence stay tied to the captured run.
	if err := readJSONFile(filepath.Join(opts.stringValue("run"), "run.json"), &runEvidence); err != nil {
		return forensic.Input{}, err
	}
	return forensic.Input{Policy: policy, Run: runEvidence}, nil
}

func loadAdapterCaptureInput(opts *flagSet) (adaptercapture.Input, error) {
	var runEvidence adaptercapture.RunEvidence
	// Adapter capture currently consumes the normalized run evidence file only.
	if err := readJSONFile(filepath.Join(opts.stringValue("run"), "run.json"), &runEvidence); err != nil {
		return adaptercapture.Input{}, err
	}
	return adaptercapture.Input{Run: runEvidence}, nil
}
