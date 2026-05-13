package repoobserver

import (
	"fmt"
	"os/exec"
	"strings"
)

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
