package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/GetStream/tikv-ui/pkg/server"
	"github.com/GetStream/tikv-ui/pkg/types"
	"github.com/GetStream/tikv-ui/pkg/utils"
	"github.com/tikv/client-go/v2/rawkv"
)

// Scan handles SCAN requests to retrieve a range of key-value pairs from TiKV
func Scan(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.MethodNotAllowed(w)
			return
		}

		var req types.ScanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if req.Limit <= 0 || req.Limit > rawkv.MaxRawKVScanLimit {
			req.Limit = 100
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		keys, values, err := s.GetActiveClient().Scan(ctx, []byte(req.StartKey), []byte(req.EndKey), req.Limit)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "TiKV Scan error: "+err.Error())
			return
		}

		items := make([]types.ScanItem, 0, len(keys))
		for i := range keys {
			parsed, raw := utils.ParseValue(values[i])
			items = append(items, types.ScanItem{
				Key:      string(keys[i]),
				Value:    parsed,
				RawValue: raw,
			})
		}

		utils.WriteJSON(w, http.StatusOK, types.ScanResponse{Items: items})
	}
}
