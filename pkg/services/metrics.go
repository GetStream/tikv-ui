package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/GetStream/tikv-ui/pkg/types"
	"github.com/GetStream/tikv-ui/pkg/utils"
)

// ClusterInfo holds basic cluster information for metrics polling
type ClusterInfo struct {
	Name   string
	PDAddr string
}

type Monitor struct {
	getClusters func() []ClusterInfo
	interval    time.Duration
	client      *http.Client
	cache       *utils.Cache
}

func NewMonitor(getClusters func() []ClusterInfo, interval time.Duration, cache *utils.Cache) *Monitor {
	return &Monitor{
		getClusters: getClusters,
		interval:    interval,
		cache:       cache,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (m *Monitor) Start(ctx context.Context) {
	ticker := time.NewTicker(m.interval)

	m.pollAllClusters(ctx)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				go m.pollAllClusters(ctx)
			}
		}
	}()
}

func (m *Monitor) pollAllClusters(ctx context.Context) {
	clusters := m.getClusters()
	if len(clusters) == 0 {
		log.Printf("metrics: no clusters available")
		return
	}

	for _, cluster := range clusters {
		m.pollStores(ctx, cluster)
		m.pollTiKVMetrics(ctx, cluster.Name)
	}
}

func (m *Monitor) pollStores(ctx context.Context, cluster ClusterInfo) {
	if cluster.PDAddr == "" {
		log.Printf("pd metrics [%s]: no PD address", cluster.Name)
		return
	}

	storesURL := cluster.PDAddr + "/pd/api/v1/stores"
	if !strings.HasPrefix(cluster.PDAddr, "http") {
		storesURL = "http://" + cluster.PDAddr + "/pd/api/v1/stores"
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		storesURL,
		nil,
	)
	if err != nil {
		log.Printf("pd metrics [%s]: request error: %v", cluster.Name, err)
		return
	}

	resp, err := m.client.Do(req)
	if err != nil {
		log.Printf("pd metrics [%s]: http error: %v", cluster.Name, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("pd metrics [%s]: bad status %d", cluster.Name, resp.StatusCode)
		return
	}

	var data types.PDStoresResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("pd metrics [%s]: decode error: %v", cluster.Name, err)
		return
	}

	sort.Slice(data.Stores, func(i, j int) bool {
		return data.Stores[i].Store.ID < data.Stores[j].Store.ID
	})

	m.cache.Set("metrics:"+cluster.Name, "pd", data)
}

func (m *Monitor) pollTiKVMetrics(ctx context.Context, clusterName string) {
	// Get stores from cache
	cached, ok := m.cache.Get("metrics:"+clusterName, "pd")
	if !ok {
		log.Printf("tikv metrics [%s]: no PD stores in cache", clusterName)
		return
	}

	pdData, ok := cached.(types.PDStoresResponse)
	if !ok {
		log.Printf("tikv metrics [%s]: invalid PD cache data", clusterName)
		return
	}

	// Retrieve existing metrics from cache
	var existing utils.ScrapeResponse
	if cached, ok := m.cache.Get("metrics:"+clusterName, "tikv"); ok {
		if val, ok := cached.(utils.ScrapeResponse); ok {
			existing = val
		}
	}

	// Fetch metrics from each TiKV node
	for _, store := range pdData.Stores {
		statusAddr := store.Store.StatusAddress
		if statusAddr == "" {
			continue
		}

		metricsURL := statusAddr + "/metrics"
		if !strings.HasPrefix(statusAddr, "http") {
			metricsURL = "http://" + statusAddr + "/metrics"
		}
		newMetrics, err := m.fetchNodeMetrics(ctx, metricsURL, statusAddr)
		if err != nil {
			log.Printf("tikv metrics [%s]: node %s error: %v", clusterName, statusAddr, err)
			continue
		}

		// Merge with 30 mins retention, TODO: make this configurable
		existing.Merge(newMetrics, 30*time.Minute)
	}

	sort.Slice(existing.Gauges, func(i, j int) bool {
		return existing.Gauges[i].Name < existing.Gauges[j].Name
	})

	m.cache.Set("metrics:"+clusterName, "tikv", existing)
}

func (m *Monitor) fetchNodeMetrics(ctx context.Context, url string, instance string) (utils.ScrapeResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return utils.ScrapeResponse{}, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return utils.ScrapeResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return utils.ScrapeResponse{}, err
	}

	return utils.ParseMetrics(resp.Body, instance)
}
