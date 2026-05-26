package main

import "fmt"

func compareRegistryWithDocs(doc string, registry []string) error {
	docCommands := documentedCommands(doc)
	registrySet := stringSliceToSet(registry)
	stale := diffStringSliceAgainstSet(docCommands, registrySet, "sdp-trace --help")
	missing := missingFromDocs(registry, docCommands)
	if len(missing) == 0 && len(stale) == 0 {
		return nil
	}
	return fmt.Errorf("registry/docs drift: %s", commandSurfaceDrift(missing, stale))
}
