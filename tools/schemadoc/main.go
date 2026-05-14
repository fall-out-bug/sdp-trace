package main

import (
	"encoding/json"
	"errors"
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
	indexPath        = "schema/index.json"
	schemaGlob       = "schema/*.schema.json"
	readmePath       = "schema/README.md"
	statusCurrent    = "current"
	statusHistorical = "historical"
	statusNotAssessed = "not_assessed"
	examplePresent    = "present"
	exampleNotAssessed = "not_assessed"
	indexVersion      = "1"
)

var validStatuses = map[string]bool{
	statusCurrent:     true,
	statusHistorical:  true,
	statusNotAssessed: true,
}

var validExampleCoverages = map[string]bool{
	examplePresent:    true,
	exampleNotAssessed: true,
	"":                true, // absent means not_assessed
}

func main() {
	var generate, verifyReadme bool
	flag.BoolVar(&generate, "generate", false, "print generated README table section to stdout")
	flag.BoolVar(&verifyReadme, "verify-readme", false, "verify README.md schema table matches index.json")
	flag.Parse()

	os.Exit(exitCode(run(generate, verifyReadme), os.Stderr))
}

func exitCode(err error, stderr io.Writer) int {
	if err == nil {
		return 0
	}
	fmt.Fprintln(stderr, err)
	return 1
}

func run(generate, verifyReadme bool) error {
	root := repoRoot()
	idx, err := readIndex(filepath.Join(root, indexPath))
	if err != nil {
		return err
	}

	if generate {
		_, err := io.WriteString(os.Stdout, generateTable(idx))
		return err
	}

	if verifyReadme {
		return checkReadme(root, idx)
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
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	Purpose         string   `json:"purpose"`
	ExampleCoverage string   `json:"example_coverage,omitempty"`
	Examples        []string `json:"examples,omitempty"`
}

func readIndex(path string) (*Index, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var idx Index
	if err := json.Unmarshal(b, &idx); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if idx.Version != indexVersion {
		return nil, fmt.Errorf("unsupported %s version %q, want %q", path, idx.Version, indexVersion)
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

	var issues []string
	indexed := buildIndexedMap(idx.Schemas, &issues)
	issues = append(issues, findMissingFromIndex(files, indexed)...)
	for i := range idx.Schemas {
		issues = append(issues, validateEntry(root, &idx.Schemas[i])...)
	}

	if len(issues) > 0 {
		sort.Strings(issues)
		return fmt.Errorf("schema doc drift:\n- %s", strings.Join(issues, "\n- "))
	}
	return nil
}

func buildIndexedMap(entries []SchemaEntry, issues *[]string) map[string]*SchemaEntry {
	indexed := make(map[string]*SchemaEntry, len(entries))
	for i := range entries {
		entry := &entries[i]
		if _, exists := indexed[entry.Name]; exists {
			*issues = append(*issues, fmt.Sprintf("duplicate index entry: %s", entry.Name))
			continue
		}
		indexed[entry.Name] = entry
	}
	return indexed
}

func findMissingFromIndex(files []string, indexed map[string]*SchemaEntry) []string {
	var issues []string
	for _, f := range files {
		if _, ok := indexed[f]; !ok {
			issues = append(issues, fmt.Sprintf("missing index entry: %s", f))
		}
	}
	return issues
}

func validateEntry(root string, entry *SchemaEntry) []string {
	var issues []string
	if !strings.HasSuffix(entry.Name, ".schema.json") {
		issues = append(issues, fmt.Sprintf("%s: name must end with .schema.json", entry.Name))
	}
	issues = append(issues, validateFilePresence(root, entry)...)
	issues = append(issues, validateEntryStatus(entry)...)
	issues = append(issues, validateEntryPurpose(entry)...)
	issues = append(issues, validateExampleCoverage(entry)...)
	issues = append(issues, validateExampleRefs(root, entry)...)
	return issues
}

func validateFilePresence(root string, entry *SchemaEntry) []string {
	path := filepath.Join(root, "schema", entry.Name)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return []string{fmt.Sprintf("extra index entry (file missing): %s", entry.Name)}
		}
		return []string{fmt.Sprintf("%s: cannot access file: %v", entry.Name, err)}
	}
	return nil
}

func validateEntryStatus(entry *SchemaEntry) []string {
	if entry.Status == "" {
		return []string{fmt.Sprintf("%s: missing status", entry.Name)}
	}
	if !validStatuses[entry.Status] {
		return []string{fmt.Sprintf("%s: invalid status %q", entry.Name, entry.Status)}
	}
	return nil
}

func validateEntryPurpose(entry *SchemaEntry) []string {
	if strings.TrimSpace(entry.Purpose) == "" {
		return []string{fmt.Sprintf("%s: missing purpose", entry.Name)}
	}
	return nil
}

func validateExampleCoverage(entry *SchemaEntry) []string {
	var issues []string
	if !validExampleCoverages[entry.ExampleCoverage] {
		issues = append(issues, fmt.Sprintf("%s: invalid example_coverage %q", entry.Name, entry.ExampleCoverage))
	}
	if entry.ExampleCoverage == examplePresent && len(entry.Examples) == 0 {
		issues = append(issues, fmt.Sprintf("%s: example_coverage is present but examples list is empty", entry.Name))
	}
	if entry.ExampleCoverage != examplePresent && len(entry.Examples) > 0 {
		issues = append(issues, fmt.Sprintf("%s: example_coverage is %q but examples list is non-empty", entry.Name, entry.ExampleCoverage))
	}
	return issues
}

func validateExampleRefs(root string, entry *SchemaEntry) []string {
	var issues []string
	for _, ex := range entry.Examples {
		exPath := filepath.Join(root, ex)
		if _, err := os.Stat(exPath); err != nil {
			if os.IsNotExist(err) {
				issues = append(issues, fmt.Sprintf("%s: broken example ref: %s", entry.Name, ex))
			} else {
				issues = append(issues, fmt.Sprintf("%s: cannot access example %s: %v", entry.Name, ex, err))
			}
		}
	}
	return issues
}

func generateTable(idx *Index) string {
	var b strings.Builder
	b.WriteString("| Schema | Status | Purpose |\n")
	b.WriteString("|---|---|---|\n")
	for _, entry := range idx.Schemas {
		name := strings.ReplaceAll(entry.Name, "|", "\\|")
		purpose := strings.ReplaceAll(entry.Purpose, "|", "\\|")
		b.WriteString(fmt.Sprintf("| `%s` | %s | %s |\n", name, entry.Status, purpose))
	}
	return b.String()
}

func checkReadme(root string, idx *Index) error {
	return checkReadmeAt(filepath.Join(root, readmePath), idx)
}

func checkReadmeAt(path string, idx *Index) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	content := string(data)
	const startMarker = "<!-- schemadoc-start -->\n"
	const endMarker = "\n<!-- schemadoc-end -->"
	startIdx := strings.Index(content, startMarker)
	endIdx := strings.Index(content, endMarker)
	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		return errors.New("README missing schemadoc markers")
	}
	want := generateTable(idx)
	got := content[startIdx+len(startMarker) : endIdx]
	if strings.TrimSpace(got) != strings.TrimSpace(want) {
		return fmt.Errorf("README schema table drift: run 'go run ./tools/schemadoc -generate' and update the section between %s and %s", startMarker, endMarker)
	}
	return nil
}
