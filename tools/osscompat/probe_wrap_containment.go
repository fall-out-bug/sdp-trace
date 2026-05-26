package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func requireResolvedPathUnderTmp(resolvedPath, resolvedTmp string) error {
	if !strings.HasPrefix(resolvedPath, resolvedTmp+string(filepath.Separator)) {
		return fmt.Errorf("run.json resolved outside tmpDir: %s", resolvedPath)
	}
	return nil
}
