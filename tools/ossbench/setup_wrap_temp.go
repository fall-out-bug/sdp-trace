package main

import (
	"fmt"
	"os"
)

func setupWrapTempDir(def *benchmarkDef) error {
	tmpDir, err := os.MkdirTemp("", "ossbench-wrap-*")
	if err != nil {
		return fmt.Errorf("mkdir temp: %w", err)
	}
	def.Dir = tmpDir
	def.Cleanup = func() {
		_ = os.RemoveAll(tmpDir)
	}
	return nil
}
