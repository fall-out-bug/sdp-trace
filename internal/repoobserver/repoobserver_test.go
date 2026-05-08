package repoobserver

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/feedback"
)

func TestDoctorSeparatesInstallStateFromProofState(t *testing.T) {
	repo := initRepo(t)
	status, err := Doctor(Options{
		RepoRoot:     repo,
		Profile:      ProfileGithubActionsGitHooksV1,
		RepositoryID: "demo_repo",
		Now:          time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if status.InstallState != StateFail {
		t.Fatalf("install state = %s", status.InstallState)
	}
	if status.ProofState != StateNotAssessed {
		t.Fatalf("proof state = %s", status.ProofState)
	}
	if status.RepositoryRootRef != "current_repository" {
		t.Fatalf("repository root ref leaked path: %s", status.RepositoryRootRef)
	}
	if surfaceByID(t, status, SurfaceHooksPath).ReasonCode != ReasonHooksPathAbsent {
		t.Fatalf("missing hooks path reason")
	}
	if surfaceByID(t, status, SurfaceAgentPrompt).TrustScope != ScopeAgentReported {
		t.Fatalf("agent prompt scope not represented")
	}
}

func TestInstallWritesIdempotentObserverFiles(t *testing.T) {
	repo := initRepo(t)
	status, err := Install(Options{
		RepoRoot:     repo,
		Profile:      ProfileGithubActionsGitHooksV1,
		RepositoryID: "demo_repo",
		Write:        true,
		Now:          time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if status.InstallState != StatePass {
		t.Fatalf("install state = %s", status.InstallState)
	}
	if got := strings.TrimSpace(gitOutputForTest(t, repo, "config", "--get", "core.hooksPath")); got != ".githooks" {
		t.Fatalf("hooksPath = %s", got)
	}
	for _, rel := range []string{
		".sdp-trace/README.md",
		".sdp-trace/config.json",
		".githooks/pre-commit",
		".githooks/post-commit",
		".githooks/pre-push",
		".github/workflows/sdp-trace-observe.yml",
	} {
		if _, err := os.Stat(filepath.Join(repo, rel)); err != nil {
			t.Fatalf("missing %s: %v", rel, err)
		}
	}
	var config struct {
		RepositoryID    string            `json:"repository_id"`
		InstalledFiles  []string          `json:"installed_files"`
		InstallMetadata map[string]string `json:"install_metadata"`
	}
	readJSONForTest(t, filepath.Join(repo, ".sdp-trace", "config.json"), &config)
	if config.RepositoryID != "demo_repo" {
		t.Fatalf("config repository_id = %s", config.RepositoryID)
	}
	if len(config.InstalledFiles) == 0 || config.InstallMetadata["generated_by"] == "" {
		t.Fatalf("config missing manifest or metadata: %+v", config)
	}
	if _, err := Install(Options{RepoRoot: repo, Profile: ProfileGithubActionsGitHooksV1, RepositoryID: "demo_repo", Write: true}); err != nil {
		t.Fatalf("idempotent install failed: %v", err)
	}
}

func TestDoctorUsesInstalledConfigRepositoryID(t *testing.T) {
	repo := initRepo(t)
	if _, err := Install(Options{
		RepoRoot:     repo,
		Profile:      ProfileGithubActionsGitHooksV1,
		RepositoryID: "demo_repo",
		Write:        true,
		Now:          time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatal(err)
	}
	status, err := Doctor(Options{
		RepoRoot: repo,
		Profile:  ProfileGithubActionsGitHooksV1,
		Now:      time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if status.RepositoryID != "demo_repo" {
		t.Fatalf("repository id = %s", status.RepositoryID)
	}
}

func TestInstallForceProducesSafeDiffSummary(t *testing.T) {
	repo := initRepo(t)
	writeFileForTest(t, filepath.Join(repo, ".githooks", "pre-commit"), "#!/usr/bin/env bash\nprintf old\n")
	status, err := Install(Options{
		RepoRoot:     repo,
		Profile:      ProfileGithubActionsGitHooksV1,
		RepositoryID: "demo_repo",
		Write:        true,
		Force:        true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(status.ForceDiffSummary) == 0 {
		t.Fatalf("missing force diff summary")
	}
	if strings.Contains(status.ForceDiffSummary[0].Path, repo) || strings.Contains(status.ForceDiffSummary[0].Backup, repo) {
		t.Fatalf("force summary leaked absolute path: %+v", status.ForceDiffSummary[0])
	}
	if _, err := os.Stat(filepath.Join(repo, ".githooks", "pre-commit.bak")); err != nil {
		t.Fatalf("missing backup: %v", err)
	}
}

func TestInstallRefusesHooksPathMismatchWithoutForce(t *testing.T) {
	repo := initRepo(t)
	runGitForTest(t, repo, "config", "core.hooksPath", "custom-hooks")
	_, err := Install(Options{
		RepoRoot:     repo,
		Profile:      ProfileGithubActionsGitHooksV1,
		RepositoryID: "demo_repo",
		Write:        true,
	})
	if err == nil || !strings.Contains(err.Error(), ReasonHooksPathMismatch) {
		t.Fatalf("expected hooks path mismatch, got %v", err)
	}
}

func TestInvalidRepositoryIDRejected(t *testing.T) {
	repo := initRepo(t)
	_, err := Doctor(Options{
		RepoRoot:     repo,
		Profile:      ProfileGithubActionsGitHooksV1,
		RepositoryID: "/tmp/private/repo",
	})
	if err == nil || !strings.Contains(err.Error(), ReasonUnsafeOutputRefused) {
		t.Fatalf("expected invalid repository id, got %v", err)
	}
}

func TestDoctorObservesOnlyValidFeedbackEvents(t *testing.T) {
	repo := initRepo(t)
	messagePath := filepath.Join(repo, "feedback.md")
	writeFileForTest(t, messagePath, "Do not edit the GSD plan; record feedback separately.\n")
	event, err := feedback.Record(feedback.Options{
		Kind:        "corrective_feedback",
		From:        "human",
		To:          "gsd",
		SourceRef:   "codex-thread",
		Summary:     "Boundary correction for GSD plan observation.",
		MessageFile: messagePath,
		Now:         time.Date(2026, 5, 8, 20, 30, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := feedback.WriteJSON(filepath.Join(repo, ".sdp-trace", "feedback", "event.json"), event); err != nil {
		t.Fatal(err)
	}
	status, err := Doctor(Options{
		RepoRoot:     repo,
		Profile:      ProfileGithubActionsGitHooksV1,
		RepositoryID: "demo_repo",
	})
	if err != nil {
		t.Fatal(err)
	}
	surface := surfaceByID(t, status, SurfaceFeedbackEvents)
	if surface.InstallState != StateNotAssessed || surface.ProofState != StateNotAssessed || surface.ReasonCode != ReasonFeedbackEventsObserved {
		t.Fatalf("feedback surface = %+v", surface)
	}
	if surface.ObservedRef != ".sdp-trace/feedback/event.json" {
		t.Fatalf("feedback observed ref = %s", surface.ObservedRef)
	}
}

func TestDoctorRejectsMalformedFeedbackEventSurface(t *testing.T) {
	repo := initRepo(t)
	writeFileForTest(t, filepath.Join(repo, ".sdp-trace", "feedback", "event.json"), "{}\n")
	status, err := Doctor(Options{
		RepoRoot:     repo,
		Profile:      ProfileGithubActionsGitHooksV1,
		RepositoryID: "demo_repo",
	})
	if err != nil {
		t.Fatal(err)
	}
	surface := surfaceByID(t, status, SurfaceFeedbackEvents)
	if surface.InstallState != StateCannotVerify || surface.ProofState != StateCannotVerify || surface.ReasonCode != ReasonUnsafeOutputRefused {
		t.Fatalf("malformed feedback surface = %+v", surface)
	}
}

func TestBlock28ExampleStatusesUseClosedReasonCodes(t *testing.T) {
	allowed := map[string]bool{
		ReasonHooksPathAbsent:             true,
		ReasonHooksPathSet:                true,
		ReasonHooksPathMismatch:           true,
		ReasonHookScriptAbsent:            true,
		ReasonHookScriptPresent:           true,
		ReasonHookOutputNotObserved:       true,
		ReasonLocalHooksBypassable:        true,
		ReasonAlreadyInstalled:            true,
		ReasonCIWorkflowAbsent:            true,
		ReasonCIWorkflowPresent:           true,
		ReasonCIArtifactUploadAbsent:      true,
		ReasonCIArtifactUploadPresent:     true,
		ReasonCIArtifactBundleNotObserved: true,
		ReasonCIArtifactBundleObserved:    true,
		ReasonFeedbackEventsObserved:      true,
		ReasonFeedbackEventsNotObserved:   true,
		ReasonAgentReportedNotProof:       true,
		ReasonOutsideProfileScope:         true,
		ReasonUnsafeOutputRefused:         true,
		ReasonManualStepRequired:          true,
	}
	for _, rel := range []string{
		"../../examples/block28-repo-observer/valid-installed-local/repo-observer-status.json",
		"../../examples/block28-repo-observer/missing-ci-workflow/repo-observer-status.json",
	} {
		var status Status
		readJSONForTest(t, rel, &status)
		if status.SchemaVersion != SchemaVersion {
			t.Fatalf("%s schema_version = %s", rel, status.SchemaVersion)
		}
		if len(status.Surfaces) != 14 {
			t.Fatalf("%s surfaces = %d", rel, len(status.Surfaces))
		}
		for _, surface := range status.Surfaces {
			if !allowed[surface.ReasonCode] {
				t.Fatalf("%s uses non-closed reason code %q", rel, surface.ReasonCode)
			}
		}
	}
}

func initRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	runGitForTest(t, repo, "init")
	return repo
}

func runGitForTest(t *testing.T, dir string, args ...string) {
	t.Helper()
	_ = gitOutputForTest(t, dir, args...)
}

func gitOutputForTest(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, string(output))
	}
	return string(output)
}

func writeFileForTest(t *testing.T, path string, value string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(value), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func readJSONForTest(t *testing.T, path string, target any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read json: %v", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("parse json: %v", err)
	}
}

func surfaceByID(t *testing.T, status Status, id string) Surface {
	t.Helper()
	for _, surface := range status.Surfaces {
		if surface.SurfaceID == id {
			return surface
		}
	}
	t.Fatalf("surface %s not found in %+v", id, status.Surfaces)
	return Surface{}
}
