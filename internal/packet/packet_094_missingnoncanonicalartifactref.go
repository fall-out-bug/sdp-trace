package packet

import (
	"strings"
)

func missingNonCanonicalArtifactRef(projection Projection) bool {
	return !projection.Canonical && strings.TrimSpace(projection.ArtifactRef) == ""
}
