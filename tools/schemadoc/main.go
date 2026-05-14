package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

const (
	indexPath   = "schema/index.json"
	schemaGlob  = "schema/*.schema.json"
	readmePath  = "schema/README.md"
	statusCurrent     = "current"
	statusHistorical  = "historical"
	statusNotAssessed = "not_assessed"
)

var validStatuses = map[string]bool{
	statusCurrent:     true,
	statusHistorical:  true,
	statusNotAssessed: true,
}

func main() {
	var generate bool
	flag.BoolVar(&generate, "generate", false, "print generated README table section to stdout")
	flag.Parse()

	os.Exit(exitCode(run(generate), os.Stderr))
}

func exitCode(err error, stderr io.Writer) int {
	if err == nil {
		return 0
	}
	fmt.Fprintln(stderr, err)
	return 1
}

func run(generate bool) error {
	root := repoRoot()
	idx, err := readIndex(filepath.Join(root, indexPath))
	if err != nil {
		return err
	}

	if generate {
		_, err := io.WriteString(os.Stdout, generateTable(idx))
		return err
	}

	return check(root, idx)
}

func repoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// Index is the machine-readable schema documentation source of truth.
type Index struct {
	Version string        `json:"version"`
	Schemas []SchemaEntry `json:"schemas"`
}

// SchemaEntry describes one schema file.
type SchemaEntry struct {
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	Purpose  string   `json:"purpose"`
	Examples []string `json:"examples,omitempty"`
}

func readIndex(path string) (*Index, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", indexPath, err)
	}
	var idx Index
	if err := json.Unmarshal(b, &idx); err != nil {
		return nil, fmt.Errorf("parse %s: %w", indexPath, err)
	}
	return &idx, nil
}

func listSchemaFiles(root string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(root, schemaGlob))
	if err != nil {
		return nil, err
	}
	var names []string
	for _, m := range matches {
		names = append(names, filepath.Base(m))
	}
	sort.Strings(names)
	return names, nil
}

func check(root string, idx *Index) error {
	files, err := listSchemaFiles(root)
	if err != nil {
		return err
	}

	indexed := make(map[string]*SchemaEntry, len(idx.Schemas))
	for i := range idx.Schemas {
		entry := &idx.Schemas[i]
		indexed[entry.Name] = entry
	}

	var issues []string

	// Missing from index
	for _, f := range files {
		if _, ok := indexed[f]; !ok {
			issues = append(issues, fmt.Sprintf("missing index entry: %s", f))
		}
	}

	// Extra or invalid index entries
	for _, entry := range idx.Schemas {
		path := filepath.Join(root, "schema", entry.Name)
		if _, err := os.Stat(path); err != nil {
			issues = append(issues, fmt.Sprintf("extra index entry (file missing): %s", entry.Name))
			continue
		}

		if entry.Status == "" {
			issues = append(issues, fmt.Sprintf("%s: missing status", entry.Name))
		} else if !validStatuses[entry.Status] {
			issues = append(issues, fmt.Sprintf("%s: invalid status %q", entry.Name, entry.Status))
		}

		if strings.TrimSpace(entry.Purpose) == "" {
			issues = append(issues, fmt.Sprintf("%s: missing purpose", entry.Name))
		}

		for _, ex := range entry.Examples {
			exPath := filepath.Join(root, ex)
			if _, err := os.Stat(exPath); err != nil {
				issues = append(issues, fmt.Sprintf("%s: broken example ref: %s", entry.Name, ex))
			}
		}
	}

	if len(issues) > 0 {
		sort.Strings(issues)
		return fmt.Errorf("schema doc drift:\n- %s", strings.Join(issues, "\n- "))
	}
	return nil
}

func generateTable(idx *Index) string {
	var b strings.Builder
	b.WriteString("| Schema | Status | Purpose |\n")
	b.WriteString("|---|---|---|\n")
	for _, entry := range idx.Schemas {
		b.WriteString(fmt.Sprintf("| `%s` | %s | %s |\n", entry.Name, entry.Status, entry.Purpose))
	}
	return b.String()
}
