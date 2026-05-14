package main

// readme.go generates the Markdown schema table and verifies that README.md
// stays synchronized with index.json via <!-- schemadoc-start/end --> markers.
// The markers allow CI to detect drift between the human-readable table and
// the machine-readable index.json source of truth.
//
// Table generation is deterministic: rows follow the index order and pipe
// characters in purpose strings are escaped so Markdown remains valid.
//
// checkReadmeAt looks for literal marker strings rather than parsing Markdown
// so the tool stays simple and does not depend on a Markdown parser.
//
// The comparison uses strings.TrimSpace on both sides so that harmless
// whitespace changes do not cause false drift failures.
//
// If the markers are missing or out of order the error tells the user exactly
// which marker is wrong so the fix is obvious.

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateTable produces a deterministic Markdown table from the index entries.
func generateTable(idx *Index) string {
	b := strings.Builder{}
	b.WriteString("| Schema | Status | Purpose |\n|---|---|---|\n")
	for _, e := range idx.Schemas {
		n := strings.ReplaceAll(e.Name, "|", "\\|")
		p := strings.ReplaceAll(e.Purpose, "|", "\\|")
		fmt.Fprintf(&b, "| `%s` | %s | %s |\n", n, e.Status, p)
	}
	return b.String()
}

// checkReadme verifies that the README.md table section matches the current index.
func checkReadme(root string, idx *Index) error {
	return checkReadmeAt(filepath.Join(root, readmePath), idx)
}

// checkReadmeAt extracts the table between markers and compares it with the generated table.
// It returns an error if markers are missing or if the table does not match the index.
func checkReadmeAt(path string, idx *Index) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	c := string(data)
	start := "<!-- schemadoc-start -->\n"
	end := "\n<!-- schemadoc-end -->"
	si := strings.Index(c, start)
	ei := strings.Index(c, end)
	if err := checkMarkers(si, ei); err != nil {
		return err
	}
	want := generateTable(idx)
	got := c[si+len(start) : ei]
	if strings.TrimSpace(got) != strings.TrimSpace(want) {
		return fmt.Errorf("README schema table drift")
	}
	return nil
}

// checkMarkers validates that both markers exist and are in the correct order.
func checkMarkers(startIdx, endIdx int) error {
	if startIdx < 0 {
		return errors.New("README missing schemadoc-start marker")
	}
	if endIdx < 0 {
		return errors.New("README missing schemadoc-end marker")
	}
	if endIdx <= startIdx {
		return errors.New("invalid schemadoc marker order")
	}
	return nil
}
