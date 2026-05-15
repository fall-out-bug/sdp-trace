package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// checkRootExecutables returns findings for root files that have the
// executable bit set in the git index or that match a known binary format.
func checkRootExecutables(root string, tracked []string) []string {
	var findings []string
	for _, f := range tracked {
		if finding := rootExecutableFinding(root, f); finding != "" {
			findings = append(findings, finding)
		}
	}
	return findings
}

// rootExecutableFinding checks one root file. The binary test is scoped to
// root only because compiled artifacts under cmd/ or internal/ are legitimate
// tracked files.
func rootExecutableFinding(root, f string) string {
	if strings.Contains(f, "/") {
		return ""
	}
	mode := fileModeInIndex(root, f)
	if mode == "100755" {
		return fmt.Sprintf("root executable: %s", f)
	}
	if isBinaryFile(filepath.Join(root, f)) {
		return fmt.Sprintf("root binary artifact: %s", f)
	}
	return ""
}

// fileModeInIndex returns the git index mode for path, or empty string on
// error. The expected output format is "<mode> <object> <stage>\t<file>".
func fileModeInIndex(root, path string) string {
	cmd := exec.Command("git", "ls-files", "--stage", path)
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(out))
	if len(fields) >= 1 {
		return fields[0]
	}
	return ""
}
