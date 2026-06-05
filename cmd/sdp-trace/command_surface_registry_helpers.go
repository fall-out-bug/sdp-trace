package main

// Command-surface registry helpers grouped from transitional numbered shards.
func commandSurfaceCoreCommands() []commandSurfaceCmd {
	return commandSurfaceCore
}

func commandSurfaceObserveCommands() []commandSurfaceCmd {
	return commandSurfaceObserveGroup
}

func commandSurfaceAssessCommands() []commandSurfaceCmd {
	return commandSurfaceAssessGroup
}

func commandSurfaceOtherCommands() []commandSurfaceCmd {
	return commandSurfaceOtherGroup
}

func commandSurfacePacketCommands() []commandSurfaceCmd {
	return commandSurfacePacketGroup
}

func buildCommandSurface() commandSurfaceSchema {
	return commandSurfaceRegistry
}
