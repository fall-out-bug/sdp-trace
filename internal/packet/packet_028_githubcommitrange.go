package packet

type GitHubCommitRange struct {
	Base            string `json:"base"`
	Head            string `json:"head"`
	ChangedFilesRef string `json:"changed_files_ref,omitempty"`
}
