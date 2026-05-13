package packet

var redactionStatuses = map[string]bool{

	"not_needed":      true,
	"redacted":        true,
	"digest_only":     true,
	"withheld":        true,
	StateCannotVerify: true,
}
