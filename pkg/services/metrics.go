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

type Monitor struct {
	getPDAddr      func() string
	getClusterName func() string
	interval       time.Duration
	client         *http.Client
	cache          *utils.Cache
}

func NewMonitor(getPDAddr func() string, getClusterName func() string, interval time.Duration, cache *utils.Cache) *Monitor {
	return &Monitor{
		getPDAddr:      getPDAddr,
		getClusterName: getClusterName,
		interval:       interval,
		cache:          cache,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (m *Monitor) Start(ctx context.Context) {
	ticker := time.NewTicker(m.interval)

	m.pollStores(ctx)
	m.pollTiKVMetrics(ctx)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				go m.pollStores(ctx)
				go m.pollTiKVMetrics(ctx)
			}
		}
	}()
}

func (m *Monitor) pollStores(ctx context.Context) {
	pdAddr := m.getPDAddr()
	if pdAddr == "" {
		log.Printf("pd metrics: no active PD address")
		return
	}

	clusterName := m.getClusterName()
	if clusterName == "" {
		log.Printf("pd metrics: no active cluster name")
		return
	}
	storesURL := pdAddr + "/pd/api/v1/stores"
	if !strings.HasPrefix(pdAddr, "http") {
		storesURL = "http://" + pdAddr + "/pd/api/v1/stores"
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		storesURL,
		nil,
	)
	if err != nil {
		log.Printf("pd metrics: request error: %v", err)
		return
	}

	resp, err := m.client.Do(req)
	if err != nil {
		log.Printf("pd metrics: http error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("pd metrics: bad status %d", resp.StatusCode)
		return
	}

	var data types.PDStoresResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("pd metrics: decode error: %v", err)
		return
	}

	sort.Slice(data.Stores, func(i, j int) bool {
		return data.Stores[i].Store.ID < data.Stores[j].Store.ID
	})

	m.cache.Set("metrics:"+clusterName, "pd", data)
}

func (m *Monitor) pollTiKVMetrics(ctx context.Context) {
	clusterName := m.getClusterName()
	if clusterName == "" {
		log.Printf("tikv metrics: no active cluster name")
		return
	}

	// Get stores from cache
	cached, ok := m.cache.Get("metrics:"+clusterName, "pd")
	if !ok {
		log.Printf("tikv metrics: no PD stores in cache for cluster %s", clusterName)
		return
	}

	pdData, ok := cached.(types.PDStoresResponse)
	if !ok {
		log.Printf("tikv metrics: invalid PD cache data")
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
			log.Printf("tikv metrics: node %s error: %v", statusAddr, err)
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
