package interaction

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRelayWritesBeforeForwardingAndSummarizesFriction(t *testing.T) {
	dir := t.TempDir()
	forwardPath := filepath.Join(dir, "forwarded.txt")
	tracePath := filepath.Join(dir, "trace.json")
	forward := shellCommand("cat > " + shellQuote(forwardPath))
	trace, exitCode, err := Relay(context.Background(), RelayOptions{
		TaskID:      "task-1",
		ActorType:   "human_user",
		ActorID:     "human",
		Target:      "gsd",
		EventType:   "corrective_feedback",
		OperationID: "op-1",
		StageID:     "plan",
		Out:         tracePath,
		Command:     forward,
		Now:         time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
	}, strings.NewReader("Do not edit the plan before approval.\n"), &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if exitCode != 0 {
		t.Fatalf("exit = %d", exitCode)
	}
	if _, err := os.Stat(tracePath); err != nil {
		t.Fatalf("trace not written: %v", err)
	}
	body, err := os.ReadFile(forwardPath)
	if err != nil {
		t.Fatalf("forwarded content missing: %v", err)
	}
	if !strings.Contains(string(body), "Do not edit") {
		t.Fatalf("forwarded body = %q", string(body))
	}
	if _, err := os.Stat(contentBlobPath(tracePath, trace.Events[0])); err != nil {
		t.Fatalf("content blob missing: %v", err)
	}
	if len(trace.Events) != 1 {
		t.Fatalf("events = %d", len(trace.Events))
	}
	event := trace.Events[0]
	if !event.ObservedBeforeDelivery || event.ChannelExclusivity != StateNotAssessed {
		t.Fatalf("observed boundary = %+v", event)
	}
	if event.Retention != RetentionSanitizedExcerpt || event.ContentRef == "" || event.ContentDigest == "" {
		t.Fatalf("retention/digest missing: %+v", event)
	}
	summary := SummarizeTrace(trace)
	if summary.CorrectionAfterTask != 0 || summary.FrictionCounts["correction"] != 1 || len(summary.NotAssessed) == 0 {
		t.Fatalf("summary = %+v", summary)
	}
}

func TestRelayRefusesUnsafeContentBeforeForwarding(t *testing.T) {
	dir := t.TempDir()
	forwardPath := filepath.Join(dir, "forwarded.txt")
	tracePath := filepath.Join(dir, "trace.json")
	forward := shellCommand("cat > " + shellQuote(forwardPath))
	_, _, err := Relay(context.Background(), RelayOptions{
		TaskID:    "task-1",
		ActorType: "human_user",
		ActorID:   "human",
		Target:    "gsd",
		EventType: "corrective_feedback",
		Out:       tracePath,
		Command:   forward,
	}, strings.NewReader("api_key=SECRET_TOKEN_SHOULD_NOT_APPEAR\n"), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "unsafe") {
		t.Fatalf("expected unsafe content error, got %v", err)
	}
	if _, err := os.Stat(tracePath); !os.IsNotExist(err) {
		t.Fatalf("trace should not be written, stat=%v", err)
	}
	if _, err := os.Stat(forwardPath); !os.IsNotExist(err) {
		t.Fatalf("forward command should not run, stat=%v", err)
	}
}

func TestImportTranscriptRejectsAgentReportedAndNonMonotonicOrdering(t *testing.T) {
	dir := t.TempDir()
	eventsPath := filepath.Join(dir, "events.jsonl")
	event := validImportedEvent("ix-1", 2)
	event.Source.SourceType = SourceAgentReported
	writeJSONL(t, eventsPath, event)
	_, err := ImportTranscript(ImportOptions{
		TaskID:      "task-1",
		Source:      SourcePreclassifiedTranscript,
		SourceRef:   "export-1",
		EventsJSONL: eventsPath,
		Out:         filepath.Join(dir, "trace.json"),
	})
	if err == nil || !strings.Contains(err.Error(), "agent-reported") {
		t.Fatalf("expected agent reported rejection, got %v", err)
	}

	event = validImportedEvent("ix-1", 2)
	event2 := validImportedEvent("ix-2", 1)
	writeJSONL(t, eventsPath, event, event2)
	_, err = ImportTranscript(ImportOptions{
		TaskID:      "task-1",
		Source:      SourcePreclassifiedTranscript,
		SourceRef:   "export-1",
		EventsJSONL: eventsPath,
		Out:         filepath.Join(dir, "trace.json"),
	})
	if err == nil || !strings.Contains(err.Error(), "non-monotonic") {
		t.Fatalf("expected ordering rejection, got %v", err)
	}
}

func TestImportTranscriptSuccessSummarizesCatalogMetrics(t *testing.T) {
	dir := t.TempDir()
	eventsPath := filepath.Join(dir, "events.jsonl")
	assignment := validImportedEvent("ix-0", 0)
	assignment.EventType = "task_assignment"
	assignment.FrictionClass = "none"
	clarification := validImportedEvent("ix-1", 1)
	clarification.EventType = "clarification_request"
	clarification.FrictionClass = "clarification"
	rejection := validImportedEvent("ix-2", 2)
	rejection.EventType = "plan_rejected"
	rejection.FrictionClass = "correction"
	writeJSONL(t, eventsPath, assignment, clarification, rejection)

	trace, err := ImportTranscript(ImportOptions{
		TaskID:      "task-1",
		Source:      SourcePreclassifiedTranscript,
		SourceRef:   "export-1",
		EventsJSONL: eventsPath,
		Out:         filepath.Join(dir, "trace.json"),
	})
	if err != nil {
		t.Fatal(err)
	}
	summary := SummarizeTrace(trace)
	if summary.PlanRejectionCount != 1 || summary.ClarificationTurnCount != 1 || summary.CorrectionAfterTask != 0 {
		t.Fatalf("summary = %+v", summary)
	}
}

func TestEnvelopeSummaryCountsRefsAndRejectsUnsafeRunRef(t *testing.T) {
	envelope := Envelope{
		SchemaVersion:   SchemaVersion,
		TaskID:          "task-1",
		EnvelopeID:      "env-1",
		RunRefs:         []string{"recorder:run-1"},
		SourceRefs:      []string{"evidence:source-1"},
		TaskRefs:        []string{"evidence:task-1"},
		PromiseRefs:     []string{"evidence:promise-1"},
		InteractionRefs: []string{"sdp://interaction/task-1/ix-1"},
		OperationRefs:   []string{"recorder:run-1/event:1"},
		LLMRefs:         []string{"recorder:run-1/event:2"},
		ToolRefs:        []string{"recorder:run-1/event:3"},
		MutationRefs:    []string{"evidence:mutation-1"},
		StageRefs:       []string{"evidence:stage-1"},
		FrictionRefs:    []string{"sdp://interaction/task-1/ix-1"},
		AssessmentState: "partial",
		NotAssessed:     []string{"gateway linkage absent"},
		CreatedAt:       "2026-05-09T10:00:00Z",
		UpdatedAt:       "2026-05-09T10:00:00Z",
	}
	if err := ValidateEnvelope(envelope); err != nil {
		t.Fatal(err)
	}
	summary := SummarizeEnvelope(envelope)
	if summary.RunRefCount != 1 || summary.SourceRefCount != 1 || summary.TaskRefCount != 1 || summary.PromiseRefCount != 1 ||
		summary.InteractionRefCount != 1 || summary.OperationRefCount != 1 || summary.LLMRefCount != 1 || summary.ToolRefCount != 1 ||
		summary.MutationRefCount != 1 || summary.StageRefCount != 1 || summary.FrictionRefCount != 1 {
		t.Fatalf("summary = %+v", summary)
	}
	envelope.RunRefs = []string{"/tmp/run"}
	if err := ValidateEnvelope(envelope); err == nil {
		t.Fatalf("expected unsafe run ref rejection")
	}
}

func TestValidateEventRejectsInvalidFields(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*Event)
		wantErr string
	}{
		{
			name: "friction class mismatch",
			mutate: func(event *Event) {
				event.FrictionClass = "none"
			},
			wantErr: "does not match event_type",
		},
		{
			name: "missing retained content reason",
			mutate: func(event *Event) {
				event.ContentRef = ""
				event.NotRetainedReason = ""
			},
			wantErr: "without content_ref requires not_retained_reason",
		},
		{
			name: "unsupported source type",
			mutate: func(event *Event) {
				event.Source.SourceType = "agent-memory"
			},
			wantErr: "unsupported source_type",
		},
		{
			name: "unsafe source id",
			mutate: func(event *Event) {
				event.SourceID = "../transcript"
			},
			wantErr: "source_id must match",
		},
		{
			name: "invalid reference ref",
			mutate: func(event *Event) {
				event.ReferenceRefs = []string{"/tmp/local-proof"}
			},
			wantErr: "unsupported reference_ref",
		},
		{
			name: "invalid llm linkage state",
			mutate: func(event *Event) {
				event.LLMRefs = []LLMRef{{LinkageState: "passed"}}
			},
			wantErr: "unsupported llm linkage_state",
		},
		{
			name: "negative source sequence",
			mutate: func(event *Event) {
				event.SourceSequence = -1
			},
			wantErr: "source_sequence must be non-negative",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			event := validImportedEvent("ix-1", 1)
			tc.mutate(&event)
			err := ValidateEvent(event)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("ValidateEvent() error = %v, want %q", err, tc.wantErr)
			}
		})
	}
}

func validImportedEvent(id string, sequence int) Event {
	body := []byte("Please fix stale evidence.\n")
	sum := sha256Hex(body)
	return Event{
		SchemaVersion:          SchemaVersion,
		InteractionID:          id,
		TaskID:                 "task-1",
		SourceID:               "transcript-1",
		SourceSequence:         sequence,
		EventType:              "evidence_correction",
		FrictionClass:          "evidence",
		Actor:                  Actor{ID: "human", ActorType: "human_user"},
		Target:                 "agent",
		Source:                 Source{SourceType: SourcePreclassifiedTranscript, SourceID: "transcript-1", SourceRef: "export-1"},
		ContentRef:             "external:export-1",
		ContentDigest:          sum,
		DigestAlgorithm:        DigestAlgorithmSHA256,
		Retention:              RetentionExternalArtifactRef,
		State:                  StateUnreferenced,
		ObservedBeforeDelivery: false,
		ChannelExclusivity:     StateNotAssessed,
		CompletenessState:      CompletenessPartial,
		Redaction:              Redaction{PolicyRef: DefaultRedactionPolicyRef},
		ObservedAt:             "2026-05-09T10:00:00Z",
		CreatedAt:              "2026-05-09T10:00:00Z",
	}
}

func writeJSONL(t *testing.T, path string, events ...Event) {
	t.Helper()
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, event := range events {
		if err := enc.Encode(event); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

func sha256Hex(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func shellCommand(script string) []string {
	return []string{"sh", "-c", script}
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
