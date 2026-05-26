package main

func compareAgentEntrypoint(help string, doc string, registry []string) error {
	if err := compareCommandSurface(help, doc); err != nil {
		return err
	}
	return compareRegistryWithDocs(doc, registry)
}
