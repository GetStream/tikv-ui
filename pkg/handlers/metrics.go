package handlers

import (
	"net/http"

	"github.com/GetStream/tikv-ui/pkg/server"
	"github.com/GetStream/tikv-ui/pkg/utils"
)

func Metrics(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clusterName := s.GetActiveClusterName()
		cacheKey := "metrics:" + clusterName

		pd, _ := s.Cache.Get(cacheKey, "pd")
		tikv, _ := s.Cache.Get(cacheKey, "tikv")

		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"pd":   pd,
			"tikv": tikv,
		})
	}
}
