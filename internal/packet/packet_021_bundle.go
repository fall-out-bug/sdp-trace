package packet

type Bundle struct {
	Packet   Packet         `json:"packet"`
	Manifest BundleManifest `json:"manifest"`
}
