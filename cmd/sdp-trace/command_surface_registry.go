package main

import "slices"

var commandSurfaceRegistry = commandSurfaceSchema{
	SchemaVersion: commandSurfaceSchemaVersion,
	Commands: slices.Concat(
		commandSurfaceCoreCommands(),
		commandSurfaceObserveCommands(),
		commandSurfaceAssessCommands(),
		commandSurfaceOtherCommands(),
		commandSurfacePacketCommands(),
	),
	Profiles:     knownAssessmentProfiles(),
	WitnessKinds: knownWitnessKinds(),
	States:       knownStates(),
}
