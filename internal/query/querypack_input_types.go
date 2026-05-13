package query

type packInputs struct {
	run              runArtifact
	runArtifact      QueryPackInputArtifact
	runErr           error
	forensic         assessmentEnvelope
	forensicArtifact *QueryPackInputArtifact
	forensicPresent  bool
	forensicErr      error
	adapter          assessmentEnvelope
	adapterArtifact  *QueryPackInputArtifact
	adapterPresent   bool
	adapterErr       error
}
