package ciartifact

func Evaluate(manifest Manifest) ObservationResult {
	// Evaluation starts from the manifest inputs and produces evidence states,
	// not a scalar health score or an implicit pass claim.
	// Each later helper keeps missing, unsafe, and contradictory evidence separate.

	source, sourceSafe := sanitizeSource(manifest.SelectedSource)
	run, runSafe := sanitizeRun(manifest.SelectedRun)
	inputs := evaluatedManifestInputs(manifest)
	identityCannotVerify := !sourceSafe || !runSafe
	state := topLevel(inputs.families, inputs.index, inputs.safety, len(inputs.reqs), identityCannotVerify)
	return observationResult(manifest, source, run, inputs, state, identityCannotVerify)
}
