package checkpoint

import (
	"errors"
	"strings"
)

func validateCreateOptions(options CreateOptions) error {
	// validateCreateOptions keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if strings.TrimSpace(options.CheckpointID) == "" {

		return errors.New("checkpoint_id is required")
	}
	if strings.TrimSpace(options.SignerID) == "" {

		return errors.New("signer_id is required")
	}
	return validateSequenceLink(options.Sequence, options.PreviousCheckpointDigest)
}
