package packet

type PromptBoundaryClassification struct {
	Verdict          string   `json:"verdict"`
	RouteProofEffect string   `json:"route_proof_effect"`
	Reasons          []string `json:"reasons,omitempty"`
}
