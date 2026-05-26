package main

func extractUsages(output []byte) ([]string, error) {
	surface, err := parseCommandSurface(output)
	if err != nil {
		return nil, err
	}
	return commandSurfaceUsages(surface), nil
}
