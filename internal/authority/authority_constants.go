package authority

const (
	PackageSchemaVersion = "authority-observation-package-v1"
	ResultSchemaVersion  = "authority-evaluation-result-v1"
	Profile              = "authority-envelope"

	StateWithinAuthority  = "within_authority"
	StateOutsideAuthority = "outside_authority"
	StateNotAssessed      = "not_assessed"
	StateCannotVerify     = "cannot_verify"

	AttributionVerified     = "verified"
	AttributionNotAssessed  = "not_assessed"
	AttributionCannotVerify = "cannot_verify"

	BindingVerified     = "verified"
	BindingNotAssessed  = "not_assessed"
	BindingCannotVerify = "cannot_verify"
)
