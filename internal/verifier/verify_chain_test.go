package verifier

import (
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func TestVerifyChainPreservesBehavioralContracts(t *testing.T) {
	base := []trace.Event{
		makeTestEvent(t, 0, trace.NullEventHash, "evt-0"),
	}
	base = append(base, makeTestEvent(t, 1, base[0].EventHash, "evt-1"))
	base = append(base, makeTestEvent(t, 2, base[1].EventHash, "evt-2"))

	t.Run("valid_chain", func(t *testing.T) {
		ok, issue := verifyChain(base, base[2].EventHash)
		if !ok {
			t.Fatalf("expected valid chain, got issue=%q", issue)
		}
	})

	t.Run("valid_chain_with_empty_expected_head", func(t *testing.T) {
		ok, issue := verifyChain(base, "")
		if !ok {
			t.Fatalf("expected valid chain, got issue=%q", issue)
		}
	})

	t.Run("sequence_mismatch", func(t *testing.T) {
		events := copyEvents(base)
		events[1].Sequence = 99
		ok, issue := verifyChain(events, "")
		if ok {
			t.Fatalf("expected invalid sequence, got ok=true")
		}
		expect := "sequence mismatch at evt-1"
		if issue != expect {
			t.Fatalf("expected issue %q, got %q", expect, issue)
		}
	})

	t.Run("first_event_prev_hash", func(t *testing.T) {
		events := copyEvents(base)
		events[0].PrevEventHash = "bad"
		ok, issue := verifyChain(events, "")
		if ok {
			t.Fatalf("expected invalid first-event prev hash, got ok=true")
		}
		expect := "first event has non-empty prev_event_hash"
		if issue != expect {
			t.Fatalf("expected issue %q, got %q", expect, issue)
		}
	})

	t.Run("broken_chain", func(t *testing.T) {
		events := copyEvents(base)
		events[1].PrevEventHash = "bad"
		ok, issue := verifyChain(events, "")
		if ok {
			t.Fatalf("expected broken chain, got ok=true")
		}
		expect := "broken chain at 2 (evt-1)"
		if issue != expect {
			t.Fatalf("expected issue %q, got %q", expect, issue)
		}
	})

	t.Run("invalid_payload_digest", func(t *testing.T) {
		events := copyEvents(base)
		events[2].PayloadDigest = "sha256:invalid"
		ok, issue := verifyChain(events, "")
		if ok {
			t.Fatalf("expected invalid payload digest, got ok=true")
		}
		expectPrefix := "invalid payload digest for evt-2: "
		if !strings.HasPrefix(issue, expectPrefix) {
			t.Fatalf("expected issue with prefix %q, got %q", expectPrefix, issue)
		}
	})

	t.Run("hash_mismatch", func(t *testing.T) {
		events := copyEvents(base)
		events[2].EventHash = "sha256:bad"
		ok, issue := verifyChain(events, "")
		if ok {
			t.Fatalf("expected hash mismatch, got ok=true")
		}
		expect := "hash mismatch for evt-2"
		if issue != expect {
			t.Fatalf("expected issue %q, got %q", expect, issue)
		}
	})

	t.Run("expected_head_mismatch", func(t *testing.T) {
		ok, issue := verifyChain(base, "sha256:bad")
		if ok {
			t.Fatalf("expected head mismatch, got ok=true")
		}
		expect := "run head does not match manifest final_chain_head"
		if issue != expect {
			t.Fatalf("expected issue %q, got %q", expect, issue)
		}
	})
}

func makeTestEvent(t *testing.T, sequence int, prevHash, eventID string) trace.Event {
	t.Helper()
	event, err := trace.Event{
		SchemaVersion: trace.SchemaVersion,
		RunID:         "chain-test",
		EventID:       eventID,
		Sequence:      sequence,
		EventType:     trace.EventType("unit_test"),
		Timestamp:     "2026-05-10T00:00:00Z",
		PrevEventHash: prevHash,
		HashAlgorithm: trace.HashAlgSHA256,
		EventPayload:  map[string]any{"value": eventID},
	}.WithComputedEventHash()
	if err != nil {
		t.Fatalf("build test event: %v", err)
	}
	return event
}

func copyEvents(events []trace.Event) []trace.Event {
	cp := make([]trace.Event, len(events))
	copy(cp, events)
	return cp
}
