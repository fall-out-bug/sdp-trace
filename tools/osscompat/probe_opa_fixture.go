package main

import (
	"path/filepath"
)

const ossPolicyDir = "examples/oss-policy"

func opaFixturePaths(name string) (regoPath, fixturePath string, err error) {
	root := filepath.Join(repoRoot(), ossPolicyDir)
	regoPath = filepath.Join(root, "adapter.rego")
	fixturePath = filepath.Join(root, name)
	return checkedOPAFixturePaths(regoPath, fixturePath, name)
}
