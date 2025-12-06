package handlers

import (
	"net/http"

	"github.com/darkoatanasovski/tikv-ui/pkg/server"
	"github.com/darkoatanasovski/tikv-ui/pkg/utils"
)

// Health returns a simple health check response
func Health(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
