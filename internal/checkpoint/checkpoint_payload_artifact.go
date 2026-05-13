package checkpoint

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func validatedRunArtifact(runDir string) (trace.RunArtifact, error) {
	// validatedRunArtifact keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	artifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {
		return trace.RunArtifact{}, err
	}
	if err := artifact.Manifest.Validate(); err != nil {
		return trace.RunArtifact{}, err
	}
	if artifact.Manifest.EventCount != len(artifact.Events) {

		return trace.RunArtifact{}, fmt.Errorf("event_count mismatch: run.json=%d files=%d", artifact.Manifest.EventCount, len(artifact.Events))
	}
	return artifact, trace.ValidateEventChain(artifact.Events)
}
