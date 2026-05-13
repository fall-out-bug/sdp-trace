package repoobserver

// Config is the local observer manifest written by install mode and read by
// doctor mode for repository identity.
type Config struct {
	SchemaVersion   string            `json:"schema_version"`
	Profile         string            `json:"profile"`
	RepositoryID    string            `json:"repository_id"`
	TrustBoundary   string            `json:"trust_boundary"`
	InstalledFiles  []string          `json:"installed_files"`
	InstallMetadata map[string]string `json:"install_metadata"`
}
