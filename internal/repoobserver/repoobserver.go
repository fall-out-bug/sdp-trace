package repoobserver

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	ProfileGithubActionsGitHooksV1 = "github-actions-git-hooks-v1"
	SchemaVersion                  = "block28-repo-observer-status-v1"

	StatePass         = "pass"
	StateFail         = "fail"
	StateNotAssessed  = "not_assessed"
	StateCannotVerify = "cannot_verify"

	ScopeLocalStructural = "local_structural"
	ScopeCIUploaded      = "ci_uploaded"
	ScopeExternalWitness = "external_witnessed"
	ScopeAgentReported   = "agent_reported"
	ScopeNotApplicable   = "not_applicable"

	ReasonHooksPathAbsent              = "hooks_path_absent"
	ReasonHooksPathSet                 = "hooks_path_set"
	ReasonHooksPathMismatch            = "hooks_path_mismatch"
	ReasonHookScriptAbsent             = "hook_script_absent"
	ReasonHookScriptPresent            = "hook_script_present"
	ReasonHookOutputNotObserved        = "hook_output_not_observed"
	ReasonLocalHooksBypassable         = "local_hooks_bypassable"
	ReasonAlreadyInstalled             = "already_installed"
	ReasonCIWorkflowAbsent             = "ci_workflow_absent"
	ReasonCIWorkflowPresent            = "ci_workflow_present"
	ReasonCIArtifactUploadAbsent       = "ci_artifact_upload_absent"
	ReasonCIArtifactUploadPresent      = "ci_artifact_upload_present"
	ReasonCIArtifactBundleNotObserved  = "ci_artifact_bundle_not_observed"
	ReasonCIArtifactBundleObserved     = "ci_artifact_bundle_observed"
	ReasonAgentReportedNotProof        = "agent_reported_not_proof"
	ReasonOutsideProfileScope          = "outside_profile_scope"
	ReasonUnsafeOutputRefused          = "unsafe_output_refused"
	ReasonManualStepRequired           = "manual_step_required"
	SurfaceHooksPath                   = "git_hooks_path"
	SurfacePreCommitHook               = "git_hook_pre_commit"
	SurfacePostCommitHook              = "git_hook_post_commit"
	SurfacePrePushHook                 = "git_hook_pre_push"
	SurfaceCIWorkflow                  = "github_actions_workflow"
	SurfaceCIArtifactUpload            = "github_actions_artifact_upload"
	SurfaceCIArtifactBundleObservation = "github_actions_artifact_bundle"
	SurfacePRCheckBinding              = "pr_check_binding"
	SurfaceLocalWrappedCommands        = "local_wrapped_commands"
	SurfaceAgentPrompt                 = "agent_prompt"
	SurfaceRepositoryIdentity          = "repository_identity"
	SurfaceConfig                      = "sdp_trace_config"
	SurfaceGitignore                   = "sdp_trace_gitignore"
	gitignoreBeginMarker               = "# sdp-trace begin"
	gitignoreEndMarker                 = "# sdp-trace end"
	gitignoreBlock                     = "# sdp-trace begin\n.sdp-trace/hooks/\n.sdp-trace/ci/\n.sdp-trace/install-diff.txt\n# sdp-trace end\n"
)

var safeIDPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

type Options struct {
	RepoRoot     string
	Profile      string
	RepositoryID string
	Write        bool
	Force        bool
	Now          time.Time
}

type Config struct {
	SchemaVersion   string            `json:"schema_version"`
	Profile         string            `json:"profile"`
	RepositoryID    string            `json:"repository_id"`
	TrustBoundary   string            `json:"trust_boundary"`
	InstalledFiles  []string          `json:"installed_files"`
	InstallMetadata map[string]string `json:"install_metadata"`
}

type Status struct {
	SchemaVersion     string        `json:"schema_version"`
	Profile           string        `json:"profile"`
	RepositoryID      string        `json:"repository_id"`
	RepositoryRootRef string        `json:"repository_root_ref"`
	GitHead           string        `json:"git_head"`
	GitBranch         string        `json:"git_branch"`
	InstallState      string        `json:"install_state"`
	ProofState        string        `json:"proof_state"`
	Surfaces          []Surface     `json:"surfaces"`
	Gaps              []Gap         `json:"gaps"`
	NextActions       []NextAction  `json:"next_actions"`
	ForceDiffSummary  []DiffSummary `json:"force_diff_summary,omitempty"`
	GeneratedAt       string        `json:"generated_at"`
}

type Surface struct {
	SurfaceID      string `json:"surface_id"`
	InstallState   string `json:"install_state"`
	ProofState     string `json:"proof_state"`
	TrustScope     string `json:"trust_scope"`
	EvidenceSource string `json:"evidence_source"`
	ObservedRef    string `json:"observed_ref,omitempty"`
	ReasonCode     string `json:"reason_code"`
	NextAction     string `json:"next_action,omitempty"`
}

type Gap struct {
	SurfaceID  string `json:"surface_id"`
	ReasonCode string `json:"reason_code"`
	Detail     string `json:"detail"`
}

type NextAction struct {
	SurfaceID  string `json:"surface_id"`
	ActionText string `json:"action_text"`
	Blocking   bool   `json:"blocking"`
}

type DiffSummary struct {
	Path    string `json:"path"`
	Action  string `json:"action"`
	Before  string `json:"before,omitempty"`
	After   string `json:"after,omitempty"`
	Summary string `json:"summary"`
	Backup  string `json:"backup,omitempty"`
}

type targetFile struct {
	path       string
	content    string
	executable bool
}

func Doctor(opts Options) (Status, error) {
	// Doctor is read-only: it reports local structural setup without installing
	// generated observer files.
	opts, err := normalizeOptions(opts)
	if err != nil {
		return Status{}, err
	}
	opts, err = withConfiguredRepositoryID(opts)
	if err != nil {
		return Status{}, err
	}
	return buildStatus(opts, false)
}

func Install(opts Options) (Status, error) {
	// Install builds a preview first so --write callers can still receive the
	// same surface model after mutations are applied.
	opts, err := normalizeOptions(opts)
	if err != nil {
		return Status{}, err
	}
	status, err := buildStatus(opts, true)
	if err != nil {
		return status, err
	}
	if !opts.Write {
		return status, nil
	}
	return installWriteMode(opts, status)
}

func installWriteMode(opts Options, status Status) (Status, error) {
	// Preserve mutation summaries across the post-write rescan so the final
	// status shows both resulting surfaces and what changed.
	summary, err := writeInstallFiles(opts)
	status.ForceDiffSummary = summary
	if err != nil {
		return status, err
	}
	status, err = buildStatus(opts, false)
	status.ForceDiffSummary = summary
	return status, err
}

func WriteJSON(path string, status Status) error {
	// Empty output paths are optional CLI sinks; nonempty paths get stable,
	// newline-terminated JSON for review artifacts.
	// The writer does not re-evaluate status; callers own the observation or
	// install pass that produced it.
	if strings.TrimSpace(path) == "" {
		return nil
	}
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func LoadConfig(repoRoot string) (Config, error) {
	// Config is local structural input; malformed config is unsafe because it
	// could mislabel repository identity or profile.
	data, err := os.ReadFile(filepath.Join(repoRoot, ".sdp-trace", "config.json"))
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("%s: .sdp-trace/config.json is malformed", ReasonUnsafeOutputRefused)
	}
	return validateConfig(config)
}

func HumanTable(status Status) string {
	// Human output keeps install and proof state separate instead of compressing
	// them into a single health score.
	var b strings.Builder
	writeHumanTableHeader(&b, status)
	writeHumanTableSurfaces(&b, status.Surfaces)
	writeHumanTableDiffSummary(&b, status.ForceDiffSummary)
	b.WriteString("\nNote: core.hooksPath is local checkout configuration and is not committed into the repository.\n")
	return b.String()
}

func normalizeOptions(opts Options) (Options, error) {
	// Defaults, absolute paths, and repository identity all pass through one
	// normalization path before any observation or install work.
	opts = withDefaultProfile(opts)
	if err := validateProfile(opts.Profile); err != nil {
		return Options{}, err
	}
	opts, err := withAbsoluteRepoRoot(opts)
	if err != nil {
		return Options{}, err
	}
	opts = withDefaultNow(opts)
	return opts, validateRepositoryID(opts.RepositoryID, "repository id must match [A-Za-z0-9_.-]+")
}

func withDefaultProfile(opts Options) Options {
	if strings.TrimSpace(opts.Profile) == "" {
		// Empty profile selects the only portable repo-observer contract rather
		// than leaving behavior harness-specific.
		opts.Profile = ProfileGithubActionsGitHooksV1
	}
	return opts
}

func validateProfile(profile string) error {
	if profile != ProfileGithubActionsGitHooksV1 {
		// Unknown profiles would imply setup semantics this portable tool cannot
		// verify.
		return fmt.Errorf("repo observer requires --profile %s", ProfileGithubActionsGitHooksV1)
	}
	return nil
}

func withAbsoluteRepoRoot(opts Options) (Options, error) {
	// Blank roots are resolved through git and then made absolute for later
	// containment checks.
	// The absolute path is local structural evidence only and is rendered as an
	// abstract repository root in user-facing status.
	if strings.TrimSpace(opts.RepoRoot) == "" {
		root, err := repoRoot(".")
		if err != nil {
			return Options{}, err
		}
		opts.RepoRoot = root
	}
	abs, err := filepath.Abs(opts.RepoRoot)
	if err != nil {
		return Options{}, err
	}
	opts.RepoRoot = abs
	return opts, nil
}

func withDefaultNow(opts Options) Options {
	if opts.Now.IsZero() {
		// Repo observations use wall-clock UTC in normal runs, while tests inject
		// time for deterministic status output.
		opts.Now = time.Now().UTC()
	}
	return opts
}

func withConfiguredRepositoryID(opts Options) (Options, error) {
	if opts.RepositoryID != "" {
		// Explicit repository identity wins over local config so callers can
		// bind observations to the intended source.
		return opts, nil
	}
	return withConfigFileRepositoryID(opts)
}

func withConfigFileRepositoryID(opts Options) (Options, error) {
	// Missing config is allowed for first install/doctor runs; the caller can
	// still use a derived repository ID.
	config, err := LoadConfig(opts.RepoRoot)
	if errors.Is(err, os.ErrNotExist) {
		return opts, nil
	}
	if err != nil {
		return Options{}, err
	}
	if config.RepositoryID != "" {
		opts.RepositoryID = config.RepositoryID
	}
	return opts, nil
}

func validateConfig(config Config) (Config, error) {
	if config.Profile != "" && config.Profile != ProfileGithubActionsGitHooksV1 {
		// Config files cannot opt into unimplemented observer profiles.
		return Config{}, fmt.Errorf("%s: unsupported repo observer profile in .sdp-trace/config.json", ReasonUnsafeOutputRefused)
	}
	return config, validateRepositoryID(config.RepositoryID, "repository id in .sdp-trace/config.json must match [A-Za-z0-9_.-]+")
}

func validateRepositoryID(repositoryID, message string) error {
	if repositoryID != "" && !safeIDPattern.MatchString(repositoryID) {
		// Repository IDs are rendered in reports and config, so reject unsafe
		// labels at the boundary.
		return fmt.Errorf("%s: %s", ReasonUnsafeOutputRefused, message)
	}
	return nil
}

func writeHumanTableHeader(b *strings.Builder, status Status) {
	// The repository root is represented abstractly; absolute checkout paths are
	// not part of the portable status table.
	fmt.Fprintf(b, "Profile: %s\n", status.Profile)
	fmt.Fprintf(b, "Repository: %s\n", status.RepositoryID)
	fmt.Fprintf(b, "Install state: %s\n", status.InstallState)
	fmt.Fprintf(b, "Proof state: %s\n\n", status.ProofState)
	b.WriteString("Surface | Install state | Proof state | Trust scope | Evidence source | Next action\n")
	b.WriteString("--- | --- | --- | --- | --- | ---\n")
}

func writeHumanTableSurfaces(b *strings.Builder, surfaces []Surface) {
	for _, surface := range surfaces {
		// Render every surface independently so install and proof gaps remain
		// inspectable instead of being collapsed into a health score.
		writeHumanTableSurface(b, surface)
	}
}

func writeHumanTableSurface(b *strings.Builder, surface Surface) {
	// Empty remediation renders as "-" so absence of an action is explicit.
	// The table prints install and proof states side by side to avoid implying
	// that installed means verified.
	action := surface.NextAction
	if action == "" {
		action = "-"
	}
	fmt.Fprintf(b, "%s | %s | %s | %s | %s | %s\n",
		surface.SurfaceID,
		surface.InstallState,
		surface.ProofState,
		surface.TrustScope,
		surface.EvidenceSource,
		action,
	)
}

func writeHumanTableDiffSummary(b *strings.Builder, summary []DiffSummary) {
	// Force summaries intentionally expose safe digests/counts, not full file
	// content.
	if len(summary) == 0 {
		return
	}
	b.WriteString("\nForce diff summary\n")
	for _, item := range summary {
		writeHumanTableDiffItem(b, item)
	}
}

func writeHumanTableDiffItem(b *strings.Builder, item DiffSummary) {
	// Diff summaries contain safe action/digest metadata only; file contents are
	// intentionally omitted from the human table.
	fmt.Fprintf(b, "- %s: %s", item.Path, item.Action)
	if item.Before != "" || item.After != "" {
		fmt.Fprintf(b, " [%s -> %s]", item.Before, item.After)
	}
	if item.Backup != "" {
		fmt.Fprintf(b, " (backup: %s)", item.Backup)
	}
	b.WriteString("\n")
}

func buildStatus(opts Options, installPreview bool) (Status, error) {
	// Derived repository identity is a sanitized local fallback, not externally
	// signed source proof.
	// Status generation is a snapshot; it does not persist files unless Install
	// enters write mode.
	// Git metadata is captured as strings and may be empty when local git cannot
	// provide it.
	// Aggregate states are derived from surface rows after any preview action
	// text is applied.
	// Gaps and next actions are derived from the same surface slice to keep JSON
	// fields consistent.
	// GeneratedAt uses the normalized clock so tests and CLI output can replay
	// deterministic snapshots.
	repoID := opts.RepositoryID
	if repoID == "" {
		repoID = derivedRepositoryID(opts.RepoRoot)
	}
	surfaces := buildSurfaces(opts)
	if installPreview {
		applyInstallPreviewActions(surfaces)
	}
	return statusFromSurfaces(opts, repoID, surfaces), nil
}

func statusFromSurfaces(opts Options, repoID string, surfaces []Surface) Status {
	// Aggregate install/proof states, gaps, and actions are all derived from the
	// same surface snapshot to avoid prose-only closure over missing evidence.
	gitHead, gitBranch := statusGitRefs(opts.RepoRoot)
	state := surfaceStatusState(surfaces)
	return Status{
		SchemaVersion:     SchemaVersion,
		Profile:           opts.Profile,
		RepositoryID:      repoID,
		RepositoryRootRef: "current_repository",
		// Git fields are observations from this checkout, not immutable source
		// proof.
		GitHead:   gitHead,
		GitBranch: gitBranch,
		// Aggregate states are derived, never hand-authored.
		InstallState: state.install,
		ProofState:   state.proof,
		Surfaces:     surfaces,
		// Gaps and actions preserve the same surface evidence used for verdicts.
		Gaps:        state.gaps,
		NextActions: state.actions,
		GeneratedAt: opts.Now.Format(time.RFC3339),
	}
}

type surfaceState struct {
	install string
	proof   string
	gaps    []Gap
	actions []NextAction
}

func surfaceStatusState(surfaces []Surface) surfaceState {
	// Every aggregate field is recomputed from measured surfaces, not from task
	// checkboxes or generated prose.
	return surfaceState{
		install: aggregateInstallState(surfaces),
		proof:   aggregateProofState(surfaces),
		gaps:    gapsFor(surfaces),
		actions: nextActionsFor(surfaces),
	}
}

func statusGitRefs(repoRoot string) (string, string) {
	// Git refs are local structural observations; empty strings mean git could
	// not provide that field during this snapshot.
	return gitOutput(repoRoot, "rev-parse", "--verify", "HEAD"),
		gitOutput(repoRoot, "branch", "--show-current")
}

func buildSurfaces(opts Options) []Surface {
	// Stable ordering keeps JSON and human output diffable across repeated runs.
	// Each surface reports a separate trust boundary instead of contributing to
	// an opaque health score.
	// Local hook, config, CI workflow, and non-applicable profile surfaces are
	// listed together so gaps stay visible.
	return []Surface{
		repositoryIdentitySurface(opts),
		hooksPathSurface(opts),
		hookSurface(opts, "pre-commit", SurfacePreCommitHook),
		hookSurface(opts, "post-commit", SurfacePostCommitHook),
		hookSurface(opts, "pre-push", SurfacePrePushHook),
		generatedFileSurface(opts, ".sdp-trace/config.json", SurfaceConfig),
		gitignoreSurface(opts),
		ciWorkflowSurface(opts),
		ciArtifactUploadSurface(opts),
		ciArtifactBundleSurface(opts),
		prCheckBindingSurface(),
		localWrappedCommandsSurface(),
		agentPromptSurface(),
	}
}

func applyInstallPreviewActions(surfaces []Surface) {
	// Preview rewrites only remediation text; it does not change measured
	// install/proof state.
	for i := range surfaces {
		if surfaces[i].InstallState == StateFail && strings.HasPrefix(surfaces[i].NextAction, "run install") {
			surfaces[i].NextAction = "rerun with --write to install this surface"
		}
	}
}

func repositoryIdentitySurface(opts Options) Surface {
	if opts.RepositoryID != "" {
		// Caller-supplied identity binds the observation to an intended source,
		// but it is still local structural evidence rather than external proof.
		return surface(SurfaceRepositoryIdentity, StatePass, StateNotAssessed, ScopeLocalStructural, "caller_supplied_repository_id", ReasonManualStepRequired, "", "")
	}
	// A sanitized origin hash avoids leaking remotes while keeping the identity
	// surface inspectable for follow-up binding work.
	return surface(SurfaceRepositoryIdentity, StatePass, StateNotAssessed, ScopeLocalStructural, "sanitized_origin_hash", ReasonManualStepRequired, "", "")
}

func hooksPathSurface(opts Options) Surface {
	// core.hooksPath is checkout-local git config, so it can prove installation
	// shape but not committed repository proof.
	value := strings.TrimSpace(gitOutput(opts.RepoRoot, "config", "--get", "core.hooksPath"))
	if value == ".githooks" {
		return surface(SurfaceHooksPath, StatePass, StateNotAssessed, ScopeLocalStructural, "git_config:core.hooksPath", ReasonHooksPathSet, "core.hooksPath=.githooks", "")
	}
	if value == "" {
		return surface(SurfaceHooksPath, StateFail, StateNotAssessed, ScopeLocalStructural, "git_config:core.hooksPath", ReasonHooksPathAbsent, "", "set core.hooksPath to .githooks")
	}
	return surface(SurfaceHooksPath, StateFail, StateCannotVerify, ScopeLocalStructural, "git_config:core.hooksPath", ReasonHooksPathMismatch, "core.hooksPath="+safeRef(value), "inspect existing hooks path before replacing it")
}

func hookSurface(opts Options, name, surfaceID string) Surface {
	// Hook existence is local structure; proof remains not_assessed until hook
	// output is observed from a git operation.
	rel := filepath.Join(".githooks", name)
	path := filepath.Join(opts.RepoRoot, rel)
	info, err := os.Stat(path)
	if err == nil {
		return presentHookSurface(info, surfaceID, rel)
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return surface(surfaceID, StateCannotVerify, StateCannotVerify, ScopeLocalStructural, "filesystem:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable hook path")
	}
	return surface(surfaceID, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonHookScriptAbsent, rel, "install generated hook script")
}

func presentHookSurface(info os.FileInfo, surfaceID, rel string) Surface {
	// A non-executable hook is present but cannot be verified as runnable.
	if info.IsDir() {
		return surface(surfaceID, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonHookScriptAbsent, rel, "install generated hook script")
	}
	state := StatePass
	if info.Mode()&0o111 == 0 {
		state = StateCannotVerify
	}
	return surface(surfaceID, state, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonHookScriptPresent, rel, "run a git operation to observe hook output")
}

func generatedFileSurface(opts Options, rel, surfaceID string) Surface {
	// Generated file presence is an install signal only; proof still requires an
	// observed run or CI artifact.
	path := filepath.Join(opts.RepoRoot, rel)
	_, err := os.Stat(path)
	if err == nil {
		return surface(surfaceID, StatePass, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonAlreadyInstalled, rel, "")
	}
	if errors.Is(err, os.ErrNotExist) {
		return surface(surfaceID, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonManualStepRequired, rel, "write generated observer configuration")
	}
	return surface(surfaceID, StateCannotVerify, StateCannotVerify, ScopeLocalStructural, "filesystem:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable generated file path")
}

func gitignoreSurface(opts Options) Surface {
	// Only the managed sdp-trace marker block is inspected; unrelated ignore
	// rules are outside this surface.
	rel := ".gitignore"
	data, err := os.ReadFile(filepath.Join(opts.RepoRoot, rel))
	if err == nil {
		return gitignoreContentSurface(rel, string(data))
	}
	if errors.Is(err, os.ErrNotExist) {
		return missingGitignoreSurface(rel)
	}
	return surface(SurfaceGitignore, StateCannotVerify, StateCannotVerify, ScopeLocalStructural, "filesystem:.gitignore", ReasonUnsafeOutputRefused, rel, "fix unreadable .gitignore")
}

func gitignoreContentSurface(rel, data string) Surface {
	if strings.Contains(data, "# sdp-trace begin") && strings.Contains(data, "# sdp-trace end") {
		// Marker presence proves only local ignore configuration, not CI proof.
		return surface(SurfaceGitignore, StatePass, StateNotAssessed, ScopeLocalStructural, "filesystem:.gitignore", ReasonAlreadyInstalled, rel, "")
	}
	return missingGitignoreSurface(rel)
}

func missingGitignoreSurface(rel string) Surface {
	return surface(SurfaceGitignore, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:.gitignore", ReasonManualStepRequired, rel, "add sdp-trace ignore block")
}

func ciWorkflowSurface(opts Options) Surface {
	// A checked-in workflow is not proof that CI has executed it.
	rel := filepath.Join(".github", "workflows", "sdp-trace-observe.yml")
	_, err := os.ReadFile(filepath.Join(opts.RepoRoot, rel))
	if err == nil {
		reason := ReasonCIWorkflowPresent
		proof := StateNotAssessed
		return surface(SurfaceCIWorkflow, StatePass, proof, ScopeLocalStructural, "filesystem:"+rel, reason, rel, "observe a CI run artifact before treating workflow as proof")
	}
	return missingCIWorkflowSurface(rel, err)
}

func missingCIWorkflowSurface(rel string, err error) Surface {
	if errors.Is(err, os.ErrNotExist) {
		// Missing workflow is an install gap; unreadable workflow paths below are
		// cannot_verify because the local filesystem state could not be replayed.
		return surface(SurfaceCIWorkflow, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonCIWorkflowAbsent, rel, "install GitHub Actions observer workflow")
	}
	return surface(SurfaceCIWorkflow, StateCannotVerify, StateCannotVerify, ScopeLocalStructural, "filesystem:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable workflow path")
}

func ciArtifactUploadSurface(opts Options) Surface {
	// Workflow upload declaration is local structure; uploaded artifacts need a
	// real CI run inspection.
	rel := filepath.Join(".github", "workflows", "sdp-trace-observe.yml")
	data, err := os.ReadFile(filepath.Join(opts.RepoRoot, rel))
	if err == nil && strings.Contains(string(data), "actions/upload-artifact") {
		return surface(SurfaceCIArtifactUpload, StatePass, StateNotAssessed, ScopeCIUploaded, "workflow_declaration:"+rel, ReasonCIArtifactUploadPresent, rel, "inspect uploaded artifact bundle from a real CI run")
	}
	return missingCIArtifactUploadSurface(rel, err)
}

func missingCIArtifactUploadSurface(rel string, err error) Surface {
	// Absence of upload-artifact configuration is an install gap; unreadable
	// workflow state is cannot_verify because local structure could not replay.
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return surface(SurfaceCIArtifactUpload, StateFail, StateNotAssessed, ScopeCIUploaded, "workflow_declaration:"+rel, ReasonCIArtifactUploadAbsent, rel, "declare CI artifact upload in observer workflow")
	}
	return surface(SurfaceCIArtifactUpload, StateCannotVerify, StateCannotVerify, ScopeCIUploaded, "workflow_declaration:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable workflow path")
}

func ciArtifactBundleSurface(opts Options) Surface {
	// Local CI artifact directories are only structural unless downloaded and
	// bound to CI artifact storage.
	rel := filepath.Join(".sdp-trace", "ci")
	path := filepath.Join(opts.RepoRoot, rel)
	entries, err := os.ReadDir(path)
	if err == nil && len(entries) > 0 {
		return surface(SurfaceCIArtifactBundleObservation, StatePass, StateNotAssessed, ScopeCIUploaded, "filesystem:"+rel, ReasonCIArtifactBundleObserved, rel, "treat as local structural only unless downloaded from CI artifact storage")
	}
	return missingCIArtifactBundleSurface(rel, err)
}

func missingCIArtifactBundleSurface(rel string, err error) Surface {
	// Missing local CI artifacts is not a failure by itself; it means CI evidence
	// has not been inspected for this profile.
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return surface(SurfaceCIArtifactBundleObservation, StateNotAssessed, StateNotAssessed, ScopeCIUploaded, "filesystem:"+rel, ReasonCIArtifactBundleNotObserved, rel, "run CI and inspect uploaded artifact bundle")
	}
	return surface(SurfaceCIArtifactBundleObservation, StateCannotVerify, StateCannotVerify, ScopeCIUploaded, "filesystem:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable CI artifact path")
}

func agentPromptSurface() Surface {
	return surface(SurfaceAgentPrompt, StateNotAssessed, StateNotAssessed, ScopeAgentReported, "agent_prompt:not_inspected", ReasonAgentReportedNotProof, "", "do not rely on prompt instructions as setup proof")
}

func prCheckBindingSurface() Surface {
	return surface(SurfacePRCheckBinding, StateNotAssessed, StateNotAssessed, ScopeNotApplicable, "github_pr_checks:not_inspected", ReasonOutsideProfileScope, "", "outside selected profile; no action required")
}

func localWrappedCommandsSurface() Surface {
	return surface(SurfaceLocalWrappedCommands, StateNotAssessed, StateNotAssessed, ScopeNotApplicable, "sdp_trace_runs:not_inspected", ReasonOutsideProfileScope, "", "outside selected profile; no action required")
}

func surface(id, install, proof, scope, source, reason, ref, action string) Surface {
	// Sanitization happens at construction so every renderer receives the same
	// safe observed ref.
	return Surface{
		SurfaceID:      id,
		InstallState:   install,
		ProofState:     proof,
		TrustScope:     scope,
		EvidenceSource: source,
		ObservedRef:    safeRef(ref),
		ReasonCode:     reason,
		NextAction:     action,
	}
}

func aggregateInstallState(surfaces []Surface) string {
	// Cannot_verify dominates install aggregation because unsafe surfaces could
	// not be inspected.
	state := StatePass
	for _, s := range surfaces {
		if s.InstallState == StateCannotVerify {
			return StateCannotVerify
		}
		if s.InstallState == StateFail {
			state = StateFail
		}
	}
	return state
}

func aggregateProofState(surfaces []Surface) string {
	// Empty or partially unassessed proof surfaces keep aggregate proof from
	// passing.
	if len(surfaces) == 0 {
		return StateNotAssessed
	}
	state := StatePass
	for _, s := range surfaces {
		state = combineProofState(state, s.ProofState)
		if state == StateCannotVerify {
			return state
		}
	}
	return state
}

func gapsFor(surfaces []Surface) []Gap {
	// Every non-passing surface remains visible as a gap with a concrete reason.
	gaps := make([]Gap, 0)
	for _, s := range surfaces {
		gap, ok := gapForSurface(s)
		if !ok {
			continue
		}
		gaps = append(gaps, gap)
	}
	return gaps
}

func combineProofState(current, next string) string {
	return combinedProofState(current, next)
}

func combinedProofState(current, next string) string {
	// Proof aggregation is monotonic: cannot_verify outranks fail, and
	// not_assessed prevents a clean pass.
	if proofStateDominates(next) {
		return StateCannotVerify
	}
	if proofStateFails(next) {
		return StateFail
	}
	return combineNonFailingProofState(current, next)
}

func proofStateDominates(state string) bool {
	return state == StateCannotVerify
}

func proofStateFails(state string) bool {
	return state == StateFail
}

func combineNonFailingProofState(current, next string) string {
	if current == StatePass && next == StateNotAssessed {
		// Any not_assessed proof surface prevents an aggregate pass unless a
		// stronger fail/cannot_verify state already dominated.
		return StateNotAssessed
	}
	return current
}

func gapForSurface(s Surface) (Gap, bool) {
	// Agent-prompt gaps get a custom explanation because prompt cooperation is
	// not repository setup proof.
	if s.InstallState == StatePass && s.ProofState == StatePass {
		return Gap{}, false
	}
	if agentPromptNotAssessed(s) {
		return agentPromptGap(s), true
	}
	return Gap{SurfaceID: s.SurfaceID, ReasonCode: s.ReasonCode, Detail: gapDetail(s)}, true
}

func agentPromptNotAssessed(s Surface) bool {
	return s.InstallState == StateNotAssessed && s.ProofState == StateNotAssessed && s.SurfaceID == SurfaceAgentPrompt
}

func agentPromptGap(s Surface) Gap {
	return Gap{SurfaceID: s.SurfaceID, ReasonCode: s.ReasonCode, Detail: "agent prompt cooperation is not repository setup proof"}
}

func gapDetail(s Surface) string {
	if s.NextAction != "" {
		// Prefer concrete remediation over terse reason codes in human summaries.
		return s.NextAction
	}
	return s.ReasonCode
}

func nextActionsFor(surfaces []Surface) []NextAction {
	// Stable sorting makes remediation output deterministic for docs and tests.
	// Empty actions are omitted so already-satisfied surfaces do not create
	// misleading follow-up work.
	actions := make([]NextAction, 0)
	for _, s := range surfaces {
		if s.NextAction == "" {
			continue
		}
		actions = append(actions, nextActionForSurface(s))
	}
	sort.SliceStable(actions, func(i, j int) bool {
		return actions[i].SurfaceID < actions[j].SurfaceID
	})
	return actions
}

func nextActionForSurface(s Surface) NextAction {
	return NextAction{SurfaceID: s.SurfaceID, ActionText: s.NextAction, Blocking: surfaceActionBlocking(s)}
}

func surfaceActionBlocking(s Surface) bool {
	return s.InstallState == StateFail || s.InstallState == StateCannotVerify || s.ProofState == StateCannotVerify
}

func writeInstallFiles(opts Options) ([]DiffSummary, error) {
	// Validate hooksPath before writing any files because changing it affects
	// all local git hook execution.
	// Repository files are written before hooksPath is changed so a failed file
	// write does not partially redirect local hooks.
	if err := ensureNoUnsafeHooksPath(opts); err != nil {
		return nil, err
	}
	summaries, err := writeInstallTargets(opts)
	if err != nil {
		return summaries, err
	}
	summary, err := updateGitignore(opts)
	if err != nil {
		return summaries, err
	}
	summaries = append(summaries, summary...)
	return appendHooksPathSummary(opts, summaries)
}

func writeInstallTargets(opts Options) ([]DiffSummary, error) {
	// Generated config and generated docs share one repository ID for internal
	// consistency.
	// Summaries accumulate only safe path/action/digest facts from each target.
	repoID := opts.RepositoryID
	if repoID == "" {
		repoID = derivedRepositoryID(opts.RepoRoot)
	}
	summaries := make([]DiffSummary, 0)
	for _, target := range installTargets(opts, repoID) {
		summary, err := writeTarget(opts, target)
		if err != nil {
			return summaries, err
		}
		summaries = append(summaries, summary...)
	}
	return summaries, nil
}

func appendHooksPathSummary(opts Options, summaries []DiffSummary) ([]DiffSummary, error) {
	// hooksPath changes are local git config mutations and are summarized
	// separately from repository file writes.
	// Existing non-default hook paths are summarized only in force mode after the
	// caller has accepted replacement.
	previousHooksPath := strings.TrimSpace(gitOutput(opts.RepoRoot, "config", "--get", "core.hooksPath"))
	if opts.Force && isDifferentHooksPath(previousHooksPath) {
		summaries = append(summaries, DiffSummary{
			Path:    "git_config:core.hooksPath",
			Action:  "overwrite_hooks_path",
			Before:  safeRef(previousHooksPath),
			After:   ".githooks",
			Summary: "replace local checkout hooks path reference",
		})
	}
	if err := runGit(opts.RepoRoot, "config", "core.hooksPath", ".githooks"); err != nil {
		return summaries, err
	}
	return summaries, nil
}

func isDifferentHooksPath(path string) bool {
	return path != "" && path != ".githooks"
}

func ensureNoUnsafeHooksPath(opts Options) error {
	// Existing non-.githooks values require --force so user hook configuration is
	// not silently replaced.
	value := strings.TrimSpace(gitOutput(opts.RepoRoot, "config", "--get", "core.hooksPath"))
	if value == "" || value == ".githooks" || opts.Force {
		return nil
	}
	return fmt.Errorf("%s: core.hooksPath is %s; use --force only after reviewing existing hooks", ReasonHooksPathMismatch, safeRef(value))
}

func writeTarget(opts Options, target targetFile) ([]DiffSummary, error) {
	// Every generated target is resolved through containment checks before any
	// read or write.
	// Existing files and new files take separate paths so force-mode overwrite
	// policy cannot affect first-time installs.
	path, err := safeTargetPath(opts, target)
	if err != nil {
		return nil, err
	}
	mode := targetMode(target)
	data := []byte(target.content)
	if existing, err := os.ReadFile(path); err == nil {
		return writeExistingTarget(opts, target, path, mode, existing, data)
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return writeNewTarget(path, data, mode)
}

func safeTargetPath(opts Options, target targetFile) (string, error) {
	// Clean plus Rel prevents generated target paths from escaping the selected
	// repository root.
	path := filepath.Clean(filepath.Join(opts.RepoRoot, target.path))
	rel, relErr := filepath.Rel(opts.RepoRoot, path)
	if targetPathEscapes(rel, relErr) {
		return "", fmt.Errorf("%s: target outside repository", ReasonUnsafeOutputRefused)
	}
	return path, nil
}

func targetPathEscapes(rel string, relErr error) bool {
	if relErr != nil {
		// Failed relative-path calculation is treated as containment failure.
		return true
	}
	return invalidRelativeTarget(rel)
}

func invalidRelativeTarget(rel string) bool {
	// Install targets must resolve to concrete files below the repository root.
	if rel == "." || rel == ".." {
		return true
	}
	if filepath.IsAbs(rel) {
		return true
	}
	return strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

func targetMode(target targetFile) os.FileMode {
	if target.executable {
		// Hook targets must retain executable mode; documentation targets must
		// not be made executable.
		return 0o755
	}
	return 0o644
}

func writeExistingTarget(opts Options, target targetFile, path string, mode os.FileMode, existing, data []byte) ([]DiffSummary, error) {
	// Differing files are overwritten only with --force after backup/diff summary
	// protections are available.
	if string(existing) == target.content {
		if target.executable {
			return nil, os.Chmod(path, mode)
		}
		return nil, nil
	}
	if !opts.Force {
		return nil, fmt.Errorf("%s: %s exists and differs; use --force after reviewing safe diff", ReasonManualStepRequired, target.path)
	}
	return overwriteTarget(target, path, mode, existing, data)
}

func overwriteTarget(target targetFile, path string, mode os.FileMode, existing, data []byte) ([]DiffSummary, error) {
	// Backup first, then write; this gives force mode a local recovery path.
	// The summary stores only digests, byte counts, and line counts for safe
	// review.
	// Directory creation still happens after backup so existing content is
	// protected before the replacement write.
	// The replacement write uses the same target mode as first-time generation.
	if err := os.WriteFile(path+".bak", existing, 0o644); err != nil {
		return nil, fmt.Errorf("%s: backup failed for %s", ReasonUnsafeOutputRefused, target.path)
	}
	summary := DiffSummary{
		Path:    target.path,
		Action:  "overwrite_existing_file",
		Before:  contentSummary(existing),
		After:   contentSummary(data),
		Summary: "replace generated file content using safe byte and line counts",
		Backup:  target.path + ".bak",
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return []DiffSummary{summary}, os.WriteFile(path, data, mode)
}

func writeNewTarget(path string, data []byte, mode os.FileMode) ([]DiffSummary, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	// New generated files do not need a force diff summary because there is no
	// previous repository content to overwrite.
	return nil, os.WriteFile(path, data, mode)
}

func updateGitignore(opts Options) ([]DiffSummary, error) {
	// .gitignore handling is limited to the managed block and never interprets
	// unrelated patterns as proof.
	// Missing .gitignore can be created without force because no user content is
	// overwritten.
	path := filepath.Join(opts.RepoRoot, ".gitignore")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, os.WriteFile(path, []byte(gitignoreBlock), 0o644)
	}
	if err != nil {
		return nil, err
	}
	text := string(data)
	start, end := locateGitignoreBlock(text)
	if start >= 0 {
		return updateSdpTraceGitignoreBlock(opts, path, text, data, start, end)
	}
	return appendGitignoreMarker(path, text)
}

func locateGitignoreBlock(text string) (int, int) {
	// Both markers must be ordered correctly before the block is considered
	// manageable.
	start := strings.Index(text, gitignoreBeginMarker)
	if start < 0 {
		return -1, -1
	}
	end := strings.Index(text, gitignoreEndMarker)
	if end < start {
		return -1, -1
	}
	return start, end + len(gitignoreEndMarker)
}

func updateSdpTraceGitignoreBlock(opts Options, path, text string, data []byte, start, end int) ([]DiffSummary, error) {
	// Replacing an existing managed block requires --force and a backup because
	// the user may have edited it.
	// Exact generated block matches are idempotent and produce no summary.
	// Replacement preserves all text before and after the managed marker range.
	// The backup stores the full previous .gitignore, not just the managed block.
	// Diff summaries again use digest/size metadata rather than raw ignore
	// content.
	// Force mode is the only path that can replace a divergent managed block.
	current := text[start:end]
	if current == strings.TrimSuffix(gitignoreBlock, "\n") {
		return nil, nil
	}
	if !opts.Force {
		return nil, fmt.Errorf("%s: .gitignore sdp-trace block differs; use --force after reviewing safe diff", ReasonManualStepRequired)
	}
	if err := os.WriteFile(path+".bak", data, 0o644); err != nil {
		return nil, fmt.Errorf("%s: backup failed for .gitignore", ReasonUnsafeOutputRefused)
	}
	next := text[:start] + strings.TrimSuffix(gitignoreBlock, "\n") + text[end:]
	return replacedGitignoreBlockSummary(data, next), os.WriteFile(path, []byte(next), 0o644)
}

func replacedGitignoreBlockSummary(before []byte, next string) []DiffSummary {
	// The summary is digest and count metadata only; the managed ignore content
	// is not copied into force-mode review output.
	return []DiffSummary{{
		Path:    ".gitignore",
		Action:  "replace_sdp_trace_block",
		Before:  contentSummary(before),
		After:   contentSummary([]byte(next)),
		Summary: "replace marked sdp-trace gitignore block using safe byte and line counts",
		Backup:  ".gitignore.bak",
	}}
}

func appendGitignoreMarker(path, text string) ([]DiffSummary, error) {
	// Preserve existing ignore content and append a clean newline boundary before
	// the managed block.
	if strings.TrimSpace(text) != "" && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	text += gitignoreBlock
	return nil, os.WriteFile(path, []byte(text), 0o644)
}

func installTargets(opts Options, repoID string) []targetFile {
	// Generated files stay portable: hooks, local config/docs, and a GitHub
	// workflow declaration only.
	return []targetFile{
		{path: ".sdp-trace/README.md", content: sdpTraceReadme()},
		{path: ".sdp-trace/config.json", content: sdpTraceConfig(opts, repoID)},
		{path: ".githooks/pre-commit", content: hookScript("pre-commit"), executable: true},
		{path: ".githooks/post-commit", content: hookScript("post-commit"), executable: true},
		{path: ".githooks/pre-push", content: hookScript("pre-push"), executable: true},
		{path: ".github/workflows/sdp-trace-observe.yml", content: githubWorkflow()},
	}
}

func sdpTraceReadme() string {
	// The generated README states the evidence boundary beside generated files so
	// future agents do not treat local hooks as external proof.
	return `# sdp-trace repository observer

This directory is generated by ` + "`sdp-trace install repo-observer --profile github-actions-git-hooks-v1 --write`" + `.

It stores local structural observations only. Local hook output and checked-in
configuration are not external proof. CI-uploaded artifacts become stronger
evidence only when the selected profile can inspect the uploaded artifact bundle.
`
}

func sdpTraceConfig(opts Options, repoID string) string {
	// Config records installed paths and trust boundary metadata, not a proof
	// claim.
	// Installed file paths are sorted so repeated installs produce stable config
	// output.
	// The local_config_note calls out core.hooksPath as checkout-local state.
	// InstalledFiles is generated from a manifest-only target list so content
	// changes do not affect config shape.
	// JSON marshal errors are impossible for this map shape, so generation keeps
	// the helper signature simple.
	// The generated config is an install manifest, not a verifier result.
	payload := sdpTraceConfigPayload(opts, repoID)
	data, _ := json.MarshalIndent(payload, "", "  ")
	return string(data) + "\n"
}

func sdpTraceConfigPayload(opts Options, repoID string) map[string]any {
	// Generated config is local structural metadata. CI proof remains
	// not_assessed until a later artifact-observation path binds uploaded output.
	return map[string]any{
		"schema_version":   "sdp-trace-repo-observer-config-v1",
		"profile":          opts.Profile,
		"repository_id":    repoID,
		"trust_boundary":   "local_structural_until_ci_artifact_observed",
		"installed_files":  sdpTraceConfigPaths(),
		"install_metadata": sdpTraceInstallMetadata(),
	}
}

func sdpTraceConfigPaths() []string {
	// Installed file paths are an index of generated surfaces, not proof that
	// those surfaces have executed.
	paths := make([]string, 0)
	for _, target := range installTargetsForManifest() {
		paths = append(paths, target.path)
	}
	paths = append(paths, ".gitignore:# sdp-trace begin")
	sort.Strings(paths)
	return paths
}

func sdpTraceInstallMetadata() map[string]string {
	// Local core.hooksPath is called out explicitly because it is checkout-local
	// configuration and not committed repository evidence.
	return map[string]string{
		"generated_by":      "sdp-trace install repo-observer",
		"template_version":  SchemaVersion,
		"local_config_note": "core.hooksPath is local checkout configuration",
	}
}

func installTargetsForManifest() []targetFile {
	// Manifest targets omit content so the config stays a compact structural
	// index.
	return []targetFile{
		{path: ".sdp-trace/README.md"},
		{path: ".sdp-trace/config.json"},
		{path: ".githooks/pre-commit"},
		{path: ".githooks/post-commit"},
		{path: ".githooks/pre-push"},
		{path: ".github/workflows/sdp-trace-observe.yml"},
	}
}

func hookScript(name string) string {
	// Hooks collect metadata and diagnostics only; they do not enforce
	// sdp-trace policy or block user operations.
	return `#!/usr/bin/env bash
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
cd "$repo_root"
event="` + name + `"
ts="$(date -u +%Y%m%dT%H%M%SZ)"
out_dir=".sdp-trace/hooks/${event}/${ts}"
mkdir -p "$out_dir"

{
  printf 'event=%s\n' "$event"
  printf 'git_head=%s\n' "$(git rev-parse --verify HEAD 2>/dev/null || printf unknown)"
  printf 'git_branch=%s\n' "$(git branch --show-current 2>/dev/null || printf unknown)"
} > "$out_dir/metadata.env"

git status --short > "$out_dir/status.txt" || true
git diff --cached --name-status > "$out_dir/staged-files.txt" || true
git diff --check > "$out_dir/diff-check.txt" 2>&1 || true
`
}

func contentSummary(data []byte) string {
	// Force summaries expose digest/size metadata instead of raw file content.
	sum := sha256.Sum256(data)
	lines := 0
	if len(data) > 0 {
		lines = strings.Count(string(data), "\n")
		if data[len(data)-1] != '\n' {
			lines++
		}
	}
	return fmt.Sprintf("sha256:%s bytes:%d lines:%d", hex.EncodeToString(sum[:])[:16], len(data), lines)
}

func githubWorkflow() string {
	// The generated workflow uploads observations; proof is established only
	// after inspecting a real CI artifact.
	// It captures repository metadata and safe status output, not raw secrets or
	// local hook output.
	// The workflow is intentionally small and uses shell commands only as a thin
	// CI metadata launcher.
	// It uploads `.sdp-trace/ci` as an artifact source for later inspection.
	// Repository proof remains not_assessed until that uploaded artifact is
	// inspected by a separate evidence path.
	// The checkout step is the only third-party action used by this template.
	return `name: sdp-trace observe

on:
  pull_request:
  push:

jobs:
  observe:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      actions: read
    steps:
      - uses: actions/checkout@v4
      - name: Capture safe repository metadata
        shell: bash
        run: |
          set -euo pipefail
          mkdir -p .sdp-trace/ci
          {
            printf 'github_repository=%s\n' "$GITHUB_REPOSITORY"
            printf 'github_run_id=%s\n' "$GITHUB_RUN_ID"
            printf 'github_sha=%s\n' "$GITHUB_SHA"
            printf 'github_ref=%s\n' "$GITHUB_REF"
          } > .sdp-trace/ci/metadata.env
          git status --short > .sdp-trace/ci/status.txt
          git diff --check > .sdp-trace/ci/diff-check.txt 2>&1 || true
      - name: Optional Bazel test smoke
        shell: bash
        run: |
          set -euo pipefail
          if [ -f MODULE.bazel ] || [ -f WORKSPACE ] || [ -f WORKSPACE.bazel ]; then
            printf 'bazel_config_present\n' > .sdp-trace/ci/bazel-test.txt
          else
            printf 'bazel_not_configured\n' > .sdp-trace/ci/bazel-test.txt
          fi
      - uses: actions/upload-artifact@v4
        with:
          name: sdp-trace-observer
          path: .sdp-trace/ci/
`
}

func repoRoot(start string) (string, error) {
	// Git is the source of truth for repository root discovery.
	out, err := exec.Command("git", "-C", start, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("cannot locate git repository root")
	}
	return strings.TrimSpace(string(out)), nil
}

func gitOutput(root string, args ...string) string {
	// Missing git data stays an empty observation so doctor can still report the
	// rest of the surfaces.
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func runGit(root string, args ...string) error {
	// Mutating git commands include command output in errors for actionable DX.
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %s failed: %s", strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
	return nil
}

func derivedRepositoryID(root string) string {
	// Derived IDs hash sanitized remote data to avoid rendering credentials.
	origin := gitOutput(root, "config", "--get", "remote.origin.url")
	if strings.TrimSpace(origin) == "" {
		origin = "current_repository"
	}
	sanitized := sanitizeOrigin(origin)
	sum := sha256.Sum256([]byte(sanitized))
	return "repo_" + hex.EncodeToString(sum[:])[:16]
}

func sanitizeOrigin(origin string) string {
	// URL credentials are stripped before repository identity hashing.
	origin = removeOriginFragment(strings.TrimSpace(origin))
	hadURLCredentials := hasURLCredentials(origin)
	origin = removeOriginCredentials(origin)
	if hadURLCredentials {
		return origin
	}
	return originTail(origin)
}

func removeOriginFragment(origin string) string {
	if idx := strings.Index(origin, "#"); idx >= 0 {
		// URL fragments are local navigation hints and not repository identity.
		return origin[:idx]
	}
	return origin
}

func removeOriginCredentials(origin string) string {
	// SCP-like and URL remotes encode userinfo differently, so redact both forms.
	if strings.Contains(origin, "@") && !strings.Contains(origin, "://") {
		return origin[strings.LastIndex(origin, "@")+1:]
	}
	if originHasURLCredentials(origin) {
		return originWithoutURLCredentials(origin)
	}
	return origin
}

func originHasURLCredentials(origin string) bool {
	// Treat @ as URL credentials only when it appears after a URL scheme.
	at := strings.LastIndex(origin, "@")
	return at >= 0 && strings.Contains(origin[:max(at, 0)], "://")
}

func originWithoutURLCredentials(origin string) string {
	at := strings.LastIndex(origin, "@")
	schemeEnd := strings.Index(origin, "://")
	// Preserve scheme and host while dropping userinfo from the rendered remote.
	return origin[:schemeEnd+3] + origin[at+1:]
}

func hasURLCredentials(origin string) bool {
	return originHasURLCredentials(origin)
}

func originTail(origin string) string {
	// Non-credential origins use the final owner/repo-ish path components.
	origin = strings.ReplaceAll(origin, "\\", "/")
	parts := strings.Split(origin, "/")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return origin
}

func safeRef(ref string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	// Normalize path separators before redacting unsafe config values in
	// human-facing output.
	ref = strings.ReplaceAll(ref, "\\", "/")
	if strings.HasPrefix(ref, "/") || strings.Contains(ref, ":/") {
		return "unsafe_absolute_path_redacted"
	}
	return ref
}
