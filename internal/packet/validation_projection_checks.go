package packet

import "strings"

func invalidCanonicalProjection(projection Projection) bool {
	return projection.Canonical && projection.Kind != ProjectionCanonical
}

func missingNonCanonicalArtifactRef(projection Projection) bool {
	return !projection.Canonical && strings.TrimSpace(projection.ArtifactRef) == ""
}
