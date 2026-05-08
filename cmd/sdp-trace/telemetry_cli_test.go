package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestTelemetryExportWritesPrometheusTextToStdoutAndFile(t *testing.T) {
	fixture := "../../examples/block21-cross-repo-posture/valid-movement/cross-repo-posture-export.json"

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"export", "telemetry",
		"--profile", "prometheus-text-v1",
		"--cross-repo-posture", fixture,
		"--out", "-",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("stdout export exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	rendered := out.String()
	for _, want := range []string{
		"# HELP sdp_trace_posture_metric_numerator ",
		"# TYPE sdp_trace_posture_metric_numerator gauge",
		"sdp_trace_posture_metric_numerator",
		"sdp_trace_posture_movement_delta",
		"sdp_trace_posture_input",
		"# EOF\n",
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("missing %q in telemetry output", want)
		}
	}
	if strings.Contains(rendered, "_total") ||
		strings.Contains(rendered, "https://") ||
		strings.Contains(rendered, "secret") ||
		strings.Contains(rendered, "/private") {
		t.Fatalf("unsafe or convention-breaking telemetry output")
	}

	out.Reset()
	errOut.Reset()
	file := t.TempDir() + "/telemetry.prom"
	exit = run([]string{
		"export", "telemetry",
		"--profile", "prometheus-text-v1",
		"--cross-repo-posture", fixture,
		"--out", file,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("file export exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if strings.TrimSpace(out.String()) != "" {
		t.Fatalf("file export wrote stdout: %s", out.String())
	}
	payload, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read telemetry file: %v", err)
	}
	if string(payload) != rendered {
		t.Fatalf("file output differs from stdout output")
	}
	committed, err := os.ReadFile("../../examples/block27-observability-telemetry/prometheus/metrics.prom")
	if err != nil {
		t.Fatalf("read committed telemetry fixture: %v", err)
	}
	if string(committed) != rendered {
		t.Fatalf("committed telemetry fixture drifted from live export")
	}
}

func TestTelemetryExportRejectsBadInputs(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"export", "telemetry", "--profile", "wrong"}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("exit=%d err=%s", exit, errOut.String())
	}
	if !strings.Contains(errOut.String(), "requires --profile prometheus-text-v1") {
		t.Fatalf("unexpected profile error: %s", errOut.String())
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{
		"export", "telemetry",
		"--profile", "prometheus-text-v1",
		"--cross-repo-posture", "missing.json",
		"--out", "-",
	}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("exit=%d err=%s", exit, errOut.String())
	}
	if !strings.Contains(errOut.String(), "posture_unreadable") {
		t.Fatalf("unexpected unreadable error: %s", errOut.String())
	}
}
