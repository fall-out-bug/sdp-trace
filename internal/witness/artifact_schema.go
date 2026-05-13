package witness

// ArtifactDigest records a path-local SHA-256 binding. The path stays relative
// to the evidence root so the artifact can be moved with its trace packet.
type ArtifactDigest struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
