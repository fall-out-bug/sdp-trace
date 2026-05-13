package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func safeRetentionModes() []string {
	// Doctor publishes retention modes that preserve replay without raw secret
	// exposure by default.
	return []string{
		string(trace.RetentionModeDigestOnly),
		string(trace.RetentionModeSanitizedExcerpt),
		string(trace.RetentionModeEncryptedRawRef),
		string(trace.RetentionModeExternalArtifactRef),
		string(trace.RetentionModeNotAssessed),
	}
}
