package checkpoint

import "errors"

func validateSequenceLink(sequence int, previousDigest string) error {
	// validateSequenceLink keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if sequence < 0 {
		return errors.New("checkpoint sequence must be >= 0")
	}

	return validatePreviousDigestForSequence(sequence, previousDigest)
}

func validatePreviousDigestForSequence(sequence int, previousDigest string) error {
	// validatePreviousDigestForSequence keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if sequence == 0 {

		return validateGenesisPreviousDigest(previousDigest)
	}
	if previousDigest == "" {

		return errors.New("sequence > 0 checkpoint requires previous_checkpoint_digest")
	}
	return nil
}

func validateGenesisPreviousDigest(previousDigest string) error {
	// validateGenesisPreviousDigest keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if previousDigest != "" {

		return errors.New("sequence 0 checkpoint must not declare previous_checkpoint_digest")
	}
	return nil
}
