package main

func checkActionPins(root string, tracked []string) []string {
	var findings []string
	for _, f := range tracked {
		if !actionFile(f) {
			// Only workflow YAML can contain GitHub Action refs.
			continue
		}
		findings = append(findings, actionPinFindings(root, f)...)
	}
	return findings
}
