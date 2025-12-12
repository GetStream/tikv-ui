package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/tikv/client-go/v2/rawkv"
)

// ClusterConnection represents a connection to a TiKV cluster
type ClusterConnection struct {
	Name      string
	PDAddrs   []string
	Client    *rawkv.Client
	ClusterID uint64
}

// Server holds TiKV client connections and provides HTTP handlers
type Server struct {
	mu             sync.RWMutex
	clusters       map[string]*ClusterConnection
	activeCluster  string
	defaultPDAddrs []string
}

// New creates a new Server instance with an initial connection
func New(client *rawkv.Client, pdAddrs []string, name string) *Server {
	defaultConn := &ClusterConnection{
		Name:      name,
		PDAddrs:   pdAddrs,
		Client:    client,
		ClusterID: client.ClusterID(),
	}

	return &Server{
		clusters: map[string]*ClusterConnection{
			name: defaultConn,
		},
		activeCluster:  name,
		defaultPDAddrs: pdAddrs,
	}
}

// GetActiveClient returns the currently active TiKV client
func (s *Server) GetActiveClient() *rawkv.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if conn, ok := s.clusters[s.activeCluster]; ok {
		return conn.Client
	}
	return nil
}

// AddCluster adds a new cluster connection
func (s *Server) AddCluster(ctx context.Context, name string, pdAddrs []string) (*ClusterConnection, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if cluster already exists
	if _, exists := s.clusters[name]; exists {
		return nil, fmt.Errorf("cluster '%s' already exists", name)
	}

	// Create new client
	client, err := rawkv.NewClientWithOpts(ctx, pdAddrs)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cluster: %w", err)
	}

	conn := &ClusterConnection{
		Name:      name,
		PDAddrs:   pdAddrs,
		Client:    client,
		ClusterID: client.ClusterID(),
	}

	s.clusters[name] = conn
	return conn, nil
}

// SwitchCluster switches the active cluster
func (s *Server) SwitchCluster(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clusters[name]; !exists {
		return fmt.Errorf("cluster '%s' not found", name)
	}

	s.activeCluster = name
	return nil
}

// ListClusters returns all cluster connections
func (s *Server) ListClusters() map[string]*ClusterConnection {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*ClusterConnection, len(s.clusters))
	for k, v := range s.clusters {
		result[k] = v
	}
	return result
}

// GetActiveClusterName returns the name of the active cluster
func (s *Server) GetActiveClusterName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.activeCluster
}

// Close closes all cluster connections
func (s *Server) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, conn := range s.clusters {
		if conn.Client != nil {
			conn.Client.Close()
		}
	}
}
