package utils

import (
	"bufio"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/GetStream/tikv-ui/pkg/types"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
)

type ScrapeResponse struct {
	Gauges []types.GaugeSeries `json:"gauges"`
}

func ParseMetrics(r io.Reader, instance string) (ScrapeResponse, error) {
	metrics := ScrapeResponse{}
	parser := expfmt.NewTextParser(model.LegacyValidation)
	mfs, err := parser.TextToMetricFamilies(bufio.NewReader(r))
	if err != nil {
		return ScrapeResponse{}, err
	}

	for _, mf := range mfs {
		metric := MetricsMap[mf.GetName()]
		if mf.GetType() == dto.MetricType_GAUGE && metric.Enabled {
			gaugeSerie := types.GaugeSeries{
				ID:          mf.GetName(),
				Name:        metric.Label,
				Description: mf.GetHelp(),
				Unit:        metric.Unit,
			}

			for _, m := range mf.GetMetric() {
				g := m.GetGauge()
				labels := toLabelMapWithInstance(m.GetLabel(), instance)
				if len(labels) > 0 {
					gaugeSerie.Labels = append(gaugeSerie.Labels, labels)
				}
				gaugeSerie.Points = append(gaugeSerie.Points, types.TimePoint{
					Labels: labels,
					Ts:     roundToMinute(time.Now().UnixMilli()),
					Value:  g.GetValue(),
				})
			}
			metrics.Gauges = append(metrics.Gauges, gaugeSerie)
		}
	}
	return metrics, nil
}

func (dest *ScrapeResponse) Merge(fresh ScrapeResponse, window time.Duration) {
	cutoff := time.Now().Add(-window).UnixMilli()

	index := make(map[string]int)
	for i, g := range dest.Gauges {
		index[g.ID] = i
	}

	for _, freshGauge := range fresh.Gauges {
		idx, exists := index[freshGauge.ID]
		if !exists {
			dest.Gauges = append(dest.Gauges, freshGauge)
			index[freshGauge.ID] = len(dest.Gauges) - 1
			continue
		}

		destGauge := &dest.Gauges[idx]

		minuteIndex := make(map[string]int)
		for i, pt := range destGauge.Points {
			key := minuteKey(pt.Labels, pt.Ts)
			minuteIndex[key] = i
		}

		for _, newPt := range freshGauge.Points {
			key := minuteKey(newPt.Labels, newPt.Ts)
			if existingIdx, ok := minuteIndex[key]; ok {
				destGauge.Points[existingIdx].Value = newPt.Value
			} else {
				destGauge.Points = append(destGauge.Points, newPt)
				minuteIndex[key] = len(destGauge.Points) - 1
			}
		}
	}

	for i := range dest.Gauges {
		validPoints := make([]types.TimePoint, 0, len(dest.Gauges[i].Points))
		for _, p := range dest.Gauges[i].Points {
			if p.Ts >= cutoff {
				validPoints = append(validPoints, p)
			}
		}
		dest.Gauges[i].Points = validPoints
	}
}

// minuteKey creates a unique key for a point based on its labels and minute bucket
func minuteKey(labels map[string]string, ts int64) string {
	minute := ts / 60000
	return labelSignature(labels) + "|" + strconv.FormatInt(minute, 10)
}

// roundToMinute rounds a millisecond timestamp to the start of the minute
func roundToMinute(ts int64) int64 {
	return (ts / 60000) * 60000
}

func toLabelMapWithInstance(labels []*dto.LabelPair, instance string) types.LabelSet {
	out := make(types.LabelSet, len(labels)+1)
	for _, lp := range labels {
		out[lp.GetName()] = lp.GetValue()
	}
	out["instance"] = instance
	return out
}

func labelSignature(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}
	// Sort keys for stable signature
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(labels[k])
	}
	return sb.String()
}
