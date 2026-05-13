package main

import (
	"fmt"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/telemetry"
)

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
