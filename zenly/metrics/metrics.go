package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	EnricherTimeMS = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "enricher_time_ms",
		Buckets: []float64{25, 50, 75, 100},
	}, []string{"enricher"})
)

func init() {
	prometheus.MustRegister(EnricherTimeMS)
}
