package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"os/exec"
	"sort"
	"strings"
)

func registryUsages() ([]string, error) {
	output, err := runCommandSurface()
	if err != nil {
		return nil, err
	}
	return extractUsages(output)
}

func runCommandSurface() ([]byte, error) {
	cmd := exec.Command("go", "run", "./cmd/sdp-trace", "command-surface")
	cmd.Dir = repoRoot()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run sdp-trace command-surface: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return output, nil
}

func extractUsages(output []byte) ([]string, error) {
	var surface struct {
		Commands []struct {
			Usage      string   `json:"usage"`
			Variations []string `json:"variations,omitempty"`
		} `json:"commands"`
	}
	if err := json.Unmarshal(output, &surface); err != nil {
		return nil, fmt.Errorf("parse command surface: %w", err)
	}
	usages := map[string]bool{}
	for _, cmd := range surface.Commands {
		addUsage(usages, cmd.Usage)
		addUsages(usages, cmd.Variations)
	}
	return sortedStringKeys(usages), nil
}

func addUsage(usages map[string]bool, usage string) {
	if usage != "" {
		usages[html.UnescapeString(usage)] = true
	}
}

func addUsages(usages map[string]bool, values []string) {
	for _, v := range values {
		addUsage(usages, v)
	}
}

func sortedStringKeys(m map[string]bool) []string {
	var list []string
	for k := range m {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

func compareRegistryWithDocs(doc string) error {
	registry, err := registryUsages()
	if err != nil {
		return err
	}
	docCommands := documentedCommands(doc)
	registrySet := stringSliceToSet(registry)
	stale := diffStringSliceAgainstSet(docCommands, registrySet, "sdp-trace --help")
	missing := missingFromDocs(registry, docCommands)
	if len(missing) == 0 && len(stale) == 0 {
		return nil
	}
	return fmt.Errorf("registry/docs drift: %s", commandSurfaceDrift(missing, stale))
}

func stringSliceToSet(s []string) map[string]bool {
	set := map[string]bool{}
	for _, v := range s {
		set[v] = true
	}
	return set
}

func diffStringSliceAgainstSet(slice []string, set map[string]bool, skip string) []string {
	var diff []string
	for _, v := range slice {
		if v == skip {
			continue
		}
		if !set[v] {
			diff = append(diff, v)
		}
	}
	return diff
}

func missingFromDocs(registry, docs []string) []string {
	docSet := stringSliceToSet(docs)
	var missing []string
	for _, r := range registry {
		if !docSet[r] {
			missing = append(missing, r)
		}
	}
	return missing
}
