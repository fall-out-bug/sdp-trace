package witness

import (
	"time"
)

func baseRecord(kind string) Record {
	// Base records start as cannot_verify/local_observed until profile-specific
	// evidence raises or fails the established trust scope.
	// Empty artifact slices are intentional; nil would make generated records
	// shape-shift across profiles.
	// OutputSafety starts optimistic but is recalculated before write.
	// GeneratedAt is created locally and does not claim external witness time.
	return Record{
		SchemaVersion:         "sdp-trace-witness-profile-result/v1",
		Kind:                  kind,
		ProfileID:             kind + "-v1",
		ProfileVersion:        "1.0",
		ProviderKind:          kind,
		Status:                StatusCannotVerify,
		TrustScope:            TrustScopeLocalObserved,
		RequestedTrustScope:   TrustScopeCIWitnessed,
		EstablishedTrustScope: stateCannotVerify,
		GeneratedAt:           time.Now().UTC().Format(time.RFC3339Nano),
		RunArtifacts:          []ArtifactDigest{},
		ReportArtifacts:       []ArtifactDigest{},
		OutputSafety:          passingOutputSafety(),
	}
}

func passingOutputSafety() *OutputSafety {
	// The pass state means the known unsafe classes were checked for absence; it
	// does not make any claim about the profile's underlying trust verdict.
	return &OutputSafety{
		State:                 statePass,
		VerifiedAbsentClasses: safetyClasses,
	}
}
