package packet

type IntegrationAction struct {
	Kind     string `json:"kind"`
	Actor    string `json:"actor"`
	Resolver string `json:"resolver"`
}
