package packet

var retainedForms = map[string]bool{

	"raw":          true,
	"redacted":     true,
	"digest_only":  true,
	"external_ref": true,
	"not_retained": true,
}
