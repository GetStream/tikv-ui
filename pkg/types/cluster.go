package types

type Cluster struct {
	Name    string   `json:"name"`
	PDAddrs []string `json:"pd_addrs"`
}
