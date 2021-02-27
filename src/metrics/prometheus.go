package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	MongoDBDurationsSumary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "mongodb_durations_seconds",
			Help: "MongoDB operation duration in seconds",
		},
		[]string{"path"},
	)

	MongoDBDurationsHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "mongodb_durations_histogram_seconds",
			Help: "MongoDB operation duration in seconds",
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(MongoDBDurationsSumary)
	prometheus.MustRegister(MongoDBDurationsHistogram)
}
