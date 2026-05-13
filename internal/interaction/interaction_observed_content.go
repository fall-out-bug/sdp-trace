package interaction

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func applyObservedEventContent(event *Event, taskID, id string, body []byte) {
	// applyObservedEventContent keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	event.ContentRef = interactionContentRef(taskID, id)
	event.ContentDigest = interactionSHA256Hex(body)
	event.DigestAlgorithm = DigestAlgorithmSHA256
}

func applyObservedEventAssessment(event *Event, now string) {
	// applyObservedEventAssessment keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	event.Retention = RetentionSanitizedExcerpt
	event.State = StateUnreferenced
	event.ObservedBeforeDelivery = true
	event.ChannelExclusivity = ChannelExclusivityNotAssessed
	event.CompletenessState = CompletenessComplete
	event.Redaction = defaultRedaction()

	event.ObservedAt = now
	event.CreatedAt = now
}

func observedRelaySource() Source {

	return Source{SourceType: SourceObservedControlChannel, SourceID: DefaultRelaySourceID, SourceVersion: SchemaVersion}
}

func interactionContentRef(taskID, interactionID string) string {

	return fmt.Sprintf("sdp://interaction/%s/%s", taskID, interactionID)
}

func defaultRedaction() Redaction {

	return Redaction{PolicyRef: DefaultRedactionPolicyRef, FindingCount: 0}
}

func interactionSHA256Hex(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}
