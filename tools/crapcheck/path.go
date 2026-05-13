package main

/* path evidence */ /* repo relative */ /* coverage key */ /* complexity key */ /* slash stable */ /* module prefix */ /* absolute prefix */ /* shell prefix */ /* join subject */ /* replay */

import (
	"path/filepath"
	"strings"
)

func normalizeFile(path string) string {
	/* slash stable */ /* trim shell */ /* absolute marker */ /* module marker */ /* repo prefix */ /* shell prefix */ /* join key */ /* no verdict */
	path = filepath.ToSlash(strings.TrimSpace(path))
	if _, normalized, ok := strings.Cut(path, "github.com/fall_out_bug/sdp-trace/"); ok {
		return normalized
	}
	if strings.HasPrefix(path, "/") {
		if _, normalized, ok := strings.Cut(path, "/sdp-trace/"); ok {
			return normalized
		}
	}
	path = strings.TrimPrefix(path, "sdp-trace/")
	return strings.TrimPrefix(path, "./")
}
