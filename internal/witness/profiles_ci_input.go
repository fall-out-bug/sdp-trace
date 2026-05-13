package witness

import (
	"strings"
)

func applyCIEnvelopeInputState(record *Record, kind, runsRoot, envelopePath string) bool {
	if strings.TrimSpace(envelopePath) == "" {
		// Ambient CI variables without a portable envelope are observation only;
		// they cannot establish a replayable provider witness.
		applyCIMissingEnvelopeState(record, kind)
		return false
	}
	envelope, ok := loadSafeCIEnvelopeRecord(record, envelopePath)
	if !ok {
		return false
	}
	applyCIEnvelopeRecordValues(record, kind, envelope)
	return applyCIEnvelopeTrustDecision(record, kind, runsRoot, envelope)
}
