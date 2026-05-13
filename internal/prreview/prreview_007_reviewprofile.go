package prreview

type ReviewProfile struct {
	SchemaVersion  string       `json:"schema_version"`
	ProfileID      string       `json:"profile_id"`
	RequiredPlanes []string     `json:"required_planes"`
	Roles          []ReviewRole `json:"roles"`
}

// ReviewRole binds one reviewer runner to one trust plane and its declared
// execution contract.
