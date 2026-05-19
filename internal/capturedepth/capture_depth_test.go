package capturedepth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
)

func TestCaptureDepthQueryExposesReadOnlyFacts(t *testing.T) {
	run := adaptercapture.ValidTestInput().Run
	run.UnverifiedTaskExpanded = true
	run.TaskSupersessionCount = 2
	run.UnsupportedEventTypes = []string{"tool_call"}
	filtered := []adaptercapture.AdapterEvent{}
	for _, event := range run.AdapterEvents {
		if event.EventType != "tool_call" {
			filtered = append(filtered, event)
		}
	}
	run.AdapterEvents = filtered

	runDir := t.TempDir()
	payload, err := json.Marshal(run)
	if err != nil {
		t.Fatalf("marshal run: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "run.json"), payload, 0o644); err != nil {
		t.Fatalf("write run: %v", err)
	}

	queryPayload, err := CaptureDepth(runDir)
	if err != nil {
		t.Fatalf("query capture depth: %v", err)
	}
	var summary CaptureDepthSummary
	if err := json.Unmarshal(queryPayload, &summary); err != nil {
		t.Fatalf("decode capture depth: %v", err)
	}
	if summary.Query != QueryName || summary.TopLevelAssessment != "not_emitted_for_query" {
		t.Fatalf("summary identity = %+v", summary)
	}
	if !summary.UnverifiedTaskExpanded || summary.TaskSupersessionCount != 2 {
		t.Fatalf("task facts = %+v", summary)
	}
	if len(summary.MissingAdapterEvents) != 1 || summary.MissingAdapterEvents[0] != "tool_call" {
		t.Fatalf("missing events = %+v", summary.MissingAdapterEvents)
	}
	if len(summary.UnsupportedObservers) != 1 || summary.UnsupportedObservers[0] != "tool_call" {
		t.Fatalf("unsupported observers = %+v", summary.UnsupportedObservers)
	}
	if len(summary.UnverifiedClaims) == 0 {
		t.Fatalf("expected unverified claims in summary")
	}
}

func TestCaptureDepthDoesNotEchoProviderRefs(t *testing.T) {
	run := adaptercapture.ValidTestInput().Run
	run.ProviderRefs = []adaptercapture.ProviderRef{{
		ReviewRef: "https://review.invalid/7?token=secret-token",
	}}

	runDir := t.TempDir()
	payload, err := json.Marshal(run)
	if err != nil {
		t.Fatalf("marshal run: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "run.json"), payload, 0o644); err != nil {
		t.Fatalf("write run: %v", err)
	}

	queryPayload, err := CaptureDepth(runDir)
	if err != nil {
		t.Fatalf("query capture depth: %v", err)
	}
	if strings.Contains(string(queryPayload), "secret-token") ||
		strings.Contains(string(queryPayload), "review.invalid") {
		t.Fatalf("capture-depth output echoed provider refs: %s", string(queryPayload))
	}
}
