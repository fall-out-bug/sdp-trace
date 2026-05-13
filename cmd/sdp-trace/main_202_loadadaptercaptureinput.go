package main

import (
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
)

func loadAdapterCaptureInput(opts *flagSet) (adaptercapture.Input, error) {
	var runEvidence adaptercapture.RunEvidence
	// Adapter capture currently consumes the normalized run evidence file only.
	if err := readJSONFile(filepath.Join(opts.stringValue("run"), "run.json"), &runEvidence); err != nil {
		return adaptercapture.Input{}, err
	}
	return adaptercapture.Input{Run: runEvidence}, nil
}
