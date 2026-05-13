package main

type githubActionsArtifact struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Expired   bool   `json:"expired"`
	ExpiresAt string `json:"expires_at"`
	URL       string `json:"archive_download_url"`
}
