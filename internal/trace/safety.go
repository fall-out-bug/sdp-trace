package trace

import (
	"encoding/json"
	"path/filepath"
)

// ObservationState classifies why an expected observation is or is not present.
type ObservationState string

// ObservationState constants are stable machine-enumerable values.
const (
	ObservationStateUnsupported      ObservationState = "unsupported"
	ObservationStateNotIntegrated    ObservationState = "not_integrated"
	ObservationStateSuppressed       ObservationState = "suppressed"
	ObservationStateMissingTelemetry ObservationState = "missing_telemetry"
	ObservationStateNotAssessed      ObservationState = "not_assessed"
	ObservationStateCannotVerify     ObservationState = "cannot_verify"
	ObservationStateOfflineDev       ObservationState = "offline_dev"
)

// ObservationBoundary names the portable trust boundary for an observation.
type ObservationBoundary string

// ObservationBoundary constants are stable machine-enumerable values.
const (
	ObservationBoundaryProcessWrapper  ObservationBoundary = "process_wrapper"
	ObservationBoundaryAdapterSocket   ObservationBoundary = "adapter_socket"
	ObservationBoundaryToolWrapper     ObservationBoundary = "tool_wrapper"
	ObservationBoundaryVCSPRObserver   ObservationBoundary = "vcs_pr_observer"
	ObservationBoundaryCIObserver      ObservationBoundary = "ci_observer"
	ObservationBoundaryExternalWitness ObservationBoundary = "external_witness"
)

// RetentionMode describes how retained observation material is stored.
type RetentionMode string

// RetentionMode constants are stable machine-enumerable values.
const (
	RetentionModeDigestOnly          RetentionMode = "digest_only"
	RetentionModeSanitizedExcerpt    RetentionMode = "sanitized_excerpt"
	RetentionModeEncryptedRawRef     RetentionMode = "encrypted_raw_ref"
	RetentionModeExternalArtifactRef RetentionMode = "external_artifact_ref"
	RetentionModeNotAssessed         RetentionMode = "not_assessed"
)

// RetentionDescriptor records safe, machine-readable retention metadata.
type RetentionDescriptor struct {
	Mode        RetentionMode `json:"mode"`
	Description string        `json:"description,omitempty"`
	Digest      string        `json:"digest,omitempty"`
	Reference   string        `json:"reference,omitempty"`
}

// CommandDescriptor identifies a command without retaining raw argv.
type CommandDescriptor struct {
	Executable string              `json:"executable"`
	Argc       int                 `json:"argc"`
	ArgvDigest string              `json:"argv_digest,omitempty"`
	Retention  RetentionDescriptor `json:"retention"`
}

// NewCommandDescriptor returns a descriptor safe for payload retention. The
// raw argv is represented only by argument count and a deterministic digest.
func NewCommandDescriptor(argv []string) CommandDescriptor {
	descriptor := CommandDescriptor{
		Argc: len(argv),
		Retention: RetentionDescriptor{
			Mode:        RetentionModeDigestOnly,
			Description: "argv_digest stores a sha256 digest of the command argv; raw argv is not retained",
		},
	}
	if len(argv) == 0 {
		return descriptor
	}
	descriptor.Executable = filepath.Base(argv[0])
	descriptor.ArgvDigest = SHA256Hex(mustMarshalJSONStringSlice(argv))
	return descriptor
}

func mustMarshalJSONStringSlice(values []string) string {
	encoded, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}
	return string(encoded)
}
