package main

func commandSurfaceUsages(surface commandSurface) []string {
	usages := map[string]bool{}
	for _, cmd := range surface.Commands {
		addUsage(usages, cmd.Usage)
		addUsages(usages, cmd.Variations)
	}
	return sortedStringKeys(usages)
}
