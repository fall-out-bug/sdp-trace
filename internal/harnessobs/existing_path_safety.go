package harnessobs

func safeExistingFile(path string) (string, error) {
	// Existing file inputs must already exist inside the working tree; this
	// helper only selects the file-specific policy.
	return safeExistingPath(path, existingPathSpec{
		traversalError: "path must be relative local file without traversal",
		requireDir:     false,
		typeError:      "path must be a file",
	})
}

func safeExistingDir(path string) (string, error) {
	// Existing directory inputs share traversal handling with files but require
	// directory type confirmation after symlink resolution.
	return safeExistingPath(path, existingPathSpec{
		traversalError: "path must be relative local directory without traversal",
		requireDir:     true,
		typeError:      "path must be a directory",
	})
}

func safeExistingPath(path string, spec existingPathSpec) (string, error) {
	// Existing paths are normalized, resolved, converted back to a working-tree
	// relative path, and only then type-checked for the requested kind.
	cleanPath, err := sanitizeExistingPath(path, spec.traversalError)
	if err != nil {
		return "", err
	}

	abs, err := resolveExistingAbsolutePath(cleanPath)
	if err != nil {
		return "", err
	}
	rel, err := relativeWorkingDirectoryPath(abs)
	if err != nil {
		return "", err
	}
	return rel, ensureExpectedPathType(abs, spec)
}
