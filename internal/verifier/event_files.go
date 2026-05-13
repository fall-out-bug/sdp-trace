package verifier

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func loadRunEventFiles(eventDir string) ([]string, error) {
	files, err := os.ReadDir(eventDir)
	if err != nil {
		return nil, err
	}

	// Directory order is not trusted. The caller sorts the accepted filenames
	// before replay so filesystem enumeration cannot influence the chain.
	eventFiles := make([]string, 0, len(files))
	for _, entry := range files {
		eventFiles = appendRunEventFile(eventFiles, eventDir, entry)
	}
	return eventFiles, nil
}

func appendRunEventFile(eventFiles []string, eventDir string, entry os.DirEntry) []string {
	if !isValidEventFile(entry) {
		// Ignore non-event files so auxiliary artifacts cannot affect replay
		// order.
		return eventFiles
	}
	return append(eventFiles, filepath.Join(eventDir, entry.Name()))
}

func isValidEventFile(entry os.DirEntry) bool {
	if entry.IsDir() {
		// Event chains are flat files only.
		return false
	}
	name := entry.Name()
	if !strings.HasSuffix(name, ".json") {
		// Only JSON event files participate in chain replay.
		return false
	}
	_, err := trace.EventSeqFromFilename(name)
	// Filename sequence parsing rejects unrelated JSON files.
	return err == nil
}
