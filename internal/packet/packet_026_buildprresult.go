package packet

type BuildPRResult struct {
	State      string   `json:"state"`
	BundlePath string   `json:"bundle_path,omitempty"`
	PacketPath string   `json:"packet_path,omitempty"`
	ResultPath string   `json:"result_path,omitempty"`
	Errors     []string `json:"errors,omitempty"`
}
