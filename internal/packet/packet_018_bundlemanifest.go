package packet

type BundleManifest struct {
	SchemaVersion string          `json:"schema_version"`
	BundleID      string          `json:"bundle_id"`
	PacketDigest  string          `json:"packet_digest,omitempty"`
	Entries       []BundleEntry   `json:"entries"`
	Resolvers     []ResolverEntry `json:"resolvers,omitempty"`
}
