package main

import "path/filepath"

func parentDir(dir string) string {
	parent := filepath.Dir(dir)
	if parent == dir {
		return ""
	}
	return parent
}
