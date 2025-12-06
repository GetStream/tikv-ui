package types

// GetRequest represents a request to get a value by key
type GetRequest struct {
	Key string `json:"key"`
}

// PutRequest represents a request to put a key-value pair
type PutRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// DeleteRequest represents a request to delete a key
type DeleteRequest struct {
	Key string `json:"key"`
}

// ScanRequest represents a request to scan a range of keys
type ScanRequest struct {
	StartKey string `json:"start_key"`
	EndKey   string `json:"end_key"`
	Limit    int    `json:"limit"`
}

// ConnectRequest represents a request to connect to a TiKV cluster
type ConnectRequest struct {
	PDAddrs []string `json:"pd_addrs"`
	Name    string   `json:"name,omitempty"`
}
