package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/telemetry"
)

func parseTelemetryExportArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "export telemetry"}
	// Telemetry export accepts only the renderer profile, posture artifact, and
	// output target needed to replay metric generation.
	opts.setString("profile", "")
	opts.setString("cross-repo-posture", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Profile, posture artifact, and destination stay named so exports are
	// auditable from the command line alone.
	if rejectRest(opts, stderr, "export telemetry accepts only flags") {
		return nil, exitUsage, false
	}
	return requireTelemetryExportArgs(opts, stderr)
}

func requireTelemetryExportArgs(opts *flagSet, stderr io.Writer) (*flagSet, int, bool) {
	if err := requireTelemetryExportInputs(opts); err != nil {
		// Required input checks keep unsupported profiles and missing artifacts
		// as usage errors before any metric bytes are emitted.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func requireTelemetryExportInputs(opts *flagSet) error {
	if strings.TrimSpace(opts.stringValue("profile")) != telemetry.ProfilePrometheusTextV1 {
		// Telemetry export is intentionally profile-locked so future renderers
		// cannot be selected by typo or stale docs.
		return fmt.Errorf("export telemetry requires --profile prometheus-text-v1")
	}
	if strings.TrimSpace(opts.stringValue("cross-repo-posture")) == "" {
		// The posture artifact is the sole metric source; the CLI does not infer
		// repository posture from the working tree.
		return fmt.Errorf("export telemetry requires --cross-repo-posture")
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		// Metrics must either be explicitly written or deliberately streamed with
		// `--out -`.
		return fmt.Errorf("export telemetry requires --out")
	}
	return nil
}
