package checkpoint

import (
	"errors"
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func validateEnvelope(checkpoint SignedCheckpoint) error {
	// validateEnvelope keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	return firstError(
		validateEnvelopeField(checkpoint.SchemaVersion == CheckpointSchemaVersion, "unsupported checkpoint schema_version %s", checkpoint.SchemaVersion),
		validateEnvelopeField(checkpoint.Profile == ProfileEd25519Detached, "unsupported checkpoint profile %s", checkpoint.Profile),
		validateEnvelopeField(checkpoint.HashAlgorithm == HashAlgorithmSHA256, "unsupported checkpoint hash_algorithm %s", checkpoint.HashAlgorithm),
		validateCanonicalization(checkpoint.Canonical),
	)
}

func validateEnvelopeField(ok bool, format, value string) error {
	// validateEnvelopeField keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if ok {
		return nil
	}

	return fmt.Errorf(format, value)
}

func validateCanonicalization(canonical trace.Canonicalization) error {
	// validateCanonicalization keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if canonical.Algorithm == trace.CanonicalSchemaAlgo && canonical.Version == trace.CanonicalAlgoVersion {
		return nil
	}

	return errors.New("unsupported checkpoint canonicalization")
}

func firstError(errs ...error) error {
	// firstError keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	for _, err := range errs {
		if err != nil {

			return err
		}
	}
	return nil
}
