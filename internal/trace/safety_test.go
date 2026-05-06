package trace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
)

func TestObservationTaxonomyValuesAreStable(t *testing.T) {
	tests := []struct {
		name   ObservationState
		stable string
	}{
		{ObservationStateUnsupported, "unsupported"},
		{ObservationStateNotIntegrated, "not_integrated"},
		{ObservationStateSuppressed, "suppressed"},
		{ObservationStateMissingTelemetry, "missing_telemetry"},
		{ObservationStateNotAssessed, "not_assessed"},
		{ObservationStateCannotVerify, "cannot_verify"},
		{ObservationStateOfflineDev, "offline_dev"},
	}

	for _, test := range tests {
		if string(test.name) != test.stable {
			t.Fatalf("observation state changed: got %q want %q", test.name, test.stable)
		}
	}
}

func TestObservationBoundaryValuesAreStable(t *testing.T) {
	tests := []struct {
		name   ObservationBoundary
		stable string
	}{
		{ObservationBoundaryProcessWrapper, "process_wrapper"},
		{ObservationBoundaryAdapterSocket, "adapter_socket"},
		{ObservationBoundaryToolWrapper, "tool_wrapper"},
		{ObservationBoundaryVCSPRObserver, "vcs_pr_observer"},
		{ObservationBoundaryCIObserver, "ci_observer"},
		{ObservationBoundaryExternalWitness, "external_witness"},
	}

	for _, test := range tests {
		if string(test.name) != test.stable {
			t.Fatalf("observation boundary changed: got %q want %q", test.name, test.stable)
		}
	}
}

func TestRetentionModeValuesAreStable(t *testing.T) {
	tests := []struct {
		name   RetentionMode
		stable string
	}{
		{RetentionModeDigestOnly, "digest_only"},
		{RetentionModeSanitizedExcerpt, "sanitized_excerpt"},
		{RetentionModeEncryptedRawRef, "encrypted_raw_ref"},
		{RetentionModeExternalArtifactRef, "external_artifact_ref"},
		{RetentionModeNotAssessed, "not_assessed"},
	}

	for _, test := range tests {
		if string(test.name) != test.stable {
			t.Fatalf("retention mode changed: got %q want %q", test.name, test.stable)
		}
	}
}

func TestCommandDescriptorDoesNotRetainRawArgv(t *testing.T) {
	argv := []string{"/usr/bin/git", "commit", "-m", "secret message"}
	descriptor := NewCommandDescriptor(argv)

	if descriptor.Executable != "git" {
		t.Fatalf("executable = %q, want git", descriptor.Executable)
	}
	if descriptor.Argc != len(argv) {
		t.Fatalf("argc = %d, want %d", descriptor.Argc, len(argv))
	}
	if descriptor.Retention.Mode != RetentionModeDigestOnly {
		t.Fatalf("retention mode = %q, want %q", descriptor.Retention.Mode, RetentionModeDigestOnly)
	}

	expectedDigest := sha256HexString(`["/usr/bin/git","commit","-m","secret message"]`)
	if descriptor.ArgvDigest != expectedDigest {
		t.Fatalf("argv digest = %q, want %q", descriptor.ArgvDigest, expectedDigest)
	}

	encoded, err := json.Marshal(descriptor)
	if err != nil {
		t.Fatal(err)
	}
	payload := string(encoded)
	for _, leaked := range []string{"/usr/bin/git", "commit", "-m", "secret message"} {
		if strings.Contains(payload, leaked) {
			t.Fatalf("descriptor leaked raw argv %q in %s", leaked, payload)
		}
	}
}

func TestEmptyCommandDescriptorIsSafeNotInvented(t *testing.T) {
	descriptor := NewCommandDescriptor(nil)

	if descriptor.Executable != "" {
		t.Fatalf("executable = %q, want empty", descriptor.Executable)
	}
	if descriptor.Argc != 0 {
		t.Fatalf("argc = %d, want 0", descriptor.Argc)
	}
	if descriptor.ArgvDigest != "" {
		t.Fatalf("argv digest = %q, want empty", descriptor.ArgvDigest)
	}
	if descriptor.Retention.Mode != RetentionModeDigestOnly {
		t.Fatalf("retention mode = %q, want %q", descriptor.Retention.Mode, RetentionModeDigestOnly)
	}
}

func sha256HexString(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
