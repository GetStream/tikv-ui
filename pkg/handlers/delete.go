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

// Delete handles DELETE requests to remove a key from TiKV
func Delete(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.MethodNotAllowed(w)
			return
		}

		var req types.DeleteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if req.Key == "" {
			utils.WriteError(w, http.StatusBadRequest, "key is required")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := s.GetActiveClient().Delete(ctx, []byte(req.Key)); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "TiKV Delete error: "+err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
	}
}
