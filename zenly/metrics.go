package zenly

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	EnricherTimeMS *prometheus.HistogramVec
}
