package trace

import "path/filepath"

// NewCommandDescriptor returns a descriptor safe for payload retention. The
// raw argv is represented only by argument count and a deterministic digest.
func NewCommandDescriptor(argv []string) CommandDescriptor {
	// The descriptor never stores raw argv; even an empty command keeps the
	// retention policy visible to downstream evidence readers.
	descriptor := CommandDescriptor{
		Argc: len(argv),
		Retention: RetentionDescriptor{
			Mode:        RetentionModeDigestOnly,
			Description: "argv_digest stores a sha256 digest of the command argv; raw argv is not retained",
		},
	}
	if len(argv) == 0 {
		// Empty argv still yields an explicit digest-only retention descriptor.
		return descriptor
	}
	// The executable basename is retained for operator readability; the full
	// argument vector is digest-only.
	descriptor.Executable = filepath.Base(argv[0])
	// The digest binds the complete command line without retaining sensitive
	// argument values in trace payloads.
	descriptor.ArgvDigest = SHA256Hex(mustMarshalJSONStringSlice(argv))
	return descriptor
}
