package packet

type Projection struct {
	Kind        string `json:"kind"`
	Canonical   bool   `json:"canonical"`
	ArtifactRef string `json:"artifact_ref,omitempty"`
}
