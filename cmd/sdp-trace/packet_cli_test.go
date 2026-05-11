package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func TestPacketValidateAndRenderCLI(t *testing.T) {
	root := t.TempDir()
	bundlePath := filepath.Join(root, "bundle.json")
	writeTestJSON(t, bundlePath, validPacketBundleForCLI())
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"packet", "validate", "--bundle", bundlePath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("validate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"state": "pass"`) {
		t.Fatalf("validate output missing pass: %s", out.String())
	}

	out.Reset()
	errOut.Reset()
	packetPath := filepath.Join(root, "change-evidence-packet.md")
	exit = run([]string{"packet", "render", "--bundle", bundlePath, "--out", packetPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("render exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	raw, err := os.ReadFile(packetPath)
	if err != nil {
		t.Fatalf("read rendered packet: %v", err)
	}
	if !strings.Contains(string(raw), "## Required Rows") ||
		!strings.Contains(string(raw), "| PC-THEATER | pass | No P0 theater finding triggered.") {
		t.Fatalf("rendered packet missing required content:\n%s", string(raw))
	}
}

func TestPacketCommandAppearsInTopLevelHelp(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"--help"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("help exit %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "sdp-trace packet validate --bundle <file>") ||
		!strings.Contains(out.String(), "sdp-trace packet build-github --github-input <file> --out <file>") ||
		!strings.Contains(out.String(), "sdp-trace packet render --bundle <file> --out <file>") {
		t.Fatalf("help missing packet commands:\n%s", out.String())
	}
}

func TestPacketBuildGitHubCLI(t *testing.T) {
	root := repoRootForTest(t)
	inputPath := filepath.Join(root, "examples", "change-evidence-packet", "github-input.json")
	outPath := filepath.Join(t.TempDir(), "bundle.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"packet", "build-github", "--github-input", inputPath, "--out", outPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("build-github exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	var bundle packet.Bundle
	readTestJSON(t, outPath, &bundle)
	if bundle.Packet.AuthoringMethod != packet.AuthoringToolGenerated ||
		bundle.Packet.BundleRef != bundle.Manifest.BundleID ||
		bundle.Manifest.PacketDigest == "" {
		t.Fatalf("generated bundle not bound/tool-generated: %+v", bundle)
	}
}

func TestPacketCLIValidatesCommittedHappyPathFixture(t *testing.T) {
	root := repoRootForTest(t)
	bundlePath := filepath.Join(root, "examples", "change-evidence-packet", "happy-path.bundle.json")
	packetPath := filepath.Join(t.TempDir(), "packet.md")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"packet", "validate", "--bundle", bundlePath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("validate fixture exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"packet", "render", "--bundle", bundlePath, "--out", packetPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("render fixture exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	rendered, err := os.ReadFile(packetPath)
	if err != nil {
		t.Fatalf("read rendered fixture packet: %v", err)
	}
	if !strings.Contains(string(rendered), "PC-THEATER | pass") {
		t.Fatalf("rendered fixture missing clean theater row:\n%s", string(rendered))
	}
}

func TestPacketValidateCLIRejectsInvalidBundle(t *testing.T) {
	root := t.TempDir()
	bundle := validPacketBundleForCLI()
	bundle.Packet.Rows[0].EvidenceRefs = []string{"missing:ref"}
	bundlePath := filepath.Join(root, "bundle.json")
	writeTestJSON(t, bundlePath, bundle)
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"packet", "validate", "--bundle", bundlePath}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("validate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "missing:ref") {
		t.Fatalf("validate output missing error: %s", out.String())
	}
}

func repoRootForTest(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func validPacketBundleForCLI() packet.Bundle {
	refs := map[string]string{
		"git:change":        "git",
		"task:pr-body":      "change_host",
		"harness:route":     "harness",
		"ci:run":            "ci",
		"review:packet":     "review",
		"authority:profile": "manual",
		"theater:clean":     "witness",
		"gap:residual":      "manual",
		"decision:owner":    "manual",
	}
	entries := []packet.BundleEntry{}
	for ref, sourceClass := range refs {
		entries = append(entries, packet.BundleEntry{
			Ref:             ref,
			SourceClass:     sourceClass,
			Digest:          "sha256:1111111111111111111111111111111111111111111111111111111111111111",
			RetainedForm:    "raw",
			RedactionStatus: "not_needed",
			Resolver:        "examples/change-evidence-packet/" + ref,
		})
	}
	bundle := packet.Bundle{
		Packet: packet.Packet{
			PacketVersion:   packet.PacketSchemaVersion,
			PacketID:        "packet-cli-test",
			SourceChange:    packet.SourceChange{Repository: "fall-out-bug/sdp-trace", ChangeID: "PR-38", URL: "https://github.com/fall-out-bug/sdp-trace/pull/38"},
			GeneratedAt:     "2026-05-11T12:00:00Z",
			AuthoringMethod: packet.AuthoringToolGenerated,
			SelectedProfile: "change-host-rich-v0",
			RedactionPolicy: "not_needed",
			BundleRef:       "bundle-cli-test",
			PacketState:     "review_ready",
			Projection:      packet.Projection{Kind: packet.ProjectionCanonical, Canonical: true, ArtifactRef: "packet:markdown"},
			Rows: []packet.Row{
				packetCLIRow("PC-CHANGE", packet.StatePass, "Change metadata is retained.", []string{"git:change"}),
				packetCLIRow("PC-INITIATOR", packet.StatePartial, "PR body task source is retained.", []string{"task:pr-body"}),
				packetCLIRow("PC-AGENT-ROUTE", packet.StatePartial, "OpenCode/GSD route is retained.", []string{"harness:route"}),
				packetCLIRow("PC-MUTATION", packet.StatePass, "Commit range and changed files are retained.", []string{"git:change"}),
				packetCLIRow("PC-VERIFICATION", packet.StatePass, "CI run evidence is retained.", []string{"ci:run"}),
				packetCLIRow("PC-REVIEW", packet.StatePartial, "Review evidence is retained.", []string{"review:packet"}),
				packetCLIRow("PC-AUTHORITY", packet.StateNotAssessed, "Authority was not assessed.", nil),
				packetCLIRow("PC-THEATER", packet.StatePass, "No P0 theater finding triggered.", []string{"theater:clean"}),
				packetCLIRow("PC-ATTESTATION", packet.StateNotAssessed, "Signed or external attestation was not assessed.", nil),
				packetCLIRow("PC-DECISION", packet.StatePass, "Decision owner is retained.", []string{"decision:owner"}),
				packetCLIRow("PC-RESIDUAL-GAPS", packet.StatePartial, "Partial rows remain explicit.", []string{"gap:residual"}),
			},
			ResidualGaps: []packet.ResidualGap{
				{RowID: "PC-INITIATOR", State: packet.StatePartial, Reason: "PR body is weaker than issue binding", ClosureEvidence: "GitHub issue"},
				{RowID: "PC-AGENT-ROUTE", State: packet.StatePartial, Reason: "route is observation-only", ClosureEvidence: "complete route chain"},
				{RowID: "PC-REVIEW", State: packet.StatePartial, Reason: "review is external", ClosureEvidence: "GitHub review"},
				{RowID: "PC-AUTHORITY", State: packet.StateNotAssessed, Reason: "authority profile is outside fixture", ClosureEvidence: "authority profile"},
				{RowID: "PC-ATTESTATION", State: packet.StateNotAssessed, Reason: "signed trust is outside fixture", ClosureEvidence: "signed packet"},
			},
			DecisionOwners: []packet.DecisionOwner{
				{Decision: "merge", Owner: "maintainer", State: packet.StateNotAssessed, Reason: "packet is not approval"},
				{Decision: "release", Owner: "release owner", State: packet.StateNotAssessed, Reason: "packet is not release approval"},
				{Decision: "risk_acceptance", Owner: "risk owner", State: packet.StateNotAssessed, Reason: "packet is not risk acceptance"},
				{Decision: "security_review", Owner: "security owner", State: packet.StateNotAssessed, Reason: "packet is not security review"},
			},
			NonApproval: "This packet does not approve merge, release, compliance, production trust, semantic correctness, or signed external trust.",
		},
		Manifest: packet.BundleManifest{
			SchemaVersion: packet.BundleSchemaVersion,
			BundleID:      "bundle-cli-test",
			Entries:       entries,
		},
	}
	bundle.Manifest.PacketDigest = packet.PacketDigest(bundle.Packet)
	return bundle
}

func packetCLIRow(id, state, summary string, refs []string) packet.Row {
	row := packet.Row{ID: id, State: state, Summary: summary, EvidenceRefs: refs, Owner: "maintainer"}
	if state != packet.StatePass {
		row.Reason = "non-pass state is explicit"
	}
	return row
}
