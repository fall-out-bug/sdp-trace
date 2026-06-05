package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

func runLocalDoctor(opts *flagSet, stdout, stderr io.Writer) int {
	// Local doctor checks CLI defaults and environment-derived witness context;
	// it does not inspect repository install hooks.
	report, exitCode := buildDoctorReport(doctorOptions{
		ContractPath: opts.stringValue("contract"),
		OutputDir:    opts.stringValue("output-dir"),
		ReportDir:    opts.stringValue("report-dir"),
		Env:          witness.EnvironmentFromOS(),
	})
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		// A doctor report that cannot be serialized cannot be trusted by
		// automation even if the underlying checks ran.
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "%s\n", data)
	return exitCode
}
