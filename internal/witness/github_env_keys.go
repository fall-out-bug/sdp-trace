package witness

// githubIdentityEnvKeys names the GitHub Actions fields that bind a witness to
// a specific source revision and CI run.
var githubIdentityEnvKeys = []string{
	"GITHUB_ACTIONS",
	"GITHUB_SHA",
	"GITHUB_RUN_ID",
	"GITHUB_RUN_ATTEMPT",
	"GITHUB_WORKFLOW",
	"GITHUB_JOB",
	"GITHUB_ACTOR",
	"GITHUB_REPOSITORY",
	"GITHUB_REF",
	"GITHUB_SERVER_URL",
}
