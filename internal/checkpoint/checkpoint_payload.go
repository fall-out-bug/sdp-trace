package checkpoint

import (
	"errors"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func BuildPayload(runDir, previousCheckpointDigest string) (Payload, error) {
	// BuildPayload keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	artifact, err := validatedRunArtifact(runDir)
	if err != nil {
		return Payload{}, err
	}
	nonce := runNonce(artifact.Events)
	if nonce == "" {

		return Payload{}, errors.New("run nonce missing from recorder_attached event")
	}
	return payloadFromArtifact(artifact, nonce, previousCheckpointDigest), nil
}

func payloadFromArtifact(artifact trace.RunArtifact, nonce, previousCheckpointDigest string) Payload {
	// payloadFromArtifact keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	manifest := artifact.Manifest

	return Payload{
		RunID:                    manifest.RunID,
		RunNonce:                 nonce,
		EventChainHead:           manifestChainHead(manifest),
		EventCount:               manifest.EventCount,
		SourceSnapshotDigest:     manifest.SourceSnapshot,
		SourceSnapshotState:      manifest.SourceState,
		TaskHash:                 trace.SHA256Hex(manifest.Task),
		ContractDigest:           manifest.ContractDigest,
		PreviousCheckpointDigest: previousCheckpointDigest,
		ReplayContext:            notAssessedReplayContext(),
	}
}

func notAssessedReplayContext() ReplayContext {
	// notAssessedReplayContext keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	return ReplayContext{
		Repository: "not_assessed",
		Ref:        "not_assessed",
		CommitSHA:  "not_assessed",
	}
}

func manifestChainHead(manifest trace.RunManifest) string {
	// manifestChainHead keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if manifest.EventChainHead != "" {
		return manifest.EventChainHead
	}

	return manifest.FinalChainHead
}
