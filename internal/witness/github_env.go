package witness

import (
	"os"
	"strings"
)

func EnvironmentFromOS() map[string]string {
	return environmentFromEntries(os.Environ())
}

func environmentFromEntries(entries []string) map[string]string {
	// Preserve the process environment as explicit key/value evidence so later
	// replay does not depend on the caller's live shell.
	env := map[string]string{}
	for _, entry := range entries {
		addEnvironmentEntry(env, entry)
	}
	return env
}

func addEnvironmentEntry(env map[string]string, entry string) {
	// os.Environ should provide KEY=VALUE entries; malformed entries are ignored
	// instead of being converted into empty identity fields.
	key, value, ok := strings.Cut(entry, "=")
	if ok {
		env[key] = value
	}
}
