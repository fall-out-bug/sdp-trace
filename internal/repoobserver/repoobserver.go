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
	opts, err := normalizeOptions(opts)
	if err != nil {
		return Status{}, err
	}
	return buildStatus(opts, false)
}

func Install(opts Options) (Status, error) {
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

func HumanTable(status Status) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Profile: %s\n", status.Profile)
	fmt.Fprintf(&b, "Repository: %s\n", status.RepositoryID)
	fmt.Fprintf(&b, "Install state: %s\n", status.InstallState)
	fmt.Fprintf(&b, "Proof state: %s\n\n", status.ProofState)
	b.WriteString("Surface | Install state | Proof state | Trust scope | Evidence source | Next action\n")
	b.WriteString("--- | --- | --- | --- | --- | ---\n")
	for _, surface := range status.Surfaces {
		action := surface.NextAction
		if action == "" {
			action = "-"
		}
		fmt.Fprintf(&b, "%s | %s | %s | %s | %s | %s\n",
			surface.SurfaceID,
			surface.InstallState,
			surface.ProofState,
			surface.TrustScope,
			surface.EvidenceSource,
			action,
		)
	}
	if len(status.ForceDiffSummary) > 0 {
		b.WriteString("\nForce diff summary\n")
		for _, item := range status.ForceDiffSummary {
			fmt.Fprintf(&b, "- %s: %s", item.Path, item.Action)
			if item.Before != "" || item.After != "" {
				fmt.Fprintf(&b, " [%s -> %s]", item.Before, item.After)
			}
			if item.Backup != "" {
				fmt.Fprintf(&b, " (backup: %s)", item.Backup)
			}
			b.WriteString("\n")
		}
	}
	b.WriteString("\nNote: core.hooksPath is local checkout configuration and is not committed into the repository.\n")
	return b.String()
}

func normalizeOptions(opts Options) (Options, error) {
	if strings.TrimSpace(opts.Profile) == "" {
		opts.Profile = ProfileGithubActionsGitHooksV1
	}
	if opts.Profile != ProfileGithubActionsGitHooksV1 {
		return Options{}, fmt.Errorf("repo observer requires --profile %s", ProfileGithubActionsGitHooksV1)
	}
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
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	if opts.RepositoryID != "" && !safeIDPattern.MatchString(opts.RepositoryID) {
		return Options{}, fmt.Errorf("%s: repository id must match [A-Za-z0-9_.-]+", ReasonUnsafeOutputRefused)
	}
	return opts, nil
}

func buildStatus(opts Options, installPreview bool) (Status, error) {
	repoID := opts.RepositoryID
	if repoID == "" {
		repoID = derivedRepositoryID(opts.RepoRoot)
	}
	surfaces := []Surface{
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
	if installPreview {
		for i := range surfaces {
			if surfaces[i].InstallState == StateFail && strings.HasPrefix(surfaces[i].NextAction, "run install") {
				surfaces[i].NextAction = "rerun with --write to install this surface"
			}
		}
	}
	status := Status{
		SchemaVersion:     SchemaVersion,
		Profile:           opts.Profile,
		RepositoryID:      repoID,
		RepositoryRootRef: "current_repository",
		GitHead:           gitOutput(opts.RepoRoot, "rev-parse", "--verify", "HEAD"),
		GitBranch:         gitOutput(opts.RepoRoot, "branch", "--show-current"),
		InstallState:      aggregateInstallState(surfaces),
		ProofState:        aggregateProofState(surfaces),
		Surfaces:          surfaces,
		Gaps:              gapsFor(surfaces),
		NextActions:       nextActionsFor(surfaces),
		GeneratedAt:       opts.Now.Format(time.RFC3339),
	}
	return status, nil
}

func repositoryIdentitySurface(opts Options) Surface {
	if opts.RepositoryID != "" {
		return surface(SurfaceRepositoryIdentity, StatePass, StateNotAssessed, ScopeLocalStructural, "caller_supplied_repository_id", ReasonManualStepRequired, "", "")
	}
	return surface(SurfaceRepositoryIdentity, StatePass, StateNotAssessed, ScopeLocalStructural, "sanitized_origin_hash", ReasonManualStepRequired, "", "")
}

func hooksPathSurface(opts Options) Surface {
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
	rel := filepath.Join(".githooks", name)
	path := filepath.Join(opts.RepoRoot, rel)
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		state := StatePass
		if info.Mode()&0o111 == 0 {
			state = StateCannotVerify
		}
		return surface(surfaceID, state, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonHookScriptPresent, rel, "run a git operation to observe hook output")
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return surface(surfaceID, StateCannotVerify, StateCannotVerify, ScopeLocalStructural, "filesystem:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable hook path")
	}
	return surface(surfaceID, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonHookScriptAbsent, rel, "install generated hook script")
}

func generatedFileSurface(opts Options, rel, surfaceID string) Surface {
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
	rel := ".gitignore"
	data, err := os.ReadFile(filepath.Join(opts.RepoRoot, rel))
	if err == nil && strings.Contains(string(data), "# sdp-trace begin") && strings.Contains(string(data), "# sdp-trace end") {
		return surface(SurfaceGitignore, StatePass, StateNotAssessed, ScopeLocalStructural, "filesystem:.gitignore", ReasonAlreadyInstalled, rel, "")
	}
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return surface(SurfaceGitignore, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:.gitignore", ReasonManualStepRequired, rel, "add sdp-trace ignore block")
	}
	return surface(SurfaceGitignore, StateCannotVerify, StateCannotVerify, ScopeLocalStructural, "filesystem:.gitignore", ReasonUnsafeOutputRefused, rel, "fix unreadable .gitignore")
}

func ciWorkflowSurface(opts Options) Surface {
	rel := filepath.Join(".github", "workflows", "sdp-trace-observe.yml")
	_, err := os.ReadFile(filepath.Join(opts.RepoRoot, rel))
	if err == nil {
		reason := ReasonCIWorkflowPresent
		proof := StateNotAssessed
		return surface(SurfaceCIWorkflow, StatePass, proof, ScopeLocalStructural, "filesystem:"+rel, reason, rel, "observe a CI run artifact before treating workflow as proof")
	}
	if errors.Is(err, os.ErrNotExist) {
		return surface(SurfaceCIWorkflow, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonCIWorkflowAbsent, rel, "install GitHub Actions observer workflow")
	}
	return surface(SurfaceCIWorkflow, StateCannotVerify, StateCannotVerify, ScopeLocalStructural, "filesystem:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable workflow path")
}

func ciArtifactUploadSurface(opts Options) Surface {
	rel := filepath.Join(".github", "workflows", "sdp-trace-observe.yml")
	data, err := os.ReadFile(filepath.Join(opts.RepoRoot, rel))
	if err == nil && strings.Contains(string(data), "actions/upload-artifact") {
		return surface(SurfaceCIArtifactUpload, StatePass, StateNotAssessed, ScopeCIUploaded, "workflow_declaration:"+rel, ReasonCIArtifactUploadPresent, rel, "inspect uploaded artifact bundle from a real CI run")
	}
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return surface(SurfaceCIArtifactUpload, StateFail, StateNotAssessed, ScopeCIUploaded, "workflow_declaration:"+rel, ReasonCIArtifactUploadAbsent, rel, "declare CI artifact upload in observer workflow")
	}
	return surface(SurfaceCIArtifactUpload, StateCannotVerify, StateCannotVerify, ScopeCIUploaded, "workflow_declaration:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable workflow path")
}

func ciArtifactBundleSurface(opts Options) Surface {
	rel := filepath.Join(".sdp-trace", "ci")
	path := filepath.Join(opts.RepoRoot, rel)
	entries, err := os.ReadDir(path)
	if err == nil && len(entries) > 0 {
		return surface(SurfaceCIArtifactBundleObservation, StatePass, StateNotAssessed, ScopeCIUploaded, "filesystem:"+rel, ReasonCIArtifactBundleObserved, rel, "treat as local structural only unless downloaded from CI artifact storage")
	}
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
	if len(surfaces) == 0 {
		return StateNotAssessed
	}
	state := StatePass
	for _, s := range surfaces {
		if s.ProofState == StateCannotVerify {
			return StateCannotVerify
		}
		if s.ProofState == StateFail {
			state = StateFail
		}
		if s.ProofState == StateNotAssessed && state == StatePass {
			state = StateNotAssessed
		}
	}
	return state
}

func gapsFor(surfaces []Surface) []Gap {
	gaps := make([]Gap, 0)
	for _, s := range surfaces {
		if s.InstallState == StatePass && s.ProofState == StatePass {
			continue
		}
		if s.InstallState == StateNotAssessed && s.ProofState == StateNotAssessed && s.SurfaceID == SurfaceAgentPrompt {
			gaps = append(gaps, Gap{SurfaceID: s.SurfaceID, ReasonCode: s.ReasonCode, Detail: "agent prompt cooperation is not repository setup proof"})
			continue
		}
		detail := s.NextAction
		if detail == "" {
			detail = s.ReasonCode
		}
		gaps = append(gaps, Gap{SurfaceID: s.SurfaceID, ReasonCode: s.ReasonCode, Detail: detail})
	}
	return gaps
}

func nextActionsFor(surfaces []Surface) []NextAction {
	actions := make([]NextAction, 0)
	for _, s := range surfaces {
		if s.NextAction == "" {
			continue
		}
		blocking := s.InstallState == StateFail || s.InstallState == StateCannotVerify || s.ProofState == StateCannotVerify
		actions = append(actions, NextAction{SurfaceID: s.SurfaceID, ActionText: s.NextAction, Blocking: blocking})
	}
	sort.SliceStable(actions, func(i, j int) bool {
		return actions[i].SurfaceID < actions[j].SurfaceID
	})
	return actions
}

func writeInstallFiles(opts Options) ([]DiffSummary, error) {
	if err := ensureNoUnsafeHooksPath(opts); err != nil {
		return nil, err
	}
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
	summary, err := updateGitignore(opts)
	if err != nil {
		return summaries, err
	}
	summaries = append(summaries, summary...)
	previousHooksPath := strings.TrimSpace(gitOutput(opts.RepoRoot, "config", "--get", "core.hooksPath"))
	if opts.Force && previousHooksPath != "" && previousHooksPath != ".githooks" {
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

func ensureNoUnsafeHooksPath(opts Options) error {
	value := strings.TrimSpace(gitOutput(opts.RepoRoot, "config", "--get", "core.hooksPath"))
	if value == "" || value == ".githooks" || opts.Force {
		return nil
	}
	return fmt.Errorf("%s: core.hooksPath is %s; use --force only after reviewing existing hooks", ReasonHooksPathMismatch, safeRef(value))
}

func writeTarget(opts Options, target targetFile) ([]DiffSummary, error) {
	path := filepath.Clean(filepath.Join(opts.RepoRoot, target.path))
	rel, relErr := filepath.Rel(opts.RepoRoot, path)
	if relErr != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || filepath.IsAbs(rel) {
		return nil, fmt.Errorf("%s: target outside repository", ReasonUnsafeOutputRefused)
	}
	mode := os.FileMode(0o644)
	if target.executable {
		mode = 0o755
	}
	data := []byte(target.content)
	if existing, err := os.ReadFile(path); err == nil {
		if string(existing) == target.content {
			if target.executable {
				return nil, os.Chmod(path, mode)
			}
			return nil, nil
		}
		if !opts.Force {
			return nil, fmt.Errorf("%s: %s exists and differs; use --force after reviewing safe diff", ReasonManualStepRequired, target.path)
		}
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
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return nil, os.WriteFile(path, data, mode)
}

func updateGitignore(opts Options) ([]DiffSummary, error) {
	path := filepath.Join(opts.RepoRoot, ".gitignore")
	block := "# sdp-trace begin\n.sdp-trace/hooks/\n.sdp-trace/ci/\n.sdp-trace/install-diff.txt\n# sdp-trace end\n"
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, os.WriteFile(path, []byte(block), 0o644)
	}
	if err != nil {
		return nil, err
	}
	text := string(data)
	start := strings.Index(text, "# sdp-trace begin")
	end := strings.Index(text, "# sdp-trace end")
	if start >= 0 && end >= start {
		end += len("# sdp-trace end")
		current := text[start:end]
		if current == strings.TrimSuffix(block, "\n") {
			return nil, nil
		}
		if !opts.Force {
			return nil, fmt.Errorf("%s: .gitignore sdp-trace block differs; use --force after reviewing safe diff", ReasonManualStepRequired)
		}
		if err := os.WriteFile(path+".bak", data, 0o644); err != nil {
			return nil, fmt.Errorf("%s: backup failed for .gitignore", ReasonUnsafeOutputRefused)
		}
		next := text[:start] + strings.TrimSuffix(block, "\n") + text[end:]
		return []DiffSummary{{
			Path:    ".gitignore",
			Action:  "replace_sdp_trace_block",
			Before:  contentSummary(data),
			After:   contentSummary([]byte(next)),
			Summary: "replace marked sdp-trace gitignore block using safe byte and line counts",
			Backup:  ".gitignore.bak",
		}}, os.WriteFile(path, []byte(next), 0o644)
	}
	if strings.TrimSpace(text) != "" && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	text += block
	return nil, os.WriteFile(path, []byte(text), 0o644)
}

func installTargets(opts Options, repoID string) []targetFile {
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
	return `# sdp-trace repository observer

This directory is generated by ` + "`sdp-trace install repo-observer --profile github-actions-git-hooks-v1 --write`" + `.

It stores local structural observations only. Local hook output and checked-in
configuration are not external proof. CI-uploaded artifacts become stronger
evidence only when the selected profile can inspect the uploaded artifact bundle.
`
}

func sdpTraceConfig(opts Options, repoID string) string {
	paths := make([]string, 0)
	for _, target := range installTargetsForManifest() {
		paths = append(paths, target.path)
	}
	paths = append(paths, ".gitignore:# sdp-trace begin")
	sort.Strings(paths)
	payload := map[string]any{
		"schema_version":  "sdp-trace-repo-observer-config-v1",
		"profile":         opts.Profile,
		"repository_id":   repoID,
		"trust_boundary":  "local_structural_until_ci_artifact_observed",
		"installed_files": paths,
		"install_metadata": map[string]string{
			"generated_by":      "sdp-trace install repo-observer",
			"template_version":  SchemaVersion,
			"local_config_note": "core.hooksPath is local checkout configuration",
		},
	}
	data, _ := json.MarshalIndent(payload, "", "  ")
	return string(data) + "\n"
}

func installTargetsForManifest() []targetFile {
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
	out, err := exec.Command("git", "-C", start, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("cannot locate git repository root")
	}
	return strings.TrimSpace(string(out)), nil
}

func gitOutput(root string, args ...string) string {
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func runGit(root string, args ...string) error {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %s failed: %s", strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
	return nil
}

func derivedRepositoryID(root string) string {
	origin := gitOutput(root, "config", "--get", "remote.origin.url")
	if strings.TrimSpace(origin) == "" {
		origin = "current_repository"
	}
	sanitized := sanitizeOrigin(origin)
	sum := sha256.Sum256([]byte(sanitized))
	return "repo_" + hex.EncodeToString(sum[:])[:16]
}

func sanitizeOrigin(origin string) string {
	origin = strings.TrimSpace(origin)
	if idx := strings.Index(origin, "#"); idx >= 0 {
		origin = origin[:idx]
	}
	if strings.Contains(origin, "@") && !strings.Contains(origin, "://") {
		origin = origin[strings.LastIndex(origin, "@")+1:]
	}
	if at := strings.LastIndex(origin, "@"); strings.Contains(origin[:max(at, 0)], "://") && at >= 0 {
		schemeEnd := strings.Index(origin, "://")
		return origin[:schemeEnd+3] + origin[at+1:]
	}
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
	ref = strings.ReplaceAll(ref, "\\", "/")
	if strings.HasPrefix(ref, "/") || strings.Contains(ref, ":/") {
		return "unsafe_absolute_path_redacted"
	}
	return ref
}
