package witness

const (
	// KindGitHubActions identifies a witness built from GitHub Actions identity
	// and optional OIDC evidence.
	KindGitHubActions = "github-actions"
	// KindGitLabCI and KindBuildkite are reserved provider labels for portable
	// schema compatibility before verifier support exists.
	KindGitLabCI  = "gitlab-ci"
	KindBuildkite = "buildkite"
	// KindCustomerPKI identifies customer-managed signing authority evidence.
	KindCustomerPKI = "customer-pki"
)
