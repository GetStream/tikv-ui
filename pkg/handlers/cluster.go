package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/darkoatanasovski/tikv-ui/pkg/server"
	"github.com/darkoatanasovski/tikv-ui/pkg/types"
	"github.com/darkoatanasovski/tikv-ui/pkg/utils"
)

// Connect handles requests to connect to a new TiKV cluster
func Connect(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.MethodNotAllowed(w)
			return
		}

		var req types.ConnectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if len(req.PDAddrs) == 0 {
			utils.WriteError(w, http.StatusBadRequest, "pd_addrs is required")
			return
		}
		if req.Name == "" {
			req.Name = "cluster-" + time.Now().Format("20060102-150405")
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		conn, err := s.AddCluster(ctx, req.Name, req.PDAddrs)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Automatically switch to the new cluster
		if err := s.SwitchCluster(req.Name); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, types.ClusterInfo{
			Name:      conn.Name,
			ClusterID: conn.ClusterID,
			PDAddrs:   conn.PDAddrs,
			Active:    true,
		})
	}
}

// ListClusters handles requests to list all connected clusters
func ListClusters(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.MethodNotAllowed(w)
			return
		}

		clusters := s.ListClusters()
		activeCluster := s.GetActiveClusterName()

		infos := make([]types.ClusterInfo, 0, len(clusters))
		for _, conn := range clusters {
			infos = append(infos, types.ClusterInfo{
				Name:      conn.Name,
				ClusterID: conn.ClusterID,
				PDAddrs:   conn.PDAddrs,
				Active:    conn.Name == activeCluster,
			})
		}

		utils.WriteJSON(w, http.StatusOK, types.ClustersResponse{
			Clusters: infos,
		})
	}
}

// SwitchCluster handles requests to switch the active cluster
func SwitchCluster(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.MethodNotAllowed(w)
			return
		}

		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if req.Name == "" {
			utils.WriteError(w, http.StatusBadRequest, "name is required")
			return
		}

		if err := s.SwitchCluster(req.Name); err != nil {
			utils.WriteError(w, http.StatusNotFound, err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"active_cluster": req.Name,
		})
	}
}
