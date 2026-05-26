package main

func registryUsages() ([]string, error) {
	output, err := runCommandSurface()
	if err != nil {
		return nil, err
	}
	return extractUsages(output)
}
