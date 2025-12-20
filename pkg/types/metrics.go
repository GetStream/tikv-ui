package types

type PDStoresResponse struct {
	Count  int       `json:"count"`
	Stores []PDStore `json:"stores"`
}

type PDStore struct {
	Store  StoreMeta `json:"store"`
	Status Status    `json:"status"`
}

type StoreMeta struct {
	ID            uint64 `json:"id"`
	Address       string `json:"address"`
	State         int    `json:"state"`
	LastHeartbeat int64  `json:"last_heartbeat"`
	StateName     string `json:"state_name"`
	StatusAddress string `json:"status_address"`
}

type Status struct {
	Capacity        string `json:"capacity"`
	Available       string `json:"available"`
	UsedSize        string `json:"used_size"`
	LeaderCount     int    `json:"leader_count"`
	RegionCount     int    `json:"region_count"`
	StartTS         string `json:"start_ts"`
	LastHeartBeatTS string `json:"last_heartbeat_ts"`
	Uptime          string `json:"uptime"`
}

type LabelSet map[string]string

type TimePoint struct {
	Labels map[string]string `json:"labels,omitempty"`
	Ts     int64             `json:"ts"`
	Value  float64           `json:"value"`
}

type GaugeSeries struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Labels      []LabelSet  `json:"labels,omitempty"`
	Unit        string      `json:"unit"`
	Points      []TimePoint `json:"points"`
}

type MetricStatus struct {
	Enabled bool   `json:"enabled"`
	Label   string `json:"label"`
	Unit    string `json:"unit"`
}
