package witness

const (
	statePass         = "pass"
	stateFail         = "fail"
	stateCannotVerify = "cannot_verify"
	stateNotAssessed  = "not_assessed"

	independenceExternal = "external_independent"
	independenceCIJob    = "ci_isolated_job"
	independenceSameJob  = "ci_same_job"
)

var safetyClasses = []string{
	"ci_token",
	"oidc_token",
	"jwt_body",
	"private_key_material",
	"provider_token",
	"authenticated_provider_url",
	"raw_job_log",
	"private_filesystem_path",
	"unsafe_personal_identifier",
	"free_text_parser_error_with_input",
	"customer_directory_payload",
	"ldap_payload",
	"saml_payload",
	"cloud_payload",
	"vault_payload",
	"hsm_payload",
	"kms_payload",
	"pki_payload",
}

var outputSafetyMarkers = []string{
	"https://user:",
	"https://token:",
	"bearer ",
	"ghp_",
	"glpat-",
	"xoxb-",
	"/private/",
	"raw_job_log_sentinel",
	"oidc.jwt.",
	"customer_directory_payload",
	"ldap_payload",
	"saml_payload",
	"cloud_payload",
	"vault_payload",
	"hsm_payload",
	"kms_payload",
	"pki_payload",
	"free_text_parser_error_with_input",
}

var secretSafetyMarkers = []string{
	"-----begin private key-----",
	"-----begin rsa private key-----",
	"-----begin ec private key-----",
	"token_secret_",
	"jwt_secret_",
	"oidc.jwt.",
	"bearer ",
	"ghp_",
	"glpat-",
	"xoxb-",
	"https://user:",
	"https://token:",
	"raw_job_log_sentinel",
	"customer_directory_payload",
	"ldap_payload",
	"saml_payload",
	"cloud_payload",
	"vault_payload",
	"hsm_payload",
	"kms_payload",
	"pki_payload",
	"free_text_parser_error_with_input",
}
