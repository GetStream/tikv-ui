package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/GetStream/tikv-ui/pkg/types"
)

// SplitAndTrim splits a string by separator and trims whitespace from each part
func SplitAndTrim(s, sep string) []string {
	raw := strings.Split(s, sep)
	out := make([]string, 0, len(raw))
	for _, part := range raw {
		p := strings.TrimSpace(part)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func GetClusters(s string) []types.Cluster {
	parts := SplitAndTrim(s, ";")
	clusters := make([]types.Cluster, 0, len(parts))
	for _, part := range parts {
		clusters = append(clusters, GetCluster(part))
	}
	return clusters
}

func GetCluster(s string) types.Cluster {
	parts := SplitAndTrim(s, "|")
	if len(parts) < 2 {
		hex, _ := RandHex(5)
		return types.Cluster{
			Name:    fmt.Sprintf("cluster-%s", hex),
			PDAddrs: SplitAndTrim(parts[0], ","),
		}
	}
	return types.Cluster{
		Name:    parts[1],
		PDAddrs: SplitAndTrim(parts[0], ","),
	}
}

func RandHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
