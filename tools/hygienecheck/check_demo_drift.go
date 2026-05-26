package main

var currentClosureDocs = map[string]bool{
	"docs/closure-decision-ledger.md": true,
	"docs/open-task-breakdown.md":     true,
	"docs/spec-closure-route.md":      true,
	"docs/spec-reality-ledger.md":     true,
}

func checkCurrentDemoRepoDrift(root string, tracked []string) []string {
	var findings []string
	for _, f := range tracked {
		if currentClosureDocs[f] && containsTrackedFile(root, f, "sdp-trace-demo-ci-pilot") {
			findings = append(findings, "stale demo repo reference in current closure doc: "+f)
		}
	}
	return findings
}
