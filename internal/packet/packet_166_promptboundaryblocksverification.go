package packet

func promptBoundaryBlocksVerification(required bool, classification PromptBoundaryClassification) bool {

	return required && (classification.RouteProofEffect == StateFail || classification.RouteProofEffect == StateCannotVerify)
}
