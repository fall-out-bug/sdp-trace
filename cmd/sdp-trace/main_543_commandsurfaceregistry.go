package main

func buildCommandSurface() commandSurfaceSchema {
	var cmds []commandSurfaceCmd
	cmds = append(cmds, commandSurfaceCoreCommands()...)
	cmds = append(cmds, commandSurfaceObserveCommands()...)
	cmds = append(cmds, commandSurfaceAssessCommands()...)
	cmds = append(cmds, commandSurfaceOtherCommands()...)
	cmds = append(cmds, commandSurfacePacketCommands()...)
	return commandSurfaceSchema{
		SchemaVersion: commandSurfaceSchemaVersion,
		Commands:      cmds,
		Profiles:      knownAssessmentProfiles(),
		WitnessKinds:  knownWitnessKinds(),
		States:        knownStates(),
	}
}
