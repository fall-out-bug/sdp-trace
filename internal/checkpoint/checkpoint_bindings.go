package checkpoint

import "fmt"

func compareBindings(result *VerificationResult, expected, actual Payload) {
	// compareBindings keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	compareRunBinding(result, expected, actual)
	compareChainBinding(result, expected, actual)
	compareSourceBinding(result, expected, actual)
	compareNonceBinding(result, expected, actual)
}

func compareRunBinding(result *VerificationResult, expected, actual Payload) {
	// compareRunBinding keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if actual.RunID == expected.RunID {
		return
	}

	result.RunBindingState = StateFail
	result.Reasons = append(result.Reasons, fmt.Sprintf("run_id mismatch: expected %s got %s", expected.RunID, actual.RunID))
}

func compareChainBinding(result *VerificationResult, expected, actual Payload) {
	// compareChainBinding keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if actual.EventChainHead == expected.EventChainHead && actual.EventCount == expected.EventCount {
		return
	}

	result.ChainBindingState = StateFail
	result.Reasons = append(result.Reasons, "event chain binding does not match selected run")
}

func compareSourceBinding(result *VerificationResult, expected, actual Payload) {
	// compareSourceBinding keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if sourceBindingMatches(expected, actual) {
		return
	}

	result.SourceBindingState = StateFail
	result.Reasons = append(result.Reasons, "source, task, or contract binding does not match selected run")
}

func sourceBindingMatches(expected, actual Payload) bool {
	// sourceBindingMatches keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	return actual.SourceSnapshotDigest == expected.SourceSnapshotDigest &&
		actual.SourceSnapshotState == expected.SourceSnapshotState &&
		actual.TaskHash == expected.TaskHash &&
		actual.ContractDigest == expected.ContractDigest
}

func compareNonceBinding(result *VerificationResult, expected, actual Payload) {
	// compareNonceBinding keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if actual.RunNonce == expected.RunNonce {
		return
	}

	result.NonceBindingState = StateFail
	result.Reasons = append(result.Reasons, "run nonce binding does not match selected run")
}
