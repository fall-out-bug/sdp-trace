package main

type commandSurface struct {
	Commands []commandSurfaceCommand `json:"commands"`
}

type commandSurfaceCommand struct {
	Usage      string   `json:"usage"`
	Variations []string `json:"variations,omitempty"`
}
