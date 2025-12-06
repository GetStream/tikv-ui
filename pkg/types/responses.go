package types

// GetResponse represents a response from a get operation
type GetResponse struct {
	Key      string `json:"key"`
	Value    any    `json:"value,omitempty"`
	RawValue string `json:"raw_value,omitempty"`
	Found    bool   `json:"found"`
}

// ScanItem represents a single key-value pair in a scan result
type ScanItem struct {
	Key      string `json:"key"`
	Value    any    `json:"value"`
	RawValue string `json:"raw_value"`
}

// ScanResponse represents a response from a scan operation
type ScanResponse struct {
	Items []ScanItem `json:"items"`
}

// ClusterInfo represents information about a connected cluster
type ClusterInfo struct {
	Name      string   `json:"name"`
	ClusterID uint64   `json:"cluster_id"`
	PDAddrs   []string `json:"pd_addrs"`
	Active    bool     `json:"active"`
}

// ClustersResponse represents a list of available clusters
type ClustersResponse struct {
	Clusters []ClusterInfo `json:"clusters"`
}
