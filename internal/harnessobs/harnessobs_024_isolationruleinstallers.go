package harnessobs

var isolationRuleInstallers = map[string]func(string, string) error{
	"ignore_line":    ensureLineFileRule,
	"json_read_deny": ensureJSONReadDenyRule,
}
