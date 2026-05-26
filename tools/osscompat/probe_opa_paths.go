package main

import "fmt"

func checkedOPAFixturePaths(regoPath, fixturePath, name string) (string, string, error) {
	if err := requireExistingPath(regoPath); err != nil {
		return "", "", fmt.Errorf("adapter.rego not found: %w", err)
	}
	if err := requireExistingPath(fixturePath); err != nil {
		return "", "", fmt.Errorf("%s not found: %w", name, err)
	}
	return regoPath, fixturePath, nil
}
