package harnessobs

func jsonReadDenyRulePresent(path, pattern string) (bool, error) {
	// Decode only the subtree needed for this isolation proof.
	var config struct {
		Permission struct {
			Read map[string]string `json:"read"`
		} `json:"permission"`
	}
	if err := readExistingJSON(path, &config); err != nil {
		return false, err
	}
	return config.Permission.Read[pattern] == "deny", nil
}
