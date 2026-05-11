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
		!strings.Contains(out.String(), "sdp-trace packet check-demo --bundle <file>") ||
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

func TestPacketCheckDemoCLIRequiresFirstPacketRouteEvidence(t *testing.T) {
	root := t.TempDir()
	bundle := validPacketBundleForCLI()
	bundlePath := filepath.Join(root, "bundle.json")
	writeTestJSON(t, bundlePath, bundle)
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"packet", "check-demo", "--bundle", bundlePath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("check-demo exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"state": "pass"`) {
		t.Fatalf("check-demo output missing pass: %s", out.String())
	}

	setPacketCLIRowState(&bundle, "PC-AGENT-ROUTE", packet.StateNotAssessed, nil, "route evidence missing")
	bundle.Manifest.PacketDigest = packet.PacketDigest(bundle.Packet)
	writeTestJSON(t, bundlePath, bundle)
	out.Reset()
	errOut.Reset()
	exit = run([]string{"packet", "check-demo", "--bundle", bundlePath}, &out, &errOut)
	if exit != exitFail {
		t.Fatalf("check-demo missing route exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "PC-AGENT-ROUTE must be pass or partial") {
		t.Fatalf("check-demo output missing route error: %s", out.String())
	}
}

func TestPacketCheckDemoCLIRejectsWeakDemoProof(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name  string
		edit  func(*packet.Bundle)
		error string
	}{
		{
			name: "hand authored",
			edit: func(bundle *packet.Bundle) {
				bundle.Packet.AuthoringMethod = "hand_authored_before_tooling"
			},
			error: "tool_generated authoring_method",
		},
		{
			name: "generic harness route",
			edit: func(bundle *packet.Bundle) {
				setPacketCLIEntryResolver(bundle, "harness:route", "generic harness observation artifact")
				setPacketCLIEntryRouteObservation(bundle, "harness:route", "", nil)
			},
			error: "structured OpenCode/GSD/MiniMax harness route observation",
		},
		{
			name: "expired change evidence access",
			edit: func(bundle *packet.Bundle) {
				setPacketCLIEntryAccess(bundle, "git:change", "expired")
			},
			error: "PC-CHANGE evidence ref",
		},
		{
			name: "expired partial change evidence",
			edit: func(bundle *packet.Bundle) {
				setPacketCLIRowState(bundle, "PC-CHANGE", packet.StatePartial, []string{"git:change"}, "change evidence is partial")
				setPacketCLIEntryExpiresAt(bundle, "git:change", "2026-05-10T12:00:00Z")
			},
			error: "PC-CHANGE evidence ref",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bundle := validPacketBundleForCLI()
			tt.edit(&bundle)
			bundle.Manifest.PacketDigest = packet.PacketDigest(bundle.Packet)
			bundlePath := filepath.Join(root, strings.ReplaceAll(tt.name, " ", "-")+".json")
			writeTestJSON(t, bundlePath, bundle)
			var out bytes.Buffer
			var errOut bytes.Buffer
			exit := run([]string{"packet", "check-demo", "--bundle", bundlePath}, &out, &errOut)
			if exit != exitFail {
				t.Fatalf("check-demo exit %d err=%s out=%s", exit, errOut.String(), out.String())
			}
			if !strings.Contains(out.String(), tt.error) {
				t.Fatalf("check-demo output missing %q: %s", tt.error, out.String())
			}
		})
	}
}

func TestPacketCheckDemoCLIRejectsMinimumBarRegressions(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name  string
		edit  func(*packet.Bundle)
		error string
	}{
		{
			name: "too few pass partial rows",
			edit: func(bundle *packet.Bundle) {
				for _, id := range []string{"PC-INITIATOR", "PC-VERIFICATION", "PC-REVIEW", "PC-THEATER", "PC-DECISION", "PC-RESIDUAL-GAPS"} {
					setPacketCLIRowState(bundle, id, packet.StateNotAssessed, nil, "not assessed in regression fixture")
					addPacketCLIGap(bundle, id, packet.StateNotAssessed, "not assessed in regression fixture")
				}
			},
			error: "at least 4 pass or partial rows",
		},
		{
			name: "missing change evidence",
			edit: func(bundle *packet.Bundle) {
				setPacketCLIRowState(bundle, "PC-CHANGE", packet.StatePartial, nil, "missing retained change evidence")
			},
			error: "PC-CHANGE retained evidence refs",
		},
		{
			name: "too many unclosed cannot verify rows",
			edit: func(bundle *packet.Bundle) {
				setPacketCLIRowState(bundle, "PC-AUTHORITY", packet.StateCannotVerify, nil, "cannot verify authority")
				setPacketCLIRowState(bundle, "PC-ATTESTATION", packet.StateCannotVerify, nil, "cannot verify attestation")
				bundle.Packet.ResidualGaps = withoutPacketCLIGaps(bundle.Packet.ResidualGaps, "PC-AUTHORITY", "PC-ATTESTATION")
			},
			error: "at most one cannot_verify row without closure path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bundle := validPacketBundleForCLI()
			tt.edit(&bundle)
			bundle.Manifest.PacketDigest = packet.PacketDigest(bundle.Packet)
			bundlePath := filepath.Join(root, strings.ReplaceAll(tt.name, " ", "-")+".json")
			writeTestJSON(t, bundlePath, bundle)
			var out bytes.Buffer
			var errOut bytes.Buffer
			exit := run([]string{"packet", "check-demo", "--bundle", bundlePath}, &out, &errOut)
			if exit != exitFail {
				t.Fatalf("check-demo exit %d err=%s out=%s", exit, errOut.String(), out.String())
			}
			if !strings.Contains(out.String(), tt.error) {
				t.Fatalf("check-demo output missing %q: %s", tt.error, out.String())
			}
		})
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
	for i := range entries {
		if entries[i].Ref == "harness:route" {
			entries[i].Resolver = "OpenCode/GSD MiniMax-M2.5 observation artifact"
			entries[i].EvidenceKind = "harness_route_observation"
			entries[i].ObservedComponents = []string{"opencode", "gsd", "minimax-m2.5"}
		}
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

func setPacketCLIRowState(bundle *packet.Bundle, id, state string, refs []string, reason string) {
	for i := range bundle.Packet.Rows {
		if bundle.Packet.Rows[i].ID == id {
			bundle.Packet.Rows[i].State = state
			bundle.Packet.Rows[i].EvidenceRefs = refs
			bundle.Packet.Rows[i].Reason = reason
			return
		}
	}
}

func setPacketCLIEntryResolver(bundle *packet.Bundle, ref, resolver string) {
	for i := range bundle.Manifest.Entries {
		if bundle.Manifest.Entries[i].Ref == ref {
			bundle.Manifest.Entries[i].Resolver = resolver
			return
		}
	}
}

func setPacketCLIEntryAccess(bundle *packet.Bundle, ref, access string) {
	for i := range bundle.Manifest.Entries {
		if bundle.Manifest.Entries[i].Ref == ref {
			bundle.Manifest.Entries[i].ArtifactAccess = access
			return
		}
	}
}

func setPacketCLIEntryExpiresAt(bundle *packet.Bundle, ref, expiresAt string) {
	for i := range bundle.Manifest.Entries {
		if bundle.Manifest.Entries[i].Ref == ref {
			bundle.Manifest.Entries[i].ExpiresAt = expiresAt
			return
		}
	}
}

func setPacketCLIEntryRouteObservation(bundle *packet.Bundle, ref, evidenceKind string, components []string) {
	for i := range bundle.Manifest.Entries {
		if bundle.Manifest.Entries[i].Ref == ref {
			bundle.Manifest.Entries[i].EvidenceKind = evidenceKind
			bundle.Manifest.Entries[i].ObservedComponents = components
			return
		}
	}
}

func withoutPacketCLIGaps(gaps []packet.ResidualGap, rowIDs ...string) []packet.ResidualGap {
	remove := map[string]bool{}
	for _, rowID := range rowIDs {
		remove[rowID] = true
	}
	filtered := []packet.ResidualGap{}
	for _, gap := range gaps {
		if !remove[gap.RowID] {
			filtered = append(filtered, gap)
		}
	}
	return filtered
}

func addPacketCLIGap(bundle *packet.Bundle, rowID, state, reason string) {
	if rowID == "PC-RESIDUAL-GAPS" {
		return
	}
	for _, gap := range bundle.Packet.ResidualGaps {
		if gap.RowID == rowID {
			return
		}
	}
	bundle.Packet.ResidualGaps = append(bundle.Packet.ResidualGaps, packet.ResidualGap{
		RowID:           rowID,
		State:           state,
		Reason:          reason,
		ClosureEvidence: "regression fixture closure evidence",
	})
}
