package main

var commandSurfaceStates = []stateMeta{
	{Name: "observed", Description: "Verifier evidence met required checks for the selected local profile."},
	{Name: "pass", Description: "Selected profile concluded successfully where the command contract uses pass/fail states."},
	{Name: "fail", Description: "Verifier evidence conflicted or was insufficient for required checks."},
	{Name: "not_assessed", Description: "State was outside the run scope; it does not imply success or evidence."},
	{Name: "cannot_verify", Description: "Verifier could not execute a required check or lacked required evidence."},
}
