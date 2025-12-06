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

// Get handles GET requests to retrieve a value by key from TiKV
func Get(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.MethodNotAllowed(w)
			return
		}

		var req types.GetRequest
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

		val, err := s.GetActiveClient().Get(ctx, []byte(req.Key))
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "TiKV Get error: "+err.Error())
			return
		}

		resp := types.GetResponse{
			Key: req.Key,
		}
		if val != nil {
			resp.Value, resp.RawValue = utils.ParseValue(val)
			resp.Found = true
		}

		utils.WriteJSON(w, http.StatusOK, resp)
	}
}
