package witness

import (
	"strings"
)

func applyCIEnvelopeRecordValues(record *Record, kind string, envelope EnvelopeInput) {
	// Envelope metadata is copied only after safe parsing. Empty metadata keeps
	// deterministic defaults so missing fields do not become implicit claims.
	record.SchemaVersion = defaultString(envelope.SchemaVersion, "sdp-trace-witness-profile-result/v1")
	record.ProfileID = defaultString(envelope.ProfileID, kind+"-v1")
	record.ProfileVersion = defaultString(envelope.ProfileVersion, "1.0")
	record.ProviderKind = defaultString(envelope.ProviderKind, kind)
	record.RequestedTrustScope = defaultString(envelope.RequestedTrustScope, TrustScopeCIWitnessed)
	record.Source = envelope.Source
	record.CI = envelope.CI
	record.ProfileStates = &envelope.ProfileStates
	if strings.TrimSpace(envelope.GeneratedAt) != "" {
		record.GeneratedAt = envelope.GeneratedAt
	}
}
